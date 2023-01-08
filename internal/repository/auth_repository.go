package repository

import (
	"fmt"
	"forum/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByUsername(username string) (models.User, error)
	AddSessionToken(email, token string, expiresAt time.Time) error
	GetSessionToken(token string) (models.User, error)
	DeleteSessionToken(token string) error
}

type AuthStorage struct {
	db *sqlx.DB
}

func NewAuthSqlite(db *sqlx.DB) *AuthStorage {
	return &AuthStorage{db: db}
}

func (r *AuthStorage) CreateUser(user *models.User) error {
	query := fmt.Sprintf("INSERT INTO users (username, email, password) values ($1, $2, $3)")
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthStorage) GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE email=$1;`
	row := s.db.QueryRow(query, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by login: %w", err)
	}
	return user, nil
}

func (s *AuthStorage) GetUserByUsername(username string) (models.User, error) {
	query := `SELECT id, email, username, password FROM users WHERE username=$1;`
	row := s.db.QueryRow(query, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by login: %w", err)
	}
	return user, nil
}

func (s *AuthStorage) AddSessionToken(email, token string, expiresAt time.Time) error {
	query := `UPDATE users SET token = $1, expiresAt = $2 WHERE email = $3;`
	_, err := s.db.Exec(query, token, expiresAt, email)
	if err != nil {
		return fmt.Errorf("storage: save session token: %w", err)
	}
	return nil
}

func (s *AuthStorage) GetSessionToken(token string) (models.User, error) {
	query := `SELECT id, email, username, password, token, expiresAt FROM users WHERE token=$1;`

	row := s.db.QueryRow(query, token)
	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Token, &user.ExpiresAt)
	if err != nil {
		return models.User{}, fmt.Errorf("storage: get user by login: %w", err)
	}
	return user, nil
}

func (s *AuthStorage) DeleteSessionToken(token string) error {
	query := `UPDATE users SET token = NULL, expiresAt = NULL WHERE token = $1;`
	_, err := s.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("storage: delete session token: %w", err)
	}
	return nil
}
