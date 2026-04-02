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

func IsUsernameExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	err := DB.QueryRowContext(ctx, query, username).Scan(&exists)
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

func SetUserVerificationCode(ctx context.Context, email, code string, expires time.Time) error {
	query := `
		UPDATE users
		SET verification_code = $1, verification_expires = $2
		WHERE email = $3
	`
	_, err := DB.ExecContext(ctx, query, code, expires, email)
	return err
}

func VerifyResetCode(ctx context.Context, email, code string) (bool, error) {
	var expires time.Time
	query := `SELECT verification_expires FROM users WHERE email = $1 AND verification_code = $2`
	err := DB.QueryRowContext(ctx, query, email, code).Scan(&expires)
	if err != nil {
		return false, nil 
	}
	if time.Now().After(expires) {
		return false, nil
	}
	return true, nil
}

func UpdateUserPasswordAndClearCode(ctx context.Context, email, hashedPassword string) error {
	query := `UPDATE users SET password = $1, verification_code = NULL, verification_expires = NULL WHERE email = $2`
	_, err := DB.ExecContext(ctx, query, hashedPassword, email)
	return err
}

func GetProfile(ctx context.Context, userID int) (models.User, error) {
	var u models.User
	query := `
		SELECT id, name, username, email, role, created_at
		FROM users
		WHERE id = $1
	`
	err := DB.QueryRowContext(ctx, query, userID).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Email,
		&u.Role,
		&u.CreateAt,
	)

	return u, err
}

func UpdateUser(ctx context.Context, userID int, req models.UpdateUserRequest) error {

	query := `
		UPDATE users SET
			name = COALESCE($1, name),
			username = COALESCE($2, username),
			email = COALESCE($3, email),
			password = COALESCE($4, password)
		WHERE id = $5
	`

	_, err := DB.ExecContext(
		ctx,
		query,
		req.Name,
		req.Username,
		req.Email,
		req.Password,
		userID,
	)

	return err
}

func GetUserEmailByID(userID int) (string, error) {
	var email string
	query := `SELECT email FROM users WHERE id = $1`
	err := DB.QueryRowContext(context.Background(), query, userID).Scan(&email)
	return email, err
}