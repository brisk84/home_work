package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/server/http"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
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

	// viper.SetConfigName("config")
	// viper.AddConfigPath("../../configs/")

	cfg := NewConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	logg := logger.New(cfg.Logger.Path, cfg.Logger.Level)

	var stor storage.Calendar
	if cfg.App.Storage == "memory" {
		stor = memorystorage.New()
	} else {
		stor = sqlstorage.New(cfg.Database.DBType, cfg.Database.ConnStr, cfg.Database.MaxConns)
	}

	calendar := app.New(logg, &stor)
	server := internalhttp.NewServer(logg, calendar, net.JoinHostPort(cfg.Server.Host, cfg.Server.Port))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()
	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
