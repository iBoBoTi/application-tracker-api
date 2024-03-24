package service

import (
	"fmt"
	"time"

	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/security"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/repository"
)

var AccessTokenDuration = 30 * time.Minute

type UserService interface {
	CreateUser(userDto *dtos.User) (*dtos.User, error)
	Login(loginRequest *dtos.LoginRequest) (*security.AuthPayload, error)
}

type userService struct {
	tokenMaker     security.Maker
	userRepository repository.UserRepository
}

func NewUserService(tokenMaker security.Maker, userRepository repository.UserRepository) *userService {
	return &userService{
		tokenMaker:     tokenMaker,
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

func (u *userService) Login(loginRequest *dtos.LoginRequest) (*security.AuthPayload, error) {
	v := validator.NewValidator()
	if !loginRequest.Validate(v) {
		return nil, validator.NewValidationError("validation failed", v.Errors)
	}

	emailExist, err := u.userRepository.EmailExist(loginRequest.Email)
	if err != nil {
		return nil, err
	}

	if emailExist {
		return nil, fmt.Errorf("%v: email not found", appError.ErrNotFound)
	}

	foundUser, err := u.userRepository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		return nil, appError.ErrInternalServer
	}

	if err := loginRequest.CheckPassword(foundUser.PasswordHash); err != nil {
		return nil, appError.ErrInvalidCredential
	}

	var res security.AuthPayload
	res.Data = make(map[string]any)

	if err := u.tokenMaker.GenerateAuthAccessToken(foundUser, &res, AccessTokenDuration); err != nil {
		return nil, err
	}

	// validate if user exist then also validate the password
	return &res, nil
}
