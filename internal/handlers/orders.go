package handlers

import (
	"github.com/akram-fattah/food-delivery/internal/helper"
	"github.com/akram-fattah/food-delivery/internal/database"
	"github.com/akram-fattah/food-delivery/internal/models"
	"encoding/json"
	"net/http"
	"fmt"
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

	email, _ := database.GetUserEmailByID(userID)
	var itemsWithDetails []models.OrderItem
	var totalPrice float64
	
	for _, item := range req.Items {
		var price float64
		var name string

		err := database.DB.QueryRowContext(
			r.Context(),
			"SELECT name, price FROM meals WHERE id = $1",
			item.MealID,
		).Scan(&name, &price)

		if err != nil {
			helper.SendError(w, fmt.Sprintf("meal not found: %d", item.MealID), http.StatusBadRequest)
			return
		}

		totalPrice += price * float64(item.Quantity)

		itemsWithDetails = append(itemsWithDetails, models.OrderItem{
			MealID:   item.MealID,
			Name:     name,
			Quantity: item.Quantity,
			Price:    price,
		})
	}
	SendNotification(email, itemsWithDetails)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"message":  "Order created successfully",
	})
}


func SendNotification(email string, items []models.OrderItem) {
body := `
<html box-sizing: border-box;>
	<head>
        <link href="https://fonts.googleapis.com/css2?family=Tajawal:wght@400;700&display=swap" rel="stylesheet">
    </head>
	<body style="margin:0; padding:0 5px; background:#f8f9fa; direction:rtl; text-align:center; font-family:'Tajawal', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; ">

<div style="width: 330px; background:#f4f6f8; padding:20px 10px; box-sizing:border-box; border-box: 8px; text-align: center; margin: 0 auto; border: 1px solid #e1e1e1; border-radius: 12px">
                <h1 style="color: #e87224; margin: 0; font-size: 28px; font-weight: bold; letter-spacing: 1px; font-family: 'Tajawal', sans-serif;">مطعم الوحدة</h1>
                <p style="color: #333; margin: 5px 0 0 0; font-size: 14px; font-family: 'Tajawal', sans-serif;">أصل المذاق اليمني</p>
<table role="presentation" max-width="370px" cellspacing="0" cellpadding="0">
<tr>
<td align="center">

<!-- CARD -->
<table role="presentation" width="310px" cellspacing="0" cellpadding="0"
style="width:100%; max-width:100%; background:#ffffff; border-radius:14px; overflow:hidden; box-shadow:0 6px 20px rgba(0,0,0,0.08);">

<!-- CONTENT -->
<tr>
<td style="padding:25px; direction:rtl; text-align:center; font-family:Tajawal, sans-serif;">

	<h3 style="color:#e87224; margin-bottom:10px;">شكراً لطلبك ❤️</h3>

	<p style="color:#555; line-height:1.8; margin-bottom:20px;">
		تم استلام طلبك بنجاح ونحن نعمل على تحضيره الآن.<br>
		سنقوم بإعلامك بمجرد أن يكون جاهزاً للتوصيل.
	</p>

	<h4 style="margin:15px 0; color:#333;">تفاصيل الطلب</h4>

	<table role="presentation" width="100%" cellspacing="0" cellpadding="0"
	style="border-collapse:collapse; width:100%; font-size:14px;">

		<tr style="background:#f1f3f5;">
			<th style="padding:12px; border-bottom:1px solid #ddd;">الوجبة</th>
			<th style="padding:12px; border-bottom:1px solid #ddd;">الكمية</th>
			<th style="padding:12px; border-bottom:1px solid #ddd;">السعر</th>
			<th style="padding:12px; border-bottom:1px solid #ddd;">الإجمالي</th>
		</tr>
`
var total float64

for _, item := range items {
	itemTotal := item.Price * float64(item.Quantity)
	total += itemTotal

	body += fmt.Sprintf(`
	<tr>
		<td style="padding:12px; border-bottom:1px solid #eee; text-align:center;">%s</td>
		<td style="padding:12px; border-bottom:1px solid #eee; text-align:center;">%d</td>
		<td style="padding:12px; border-bottom:1px solid #eee; text-align:center;">%.2f</td>
		<td style="padding:12px; border-bottom:1px solid #eee; text-align:center; font-weight:bold; color:#e87224;">%.2f</td>
	</tr>
	`, item.Name, item.Quantity, item.Price, itemTotal)
}
body += fmt.Sprintf(`
	</table>

	<!-- TOTAL -->
	<div style="margin-top:20px; text-align:center;">
		<div style="display:inline-block; background:#e87224; color:#fff; padding:10px 18px; border-radius:20px; font-weight:bold;">
			الإجمالي الكلي: %.2f
		</div>
	</div>

	<p style="color:#555; margin-top:20px;">
		إذا كان لديك أي استفسار، تواصل معنا
	</p>

	<p style="color:#888;">
		نتمنى لك وجبة شهية 🍽️
	</p>

</td>
</tr style="background:#f43;">


</table>

</td>
</tr>
</table>
	<p style="color: #bdc3c7; font-size: 11px; margin: 0; font-family: 'Tajawal', sans-serif;">
		جميع الحقوق محفوظة © 2026 - فريق تطوير اكرم سوفت
	</p>
</div>
</body>
</html>
`, total)
	helper.SendEmail("تأكيد طلبك من مطعم الوحدة",email ,body, " ")
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