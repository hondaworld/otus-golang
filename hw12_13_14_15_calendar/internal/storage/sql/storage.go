package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hondaworld/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage struct {
	db  *sql.DB
	dsn string // DSN - Data Source Name, строка подключения
}

// New создает новый экземпляр Storage
func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

// Connect устанавливает соединение с базой данных
func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open("postgres", s.dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Проверка соединения с базой данных
	if err := s.db.PingContext(ctx); err != nil {
		s.db.Close()
		return fmt.Errorf("error pinging database: %w", err)
	}

	return nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return fmt.Errorf("error closing database: %w", err)
		}
	}
	return nil
}

func (s *Storage) CreateEvent(event storage.Event) (string, error) {
	event.ID = uuid.New().String()
	query := `INSERT INTO events (id, title, start_time, duration, description, user_id, notify_before)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(query, event.ID, event.Title, event.StartTime, event.Duration, event.Description, event.UserID, event.NotifyBefore)
	if err != nil {
		return "", err
	}
	return event.ID, nil
}

func (s *Storage) UpdateEvent(eventID string, event storage.Event) error {
	query := `UPDATE events SET title = $2, start_time = $3, duration = $4, description = $5, user_id = $6, notify_before = $7 WHERE id = $1`
	_, err := s.db.Exec(query, eventID, event.Title, event.StartTime, event.Duration, event.Description, event.UserID, event.NotifyBefore)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(eventID string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.db.Exec(query, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListEventsForDay(date time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, start_time, duration, description, user_id, notify_before 
           FROM events 
           WHERE start_time::date = $1::date`
	return s.fetchEvents(query, date)
}

func (s *Storage) ListEventsForWeek(startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.AddDate(0, 0, 7)
	query := `SELECT id, title, start_time, duration, description, user_id, notify_before 
           FROM events 
           WHERE start_time >= $1 AND start_time < $2`
	return s.fetchEvents(query, startDate, endDate)
}

func (s *Storage) ListEventsForMonth(startDate time.Time) ([]storage.Event, error) {
	endDate := startDate.AddDate(0, 1, 0)
	query := `SELECT id, title, start_time, duration, description, user_id, notify_before 
           FROM events 
           WHERE start_time >= $1 AND start_time < $2`
	return s.fetchEvents(query, startDate, endDate)
}

func (s *Storage) fetchEvents(query string, args ...interface{}) ([]storage.Event, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(&event.ID, &event.Title, &event.StartTime, &event.Duration, &event.Description, &event.UserID, &event.NotifyBefore)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
