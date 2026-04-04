package database

import (
	"context"
	"database/sql"
	"time"
)

func CreateCategory(ctx context.Context, categoryName string) (int, error) {
	query := `
  INSERT INTO categories (category_name)
  VALUES ($1)
  RETURNING id`
	var id int
	err := DB.QueryRowContext(ctx, query, categoryName).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetCategories(ctx context.Context) ([]map[string]interface{}, error) {
	query := `SELECT id, category_name, created_at FROM categories`
	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var id int
		var categoryName string
		var createdAt time.Time
		if err := rows.Scan(&id, &categoryName, &createdAt); err != nil {
			return nil, err
		}
		categories = append(categories, map[string]interface{}{
			"id":            id,
			"category_name": categoryName,
			"created_at":    createdAt,
		})
	}
	return categories, nil
}

func GetCategoryByID(ctx context.Context, id int) (map[string]interface{}, error) {
	query := `SELECT id, category_name, created_at FROM categories WHERE id = $1`
	var categoryName string
	var createdAt time.Time
	if err := DB.QueryRowContext(ctx, query, id).Scan(&id, &categoryName, &createdAt); err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"category_name": categoryName,
		"created_at":    createdAt,
	}, nil
}


func UpdateCategory(ctx context.Context, id int, name string) (map[string]interface{}, error) {
	query := `
		UPDATE categories
		SET category_name = $1
		WHERE id = $2
		RETURNING id, category_name, created_at
	`

	var (
		categoryID   int
		categoryName string
		createdAt    time.Time
	)

	err := DB.QueryRowContext(ctx, query, name, id).Scan(
		&categoryID,
		&categoryName,
		&createdAt,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":            categoryID,
		"category_name": categoryName,
		"created_at":    createdAt,
	}, nil
}

func DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	res, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}


func CategoryExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)`
	var exists bool

	err := DB.QueryRowContext(ctx, query, id).Scan(&exists)
	return exists, err
}