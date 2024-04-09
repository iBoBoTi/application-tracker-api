package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/repository"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"gorm.io/gorm"
)

type ApplicantService interface {
	CreateApplicant(applicant *dtos.Applicant, file multipart.File) error
	GetAllApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error)
	GetQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error)
	GetUnQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error)
	GetApplicantByJobID(id, jobID uuid.UUID) (*dtos.Applicant, error)
}

type applicantService struct {
	applicantRepository repository.ApplicantRepository
	jobPostRepository   repository.JobPostRepository
}

func NewApplicantService(applicantRepository repository.ApplicantRepository, jobPostRepository repository.JobPostRepository) *applicantService {
	return &applicantService{
		applicantRepository: applicantRepository,
		jobPostRepository:   jobPostRepository,
	}
}

func (a *applicantService) CreateApplicant(applicant *dtos.Applicant, file multipart.File) error {
	v := validator.NewValidator()

	jobPost, err := a.jobPostRepository.GetJobPostByID(applicant.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			v.Check(false, "job_id", "job post with id doesn't exist")
			return validator.NewValidationError("validation failed", v.Errors)
		}
		return err
	}

	if !applicant.Validate(v) {
		return validator.NewValidationError("validation failed", v.Errors)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return appError.ErrInternalServer
	}

	text, err := ExtractTextFromPDF(file)
	if err != nil {
		log.Println(err)
		return appError.ErrInternalServer
	}

	applicantModel := mappers.ApplicantDtoMapToApplicantModel(applicant, file)

	jobPostKeywords := strings.Split(jobPost.Keywords, ",")
	minimumKeywordsPermitted := len(jobPostKeywords) / 2

	if applicant.YearsOfExperience >= int(jobPost.MinYearsOfExperience) {
		numOfFoundKeywords := 0
		for _, keywords := range jobPostKeywords {
			if strings.Contains(strings.ToLower(text), strings.ToLower(keywords)) {
				numOfFoundKeywords++
			}
		}
		if numOfFoundKeywords >= minimumKeywordsPermitted {
			applicantModel.IsQualified = true
		}
	}

	applicantModel.Resume = content

	return a.applicantRepository.CreateApplicant(applicantModel)
}

func (a *applicantService) GetAllApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error) {

	paginateReq.Normalize()

	applicantsModels, err := a.applicantRepository.GetAllApplicantsByJobIDPaginated(jobID, int(paginateReq.Limit), int(paginateReq.Page))
	if err != nil {
		return nil, err
	}

	applicantDtos := make([]dtos.Applicant, 0)
	if len(applicantsModels) > 0 {
		for _, applicant := range applicantsModels {
			applicantDto := mappers.ApplicantModelMapToApplicantDto(&applicant)
			applicantDtos = append(applicantDtos, *applicantDto)
		}
	}

	return applicantDtos, nil
}

func (a *applicantService) GetQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error) {

	paginateReq.Normalize()

	applicantsModels, err := a.applicantRepository.GetQualifiedApplicantsByJobIDPaginated(jobID, int(paginateReq.Limit), int(paginateReq.Page))
	if err != nil {
		return nil, err
	}

	applicantDtos := make([]dtos.Applicant, 0)
	if len(applicantsModels) > 0 {
		for _, applicant := range applicantsModels {
			applicantDto := mappers.ApplicantModelMapToApplicantDto(&applicant)
			applicantDtos = append(applicantDtos, *applicantDto)
		}
	}

	return applicantDtos, nil
}

func (a *applicantService) GetUnQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.Applicant, error) {

	paginateReq.Normalize()

	applicantsModels, err := a.applicantRepository.GetUnQualifiedApplicantsByJobIDPaginated(jobID, int(paginateReq.Limit), int(paginateReq.Page))
	if err != nil {
		return nil, err
	}

	applicantDtos := make([]dtos.Applicant, 0)
	if len(applicantsModels) > 0 {
		for _, applicant := range applicantsModels {
			applicantDto := mappers.ApplicantModelMapToApplicantDto(&applicant)
			applicantDtos = append(applicantDtos, *applicantDto)
		}
	}

	return applicantDtos, nil
}

func (a *applicantService) GetApplicantByJobID(id, jobID uuid.UUID) (*dtos.Applicant, error) {

	applicantModel, err := a.applicantRepository.GetApplicantByJobID(id, jobID)
	if err != nil {
		return nil, err
	}
	log.Println("GOT HERE TOO")

	return mappers.ApplicantModelMapToApplicantDto(applicantModel), nil
}

func ExtractTextFromPDF(file multipart.File) (string, error) {
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		return "", fmt.Errorf("error reading PDF file: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("error getting number of pages in PDF file: %v", err)
	}

	var fullText string

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", fmt.Errorf("error getting number of page in PDF file: %v", err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return "", fmt.Errorf("error extracting pages in PDF file: %v", err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return "", fmt.Errorf("error extracting text from page in PDF file: %v", err)
		}

		fullText = fullText + "\n" + text
	}
	return fullText, nil
}
