package router

import (
	"github.com/san035/basicApiGo/pkg/routerbasic"
)

// InitEndPoint добавление всех маршрутов
func InitEndPoint() {
	app := routerbasic.AppFiber
	app.Get("/", routerbasic.Stat)
	app.Get("/stat/", routerbasic.Stat)
}
