package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	"forum/internal/repository"
	"strings"
)

var ErrInvalidPost = errors.New("invalid post")

type PostItem interface {
	CreatePost(post *models.Post) error
	GetAllPosts() (posts []models.Post, err error)
	GetPostsByCategory(category string) ([]models.Post, error)
	GetCreatedPosts(userID int) ([]models.Post, error)
	GetLikedPosts(username string) ([]models.Post, error)
	GetPostByID(id int) (models.Post, error)
	UpdatePost(id, like, dislike int, title, content string) error
	DeletePost(id int) error
	LikePost(username string, postid int) error
	DisLikePost(username string, postid int) error
}

type PostService struct {
	repo repository.PostItem
}

func NewPostService(repo repository.PostItem) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(post *models.Post) error {
	post.Category = strings.Split(post.Category[0], ",")

	if err := isValidPost(post); err != nil {
		return err
	}

	return p.repo.CreatePost(post)
}

func (p *PostService) GetAllPosts() ([]models.Post, error) {
	posts, err := p.repo.GetAllPosts()
	if err != nil {
		return []models.Post{}, err
	}

	for i := range posts {
		category, err := p.repo.GetCategoriesByPostID(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post: %w", err)
		}
		posts[i].Category = category
	}

	return posts, nil
}

func (p *PostService) GetPostsByCategory(category string) ([]models.Post, error) {
	posts, err := p.repo.GetPostsByCategory(category)
	if err != nil {
		return []models.Post{}, err
	}

	for i := range posts {
		category, err := p.repo.GetCategoriesByPostID(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post: %w", err)
		}
		posts[i].Category = category
	}

	return posts, nil
}

func (p *PostService) GetCreatedPosts(userID int) ([]models.Post, error) {
	posts, err := p.repo.GetCreatedPosts(userID)
	if err != nil {
		return []models.Post{}, err
	}

	for i := range posts {
		category, err := p.repo.GetCategoriesByPostID(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post: %w", err)
		}
		posts[i].Category = category
	}

	return posts, nil
}

func (p *PostService) GetLikedPosts(username string) ([]models.Post, error) {
	posts, err := p.repo.GetLikedPosts(username)
	if err != nil {
		return []models.Post{}, err
	}

	for i := range posts {
		category, err := p.repo.GetCategoriesByPostID(posts[i].Id)
		if err != nil {
			return nil, fmt.Errorf("service: get all post: %w", err)
		}
		posts[i].Category = category
	}
	return posts, nil
}

func (p *PostService) GetPostByID(id int) (posts models.Post, err error) {
	post, err := p.repo.GetPostByID(id)
	if err != nil {
		return models.Post{}, err
	}

	post.Category, err = p.repo.GetCategoriesByPostID(id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (p *PostService) LikePost(username string, postid int) error {
	if err := p.repo.HasUserLiked(username, postid); err != nil {
		if err = p.repo.HasUserDislike(username, postid); err == nil {
			if err = p.repo.RemoveDisLikePost(postid); err != nil {
				return err
			}
		}
		return p.repo.LikePost(username, postid)

	}

	return p.repo.RemoveLikePost(postid)
}

func (p *PostService) DisLikePost(username string, postid int) error {
	if err := p.repo.HasUserDislike(username, postid); err != nil {
		if err := p.repo.HasUserLiked(username, postid); err == nil {
			if err = p.repo.RemoveLikePost(postid); err != nil {
				return err
			}
		}
		return p.repo.DisLikePost(username, postid)

	}
	return p.repo.RemoveDisLikePost(postid)
}

func isValidPost(post *models.Post) error {
	if len(post.Title) > 100 {
		return errors.New("title length out of range")
	}

	if len(post.About) > 500 {
		return errors.New("content length out of range")
	}

	post.Title = strings.Trim(post.Title, " \n\r")

	for _, char := range post.Title {
		if (char != 13 && char != 10) && (char < 32 || char > 126) {
			return ErrInvalidPost
		}
	}

	if post.Title == "" {
		return ErrInvalidPost
	}

	post.About = strings.Trim(post.About, " \n\r")

	for _, char := range post.About {
		if (char != 13 && char != 10) && (char < 32 || char > 126) {
			return ErrInvalidPost
		}
	}

	if post.About == "" {
		return ErrInvalidPost
	}

	post.Content = strings.Trim(post.Content, " \n\r")

	for _, char := range post.Content {
		if (char != 13 && char != 10) && (char < 32 || char > 126) {
			return ErrInvalidPost
		}
	}

	if post.Content == "" {
		return ErrInvalidPost
	}

	return nil
}

func (p *PostService) UpdatePost(id, like, dislike int, title, content string) error {
	return p.repo.UpdatePost(id, like, dislike, title, content)
}

func (p *PostService) DeletePost(id int) error {
	return p.repo.DeletePost(id)
}
