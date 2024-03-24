package service

import (
	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/repository"
)

type UserService interface {
	CreateUser(userDto *dtos.User) (*dtos.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *userService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) CreateUser(userDto *dtos.User) (*dtos.User, error) {

	v := validator.NewValidator()

	emailExist, err := u.userRepository.EmailExist(userDto.Email)
	if err != nil {
		return nil, err
	}
	v.Check(emailExist, "email", "email already exist")

	if !userDto.Validate(v) {
		return nil, validator.NewValidationError("validation failed", v.Errors)
	}

	if err := userDto.HashPassword(); err != nil {
		return nil, err
	}

	user := mappers.UserDtoMapToUserModel(userDto)

	if err := u.userRepository.CreateUser(user); err != nil {
		return nil, err
	}

	return mappers.UserModelMapToUserDto(user), nil
}
