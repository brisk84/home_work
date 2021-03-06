package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brisk84/home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Storage struct {
	db       *sqlx.DB
	DBType   string
	ConnStr  string
	MaxConns int
}

func New(dbType string, connStr string, maxConns int) *Storage {
	return &Storage{DBType: dbType, ConnStr: connStr, MaxConns: maxConns}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sqlx.Open(s.DBType, s.ConnStr)
	if err != nil {
		return err
	}
	s.db.DB.SetMaxOpenConns(s.MaxConns)
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) AddEvent(event storage.Event) error {
	_, err := s.db.Exec("insert into "+
		"events(id, title, time_start, time_end, description, user_id, notify_before) "+
		"values ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.TimeStart, event.TimeEnd, event.Description, event.UserID, event.NotifyBefore)
	if err != nil {
		var pqErr *pq.Error
		if ok := errors.As(err, &pqErr); ok && pqErr.Code == pq.ErrorCode("23505") {
			err = fmt.Errorf("%w: %v", storage.ErrUUIDBusy, err)
		}
	}
	return err
}

func (s *Storage) GetEvent(id string) (storage.Event, error) {
	var ev []storage.Event
	err := s.db.Select(&ev, "select * from events where id=$1", id)
	if err != nil {
		return storage.Event{}, err
	}
	ev[0].TimeStart = ev[0].TimeStart.Local()
	ev[0].TimeEnd = ev[0].TimeEnd.Local()
	ev[0].NotifyBefore = ev[0].NotifyBefore.Local()
	return ev[0], nil
}

func (s *Storage) EditEvent(event storage.Event) error {
	evs := s.ListEvents()
	for _, v := range evs {
		if v.ID == event.ID {
			continue
		}
		if v.TimeStart.Equal(event.TimeStart) {
			return storage.ErrDateBusy
		}
	}
	res, err := s.db.Exec("update events set "+
		"title=$2, time_start=$3, time_end=$4, description=$5, user_id=$6, notify_before=$7 where id=$1",
		event.ID, event.Title, event.TimeStart, event.TimeEnd, event.Description, event.UserID, event.NotifyBefore)
	if num, _ := res.RowsAffected(); num == 0 {
		err = fmt.Errorf("%w: %v", storage.ErrNotFound, err)
	}
	return err
}

func (s *Storage) DeleteEvent(id string) error {
	res, err := s.db.Exec("delete from events where id=$1", id)
	if num, _ := res.RowsAffected(); num == 0 {
		err = fmt.Errorf("%w: %v", storage.ErrNotFound, err)
	}
	return err
}

func (s *Storage) ListEvents() []storage.Event {
	var evs []storage.Event
	err := s.db.Select(&evs, "select * from events")
	if err != nil {
		return nil
	}
	for _, ev := range evs {
		ev.TimeStart = ev.TimeStart.Local()
		ev.TimeEnd = ev.TimeEnd.Local()
		ev.NotifyBefore = ev.NotifyBefore.Local()
	}
	return evs
}

func (s *Storage) ClearCalendar() error {
	_, err := s.db.Exec("delete from events *")
	return err
}
