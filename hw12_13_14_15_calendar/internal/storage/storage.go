package storage

import (
	"context"
	"errors"
)

var (
	ErrDateBusy = errors.New("date is busy")
	ErrNotFound = errors.New("event not found")
	ErrUUIDBusy = errors.New("uuid is busy")
)

type Calendar interface {
	AddEvent(ctx context.Context, event Event) error
	EditEvent(ctx context.Context, event Event) error
	GetEvent(ctx context.Context, id string) (Event, error)
	DeleteEvent(ctx context.Context, id string) error
	ListEvents(ctx context.Context) []Event
}
