package common

import (
	"fmt"
	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"os"
	"path/filepath"
)

const NAME_FILE_CONFIG_YAML = "config.yml"

// GetFilesConfig возвращает список с файлом конфигурации
func GetFilesConfig() (filesConfigYaml []string) {
	fileConfig := filepath.Dir(os.Args[0]) + string(filepath.Separator) + NAME_FILE_CONFIG_YAML
	if _, err := os.Stat(fileConfig); err == nil {
		filesConfigYaml = append(filesConfigYaml, fileConfig)
		log.Debug().Str("file", fileConfig).Msg("Найден файл " + NAME_FILE_CONFIG_YAML)
	} else {
		log.Debug().Str("file", fileConfig).Msg("Не найден файл " + NAME_FILE_CONFIG_YAML)
	}
	return
}

// LoadConfig загрузка config.yml и env
func LoadConfig(conf interface{}) (err error) {

	//поиск файла config.yml
	filesConfigYaml := GetFilesConfig()

	// загрузка
	err = configor.Load(conf, filesConfigYaml...)
	if err != nil {
		err = logger.Wrap(&err, "config.LoadConfig")
		return
	}

	log.Debug().Str("Config", fmt.Sprintf("%#v", conf)).Msg("Load configs+")
	return
}
