package models

type Post struct {
	Id       int
	UserID   int
	Category []string
	Title    string
	Content  string
	About    string
	Comments int
	Like     int
	DisLike  int
}

func NewPost(id, like, dislike, userID, comments int, title, content, about string, category []string) *Post {
	return &Post{id, userID, category, title, content, about, comments, like, dislike}
}
