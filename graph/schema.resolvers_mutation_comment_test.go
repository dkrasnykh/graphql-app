package graph

import (
	"context"
	"math/rand"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/service"
)

func (ts *ResolverTestSuite) TestCreateComment_BodyTooBig() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	body := randString(2001)
	inputComment := model.NewComment{Text: body, PostID: post.ID, UserID: "1"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrCommentBodyTooBig)
}

func (ts *ResolverTestSuite) TestCreateComment_InvalidPostID() {
	input := model.NewComment{Text: "comment 1", PostID: "text", UserID: "1"}
	_, err := ts.mutation.CreateComment(context.Background(), input)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestCreateComment_InvalidUserID() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "comment 1", PostID: post.ID, UserID: "text"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestCreateComment_InvalidParentID() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	parentID := "text"
	inputComment := model.NewComment{Text: "comment 1", ParentCommentID: &parentID, PostID: post.ID, UserID: "text"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestCreateComment_PostNotFound() {
	input := model.NewComment{Text: "comment 1", PostID: "1", UserID: "1"}
	_, err := ts.mutation.CreateComment(context.Background(), input)
	ts.ErrorIs(err, service.ErrPostNotFound)
}

func (ts *ResolverTestSuite) TestCreateComment_PostCommentsDisabled() {
	disabled := true
	inputPost := model.NewPost{Text: "awesome post", UserID: "1", CommentsOff: &disabled}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "comment 1", PostID: post.ID, UserID: "1"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrPostCommentsDisabled)
}

func (ts *ResolverTestSuite) TestCreateComment_ParentCommentNotFound() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	parentCommentID := "1"
	inputComment := model.NewComment{Text: "comment 1", ParentCommentID: &parentCommentID, PostID: post.ID, UserID: "1"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrInvalidParentCommentID)
}

func (ts *ResolverTestSuite) TestCreateComment_ParentCommentBelongAnotherPost() {
	inputPost1 := model.NewPost{Text: "awesome post 1", UserID: "1"}
	post1, err := ts.mutation.CreatePost(context.Background(), inputPost1)
	ts.NoError(err)

	inputPost2 := model.NewPost{Text: "awesome post 2", UserID: "1"}
	post2, err := ts.mutation.CreatePost(context.Background(), inputPost2)
	ts.NoError(err)

	parentCommentInput := model.NewComment{Text: "post2 comment", PostID: post2.ID, UserID: "1"}
	parentComment, err := ts.mutation.CreateComment(context.Background(), parentCommentInput)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "post1 comment", ParentCommentID: &parentComment.ID, PostID: post1.ID, UserID: "1"}
	_, err = ts.mutation.CreateComment(context.Background(), inputComment)
	ts.ErrorIs(err, service.ErrParentCommentBelongAnotherPost)
}

func (ts *ResolverTestSuite) TestCreateComment_RootOK() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "post1 comment", PostID: post.ID, UserID: "1"}
	comment, err := ts.mutation.CreateComment(context.Background(), inputComment)
	ts.NoError(err)
	ts.Equal(inputComment.Text, comment.Text)
	ts.Equal(inputComment.PostID, comment.PostID)
	ts.Equal(inputComment.UserID, comment.UserID)
	ts.Nil(comment.ParentCommentID)
}

func (ts *ResolverTestSuite) TestCreateComment_WithParentOK() {
	inputPost := model.NewPost{Text: "awesome post", UserID: "1"}
	post1, err := ts.mutation.CreatePost(context.Background(), inputPost)
	ts.NoError(err)

	parentCommentInput := model.NewComment{Text: "post1 comment", PostID: post1.ID, UserID: "1"}
	parentComment, err := ts.mutation.CreateComment(context.Background(), parentCommentInput)
	ts.NoError(err)

	inputComment := model.NewComment{Text: "post1 comment", ParentCommentID: &parentComment.ID, PostID: post1.ID, UserID: "1"}
	comment, err := ts.mutation.CreateComment(context.Background(), inputComment)
	ts.NoError(err)
	ts.Equal(inputComment.Text, comment.Text)
	ts.Equal(inputComment.PostID, comment.PostID)
	ts.Equal(inputComment.UserID, comment.UserID)
	ts.Equal(*inputComment.ParentCommentID, *comment.ParentCommentID)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
