package application

import (
	"context"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"log/slog"
)

type ConfigLoader interface {
	LoadConfig(path string) (*models.Config, error)
}

type EventParser interface {
	ParseEvents(path string) ([]models.Event, error)
}

type EventHandler interface {
	HandleEvent(event models.Event) error
	GetCompetitors() []*models.Competitor
}

type ReportGenerator interface {
	GenerateReport(competitors []*models.Competitor, config *models.Config) string
}

type App struct {
	configLoader    ConfigLoader
	eventParser     EventParser
	eventProcessor  EventHandler
	reportGenerator ReportGenerator
	logger          *slog.Logger
}

func NewApp(
	configLoader ConfigLoader,
	eventParser EventParser,
	processor EventHandler,
	report ReportGenerator,
	logger *slog.Logger,
) *App {
	return &App{
		configLoader:    configLoader,
		eventParser:     eventParser,
		eventProcessor:  processor,
		reportGenerator: report,
		logger:          logger,
	}
}

func (a *App) Run(configPath, eventsPath string, fullOutput bool) (string, error) {
	// Загрузка конфигурации
	if a.logger.Enabled(context.Background(), slog.LevelDebug) {
		a.logger.Debug("Loading configuration", "path", configPath)
	}
	config, err := a.configLoader.LoadConfig(configPath)
	if err != nil {
		a.logger.Error("Failed to load config", "path", configPath, "error", err)
		return "", err
	}

	a.reportGenerator = NewReportService(config, fullOutput, a.logger)
	a.eventProcessor = NewEventProcessor(config, a.logger)

	// Парсинг событий
	if a.logger.Enabled(context.Background(), slog.LevelDebug) {
		a.logger.Debug("Parsing events", "path", eventsPath)
	}
	events, err := a.eventParser.ParseEvents(eventsPath)
	if err != nil {
		a.logger.Error("Failed to read events", "path", eventsPath, "error", err)
		return "", err
	}

	// Обработка событий
	if a.logger.Enabled(context.Background(), slog.LevelDebug) {
		a.logger.Debug("Processing events", "count", len(events))
	}
	for _, event := range events {
		if err := a.eventProcessor.HandleEvent(event); err != nil {
			a.logger.Error("Event processing failed",
				"eventTime", event.Time,
				"competitorID", event.CompetitorID,
				"error", err,
			)
			a.logger.Error("Failed to event processing", "error", err)
			return "", err
		}
	}

	// Генерация отчёта
	a.logger.Info("Generating final report")
	competitors := a.eventProcessor.GetCompetitors()
	return a.reportGenerator.GenerateReport(competitors, config), nil
}
