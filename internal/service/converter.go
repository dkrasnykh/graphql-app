package service

import (
	"strconv"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/entity"
)

func convertNewPostModelIntoEntity(newPost model.NewPost) *entity.Post {
	var disabled bool
	if newPost.CommentsOff != nil && *newPost.CommentsOff {
		disabled = true
	}
	//parse error already checked (Validate method)
	userID, _ := strconv.ParseInt(newPost.UserID, 10, 64)
	return &entity.Post{
		Text:        newPost.Text,
		User:        userID,
		CommentsOFF: disabled,
	}
}

func convertPostEntityIntoModel(post entity.Post) *model.Post {
	return &model.Post{
		ID:          strconv.FormatInt(post.ID, 10),
		Text:        post.Text,
		UserID:      strconv.FormatInt(post.User, 10),
		CommentsOff: post.CommentsOFF,
	}
}

func convertNewCommentModelIntoEntity(newComment model.NewComment) *entity.Comment {
	var comment entity.Comment
	if newComment.ParentCommentID != nil {
		//parse error already checked (Validate method)
		parentCommentID, _ := strconv.ParseInt(*newComment.ParentCommentID, 10, 64)
		comment.ParentCommentID = &parentCommentID
	}
	//parse error already checked (Validate method)
	comment.UserID, _ = strconv.ParseInt(newComment.UserID, 10, 64)
	//parse error already checked (Validate method)
	comment.PostID, _ = strconv.ParseInt(newComment.PostID, 10, 64)
	comment.Text = newComment.Text
	return &comment
}

func convertCommentEntityIntoModel(comment entity.Comment) *model.Comment {
	var parentCommentID *string
	if comment.ParentCommentID != nil {
		value := strconv.FormatInt(*comment.ParentCommentID, 10)
		parentCommentID = &value
	}
	return &model.Comment{
		ID:              strconv.FormatInt(comment.ID, 10),
		Text:            comment.Text,
		ParentCommentID: parentCommentID,
		PostID:          strconv.FormatInt(comment.PostID, 10),
		UserID:          strconv.FormatInt(comment.UserID, 10),
	}
}
