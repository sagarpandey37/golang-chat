package utils

import "github.com/gorilla/websocket"

// User Details
type User struct {
	UserName             string `json:"userName"`
	UserID               int    `json:"userID"`
	UserStatus           bool   `json:"userStatus"` // True -> Online False-> Offline
	UserLastActivityTime string `json:"userLastActivityTime"`
}

// ClientsMeta details
type ClientsMeta struct {
	ChannelKey  int     `json:"channelKey"`
	ChannelType int     `json:"channelType"` // 0 -> one to one chat  // 1 -> one to many chat
	Sender      User    `json:"sender"`
	Reciever    User    `json:"reciever"`
	CreateDate  string  `json:"createDate"`
	Data        Message `json:"message"`
}

// Message Details
type Message struct {
	Text       string `json:"text,omitempty"`
	Author     User   `json:"author"`
	CreateDate string `json:"createDate"`
}

// Client socket
type Client struct {
	WebSocketConn *websocket.Conn
}
