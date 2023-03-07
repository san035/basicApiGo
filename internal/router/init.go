package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/san035/basicApiGo/pkg/routerbasic"
)

// InitEndPoint добавление всех маршрутов
func InitEndPoint(app *fiber.App) {
	app.Use("/", routerbasic.BerforEndpont)

	app.Get("/", routerbasic.Stat)
	app.Get("/stat/", routerbasic.Stat)
}
