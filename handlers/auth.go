package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/server"
	"github.com/iBoBoTi/ats/service"
)

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
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

	server.SuccessJSONResponse(ctx, http.StatusCreated, "user signup successfully", createdUser)
}

func (a *authHandler) Login(ctx *gin.Context) {
	var loginRequest dtos.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := a.userService.Login(&loginRequest)
	if err != nil {
		var e *validator.ValidationError
		switch {
		case errors.As(err, &e):
			server.SendValidationError(ctx, e)
		default:
			server.ErrorJSONResponse(ctx, appError.ErrStatusCode(err), err)
		}
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "login successful", res.Data)
}
