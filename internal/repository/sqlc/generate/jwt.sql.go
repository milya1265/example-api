// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: jwt.sql

package repository

import (
	"context"
)

const createTokens = `-- name: CreateTokens :exec
INSERT INTO jwt (user_id, access_token, refresh_token) VALUES ($1, $2, $3)
`

type CreateTokensParams struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

func (q *Queries) CreateTokens(ctx context.Context, arg CreateTokensParams) error {
	_, err := q.db.ExecContext(ctx, createTokens, arg.UserID, arg.AccessToken, arg.RefreshToken)
	return err
}

const deleteTokens = `-- name: DeleteTokens :exec
DELETE FROM jwt WHERE user_id = $1
`

func (q *Queries) DeleteTokens(ctx context.Context, userID string) error {
	_, err := q.db.ExecContext(ctx, deleteTokens, userID)
	return err
}

const getTokens = `-- name: GetTokens :one
SELECT user_id, access_token, refresh_token FROM jwt WHERE user_id = $1
`

func (q *Queries) GetTokens(ctx context.Context, userID string) (Jwt, error) {
	row := q.db.QueryRowContext(ctx, getTokens, userID)
	var i Jwt
	err := row.Scan(&i.UserID, &i.AccessToken, &i.RefreshToken)
	return i, err
}

const updateAccessToken = `-- name: UpdateAccessToken :exec
UPDATE jwt SET access_token = $1 WHERE user_id = $2
`

type UpdateAccessTokenParams struct {
	AccessToken string
	UserID      string
}

func (q *Queries) UpdateAccessToken(ctx context.Context, arg UpdateAccessTokenParams) error {
	_, err := q.db.ExecContext(ctx, updateAccessToken, arg.AccessToken, arg.UserID)
	return err
}

const updateTokens = `-- name: UpdateTokens :exec
UPDATE jwt SET access_token = $1, refresh_token = $2 WHERE user_id = $3
`

type UpdateTokensParams struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

func (q *Queries) UpdateTokens(ctx context.Context, arg UpdateTokensParams) error {
	_, err := q.db.ExecContext(ctx, updateTokens, arg.AccessToken, arg.RefreshToken, arg.UserID)
	return err
}
