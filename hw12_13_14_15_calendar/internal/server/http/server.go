package internalhttp

import (
	"context"
)

type Server struct {
	logger Logger
	app    Application
}

type Logger interface {
	Info(message string)
	Warn(message string)
	Error(message string)
	Debug(message string)
}

type Application interface {
	Init() error
	Shutdown() error
	Run() error
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		logger: logger,
		app:    app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Запуск сервера...")

	// Инициализация приложения
	if err := s.app.Init(); err != nil {
		s.logger.Error("Не удалось инициализировать приложение")
		return err
	}

	// Запуск приложения
	go func() {
		if err := s.app.Run(); err != nil {
			s.logger.Error("Ошибка выполнения приложения: " + err.Error())
		}
	}()

	// Ожидание отмены контекста
	<-ctx.Done()
	s.logger.Info("Контекст сервера завершен. Остановка сервера...")
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Остановка сервера...")

	// Завершение работы приложения
	if err := s.app.Shutdown(); err != nil {
		s.logger.Error("Не удалось завершить работу приложения: " + err.Error())
		return err
	}

	s.logger.Info("Сервер успешно остановлен.")
	return nil
}

// TODO
