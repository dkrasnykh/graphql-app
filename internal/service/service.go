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

type Storager interface {
	SavePost(ctx context.Context, post entity.Post) (int64, error)
	PostByID(ctx context.Context, id int64) (*entity.Post, error)
	AllPosts(ctx context.Context) ([]*entity.Post, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error

	SaveComment(ctx context.Context, comment entity.Comment) (int64, error)
	AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error)
}

type Service struct {
	storage Storager
}

func New(storage Storager) *Service {
	return &Service{
		storage: storage,
	}
}
