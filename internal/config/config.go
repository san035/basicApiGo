package config

import (
	"github.com/san035/basicApiGo/pkg/token"
)

type DBConfig struct {
	ListUri  []string `env:"DB_LIST_URI" default:"[\"bpm.dev.itkn.ru:3301\"]"` // Список uri БД через запятую, пример env DB_LIST_URI=["bpm.dev.itkn.ru:3301"]
	User     string   `env:"DB_USER" default:"user_service" yaml:"User"`       // Пользователь БД
	Password string   `env:"DB_PASSWORD" default:"hgsFy23_jW" yaml:"Password"` // Пароль пользователя
}

var Config = struct {
	API struct {
		ListPort      []int  `default:"[8091,8092,8093]" env:"API_LIST_PORT" yaml:"ListPort"` // порт микросервиса
		CertFileHtpps string `default:"" env:"API_CERT_FILE_HTPPS" yaml:"CertFileHtpps"`      // Ссылка на файл cert.pem, при отсутствии протокол http
		KeyFileHtpps  string `default:"" env:"API_KEY_FILE_HTPPS" yaml:"KeyFileHtpps"`        // Ссылка на файл key.pem, при отсутствии протокол http
		LevelLog      string `default:"info" env:"API_LEVELlOG" yaml:"LevelLog"`              // Режим логгирования
	}
	DB  DBConfig
	JWT token.JWTConfig
}{}
