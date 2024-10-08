package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Room struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	rooms = make(map[string]*Room)
	mutex sync.Mutex
)

func newRoom() *Room {
	return &Room{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				close(client.send)
			}
		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

func handleWebSocket(c echo.Context) error {
	roomID := c.Param("roomID")
	mutex.Lock()
	room, ok := rooms[roomID]
	if !ok {
		room = newRoom()
		rooms[roomID] = room
		go room.run()
	}
	mutex.Unlock()

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	room.register <- client

	go func() {
		defer func() {
			room.unregister <- client
			client.conn.Close()
		}()
		for {
			_, message, err := client.conn.ReadMessage()
			if err != nil {
				break
			}
			room.broadcast <- message
		}
	}()

	go func() {
		for {
			message, ok := <-client.send
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.conn.WriteMessage(websocket.TextMessage, message)
		}
	}()

	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8081"},
	}))
	e.GET("/ws/:roomID", handleWebSocket)
	e.Logger.Fatal(e.Start(":8080"))
}
