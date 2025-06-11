package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	Order     Order   `json:"-"` // Avoid recursion in JSON
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}
