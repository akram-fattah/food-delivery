package database

import (
	"context"
	"time"

	"github.com/akram-fattah/food-delivery/internal/models"
)

func CreateUser(ctx context.Context, u *models.User, code string, expires time.Time) error {

	query := `
	INSERT INTO users (username, name, email, password, verification_code, verification_expires, is_verified)
	VALUES ($1, $2, $3, $4, $5, $6, false)
	`

	_, err := DB.ExecContext(
		ctx,
		query,
		u.Username,
		u.Name,
		u.Email,
		u.Password,
		code,
		expires,
	)

	return err
}

func IsEmailExists(ctx context.Context, email string) (bool, error) {

	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`

	err := DB.QueryRowContext(ctx, query, email).Scan(&exists)

	return exists, err
}

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User

	query := `
		SELECT id, name, email, password, is_verified
		FROM users
		WHERE email=$1
	`

	err := DB.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsVerified)

	return u, err
}

func VerifyUserEmail(ctx context.Context, email, code string) (bool, error) {

	query := `
		UPDATE users
		SET is_verified = true,
		    verification_code = NULL
		WHERE email = $1
		  AND verification_code = $2
		  AND is_verified = false
	`

	result, err := DB.ExecContext(ctx, query, email, code)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}
