package config

import (
	"encoding/json"
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"os"
)

type RawConfig struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type JSONConfigLoader struct{}

func NewJSONConfigLoader() *JSONConfigLoader {
	return &JSONConfigLoader{}
}

func (l *JSONConfigLoader) LoadConfig(path string) (*models.Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var raw RawConfig
	if err := json.Unmarshal(file, &raw); err != nil {
		return nil, fmt.Errorf("invalid config format: %w", err)
	}

	return NewConfigAdapter().ToDomain(raw)
}
