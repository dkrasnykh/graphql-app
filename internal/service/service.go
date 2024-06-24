package service

import (
	"context"
	"errors"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/subscription"
)

var (
	ErrInvalidID                      = errors.New("error converting post id into int64 type")
	ErrInternal                       = errors.New("internal error")
	ErrEmptyBody                      = errors.New("text value should not be empty")
	ErrCommentBodyTooBig              = errors.New("comment text should not exceed 2000 characters")
	ErrPostNotFound                   = errors.New("post with id does not exist")
	ErrAccess                         = errors.New("post keeper is another user")
	ErrPostCommentsDisabled           = errors.New("сomments are turned off")
	ErrInvalidParentCommentID         = errors.New("there is no comment with ParentCommentID for this post")
	ErrParentCommentBelongAnotherPost = errors.New("parent comment belong another post")
)

type Storager interface {
	SavePost(ctx context.Context, post entity.Post) (int64, error)
	PostByID(ctx context.Context, id int64) (*entity.Post, error)
	AllPosts(ctx context.Context) ([]*entity.Post, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error

	SaveComment(ctx context.Context, comment entity.Comment) (int64, error)
	AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error)

	Clear() // for unit tests (implemented only for memory storage)
}

type Service struct {
	storage       Storager
	subscriptions *subscription.Subscription
}

func New(storage Storager, subscriptions *subscription.Subscription) *Service {
	return &Service{
		storage:       storage,
		subscriptions: subscriptions,
	}
}
