package models

import (
	"time"
)

type Config struct {
	Laps        int           // Количество кругов основной дистанции
	LapLen      int           // Длина каждого основного круга
	PenaltyLen  int           // Длина каждого штрафного круга
	FiringLines int           // Количество стрелковых рубежей на круг
	Start       time.Time     // Планируемое время старта первого участника
	StartDelta  time.Duration // Планируемый интервал между стартами
}

func NewConfig(
	laps int,
	lapLen int,
	penaltyLen int,
	firingLines int,
	start time.Time,
	startDelta time.Duration,
) *Config {
	return &Config{
		Laps:        laps,
		LapLen:      lapLen,
		PenaltyLen:  penaltyLen,
		FiringLines: firingLines,
		Start:       start,
		StartDelta:  startDelta,
	}
}
