package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type InvoiceItem struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type Invoice struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InvoiceNo   string             `bson:"invoice_number" json:"invoice_number"`
	OrderID     primitive.ObjectID `bson:"order_id" json:"order_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Amount      float64            `bson:"amount" json:"amount"`
	Total       float64            `bson:"total" json:"total"`
	TotalPrice  float64            `bson:"total_price" json:"total_price"`
	Date        time.Time          `bson:"date" json:"date"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	Items       []InvoiceItem      `bson:"items" json:"items"`
}
