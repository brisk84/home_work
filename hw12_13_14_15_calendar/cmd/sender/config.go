package main

type Config struct {
	Logger   LoggerConf
	Rabbit   RabbitConf
	Database DatabaseConf
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

type RabbitConf struct {
	URL string
}

func NewConfig() Config {
	return Config{}
}
