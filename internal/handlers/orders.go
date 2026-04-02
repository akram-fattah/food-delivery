package handlers

import (
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/models"
	"encoding/json"
	"net/http"
)


func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	userID, _, err := helper.ParseJWT(authHeader)
	if err != nil {
		helper.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		helper.SendError(w, "No items", http.StatusBadRequest)
		return
	}

	orderID, err := database.CreateOrder(r.Context(), userID, req)
	if err != nil {
		helper.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"message":  "Order created successfully",
	})
}


func GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")

	userID, role, err := helper.ParseJWT(authHeader)
	if err != nil {
		helper.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	orders, err := database.GetOrders(r.Context(), userID, role)
	if err != nil {
		helper.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		helper.SendError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	_, role, err := helper.ParseJWT(authHeader)
	if err != nil {
		helper.SendError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if role != "admin" {
		helper.SendError(w, "Forbidden", http.StatusForbidden)
		return
	}

	orderID, err := helper.GetIDFromURL(r.URL.Path)

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		helper.SendError(w, "status is required", http.StatusBadRequest)
		return
	}

	err = database.UpdateOrderStatus(r.Context(), orderID, req.Status)
	if err != nil {
		helper.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Order status updated successfully",
		"order_id": orderID,
		"status":   req.Status,
	})
}