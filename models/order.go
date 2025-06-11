package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID      uint        `json:"user_id"`
	User        User        `json:"user"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"` // pending, paid, shipped, etc.
	OrderItems  []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}
