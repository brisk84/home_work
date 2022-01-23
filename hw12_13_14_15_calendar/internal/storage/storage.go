package storage

import "errors"

var (
	ErrDateBusy = errors.New("date is busy")
	ErrNotFound = errors.New("event not found")
	ErrUUIDBusy = errors.New("uuid is busy")
)

type Calendar interface {
	AddEvent(event Event) error
	EditEvent(event Event) error
	GetEvent(id string) (Event, error)
	DeleteEvent(id string) error
	ListEvents() []Event
}
