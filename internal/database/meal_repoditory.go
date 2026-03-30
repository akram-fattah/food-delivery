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


func GetMealByID(ctx context.Context, id int) (*models.Meal, error) {
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
		WHERE meals.id = $1
	`

	var meal models.Meal
	err := DB.QueryRowContext(ctx, query, id).Scan(
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

	return &meal, nil
}

func DeleteMeal(ctx context.Context, id int) error {
	query := `DELETE FROM meals WHERE id = $1`
	_, err := DB.ExecContext(ctx, query, id)
	return err
}


func UpdateMeal(ctx context.Context, meal *models.Meal) error {
	query := `
		UPDATE meals
		SET 
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			price = COALESCE($3, price),
			image_url = COALESCE($4, image_url),
			category_id = COALESCE($5, category_id),
			is_available = COALESCE($6, is_available)
		WHERE id = $7
	`

	_, err := DB.ExecContext(ctx, query,
		meal.Name,
		meal.Description,
		meal.Price,
		meal.ImageURL,
		meal.CategoryID,
		meal.IsAvailable,
		meal.ID,
	)

	return err
}

func GetMealsByCategoryID(ctx context.Context, categoryID int) ([]models.Meal, error) {
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
		WHERE meals.category_id = $1
	`

	rows, err := DB.QueryContext(ctx, query, categoryID)
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

func SearchMeals(ctx context.Context, keyword string) ([]models.Meal, error) {
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
		WHERE meals.name ILIKE $1 OR meals.description ILIKE $1
	`

	rows, err := DB.QueryContext(ctx, query, "%"+keyword+"%")
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
