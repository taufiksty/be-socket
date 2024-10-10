package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/domain/repository"
)

type roomClient struct {
	RoomId uuid.UUID `json:"room_id"`
	UserId uuid.UUID `json:"user_id"`
}

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) repository.RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(room *entity.Room) error {
	query := `INSERT INTO rooms (id, name) VALUES ($1, $2)`
	_, err := r.db.Exec(query, room.ID, room.Name)
	return err
}

func (r *RoomRepository) GetRooms() ([]*entity.Room, error) {
	query := `SELECT id, name FROM rooms`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*entity.Room
	for rows.Next() {
		var room entity.Room
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, nil
}

func (r *RoomRepository) GetRoomByName(name string) (*entity.Room, error) {
	query := `SELECT id, name FROM rooms WHERE name = $1`
	row := r.db.QueryRow(query, name)

	var room entity.Room
	err := row.Scan(&room.ID, &room.Name)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepository) GetRoomById(roomID uuid.UUID) (*entity.Room, error) {
	query := `SELECT id, name FROM rooms WHERE id = $1`
	row := r.db.QueryRow(query, roomID)

	var room entity.Room
	err := row.Scan(&room.ID, &room.Name)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomRepository) AddClient(roomId, userId uuid.UUID) error {
	query := `INSERT INTO room_clients (room_id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, roomId, userId)
	return err
}

func (r *RoomRepository) VerifyClientRoom(roomId, userId uuid.UUID) (bool, error) {
	query := `SELECT room_id, user_id FROM room_clients WHERE room_id = $1 AND user_id = $2`
	row := r.db.QueryRow(query, roomId, userId)

	var rc roomClient
	err := row.Scan(&rc.RoomId, &rc.UserId)
	if err != nil {
		return false, err
	}

	if rc.RoomId != uuid.Nil || rc.UserId != uuid.Nil {
		return true, nil
	}

	return false, nil
}
