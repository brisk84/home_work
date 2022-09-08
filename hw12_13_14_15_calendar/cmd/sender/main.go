package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/ver"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	// flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		ver.PrintVersion()
		return
	}

	cfgPath, cfgFile := filepath.Split(configFile)
	cfgFile = strings.TrimSuffix(cfgFile, filepath.Ext(cfgFile))
	fmt.Println(cfgPath, " - ", cfgFile)

	viper.SetConfigName(cfgFile)
	viper.AddConfigPath(cfgPath)

	cfg := NewConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logg := logger.New(cfg.Logger.Path, cfg.Logger.Level)

	rabb := rabbit.NewRabbit(cfg.Rabbit.URL)
	err = rabb.Connect()
	if err != nil {
		logg.Error(err.Error())
		return
	}

	ch, err := rabb.Get()
	if err != nil {
		logg.Error(err.Error())
		return
	}

	for d := range ch {
		log.Printf("Received a message: %s\n", d.Body)
	}
}
