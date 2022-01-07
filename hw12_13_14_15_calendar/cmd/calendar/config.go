package main

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Server   ServerConf
	App      AppConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type DatabaseConf struct {
	DBType   string `mapstructure:"db_type"`
	ConnStr  string `mapstructure:"conn_str"`
	MaxConns int    `mapstructure:"max_conns"`
}

type ServerConf struct {
	Host string
	Port string
}

type AppConf struct {
	Storage string
}

func NewConfig() Config {
	return Config{}
}
