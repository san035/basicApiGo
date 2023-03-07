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
var timeBuild string

func init() {
	startApp = time.Now()
}

// SetBuildApp установка времени компиляции приложения
func SetTimeBuildApp(newTimeBuild string) {
	timeBuild = newTimeBuild
}

// Stat информацию о микросервисе
func Stat(ctx *fiber.Ctx) (err error) {
	// Логгирование и перехват фатальных ошибок
	defer EndRequest(ctx, &err)

	mapaboutAPI := map[string]interface{}{}

	//Дата созданя API
	if timeBuild == "" {
		var fileInfo os.FileInfo
		fileInfo, err = os.Stat(os.Args[0])
		if err != nil {
			return err
		}
		timeBuild = fileInfo.ModTime().Format(time.RFC3339)
	}
	mapaboutAPI["Дата созданя API:"] = timeBuild
	mapaboutAPI["Дата запуска API:"] = startApp.Format(time.RFC3339)
	mapaboutAPI["Параметры запуска:"] = os.Args

	mapaboutAPI["Горутин:"] = runtime.NumGoroutine()
	mapaboutAPI["Версия компилятора:"] = runtime.Version()
	mapaboutAPI["pprof:"] = `/debug/pprof/`

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

	// Получение userToken из контекста:
	c := ctx.UserContext()
	userToken, ok := c.Value(KeyContextUserToken).(*userclass.User)
	if !ok {
		mapaboutAPI["Check token"] = ErrorNoFindUserInContext
		ctx.Status(fiber.StatusUnauthorized)
	} else {
		mapaboutAPI["userToken token"] = struct {
			ID, Email, Role, Exp string
		}{ID: userToken.ID, Email: userToken.Email, Role: userToken.Role, Exp: time.Unix(userToken.Exp, 0).String()}

		// Инфо для админов
		if userToken.Role == userclass.RoleAdmin {
			mapaboutAPI["JWT_EXPIRES_MINUTES"] = config.Config.JWT.ExpiresMinutes
			mapaboutAPI["JWT_FILE_PUBLIC_KEY_RSA"], _ = os.ReadFile(config.Config.JWT.PublicKeyFile)
			mapaboutAPI["DB_LIST_URI"] = config.Config.DB.ListUri
		}
	}

	err = ctx.JSON(mapaboutAPI)

	return err
}
