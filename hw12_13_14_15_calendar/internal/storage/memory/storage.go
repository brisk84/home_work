package memorystorage

import (
	"fmt"
	"sync"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.CheckEvent(event); err != nil {
		return err
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) EditEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrNotFound
	}
	for _, v := range s.events {
		if v.ID == event.ID {
			continue
		}
		if v.TimeStart.Equal(event.TimeStart) {
			return storage.ErrDateBusy
		}
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) CheckEvent(event storage.Event) error {
	var err error
	if _, ok := s.events[event.ID]; ok {
		err = storage.ErrUUIDBusy
	}
	for _, v := range s.events {
		if v.TimeStart.Equal(event.TimeStart) {
			err = fmt.Errorf("%w; %v", err, storage.ErrDateBusy)
		}
	}
	return err
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ev, ok := s.events[id]; ok {
		return ev, nil
	}
	return storage.Event{}, storage.ErrNotFound
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) ListEvents() []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := make([]storage.Event, len(s.events))
	i := 0
	for _, v := range s.events {
		ret[i] = v
		i++
	}
	return ret
}
