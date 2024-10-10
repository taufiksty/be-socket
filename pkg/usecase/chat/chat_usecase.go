package chat

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/domain/repository"
)

type ChatUsecase struct {
	roomRepo repository.RoomRepository
	chatRepo repository.ChatRepository
}

func NewChatUsecase(roomRepo repository.RoomRepository, chatRepo repository.ChatRepository) *ChatUsecase {
	return &ChatUsecase{
		roomRepo: roomRepo,
		chatRepo: chatRepo,
	}
}

func (uc *ChatUsecase) CreateRoom(name string, roomId, userId uuid.UUID) error {
	fmt.Println("Check create room use case called")
	if roomGet, err := uc.roomRepo.GetRoomByName(name); (err != nil && err != sql.ErrNoRows) || roomGet != nil {
		return err
	}
	fmt.Println("Check room name completely")

	room := &entity.Room{
		ID:   roomId,
		Name: name,
	}

	err := uc.roomRepo.CreateRoom(room)
	if err != nil {
		return err
	}
	fmt.Println("Check create room completely")
	err = uc.roomRepo.AddClient(room.ID, userId)
	if err != nil {
		return err
	}
	fmt.Println("Check add client completely")
	return nil
}

func (uc *ChatUsecase) JoinRoom(roomID, userID uuid.UUID) error {
	if ok, _ := uc.roomRepo.VerifyClientRoom(roomID, userID); !ok {
		err := uc.roomRepo.AddClient(roomID, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *ChatUsecase) GetRooms() ([]*entity.Room, error) {
	return uc.roomRepo.GetRooms()
}

func (uc *ChatUsecase) GetRoomById(roomID uuid.UUID) (*entity.Room, error) {
	return uc.roomRepo.GetRoomById(roomID)
}

func (uc *ChatUsecase) SendMessage(message entity.Message) error {
	return uc.chatRepo.SaveMessage(&message)
}

func (uc *ChatUsecase) GetMessagesByRoom(roomID uuid.UUID) ([]*entity.Message, error) {
	return uc.chatRepo.GetMessagesByRoom(roomID)
}
