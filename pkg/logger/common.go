package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
)

// Recover Восстановление при падении
func Recover() {
	r := recover()
	if r != nil {
		err := errors.New(fmt.Sprint(r))
		log.Error().Stack().Err(err).Msg("Recovered")
	}
}

// EndMain вызов в main.go:
// defer  logger.EndMain(&err)
func EndMain(err *error) {
	//Восстановление при падении
	Recover()

	// ошибку в лог
	if *err != nil {
		logMsg := log.Error().Err(*err)
		logMsg = AddAddition(err, logMsg)
		logMsg.Msg("finish main-")
		os.Exit(1)
		return
	}

	log.Info().Msg("finish main+")
}
