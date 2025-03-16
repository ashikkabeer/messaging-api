package models

import (
	"time"
)

type RequestBody struct {
    SenderID   string `json:"sender_id"`
    ReceiverID string `json:"receiver_id"`
    Content    string `json:"content"`
}

type Message struct {
    ID         string    `json:"message_id"`
    SenderID   string    `json:"sender_id"`
    ReceiverID string    `json:"receiver_id"`
    Content    string    `json:"content"`
    CreatedAt  time.Time `json:"timestamp"`
    Read       bool      `json:"read"`
}