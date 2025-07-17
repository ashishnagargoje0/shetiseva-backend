package chatbot

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// HandleMessage processes incoming chat messages
func HandleMessage(c *gin.Context) {
	var req ChatRequest // ✅ Use type directly, no chatbot. prefix

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message is required"})
		return
	}

	// 🧠 TODO: Add AI logic here
	reply := "Thanks for your message: " + req.Message

	// Respond to client
	c.JSON(http.StatusOK, ChatResponse{ // ✅ Use type directly
		Reply: reply,
	})
}
