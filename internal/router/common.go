// функции используемые в разных endpoint
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/san035/basicApiGo/internal/userclass"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/token"
)

// GetUserByTokenRequest возвращает user по токену из заголовка
func GetUserByTokenRequest(ctx *fiber.Ctx) (user *userclass.User, err error) {
	tokenString := ctx.Get("Authorization", "")

	// Валидация токена
	mapClaims, err := token.Validate(&tokenString)
	if err != nil {
		return
	}

	// Проверка данных в map и сохранение в user
	var valueString string
	var valueFloat64 float64
	var ok bool
	user = new(userclass.User)
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
			user.ID = valueString
		case key == "email":
			user.Email = valueString
		case key == "role":
			user.Role = valueString
		case key == "exp":
			user.Exp = int64(valueFloat64)
		}
	}
	return
}
