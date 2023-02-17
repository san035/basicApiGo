package routerbasic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		//AllowHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, X-Secret, Access-Control-Allow-Origin, Access-Control-Allow-Headers, X-CSRF-Token, authorization-token",
	}))
}
