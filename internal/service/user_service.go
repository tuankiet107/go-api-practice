package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"go-api-practice/internal/model"
	"go-api-practice/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNameRequired = errors.New("name is required")
	ErrNameTooLong  = errors.New("name must not exceed 100 characters")
	ErrInvalidEmail = errors.New("email is invalid")
	ErrEmailTooLong = errors.New("email must not exceed 255 characters")
	ErrEmailExists  = errors.New("email already exists")
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (model.User, error) {
	user, err := s.repository.FindByID(ctx, id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return model.User{}, ErrUserNotFound
	}
	if err != nil {
		return model.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, name, email string) (model.User, error) {
	user, err := validateUser(name, email)
	if err != nil {
		return model.User{}, err
	}

	createdUser, err := s.repository.Create(ctx, user)
	if errors.Is(err, repository.ErrEmailExists) {
		return model.User{}, ErrEmailExists
	}
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) Update(ctx context.Context, id int, name, email string) (model.User, error) {
	user, err := validateUser(name, email)
	if err != nil {
		return model.User{}, err
	}

	updatedUser, err := s.repository.Update(ctx, id, user)
	if errors.Is(err, repository.ErrUserNotFound) {
		return model.User{}, ErrUserNotFound
	}
	if errors.Is(err, repository.ErrEmailExists) {
		return model.User{}, ErrEmailExists
	}
	if err != nil {
		return model.User{}, fmt.Errorf("update user: %w", err)
	}

	return updatedUser, nil
}

func (s *UserService) Delete(ctx context.Context, id int) error {
	err := s.repository.Delete(ctx, id)
	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func validateUser(name, email string) (model.User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return model.User{}, ErrNameRequired
	}
	if utf8.RuneCountInString(name) > 100 {
		return model.User{}, ErrNameTooLong
	}
	if utf8.RuneCountInString(email) > 255 {
		return model.User{}, ErrEmailTooLong
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return model.User{}, ErrInvalidEmail
	}

	return model.User{Name: name, Email: email}, nil
}
