package models

import "github.com/google/uuid"

type JobPost struct {
	Model
	Title                string      `gorm:"not null"`
	CompanyName          string      `gorm:"not null"`
	Description          string      `gorm:"not null"`
	MinYearsOfExperience int64       `gorm:"not null"`
	MaxYearsOfExperience int64       `gorm:"not null"`
	UserID               uuid.UUID   `gorm:"not null"`
	Keywords             string      `gorm:"not null"`
	Applicants           []Applicant `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}
