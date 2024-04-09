package mappers

import (
	"mime/multipart"

	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/models"
)

func ApplicantDtoMapToApplicantModel(applicantDto *dtos.Applicant, file multipart.File) *models.Applicant {
	return &models.Applicant{
		FirstName:         applicantDto.FirstName,
		LastName:          applicantDto.LastName,
		Email:             applicantDto.Email,
		YearsOfExperience: applicantDto.YearsOfExperience,
		JobID:             applicantDto.JobID,
		ResumeFileName:    applicantDto.ResumeFileDetails.NewFileName,
	}
}
