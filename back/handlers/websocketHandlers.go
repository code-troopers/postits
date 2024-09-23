package handlers

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/code-troopers/postitsonline/database"
	"github.com/code-troopers/postitsonline/webtoken"
	"github.com/gofiber/contrib/websocket"
)

const (
	NEW_BOARD      string = "NEW_BOARD"      // params name
	RENAME_BOARD   string = "RENAME_BOARD"   // params id, text
	DELETE_BOARD   string = "DELETE_BOARD"   // params id
	NEW_POSTIT     string = "NEW_POSTIT"     // params {}
	UPDATE_CONTENT string = "UPDATE_CONTENT" // params {id, text}
	MOVE_POSTIT    string = "MOVE_POSTIT"    // params {id, posX, posY}
	DELETE_POSTIT  string = "DELETE_POSTIT"  // params id
	ADD_VOTE       string = "ADD_VOTE"       // params boardId, id
	REMOVE_VOTE    string = "REMOVE_VOTE"    // params boardId, id
	SHOW_POSTITS   string = "SHOW_POSTITS"   // params boardId
	HIDE_POSTITS   string = "HIDE_POSTITS"   // params boardId
)

type WebSocketHub struct {
	clients   map[*websocket.Conn]bool // Map des connexions WebSocket actives
	broadcast chan []byte              // Canal pour diffuser les messages
	lock      sync.Mutex               // Verrou pour accéder aux clients de manière concurrente
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (hub *WebSocketHub) addClient(client *websocket.Conn) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	hub.clients[client] = true
}

func (hub *WebSocketHub) removeClient(client *websocket.Conn) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	delete(hub.clients, client)
	client.Close()
}

func (hub *WebSocketHub) broadcastMessage(message []byte) {
	hub.lock.Lock()
	defer hub.lock.Unlock()
	for client := range hub.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Erreur lors de l'envoi du message :", err)
			client.Close()
			delete(hub.clients, client)
		}
	}
}

func (hub *WebSocketHub) Run() {
	for {
		message := <-hub.broadcast
		hub.broadcastMessage(message)
	}
}

type Message struct {
	Action   string        `json:"action"`
	ID       string        `json:"id"`
	BoardId  string        `json:"boardId"`
	AuthorId string        `json:"authorId"`
	Text     string        `json:"text"`
	PosX     int           `json:"posX"`
	PosY     int           `json:"posY"`
	Token    string        `json:"token"`
	Author   database.User `json:"author"`
	Weight   int           `json:"weight"`
}

func handleAction(message *Message) {
	switch message.Action {
	case NEW_BOARD:
		board, err := createBoard(message.Text)
		if err != nil {
			log.Printf("Erreur lors de la création du board : %v", err)
			return
		}
		message.ID = board.ID

	case RENAME_BOARD:
		renameBoard(message.BoardId, message.Text)

	case DELETE_BOARD:
		deleteBoard(message.BoardId)

	case DELETE_POSTIT:
		deletePostit(message.ID)

	case NEW_POSTIT:
		user, err := webtoken.DecodeJWT(message.Token)
		if err != nil {
			log.Printf("Erreur lors du décodage du token : %v", err)
			return
		}
		message.AuthorId = user.ID
		message.Author = user
		postit, err := CreatePostit(message.BoardId, message.Text, message.PosX, message.PosY, user.ID)
		if err != nil {
			log.Printf("Erreur lors de la création du postit : %v", err)
			return
		}
		message.ID = postit.ID
		message.Weight = postit.Weight

	case UPDATE_CONTENT:
		user, err := webtoken.DecodeJWT(message.Token)
		if err != nil {
			log.Printf("Erreur lors du décodage du token : %v", err)
			return
		}
		message.Author = user
		updatePostitContent(message.ID, message.Text, user.ID)

	case MOVE_POSTIT:
		w, _ := movePostit(message.ID, message.PosX, message.PosY)
		message.Weight = w

	case ADD_VOTE:
		addVote(message.ID)

	case REMOVE_VOTE:
		removeVote(message.ID)

	case SHOW_POSTITS:
		user, err := webtoken.DecodeJWT(message.Token)
		if err != nil {
			log.Printf("Erreur lors du décodage du token : %v", err)
			return
		}
		message.AuthorId = user.ID
		showPostits(message.BoardId, user.ID, true)

	case HIDE_POSTITS:
		user, err := webtoken.DecodeJWT(message.Token)
		if err != nil {
			log.Printf("Erreur lors du décodage du token : %v", err)
			return
		}
		message.AuthorId = user.ID
		showPostits(message.BoardId, user.ID, false)
	}
}

func WebsocketHandler(c *websocket.Conn, hub *WebSocketHub) {
	defer hub.removeClient(c)
	hub.addClient(c)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Erreur de lecture:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Erreur lors du décodage du message :", err)
			continue
		}

		log.Printf("Message reçu : %+v", message)

		handleAction(&message)
		message.Token = ""
		msgJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("Erreur lors de l'encodage du message :", err)
			continue
		}

		hub.broadcast <- msgJSON
	}
}
