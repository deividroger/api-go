package database

import "github.com/deividroger/api-go/internal/entity"

type UserInterFace interface {
	Create(user *entity.User) error
	Find(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}

type ProductInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindById(id string) (entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}
