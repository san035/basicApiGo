package test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/san035/basicApiGo/internal/config"
	"github.com/san035/basicApiGo/internal/db"
	"github.com/san035/basicApiGo/internal/router"
	"github.com/san035/basicApiGo/pkg/logger"
	"github.com/san035/basicApiGo/pkg/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math/rand"
	"net/http/httptest"
	"path"
	"runtime"
	"testing"
	"time"
)

const FormatData = `02.01.2006`

type ItemTest struct {
	Name               string
	Method             string
	Route              string
	Headers            map[string]string
	Body               interface{}
	ExpectedStatusCode int
	ExpectedBody       interface{}
	FuncExpectedBody   func(respBody *[]byte) (respBodyTrue string, err error) // возвращает расчетный expectedBody
	FuncBefore         func(*ItemTest, *testing.T) error                       // Init функция, например подготовка в БД новых записей
	FuncEnd            func(*ItemTest, *testing.T) error                       // Завершающая функция, например удаление из БД созданных записей
	SkipTest           bool                                                    // Признак пропускать тест
}

var AppFiber = fiber.New()

func init() {
	var err error
	defer func(err *error) {
		if *err != nil {
			log.Fatal().Err(*err).Msg("test.init-")
		}
	}(&err)

	//добавление всех маршрутов
	router.InitEndPoint(AppFiber)

	// Инициализация логгера
	logger.Init("debug")

	// Загрузка настроек из env
	err = config.LoadConfig()
	if err != nil {
		return
	}

	// Загрузка RSA ключей
	err = setTestJWT(&config.Config.JWT)
	if err != nil {
		return
	}
	err = token.Init(&config.Config.JWT)
	if err != nil {
		return
	}

	// Тест всех uri БД
	err = db.Init(&config.Config.DB)
	if err != nil {
		log.Error().Err(err).Msg("db.Init-")
	}

}

// DoListTest Общий исполнитель cписка тестов
func DoListTest(t *testing.T, listTests *[]ItemTest) {
	for idTest, testReq := range *listTests {
		if testReq.SkipTest {
			continue
		}

		// запуск init функции
		if testReq.FuncBefore != nil {
			err := testReq.FuncBefore(&testReq, t)
			if err != nil {
				t.Error(err)
			}
		}

		bodyReq, err := json.Marshal(testReq.Body)
		t.Logf("\n\n%d. Начало теста '%s' %s %s Body: %s ожидаемы статус: %d ожидаемый ответ: %s", idTest, testReq.Name, testReq.Method, testReq.Route, string(bodyReq), testReq.ExpectedStatusCode, testReq.ExpectedBody)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(testReq.Method, testReq.Route, bytes.NewReader(bodyReq))

		// Установка заголовков
		for k, v := range testReq.Headers {
			req.Header.Add(k, v)
		}

		// Выполнение запроса
		resp, err := AppFiber.Test(req, 50000)
		if err != nil {
			t.Fatal(err)
		}

		// Сверка статуса ответа с ожидаемым
		assert.Equal(t, testReq.ExpectedStatusCode, resp.StatusCode, "Cтатус ответа")

		// Получение body ответа
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		// Заполнение ожидаемого body
		if testReq.FuncExpectedBody != nil {
			testReq.ExpectedBody, err = testReq.FuncExpectedBody(&respBody)
			if err != nil {
				t.Fatal(err)
			}
		}

		// Сверка body ответа с ожидаемым
		switch expectedBody := testReq.ExpectedBody.(type) {
		case string:
			require.Equal(t, expectedBody, string(respBody), testReq.Name)
		case []string:
			require.Contains(t, expectedBody, string(respBody), testReq.Name)
		default:
			t.Fatal("Не известный тип body: ", testReq.ExpectedBody)
		}
		//t.Logf("body %t", ok)

		// запуск завершающей функции
		if testReq.FuncEnd != nil {
			err = testReq.FuncEnd(&testReq, t)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

// setTestJWT установка тестовых rsa_key файлов
func setTestJWT(config *config.JWTConfig) error {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return logger.New("Ошибка получения текущей папки тестов")
	}

	folderJWT := path.Join(path.Dir(file), `rsa_key`)

	config.PrivateKeyFile = path.Join(folderJWT, `jwt_privat_key_rsa`)
	config.PublicKeyFile = path.Join(folderJWT, `jwt_public_key_rsa`)
	return nil
}

// GetRandomDay - случайный день строкой
func GetRandomDay() string {
	return time.Unix(rand.Int63n(32503680000), 0).Format(FormatData)
}

// GetStampMilli - возвращает секунды + милисекунды
// пример возврата "10:24:10.433"
func GetStampMilli() string {
	return time.Now().Format(`15:04:05.000`)
}

// GetTimeNow текущее время строкой
func GetTimeNow() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}

// получение времени в
//var tv syscall.Timeval
//syscall.Gettimeofday(&tv)
//return strconv.Itoa(int(tv.Sec)) + ":"  + strconv.Itoa(int(tv.Usec)/1e3))

//now := time.Now().UnixMilli()
//localTime := time.UnixMilli() // .SecondsToLocalTime(now / 1e9)
//miliSeconds := (now % 1e9) / 1e6
//return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%03d", localTime.Year, localTime.Month, localTime.Day, localTime.Hour, localTime.Minute, localTime.Second, miliSeconds)
