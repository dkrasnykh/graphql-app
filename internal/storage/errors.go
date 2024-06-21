package storage

import "errors"

var (
	ErrPostNotFound                   = errors.New("post with id does not exist")
	ErrAccess                         = errors.New("post keeper is another user")
	ErrPostCommentsDisabled           = errors.New("сomments are turned off")
	ErrInvalidParentCommentID         = errors.New("there is no comment with ParentCommentID for this post")
	ErrInternal                       = errors.New("database connection failed")
	ErrParentCommentBelongAnotherPost = errors.New("parent comment belong another post")
)
