package graph

import (
	"context"

	"github.com/dkrasnykh/graphql-app/graph/model"
	"github.com/dkrasnykh/graphql-app/internal/service"
)

/*
EXAMPLE for TestComments_OK

"awesome post 1" (id: 1)
|
|-> "comment 1" (id: 1, parent_comment_id: nil)
|	|
|	|-> "comment 3" (id: 3, parent_comment_id: 1)
|	|	|
|	|	|-> "comment 5" (id: 5, parent_comment_id: 3)
|	|
|	|-> "comment 4" (id: 4, parent_comment_id: 1)
|	|	|
|	|	|-> "comment 6" (id: 6, parent_comment_id: 4)
|
|-> "comment 2" (id: 2, parent_comment_id: nil)
|
"awesome post 2" (id: 2)
|
|-> "comment 7" (id: 7, parent_comment_id: nil)

correct result list order for pagination:
1. "awesome post 1" (id: 1)
["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"]
2. "awesome post 2" (id: 2)
["comment 7"]
*/
func (ts *ResolverTestSuite) TestComments_OK() {
	ctx := context.Background()

	inputPost1 := model.NewPost{Text: "awesome post 1", UserID: "1"}
	inputPost2 := model.NewPost{Text: "awesome post 2", UserID: "1"}

	post1, err := ts.mutation.CreatePost(ctx, inputPost1)
	ts.NoError(err)
	post2, err := ts.mutation.CreatePost(ctx, inputPost2)
	ts.NoError(err)

	inputComment1 := model.NewComment{Text: "comment 1", UserID: "1", PostID: post1.ID}
	comment1, err := ts.mutation.CreateComment(ctx, inputComment1)
	ts.NoError(err)

	inputComment2 := model.NewComment{Text: "comment 2", UserID: "1", PostID: post1.ID}
	comment2, err := ts.mutation.CreateComment(ctx, inputComment2)
	ts.NoError(err)

	inputComment3 := model.NewComment{Text: "comment 3", ParentCommentID: &comment1.ID, UserID: "1", PostID: post1.ID}
	comment3, err := ts.mutation.CreateComment(ctx, inputComment3)
	ts.NoError(err)

	inputComment4 := model.NewComment{Text: "comment 4", ParentCommentID: &comment1.ID, UserID: "1", PostID: post1.ID}
	comment4, err := ts.mutation.CreateComment(ctx, inputComment4)
	ts.NoError(err)

	inputComment5 := model.NewComment{Text: "comment 5", ParentCommentID: &comment3.ID, UserID: "1", PostID: post1.ID}
	comment5, err := ts.mutation.CreateComment(ctx, inputComment5)
	ts.NoError(err)

	inputComment6 := model.NewComment{Text: "comment 6", ParentCommentID: &comment4.ID, UserID: "1", PostID: post1.ID}
	comment6, err := ts.mutation.CreateComment(ctx, inputComment6)
	ts.NoError(err)

	inputComment7 := model.NewComment{Text: "comment 7", UserID: "1", PostID: post2.ID}
	comment7, err := ts.mutation.CreateComment(ctx, inputComment7)
	ts.NoError(err)

	limit, offset := 10, 0
	list, err := ts.query.Comments(ctx, post1.ID, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"]
	ts.Equal(6, len(list))
	ts.Equal(comment1, list[0])
	ts.Equal(comment3, list[1])
	ts.Equal(comment5, list[2])
	ts.Equal(comment4, list[3])
	ts.Equal(comment6, list[4])
	ts.Equal(comment2, list[5])
	// ["comment 7"]
	list, err = ts.query.Comments(ctx, post2.ID, &limit, &offset)
	ts.NoError(err)
	ts.Equal(1, len(list))
	ts.Equal(comment7, list[0])

	limit, offset = 2, 1
	list, err = ts.query.Comments(ctx, post1.ID, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> ["comment 3", "comment 5"]
	ts.Equal(2, len(list))
	ts.Equal(comment3, list[0])
	ts.Equal(comment5, list[1])

	limit, offset = 10, 4
	list, err = ts.query.Comments(ctx, post1.ID, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> ["comment 6", "comment 2"]
	ts.Equal(2, len(list))
	ts.Equal(comment6, list[0])
	ts.Equal(comment2, list[1])

	limit, offset = 10, 10
	list, err = ts.query.Comments(ctx, post1.ID, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> []
	ts.Equal(0, len(list))
}

func (ts *ResolverTestSuite) TestPost_InvalidID() {
	_, err := ts.query.Post(context.Background(), "text")
	ts.ErrorIs(err, service.ErrInvalidID)
}

func (ts *ResolverTestSuite) TestPost_PostNotFound() {
	_, err := ts.query.Post(context.Background(), "1")
	ts.ErrorIs(err, service.ErrPostNotFound)
}

func (ts *ResolverTestSuite) TestPost_OK() {
	input := model.NewPost{Text: "awesome post", UserID: "1"}
	created, err := ts.mutation.CreatePost(context.Background(), input)
	ts.NoError(err)

	post, err := ts.query.Post(context.Background(), created.ID)
	ts.NoError(err)
	ts.Equal(input.Text, post.Text)
	ts.Equal(input.UserID, post.UserID)
}

func (ts *ResolverTestSuite) TestPosts_OK() {
	inputPost1 := model.NewPost{Text: "awesome post 1", UserID: "1"}
	_, err := ts.mutation.CreatePost(context.Background(), inputPost1)
	ts.NoError(err)
	inputPost2 := model.NewPost{Text: "awesome post 2", UserID: "1"}
	_, err = ts.mutation.CreatePost(context.Background(), inputPost2)
	ts.NoError(err)

	list, err := ts.query.Posts(context.Background())
	ts.NoError(err)
	ts.Equal(2, len(list))
}

func (ts *ResolverTestSuite) TestPosts_EmptyResult() {
	list, err := ts.query.Posts(context.Background())
	ts.NoError(err)
	ts.Equal(0, len(list))
}
