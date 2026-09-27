package service

import (
	"errors"
	"net/mail"
	"strings"

	"go-api-practice/internal/model"
	"go-api-practice/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNameRequired = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is invalid")
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetAll() []model.User {
	return s.repository.FindAll()
}

func (s *UserService) GetByID(id int) (model.User, error) {
	user, found := s.repository.FindByID(id)
	if !found {
		return model.User{}, ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) Create(name, email string) (model.User, error) {
	user, err := validateUser(name, email)
	if err != nil {
		return model.User{}, err
	}

	return s.repository.Create(user), nil
}

func (s *UserService) Update(id int, name, email string) (model.User, error) {
	user, err := validateUser(name, email)
	if err != nil {
		return model.User{}, err
	}

	updatedUser, found := s.repository.Update(id, user)
	if !found {
		return model.User{}, ErrUserNotFound
	}

	return updatedUser, nil
}

func (s *UserService) Delete(id int) error {
	if deleted := s.repository.Delete(id); !deleted {
		return ErrUserNotFound
	}

	return nil
}

func validateUser(name, email string) (model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return model.User{}, ErrNameRequired
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return model.User{}, ErrInvalidEmail
	}

	return model.User{Name: name, Email: email}, nil
}
