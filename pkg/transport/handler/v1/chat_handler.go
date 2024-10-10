package v1

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/taufiksty/be-socket/pkg/adapter/controller"
	"github.com/taufiksty/be-socket/pkg/shared/util"
)

type ChatHandler struct {
	chatController *controller.ChatController
}

type CreateRoomRequest struct {
	RoomID uuid.UUID `json:"id"`
	Name   string    `json:"name"`
}

func NewChatHandler(chatController *controller.ChatController) *ChatHandler {
	return &ChatHandler{chatController: chatController}
}

// func (h *ChatHandler) SendMessage(c echo.Context) error {
// 	senderIDRequest := c.FormValue("sender_id")
// 	roomIDRequest := c.FormValue("room_id")
// 	content := c.FormValue("content")

// 	senderID, err := uuid.Parse(senderIDRequest)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}
// 	roomID, err := uuid.Parse(roomIDRequest)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	err = h.chatController.SendMessage(senderID, roomID, content)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusCreated, nil)
// }

func (h *ChatHandler) CreateRoom(c echo.Context) error {
	var req CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userIdInterface := c.Get("user_id")
	userId, ok := userIdInterface.(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusInternalServerError, "invalid user_id type")
	}

	err := h.chatController.CreateRoom(req.Name, req.RoomID, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *ChatHandler) GetRooms(c echo.Context) error {
	rooms, err := h.chatController.GetRooms()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rooms)
}

func (h *ChatHandler) GetRoomById(c echo.Context) error {
	fmt.Println("Check handler called")
	roomIDRequest := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	room, err := h.chatController.GetRoomById(roomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}

func (h *ChatHandler) GetMessages(c echo.Context) error {
	roomIDRequest := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDRequest)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	messages, err := h.chatController.GetMessages(roomID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) JoinRoom(c echo.Context) error {
	// Extract JWT token from query params
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
	}

	// Validate JWT token
	claims, err := util.ParseToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	// Extract user ID from JWT claims
	userID := claims.UserId

	roomId, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// Proceed with WebSocket connection
	err = h.chatController.JoinRoom(c, roomId, userID)
	if err != nil {
		return err
	}

	return nil
}
