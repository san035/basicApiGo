package routerbasic

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/token"
	"github.com/san035/basicApiGo/pkg/userclass"
)

const KeyContextUserToken = "UserToken"

// SendReeplyOk Отправка ok в ответ запроса
func SendReeplyOk(ctx *fiber.Ctx) error {
	err := ctx.Send([]byte(`{"Status":"ok"}`))
	if err != nil {
		err = logger.WrapWithDeep1(&err)
	}
	return err
}

// Запись в лог события о запросе
func EndRequest(ctx *fiber.Ctx, err *error) {
	//Восстановление при падении
	logger.Recover()

	c := ctx.UserContext()
	addFields := c.Value("logAddition")

	// Запись в лог
	if *err == nil {
		log.Info().Str("IP", ctx.IP()).Str(ctx.Method(), ctx.OriginalURL()).Fields(addFields).Msg("request+")
		return
	}

	//Строка лога
	logData := logger.Error(err).Str("IP", ctx.IP()).Str(ctx.Method(), ctx.OriginalURL())
	if addFields != nil {
		logData = logData.Fields(addFields)
	}
	switch errExt := (*err).(type) {
	case *logger.ExtendedError:
		if errExt.Addition != nil {
			logData = logData.Interface("ext", errExt.Addition)
		}
	}
	logData.Msg("request-")

	// Возврат ошибки
	oldStatusCode := ctx.Context().Response.StatusCode()
	if oldStatusCode == fiber.StatusOK {
		ctx.Status(fiber.StatusInternalServerError)
	}

	// Отправка ответа
	errJSON := struct{ Error string }{Error: (*err).Error()}
	*err = nil // очищаем, чтобы сработал ctx.JSON
	err2 := ctx.JSON(errJSON)
	if err2 != nil {
		logger.AddAddition(&err2, log.Error()).Msg("ctx.JSON-")
	}
}

// GetUserByTokenRequest возвращает user по токену из заголовка
func GetUserByTokenRequest(ctx *fiber.Ctx) (userToken *userclass.User, err error) {
	tokenString := ctx.Get("Authorization", "")

	// Валидация токена
	mapClaims, err := token.Validate(&tokenString)
	if err != nil {
		return
	}

	// Проверка данных в map и сохранение в userToken
	var valueString string
	var valueFloat64 float64
	var ok bool
	userToken = new(userclass.User)
	for _, key := range [4]string{"id", "email", "role", "exp"} {
		if key == "exp" {
			valueFloat64, ok = mapClaims[key].(float64)
		} else {
			valueString, ok = mapClaims[key].(string)
		}
		switch {
		case !ok:
			err = logger.New(token.NotExistKeyInToken + key)
			return
		case key == "id":
			userToken.ID = valueString
		case key == "email":
			userToken.Email = valueString
		case key == "role":
			userToken.Role = valueString
		case key == "exp":
			userToken.Exp = int64(valueFloat64)
		}
	}

	// Дополнительные поля лога
	SetLogAddition(ctx, "Email", userToken.Email)
	return
}

// SetLogAddition Добавление в переменную контекста logAddition доп информации
func SetLogAddition(ctx *fiber.Ctx, key string, value interface{}) {
	c := ctx.UserContext()
	mapLogAddition := c.Value("logAddition")
	switch mapLog := mapLogAddition.(type) {
	case map[string]interface{}:
		mapLog[key] = value // Сохранение для лога
		c = context.WithValue(c, "logAddition", mapLog)
	default:
		mapLog = map[string]interface{}{key: value}
		c = context.WithValue(c, "logAddition", mapLog)
	}
	ctx.SetUserContext(c)
	return
}
