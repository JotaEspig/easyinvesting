package services

import (
	"easyinvesting/pkg/dtos"
	"easyinvesting/pkg/models"
	"easyinvesting/pkg/repositories"
	"fmt"
)

type UserService interface {
	Save(userDTO *dtos.UserDTO) error
	Login(email, password string) (*dtos.UserDTO, error)
	FindByEmail(email string) (*dtos.UserDTO, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Save(userDTO *dtos.UserDTO) error {
	u := &models.User{
		Email: userDTO.Email,
	}
	u.SetPassword(userDTO.Password)
	if !u.IsValid() {
		return fmt.Errorf("Invalid user data")
	}
	if err := s.userRepository.Save(u); err != nil {
		return fmt.Errorf("Failed to save user: %w", err)
	}
	userDTO.ID = u.ID
	return nil
}

func (s *userService) Login(email, password string) (*dtos.UserDTO, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	if !u.CheckPassword(password) {
		return nil, fmt.Errorf("Invalid password")
	}
	userDTO := &dtos.UserDTO{
		ID:    u.ID,
		Email: u.Email,
	}
	return userDTO, nil
}

func (s *userService) FindByEmail(email string) (*dtos.UserDTO, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	userDTO := &dtos.UserDTO{
		ID:    u.ID,
		Email: u.Email,
	}
	return userDTO, nil
}
