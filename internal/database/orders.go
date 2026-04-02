package database

import (
	"context"
	"fmt"
	"github.com/akram-fattah/food-delivery/internal/models"
	"database/sql"
)


func CreateOrder(ctx context.Context, userID int, req models.CreateOrderRequest) (int, error) {

	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var totalPrice float64

	for _, item := range req.Items {
		var price float64

		err := tx.QueryRowContext(ctx,
			"SELECT price FROM meals WHERE id = $1",
			item.MealID,
		).Scan(&price)

		if err != nil {
			return 0, fmt.Errorf("meal not found: %d", item.MealID)
		}

		totalPrice += price * float64(item.Quantity)
	}

	var orderID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, total_price, address, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, totalPrice, req.Address, req.Phone).Scan(&orderID)

	if err != nil {
		return 0, err
	}

	for _, item := range req.Items {
		var price float64

		err := tx.QueryRowContext(ctx,
			"SELECT price FROM meals WHERE id = $1",
			item.MealID,
		).Scan(&price)

		if err != nil {
			return 0, err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, meal_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, orderID, item.MealID, item.Quantity, price)

		if err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}


func GetOrders(ctx context.Context, userID int, role string) ([]models.Order, error) {

	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = DB.QueryContext(ctx, `
			SELECT id, user_id, status, total_price, address, phone, created_at
			FROM orders
		`)
	} else {
		rows, err = DB.QueryContext(ctx, `
			SELECT id, user_id, status, total_price, address, phone, created_at
			FROM orders
			WHERE user_id = $1
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var o models.Order

		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Status,
			&o.TotalPrice,
			&o.Address,
			&o.Phone,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items, err := getOrderItems(ctx, o.ID)
		if err != nil {
			return nil, err
		}

		o.Items = items
		orders = append(orders, o)
	}

	return orders, nil
}


func getOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	
	rows, err := DB.QueryContext(ctx, `
		SELECT 
			oi.id,
			oi.order_id,
			oi.meal_id,
			m.name AS meal_name,
			c.category_name,
			oi.quantity,
			oi.price
		FROM order_items oi
		JOIN meals m ON m.id = oi.meal_id
		JOIN categories c ON c.id = m.category_id
		WHERE oi.order_id = $1
	`, orderID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MealID,
			&item.Name,
			&item.CategoryName,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func UpdateOrderStatus(ctx context.Context, orderID int, status string) error {

	_, err := DB.ExecContext(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`, status, orderID)
	if err != nil {
		return err
	}
	return nil
}
