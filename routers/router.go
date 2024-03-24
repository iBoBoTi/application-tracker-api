package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/handlers"
	"github.com/iBoBoTi/ats/repository"
	"github.com/iBoBoTi/ats/server"
	"github.com/iBoBoTi/ats/service"
)

const (
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
	EnvironmentTesting     = "testing"
)

// SetupRouter registers all the HTTP routes in the system
// if you want to move out some routes, you can accept *gin.Engine as an argument
func SetupRouter(srv *server.Server) {
	if srv.GetConfig().Environment == EnvironmentDevelopment ||
		srv.GetConfig().Environment == EnvironmentTesting {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(server.CustomLogger(srv.Logger), gin.Recovery()).Use(server.CORS())

	v1 := router.Group("/api/v1")

	// v1.Use(srv.ApplyAuthentication())

	v1.GET("/health-check", handlers.NewHealthController(srv).HealthCheck)

	authHandler := handlers.NewAuthHandler(service.NewUserService(srv.TokenMaker, repository.NewUserRepository(srv.DB.GormDB)))
	v1.POST("/users/login", authHandler.Login)
	v1.POST("/users/signup", authHandler.SignUp)

	// authCrt := auth.NewAuthController(srv)
	// v1.POST("/users/login", authCrt.LoginUser)
	// v1.POST("/users/refresh-token", authCrt.UserRefreshToken)
	// v1.POST("/users/create", authCrt.CreateAccount)
	// v1.POST("/users/verify-by-token", authCrt.VerifyUserAccountByToken)
	// v1.POST("/users/forgot-password", authCrt.ForgotPassword)
	// v1.POST("/users/tokens/check", authCrt.ResetPasswordTokenCheck)
	// v1.POST("/users/reset-password", authCrt.ResetPassword)

	// v1.GET("/currencies", currency.NewCurrencyController(srv).GetActiveCurrencies)
	// v1.POST("/resend-token", authCrt.ResendToken)

	// v1.GET("/quotes", exchangerate.NewExchangeRateController(srv).GetQuotesForPublic)

	// internal := v1.Group("/internal")
	// internal.POST("/user/impersonate", authCrt.ImpersonateUser)

	// // register user routes/endpoint
	// // For authenticated users only
	// authRoute := v1.Group("/")
	// authRoute.Use(srv.CheckIfBanned(), srv.LogUserAction())
	// registerUserRoutes(authRoute, srv)

	srv.Router = router

}
