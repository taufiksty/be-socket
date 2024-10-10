package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/domain/repository"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) repository.ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessage(msg *entity.Message) error {
	query := `INSERT INTO messages (id, room_id, user_id, sender_username, content, timestamp) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, msg.ID, msg.RoomID, msg.UserID, msg.SenderUsername, msg.Content, msg.Timestamp)
	return err
}

func (r *ChatRepository) GetMessagesByRoom(roomID uuid.UUID) ([]*entity.Message, error) {
	query := `SELECT id, room_id, user_id, sender_username, content, timestamp 
			FROM messages WHERE room_id = $1 ORDER BY timestamp ASC`
	rows, err := r.db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.RoomID, &msg.UserID, &msg.SenderUsername, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	// Check if there was any error during rows iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
