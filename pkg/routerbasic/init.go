package routerbasic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"strconv"
)

var AppFiber = fiber.New(fiber.Config{DisableStartupMessage: true})

// Запуск web сервера
func StartWebServer(listPort *[]int) (err error) {

	ptotocol := "http"

	// Настройка роутинга
	InitCors(AppFiber)

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
