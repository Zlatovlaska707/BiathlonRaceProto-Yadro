package main

import (
	"flag"
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/application"
	"github.com/BiathlonRaceProto-Yadro/internal/infrastructure/config"
	"github.com/BiathlonRaceProto-Yadro/internal/infrastructure/event_parser"
	"github.com/BiathlonRaceProto-Yadro/internal/logging"
	"log/slog"
	"os"
)

func main() {
	logDebug := flag.Bool("debug", false, "Enable debug logs")
	logInfo := flag.Bool("info", false, "Enable info logs")
	logError := flag.Bool("error", false, "Enable error logs")
	fullOutput := flag.Bool("fullOutput", false, "Generate full report")
	flag.Parse()

	logger := logging.СonfigureLogger(*logDebug, *logInfo, *logError)

	args := flag.Args()
	if len(args) != 2 {
		logger.Error("Usage: main.go [flags] <config_path> <events_path>", "argsCount", len(args))
		os.Exit(1)
	}
	configPath := args[0]
	eventsPath := args[1]

	app := initializeApp(logger)

	report, err := app.Run(configPath, eventsPath, *fullOutput)
	if err != nil {
		logger.Error("Application failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Application completed successfully")
	fmt.Println(report)
}

func initializeApp(logger *slog.Logger) *application.App {
	configLoader := config.NewJSONConfigLoader()
	eventParser := event_parser.NewTextEventParser()

	// Создаём временные заглушки, которые будут перезаписаны в Run()
	reportService := application.NewReportService(nil, false, logger)
	processor := application.NewEventProcessor(nil, logger)

	return application.NewApp(
		configLoader,
		eventParser,
		processor,
		reportService,
		logger,
	)
}
