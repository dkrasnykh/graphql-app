package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

func (s *Service) ValidateComment(input model.NewComment) (*entity.Comment, error) {
	var errList []error
	if len(input.Text) == 0 {
		errList = append(errList, ErrEmptyBody)
	}
	if len([]rune(input.Text)) > 2000 {
		errList = append(errList, ErrCommentBodyTooBig)
	}
	if _, err := strconv.ParseInt(input.UserID, 10, 64); err != nil {
		errList = append(errList, fmt.Errorf("%w, user id: %s", ErrInvalidID, input.UserID))
	}
	if _, err := strconv.ParseInt(input.PostID, 10, 64); err != nil {
		errList = append(errList, fmt.Errorf("%w, post id: %s", ErrInvalidID, input.PostID))
	}
	if input.ParentCommentID != nil {
		_, err := strconv.ParseInt(*input.ParentCommentID, 10, 64)
		if err != nil {
			errList = append(errList, fmt.Errorf("%w, parent comment id: %s", ErrInvalidID, *input.ParentCommentID))
		}
	}
	if len(errList) > 0 {
		return nil, errors.Join(errList...)
	}

	return convertNewCommentModelIntoEntity(input), nil
}

func (s *Service) SaveComment(ctx context.Context, comment entity.Comment) (*model.Comment, error) {
	id, err := s.storage.SaveComment(ctx, comment)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrPostNotFound):
			return nil, fmt.Errorf("%w; post id: %d", ErrPostNotFound, comment.PostID)
		case errors.Is(err, storage.ErrPostCommentsDisabled):
			return nil, fmt.Errorf("%w; post id: %d", ErrPostCommentsDisabled, comment.PostID)
		case errors.Is(err, storage.ErrInvalidParentCommentID):
			return nil, fmt.Errorf("%w; post id: %d; parent comment id: %d", ErrInvalidParentCommentID, comment.PostID, *comment.ParentCommentID)
		case errors.Is(err, storage.ErrParentCommentBelongAnotherPost):
			return nil, fmt.Errorf("%w; current post id: %d", ErrParentCommentBelongAnotherPost, comment.PostID)
		default:
			return nil, ErrInternal
		}
	}
	comment.ID = id
	target := convertCommentEntityIntoModel(comment)
	s.subscriptions.Broadcast(comment.PostID, target)
	return target, nil
}

func (s *Service) AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*model.Comment, error) {
	list, err := s.storage.AllComments(ctx, postID, limit, offset)
	if err != nil {
		return nil, ErrInternal
	}

	all := make([]*model.Comment, len(list))
	for i, c := range list {
		all[i] = convertCommentEntityIntoModel(*c)
	}
	return all, nil
}
