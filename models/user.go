package models

const (
	AdminUserType string = "admin"
)

type User struct {
	Model
	FirstName   string
	LastName    string
	Email       string
	Password    string
	UserType    string
	JobPostings []JobPosting
}
