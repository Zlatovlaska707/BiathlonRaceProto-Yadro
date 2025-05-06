# Makefile

MAIN_DIR := ./cmd/run/
RUN_NAME := main.go
BINARY_NAME := biathlon
CONFIG_DIR := ./input/config/
CONFIG_NAME:= config.json
EVENTS_DIR := ./input/events/
EVENTS_NAME := events

LINTER = golangci-lint
LINTER_FLAGS = run

.PHONY: all lint build run run-fullOutput clean setup-dirs copy-resources deps

all: run-fullOutput

lint:
	$(LINTER) $(LINTER_FLAGS)

build:
	@echo "Сборка проекта..."
	go build -o $(MAIN_DIR)$(BINARY_NAME) $(MAIN_DIR)$(RUN_NAME)

run: build
	@echo "Запуск приложения в minimal версии..."
	$(MAIN_DIR)$(BINARY_NAME) $(CONFIG_DIR)$(CONFIG_NAME) $(EVENTS_DIR)$(EVENTS_NAME)

run-fullOutput: build
	@echo "Запуск приложения в full версии..."
	$(MAIN_DIR)$(BINARY_NAME) -fullOutput $(CONFIG_DIR)$(CONFIG_NAME) $(EVENTS_DIR)$(EVENTS_NAME)

clean:
	@echo "Очистка артефактов..."
	rm -f $(MAIN_DIR)$(BINARY_NAME)

# Подготовка окружения
setup-dirs:
	mkdir -p $(CONFIG_DIR) $(EVENTS_DIR)

copy-resources: setup-dirs
	@echo "Копирование ресурсов..."
	cp config.json $(CONFIG_DIR)/config.json
	cp events.txt $(EVENTS_DIR)/events

# Управление зависимостями
deps:
	@echo "Инициализация модуля и загрузка зависимостей..."
	go mod init github.com/BiathlonRaceProto-Yadro
	go mod tidy

help:
	@echo "Доступные команды:"
	@echo "  make build     - собрать проект"
	@echo "  make run       - собрать и запустить в minimal версии "
	@echo "  make deps      - установить зависимости"
	@echo "  make setup     - подготовить структуру каталогов и ресурсы"
	@echo "  make clean     - очистить артефакты сборки"
	@echo "  make help      - показать эту справку"

setup: copy-resources deps