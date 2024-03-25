package mappers

import (
	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/models"
	"strings"
)

func JobPostDtoMapToJobPostModel(jobPostDto *dtos.JobPost) *models.JobPost {
	return &models.JobPost{
		Title:                jobPostDto.Title,
		CompanyName:          jobPostDto.CompanyName,
		Description:          jobPostDto.Description,
		MinYearsOfExperience: jobPostDto.MinYearsOfExperience,
		MaxYearsOfExperience: jobPostDto.MaxYearsOfExperience,
		UserID:               jobPostDto.UserID,
		Keywords:             strings.Join(jobPostDto.Keywords, ","),
	}
}

func JobPostModelMapToJobPostDto(jobPostModel *models.JobPost) *dtos.JobPost {
	return &dtos.JobPost{
		ID:                   jobPostModel.ID,
		Title:                jobPostModel.Title,
		Description:          jobPostModel.Description,
		CompanyName:          jobPostModel.CompanyName,
		MinYearsOfExperience: jobPostModel.MinYearsOfExperience,
		MaxYearsOfExperience: jobPostModel.MaxYearsOfExperience,
		Keywords:             strings.Split(jobPostModel.Keywords, ","),
		CreatedAt:            jobPostModel.CreatedAt,
		UpdatedAt:            jobPostModel.UpdatedAt,
	}
}
