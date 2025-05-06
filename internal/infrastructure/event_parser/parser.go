package event_parser

import (
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"regexp"
)

type TextEventParser struct {
	reader  *FileReader
	adapter *EventAdapter
}

func NewTextEventParser() *TextEventParser {
	return &TextEventParser{
		reader:  NewFileReader(),
		adapter: NewEventAdapter(),
	}
}

func (p *TextEventParser) ParseEvents(path string) ([]models.Event, error) {
	lines, err := p.reader.ReadLines(path)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	for i, line := range lines {
		event, err := p.parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", i+1, err)
		}
		events = append(events, *event)
	}

	return events, nil
}

func (p *TextEventParser) parseLine(line string) (*models.Event, error) {
	pattern := `^\[(\d{2}:\d{2}:\d{2}\.\d{3})\]\s+(\d+)\s+(\d+)(?:\s+(.+))?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return nil, fmt.Errorf("invalid events format")
	}

	return p.adapter.ParseEvent(
		matches[1], // time
		matches[2], // eventID
		matches[3], // competitorID
		matches[4], // extraParams
	)
}
