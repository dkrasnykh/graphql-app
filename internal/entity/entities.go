package entity

type Comment struct {
	ID              int64
	Text            string
	ParentCommentID *int64
	PostID          int64
	UserID          int64
}

type Post struct {
	ID          int64
	Text        string
	User        int64
	CommentsOFF bool
}
