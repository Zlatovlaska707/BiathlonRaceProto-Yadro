package config

import (
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"github.com/BiathlonRaceProto-Yadro/pkg/utils"
)

type ConfigAdapter struct{}

func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

func (a *ConfigAdapter) ToDomain(raw RawConfig) (*models.Config, error) {
	startTime, err := utils.ParseTime(raw.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	delta, err := utils.ParseDuration(raw.StartDelta)
	if err != nil {
		return nil, fmt.Errorf("invalid start delta: %w", err)
	}

	if raw.Laps <= 0 {
		return nil, fmt.Errorf("laps must be positive")
	}

	return models.NewConfig(
		raw.Laps,
		raw.LapLen,
		raw.PenaltyLen,
		raw.FiringLines,
		startTime,
		delta,
	), nil
}
