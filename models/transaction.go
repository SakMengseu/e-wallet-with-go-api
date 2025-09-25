package models

import (
	"time"
)

type Transaction struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	UserID     string    `bson:"user_id,omitempty" json:"user_id"`
	SenderID   string    `bson:"sender_id,omitempty" json:"sender_id"`
	ReceiverID string    `bson:"receiver_id,omitempty" json:"receiver_id"`
	Amount     float64   `bson:"amount" json:"amount"`
	Type       string    `bson:"type" json:"type"`     // deposit, transfer, withdraw
	Status     string    `bson:"status" json:"status"` // completed, failed
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
