package router

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/internal/config"
	"github.com/san035/basicApiGo/pkg/logger"
	"strconv"
)

// Запуск web сервера
func Init() (err error) {

	ptotocol := "https"
	if config.Config.API.CertFileHtpps == "" || config.Config.API.KeyFileHtpps == "" {
		ptotocol = "http"
	}

	// Настройка конвертора данных запросов
	//InitFiberConverter()

	// Настройка роутинга
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	initCors(app)
	InitEndPoint(app)

	// Запупуск Listen на свободном порту
	for _, freePort := range config.Config.API.ListPort {
		freePorStr := strconv.Itoa(freePort)
		if ptotocol == "http" {
			log.Info().Str("Port", freePorStr).Str("ptotocol", ptotocol).Msg("Start API server " + ptotocol + "://127.0.0.1:" + freePorStr + " +")
			err = app.Listen(":" + freePorStr)
			if err == nil {
				break
			}
			log.Info().Err(err).Str("Port", freePorStr).Msg("Порт занят-")
			continue
		} else {
			err = errors.New(ErrorUnsupportedProtocol + ptotocol)
		}
	}

	if err != nil {
		err = logger.Wrap(&err)
	}
	return err
}

// InitEndPoint добавление всех маршрутов
func InitEndPoint(app *fiber.App) {
	//Информация о микросервисе
	app.Get("/", Stat)
	app.Get("/stat/", Stat)
}

func initCors(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		//AllowHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, X-Secret, Access-Control-Allow-Origin, Access-Control-Allow-Headers, X-CSRF-Token, authorization-token",
	}))
}
