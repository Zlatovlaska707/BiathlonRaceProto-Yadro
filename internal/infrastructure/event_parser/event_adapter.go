package event_parser

import (
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"github.com/BiathlonRaceProto-Yadro/pkg/utils"
	"strconv"
	"strings"
)

type EventAdapter struct{}

func NewEventAdapter() *EventAdapter {
	return &EventAdapter{}
}

func (a *EventAdapter) ParseEvent(
	timeStr string,
	eventIDStr string,
	competitorIDStr string,
	extraParams string,
) (*models.Event, error) {
	// Parse time
	eventTime, err := utils.ParseTime(timeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid events time: %w", err)
	}

	// Parse events type
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid events ID: %w", err)
	}
	eventType, err := models.ParseEventType(eventID)
	if err != nil {
		return nil, err
	}

	// Parse competitor ID
	competitorID, err := strconv.Atoi(competitorIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid competitor ID: %w", err)
	}

	// Parse extra parameters
	var params []string
	if extraParams != "" {
		params = strings.Split(extraParams, " ")
	}

	// Event-specific validations
	switch eventType {
	case models.StartTimeSet:
		if len(params) != 1 {
			return nil, fmt.Errorf("events 2 requires exactly 1 parameter")
		}
		if _, err := utils.ParseTime(params[0]); err != nil {
			return nil, fmt.Errorf("invalid start time parameter: %w", err)
		}
	case models.TargetHit:
		if len(params) != 1 {
			return nil, fmt.Errorf("events 6 requires target number")
		}
		if _, err := strconv.Atoi(params[0]); err != nil {
			return nil, fmt.Errorf("invalid target number: %w", err)
		}
	}

	return models.NewEvent(
		eventTime,
		eventType,
		competitorID,
		params,
	), nil
}
