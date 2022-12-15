package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
)

var ErrInvalidComment = errors.New("invalid comment")

type Comment interface {
	CreateComment(comment *models.Comment) error
	GetComments(postID int) ([]*models.Comment, error)
	GetCommentByID(commentID int) (models.Comment, error)
	LikeComment(commentID int, username string) error
	DislikeComment(commentID int, username string) error
}

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (c *CommentService) CreateComment(comment *models.Comment) error {
	if err := isValidComment(comment); err != nil {
		return err
	}

	return c.repo.CreateComment(comment)
}

func (c *CommentService) GetComments(postID int) ([]*models.Comment, error) {
	return c.repo.GetComments(postID)
}

func (c *CommentService) GetCommentByID(commentID int) (models.Comment, error) {
	return c.repo.GetCommentByID(commentID)
}

func (c *CommentService) LikeComment(commentID int, username string) error {
	if err := c.repo.CommentHasLike(commentID, username); err == nil {
		if err := c.repo.RemoveLikeComment(commentID, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
		return nil
	}
	fmt.Print(1)
	if err := c.repo.CommentHasDislike(commentID, username); err == nil {
		if err := c.repo.RemoveDislikeComment(commentID, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
	}

	if err := c.repo.LikeComment(commentID, username); err != nil {
		return fmt.Errorf("service: like comment: %w", err)
	}

	return nil
}

func (c *CommentService) DislikeComment(commentID int, username string) error {
	if err := c.repo.CommentHasDislike(commentID, username); err == nil {
		if err := c.repo.RemoveDislikeComment(commentID, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
		return nil
	}
	if err := c.repo.CommentHasLike(commentID, username); err == nil {
		if err := c.repo.RemoveLikeComment(commentID, username); err != nil {
			return fmt.Errorf("service: like comment: %w", err)
		}
	}

	if err := c.repo.DislikeComment(commentID, username); err != nil {
		return fmt.Errorf("service: like comment: %w", err)
	}

	return nil
}

func isValidComment(comment *models.Comment) error {
	if len(comment.Text) > 500 {
		return fmt.Errorf("service: create comment: %w", ErrInvalidComment)
	}

	comment.Text = strings.Trim(comment.Text, " \n\r")

	for _, char := range comment.Text {
		if (char != 13 && char != 10) && (char < 32 || char > 126) {
			return fmt.Errorf("service: CreatePost: isValidComment err: %v", ErrInvalidComment)
		}
	}

	if comment.Text == "" {
		return fmt.Errorf("service: CreatePost: isValidComment err: %v", ErrInvalidComment)
	}

	return nil
}
