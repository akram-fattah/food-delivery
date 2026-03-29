package database

import (
	"context"
)

func UpdateRefreshToken(ctx context.Context, userID int64, newToken string) error {
	query := `
  UPDATE refresh_tokens
  SET token = $1,
      expires_at = NOW() + INTERVAL '7 days'
  WHERE user_id = $2`

	_, err := DB.ExecContext(ctx, query, newToken, userID)
	return err
}
