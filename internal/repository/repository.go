package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	Authorization
	PostItem
	Comment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSqlite(db),
		PostItem:      NewPostSqlite(db),
		Comment:       NewCommentSqlite(db),
	}
}
