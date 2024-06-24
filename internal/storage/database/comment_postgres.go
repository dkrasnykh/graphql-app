package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

// TODO move validation logic into service layer
// all checks must be performed in one transaction
// TODO design how to begin/commit transaction into service layer (lock / unlock for memory storage)
func (s *StoragePostgres) SaveComment(ctx context.Context, comment entity.Comment) (int64, error) {
	const op = "Storage.postgresql.SaveComment"

	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(newCtx)
	if err != nil {
		return 0, storage.ErrInternal
	}

	// check that the post exists and comments are enabled
	var isDisabled bool
	row := tx.QueryRow(ctx, "SELECT is_comments_disabled FROM posts WHERE id = $1 FOR UPDATE", comment.PostID)
	err = row.Scan(&isDisabled)
	if err != nil {
		if err := tx.Rollback(newCtx); err != nil {
			slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrPostNotFound
		}
		return 0, storage.ErrInternal
	}

	if isDisabled {
		if err := tx.Rollback(newCtx); err != nil {
			slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
		}
		return 0, storage.ErrPostCommentsDisabled
	}

	// TODO ? save and compare parents IDS like array field in database
	rank := []byte{}

	if comment.ParentCommentID != nil {
		// extract parent rank (rank needed for pagination data sorting)
		// check that parent comment exists and belong the same post
		var parentPostID int64
		var parentRank string

		row := tx.QueryRow(ctx, "SELECT post_id, rank FROM comments WHERE id = $1 FOR UPDATE", *comment.ParentCommentID)
		err = row.Scan(&parentPostID, &parentRank)
		if err != nil {
			if err := tx.Rollback(newCtx); err != nil {
				slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
			}
			if errors.Is(err, pgx.ErrNoRows) {
				return 0, storage.ErrInvalidParentCommentID
			}
			return 0, storage.ErrInternal
		}
		if parentPostID != comment.PostID {
			if err := tx.Rollback(newCtx); err != nil {
				slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
			}
			return 0, fmt.Errorf("%w; parent comment post id: %d", storage.ErrParentCommentBelongAnotherPost, parentPostID)
		}

		rank = append(rank, []byte(parentRank)...)
		rank = append(rank, '-')
	}

	// insert new comment
	var id int64
	row = tx.QueryRow(context.Background(),
		"INSERT INTO comments (text, user_id, post_id, parent_comment_id) values ($1, $2, $3, $4) RETURNING id",
		comment.Text, comment.UserID, comment.PostID, comment.ParentCommentID)
	err = row.Scan(&id)
	if err != nil {
		if err := tx.Rollback(newCtx); err != nil {
			slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
		}
		return 0, storage.ErrInternal
	}

	// build rank (rank needed for pagination data sorting)
	// int64 max value = 9223372036854775807 - 19 characters -> subrank len == 19
	subRank := fmt.Sprintf("%019d", id)
	rank = append(rank, []byte(subRank)...)
	// update current comment
	_, err = tx.Exec(ctx, "UPDATE comments SET rank = $1 WHERE id = $2", string(rank), id)
	if err != nil {
		if err := tx.Rollback(newCtx); err != nil {
			slog.Log(newCtx, slog.LevelError, "%s %w transaction rollback error", op, err)
		}
		return 0, storage.ErrInternal
	}

	err = tx.Commit(newCtx)
	if err != nil {
		return 0, storage.ErrInternal
	}

	return id, nil
}

// default values limit = 10, offset = 0 (graphql schema)
func (s *StoragePostgres) AllComments(ctx context.Context, postID int64, limit *int, offset *int) ([]*entity.Comment, error) {
	const op = "Storage.postgresql.AllComments"

	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// TODO query can be more simply, if storing an additional field in the database (id root comment value)
	rows, err := s.db.Query(newCtx,
		`WITH RECURSIVE tmp(comment_id, parent_id, root) AS (
				SELECT t1.comment_id, t1.parent_id, t1.comment_id AS root
				FROM (SELECT id AS comment_id, parent_comment_id AS parent_id FROM comments WHERE parent_comment_id IS NULL AND post_id = $1) AS t1
    			UNION
    			SELECT t2.comment_id, t2.parent_id, tmp.root
    			FROM (SELECT id AS comment_id, parent_comment_id AS parent_id FROM comments WHERE post_id = $1) AS t2 JOIN tmp ON tmp.comment_id = t2.parent_id
			)
			SELECT c.id, c.text, c.user_id, c.post_id, c.parent_comment_id 
			FROM tmp LEFT JOIN comments AS c ON tmp.comment_id = c.id 
			ORDER BY tmp.root, c.rank OFFSET $2 LIMIT $3;`,
		postID, *offset, *limit)
	if err != nil {
		return nil, storage.ErrInternal
	}

	if rows.Err() != nil && !errors.Is(rows.Err(), pgx.ErrNoRows) {
		return nil, storage.ErrInternal
	}

	var list []*entity.Comment

	for rows.Next() {
		var c entity.Comment
		var parentCommentID sql.NullInt64
		err := rows.Scan(&c.ID, &c.Text, &c.UserID, &c.PostID, &parentCommentID)
		if err != nil {
			slog.Log(newCtx, slog.LevelError, "%s %w failed to parse selection row from database", op, err)
		}
		if parentCommentID.Valid {
			c.ParentCommentID = &parentCommentID.Int64
		}
		list = append(list, &c)
	}

	return list, nil
}
