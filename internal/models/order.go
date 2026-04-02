package models

type OrderItem struct {
	ID       int     `json:"id"`
	OrderID  int     `json:"order_id"`
	MealID   int     `json:"meal_id"`
	Name     string  `json:"name"`
	CategoryName string `json:"category_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}


type CreateOrderRequest struct {
	Items   []OrderItemRequest `json:"items"`
	Address string             `json:"address"`
	Phone   string             `json:"phone"`
}

type OrderItemRequest struct {
	MealID   int `json:"meal_id"`
	Quantity int `json:"quantity"`
}

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
	Address    string    `json:"address"`
	Phone      string    `json:"phone"`
	CreatedAt  string    `json:"created_at"`
	Items      []OrderItem `json:"items"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status"`
}