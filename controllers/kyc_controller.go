package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var kycCollection *mongo.Collection

func InitKYCCollection() {
	kycCollection = config.DB.Collection("users")
}

// 📤 POST /kyc/upload
func UploadKYC(c *gin.Context) {
	// 1. Auth Check
	emailVal, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userEmail := emailVal.(string)

	// 2. Nil check on collection
	if kycCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KYC collection not initialized"})
		return
	}

	// 3. Parse file
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "KYC document required"})
		return
	}

	// 4. Optional: Check file size / extension / content-type
	// if file.Size > 5*1024*1024 { ... }

	// 5. Create upload path
	uploadPath := "uploads/kyc"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload folder"})
			return
		}
	}

	// 6. Save file
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	fullPath := filepath.Join(uploadPath, filename)
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// 7. Update DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"kyc_status":  "pending",
			"kyc_doc_url": "/" + fullPath,
		},
	}

	_, err = kycCollection.UpdateOne(ctx, bson.M{"email": userEmail}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KYC status"})
		return
	}

	// 8. Success response
	c.JSON(http.StatusOK, gin.H{
		"message":     "KYC uploaded successfully",
		"documentURL": "/" + fullPath,
	})
}

// 🔍 GET /kyc/status
func GetKYCStatus(c *gin.Context) {
	emailVal, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userEmail := emailVal.(string)

	if kycCollection == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "KYC collection not initialized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result struct {
		KYCStatus string `bson:"kyc_status" json:"kyc_status"`
		KYCDocURL string `bson:"kyc_doc_url" json:"kyc_doc_url"`
	}

	err := kycCollection.FindOne(ctx, bson.M{"email": userEmail}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User or KYC data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"kyc_status": result.KYCStatus,
		"kyc_doc":    result.KYCDocURL,
	})
}
// ✅ PUT /kyc/verify/:email
func VerifyKYC(c *gin.Context) {
	email := c.Param("email")

	var payload struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || (payload.Status != "approved" && payload.Status != "rejected") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status (must be 'approved' or 'rejected')"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"kyc_status": payload.Status}}
	res, err := kycCollection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil || res.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KYC status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC status updated", "email": email, "status": payload.Status})
}
