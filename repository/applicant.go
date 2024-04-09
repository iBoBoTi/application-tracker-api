package repository

import (
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/models"
	"gorm.io/gorm"
)

type ApplicantRepository interface {
	CreateApplicant(applicant *models.Applicant) error
	GetApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error)
	GetQualifiedApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error)
	GetUnQualifiedApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error)
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

func (a *applicantRepository) GetApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("id = ?", id).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetQualifiedApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("id = ? AND is_qualified = ?", id, true).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetUnQualifiedApplicantsByJobIDPaginated(id uuid.UUID, limit, page int) ([]models.Applicant, error) {

	var applicants []models.Applicant
	if err := a.db.Model(&models.Applicant{}).Where("id = ? AND is_qualified = ?", id, false).Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&applicants).Error; err != nil {
		return nil, err
	}
	return applicants, nil
}

func (a *applicantRepository) GetApplicantByJobID(id, jobID uuid.UUID) (*models.Applicant, error) {
	var applicant models.Applicant

	if err := a.db.Model(&models.Applicant{}).Where("id = ? AND job_id = ?", id, jobID).First(applicant).Error; err != nil {
		return nil, err
	}
	return &applicant, nil
}
