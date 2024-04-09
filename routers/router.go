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

	v1.GET("/health-check", handlers.NewHealthController(srv).HealthCheck)

	authHandler := handlers.NewAuthHandler(service.NewUserService(srv.TokenMaker, repository.NewUserRepository(srv.DB.GormDB)))
	v1.POST("/users/login", authHandler.Login)
	v1.POST("/users/signup", authHandler.SignUp)

	jobPostHandler := handlers.NewJobPostingHandler(srv, service.NewJobPostingService(repository.NewJobPostRepository(srv.DB.GormDB)))

	v1.GET("/job-posts", jobPostHandler.GetAllJobPosts)
	v1.GET("/job-posts/:id", jobPostHandler.GetJobPostByID)

	applicantHandler := handlers.NewApplicantHandler(service.NewApplicantService(repository.NewApplicantRepository(srv.DB.GormDB), repository.NewJobPostRepository(srv.DB.GormDB)))
	v1.POST("/job-posts/:id/applicants", applicantHandler.CreateApplicant)

	v1.Use(srv.ApplyAuthentication())
	v1.POST("/users/job-posts", jobPostHandler.CreateJobPost)
	v1.GET("/users/job-posts", jobPostHandler.GetUserJobPosts)
	v1.GET("/users/job-posts/:id", jobPostHandler.GetUserJobPostByID)

	srv.Router = router

}
