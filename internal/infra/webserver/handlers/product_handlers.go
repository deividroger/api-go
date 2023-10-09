package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/deividroger/api-go/internal/dto"
	"github.com/deividroger/api-go/internal/entity"
	"github.com/deividroger/api-go/internal/infra/database"
	entityPkg "github.com/deividroger/api-go/pkg/entity"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDb database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDb: db,
	}
}

// Create user godoc
// @Summary 	List Products
// @Description List Products
// @Tags 		products
// @Accept  	json
// @Produce  	json
// @Param 		page	 	query 		string 	false 	"page_number"
// @Param 		limit	 	query 		string 	false 	"limit"
// @Success 	200			{array}		[]entity.Product
// @Failure 	404 		{object} 	dto.Error
// @Failure 	500 		{object} 	dto.Error
// @Router 		/products [get]
// @security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		pageInt = 0
	}
	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		limitInt = 0
	}

	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDb.FindAll(pageInt, limitInt, sort)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)

}

// Create user godoc
// @Summary 	Create Product
// @Description Create Product
// @Tags 		products
// @Accept  	json
// @Produce  	json
// @Param 		request 	body 		dto.CreateProductInput 	true 	"product request"
// @Success 	201
// @Failure 	500 		{object} 	dto.Error
// @Router 		/products [post]
// @security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	err = h.ProductDb.Create(p)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Create user godoc
// @Summary 	Get Product
// @Description Get Product
// @Tags 		products
// @Accept  	json
// @Produce  	json
// @Param 		id		 	path 		string 	true 	"product ID Format(UUID)"
// @Success 	200			{object}	entity.Product
// @Failure 	404 		{object} 	dto.Error
// @Failure 	500 		{object} 	dto.Error
// @Router 		/products/{id} [get]
// @security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: "id is required",
		})
		return
	}

	product, err := h.ProductDb.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Create user godoc
// @Summary 	Update a product
// @Description Update a product
// @Tags 		products
// @Accept  	json
// @Produce  	json
// @Param 		id		 	path 		string 					true 			"product ID Format(UUID)"
// @Param 		request 	body 		dto.CreateProductInput 	true 			"product request"
// @Success 	200
// @Failure 	404 		{object} 	dto.Error
// @Failure 	500 		{object} 	dto.Error
// @Router 		/products/{id} [put]
// @security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is required"))
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	product.ID, err = entityPkg.ParseID(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	_, err = h.ProductDb.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	err = h.ProductDb.Update(&product)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Create user godoc
// @Summary 	Delete Product
// @Description Delete Product
// @Tags 		products
// @Accept  	json
// @Produce  	json
// @Param 		id		 	path 		string 	true 	"product ID Format(UUID)"
// @Success 	204
// @Failure 	404 		{object} 	dto.Error
// @Failure 	500 		{object} 	dto.Error
// @Router 		/products/{id} [delete]
// @security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(dto.Error{
			Message: "id is required",
		})
		return
	}

	_, err := h.ProductDb.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	err = h.ProductDb.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}
