package osutils

import (
	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/process"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// KillDoubleApp Закрытие ранее открытых приложений
func KillDoubleApp() {
	selfNameApp := path.Base(strings.ReplaceAll(os.Args[0], "\\", "/"))
	selfPID := os.Getpid()

	processes, _ := process.Processes()
	for _, processApp := range processes {
		processName, _ := processApp.Name()
		if selfNameApp == processName && selfPID != int(processApp.Pid) {
			err := TerminateProcess(uint32(processApp.Pid))
			if err == nil {
				log.Info().Str("app", selfNameApp).Int32("pid", processApp.Pid).Msg("Kill old process+")
			} else {
				log.Error().Err(err).Str("app", selfNameApp).Int32("pid", processApp.Pid).Msg("Kill old process-")
			}
		}
	}
}

// Если в аргументах есть Down, то завершение работы
func ChaeckArg() {
	for i := 1; i < len(os.Args); i++ {
		if strings.ToLower(os.Args[i]) == "down" {
			log.Fatal().Str("reason", "Запуск с аргументом down").Msg("Завершение работы")
		}
	}
}

// AddPathApp добавление к file паки исполняемого файла
// формирование полных путей
func AddPathApp(file *string) {
	lenName := len(*file)
	if lenName == 0 {
		return
	}

	// Если есть : или  первый символ сепаратор, то ничего добавлять не надо
	if strings.Index(*file, ":") >= 0 {
		return
	}

	if os.IsPathSeparator((*file)[0]) {
		return
	}

	*file = filepath.Join(filepath.Dir(os.Args[0]), *file)
}
