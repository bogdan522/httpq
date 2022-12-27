package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	kind    int
	message string
}

var topicMap = make(map[string]Topic)

type Topic struct {
	messages []Message
	clients  []*websocket.Conn
	name     string
}

func reader(conn *websocket.Conn, t string) {
	for {
		// read in a message
		topic := topicMap[t]
		// This should be in another go routine
		for _, client := range topic.clients {
			// should check if client is still alive, if so remove from slice
			messageType, p, err := client.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}

			msg := Message{kind: messageType, message: string(p)}
			topic.messages = append(topic.messages, msg)
		}

		// this should be in its own section
		for _, msg := range topic.messages {
			// should check if client is still alive, if so remove from slice
			for _, client := range topic.clients {
				if err := client.WriteMessage(msg.kind, []byte(msg.message)); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	topicUrl := r.URL.Path[4:]
	t, ok := topicMap[topicUrl]

	if !ok {
		t = Topic{name: topicUrl}
		topicMap[topicUrl] = t
	}

	t.clients = append(t.clients, ws)

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws, topicUrl)
}

func setupRoutes(r *mux.Router) {

	r.HandleFunc("/", homePage)
	r.HandleFunc("/ws/{topic}", wsEndpoint)
}
