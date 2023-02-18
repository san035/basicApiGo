package routerbasic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"strconv"
)

var AppFiber = fiber.New(fiber.Config{DisableStartupMessage: true})

func init() {
	InitCors(AppFiber)
}

// StartWebServer Запуск web сервера
// funcInitEndPoint - функция в которой добавляются endpoint
func StartWebServer(listPort *[]int, funcInitEndPoint func(*fiber.App)) (err error) {

	// Настройка роутинга
	funcInitEndPoint(AppFiber)

	ptotocol := "http"

	// Запуск Listen на свободном порту
	for _, freePort := range *listPort {
		freePorStr := strconv.Itoa(freePort)
		log.Info().Str("Port", freePorStr).Str("ptotocol", ptotocol).Msg("Start API server " + ptotocol + "://127.0.0.1:" + freePorStr)
		err = AppFiber.Listen(":" + freePorStr)
		if err == nil {
			break
		}
		log.Info().Err(err).Str("Port", freePorStr).Msg("Порт занят-")
		continue
	}

	if err != nil {
		err = logger.Wrap(&err)
	}
	return err
}
