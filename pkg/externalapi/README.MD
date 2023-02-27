# Пример использования externalapi
```
package apistorage

import (
	"github.com/san035/basicApiGo/pkg/externalapi"
	"github.com/san035/basicApiGo/pkg/logger"
	"net/http"
	"strings"
)

const RoutePhotoId = "/storage/profile/photo/"

type ResponsePhoto struct {
	Url   string
	Error string
}

// GetNewURLPhoto получение свежей ссылки на фото
func GetNewURLPhoto(token, idUser *string) (urlPhoto string, err error) {

	method := http.MethodGet
	uri := RoutePhotoId + *idUser
	r := &ResponsePhoto{}
	body := new(strings.Reader)
	err = externalapi.CallApi(Config.ListApiStorage, method, &uri, token, body, r)
	if err != nil {
		return
	}

	if r.Error != "" {
		err = logger.New(r.Error + ", " + method + ": " + uri)
		return
	}

	return r.Url, nil
}
```