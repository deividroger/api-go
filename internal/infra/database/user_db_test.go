package database

import (
	"testing"

	"github.com/deividroger/api-go/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUse(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "email@email.com", "123456")

	userDB := NewUser(db)

	err = userDB.Create(user)

	assert.Nil(t, err)

	var userFound entity.User

	err = db.First(&userFound, "id = ?", user.ID.String()).Error

	assert.Nil(t, err)
	assert.Equal(t, user.ID.String(), userFound.ID.String())
	assert.Equal(t, user.Name, userFound.Name)
	assert.Equal(t, user.Email, userFound.Email)
	assert.NotNil(t, userFound.Password)

}

func TestFindByEmail(T *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		T.Error(err)
	}
	db.AutoMigrate(&entity.User{})

	user, _ := entity.NewUser("John Doe", "j@j.com", "123456")

	userDB := NewUser(db)
	err = userDB.Create(user)

	assert.Nil(T, err)

	userFound, err := userDB.FindByEmail(user.Email)

	assert.Nil(T, err)
	assert.Equal(T, user.ID, userFound.ID)
	assert.Equal(T, user.Name, userFound.Name)
	assert.Equal(T, user.Email, userFound.Email)
	assert.NotNil(T, userFound.Password)

}
