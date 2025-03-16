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
    ID         string    `json:"id"`
    SenderID   string    `json:"sender_id"`
    ReceiverID string    `json:"receiver_id"`
    Content    string    `json:"content"`
    Read       bool      `json:"read"`
    CreatedAt  time.Time `json:"created_at"`
}