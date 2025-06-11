package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"-"`                  // hide password from API responses
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	Role      string    `gorm:"default:'user'" json:"role"` // include role in API response
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
