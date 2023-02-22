package externalAPI

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"net/http"
	"strings"
)

// CallApi вызов внешнего API
func CallApi(listApi *[]ExternalAPI, method string, uri, token, body *string, dataResp interface{}) (err error) {
	// перебор api
	for idUri, apiConfig := range *listApi {

		// запрос api
		urlReq := apiConfig.Uri + *uri
		var req *http.Request
		req, err = http.NewRequest(method, urlReq, strings.NewReader(*body))
		if err != nil {
			log.Error().Err(err).Str("Method", req.Method).Str(req.Method, urlReq).Msg("api-")
			continue
		}

		// Установка заголовка и выполнение запроса
		req.Header.Set("Content-Type", `application/json`)
		req.Header.Set("Authorization", *token)
		var resp *http.Response
		resp, err = (&http.Client{}).Do(req)
		if err != nil {
			err = logger.Wrap(&err)
			log.Error().Err(err).Str(req.Method, urlReq).Msg("api-")
			continue
		}

		// перемещаем успешное uri на первое место в списке
		if idUri != 0 {
			(*listApi)[0], (*listApi)[idUri] = (*listApi)[idUri], (*listApi)[0]
		}

		// Парсинг ответа
		err = json.NewDecoder(resp.Body).Decode(&dataResp)
		if err != nil {
			return
		}

		return
	}
	return
}
