package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/dkrasnykh/graphql-app/graph/model"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	post, err := r.Service.ValidatePost(input)
	if err != nil {
		return nil, err
	}

	return r.Service.SavePost(ctx, *post)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*model.Comment, error) {
	comment, err := r.Service.ValidateComment(input)
	if err != nil {
		return nil, err
	}

	return r.Service.SaveComment(ctx, *comment)
}

// DisableComments is the resolver for the disableComments field.
func (r *mutationResolver) DisableComments(ctx context.Context, input model.DisableCommentsRequest) (bool, error) {
	userID, postID, err := r.Service.ValidateDisableCommentsRequest(input)
	if err != nil {
		return false, err
	}

	if err := r.Service.DisableComments(ctx, userID, postID); err != nil {
		return false, err
	}

	return true, nil
}

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context) ([]*model.Post, error) {
	return r.Service.AllPosts(ctx)
}

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	postID, err := r.Service.ValidateID(id)
	if err != nil {
		return nil, err
	}

	return r.Service.PostById(ctx, postID)
}

// Comments is the resolver for the comments field.
func (r *queryResolver) Comments(ctx context.Context, postID string, limit *int, offset *int) ([]*model.Comment, error) {
	id, err := r.Service.ValidateID(postID)
	if err != nil {
		return nil, err
	}

	return r.Service.AllComments(ctx, id, limit, offset)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
