package mappers

import (
	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/models"
)

const (
	AdminUserType string = "admin"
)

func UserDtoMapToUserModel(userDto *dtos.User) *models.User {
	return &models.User{
		FirstName:    userDto.FirstName,
		LastName:     userDto.LastName,
		Email:        userDto.Email,
		PasswordHash: userDto.Password,
		UserType:     AdminUserType,
	}
}

func UserModelMapToUserDto(user *models.User) *dtos.User {
	return &dtos.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
