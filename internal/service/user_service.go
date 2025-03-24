package service

import (
	"github.com/google/uuid"
	"github.com/firllyw/gogo/internal/domain"
	"github.com/firllyw/gogo/internal/repository/postgres"
)

// UserService defines the interface for user service operations
type UserService interface {
	CreateUser(req domain.CreateUserRequest) (*domain.UserResponse, error)
	GetUserByID(id uuid.UUID) (*domain.UserResponse, error)
	UpdateUser(id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(id uuid.UUID) error
}

// userService is the implementation of UserService
type userService struct {
	userRepo postgres.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo postgres.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req domain.CreateUserRequest) (*domain.UserResponse, error) {
	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by their ID
func (s *userService) GetUserByID(id uuid.UUID) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateUser modifies an existing user
func (s *userService) UpdateUser(id uuid.UUID, req domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser removes a user
func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}