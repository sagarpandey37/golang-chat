package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"chat/utils"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

var clients = make(map[string]*utils.Client)
var broadcaster = make(chan utils.ClientsMeta)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func messageClient(client *websocket.Conn, msg utils.ClientsMeta, reciever string) {

	err := client.WriteJSON(msg)

	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, reciever)
	}
}

func getRecieverID(msg utils.ClientsMeta) string {
	recieverID := msg.Reciever.UserID
	recieverOut, err := json.Marshal(recieverID)

	if err != nil {
		panic(err)
	}

	return string(recieverOut)
}

func messageClients(msg utils.ClientsMeta) {
	// send to every client currently connected
	// for client := range clients {
	// 	messageClient(client, msg)
	// }

	// messageObject, _ := json.Marshal(msg)
	// log.Printf("jsonInfo: %s\n", messageObject)

	receiverID := getRecieverID(msg)

	clientSocket, found := clients[receiverID]

	if found {
		messageClient(clientSocket.WebSocketConn, msg, receiverID)
	} else {
		log.Println(found)
	}

}

// func sendPreviousMessages(ws *websocket.Conn) {
// 	chatMessages, err := rdb.LRange(ctx, "chat_messages", 0, -1).Result()
// 	if err != nil {
// 		panic(err)
// 	}

// 	// send previous messages
// 	for _, chatMessage := range chatMessages {
// 		var msg utils.Message
// 		json.Unmarshal([]byte(chatMessage), &msg)
// 		messageClient(ws, msg)
// 	}
// }

func handleConnections(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()["userID"][0]

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	//ensure connection close when function returns
	defer ws.Close()

	// Register or store new clients ------- step-1
	clients[q] = &utils.Client{WebSocketConn: ws}

	log.Println(clients)

	// // if it's zero, no messages were ever sent/saved ------- step-2 ( send previous chat message but check any any message exist or not in redis)
	// if rdb.Get(ctx, "chat_messages").Val() == "" {
	// 	sendPreviousMessages(ws)
	// }

	//Read Message from other clients and brodcast it
	for {
		var msg utils.ClientsMeta

		//Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			receiverID := getRecieverID(msg)
			delete(clients, receiverID)
			break
		}

		// send new message  to the channel : after this handleMessage() core go routine perform rest task
		broadcaster <- msg
	}

}

func handleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster

		// store message to redis & then send to client
		// storeInRedis(msg)
		messageClients(msg)
	}

}

func main() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	// Load Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})
	_, err1 := rdb.Ping(ctx).Result()

	if err1 != nil {
		panic(err1)
	}

	// Routings
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/websocket", handleConnections)
	go handleMessages()

	// Exec Server
	log.Print("Server starting at localhost:3000")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}

}
