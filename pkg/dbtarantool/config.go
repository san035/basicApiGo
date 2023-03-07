package dbtarantool

type DBConfig struct {
	ListUri  []string `env:"DB_LIST_URI" default:"[\"bpm.dev.itkn.ru:3301\"]"` // Список uri БД через запятую, пример env DB_LIST_URI=["bpm.dev.itkn.ru:3301"]
	User     string   `env:"DB_USER" default:"user_service" yaml:"User"`       // Пользователь БД
	Password string   `env:"DB_PASSWORD" default:"hgsFy23_jW" yaml:"Password"` // Пароль пользователя
}
