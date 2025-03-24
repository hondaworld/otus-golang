package memorystorage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/hondaworld/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex //nolint:unused
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(event storage.Event) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	event.ID = id
	s.events[id] = event
	return id, nil
}

func (s *Storage) UpdateEvent(eventID string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[eventID]; !exists {
		return errors.New("event not found")
	}

	event.ID = eventID
	s.events[eventID] = event
	return nil
}

func (s *Storage) DeleteEvent(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[eventID]; !exists {
		return errors.New("event not found")
	}

	delete(s.events, eventID)
	return nil
}

func (s *Storage) ListEventsForDay(date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var events []storage.Event
	for _, event := range s.events {
		if isSameDay(event.StartTime, date) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) ListEventsForWeek(startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	endDate := startDate.AddDate(0, 0, 7)
	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(startDate) && event.StartTime.Before(endDate) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) ListEventsForMonth(startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	endDate := startDate.AddDate(0, 1, 0)
	var events []storage.Event
	for _, event := range s.events {
		if event.StartTime.After(startDate) && event.StartTime.Before(endDate) {
			events = append(events, event)
		}
	}
	return events, nil
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
