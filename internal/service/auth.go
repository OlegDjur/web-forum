package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"net/http"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExist       = errors.New("user exist")
)

type Authorization interface {
	CreateUser(user *models.User) error
	GenerateSessionToken(email, password string) (string, time.Time, error)
	GetSessionToken(token string) (models.User, error)
	GetSessionTokenFromRequest(r *http.Request) (models.User, error)
	DeleteSessionToken(token string) error
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user *models.User) error {
	var err error

	if err = isValidUser(user); err != nil {
		return fmt.Errorf("service: create user: %w", err)
	}

	if _, err = s.repo.GetUserByEmail(user.Email); err == nil {
		return ErrUserExist
	}

	if _, err = s.repo.GetUserByUsername(user.Username); err == nil {
		return ErrUserExist
	}

	user.Password, err = generateHashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("service: create user: %w", err)
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateSessionToken(email, password string) (string, time.Time, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", time.Time{}, err
	}

	passwordComparasionError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if passwordComparasionError != nil {
		return "", time.Time{}, passwordComparasionError
	}

	token := uuid.NewV4().String()
	expiresAt := time.Now().Add(time.Hour * 12)

	err = s.repo.AddSessionToken(email, token, expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (s *AuthService) GetSessionToken(token string) (models.User, error) {
	user, err := s.repo.GetSessionToken(token)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) GetSessionTokenFromRequest(r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return models.User{}, err
	}

	user, err := s.repo.GetSessionToken(cookie.Value)
	if err != nil {
		return models.User{}, err
	}

	if user.ExpiresAt.Before(time.Now()) {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) DeleteSessionToken(token string) error {
	err := s.repo.DeleteSessionToken(token)
	if err != nil {
		return fmt.Errorf("service: delete session token: %w", err)
	}
	return nil
}

func generateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func isValidUser(user *models.User) error {
	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if ok := validEmail.MatchString(user.Email); !ok {
		return ErrInvalidEmail
	}

	for _, char := range user.Email {
		if char < 33 || char > 126 {
			return ErrInvalidUsername
		}
	}

	for _, char := range user.Username {
		if char < 32 || char > 126 {
			return ErrInvalidUsername
		}
	}

	if len(user.Username) < 2 || len(user.Username) >= 20 {
		return ErrInvalidUsername
	}

	for _, char := range user.Password {
		if char < 33 || char > 126 {
			return ErrInvalidUsername
		}
	}

	if len(user.Password) > 20 || len(user.Password) < 6 {
		return ErrInvalidPassword
	}

	return nil
}
