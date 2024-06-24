package memory

import (
	"context"
	"fmt"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

func (s *StorageMemory) SaveComment(ctx context.Context, comment entity.Comment) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// post exists check
	post, ok := s.IDValuePostMap[comment.PostID]
	if !ok {
		return 0, storage.ErrPostNotFound
	}
	// comments are enabled check
	if post.CommentsOFF {
		return 0, storage.ErrPostCommentsDisabled
	}

	if comment.ParentCommentID != nil {
		// parent comment exists check
		parentComment, ok := s.IDValueCommentMap[*comment.ParentCommentID]
		if !ok {
			return 0, storage.ErrInvalidParentCommentID
		}
		// parent comment PostID check
		if parentComment.PostID != comment.PostID {
			return 0, fmt.Errorf("%w; parent comment post id: %d", storage.ErrParentCommentBelongAnotherPost, parentComment.PostID)
		}
	}

	// insert new comment
	id := s.CommentCounter
	comment.ID = id
	s.IDValueCommentMap[id] = comment

	if comment.ParentCommentID == nil {
		// root comment
		s.PostRootComments[comment.PostID] = append(s.PostRootComments[comment.PostID], id)
	} else {
		s.PostAdjList[comment.PostID][*comment.ParentCommentID] = append(s.PostAdjList[comment.PostID][*comment.ParentCommentID], id)
	}

	s.CommentCounter += 1

	return id, nil
}

// default values limit = 10, offset = 0 (graphql schema)
func (s *StorageMemory) AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commentsList := make([]*entity.Comment, 0)
	i := 0

	var dfs func(v int64)
	dfs = func(v int64) {
		if *offset <= i && i < *offset+*limit {
			comm := s.IDValueCommentMap[v]
			commentsList = append(commentsList, &comm)
		}
		i += 1
		for _, u := range s.PostAdjList[postID][v] {
			dfs(u)
		}
	}

	for _, root := range s.PostRootComments[postID] {
		dfs(root)
		if len(commentsList) == *limit {
			break
		}
	}

	return commentsList, nil
}
