package service

import (
	"context"
	"errors"

	"github.com/dkrasnykh/graphql-app/internal/entity"
)

var (
	ErrInvalidID         = errors.New("error converting post id into int64 type")
	ErrInternal          = errors.New("internal error")
	ErrEmptyBody         = errors.New("text value should not be empty")
	ErrCommentBodyTooBig = errors.New("comment text should not exceed 2000 characters")
)

type PostStorager interface {
	Save(ctx context.Context, post entity.Post) (int64, error)
	ByID(ctx context.Context, id int64) (*entity.Post, error)
	All(ctx context.Context) ([]*entity.Post, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error
}

type CommentStorager interface {
	Save(ctx context.Context, comment entity.Comment) (int64, error)
	All(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error)
}
