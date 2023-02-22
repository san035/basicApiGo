package routerbasic

import (
	"github.com/gofiber/fiber/v2"
	"github.com/san035/basicApiGo/internal/config"
	"github.com/san035/basicApiGo/pkg/userclass"
	"github.com/shirou/gopsutil/host"
	"os"
	"runtime"
	"time"
)

var startApp time.Time

func init() {
	startApp = time.Now()
}

// Stat информацию о микросервисе
func Stat(ctx *fiber.Ctx) (err error) {
	// Логгирование и перехват фатальных ошибок
	defer AddRequestToLog(ctx, &err, nil)

	mapaboutAPI := map[string]interface{}{}

	//Дата созданя API
	var fileInfo os.FileInfo
	fileInfo, err = os.Stat(os.Args[0])
	if err != nil {
		return err
	}
	mapaboutAPI["Дата созданя API:"] = fileInfo.ModTime()
	mapaboutAPI["Дата запуска API:"] = startApp
	mapaboutAPI["Параметры запуска:"] = os.Args

	mapaboutAPI["Горутин:"] = runtime.NumGoroutine()
	mapaboutAPI["Версия компилятора:"] = runtime.Version()

	// Занято памяти
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	mapaboutAPI["Память:"] = struct {
		Alloc, HeapAlloc, TotalAlloc uint64
	}{Alloc: memStat.Alloc, HeapAlloc: memStat.HeapAlloc, TotalAlloc: memStat.TotalAlloc}

	//Время запуска ОС
	timeOS, err := host.BootTime()
	if err != nil {
		return err
	}
	mapaboutAPI["Время запуска ОС"] = time.Unix(int64(timeOS), 0)

	//Проверка токена
	user, err := GetUserByTokenRequest(ctx)
	if err != nil {
		mapaboutAPI["Check token"] = err.Error()
		ctx.Status(fiber.StatusUnauthorized)

	} else {
		mapaboutAPI["user token"] = struct {
			ID, Email, Role, Exp string
		}{ID: user.ID, Email: user.Email, Role: user.Role, Exp: time.Unix(user.Exp, 0).String()}

		// Инфо для админов
		if user.Role == userclass.RoleAdmin {
			mapaboutAPI["JWT_EXPIRES_MINUTES"] = config.Config.JWT.ExpiresMinutes
			mapaboutAPI["JWT_FILE_PUBLIC_KEY_RSA"], _ = os.ReadFile(config.Config.JWT.PublicKeyFile)
			mapaboutAPI["DB_LIST_URI"] = config.Config.DB.ListUri
		}
	}

	err = ctx.JSON(mapaboutAPI)

	return err
}
