package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	s.db.PingContext(ctx)
	s.db.DB.SetMaxOpenConns(s.MaxConns)
	_, err = s.db.Exec("select count(*) from events")
	if err == nil {
		return nil
	}
	pqErr := err.(*pq.Error)
	if pqErr.Code != "42P01" {
		return err
	}
	s.db.DB.Exec(`CREATE TABLE events(
		id text NOT NULL PRIMARY KEY,
		title text not null,
		time_start TIMESTAMP with time zone,
		time_end TIMESTAMP with time zone,
		description text,
		user_id text,
		notify_before TIMESTAMP with time zone,
		notified BIT
	);`)
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	_, err := s.db.ExecContext(ctx, "insert into "+
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

func (s *Storage) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	var ev []storage.Event
	err := s.db.SelectContext(ctx, &ev, "select * from events where id=$1", id)
	if err != nil {
		return storage.Event{}, err
	}
	if len(ev) < 1 {
		return storage.Event{}, nil
	}
	ev[0].TimeStart = ev[0].TimeStart.Local()
	ev[0].TimeEnd = ev[0].TimeEnd.Local()
	ev[0].NotifyBefore = ev[0].NotifyBefore.Local()
	return ev[0], nil
}

func (s *Storage) GetNotifyEvent(ctx context.Context, notifyDate time.Time) ([]storage.Event, error) {
	var evs []storage.Event
	err := s.db.SelectContext(ctx, &evs, "select * from events where notify_before>$1 "+
		"and notify_before<$2 and notified=false", notifyDate, notifyDate.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}
	if len(evs) < 1 {
		return nil, nil
	}
	for _, v := range evs {
		v.TimeStart = v.TimeStart.Local()
		v.TimeEnd = v.TimeEnd.Local()
		v.NotifyBefore = v.NotifyBefore.Local()
	}
	return evs, nil
}

func (s *Storage) SetNotified(ctx context.Context, event storage.Event) error {
	_, err := s.db.ExecContext(ctx, "update events set notified = true where id=$1", event.ID)
	return err
}

func (s *Storage) EditEvent(ctx context.Context, event storage.Event) error {
	evs := s.ListEvents(ctx)
	for _, v := range evs {
		if v.ID == event.ID {
			continue
		}
		if v.TimeStart.Equal(event.TimeStart) {
			return storage.ErrDateBusy
		}
	}
	res, err := s.db.ExecContext(ctx, "update events set "+
		"title=$2, time_start=$3, time_end=$4, description=$5, user_id=$6, notify_before=$7 where id=$1",
		event.ID, event.Title, event.TimeStart, event.TimeEnd, event.Description, event.UserID, event.NotifyBefore)
	if num, _ := res.RowsAffected(); num == 0 {
		err = fmt.Errorf("%w: %v", storage.ErrNotFound, err)
	}
	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, "delete from events where id=$1", id)
	if num, _ := res.RowsAffected(); num == 0 {
		err = fmt.Errorf("%w: %v", storage.ErrNotFound, err)
	}
	return err
}

func (s *Storage) ListEvents(ctx context.Context) []storage.Event {
	var evs []storage.Event
	err := s.db.SelectContext(ctx, &evs, "select * from events")
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

func (s *Storage) ClearCalendar(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "delete from events *")
	return err
}
