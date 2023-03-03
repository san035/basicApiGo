package storageminio

import (
	"github.com/san035/basicApiGo/pkg/common"
)

type MinIOInstance struct {
	Uri             string `env:"MINIO_URI" default:"bpm.dev.itkn.ru:9000" yaml:"Uri"` // Список uri БД через запятую
	AccessKeyId     string `env:"MINIO_ACCESS_KEY_ID" default:"cRLI6CRKLqVtqHxg" yaml:"AccessKeyId"`
	SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY" default:"CRu3SGuVM5p7YXkinDdMWNk8y40TSeJC" yaml:"SecretAccessKey"`
}

var Config = struct {
	ListMinIO []MinIOInstance `yaml:"ListMinIO"`
}{}

func setDefaultConfig() {

	if len(Config.ListMinIO) != 0 {
		return
	}

	// Если не заданы настройки для minIO, задаем по умолчанию
	minIOInstance := MinIOInstance{}
	common.SetValueDefaultByReflect(&minIOInstance)

	Config.ListMinIO = []MinIOInstance{minIOInstance}
	return
}
