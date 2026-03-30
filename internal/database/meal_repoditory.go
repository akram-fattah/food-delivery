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

func GetMeals(ctx context.Context) ([]models.Meal, error) {
	query := `
		SELECT 
		meals.id,
		meals.name,
		meals.description,
		meals.price,
		meals.image_url,
		meals.category_id,
		meals.created_at,
		meals.is_available,
		categories.category_name
		FROM meals
		JOIN categories ON meals.category_id = categories.id
	`

	rows, err := DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		err := rows.Scan(
			&meal.ID,
			&meal.Name,
			&meal.Description,
			&meal.Price,
			&meal.ImageURL,
			&meal.CategoryID,
			&meal.CreatedAt,
			&meal.IsAvailable,
			&meal.CategoryName,
		)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}

	return meals, nil
}
