package models

import (
	"fmt"
	"strings"
	"time"
)

type EventType int

const (
	CompetitorRegistered EventType = iota + 1 // Участник зарегистрирован
	StartTimeSet                              // Время старта установлено жеребьёвкой
	OnStartLine                               // Участник находится на стартовой линии
	Started                                   // Участник начал движение
	OnFiringRange                             // Участник находится на стрелковом рубеже
	TargetHit                                 // Мишень поражена
	LeftFiringRange                           // Участник покинул стрелковый рубеж
	EnteredPenalty                            // Участник начал штрафные круги
	LeftPenalty                               // Участник завершил штрафные круги
	LapFinished                               // Участник завершил основной круг
	CannotContinue                            // Участник не может продолжить
)

const timeLayout = "15:04:05.000"

type Event struct {
	Time         time.Time
	Type         EventType
	CompetitorID int
	ExtraParams  []string
}

func NewEvent(
	eventTime time.Time,
	eventType EventType,
	competitorID int,
	params []string,
) *Event {
	return &Event{
		Time:         eventTime,
		Type:         eventType,
		CompetitorID: competitorID,
		ExtraParams:  params,
	}
}

func (e Event) Validate() error {
	if e.CompetitorID <= 0 {
		return fmt.Errorf("invalid competitor ID")
	}

	if e.Time.IsZero() {
		return fmt.Errorf("events time is required")
	}

	return nil
}

func ParseEventType(code int) (EventType, error) {
	if code < 1 || code > 11 {
		return 0, fmt.Errorf("invalid events type code")
	}
	return EventType(code), nil
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(timeLayout, strings.Trim(timeStr, "[]"))
}
