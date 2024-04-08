package dtos_test

import (
	"testing"

	"github.com/iBoBoTi/ats/internal/dtos"
	"github.com/iBoBoTi/ats/internal/validator"
	"github.com/stretchr/testify/assert"
)

func TestJobPostValidate(t *testing.T){
	
	testCases := []struct{
		name string
		jobPost *dtos.JobPost
		isValid bool
	}{
		{
			name: "valid job post",
			jobPost: &dtos.JobPost{
				Title:                "title",
				CompanyName:          "company name",
				Description:          "description",
			},
			isValid: true,
		},
		{
			name: "invalid job post",
			jobPost: &dtos.JobPost{
				Title:                "title",
				Description:          "description",
			},
			isValid: false,
		},
	}

	validator := validator.NewValidator()

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			isValid := tc.jobPost.Validate(validator)

			assert.Equal(t, tc.isValid, isValid)
		})
	}

}