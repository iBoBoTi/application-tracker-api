package repository

import (
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/models"
	"gorm.io/gorm"
)

type ApplicantRepository interface {
	CreateApplicant(applicant *models.Applicant) error
	GetAllApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error)
	GetQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error)
	GetUnQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error)
	GetApplicantByJobID(id, jobID uuid.UUID) (*models.Applicant, error)
}

type applicantRepository struct {
	db *gorm.DB
}

func NewApplicantRepository(db *gorm.DB) *applicantRepository {
	return &applicantRepository{
		db: db,
	}
}

func (a *applicantRepository) CreateApplicant(applicant *models.Applicant) error {
	return a.db.Create(applicant).Error
}

func (a *applicantRepository) GetAllApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("job_id = ?", jobID).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("job_id = ? AND is_qualified = ?", jobID, true).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetUnQualifiedApplicantsByJobIDPaginated(jobID uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("job_id = ? AND is_qualified = ?", jobID, false).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetApplicantByJobID(id, jobID uuid.UUID) (*models.Applicant, error) {
	var applicant models.Applicant

	if err := a.db.Model(&models.Applicant{}).Where("id = ? AND job_id = ?", id, jobID).First(&applicant).Error; err != nil {
		return nil, err
	}
	return &applicant, nil
}
