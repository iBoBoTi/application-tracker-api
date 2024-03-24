package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/server"
	"github.com/iBoBoTi/ats/service"
	"github.com/iBoBoTi/ats/internal/validator"
)

type AuthHandler interface {
	SignUp(ctx *gin.Context)
}

type authHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *authHandler {
	return &authHandler{
		userService: userService,
	}
}

func (a *authHandler) SignUp(ctx *gin.Context) {
	var userRequest dtos.User

	if err := ctx.ShouldBindJSON(&userRequest); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	createdUser, err := a.userService.CreateUser(&userRequest)
	if err != nil {
		var e *validator.ValidationError
		switch {
		case errors.As(err, &e):
			server.SendValidationError(ctx, e)
		default:
			server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusCreated, "user created successfully", createdUser)
}
