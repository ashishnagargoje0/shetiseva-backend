package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name,omitempty" json:"name"`
	Email      string             `bson:"email,omitempty" json:"email"`
	Phone      string             `bson:"phone,omitempty" json:"phone"`
	Password   string             `bson:"password,omitempty" json:"-"`
	IsVerified bool               `bson:"is_verified,omitempty" json:"is_verified"`
	Role       string             `bson:"role,omitempty" json:"role"`
	CreatedAt  int64              `bson:"created_at,omitempty" json:"created_at"`
	OTP        string             `bson:"otp,omitempty" json:"otp,omitempty"`

	// 🆕 KYC fields
	KYCStatus string `bson:"kyc_status,omitempty" json:"kyc_status,omitempty"`
	KYCDocURL string `bson:"kyc_doc_url,omitempty" json:"kyc_doc_url,omitempty"`
}
