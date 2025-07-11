package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/ashishnagargoje0/backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)



// ✅ Call this in main.go after MongoDB connection is established
func InitAuthCollection() {
	if config.DB == nil {
		panic("MongoDB not initialized before InitAuthCollection")
	}
	userCollection = config.DB.Collection("users")
}

// ---------------------------
// 📌 EMAIL + PASSWORD FLOW
// ---------------------------

func Signup(c *gin.Context) {
	var input models.SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": input.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Phone:     input.Phone,
		Role:      "user",
		CreatedAt: time.Now().Unix(),
	}

	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"token":   token,
	})
}

func Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ✅ Now generate token with email and role
	token, err := utils.GenerateToken(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID.Hex(),     // ✅ Ensure userID is returned as string
			"email": user.Email,
			"role":  user.Role,
		},
	})
}




func ForgotPassword(c *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required,email"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	resetToken, err := utils.GenerateResetToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Password reset token generated",
		"reset_token": resetToken,
	})
}

func ResetPassword(c *gin.Context) {
	type Request struct {
		Email       string `json:"email" binding:"required,email"`
		ResetToken  string `json:"reset_token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	email, err := utils.ValidateResetToken(req.ResetToken)
	if err != nil || email != req.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}
	_, err = userCollection.UpdateOne(ctx, bson.M{"email": req.Email}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// ---------------------------
// 📱 PHONE + OTP FLOW
// ---------------------------

func RegisterPhone(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	otp, err := utils.GenerateOTP(input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}
	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"phone": input.Phone}).Decode(&user)

	if err == mongo.ErrNoDocuments {
		newUser := models.User{
			Phone:      input.Phone,
			Name:       input.Name,
			OTP:        otp,
			IsVerified: false,
			CreatedAt:  time.Now().Unix(),
			Role:       "user",
		}
		_, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	} else {
		_, err := userCollection.UpdateOne(ctx,
			bson.M{"phone": input.Phone},
			bson.M{"$set": bson.M{"otp": otp}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OTP sent to phone",
		"otp_sent": true,
	})
}

func VerifyOTP(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"phone": input.Phone}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != input.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	_, err = userCollection.UpdateOne(ctx,
		bson.M{"phone": input.Phone},
		bson.M{
			"$set":   bson.M{"is_verified": true},
			"$unset": bson.M{"otp": ""},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Verification failed"})
		return
	}

	identifier := user.Email
	if identifier == "" {
		identifier = user.Phone
	}

	token, err := utils.GenerateToken(identifier, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OTP verified successfully",
		"token":   token,
		"user": gin.H{
			"name":  user.Name,
			"phone": user.Phone,
		},
	})
}
