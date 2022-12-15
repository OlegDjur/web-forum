package service

import (
	"forum/internal/repository"
)

type Service struct {
	Authorization
	PostItem
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		PostItem:      NewPostService(repos.PostItem),
		Comment:       NewCommentService(repos.Comment),
	}
}
