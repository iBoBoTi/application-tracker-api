package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/dtos"
	appError "github.com/iBoBoTi/ats/internal/errors"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/iBoBoTi/ats/repository"
	"gorm.io/gorm"
)

type JobPostService interface {
	CreateJobPost(jobPost *dtos.JobPost) (*dtos.JobPost, error)
	GetJobPostsPaginatedByLatest(paginateReq *dtos.PaginatedRequest) ([]dtos.JobPost, error)
	GetUserJobPostsPaginatedByLatest(userID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.JobPost, error)
	GetUserJobPostByID(id, userID uuid.UUID) (*dtos.JobPost, error)
	GetJobPostByID(id uuid.UUID) (*dtos.JobPost, error)
}

type jobPostService struct {
	jobPostRepository repository.JobPostRepository
}

func NewJobPostingService(jobPostRepository repository.JobPostRepository) *jobPostService {
	return &jobPostService{
		jobPostRepository: jobPostRepository,
	}
}

func (j *jobPostService) CreateJobPost(jobPostDto *dtos.JobPost) (*dtos.JobPost, error) {
	v := validator.NewValidator()

	if !jobPostDto.Validate(v) {
		return nil, validator.NewValidationError("validation failed", v.Errors)
	}

	jobPostModel := mappers.JobPostDtoMapToJobPostModel(jobPostDto)

	if err := j.jobPostRepository.CreateJobPost(jobPostModel); err != nil {
		return nil, err
	}

	return mappers.JobPostModelMapToJobPostDto(jobPostModel), nil
}

func (j *jobPostService) GetJobPostsPaginatedByLatest(paginateReq *dtos.PaginatedRequest) ([]dtos.JobPost, error) {
	paginateReq.Normalize()

	jobPostModels, err := j.jobPostRepository.GetJobPostsPaginatedByLatest(int(paginateReq.Limit), int(paginateReq.Page))
	if err != nil {
		return nil, err
	}

	jobPostsDto := make([]dtos.JobPost, 0)
	if len(jobPostModels) > 0 {
		for _, jobPost := range jobPostModels {
			jobPostDto := mappers.JobPostModelMapToJobPostDto(&jobPost)
			jobPostsDto = append(jobPostsDto, *jobPostDto)
		}
	}

	return jobPostsDto, nil
}

func (j *jobPostService) GetUserJobPostsPaginatedByLatest(userID uuid.UUID, paginateReq *dtos.PaginatedRequest) ([]dtos.JobPost, error) {
	paginateReq.Normalize()

	jobPostModels, err := j.jobPostRepository.GetJobPostsByOwnerPaginatedByLatest(userID, paginateReq.Limit, paginateReq.Page)
	if err != nil {
		return nil, err
	}

	jobPostsDto := make([]dtos.JobPost, 0)
	if len(jobPostModels) > 0 {
		for _, jobPost := range jobPostModels {
			jobPostDto := mappers.JobPostModelMapToJobPostDto(&jobPost)
			jobPostsDto = append(jobPostsDto, *jobPostDto)
		}
	}

	return jobPostsDto, nil
}

func (j *jobPostService) GetUserJobPostByID(id, userID uuid.UUID) (*dtos.JobPost, error) {

	jobPostModel, err := j.jobPostRepository.GetUserJobPostByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%v: email not found", appError.ErrNotFound)
		}
		return nil, err
	}

	return mappers.JobPostModelMapToJobPostDto(jobPostModel), nil
}

func (j *jobPostService) GetJobPostByID(id uuid.UUID) (*dtos.JobPost, error) {
	jobPostModel, err := j.jobPostRepository.GetJobPostByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("%v: email not found", appError.ErrNotFound)
			}
		}
		return nil, err
	}

	return mappers.JobPostModelMapToJobPostDto(jobPostModel), nil
}
