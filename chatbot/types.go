package chatbot

// ChatRequest represents the user's message input to the chatbot.
type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatResponse represents the chatbot's reply.
type ChatResponse struct {
	Reply string `json:"reply"`
}
