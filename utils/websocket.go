package utils

import (
	"context"
	db "health/models/db"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kamva/mgm/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatRoom struct {
	Clients map[*websocket.Conn]bool
	Mutex   sync.Mutex
}

var chatrooms = make(map[string]*ChatRoom)

func getChatRoom(roomID string) *ChatRoom {
	if _, exists := chatrooms[roomID]; !exists {
		chatrooms[roomID] = &ChatRoom{Clients: make(map[*websocket.Conn]bool)}
	}
	return chatrooms[roomID]
}

func HandleWebSocket(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	chatRoom := getChatRoom(roomID)
	chatRoom.Mutex.Lock()
	chatRoom.Clients[conn] = true
	chatRoom.Mutex.Unlock()
	for {
		var msg db.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		msg.RoomID = roomID
		saveMessage(msg)
		broadcast(chatRoom, msg)
	}

	chatRoom.Mutex.Lock()
	delete(chatRoom.Clients, conn)
	chatRoom.Mutex.Unlock()
}

func saveMessage(msg db.Message) {
	collection := mgm.Coll(&msg)
	_, err := collection.InsertOne(context.TODO(), msg)
	if err != nil {
		log.Println("Error saving message:", err)
	}
}

func broadcast(chatRoom *ChatRoom, msg db.Message) {
	chatRoom.Mutex.Lock()
	defer chatRoom.Mutex.Unlock()

	for client := range chatRoom.Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(chatRoom.Clients, client)
		}
	}
}
