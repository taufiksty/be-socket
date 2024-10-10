package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `json:"id"`
	RoomID         uuid.UUID `json:"room_id"`
	UserID         uuid.UUID `json:"user_id"`
	SenderUsername string    `json:"sender_username"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
}
