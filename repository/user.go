package repository

import (
	"github.com/google/uuid"
	"github.com/iBoBoTi/ats/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindUserByID(id uuid.UUID) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	EmailExist(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) CreateUser(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) FindUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	if err := u.db.Model(&models.User{}).Where("id", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := u.db.Model(&models.User{}).Where("email", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) EmailExist(email string) (bool, error) {
	var num int
	tx := u.db.Raw("SELECT CASE WHEN EXISTS (SELECT * FROM users WHERE email = ?) THEN CAST(1 AS BIT)ELSE CAST(0 AS BIT) END", email).Scan(&num)
	if num == 1 {
		return true, tx.Error
	}
	return false, tx.Error
}
