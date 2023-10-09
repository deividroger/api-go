package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deividroger/api-go/internal/dto"
	"github.com/deividroger/api-go/internal/entity"
	"github.com/deividroger/api-go/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDb database.UserInterFace
}

func NewUserHandler(db database.UserInterFace) *UserHandler {
	return &UserHandler{
		UserDb: db,
	}
}

// Create user godoc
// @Summary 	Get a user JWT
// @Description Get a user JWT
// @Tags 		users
// @Accept  	json
// @Produce  	json
// @Param 		request 	body 		dto.GetJwtInput 	true 	"user credentials"
// @Success 	200			{object}	dto.GetJwtOutput
// @Failure		404			{object}	dto.Error
// @Failure 	500 		{object} 	dto.Error
// @Router 		/users/generateToken [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	expiresIn := r.Context().Value("expiresIn").(int)
	var user dto.GetJwtInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	u, err := h.UserDb.FindByEmail(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	if !u.ComparePassword(user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub":  u.ID.String(),
		"exp":  time.Now().Add(time.Second * time.Duration(expiresIn)).Unix(),
		"name": u.Name,
	})

	accessToken := dto.GetJwtOutput{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)

}

// Create user godoc
// @Summary 	Create user
// @Description Create user
// @Tags 		users
// @Accept  	json
// @Produce  	json
// @Param 		request 	body 		dto.CreateUserInput 	true 	"user request"
// @Success 	201
// @Failure 	500 		{object} 	dto.Error
// @Router 		/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user dto.CreateUserInput

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}

	err = h.UserDb.Create(u)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.Error{
			Message: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)

}
