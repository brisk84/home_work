package app

import (
	"context"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg *logger.Logger
	stor storage.Calendar
}

func New(logger *logger.Logger, storage storage.Calendar) *App {
	return &App{logg: logger, stor: storage}
}

func (a *App) AddEvent(ctx context.Context, event storage.Event) error {
	return a.stor.AddEvent(event)
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	return a.stor.GetEvent(id)
}

func (a *App) EditEvent(ctx context.Context, event storage.Event) error {
	return a.stor.EditEvent(event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.stor.DeleteEvent(id)
}

func (a *App) ListEvents(ctx context.Context) ([]storage.Event, error) {
	return a.stor.ListEvents(), nil
}
