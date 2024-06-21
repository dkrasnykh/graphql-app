package graph

import (
	"context"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/entity"
)

type IPost interface {
	Validate(input model.NewPost) (*entity.Post, error)
	Save(ctx context.Context, post entity.Post) (*model.Post, error)
	ValidateDisableCommentsRequest(input model.DisableCommentsRequest) (int64, int64, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error
	ByID(ctx context.Context, ID int64) (*model.Post, error)
	All(ctx context.Context) ([]*model.Post, error)
	ValidateID(postID string) (int64, error)
}

type IComment interface {
	Validate(input model.NewComment) (*entity.Comment, error)
	Save(ctx context.Context, comment entity.Comment) (*model.Comment, error)
	All(ctx context.Context, postID int64, limit *int, offset *int) ([]*model.Comment, error)
}

type Resolver struct {
	Post    IPost
	Comment IComment
}
