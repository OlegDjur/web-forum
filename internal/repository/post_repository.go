package repository

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"log"
)

type PostItem interface {
	CreatePost(post *models.Post) error
	GetAllPosts() (posts []models.Post, err error)
	GetPostByID(id int) (models.Post, error)
	GetPostsByCategory(category string) ([]models.Post, error)
	GetCreatedPosts(userID int) ([]models.Post, error)
	GetLikedPosts(username string) ([]models.Post, error)
	GetCategoriesByPostID(postId int) ([]string, error)
	UpdatePost(id, like, dislike int, title, content string) error
	DeletePost(id int) error
	LikePost(username string, postid int) error
	DisLikePost(username string, postid int) error
	RemoveLikePost(id int) error
	RemoveDisLikePost(id int) error
	HasUserLiked(username string, postid int) error
	HasUserDislike(username string, postid int) error
}

type PostStorage struct {
	db *sql.DB
}

func NewPostSqlite(db *sql.DB) *PostStorage {
	return &PostStorage{db: db}
}

func (p *PostStorage) CreatePost(post *models.Post) error {
	query := fmt.Sprintf(`INSERT INTO post (userid, title, content, about) values ($1, $2, $3, $4)`)
	result, err := p.db.Exec(query, post.UserID, post.Title, post.Content, post.About)
	if err != nil {
		return fmt.Errorf("storage: create post: %w", err)
	}
	postId, err := result.LastInsertId()

	query = `INSERT INTO post_category (postId, category) VALUES ($1, $2);`
	for _, oneCategory := range post.Category {
		_, err := p.db.Exec(query, postId, oneCategory)
		if err != nil {
			return fmt.Errorf("storage: create post: %w", err)
		}
	}
	return nil
}

func (p *PostStorage) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	rows, err := p.db.Query("SELECT id, userid, title, content, about FROM post")
	if err != nil {
		return nil, fmt.Errorf("storage: get all posts: query - %w", err)
	}

	for rows.Next() {
		p := models.Post{}
		if err = rows.Scan(&p.Id, &p.UserID, &p.Title, &p.Content, &p.About); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	rows.Close()
	return posts, nil
}

func (s *PostStorage) GetPostsByCategory(category string) ([]models.Post, error) {
	var p []models.Post
	query := `SELECT id, userid, title, content, about, like, dislike FROM post WHERE id IN (SELECT postId FROM post_category WHERE category=$1);`
	rows, err := s.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("storage: get post by category: %w", err)
	}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Id, &post.UserID, &post.Title, &post.Content, &post.About, &post.Like, &post.DisLike); err != nil {
			return nil, fmt.Errorf("storage: get post by category: %w", err)
		}
		p = append(p, post)
	}
	return p, nil
}

func (p *PostStorage) GetCreatedPosts(userID int) ([]models.Post, error) {
	var posts []models.Post
	rows, err := p.db.Query("SELECT id, userid, title, content, about FROM post WHERE userid=$1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := models.Post{}
		if err := rows.Scan(&p.Id, &p.UserID, &p.Title, &p.Content, &p.About); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}
	rows.Close()
	return posts, nil
}

func (p *PostStorage) GetLikedPosts(username string) ([]models.Post, error) {
	var posts []models.Post
	rows, err := p.db.Query("SELECT id, userid, title, content, about FROM post WHERE id IN (SELECT postid FROM like WHERE username=$1);", username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := models.Post{}
		if err := rows.Scan(&p.Id, &p.UserID, &p.Title, &p.Content, &p.About); err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	rows.Close()
	return posts, nil
}

func (p *PostStorage) GetPostByID(id int) (models.Post, error) {
	query := `SELECT id, title, content, like, dislike FROM post WHERE id=$1;`
	row := p.db.QueryRow(query, id)
	var post models.Post
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.Like, &post.DisLike)
	if err != nil {
		return models.Post{}, fmt.Errorf("storage: get user by login: %w", err)
	}

	return post, nil
}

func (s *PostStorage) GetCategoriesByPostID(postId int) ([]string, error) {
	queryCategory := `SELECT category FROM post_category where postId=$1;`
	categoryRows, err := s.db.Query(queryCategory, postId)
	if err != nil {
		return nil, fmt.Errorf("storage: get all category by post id: %w", err)
	}
	var category []string
	for categoryRows.Next() {
		var oneCategory string
		if err := categoryRows.Scan(&oneCategory); err != nil {
			return nil, fmt.Errorf("storage: get all category by post id: %w", err)
		}
		category = append(category, oneCategory)
	}
	return category, nil
}

func (p *PostStorage) UpdatePost(id, like, dislike int, title, content string) error {
	query, err := p.db.Prepare(`UPDATE post SET title=?, content=?, like=?, dislike=? WHERE id=?;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec(title, content, like, dislike, id)
	if err != nil {
		fmt.Println("update", err)
		return err
	}

	return nil
}

func (p *PostStorage) DeletePost(id int) error {
	query := `DELETE FROM post WHERE id=?`
	_, err := p.db.Exec(query, id)
	if err != nil {
		fmt.Println("delete", err)
		return err
	}
	return nil
}

func (p *PostStorage) LikePost(username string, postid int) error {
	query := `INSERT INTO like (username, postid) values ($1, $2)`

	_, err := p.db.Exec(query, username, postid)
	if err != nil {
		return fmt.Errorf("repository: like post: Insert query - %w", err)
	}

	query = `UPDATE post SET like = like + 1 WHERE id = $1;`

	_, err = p.db.Exec(query, postid)
	if err != nil {
		return fmt.Errorf("repository: like post: Insert query - %w", err)
	}

	return nil
}

func (p *PostStorage) DisLikePost(username string, postid int) error {
	query := `INSERT INTO dislike (username, postid) values ($1, $2)`

	_, err := p.db.Exec(query, username, postid)
	if err != nil {
		return fmt.Errorf("repository: dislike post: Insert query - %w", err)
	}

	query = `UPDATE post SET dislike = dislike + 1 WHERE id = $1;`

	_, err = p.db.Exec(query, postid)
	if err != nil {
		return fmt.Errorf("repository: dislike post: Insert query - %w", err)
	}
	return err
}

func (p *PostStorage) RemoveLikePost(id int) error {
	stmt, err := p.db.Prepare(`UPDATE post SET like = like - 1 WHERE id = $1;`)
	if err != nil {
		return fmt.Errorf("repository: remove like from post: Delete query - %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("repository: remove like from post: Delete query - %w", err)
	}
	return nil
}

func (p *PostStorage) RemoveDisLikePost(id int) error {
	stmt, err := p.db.Prepare(`UPDATE post SET dislike = dislike - 1 WHERE id = $1;`)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from post: Delete query - %w", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("repository: remove dislike from post: Update query - %w", err)
	}
	return nil
}

func (p *PostStorage) HasUserLiked(username string, postid int) error {
	var u string
	query := `SELECT username FROM like WHERE postid=? AND username = $2`

	if err := p.db.QueryRow(query, postid, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: post has like: %w", err)
	}

	query = `DELETE FROM like WHERE postid=? AND username = $2`
	if _, err := p.db.Exec(query, postid, username); err != nil {
		return err
	}

	return nil
}

func (p *PostStorage) HasUserDislike(username string, postid int) error {
	var u string
	query := `SELECT username FROM dislike WHERE postid=? AND username = $2`

	if err := p.db.QueryRow(query, postid, username).Scan(&u); err != nil {
		return fmt.Errorf("repository: post has dislike: %w", err)
	}

	query = `DELETE FROM dislike WHERE postid=? AND username = $2`
	if _, err := p.db.Exec(query, postid, username); err != nil {
		return err
	}

	return nil
}
