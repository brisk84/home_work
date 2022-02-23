package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/ver"
	"github.com/spf13/viper"
)

var configFile string

func init() {
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	stor := sqlstorage.New(cfg.Database.DBType, cfg.Database.ConnStr, cfg.Database.MaxConns)
	err = stor.Connect(ctx)
	if err != nil {
		logg.Error("Can't connect to dabatase")
		return
	}

	viper.BindEnv("RABBIT_URL")
	rabURL := viper.Get("RABBIT_URL")
	rabbitURL := cfg.Rabbit.URL
	if rabURL != nil {
		rabbitURL = rabURL.(string)
	}

	rabb := rabbit.NewRabbit(rabbitURL)
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
		ev := storage.Event{}
		err := json.Unmarshal(d.Body, &ev)
		if err != nil {
			logg.Error(err.Error())
			continue
		}
		ev.Description = "!Notified!"
		stor.EditEvent(ctx, ev)
	}
}
