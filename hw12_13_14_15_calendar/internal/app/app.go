package app

import (
	"context"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg *logger.Logger
	stor *storage.Calendar
}

// type Logger interface { // TODO
// }

// type Storage interface { // TODO
// }

func New(logger *logger.Logger, storage *storage.Calendar) *App {
	return &App{logg: logger, stor: storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
