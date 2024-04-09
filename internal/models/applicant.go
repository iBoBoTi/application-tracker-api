package models

import "github.com/google/uuid"

type Applicant struct {
	Model
	FirstName         string    `gorm:"not null"`
	LastName          string    `gorm:"not null"`
	YearsOfExperience int       `gorm:"not null"`
	Email             string    `gorm:"unique;not null"`
	Resume            []byte    `gorm:"not null;type:bytea"`
	ResumeFileName    string    `gorm:"not null"`
	JobID             uuid.UUID `gorm:"not null"`
	IsQualified       bool      `gorm:"not null;default:false"`
}
