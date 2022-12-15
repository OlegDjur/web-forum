package models

type Comment struct {
	ID       int
	PostID   int
	Author string
	Text     string
	Likes    int
	DisLikes int
}
