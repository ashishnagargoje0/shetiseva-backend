package models

import "time"

type ContactMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Message   string    `json:"message" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
