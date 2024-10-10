package app

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taufiksty/be-socket/pkg/adapter/controller"
	dbAdapter "github.com/taufiksty/be-socket/pkg/adapter/db"
	"github.com/taufiksty/be-socket/pkg/adapter/repository/postgres"
	v1 "github.com/taufiksty/be-socket/pkg/transport/handler/v1"
	"github.com/taufiksty/be-socket/pkg/transport/route"
	"github.com/taufiksty/be-socket/pkg/usecase/auth"
	"github.com/taufiksty/be-socket/pkg/usecase/chat"
)

var dbInstance *sql.DB

func Run() {
	err := godotenv.Load("internal/config/.env")
	if err != nil {
		log.Fatal("Error load env", err)
	}

	dbInstance = dbAdapter.InitDB()
	defer dbInstance.Close()

	// Initialize repository
	userRepo := postgres.NewUserRepository(dbInstance)
	roomRepo := postgres.NewRoomRepository(dbInstance)
	chatRepo := postgres.NewChatRepository(dbInstance)

	// Initialize usecase
	authUsecase := auth.NewAuthUsecase(userRepo)
	chatUsecase := chat.NewChatUsecase(roomRepo, chatRepo)

	// Initialize controller
	userController := controller.NewUserController(authUsecase)
	chatController := controller.NewChatController(chatUsecase)

	// Initialize handler
	userHandler := v1.NewUserHandler(userController)
	chatHandler := v1.NewChatHandler(chatController)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
	}))
	e.Use(middleware.Recover())

	route.SetupRoutes(e, userHandler, chatHandler)

	e.Logger.Fatal(e.Start(":3000"))
}
