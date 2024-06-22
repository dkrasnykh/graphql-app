package memory

import (
	"context"
	"math/rand"

	"github.com/stretchr/testify/require"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

func (ts *StoragerTestSuite) TestSavePost_OK() {
	post := entity.Post{Text: "awesome post", User: int64(3)}

	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	saved, err := ts.PostByID(context.Background(), postID)
	ts.NoError(err)
	ts.Equal(post.Text, saved.Text)
	ts.Equal(post.User, saved.User)
}

func (ts *StoragerTestSuite) TestPostByID_PostNotFound() {
	randomID := rand.Int63()
	_, err := ts.PostByID(context.Background(), randomID)
	require.ErrorIs(ts.T(), err, storage.ErrPostNotFound)
}

func (ts *StoragerTestSuite) TestAllPosts_EmptyResult() {
	list, err := ts.AllPosts(context.Background())
	ts.NoError(err)
	ts.Equal(0, len(list))
}

func (ts *StoragerTestSuite) TestAllPosts_OK() {
	post1 := entity.Post{Text: "awesome post", User: int64(3)}
	_, err := ts.SavePost(context.Background(), post1)
	ts.NoError(err)
	post2 := entity.Post{Text: "awesome post 1", User: int64(4)}
	_, err = ts.SavePost(context.Background(), post2)
	ts.NoError(err)

	list, err := ts.AllPosts(context.Background())
	ts.NoError(err)
	ts.Equal(2, len(list))
}

func (ts *StoragerTestSuite) TestDisableComments_PostNotFound() {
	userID := rand.Int63()
	postID := rand.Int63()
	err := ts.DisableComments(context.Background(), userID, postID)
	ts.ErrorIs(err, storage.ErrPostNotFound)
}

func (ts *StoragerTestSuite) TestDisableComments_PostBelongAnotherUser() {
	userID := int64(1)
	post := entity.Post{Text: "awesome post", User: userID}
	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	anotherUserID := int64(2)
	err = ts.DisableComments(context.Background(), anotherUserID, postID)
	ts.ErrorIs(err, storage.ErrAccess)
}

func (ts *StoragerTestSuite) TestDisableComments_PostAlreadyDisabled() {
	userID := rand.Int63()
	post := entity.Post{Text: "awesome post", User: userID, CommentsOFF: true}
	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	err = ts.DisableComments(context.Background(), userID, postID)
	ts.ErrorIs(err, storage.ErrPostCommentsDisabled)
}

func (ts *StoragerTestSuite) TestDisableComments_OK() {
	userID := rand.Int63()
	post := entity.Post{Text: "awesome post", User: userID}
	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	err = ts.DisableComments(context.Background(), userID, postID)
	ts.NoError(err)
}
