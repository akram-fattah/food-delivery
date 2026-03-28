package database

import (
	"context"
	"database/sql"
	"time"
)

func SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := DB.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}

func GetRefreshToken(ctx context.Context, token string) (int, error) {
	var userID int
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := DB.QueryRowContext(ctx, query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		return 0, sql.ErrNoRows
	}
	return userID, nil
}

func DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := DB.ExecContext(ctx, query, token)
	return err
}
