package models

import "github.com/google/uuid"

type JobPosting struct {
	Model
	Title       string
	Description string
	UserID      uuid.UUID
	Keywords    string
	Applicants  []Applicant `gorm:"foreignKey:JobID"`
}
