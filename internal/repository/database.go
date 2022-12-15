package repository

import (
	"database/sql"
	"fmt"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	tables := []string{userTable, postTable, commentTable, likeTable, dislikeTable, postCategoryTable}
	for _, v := range tables {
		_, err := db.Exec(v)
		if err != nil {
			return err
		}
	}
	fmt.Println(true)
	return nil
}

const userTable = `CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE,
	username TEXT UNIQUE,
	password TEXT,
	token TEXT DEFAULT NULL,
	expiresAt DATETIME DEFAULT NULL
);`

const postTable = `CREATE TABLE IF NOT EXISTS post (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	userid INTEGER,
	title TEXT,
	content TEXT,
	about TEXT,
	category TEXT,
	like INTEGER DEFAULT 0,
	dislike INTEGER DEFAULT 0,
	userliked INTEGER Default 0
);`

const postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category (
	postID INTEGER,
	category TEXT,
	FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE
);`

const commentTable = `CREATE TABLE IF NOT EXISTS comment (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	author TEXT,
	postid INTEGER,
	text TEXT,
	like INTEGER DEFAULT 0,
	dislike INTEGER DEFAULT 0
);`

const likeTable = `CREATE TABLE IF NOT EXISTS like (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	postid INTEGER,
	commentId INTEGER DEFAULT NULL
);`

const dislikeTable = `CREATE TABLE IF NOT EXISTS dislike(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT,
	postid INTEGER,
	commentId INTEGER DEFAULT NULL
);`
