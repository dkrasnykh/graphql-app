-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts
(
    id                   BIGSERIAL     NOT NULL UNIQUE,
    text                 VARCHAR(2000) NOT NULL,
    user_id              BIGINT        NOT NULL,
    is_comments_disabled BOOLEAN
);

CREATE TABLE IF NOT EXISTS comments
(
    id                BIGSERIAL        NOT NULL UNIQUE,
    text              VARCHAR(2000)   NOT NULL,
    user_id           BIGINT        NOT NULL,
    post_id           BIGINT        NOT NULL,
    parent_comment_id BIGINT,
    rank              VARCHAR
);

-- +goose Down
DROP TABLE posts;
DROP TABLE comments;