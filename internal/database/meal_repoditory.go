package database

import (
	"context"
	"github.com/akram-fattah/food-delivery/internal/models"
)



func CreateMeal(ctx context.Context, meal *models.Meal) error {
	query := `
		INSERT INTO meals (name, description, price, image_url, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, is_available
	`

	return DB.QueryRowContext(
		ctx,
		query,
		meal.Name,
		meal.Description,
		meal.Price,
		meal.ImageURL,
		meal.CategoryID,
	).Scan(
		&meal.ID,
		&meal.CreatedAt,
		&meal.IsAvailable,
	)
}