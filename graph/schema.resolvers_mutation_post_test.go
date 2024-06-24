package graph

import (
	"context"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/service"
)

func (ts *ResolverTestSuite) TestCreatePost_EmptyBody() {
	input := model.NewPost{Text: "", UserID: "1"}
	_, err := ts.mutation.CreatePost(context.Background(), input)
	ts.ErrorIs(err, service.ErrEmptyBody)
}

func (ts *ResolverTestSuite) TestCreatePost_InvalidUserID() {
	input := model.NewPost{Text: "awesome post", UserID: "text"}
	_, err := ts.mutation.CreatePost(context.Background(), input)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestCreatePost_OK() {
	input := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), input)
	ts.NoError(err)
	ts.Equal(input.Text, post.Text)
	ts.Equal(input.UserID, post.UserID)
}

func (ts *ResolverTestSuite) TestPostDisableComments_InvalidUserID() {
	input := model.DisableCommentsRequest{UserID: "text", PostID: "1"}
	_, err := ts.mutation.DisableComments(context.Background(), input)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestPostDisableComments_InvalidPostID() {
	input := model.DisableCommentsRequest{UserID: "1", PostID: "text"}
	_, err := ts.mutation.DisableComments(context.Background(), input)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestPostDisableComments_PostNotFound() {
	input := model.DisableCommentsRequest{UserID: "1", PostID: "1"}
	_, err := ts.mutation.DisableComments(context.Background(), input)
	ts.ErrorIs(err, service.ErrPostNotFound)
}

func (ts *ResolverTestSuite) TestPostDisableComments_PostBelongAnotherUser() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	input := model.DisableCommentsRequest{UserID: "2", PostID: post.ID}
	_, err = ts.mutation.DisableComments(context.Background(), input)
	ts.ErrorIs(err, service.ErrAccess)
}

func (ts *ResolverTestSuite) TestPostDisableComments_CommentsAlreadyDisabled() {
	disabled := true
	inputPost := model.NewPost{Text: "awesome post", UserID: "1", CommentsOff: &disabled}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	input := model.DisableCommentsRequest{UserID: "1", PostID: post.ID}
	_, err = ts.mutation.DisableComments(context.Background(), input)
	ts.ErrorIs(err, service.ErrPostCommentsDisabled)
}

func (ts *ResolverTestSuite) TestPostDisableComments_OK() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "comment1", PostID: post.ID, UserID: "2"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.NoError(err)

	input := model.DisableCommentsRequest{UserID: "1", PostID: post.ID}
	_, err = ts.mutation.DisableComments(context.Background(), input)
	ts.NoError(err)

	inputComment = model.NewComment{Text: "comment1", PostID: post.ID, UserID: "2"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrPostCommentsDisabled)
}
