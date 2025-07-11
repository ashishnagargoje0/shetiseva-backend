package models

import "time"

// ========== Auth ==========
type SignupInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,len=10"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// ========== Return ==========
type ReturnRequestInput struct {
	OrderID    string `json:"order_id" binding:"required"`
	Reason     string `json:"reason" binding:"required"`
	ReturnType string `json:"return_type" binding:"required"`
}

// ========== Refund ==========
type RefundRequestInput struct {
	OrderID string  `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
	Reason  string  `json:"reason" binding:"required"`
}

// ========== Support ==========
type SupportTicketInput struct {
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// ========== Review ==========
type ReviewInput struct {
	ProductID string `json:"productId" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

// ========== Voice Feedback ==========
type VoiceFeedbackInput struct {
	UserID   string `json:"userId" binding:"required"`
	AudioURL string `json:"audioUrl" binding:"required,url"`
	Comment  string `json:"comment,omitempty"`
}





// ========== Weather ==========
type WeatherData struct {
	Location string `bson:"location"`
	Date     string `bson:"date"`
	Temp     string `bson:"temp"`
	Rain     string `bson:"rain"`
}

type WeatherForecast struct {
	Location string `bson:"location"`
	Day      string `bson:"day"`
	Summary  string `bson:"summary"`
}

// ========== AI Alert ==========
type AIAlert struct {
	Title     string    `bson:"title"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"timestamp"`
}

// ========== Advisory ==========
type CropAdvisory struct {
	CropID string `bson:"crop_id"`
	Advice string `bson:"advice"`
	Tips   string `bson:"tips"`
}
