package repository

import (
	"database/sql"
)

type Repository struct {
	Authorization
	PostItem
	Comment
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
		PostItem:      NewPostSqlite(db),
		Comment:       NewCommentSqlite(db),
	}
}
