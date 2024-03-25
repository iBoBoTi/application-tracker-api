package repository

import (
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/models"
	"gorm.io/gorm"
)

type JobPostRepository interface {
	CreateJobPost(jobPost *models.JobPost) error
	GetJobPostsPaginatedByLatest(limit, page int) ([]models.JobPost, error)
	GetJobPostsByOwnerPaginatedByLatest(id uuid.UUID, limit, page int) ([]models.JobPost, error)
	GetUserJobPostByID(id, userID uuid.UUID) (*models.JobPost, error)
	GetJobPostByID(id uuid.UUID) (*models.JobPost, error)
}

type jobPostRepository struct {
	db *gorm.DB
}

func NewJobPostRepository(db *gorm.DB) *jobPostRepository {
	return &jobPostRepository{
		db: db,
	}
}

func (j *jobPostRepository) CreateJobPost(jobPost *models.JobPost) error {
	return j.db.Create(jobPost).Error
}

func (j *jobPostRepository) GetJobPostsPaginatedByLatest(limit, page int) ([]models.JobPost, error) {
	var jobPosts []models.JobPost
	if err := j.db.Model(&models.JobPost{}).Order("created_at ASC").Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&jobPosts).Error; err != nil {
		return nil, err
	}
	return jobPosts, nil
}

func (j *jobPostRepository) GetJobPostsByOwnerPaginatedByLatest(id uuid.UUID, limit, page int) ([]models.JobPost, error) {
	var jobPosts []models.JobPost
	if err := j.db.Model(&models.JobPost{}).Where("user_id = ?", id).Order("created_at ASC").Scopes(models.NewPaginate(limit, page).PaginatedResult).Find(&jobPosts).Error; err != nil {
		return nil, err
	}
	return jobPosts, nil
}

func (j *jobPostRepository) GetUserJobPostByID(id, userID uuid.UUID) (*models.JobPost, error) {
	var jobPost models.JobPost
	if err := j.db.Model(&models.JobPost{}).Where("id = ? AND user_id = ?", id, userID).First(&jobPost).Error; err != nil {
		return nil, err
	}
	return &jobPost, nil
}

func (j *jobPostRepository) GetJobPostByID(id uuid.UUID) (*models.JobPost, error) {
	var jobPost models.JobPost
	if err := j.db.Model(&models.JobPost{}).Where("id = ?", id).First(&jobPost).Error; err != nil {
		return nil, err
	}
	return &jobPost, nil
}
