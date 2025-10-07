package models

type CreateChatRequest struct {
	UserID      string `json:"user_id"`
	RecipientID string `json:"recipient_id"`
}

type SendMessageRequest struct {
	UserID      string `json:"user_id"`
	RecipientID string `json:"recipient_id"`
	Message     string `json:"message"`
}
