package repository

import (
	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
)

type RoomRepository interface {
	CreateRoom(room *entity.Room) error
	GetRooms() ([]*entity.Room, error)
	GetRoomById(roomID uuid.UUID) (*entity.Room, error)
	GetRoomByName(name string) (*entity.Room, error)
	AddClient(roomId, userId uuid.UUID) error
	VerifyClientRoom(roomId, userId uuid.UUID) (bool, error)
}
