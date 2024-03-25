package dtos

import (
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/validator"
	"time"
)

type JobPost struct {
	ID                   uuid.UUID `json:"id"`
	Title                string    `json:"title" binding:"required"`
	CompanyName          string    `json:"company_name" binding:"required"`
	Description          string    `json:"description" binding:"required"`
	MinYearsOfExperience int64     `json:"min_years" binding:"required"`
	MaxYearsOfExperience int64     `json:"max_years" binding:"required"`
	UserID               uuid.UUID `json:"-"`
	Keywords             []string  `json:"keywords"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (j *JobPost) Validate(v *validator.Validator) bool {
	v.Check(j.Title != "", "title", "must not be blank")
	v.Check(j.CompanyName != "", "company_name", "must not be blank")
	v.Check(j.Description != "", "description", "must not be blank")

	return v.Valid()
}
