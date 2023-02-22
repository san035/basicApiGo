package externalapi

type ExternalAPI struct {
	Uri string `yaml:"Uri"`
}

// SetDefaultConfig установка uri по умолчанию
func SetDefaultConfig(listApiUser *[]ExternalAPI, UriDefault string) {

	if len(*listApiUser) != 0 {
		return
	}

	// Если не заданы настройки, задаем по умолчанию
	extApi := ExternalAPI{Uri: UriDefault}
	//common.SetValueDefaultByReflect(&extApi)
	listApiUser = &([]ExternalAPI{extApi})
	return
}
