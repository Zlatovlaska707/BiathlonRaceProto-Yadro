package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"github.com/BiathlonRaceProto-Yadro/pkg/utils"
	"log/slog"
	"strconv"
	"time"
)

type EventProcessor struct {
	competitors map[int]*models.Competitor
	config      *models.Config
	logger      *slog.Logger
}

func NewEventProcessor(cfg *models.Config, lg *slog.Logger) *EventProcessor {
	return &EventProcessor{
		competitors: make(map[int]*models.Competitor),
		config:      cfg,
		logger:      lg,
	}
}

func (p *EventProcessor) HandleEvent(event models.Event) error {
	c := p.getOrCreate(event.CompetitorID)

	if err := p.validateOrder(event, c); err != nil {
		return err
	}

	handlers := map[models.EventType]func(*models.Competitor, models.Event) error{
		models.CompetitorRegistered: p.handlerRegister,
		models.StartTimeSet:         p.handlerSetStartTime,
		models.OnStartLine:          p.handlerOnStartLine,
		models.Started:              p.handlerStartRace,
		models.OnFiringRange:        p.handlerEnterFiring,
		models.TargetHit:            p.handlerHitTarget,
		models.LeftFiringRange:      p.handlerLeaveFiring,
		models.EnteredPenalty:       p.handlerEnterPenalty,
		models.LeftPenalty:          p.handlerLeavePenalty,
		models.LapFinished:          p.handlerFinishLap,
		models.CannotContinue:       p.handlerCannotContinue,
	}

	if handler, ok := handlers[event.Type]; ok {
		return handler(c, event)
	}
	return fmt.Errorf("unknown event type: %d", event.Type)
}

func (p *EventProcessor) GetCompetitors() []*models.Competitor {
	list := make([]*models.Competitor, 0, len(p.competitors))
	for _, c := range p.competitors {
		list = append(list, c)
	}
	return list
}

func (p *EventProcessor) getOrCreate(id int) *models.Competitor {
	if c, exists := p.competitors[id]; exists {
		return c
	}
	c := models.NewCompetitor(id, p.logger)
	p.competitors[id] = c
	return c
}

func (p *EventProcessor) validateOrder(event models.Event, c *models.Competitor) error {
	if !c.ActualStart.IsZero() && event.Time.Before(c.ActualStart) {
		return errors.New("event time precedes actual start")
	}
	return nil
}

func (p *EventProcessor) calculateScheduled(id int) time.Time {
	base := p.config.Start
	offset := time.Duration(id-1) * p.config.StartDelta
	return base.Add(offset)
}

// Handlers:
func (p *EventProcessor) handlerRegister(c *models.Competitor, e models.Event) error {
	c.SetScheduled(p.calculateScheduled(c.ID))

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник зарегистрирован",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	return nil
}

func (p *EventProcessor) handlerSetStartTime(c *models.Competitor, e models.Event) error {
	if len(e.ExtraParams) < 1 {
		err := errors.New("missing start time")
		p.logger.Error("missing start time", "error", err)
		return err
	}

	t, err := models.ParseTime(e.ExtraParams[0])
	if err != nil {
		p.logger.Error("invalid time:", "error", err, "input", e.ExtraParams[0])
		return err
	}
	c.SetScheduled(t)

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Время старта участника установлено жеребьёвкой",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID,
			"startTime", e.ExtraParams[0])
	}
	return nil
}

func (p *EventProcessor) handlerOnStartLine(c *models.Competitor, e models.Event) error {
	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник находится на стартовой линии",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	return c.UpdateStatus(models.OnStart)
}

func (p *EventProcessor) handlerStartRace(c *models.Competitor, e models.Event) error {
	sched := p.calculateScheduled(c.ID)
	if e.Time.After(sched.Add(p.config.StartDelta)) {
		if p.logger.Enabled(context.Background(), slog.LevelInfo) {
			p.logger.Info("Участник дисквалифицирован (опоздание на старт)",
				"time", utils.FormatTimestamp(e.Time),
				"competitorID", c.ID)
		}
		return c.UpdateStatus(models.NotStarted)
	}

	c.ActualStart = e.Time
	if err := c.UpdateStatus(models.Racing); err != nil {
		return err
	}

	c.StartNewLap(false, c.Scheduled)
	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник начал движение",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	return nil
}

func (p *EventProcessor) handlerEnterFiring(c *models.Competitor, e models.Event) error {
	if len(e.ExtraParams) < 1 {
		err := errors.New("missing firing line")
		p.logger.Error("missing firing line:", "error", err,
			"competitorID", c.ID, "eventTime:", e.Time, "paramsCount:", len(e.ExtraParams))
		return err
	}

	line, err := strconv.Atoi(e.ExtraParams[0])
	if err != nil {
		p.logger.Error("strconv.Atoi:", "error", err,
			"competitorID:", c.ID, "eventTime:", e.Time, "rawInput:", e.ExtraParams[0])
		return err
	}

	c.StartFiring(line, p.config.FiringLines, e.Time)

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник находится на стрелковом рубеже",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID,
			"firingLine", line)
	}
	return c.UpdateStatus(models.InFiringRange)
}

func (p *EventProcessor) handlerHitTarget(c *models.Competitor, e models.Event) error {
	if len(e.ExtraParams) < 1 {
		err := errors.New("missing target number")
		p.logger.Error("the target number is missing:", "error", err,
			"competitorID:", c.ID, "eventTime:", e.Time, "paramsCount:", len(e.ExtraParams))
		return err
	}

	n, err := strconv.Atoi(e.ExtraParams[0])
	if err != nil {
		p.logger.Error("strconv.Atoi:", "error", err,
			"competitorID:", c.ID, "eventTime:", e.Time, "rawInput:", e.ExtraParams[0])
		return err
	}

	c.RegisterShot(n)
	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Мишень поражена участником",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID,
			"target", n)
	}
	return nil
}

func (p *EventProcessor) handlerLeaveFiring(c *models.Competitor, e models.Event) error {
	missed := c.FinishFiring(e.Time)
	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник покинул стрелковый рубеж",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	if missed > 0 {
		if err := c.UpdateStatus(models.InPenalty); err != nil {
			return err
		}
	} else {
		if err := c.UpdateStatus(models.Racing); err != nil {
			return err
		}
	}
	return nil
}

func (p *EventProcessor) handlerEnterPenalty(c *models.Competitor, e models.Event) error {
	// Начинаем новый штрафной круг
	c.StartNewLap(true, e.Time)

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник начал штрафные круги",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	return nil
}

func (p *EventProcessor) handlerLeavePenalty(c *models.Competitor, e models.Event) error {
	c.EndPenalty(e.Time)

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник завершил штрафные круги",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}
	return c.UpdateStatus(models.Racing)
}

func (p *EventProcessor) handlerFinishLap(c *models.Competitor, e models.Event) error {
	if err := c.FinishCurrentLap(e.Time); err != nil {
		p.logger.Error("Lap completion error", "error", err,
			"competitorID", c.ID, "lapNumber", len(c.Laps))
		return err
	}

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник завершил основной круг",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID)
	}

	if c.CompletedMain(p.config.Laps) {
		c.SetFinish(e.Time)
		return c.UpdateStatus(models.Finished)
	}

	c.StartNewLap(false, e.Time)
	return nil
}

func (p *EventProcessor) handlerCannotContinue(c *models.Competitor, e models.Event) error {
	reason := ""
	if len(e.ExtraParams) > 0 {
		reason = e.ExtraParams[0]
		c.DisqualificationReason = reason
	}

	if p.logger.Enabled(context.Background(), slog.LevelInfo) {
		p.logger.Info("Участник не может продолжить",
			"time", utils.FormatTimestamp(e.Time),
			"competitorID", c.ID,
			"reason", reason)
	}
	return c.UpdateStatus(models.NotFinished)
}
