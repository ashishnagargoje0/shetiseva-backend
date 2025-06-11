package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  uint    `json:"quantity"`

	OrderID *uint `json:"order_id"` // Nullable field to track checkout status
}
