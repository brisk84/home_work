package main

type Config struct {
	Logger LoggerConf
	Rabbit RabbitConf
}

type LoggerConf struct {
	Level string
	Path  string
}
type RabbitConf struct {
	URL string
}

func NewConfig() Config {
	return Config{}
}
