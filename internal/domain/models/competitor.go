package models

import (
	"fmt"
	"time"
)

type CompetitorStatus int

const (
	Registered CompetitorStatus = iota + 1
	OnStart
	Racing
	InFiringRange
	InPenalty
	Finished
	Disqualified
	NotStarted
	NotFinished
)

const (
	maxShots = 5
)

type Lap struct {
	Number    int
	Start     time.Time
	Finish    time.Time
	IsPenalty bool
}

type Competitor struct {
	ID                     int
	Status                 CompetitorStatus
	Scheduled              time.Time
	ActualStart            time.Time
	FinishTime             time.Time
	Laps                   []Lap
	Hits, Shots            int
	DisqualificationReason string
	FiringLines            []firingSession
}

type firingSession struct {
	line      int
	entryTime time.Time
	endTime   time.Time
	hits      map[int]bool
	maxShots  int
}

func NewCompetitor(id int) *Competitor {
	return &Competitor{ID: id, Status: Registered}
}

var transitions = map[CompetitorStatus][]CompetitorStatus{
	Registered:    {OnStart, NotStarted},
	OnStart:       {Racing, NotStarted},
	Racing:        {InFiringRange, InPenalty, Finished, NotFinished},
	InFiringRange: {Racing, InPenalty, NotFinished},
	InPenalty:     {Racing, NotFinished},
}

func (c *Competitor) UpdateStatus(next CompetitorStatus) error {
	if next == Disqualified {
		c.Status = next
		return nil
	}
	allowed := transitions[c.Status]
	for _, s := range allowed {
		if s == next {
			c.Status = next
			return nil
		}
	}
	return fmt.Errorf("invalid transition %v -> %v", c.Status, next)
}

func (c *Competitor) SetScheduled(t time.Time) {
	c.Scheduled = t
}

func (c *Competitor) StartNewLap(penal bool, t time.Time) {
	lap := Lap{Number: len(c.Laps) + 1, IsPenalty: penal, Start: t}
	c.Laps = append(c.Laps, lap)
}

func (c *Competitor) FinishCurrentLap(t time.Time) error {
	if len(c.Laps) == 0 {
		return fmt.Errorf("no lap in progress")
	}
	// Ищем последний незавершённый основной круг
	for i := len(c.Laps) - 1; i >= 0; i-- {
		lap := &c.Laps[i]
		if !lap.IsPenalty && lap.Finish.IsZero() {
			lap.Finish = t
			return nil
		}
	}
	return fmt.Errorf("no unfinished main lap found")
}

func (c *Competitor) EndPenalty(t time.Time) {
	for i := len(c.Laps) - 1; i >= 0; i-- {
		if c.Laps[i].IsPenalty && c.Laps[i].Finish.IsZero() {
			c.Laps[i].Finish = t
			return
		}
	}
}

func (c *Competitor) CompletedMain(total int) bool {
	count := 0
	for _, lap := range c.Laps {
		if !lap.IsPenalty {
			count++
		}
	}
	return count >= total
}

func (c *Competitor) Finish(t time.Time) {
	c.FinishTime = t
}

func (c *Competitor) StartFiring(line, spots int, t time.Time) {
	s := firingSession{
		line:      line,
		entryTime: t,
		hits:      make(map[int]bool),
		maxShots:  maxShots,
	}
	c.FiringLines = append(c.FiringLines, s)
	c.Shots += maxShots // Учитываем все выстрелы за гонку
}

func (c *Competitor) RegisterShot(target int) {
	if len(c.FiringLines) == 0 {
		return
	}
	s := &c.FiringLines[len(c.FiringLines)-1]
	if !s.hits[target] {
		s.hits[target] = true
		c.Hits++
	}
}

func (c *Competitor) FinishFiring(t time.Time) int {
	s := &c.FiringLines[len(c.FiringLines)-1]
	s.endTime = t
	missed := maxShots - len(s.hits)
	return missed
}

func (c *Competitor) TotalTime() time.Duration {
	if c.FinishTime.IsZero() {
		return 0
	}
	return c.FinishTime.Sub(c.Scheduled)
}

func (c *Competitor) AverageSpeed(distance int, laps []Lap) float64 {
	total := time.Duration(0)
	for _, lap := range laps {
		total += lap.Finish.Sub(lap.Start)
	}
	if total == 0 {
		return 0
	}
	return float64(distance*len(laps)) / total.Seconds()
}

func (c *Competitor) MainLaps() []Lap {
	var mainLaps []Lap
	for _, lap := range c.Laps {
		if !lap.IsPenalty {
			mainLaps = append(mainLaps, lap)
		}
	}
	return mainLaps
}

func (c *Competitor) PenaltyLaps() []Lap {
	var penaltyLaps []Lap
	for _, lap := range c.Laps {
		if lap.IsPenalty {
			penaltyLaps = append(penaltyLaps, lap)
		}
	}
	return penaltyLaps
}

func (c *Competitor) PenaltyMissedShots() []int {
	var missed []int
	for _, session := range c.FiringLines {
		missedShots := session.maxShots - len(session.hits)
		if missedShots >= 0 {
			missed = append(missed, missedShots)
		}
	}
	return missed
}
