package mappers_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/mappers"
	"github.com/iBoBoTi/ats/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestJobPostDtoMapToJobPostModel(t *testing.T) {
	jobPostDto := &dtos.JobPost{
		Title:                "title",
		CompanyName:          "company name",
		Description:          "description",
		MinYearsOfExperience: 1,
		MaxYearsOfExperience: 3,
		UserID:               uuid.New(),
		Keywords: []string{
			"keyword1",
			"keyword2",
			"keyword3",
		},
	}

	jobPostModel := mappers.JobPostDtoMapToJobPostModel(jobPostDto)

	assert.Equal(t, jobPostDto.Title, jobPostModel.Title)
	assert.Equal(t, jobPostDto.UserID, jobPostModel.UserID)
	assert.Equal(t, strings.Join(jobPostDto.Keywords, ","), jobPostModel.Keywords)
}

func TestJobPostModelMapToJobPostDto(t *testing.T) {

	jobPostModel := &models.JobPost{
		Title:                "title",
		CompanyName:          "company name",
		Description:          "description",
		MinYearsOfExperience: 1,
		MaxYearsOfExperience: 3,
		Keywords:             "keyword1,keyword2,keyword3",
	}

	jobPostDto := mappers.JobPostModelMapToJobPostDto(jobPostModel)

	assert.Equal(t, jobPostModel.Title, jobPostDto.Title)
	assert.Equal(t, strings.Split(jobPostModel.Keywords, ","), jobPostDto.Keywords)
}
