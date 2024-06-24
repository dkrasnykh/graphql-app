package graph

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/dkrasnykh/graphql-app/internal/service"
	"github.com/dkrasnykh/graphql-app/internal/storage/memory"
	"github.com/dkrasnykh/graphql-app/internal/subscription"
)

type ResolverTestSuite struct {
	suite.Suite
	storage  service.Storager
	mutation MutationResolver
	query    QueryResolver
}

func (ts *ResolverTestSuite) SetupSuite() {
	ts.storage = memory.New()
	s := subscription.New()
	srv := service.New(ts.storage, s)
	resolver := Resolver{Service: srv}
	ts.mutation = resolver.Mutation()
	ts.query = resolver.Query()
}

func TestResolver(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

func (ts *ResolverTestSuite) SetupTest() {
	ts.storage.Clear()
}

func (ts *ResolverTestSuite) TearDownTest() {
	ts.storage.Clear()
}
