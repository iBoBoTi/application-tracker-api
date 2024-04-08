package mappers_test

import (
	"testing"

	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUserDtoMapToUserModel(t *testing.T) {
	userDto := &dtos.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "email@email.com",
		Password:  "password",
		UserType:  "admin",
	}

	userModel := mappers.UserDtoMapToUserModel(userDto)

	assert.Equal(t, userDto.FirstName, userModel.FirstName)
	assert.Equal(t, userDto.LastName, userModel.LastName)
	assert.Equal(t, userDto.Email, userModel.Email)
}

func TestUserModelMapToUserDto(t *testing.T) {
	userModel := &models.User{
		FirstName: "FirstName",
		LastName:  "LastName",
		Email:     "email@email.com",
		PasswordHash:  "password",
		UserType:  "admin",
	}

	userDto := mappers.UserModelMapToUserDto(userModel)

	assert.Equal(t, userModel.FirstName, userDto.FirstName)
	assert.Equal(t, userModel.LastName, userDto.LastName)
	assert.Equal(t, userModel.Email, userDto.Email)
}
