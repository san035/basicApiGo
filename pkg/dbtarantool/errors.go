package dbtarantool

import "errors"

const (
	UserNotExist                 = "Пользователь не найден"
	ParseDataUser                = "Ошибка парсинга данных"
	UserAlreadyExist             = "Пользователь уже существует"
	ErrorRecoverTarantoolConnect = "recover tarantool connect"
)

var (
	ErrorUserNotExist     = errors.New(UserNotExist)
	error_parse_data_user = errors.New(ParseDataUser)
)
