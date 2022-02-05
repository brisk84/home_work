package storage

import "time"

type Event struct {
	ID           string
	Title        string
	TimeStart    time.Time `db:"time_start"`
	TimeEnd      time.Time `db:"time_end"`
	Description  string
	UserID       string    `db:"user_id"`
	NotifyBefore time.Time `db:"notify_before"`
}
