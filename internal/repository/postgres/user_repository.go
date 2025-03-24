package postgres

import (
	"errors"

	"github.com/google/uuid"
	"github.com/firllyw/gogo/internal/domain"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
}

// userRepo is the implementation of UserRepository
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

// Create adds a new user to the database
func (r *userRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by their ID
func (r *userRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// Update modifies an existing user
func (r *userRepo) Update(user *domain.User) error {
	result := r.db.Model(user).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Delete removes a user from the database
func (r *userRepo) Delete(id uuid.UUID) error {
	result := r.db.Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}