package dbtarantool

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/tarantool/go-tarantool"
)

var (
	conn      *tarantool.Connection
	lastIdUri int // последний номер uri БД
	listUri   *[]string
	opts      tarantool.Opts
)

func Init(config *DBConfig) error {
	listUri = &config.ListUri
	opts = tarantool.Opts{User: config.User, Pass: config.Password}
	return ConnectDB()
}

// Подключение к БД с перехватываением ощибки
func tryConnect() (conn *tarantool.Connection, err error) {
	defer func(err *error) {
		rec := recover()
		if rec == nil {
			return
		}
		*err = errors.New(ErrorRecoverTarantoolConnect)
		return
	}(&err)
	return tarantool.Connect((*listUri)[lastIdUri], opts)
}

// ConnectDB подключение к БД по списку config.Config.DB.ListUri
func ConnectDB() (err error) {
	if conn == nil || !conn.ConnectedNow() {
		count_uri := len(*listUri)
		for i := 0; i < count_uri; i++ {

			conn, err = tryConnect()
			if err == nil {
				log.Info().Int("uriID", i).Str("uri", (*listUri)[lastIdUri]).Msg("ConnectDB+")
				return
			}

			err = logger.Wrap(&err)
			log.Error().Err(err).Str("User", opts.User).Str("uri", (*listUri)[lastIdUri]).Msg("ConnectDB-")
			if count_uri > 1 {
				lastIdUri = (lastIdUri + 1) % count_uri
			}
		}
	}
	return
}
