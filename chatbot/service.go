package chatbot

import (
	"strings"
)

// GenerateResponse processes the user message and returns a chatbot reply.
func GenerateResponse(message string) string {
	// Normalize input
	msg := strings.ToLower(strings.TrimSpace(message))

	// Handle empty input
	if msg == "" {
		return "Please enter your query. I’m here to help you with Shetiseva."
	}

	// Sample intent recognition (expand this later)
	switch {
	case strings.Contains(msg, "hello") || strings.Contains(msg, "hi"):
		return "Hi! 👋 Welcome to Shetiseva. How can I assist you today?"

	case strings.Contains(msg, "fertilizer"):
		return "We offer a wide range of fertilizers. Would you like Urea, DAP, or Organic options?"

	case strings.Contains(msg, "seed"):
		return "We have high-quality seeds for wheat, rice, cotton, and vegetables. Let me know what you're looking for."

	case strings.Contains(msg, "order status"):
		return "To check your order status, please visit the Orders section or provide your Order ID."

	case strings.Contains(msg, "support"):
		return "For any issues or help, you can raise a support ticket from the Support section."

	case strings.Contains(msg, "soil test"):
		return "You can book a soil testing service from our Soil Test section. Would you like to proceed?"

	case strings.Contains(msg, "bye"):
		return "Thank you for visiting Shetiseva. Have a great day! 🌾"

	default:
		return "I'm still learning! 😊 Please ask about seeds, fertilizers, orders, or Shetiseva services."
	}
}
