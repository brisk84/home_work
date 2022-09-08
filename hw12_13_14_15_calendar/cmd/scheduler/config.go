package main

type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Rabbit   RabbitConf
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

type RabbitConf struct {
	URL string
}

type AppConf struct {
	CheckInterval string `mapstructure:"check_interval"`
}

func NewConfig() Config {
	return Config{}
}
