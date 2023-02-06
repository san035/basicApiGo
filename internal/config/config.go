package config

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/osutils"
	"github.com/san035/basicApiGo/pkg/token"
	"os"
	"path/filepath"
)

const NAME_FILE_CONFIG_YAML = "config.yml"

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

// LoadConfig загрузка config.yml и env
func LoadConfig() (err error) {

	//поиск файла config.yml
	filesConfigYaml := []string{}
	fileConfig := filepath.Dir(os.Args[0]) + string(filepath.Separator) + NAME_FILE_CONFIG_YAML
	if _, err = os.Stat(fileConfig); err == nil {
		filesConfigYaml = append(filesConfigYaml, fileConfig)
		log.Debug().Str("file", fileConfig).Msg("Найден файл " + NAME_FILE_CONFIG_YAML)
	} else {
		log.Debug().Str("file", fileConfig).Msg("Не найден файл " + NAME_FILE_CONFIG_YAML)
	}

	// загрузка
	err = configor.Load(&Config, filesConfigYaml...)
	if err != nil {
		err = logger.Wrap(&err, "config.LoadConfig")
		return
	}

	// формирование полных путей
	osutils.AddPathApp(&Config.JWT.PrivateKeyFile)
	osutils.AddPathApp(&Config.JWT.PublicKeyFile)

	log.Debug().Str("Config", fmt.Sprintf("%#v", Config)).Msg("Load configs+")
	return
}
