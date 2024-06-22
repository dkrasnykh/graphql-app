package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/dkrasnykh/graphql-app/internal/entity"
)

// for running db tests locally, need to run db container first (docker-compose.yml - db service)
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
	clean(ctx context.Context) error
}

type StoragerTestSuite struct {
	suite.Suite
	testStorager

	tc *tcpostgres.PostgresContainer
}

func (ts *StoragerTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:latest"),
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.WithInitScripts(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(ts.T(), err)

	host, err := pgc.Host(ctx)
	require.NoError(ts.T(), err)

	port, err := pgc.MappedPort(ctx, "5432")
	require.NoError(ts.T(), err)

	ts.tc = pgc
	databaseURL := fmt.Sprintf("postgres://postgres:postgres@%s:%s/testdb?sslmode=disable", host, port.Port())

	err = Migrate(databaseURL)
	require.NoError(ts.T(), err)
	storage, err := New(databaseURL)

	require.NoError(ts.T(), err)

	ts.testStorager = storage

	ts.T().Logf("stared postgres at %s:%d", host, port.Int())
}

func (s *StoragePostgres) clean(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	_, err := s.db.Exec(newCtx, "DELETE FROM comments")
	if err != nil {
		return err
	}
	_, err = s.db.Exec(newCtx, "DELETE FROM posts")
	return err
}

func (ts *StoragerTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func TestStoragePostgres(t *testing.T) {
	suite.Run(t, new(StoragerTestSuite))
}

func (ts *StoragerTestSuite) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *StoragerTestSuite) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}
