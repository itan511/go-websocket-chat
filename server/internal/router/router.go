package router

import (
	"go-websocket-chat/server/internal/user"
	"go-websocket-chat/server/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var app *fiber.App

func InitRouter(userHandler *user.Handler, wsHandler *ws.Handler) {
	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST",
		AllowHeaders:     "Content-Type",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	app.Post("/signup", userHandler.CreateUser)
	app.Post("/login", userHandler.Login)
	app.Get("/logout", userHandler.Logout)

	app.Post("/ws/createRoom", wsHandler.CreateRoom)
	app.Get("/ws/joinRoom/:roomId", wsHandler.JoinRoom)
	app.Get("/ws/getRooms", wsHandler.GetRooms)
	app.Get("/ws/getClients/:roomId", wsHandler.GetClients)
}

func Start(addr string) error {
	return app.Listen(addr)
}
