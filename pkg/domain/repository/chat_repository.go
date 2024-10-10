package repository

import (
	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
)

type ChatRepository interface {
	SaveMessage(msg *entity.Message) error
	GetMessagesByRoom(roomID uuid.UUID) ([]*entity.Message, error)
}
