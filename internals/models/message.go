package models

import "time"

type Message struct {
	ID          string
	Content     string
	SenderID    string
	RecipientID string
	CreatedAt   time.Time
}
