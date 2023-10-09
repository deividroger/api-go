package database

import (
	"github.com/deividroger/api-go/internal/entity"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{DB: db}
}

func (u *User) Find(id string) (*entity.User, error) {
	return nil, nil
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (U *User) FindByEmail(email string) (*entity.User, error) {

	var user entity.User
	if err := U.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil

}
