package graph

import (
	"context"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/entity"
)

type IService interface {
	ValidatePost(input model.NewPost) (*entity.Post, error)
	SavePost(ctx context.Context, post entity.Post) (*model.Post, error)
	ValidateDisableCommentsRequest(input model.DisableCommentsRequest) (int64, int64, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error
	PostById(ctx context.Context, ID int64) (*model.Post, error)
	AllPosts(ctx context.Context) ([]*model.Post, error)
	ValidateID(ID string) (int64, error)

	ValidateComment(input model.NewComment) (*entity.Comment, error)
	SaveComment(ctx context.Context, comment entity.Comment) (*model.Comment, error)
	AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*model.Comment, error)
}

type Resolver struct {
	Service IService
}
