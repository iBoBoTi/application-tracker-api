package models

type User struct {
	Model
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	UserType     string `gorm:"not null"`
	JobPostings  []JobPosting
}
