package main

import (
	"github.com/san035/basicApiGo/internal/config"
	"github.com/san035/basicApiGo/internal/router"
	"github.com/san035/basicApiGo/pkg/common"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/osutils"
	"github.com/san035/basicApiGo/pkg/routerbasic"
	"github.com/san035/basicApiGo/pkg/token"
)

func main() {
	var err error

	// Инициализация логгера
	logger.Init("info")

	// Логирование при завершении
	defer logger.EndMain(&err)

	// Закрытие ранее открытых приложений
	osutils.KillDoubleApp()

	// Проверка аргументов
	osutils.ChaeckArg()

	// Загрузка настроек
	err = common.LoadConfig(&config.Config)
	if err != nil {
		return
	}

	// формирование полных путей
	osutils.AddPathApp(&config.Config.JWT.PrivateKeyFile)
	osutils.AddPathApp(&config.Config.JWT.PublicKeyFile)

	// Обновление уровня логгирования, после загрузки конфигурации
	logger.UpdateLevelLog(config.Config.API.LevelLog)

	// Загрузка RSA ключей
	err = token.Init(&config.Config.JWT)
	if err != nil {
		return
	}

	// Запуск хостинга
	err = routerbasic.StartWebServer(&config.Config.API.ListPort, router.InitEndPoint)
}
