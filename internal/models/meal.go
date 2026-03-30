package models

import "time"

type Meal struct {
	ID           int        `json:"id"`
	Name         *string    `json:"name,omitempty"`
	Description  *string    `json:"description,omitempty"`
	Price        *float64   `json:"price,omitempty"`
	ImageURL     *string    `json:"image_url,omitempty"`
	IsAvailable  *bool      `json:"is_available,omitempty"`
	CategoryID   *int       `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	CategoryName *string    `json:"category_name,omitempty"`
}