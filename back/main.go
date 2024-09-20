package main

import (
	"log"

	"github.com/code-troopers/postitsonline/database"
	"github.com/code-troopers/postitsonline/handlers"
	"github.com/code-troopers/postitsonline/webtoken"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Erreur lors de la connexion à la base de données :", err)
	}
	defer database.CloseDB()

	pubKey, err := webtoken.GetKeycloakPublicKey()
	if err != nil {
		log.Fatal("Erreur lors de la récupération de la clé publique :", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Changez cela pour restreindre les domaines autorisés
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use("/api", webtoken.Middleware(pubKey))

	hub := handlers.NewWebSocketHub()
	go hub.Run()

	app.Get("/ws", func(c *fiber.Ctx) error {
		token := c.Query("token")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Token manquant")
		}

		_, err := webtoken.ValidateToken(token, pubKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Token invalide")
		}

		// Save the user data
		user, err := webtoken.DecodeJWT(token)
		if err != nil {
			log.Println("Erreur lors du décodage du token :", err)
			return err
		}
		err = handlers.CreateUser(user)
		if err != nil {
			log.Println("Erreur lors de la création de l'utilisateur :", err)
			return err
		}

		return websocket.New(func(ws *websocket.Conn) {
			handlers.WebsocketHandler(ws, hub)
		})(c)
	})

	app.Get("/api/boards", func(c *fiber.Ctx) error {
		boards, err := handlers.GetAllBoards()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la récupération des boards")
		}
		if boards == nil {
			return c.JSON([]database.Board{})
		}

		return c.JSON(boards)
	})

	app.Get("/api/boards/:id/postits", func(c *fiber.Ctx) error {
		boardId := c.Params("id")
		postits, err := handlers.GetAllPostitsByBoardId(boardId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Erreur lors de la récupération des boards")
		}

		if postits == nil {
			return c.JSON([]database.Postit{})
		}

		return c.JSON(postits)
	})

	log.Fatal(app.Listen(":3000"))
}
