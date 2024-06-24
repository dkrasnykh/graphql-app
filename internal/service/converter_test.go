package service

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/entity"
)

func TestConvertNewPostModelIntoEntity(t *testing.T) {
	disabled := true
	newPost := model.NewPost{Text: "awesome post", UserID: "1", CommentsOff: &disabled}

	target := convertNewPostModelIntoEntity(newPost)

	assert.Equal(t, newPost.Text, target.Text)
	assert.Equal(t, newPost.UserID, strconv.FormatInt(target.User, 10))
	assert.Equal(t, *newPost.CommentsOff, target.CommentsOFF)
}

func TestConvertPostEntityIntoModel(t *testing.T) {
	post := entity.Post{ID: int64(1), Text: "awesome post", User: int64(1), CommentsOFF: true}

	target := convertPostEntityIntoModel(post)
	assert.Equal(t, strconv.FormatInt(post.ID, 10), target.ID)
	assert.Equal(t, post.Text, target.Text)
	assert.Equal(t, strconv.FormatInt(post.User, 10), target.UserID)
	assert.Equal(t, post.CommentsOFF, target.CommentsOff)
}

func TestConvertNewCommentModelIntoEntity(t *testing.T) {
	newComment := model.NewComment{Text: "comment 1", UserID: "1", PostID: "1"}

	target := convertNewCommentModelIntoEntity(newComment)
	assert.Equal(t, newComment.Text, target.Text)
	assert.Equal(t, newComment.UserID, strconv.FormatInt(target.UserID, 10))
	assert.Equal(t, newComment.PostID, strconv.FormatInt(target.PostID, 10))
}

func TestConvertCommentEntityIntoModel(t *testing.T) {
	comment := entity.Comment{Text: "comment 1", UserID: int64(1), PostID: int64(1)}

	target := convertCommentEntityIntoModel(comment)
	assert.Equal(t, comment.Text, target.Text)
	assert.Equal(t, strconv.FormatInt(comment.PostID, 10), target.PostID)
	assert.Equal(t, strconv.FormatInt(comment.UserID, 10), target.UserID)
}
