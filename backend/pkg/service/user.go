package service

import (
	"easyinvesting/pkg/dto"
	"easyinvesting/pkg/model"
	"easyinvesting/pkg/repository"
	"fmt"
)

type UserService interface {
	Save(userDTO *dto.UserDTO) error
	Login(email, password string) (*dto.UserDTO, error)
	FindByEmail(email string) (*dto.UserDTO, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) Save(userDTO *dto.UserDTO) error {
	u := &model.User{
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

func (s *userService) Login(email, password string) (*dto.UserDTO, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	if !u.CheckPassword(password) {
		return nil, fmt.Errorf("Invalid password")
	}
	userDTO := &dto.UserDTO{
		ID:    u.ID,
		Email: u.Email,
	}
	return userDTO, nil
}

func (s *userService) FindByEmail(email string) (*dto.UserDTO, error) {
	u, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	userDTO := &dto.UserDTO{
		ID:    u.ID,
		Email: u.Email,
	}
	return userDTO, nil
}
