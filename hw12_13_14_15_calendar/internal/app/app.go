package app

import (
	"context"
	"github.com/hondaworld/otus-golang/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
	Debug(message string)
}

type Storage interface {
	CreateEvent(event storage.Event) (string, error)
	UpdateEvent(eventID string, event storage.Event) error
	DeleteEvent(eventID string) error
	ListEventsForDay(date time.Time) ([]storage.Event, error)
	ListEventsForWeek(startDate time.Time) ([]storage.Event, error)
	ListEventsForMonth(startDate time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) Init() error {
	a.logger.Info("Инициализация приложения...")

	a.logger.Info("Приложение успешно инициализировано.")
	return nil
}

// Shutdown завершают работу приложения, освобождая ресурсы.
func (a *App) Shutdown() error {
	a.logger.Info("Завершение работы приложения...")

	a.logger.Info("Приложение успешно завершило работу.")
	return nil
}

// Run выполняет основную логику приложения.
func (a *App) Run() error {
	a.logger.Info("Запуск основного выполнения приложения...")

	a.logger.Info("Основное выполнение приложения завершено.")
	return nil
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	a.logger.Info("Создание нового события с ID: " + id)

	event := storage.Event{
		ID:    id,
		Title: title,
	}

	if _, err := a.storage.CreateEvent(event); err != nil {
		a.logger.Error("Ошибка создания события: " + err.Error())
		return err
	}

	a.logger.Info("Событие успешно создано!")
	return nil
}

// TODO
