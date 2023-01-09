package repository

import (
	"fmt"
	"forum/internal/models"

	"github.com/jmoiron/sqlx"
)

type Comment interface {
	CreateComment(comment *models.Comment) error
	GetComments(postID int) ([]*models.Comment, error)
	GetCommentByID(commentID int) (models.Comment, error)
	CommentHasLike(commentID int, username string) error
	CommentHasDislike(commentID int, username string) error
	RemoveLikeComment(commentID int, username string) error
	RemoveDislikeComment(commentID int, username string) error
	LikeComment(commentID int, username string) error
	DislikeComment(commentID int, username string) error
}

type CommentStorage struct {
	db *sqlx.DB
}

func NewCommentSqlite(db *sqlx.DB) *CommentStorage {
	return &CommentStorage{db: db}
}

func (c *CommentStorage) CreateComment(comment *models.Comment) error {
	query := fmt.Sprintf(`INSERT INTO comment (author, text, postid) values ($1, $2, $3)`)
	_, err := c.db.Exec(query, comment.Author, comment.Text, comment.PostID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentStorage) GetComments(postID int) ([]*models.Comment, error) {
	var comments []*models.Comment
	query := fmt.Sprintf(`SELECT * FROM comment WHERE postid = $1 order by id asc;`)
	rows, err := c.db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("repository: get commentaries of the post: query - %w", err)
	}

	for rows.Next() {
		c := &models.Comment{}
		if err = rows.Scan(&c.ID, &c.Author, &c.PostID, &c.Text, &c.Likes, &c.DisLikes); err != nil {
			return nil, fmt.Errorf("repository: get commentaries of the post: query - %w", err)
		}
		comments = append(comments, c)
	}
	rows.Close()
	return comments, nil
}

func (c *CommentStorage) GetCommentByID(commentID int) (models.Comment, error) {
	var comment models.Comment

	query := `SELECT id, postid, author, text, likes, dislike FROM comment WHERE id=$1;`
	row := c.db.QueryRow(query, commentID)

	err := row.Scan(&comment.ID, &comment.PostID, &comment.Author, &comment.Text, &comment.Likes, &comment.DisLikes)
	if err != nil {
		return models.Comment{}, fmt.Errorf("storage: get user by login: %w", err)
	}

	return comment, nil
}

func (s *CommentStorage) RemoveLikeComment(commentID int, username string) error {
	query := `DELETE FROM likes WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	query = `UPDATE comment SET likes = likes - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	return nil
}

func (s *CommentStorage) RemoveDislikeComment(commentID int, username string) error {
	query := `DELETE FROM dislike WHERE commentId = $1 AND username = $2;`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	query = `UPDATE comment SET dislike = dislike - 1 WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: remove like from comment: %w", err)
	}
	return nil
}

func (s *CommentStorage) CommentHasLike(commentID int, username string) error {
	var u, query string
	query = `SELECT username FROM likes WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentID, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has like: %w", err)
	}
	return nil
}

func (s *CommentStorage) CommentHasDislike(commentID int, username string) error {
	var u, query string
	query = `SELECT username FROM dislike WHERE commentId = $1 AND username = $2;`
	err := s.db.QueryRow(query, commentID, username).Scan(&u)
	if err != nil {
		return fmt.Errorf("storage: comment has like: %w", err)
	}
	return nil
}

func (s *CommentStorage) LikeComment(commentID int, username string) error {
	query := `INSERT INTO likes(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	query = `UPDATE comment SET likes = likes + 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	return nil
}

func (s *CommentStorage) DislikeComment(commentID int, username string) error {
	query := `INSERT INTO dislike(commentId, username) VALUES ($1, $2);`
	_, err := s.db.Exec(query, commentID, username)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	query = `UPDATE comment SET dislike = dislike + 1  WHERE id = $1;`
	_, err = s.db.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("storage: like comment: %w", err)
	}
	return nil
}
