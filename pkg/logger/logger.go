/*
Важнор чтобы в import был не просто "errors", а	"github.com/pkg/errors"
тогда при вызове log.Error().Err(err).Stack().Msg("*") в лог попадет стек вызовов
*/

package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func Init(levelStr string) {
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack // для log.Error().Stack()

	levelLog := UpdateLevelLog(levelStr)
	log.Info().Str("Level", levelLog.String()).Strs("Args", os.Args).Msg("Start log")
}

// Установка уровня логгирования по config.Config.API.LevelLog
func UpdateLevelLog(levelStr string) zerolog.Level {
	//Установка уровня логирования
	levelLog, err := zerolog.ParseLevel(levelStr)
	if err != nil || levelLog == zerolog.NoLevel {
		levelLog = zerolog.InfoLevel
	}
	globalLevel := zerolog.GlobalLevel()
	if globalLevel != levelLog && levelLog != zerolog.NoLevel {
		zerolog.SetGlobalLevel(levelLog)
		log.Info().Str("OldLevel", globalLevel.String()).Str("NewLevel", levelLog.String()).Msg("UpdateLevelLog")
	}
	return levelLog
}
