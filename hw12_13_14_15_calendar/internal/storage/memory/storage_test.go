package memorystorage

import (
	"github.com/hondaworld/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var event1 = storage.Event{
	ID:        "id1",
	Title:     "test title 1",
	StartTime: time.Date(2025, time.April, 1, 9, 0, 0, 0, time.UTC),
}

var event2 = storage.Event{
	ID:        "id2",
	Title:     "test title 2",
	StartTime: time.Date(2025, time.April, 2, 9, 0, 0, 0, time.UTC),
}

func TestStorage(t *testing.T) {
	t.Run("list per date", func(t *testing.T) {
		s := New()

		id1, _ := s.CreateEvent(event1)
		id2, _ := s.CreateEvent(event2)

		event1.ID = id1
		event2.ID = id2

		events, err := s.ListEventsForDay(time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC))

		require.Equal(t, []storage.Event{event1}, events)
		require.Nil(t, err)
	})
	t.Run("lest per week", func(t *testing.T) {
		s := New()

		id1, _ := s.CreateEvent(event1)
		id2, _ := s.CreateEvent(event2)

		event1.ID = id1
		event2.ID = id2

		events, err := s.ListEventsForWeek(time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC))

		require.Equal(t, []storage.Event{event1, event2}, events)
		require.Nil(t, err)
	})
}
