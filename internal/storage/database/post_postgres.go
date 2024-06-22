package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/dkrasnykh/graphql-app/internal/entity"
	"github.com/dkrasnykh/graphql-app/internal/storage"
)

func (s *StoragePostgres) SavePost(ctx context.Context, post entity.Post) (int64, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var id int64
	row := s.db.QueryRow(newCtx, "INSERT INTO posts (text, user_id, is_comments_disabled) values ($1, $2, $3) RETURNING id",
		post.Text, post.User, post.CommentsOFF)
	err := row.Scan(&id)
	if err != nil {
		return 0, storage.ErrInternal
	}

	return id, nil
}

func (s *StoragePostgres) PostByID(ctx context.Context, id int64) (*entity.Post, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rows, err := s.db.Query(newCtx, "select (id, text, user_id, is_comments_disabled) from posts where id = $1", id)
	if err != nil {
		return nil, storage.ErrInternal
	}

	post, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[entity.Post])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrPostNotFound
		}

		return nil, storage.ErrInternal
	}
	return &post, nil
}

func (s *StoragePostgres) AllPosts(ctx context.Context) ([]*entity.Post, error) {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	rows, err := s.db.Query(newCtx, "select (id, text, user_id, is_comments_disabled) from posts")
	if err != nil {
		return nil, storage.ErrInternal
	}

	list, err := pgx.CollectRows(rows, pgx.RowTo[*entity.Post])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrInternal
	}

	return list, nil
}

// TODO move validation logic into service layer
// all checks must be performed in one transaction
// TODO design how to begin/commit transaction into service layer (lock / unlock for memory storage)
func (s *StoragePostgres) DisableComments(ctx context.Context, userID int64, postID int64) error {
	newCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.db.Begin(newCtx)
	if err != nil {
		return storage.ErrInternal
	}

	var currUserID int64
	var disabled bool
	row := tx.QueryRow(newCtx, "select user_id, is_comments_disabled from posts where id = $1", postID)
	err = row.Scan(&currUserID, &disabled)
	if err != nil {
		_ = tx.Rollback(newCtx)
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrPostNotFound
		}

		return storage.ErrInternal
	}
	if currUserID != userID {
		_ = tx.Rollback(newCtx)
		return fmt.Errorf("%w, keeper ID:%d", storage.ErrAccess, currUserID)
	}
	if disabled {
		_ = tx.Rollback(newCtx)
		return storage.ErrPostCommentsDisabled
	}
	_, err = tx.Exec(newCtx, "UPDATE posts SET is_comments_disabled = true WHERE id = $1", postID)
	if err != nil {
		_ = tx.Rollback(newCtx)
		return storage.ErrInternal
	}

	err = tx.Commit(newCtx)
	if err != nil {
		return storage.ErrInternal
	}

	return nil
}
