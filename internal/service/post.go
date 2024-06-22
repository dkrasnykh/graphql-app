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

func (s *Service) ValidatePost(input model.NewPost) (*entity.Post, error) {
	var errList []error
	if len(input.Text) == 0 {
		errList = append(errList, ErrEmptyBody)
	}
	_, err := strconv.ParseInt(input.UserID, 10, 64)
	if err != nil {
		errList = append(errList, fmt.Errorf("%w; user id: %s", ErrInvalidID, input.UserID))
	}
	if len(errList) > 0 {
		return nil, errors.Join(errList...)
	}

	return convertNewPostModelIntoEntity(input), nil
}

func (s *Service) SavePost(ctx context.Context, post entity.Post) (*model.Post, error) {
	postID, err := s.storage.SavePost(ctx, post)
	if err != nil {
		return nil, ErrInternal
	}

	post.ID = postID
	return convertPostEntityIntoModel(post), nil
}

func (s *Service) ValidateDisableCommentsRequest(input model.DisableCommentsRequest) (userID int64, postID int64, err error) {
	var errList []error
	if userID, err = strconv.ParseInt(input.UserID, 10, 64); err != nil {
		errList = append(errList, fmt.Errorf("%w, user id: %s", ErrInvalidID, input.UserID))
	}
	if postID, err = strconv.ParseInt(input.PostID, 10, 64); err != nil {
		errList = append(errList, fmt.Errorf("%w, post id: %s", ErrInvalidID, input.PostID))
	}
	err = errors.Join(errList...)
	return userID, postID, err
}

func (s *Service) DisableComments(ctx context.Context, userID int64, postID int64) error {
	if err := s.storage.DisableComments(ctx, userID, postID); err != nil {
		switch {
		case errors.Is(err, storage.ErrPostNotFound):
			return fmt.Errorf("%w; post id: %d", err, postID)
		case errors.Is(err, storage.ErrAccess):
			return fmt.Errorf("%w; userID: %d; postID: %d", err, userID, postID)
		case errors.Is(err, storage.ErrPostCommentsDisabled):
			return fmt.Errorf("comments already turned off; post id: %d", postID)
		default:
			return ErrInternal
		}
	}
	return nil
}

func (s *Service) PostById(ctx context.Context, ID int64) (*model.Post, error) {
	post, err := s.storage.PostByID(ctx, ID)
	if err != nil {
		return nil, ErrInternal
	}

	return convertPostEntityIntoModel(*post), nil
}

func (s *Service) AllPosts(ctx context.Context) ([]*model.Post, error) {
	list, err := s.storage.AllPosts(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	all := make([]*model.Post, len(list))
	for i, post := range list {
		all[i] = convertPostEntityIntoModel(*post)
	}
	return all, nil
}

func (s *Service) ValidateID(ID string) (int64, error) {
	var id int64
	var err error
	if id, err = strconv.ParseInt(ID, 10, 64); err != nil {
		return 0, ErrInvalidID
	}

	return id, nil
}
