package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/server"
	"github.com/iBoBoTi/ats/service"
)

type ApplicantHandler interface {
	CreateApplicant(ctx *gin.Context)
}

type applicantHandler struct {
	applicantService service.ApplicantService
}

func NewApplicantHandler(applicantService service.ApplicantService) *applicantHandler {
	return &applicantHandler{
		applicantService: applicantService,
	}
}

func (a *applicantHandler) CreateApplicant(ctx *gin.Context) {
	applicantReq, err := dtos.SetupApplicantFormFileValues(ctx)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("resume")
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	jobID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid user_id param"))
		return
	}
	applicantReq.JobID = jobID
	applicantReq.ResumeFileDetails.ContentType = fileHeader.Header.Get("Content-Type")
	applicantReq.ResumeFileDetails.Size = fileHeader.Size
	applicantReq.ResumeFileDetails.FileName = fileHeader.Filename

	if err := a.applicantService.CreateApplicant(applicantReq, file); err != nil {
		var e *validator.ValidationError
		switch {
		case errors.As(err, &e):
			server.SendValidationError(ctx, e)
		default:
			server.ErrorJSONResponse(ctx, appError.ErrStatusCode(err), err)
		}
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusCreated, "applicant created successfully", nil)
}

func (a *applicantHandler) GetAllApplicantsByJobID(ctx *gin.Context) {
	jobID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid job id param"))
		return
	}

	var req dtos.PaginatedRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	applicants, err := a.applicantService.GetAllApplicantsByJobIDPaginated(jobID, &req)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "applicants retrieved successfully", applicants)
}

func (a *applicantHandler) GetQualifiedApplicantsByJobIDPaginated(ctx *gin.Context) {
	jobID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid job_id param"))
		return
	}

	var req dtos.PaginatedRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	applicants, err := a.applicantService.GetQualifiedApplicantsByJobIDPaginated(jobID, &req)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "applicants retrieved successfully", applicants)
}

func (a *applicantHandler) GetUnQualifiedApplicantsByJobIDPaginated(ctx *gin.Context) {
	jobID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid job id param"))
		return
	}

	var req dtos.PaginatedRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.ErrorJSONResponse(ctx, http.StatusBadRequest, err)
		return
	}

	applicants, err := a.applicantService.GetUnQualifiedApplicantsByJobIDPaginated(jobID, &req)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "applicants retrieved successfully", applicants)
}

func (a *applicantHandler) GetApplicantByJobID(ctx *gin.Context) {
	jobID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid job id param"))
		return
	}

	id, err := uuid.Parse(ctx.Param("applicant_id"))
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusUnprocessableEntity, errors.New("invalid applicant id param"))
		return
	}

	log.Println("HEEEEEEEEEEEEEERRRRRR")

	applicant, err := a.applicantService.GetApplicantByJobID(id, jobID)
	if err != nil {
		server.ErrorJSONResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	server.SuccessJSONResponse(ctx, http.StatusOK, "applicant retrieved successfully", applicant)
}
