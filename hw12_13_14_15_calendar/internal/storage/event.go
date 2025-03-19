package storage

import "time"

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	Duration     time.Duration
	Description  string
	UserID       string
	NotifyBefore time.Duration
}

// Notification представляет собой уведомление
type Notification struct {
	EventID   string
	Title     string
	EventTime time.Time
	UserID    string
}

// EventStorage описывает интерфейс для работы с событиями
type EventStorage interface {
	CreateEvent(event Event) (string, error)
	UpdateEvent(eventID string, event Event) error
	DeleteEvent(eventID string) error
	ListEventsForDay(date time.Time) ([]Event, error)
	ListEventsForWeek(startDate time.Time) ([]Event, error)
	ListEventsForMonth(startDate time.Time) ([]Event, error)
}
