package main

import (
	"github.com/san035/basicApiGo/internal/config"
	"github.com/san035/basicApiGo/internal/router"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/osutils"
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

	//// Загрузка настроек из env
	err = config.LoadConfig()
	if err != nil {
		return
	}

	// Обновление уровня логгирования, после загрузки конфигурации
	logger.UpdateLevelLog(config.Config.API.LevelLog)

	// Загрузка RSA ключей
	err = token.Init(&config.Config.JWT)
	if err != nil {
		return
	}

	// Запуск хостинга
	err = router.Init()
}
