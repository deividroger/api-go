package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {

	user, err := NewUser("John Doe", "j@j.com", "123456")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "John Doe", user.Name)

	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "j@j.com", user.Email)

}

func TestUser_ValidatePassword(t *testing.T) {

	user, err := NewUser("John Doe", "j@j.com", "123456")

	assert.Nil(t, err)
	assert.True(t, user.ComparePassword("123456"))
	assert.False(t, user.ComparePassword("1234567"))

	assert.NotEqual(t, "123456", user.Password)
}
