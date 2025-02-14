package models

import "golang.org/x/net/websocket"

type Client struct {
	conn   *websocket.Conn
	userID string
	roomID string
	send   chan []byte
}
