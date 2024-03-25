package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/server"
	"github.com/iBoBoTi/ats/service"
)

type JobPostHandler interface {
	CreateJobPost(ctx *gin.Context)
	GetAllJobPosts(ctx *gin.Context)
	GetJobPostByID(ctx *gin.Context)
	GetUserJobPostByID(ctx *gin.Context)
	GetUserJobPosts(ctx *gin.Context)
}

type jobPostHandler struct {
	srv            *server.Server
	jobPostService service.JobPostService
}

func NewJobPostingHandler(srv *server.Server, jobPostService service.JobPostService) *jobPostHandler {
	return &jobPostHandler{
		srv:            srv,
		jobPostService: jobPostService,
	}
}

func (j *jobPostHandler) CreateJobPost(ctx *gin.Context) {

	user := j.srv.ContextGetUser(ctx)

	var jobPostRequest dtos.JobPost
	jobPostRequest.UserID = user.ID

	if err := ctx.ShouldBindJSON(&jobPostRequest); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	createdJobPost, err := j.jobPostService.CreateJobPost(&jobPostRequest)
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

	server.SuccessJSONResponse(ctx, http.StatusCreated, "job post created successfully", createdJobPost)

}
func (j *jobPostHandler) GetAllJobPosts(ctx *gin.Context) {
	var req dtos.PaginatedRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	jobPosts, err := j.jobPostService.GetJobPostsPaginatedByLatest(&req)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "job posts retrieved successfully", jobPosts)
}

func (j *jobPostHandler) GetJobPostByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid user_id param"))
		return
	}

	jobPost, err := j.jobPostService.GetJobPostByID(id)
	if err != nil {
		server.ErrorJSONResponse(ctx, appError.ErrStatusCode(err), errors.New("invalid user_id param"))
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "job post retrieved successfully", jobPost)
}

func (j *jobPostHandler) GetUserJobPostByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid user_id param"))
		return
	}

	user := j.srv.ContextGetUser(ctx)

	jobPost, err := j.jobPostService.GetUserJobPostByID(id, user.ID)
	if err != nil {
		server.ErrorJSONResponse(ctx, appError.ErrStatusCode(err), errors.New("invalid user_id param"))
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "job post retrieved successfully", jobPost)
}

func (j *jobPostHandler) GetUserJobPosts(ctx *gin.Context) {
	var req dtos.PaginatedRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user := j.srv.ContextGetUser(ctx)
	jobPosts, err := j.jobPostService.GetUserJobPostsPaginatedByLatest(user.ID, &req)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "job posts retrieved successfully", jobPosts)
}
