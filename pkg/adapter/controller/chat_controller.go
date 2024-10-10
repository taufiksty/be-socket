package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/taufiksty/be-socket/pkg/domain/entity"
	"github.com/taufiksty/be-socket/pkg/usecase/chat"
)

// In-memory room map to track clients in each room
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomLock sync.Mutex

type ChatController struct {
	chatUseCase chat.ChatUsecase
}

type saveMessageRequest struct {
	ID             string `json:"id"`
	RoomID         string `json:"room_id"`
	UserID         string `json:"user_id"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
}

func NewChatController(chatUseCase *chat.ChatUsecase) *ChatController {
	return &ChatController{chatUseCase: *chatUseCase}
}

// func (ctrl *ChatController) SendMessage(senderID, roomID uuid.UUID, content string) error {
// 	err := ctrl.chatUseCase.SendMessage(senderID, roomID, content)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (ctrl *ChatController) CreateRoom(name string, roomId, userId uuid.UUID) error {
	err := ctrl.chatUseCase.CreateRoom(name, roomId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (ctrl *ChatController) GetRooms() (interface{}, error) {
	rooms, err := ctrl.chatUseCase.GetRooms()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (ctrl *ChatController) GetRoomById(roomID uuid.UUID) (interface{}, error) {
	room, err := ctrl.chatUseCase.GetRoomById(roomID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (ctrl *ChatController) GetMessages(roomID uuid.UUID) (interface{}, error) {
	messages, err := ctrl.chatUseCase.GetMessagesByRoom(roomID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ctrl *ChatController) JoinRoom(c echo.Context, roomID, userID uuid.UUID) error {
	// Upgrade connection to WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	// Add user to room in-memory map
	roomLock.Lock()
	if rooms[roomID.String()] == nil {
		rooms[roomID.String()] = make(map[*websocket.Conn]bool)
	}
	rooms[roomID.String()][conn] = true
	roomLock.Unlock()

	// Add user to room in database
	err = ctrl.chatUseCase.JoinRoom(roomID, userID)
	if err != nil {
		return err
	}

	// Handle WebSocket messages
	go func() {
		for {
			_, messageBytes, err := conn.ReadMessage()
			if err != nil {
				conn.Close()

				// Remove user from room on disconnect
				roomLock.Lock()
				delete(rooms[roomID.String()], conn)
				roomLock.Unlock()

				return
			}

			// Unmarshall the message
			var messageRequest saveMessageRequest
			err = json.Unmarshal(messageBytes, &messageRequest)
			if err != nil {
				log.Println("error decoding message:", err)
				return
			}

			parsedTime, err := time.Parse(time.RFC3339, messageRequest.Timestamp)
			if err != nil {
				fmt.Println("error parsing time:", err)
				return
			}

			// Create the Message struct with the Unix timestamp
			message := entity.Message{
				ID:             uuid.MustParse(messageRequest.ID),
				RoomID:         uuid.MustParse(messageRequest.RoomID),
				UserID:         uuid.MustParse(messageRequest.UserID),
				SenderUsername: messageRequest.SenderUsername,
				Content:        messageRequest.Content,
				Timestamp:      parsedTime,
			}
			// Save message to database
			err = ctrl.chatUseCase.SendMessage(message)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Broadcast message to all clients in the room
			roomLock.Lock()
			for client := range rooms[roomID.String()] {
				client.WriteMessage(websocket.TextMessage, messageBytes)
			}
			roomLock.Unlock()
		}
	}()

	return nil
}
