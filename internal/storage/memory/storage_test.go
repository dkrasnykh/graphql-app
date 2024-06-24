package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dkrasnykh/graphql-app/internal/entity"
)

type Storager interface {
	SavePost(ctx context.Context, post entity.Post) (int64, error)
	PostByID(ctx context.Context, id int64) (*entity.Post, error)
	AllPosts(ctx context.Context) ([]*entity.Post, error)
	DisableComments(ctx context.Context, userID int64, postID int64) error

	SaveComment(ctx context.Context, comment entity.Comment) (int64, error)
	AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error)
}

type testStorager interface {
	Storager
	clean(ctx context.Context)
}

type StoragerTestSuite struct {
	suite.Suite
	testStorager
}

func (ts *StoragerTestSuite) SetupSuite() {
	ts.testStorager = New()
}

func (s *StorageMemory) clean(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.IDComment = 1
	s.IDPost = 1
	s.PostAdjList = make(map[int64]map[int64][]int64)
	s.IDValuePostMap = make(map[int64]entity.Post)
	s.IDValueCommentMap = make(map[int64]entity.Comment)
	s.PostRootComments = make(map[int64][]int64)
}

func (ts *StoragerTestSuite) TearDownSuite() {

}

func TestStoragePostgres(t *testing.T) {
	suite.Run(t, new(StoragerTestSuite))
}

func (ts *StoragerTestSuite) SetupTest() {
	ts.clean(context.Background())
}

func (ts *StoragerTestSuite) TearDownTest() {
	ts.clean(context.Background())
}
