package models

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	Token     string
	ExpiresAt time.Time
}

