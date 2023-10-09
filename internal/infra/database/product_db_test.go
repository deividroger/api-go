package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/deividroger/api-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateNewProdu(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, err := entity.NewProduct("Product 1", 10.0)

	assert.NoError(t, err)
	assert.NotEmpty(t, product)
	assert.NotEmpty(t, product.ID)

}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	for i := 1; i < 24; i++ {
		product, _ := entity.NewProduct(fmt.Sprintf("Product %d", i), float64(rand.Float64())*100)
		assert.NoError(t, err)
		db.Create(product)

	}

	prductDb := NewProduct(db)

	products, err := prductDb.FindAll(1, 10, "asc")

	assert.NoError(t, err)
	assert.Len(t, products, 10)

	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 10", products[9].Name)

	products, err = prductDb.FindAll(2, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 10)
	assert.Equal(t, "Product 11", products[0].Name)
	assert.Equal(t, "Product 20", products[9].Name)

	products, err = prductDb.FindAll(3, 10, "asc")
	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Product 21", products[0].Name)
	assert.Equal(t, "Product 23", products[2].Name)

}
func TestFindProductByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)
	db.Create(product)

	prductDb := NewProduct(db)

	productFound, err := prductDb.FindById(product.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, product.ID, productFound.ID)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
}

func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)
	db.Create(product)

	productDb := NewProduct(db)
	product.Name = "Product 2"

	err = productDb.Update(product)
	assert.NoError(t, err)

	productFound, err := productDb.FindById(product.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, "Product 2", productFound.Name)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Product{})

	product, _ := entity.NewProduct("Product 1", 10.0)
	assert.NoError(t, err)
	db.Create(product)

	productDb := NewProduct(db)

	err = productDb.Delete(product.ID.String())

	assert.NoError(t, err)

	_, err = productDb.FindById(product.ID.String())

	assert.Error(t, err)

}
