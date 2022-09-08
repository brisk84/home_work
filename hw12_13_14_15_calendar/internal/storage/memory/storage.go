package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{events: make(map[string]storage.Event)}
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.CheckEvent(ctx, event); err != nil {
		return err
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) EditEvent(ctx context.Context, event storage.Event) error {
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

func (s *Storage) CheckEvent(ctx context.Context, event storage.Event) error {
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

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ev, ok := s.events[id]; ok {
		return ev, nil
	}
	return storage.Event{}, storage.ErrNotFound
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := make([]storage.Event, len(s.events))
	i := 0
	for _, v := range s.events {
		ret[i] = v
		i++
	}
	return ret, nil
}

func (s *Storage) GetEventsOnDay(ctx context.Context, day string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := []storage.Event{}
	dateDay, err := time.Parse("2006-01-02", day)
	if err != nil {
		return nil, err
	}
	for _, v := range s.events {
		if v.TimeStart.Truncate(24 * time.Hour).Equal(dateDay) {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (s *Storage) GetEventsOnWeek(ctx context.Context, day string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := []storage.Event{}
	dateDay, err := time.Parse("2006-01-02", day)
	if err != nil {
		return nil, err
	}
	year, week := dateDay.ISOWeek()
	dateStart := dateDay.AddDate(0, 0, -7)
	for i := 0; i < 15; i++ {
		curDay := dateStart.AddDate(0, 0, i)
		nYear, nWeek := curDay.ISOWeek()
		if nYear == year && nWeek == week {
			for _, v := range s.events {
				if v.TimeStart.Truncate(24 * time.Hour).Equal(curDay) {
					ret = append(ret, v)
				}
			}
		}
	}
	return ret, nil
}

func (s *Storage) GetEventsOnMonth(ctx context.Context, day string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ret := []storage.Event{}
	dateDay, err := time.Parse("2006-01-02", day)
	if err != nil {
		return nil, err
	}
	for _, v := range s.events {
		if v.TimeStart.Month() == dateDay.Month() {
			ret = append(ret, v)
		}
	}
	return ret, nil
}
