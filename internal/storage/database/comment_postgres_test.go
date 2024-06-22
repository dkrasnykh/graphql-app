package database

import (
	"context"
	"math/rand"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

/*
EXAMPLE for TestAllComments_OK

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
func (ts *StoragerTestSuite) TestAllComments_OK() {
	ctx := context.Background()

	userID := rand.Int63()

	post1 := entity.Post{Text: "awesome post 1", User: userID}
	post2 := entity.Post{Text: "awesome post 2", User: userID}

	postID1, err := ts.SavePost(ctx, post1)
	ts.NoError(err)
	postID2, err := ts.SavePost(ctx, post2)
	ts.NoError(err)

	comment1 := entity.Comment{Text: "comment 1", UserID: userID, PostID: postID1}
	commentID1, err := ts.SaveComment(ctx, comment1)
	ts.NoError(err)

	comment2 := entity.Comment{Text: "comment 2", UserID: userID, PostID: postID1}
	commentID2, err := ts.SaveComment(ctx, comment2)
	ts.NoError(err)

	comment3 := entity.Comment{Text: "comment 3", ParentCommentID: &commentID1, UserID: userID, PostID: postID1}
	commentID3, err := ts.SaveComment(ctx, comment3)
	ts.NoError(err)

	comment4 := entity.Comment{Text: "comment 4", ParentCommentID: &commentID1, UserID: userID, PostID: postID1}
	commentID4, err := ts.SaveComment(ctx, comment4)
	ts.NoError(err)

	comment5 := entity.Comment{Text: "comment 5", ParentCommentID: &commentID3, UserID: userID, PostID: postID1}
	commentID5, err := ts.SaveComment(ctx, comment5)
	ts.NoError(err)

	comment6 := entity.Comment{Text: "comment 6", ParentCommentID: &commentID4, UserID: userID, PostID: postID1}
	commentID6, err := ts.SaveComment(ctx, comment6)
	ts.NoError(err)

	comment7 := entity.Comment{Text: "comment 7", UserID: userID, PostID: postID2}
	commentID7, err := ts.SaveComment(ctx, comment7)
	ts.NoError(err)

	comment1.ID = commentID1
	comment2.ID = commentID2
	comment3.ID = commentID3
	comment4.ID = commentID4
	comment5.ID = commentID5
	comment6.ID = commentID6
	comment7.ID = commentID7

	limit, offset := 10, 0
	list, err := ts.AllComments(ctx, postID1, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"]
	ts.Equal(6, len(list))
	ts.Equal(comment1, *list[0])
	ts.Equal(comment3, *list[1])
	ts.Equal(comment5, *list[2])
	ts.Equal(comment4, *list[3])
	ts.Equal(comment6, *list[4])
	ts.Equal(comment2, *list[5])
	// ["comment 7"]
	list, err = ts.AllComments(ctx, postID2, &limit, &offset)
	ts.NoError(err)
	ts.Equal(1, len(list))
	ts.Equal(comment7, *list[0])

	limit, offset = 2, 1
	list, err = ts.AllComments(ctx, postID1, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> ["comment 3", "comment 5"]
	ts.Equal(2, len(list))
	ts.Equal(comment3, *list[0])
	ts.Equal(comment5, *list[1])

	limit, offset = 10, 4
	list, err = ts.AllComments(ctx, postID1, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> ["comment 6", "comment 2"]
	ts.Equal(2, len(list))
	ts.Equal(comment6, *list[0])
	ts.Equal(comment2, *list[1])

	limit, offset = 10, 10
	list, err = ts.AllComments(ctx, postID1, &limit, &offset)
	ts.NoError(err)
	// ["comment 1", "comment 3", "comment 5", "comment 4", "comment 6", "comment 2"] -> []
	ts.Equal(0, len(list))
}

func (ts *StoragerTestSuite) TestAllComments_EmptyResult() {
	//userID := rand.Int63()
	postID := rand.Int63()
	limit, offset := 10, 0
	list, err := ts.AllComments(context.Background(), postID, &limit, &offset)
	ts.NoError(err)
	ts.Equal(0, len(list))
}

func (ts *StoragerTestSuite) TestSaveComment_PostNotFound() {
	userID := rand.Int63()
	postID := rand.Int63()

	comment := entity.Comment{Text: "comment", PostID: postID, UserID: userID}

	_, err := ts.SaveComment(context.Background(), comment)
	ts.ErrorIs(err, storage.ErrPostNotFound)
}

func (ts *StoragerTestSuite) TestSaveComment_PostCommentsDisabled() {
	userID := rand.Int63()
	postCommentsOFF := entity.Post{Text: "awesome post", User: userID, CommentsOFF: true}
	postID, err := ts.SavePost(context.Background(), postCommentsOFF)
	ts.NoError(err)

	comment := entity.Comment{Text: "comment", PostID: postID, UserID: int64(3)}

	_, err = ts.SaveComment(context.Background(), comment)
	ts.ErrorIs(err, storage.ErrPostCommentsDisabled)
}

func (ts *StoragerTestSuite) TestSaveComment_ParentCommentIDNotFound() {
	userID := rand.Int63()
	post := entity.Post{Text: "awesome post", User: userID}

	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	parentID := int64(10)
	comment := entity.Comment{Text: "comment", ParentCommentID: &parentID, PostID: postID, UserID: int64(3)}

	_, err = ts.SaveComment(context.Background(), comment)
	ts.ErrorIs(err, storage.ErrInvalidParentCommentID)
}

func (ts *StoragerTestSuite) TestSaveComment_ParentCommentBelongAnotherPost() {
	userID := rand.Int63()
	post1 := entity.Post{Text: "awesome post 1", User: userID}
	post2 := entity.Post{Text: "awesome post 2", User: userID}

	postID1, err := ts.SavePost(context.Background(), post1)
	ts.NoError(err)

	postID2, err := ts.SavePost(context.Background(), post2)
	ts.NoError(err)

	parentComment := entity.Comment{Text: "parent comment", PostID: postID1, UserID: int64(3)}
	parentCommentID, err := ts.SaveComment(context.Background(), parentComment)
	ts.NoError(err)

	comment := entity.Comment{Text: "parent comment", ParentCommentID: &parentCommentID, PostID: postID2, UserID: int64(5)}
	_, err = ts.SaveComment(context.Background(), comment)
	ts.ErrorIs(err, storage.ErrParentCommentBelongAnotherPost)
}

func (ts *StoragerTestSuite) TestSaveComment_OKRootComment() {
	userID := rand.Int63()
	post := entity.Post{Text: "awesome post", User: userID}

	postID, err := ts.SavePost(context.Background(), post)
	ts.NoError(err)

	comment := entity.Comment{Text: "comment", PostID: postID, UserID: int64(3)}
	_, err = ts.SaveComment(context.Background(), comment)
	ts.NoError(err)

	offset, limit := 0, 10
	list, err := ts.AllComments(context.Background(), postID, &limit, &offset)
	ts.NoError(err)
	ts.Equal(1, len(list))
	ts.Equal(comment.Text, list[0].Text)
	ts.Equal(comment.PostID, list[0].PostID)
	ts.Nil(list[0].ParentCommentID)
	ts.Equal(comment.UserID, list[0].UserID)
}
