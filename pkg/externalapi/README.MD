# Пример использования externalapi
```
package apiuser

import (
	"fmt"
	"github.com/san035/basicApiGo/pkg/externalapi"
	"github.com/san035/basicApiGo/pkg/logger"
	"net/http"
	"net/url"
)

type ResponseProfileEdit struct {
	Status string
	Error  string
}

/*
Ответы:
{"Error": "invalid character '{' in literal true (expecting 'r')"}
{"Status":"ok"}
*/

const RouteEdit = "/profile/edit/"

// UpdateUrlPhotoInDb вызов api PATCH "/profile/edit/" c новым url в body
func UpdateUrlPhotoInDb(token, idUser *string, urlPhoto *url.URL) error {

	method := http.MethodPatch
	uri := RouteEdit + *idUser
	body := fmt.Sprintf("{\"Photo\":\"%s\"}", urlPhoto.String())

	r := &ResponseProfileEdit{}
	err := externalAPI.CallApi(&Config.ListApiUser, method, &uri, token, &body, r)
	if err != nil {
		return err
	}

	if r.Status != "ok" {
		return logger.New(r.Error + ", " + method + ": " + uri)
	}

	return nil

}
```