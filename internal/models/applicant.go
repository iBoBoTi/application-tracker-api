package models

import "github.com/google/uuid"

type Applicant struct {
	Model
	FirstName         string
	LastName          string
	YearsOfExperience int64
	Email             string
	Resume            string
	JobID             uuid.UUID
}
