package main

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

// Fonction pour récupérer la clé publique depuis Keycloak
func getKeycloakPublicKey() (*rsa.PublicKey, error) {
	keycloakCerts := os.Getenv("KEYCLOAK_CERTS")
	if keycloakCerts == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
		keycloakCerts = os.Getenv("KEYCLOAK_CERTS")
	}
	resp, err := http.Get(keycloakCerts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks struct {
		Keys []struct {
			N string `json:"n"` // Modulus
			E string `json:"e"` // Exponent
		} `json:"keys"`
	}

	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, err
	}

	nBytes, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwks.Keys[0].E)
	if err != nil {
		return nil, err
	}

	e := big.NewInt(0)
	e.SetBytes(eBytes)

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(e.Int64()),
	}

	return pubKey, nil
}

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

// Gestionnaire pour stocker les connexions actives et leur diffusion
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

func (hub *WebSocketHub) run() {
	for {
		// Diffuser le message à tous les clients connectés
		message := <-hub.broadcast
		hub.broadcastMessage(message)
	}
}

type Message struct {
	Action   string `json:"action"`
	ID       string `json:"id"`
	BoardId  string `json:"boardId"`
	AuthorId string `json:"authorId"`
	Text     string `json:"text"`
	PosX     int    `json:"posX"`
	PosY     int    `json:"posY"`
}

func validateToken(tokenString string, pubKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token, nil
}

func jwtMiddleware(pubKey *rsa.PublicKey) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Token manquant")
		}

		// Supprimer le préfixe "Bearer " du token s'il est présent
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		_, err := validateToken(token, pubKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Token invalide")
		}
		return c.Next()
	}
}

func websocketHandler(c *websocket.Conn, hub *WebSocketHub) {
	defer hub.removeClient(c) // Retirer le client lors de la fermeture

	hub.addClient(c) // Ajouter le client lors de la connexion

	// Boucle pour lire les messages WebSocket
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Erreur de lecture:", err)
			break
		}

		// Décodez le message JSON en struct Message
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Erreur lors du décodage du message :", err)
			continue
		}

		// check token

		// get user

		log.Printf("Message reçu : %+v", message)

		if message.Action == NEW_BOARD {
			u2, err := uuid.NewV4()
			if err != nil {
				log.Fatalf("failed to generate UUID: %v", err)
			}
			message.ID = u2.String()
			boards = append(boards, Board{ID: message.ID, Name: message.Text, Postits: []Postit{}})
		}

		if message.Action == RENAME_BOARD {
			for i, board := range boards {
				if board.ID == message.ID {
					boards[i].Name = message.Text
				}
			}
		}

		if message.Action == DELETE_BOARD {
			for i, board := range boards {
				if board.ID == message.ID {
					boards = append(boards[:i], boards[i+1:]...)
				}
			}
		}

		if message.Action == DELETE_POSTIT {
			for _, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.ID == message.ID {
							board.Postits = append(board.Postits[:j], board.Postits[j+1:]...)
							// TODO optimize this
						}
					}
				}
			}
		}

		if message.Action == NEW_POSTIT {
			u2, err := uuid.NewV4()
			if err != nil {
				log.Fatalf("failed to generate UUID: %v", err)
			}
			message.ID = u2.String()
			// TODO : add authorId to postit
			for i, board := range boards {
				if board.ID == message.BoardId {
					boards[i].Postits = append(boards[i].Postits, Postit{ID: message.ID, BoardID: message.BoardId, Text: message.Text, PosX: message.PosX, PosY: message.PosY})
				}
			}
		}

		if message.Action == UPDATE_CONTENT {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.ID == message.ID {
							boards[i].Postits[j].Text = message.Text
						}
					}
				}
			}
		}

		if message.Action == MOVE_POSTIT {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.ID == message.ID {
							boards[i].Postits[j].PosX = message.PosX
							boards[i].Postits[j].PosY = message.PosY
						}
					}
				}
			}
		}

		if message.Action == ADD_VOTE {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.ID == message.ID {
							boards[i].Postits[j].Votes = board.Postits[j].Votes + 1
						}
					}
				}
			}
		}

		if message.Action == REMOVE_VOTE {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.ID == message.ID {
							boards[i].Postits[j].Votes = board.Postits[j].Votes - 1
						}
					}
				}
			}
		}

		if message.Action == SHOW_POSTITS {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.Author.ID == "TODO" {
							boards[i].Postits[j].Show = true
						}
					}
				}
			}
		}

		if message.Action == HIDE_POSTITS {
			for i, board := range boards {
				if board.ID == message.BoardId {
					for j, postit := range board.Postits {
						if postit.Author.ID == "TODO" {
							boards[i].Postits[j].Show = false
						}
					}
				}
			}
		}

		// Encodez la structure en JSON pour la diffusion
		msgJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("Erreur lors de l'encodage du message :", err)
			continue
		}

		// Diffuser le message à tous les autres clients
		hub.broadcast <- msgJSON
	}
}

type Board struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Postits []Postit `json:"postits"`
}

type Postit struct {
	ID      string `json:"id"`
	BoardID string `json:"boardId"`
	Text    string `json:"text"`
	PosX    int    `json:"posX"`
	PosY    int    `json:"posY"`
	Votes   int    `json:"votes"`
	Show    bool   `json:"show"`
	Author  User   `json:"author"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var boards = []Board{}

func getBoard(id string) (*Board, error) {
	for _, board := range boards {
		if board.ID == id {
			return &board, nil
		}
	}
	return nil, errors.New("Board not found")
}

func main() {
	// Récupérer la clé publique de Keycloak
	pubKey, err := getKeycloakPublicKey()
	if err != nil {
		log.Fatal("Erreur lors de la récupération de la clé publique :", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Changez cela pour restreindre les domaines autorisés
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Utiliser le middleware pour vérifier les tokens sur tous les endpoints protégés
	app.Use("/api", jwtMiddleware(pubKey))

	// Initialiser le WebSocketHub
	hub := NewWebSocketHub()

	// Lancer la routine pour gérer la diffusion
	go hub.run()

	// Endpoint WebSocket sécurisé
	app.Get("/ws", func(c *fiber.Ctx) error {
		token := c.Query("token") // Récupérer le token de la chaîne de requête

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Token manquant")
		}

		// Vérifier le token JWT
		_, err := validateToken(token, pubKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Token invalide")
		}

		// Passer la requête au handler WebSocket
		return websocket.New(func(ws *websocket.Conn) {
			websocketHandler(ws, hub)
		})(c)
	})

	// Endpoint pour récupérer la liste des boards
	app.Get("/api/boards", func(c *fiber.Ctx) error {
		return c.JSON(boards)
	})

	// Endpoint pour récupérer la liste des postits par boardId
	app.Get("/api/boards/:id/postits", func(c *fiber.Ctx) error {
		boardId := c.Params("id")
		selectedBoard, err := getBoard(boardId)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Board non trouvé")
		}
		postits := selectedBoard.Postits

		return c.JSON(postits)
	})

	log.Fatal(app.Listen(":3000"))
}
