package models

import "time"

type Meal struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	IsAvailable bool      `json:"is_available"`
	CategoryID  int       `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	CategoryName string    `json:"category_name,omitempty"`
}