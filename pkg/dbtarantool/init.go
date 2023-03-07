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

func Init(config *DBConfig) (err error) {
	listUri = &config.ListUri
	opts = tarantool.Opts{User: config.User, Pass: config.Password}
	conn, err = ConnectDB()
	return err
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
func ConnectDB() (connTarantool *tarantool.Connection, err error) {
	if conn != nil && conn.ConnectedNow() {
		return conn, nil
	}

	countUri := len(*listUri)
	for i := 0; i < countUri; i++ {

		conn, err = tryConnect()
		if err == nil {
			log.Info().Int("uriID", i).Str("uri", (*listUri)[lastIdUri]).Msg("ConnectDB+")
			return conn, nil
		}

		err = logger.Wrap(&err)
		log.Error().Err(err).Str("User", opts.User).Str("uri", (*listUri)[lastIdUri]).Msg("ConnectDB-")
		if countUri > 1 {
			lastIdUri = (lastIdUri + 1) % countUri
		}
	}
	return
}
