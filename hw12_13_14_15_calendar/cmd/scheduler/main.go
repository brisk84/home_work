package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/rabbit"
	sqlstorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	// flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
	flag.StringVar(&configFile, "config", "configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
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

	checkInterval, err := time.ParseDuration(cfg.App.CheckInterval)
	if err != nil {
		log.Fatalf("unable to parse checking interval, %v", err)
	}

	logg.Info("Check interval is: " + strconv.Itoa(int(checkInterval)))

	ticker := time.NewTicker(checkInterval)

	stor := sqlstorage.New(cfg.Database.DBType, cfg.Database.ConnStr, cfg.Database.MaxConns)
	err = stor.Connect(context.TODO())
	if err != nil {
		logg.Error("Can't connect to dabatase")
		return
	}

	rabb := rabbit.NewRabbit(cfg.Rabbit.URL)
	err = rabb.Connect()
	if err != nil {
		logg.Error(err.Error())
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Completed !")
			return
		case <-ticker.C:
			notifyDate := time.Date(time.Now().Year(), time.Now().Month(),
				time.Now().Day(), 0, 0, 0, 0, time.Local)
			ev, er := stor.GetNotifyEvent(notifyDate)
			fmt.Println(ev, er)

			msg, err := json.Marshal(ev)
			if err != nil {
				logg.Error(err.Error())
				return
			}
			rabb.Send(msg)
		}
	}
}
