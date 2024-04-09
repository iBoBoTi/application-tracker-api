package dtos

import (
	"errors"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/iBoBoTi/ats/internal/validator"
)

const MaxFileUploadSize = 1024 * 1024 * 2
const fileExtension = ".pdf"

type Applicant struct {
	ID                uuid.UUID `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	YearsOfExperience int       `json:"years_of_experience"`
	Email             string    `json:"email"`
	ResumeFileDetails struct {
		FileName     string `json:"file"`
		ContentType  string `json:"content_type"`
		Size         int64  `json:"size"`
		NewFileName  string `json:"new_file_name"`
		DocumentType string `json:"document_type"`
	} `json:"-"`
	JobID       uuid.UUID `json:"job_id"`
	IsQualified bool      `json:"is_qualified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func SetupApplicantFormFileValues(ctx *gin.Context) (*Applicant, error) {

	yearsOfExperience, err := strconv.Atoi(ctx.Request.FormValue("years_of_experience"))
	log.Println(ctx.Request.FormValue("years_of_experience"))
	if err != nil {
		return nil, errors.New("years of experience should be a number")
	}

	applicant := Applicant{
		FirstName:         ctx.Request.FormValue("first_name"),
		LastName:          ctx.Request.FormValue("last_name"),
		YearsOfExperience: yearsOfExperience,
		Email:             ctx.Request.FormValue("email"),
	}
	return &applicant, nil
}

func (req *Applicant) Validate(v *validator.Validator) bool {

	v.Check(req.ResumeFileDetails.FileName != "", "file", "must not be blank")
	v.Check(req.ResumeFileDetails.ContentType != "", "content_type", "must not be blank")
	v.Check(req.ResumeFileDetails.Size != 0, "size", "must not be blank")

	v.Check(req.ResumeFileDetails.Size <= MaxFileUploadSize, "size", "the file is too large")

	ext := filepath.Ext(req.ResumeFileDetails.FileName)
	log.Println(ext)
	v.Check(ext == fileExtension, "file_extension", "must be a pdf file")

	req.ResumeFileDetails.NewFileName = req.JobID.String() + "_" + req.FirstName + "_" + req.LastName + ext

	return v.Valid()
}
