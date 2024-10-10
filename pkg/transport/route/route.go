package route

import (
	"github.com/labstack/echo/v4"
	"github.com/taufiksty/be-socket/pkg/shared/middleware"
	v1 "github.com/taufiksty/be-socket/pkg/transport/handler/v1"
)

func SetupRoutes(e *echo.Echo, hUser *v1.UserHandler, hChat *v1.ChatHandler) {
	apiGroup := e.Group("/api/v1")

	userGroup := apiGroup.Group("/user")
	userGroup.POST("/register", hUser.Register)
	userGroup.POST("/login", hUser.Login)

	roomGroup := apiGroup.Group("/rooms", middleware.JwtMiddleware)
	roomGroup.POST("", hChat.CreateRoom)
	roomGroup.GET("", hChat.GetRooms)
	roomGroup.GET("/:room_id", hChat.GetRoomById)
	roomGroup.GET("/messages/:room_id", hChat.GetMessages)

	e.GET("/ws/:room_id", hChat.JoinRoom)
}
