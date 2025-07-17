package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ashishnagargoje0/backend/config"
	"github.com/ashishnagargoje0/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ===== Init Collections =====

var ConsultationCollection *mongo.Collection
var DroneBookingCollection *mongo.Collection

func InitConsultationCollection() {
	ConsultationCollection = config.DB.Collection("consultations")
}

func InitDroneBookingCollection() {
	DroneBookingCollection = config.DB.Collection("drone_bookings")
}

// ====== POST /consultation/book ======

func BookConsultation(c *gin.Context) {
	var input models.ConsultationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	consultation := bson.M{
		"user_id":   userID,
		"expert_id": input.ExpertID,
		"topic":     input.Topic,
		"status":    "booked",
		"timestamp": time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := ConsultationCollection.InsertOne(ctx, consultation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book consultation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consultation booked", "consultation_id": res.InsertedID})
}

// ====== GET /consultation/status/:id ======

func GetConsultationStatus(c *gin.Context) {
	idParam := c.Param("id")
	consultationID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consultation ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result bson.M
	err = ConsultationCollection.FindOne(ctx, bson.M{"_id": consultationID}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consultation not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ====== POST /consultation/feedback ======

func SubmitConsultationFeedback(c *gin.Context) {
	var input models.ConsultationFeedback
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid feedback data"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	input.UserID = userIDRaw.(primitive.ObjectID)
	input.Timestamp = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := config.DB.Collection("consultation_feedback").InsertOne(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Feedback submitted successfully"})
}

// ====== POST /drone/book ======

func BookDrone(c *gin.Context) {
	var input models.DroneBookingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drone booking data"})
		return
	}

	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	drone := bson.M{
		"user_id":    userID,
		"field_area": input.FieldArea,
		"location":   input.Location,
		"purpose":    input.Purpose,
		"date":       input.Date,
		"status":     "pending",
		"timestamp":  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := DroneBookingCollection.InsertOne(ctx, drone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Drone booking failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking submitted", "booking_id": res.InsertedID})
}

// ====== GET /drone/availability ======

func GetDroneAvailability(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	today := time.Now().Format("2006-01-02")
	filter := bson.M{
		"status": "available",
		"date":   bson.M{"$gte": today},
	}

	cursor, err := config.DB.Collection("drone_bookings").Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching drone availability"})
		return
	}
	defer cursor.Close(ctx)

	var drones []bson.M
	if err = cursor.All(ctx, &drones); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Decoding failed"})
		return
	}

	c.JSON(http.StatusOK, drones)
}


// ====== POST /drone/cancel ======

func CancelDroneBooking(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "status": bson.M{"$in": []string{"pending", "booked"}}}
	update := bson.M{"$set": bson.M{"status": "cancelled"}}

	result, err := DroneBookingCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel drone booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking cancelled"})
}

// GET /drone/status (User)
func GetMyDroneBookings(c *gin.Context) {
	userIDRaw, _ := c.Get("user_id")
	userID := userIDRaw.(primitive.ObjectID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("drone_bookings").Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drone bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []bson.M
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding drone bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// GET /admin/drone/bookings (Admin)
func AdminGetAllDroneBookings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := config.DB.Collection("drone_bookings").Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch drone bookings"})
		return
	}
	defer cursor.Close(ctx)

	var bookings []bson.M
	if err := cursor.All(ctx, &bookings); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding drone bookings"})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// PUT /admin/drone/approve/:id
func AdminApproveDroneBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "approved"}}
	result, err := config.DB.Collection("drone_bookings").UpdateByID(ctx, bookingID, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve drone booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking approved"})
}

// PUT /admin/drone/reject/:id
func AdminRejectDroneBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "rejected"}}
	result, err := config.DB.Collection("drone_bookings").UpdateByID(ctx, bookingID, update)
	if err != nil || result.ModifiedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject drone booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking rejected"})
}

// POST /drone/approve/:id
func ApproveDroneBooking(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drone booking ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "available"}}
	_, err = config.DB.Collection("drone_bookings").UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Approval failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking approved"})
}

// POST /drone/reject/:id
func RejectDroneBooking(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid drone booking ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "rejected"}}
	_, err = config.DB.Collection("drone_bookings").UpdateByID(ctx, objID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Rejection failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Drone booking rejected"})
}
