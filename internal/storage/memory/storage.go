package memory

import (
	"sync"

	"github.com/dkrasnykh/graphql-app/internal/entity"
)

// all structures are under one mutex, because a possible case:
// when one goroutine adds a comment for post, and another disables comments for this post
type StorageMemory struct {
	mu                sync.RWMutex
	PostCounter       int64
	CommentCounter    int64
	IDValuePostMap    map[int64]entity.Post
	IDValueCommentMap map[int64]entity.Comment
	// for each comments store root comments
	PostRootComments map[int64][]int64
	// for each post store comments adjacency list
	PostAdjList map[int64]map[int64][]int64
}

func New() *StorageMemory {
	return &StorageMemory{
		mu:                sync.RWMutex{},
		PostCounter:       1,
		CommentCounter:    1,
		IDValuePostMap:    make(map[int64]entity.Post),
		IDValueCommentMap: make(map[int64]entity.Comment),
		PostRootComments:  make(map[int64][]int64),
		PostAdjList:       make(map[int64]map[int64][]int64),
	}
}

// used only in unit tests
func (s *StorageMemory) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CommentCounter = 1
	s.PostCounter = 1
	s.PostAdjList = make(map[int64]map[int64][]int64)
	s.IDValuePostMap = make(map[int64]entity.Post)
	s.IDValueCommentMap = make(map[int64]entity.Comment)
	s.PostRootComments = make(map[int64][]int64)
}
