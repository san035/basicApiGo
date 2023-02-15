package common

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"reflect"
)

// Заполнение структуры valueStruct значениями из тега default
func SetValueDefaultByReflect(valueStruct interface{}) {
	itemsType := reflect.TypeOf(valueStruct).Elem()
	itemsValue := reflect.ValueOf(valueStruct)
	for i := 0; i < itemsType.NumField(); i++ {
		itemValue := itemsValue.Elem().Field(i)
		if itemValue.CanSet() {
			newValue := itemsType.Field(i).Tag.Get("default")
			itemValue.SetString(newValue)
		}
	}
	return
}

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
