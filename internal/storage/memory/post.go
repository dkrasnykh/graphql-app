package memory

import (
	"context"
	"fmt"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

func (s *StorageMemory) SavePost(ctx context.Context, post entity.Post) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.IDPost
	post.ID = id
	s.IDValuePostMap[id] = post
	s.PostAdjList[id] = make(map[int64][]int64)
	s.IDPost += 1

	return id, nil
}

func (s *StorageMemory) PostByID(ctx context.Context, id int64) (*entity.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.IDValuePostMap[id]
	if !ok {
		return nil, storage.ErrPostNotFound
	}

	return &post, nil
}

func (s *StorageMemory) AllPosts(ctx context.Context) ([]*entity.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]*entity.Post, 0, len(s.IDValuePostMap))
	for _, v := range s.IDValuePostMap {
		post := v
		posts = append(posts, &post)
	}

	return posts, nil
}

func (s *StorageMemory) DisableComments(ctx context.Context, userID int64, postID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.IDValuePostMap[postID]; !ok {
		return storage.ErrPostNotFound
	}

	if s.IDValuePostMap[postID].User != userID {
		return fmt.Errorf("%w, keeper ID:%d", storage.ErrAccess, s.IDValuePostMap[postID].User)
	}

	// comments alredy disabled
	if s.IDValuePostMap[postID].CommentsOFF {
		return storage.ErrPostCommentsDisabled
	}

	post := s.IDValuePostMap[postID]
	post.CommentsOFF = true
	s.IDValuePostMap[postID] = post

	return nil
}
