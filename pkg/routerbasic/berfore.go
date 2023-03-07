package routerbasic

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/userclass"
)

// BerforEndpont общий роутер для всех запросов /
// Проверка токена + сохранение пользователя в контексте
func BerforEndpont(ctx *fiber.Ctx) (err error) {
	// Логгирование и перехват фатальных ошибок
	defer func(ctx *fiber.Ctx, err *error) {
		if err != nil {
			EndRequest(ctx, err)
		}
	}(ctx, &err)

	//Проверка токена
	var userToken *userclass.User
	userToken, err = GetUserByTokenRequest(ctx)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
		return logger.Wrap(&err)
	}

	// Сохранение userToken в контексте
	c := context.WithValue(ctx.UserContext(), KeyContextUserToken, userToken)
	ctx.SetUserContext(c)

	// Получение userToken из контекста:
	//c = ctx.UserContext()
	//userToken, ok := c.Value(routerbasic.KeyContextUserToken).(*userclass.User)
	//if !ok {
	//	err = logger.New("Не найден пользователь в контектсе")
	//	return
	//}

	// создание соединения с minio
	//err = storageminio.CreateConnect()
	//if err != nil {
	//	return
	//}

	err = ctx.Next()
	return
}
