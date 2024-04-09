package dtos_test

import (
	"testing"

	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {

	testCases := []struct {
		name    string
		user    *dtos.User
		isValid bool
	}{
		{
			name: "valid user details",
			user: &dtos.User{
				FirstName:       "FirstName",
				LastName:        "LastName",
				Email:           "email@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			isValid: true,
		},
		{
			name: "invalid user details caused by empty first name",
			user: &dtos.User{
				LastName:        "LastName",
				Email:           "email@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			isValid: false,
		},
		{
			name: "invalid user details caused by invalid email",
			user: &dtos.User{
				FirstName:       "FirstName",
				LastName:        "LastName",
				Email:           "email.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			isValid: false,
		},
		{
			name: "invalid user details caused by password not equal to confirm password",
			user: &dtos.User{
				FirstName:       "FirstName",
				LastName:        "LastName",
				Email:           "email@email.com",
				Password:        "password",
				ConfirmPassword: "confirm",
			},
			isValid: false,
		},
	}

	validator := validator.NewValidator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.user.Validate(validator)

			assert.Equal(t, tc.isValid, isValid)
		})
	}
}

func TestLoginValidate(t *testing.T) {
	testCases := []struct {
		name     string
		loginReq *dtos.LoginRequest
		isValid  bool
	}{
		{
			name: "valid login request",
			loginReq: &dtos.LoginRequest{
				Email:    "email@email.com",
				Password: "password",
			},
			isValid: true,
		},
		{
			name: "invalid login request caused by missing email",
			loginReq: &dtos.LoginRequest{
				Password: "password",
			},
			isValid: false,
		},
		{
			name: "valid login request caused by missing password",
			loginReq: &dtos.LoginRequest{
				Email: "email@email.com",
			},
			isValid: false,
		},
	}

	validator := validator.NewValidator()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.loginReq.Validate(validator)

			assert.Equal(t, tc.isValid, isValid)
		})
	}
}

func TestIsEmail(t *testing.T) {
	testCases := []struct {
		name         string
		email        string
		isValidEmail bool
	}{
		{
			name:         "valid email",
			email:        "email@email.com",
			isValidEmail: true,
		},
		{
			name:         "invalid email",
			email:        "randomemail",
			isValidEmail: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := dtos.IsEmail(tc.email)

			assert.Equal(t, tc.isValidEmail, isValid)
		})
	}
}
