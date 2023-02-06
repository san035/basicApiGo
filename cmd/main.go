package main

import (
	"cmd/main.go/internal/config"
	"cmd/main.go/internal/router"
	"cmd/main.go/pkg/logger"
	"cmd/main.go/pkg/osutils"
	"cmd/main.go/pkg/token"
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
