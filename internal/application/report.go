package application

import (
	"context"
	"fmt"
	"github.com/BiathlonRaceProto-Yadro/internal/domain/models"
	"github.com/BiathlonRaceProto-Yadro/pkg/utils"
	"log/slog"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

type ReportService struct {
	config     *models.Config
	fullOutput bool
	logger     *slog.Logger
}

func NewReportService(config *models.Config, fullOutput bool, logger *slog.Logger) ReportGenerator {
	return &ReportService{
		config:     config,
		fullOutput: fullOutput,
		logger:     logger,
	}
}

func (r *ReportService) GenerateReport(competitors []*models.Competitor, _ *models.Config) string {
	if r.fullOutput {
		if r.logger.Enabled(context.Background(), slog.LevelDebug) {
			r.logger.Debug("Generating full report", "competitorsCount", len(competitors))
		}
		return r.generateFullReport(competitors)
	}

	if r.logger.Enabled(context.Background(), slog.LevelDebug) {
		r.logger.Debug("Generating short report", "competitorsCount", len(competitors))
	}
	return r.generateShortReport(competitors)
}

// Короткий отчёт
func (r *ReportService) generateShortReport(competitors []*models.Competitor) string {
	sort.Slice(competitors, func(i, j int) bool {
		iTime, jTime := competitors[i].TotalTime(), competitors[j].TotalTime()
		if iTime == 0 || jTime == 0 {
			return iTime > jTime // push non-finishers to bottom
		}
		return iTime < jTime
	})

	var sb strings.Builder
	sb.WriteString("Final Results:\n")
	for _, c := range competitors {
		status := "[" + r.getStatusString(c) + "]"
		id := c.ID

		var lapsInfo []string
		mainLaps := c.MainLaps()
		for _, lap := range mainLaps {
			if lap.Finish.IsZero() {
				lapsInfo = append(lapsInfo, "{,}")
				continue
			}
			dur := lap.Finish.Sub(lap.Start)
			speed := float64(r.config.LapLen) / dur.Seconds()
			lapsInfo = append(lapsInfo, fmt.Sprintf("{%s, %.3f}", utils.FormatTimestamp(lap.Finish), speed))
		}

		penaltyLaps := c.PenaltyLaps()
		var totalPenaltyTime time.Duration
		var totalPenaltyDistance float64
		for i, lap := range penaltyLaps {
			if lap.Finish.IsZero() {
				continue
			}
			duration := lap.Finish.Sub(lap.Start)
			totalPenaltyTime += duration
			totalPenaltyDistance += float64(c.PenaltyMissedShots()[i] * r.config.PenaltyLen)
		}
		penaltyTimeStr := "-"
		penaltySpeedStr := "-"
		if totalPenaltyTime > 0 {
			penaltyTimeStr = utils.FormatDuration(totalPenaltyTime)
			penaltySpeed := totalPenaltyDistance / totalPenaltyTime.Seconds()
			penaltySpeedStr = fmt.Sprintf("%.3f", penaltySpeed)
		}

		sb.WriteString(fmt.Sprintf(
			"%s %d [%s] {%s, %s} %d/%d\n",
			status,
			id,
			strings.Join(lapsInfo, ", "),
			penaltyTimeStr,
			penaltySpeedStr,
			c.Hits,
			c.Shots,
		))
	}
	return sb.String()
}

// Полный табличный отчёт
func (r *ReportService) generateFullReport(competitors []*models.Competitor) string {
	sort.Slice(competitors, func(i, j int) bool {
		iTime, jTime := competitors[i].TotalTime(), competitors[j].TotalTime()
		if iTime == 0 || jTime == 0 {
			return iTime > jTime // Опускаем в них тех, кто не финишировал
		}
		return iTime < jTime
	})

	var sb strings.Builder
	sb.WriteString("Final Results:\n")
	w := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(w, "ID\tStatus\tTotal Time\tLaps Times\tSpeed Laps\tPenalty Times\tSpeed Penalty\tHits/Shots"); err != nil {
		r.logger.Error("failed to write header", "error", err)
	}
	if _, err := fmt.Fprintln(w, "--\t------\t----------\t----------\t----------\t-------------\t-------------\t----------"); err != nil {
		r.logger.Error("failed to write separator", "error", err)
	}

	for _, c := range competitors {
		timeStr := "-"
		if d := c.TotalTime(); d > 0 {
			timeStr = utils.FormatDuration(d)
		}

		status := r.getStatusString(c)
		mainTimes := r.formatMainLapsDirty(c)
		mainSpeeds := r.formatMainSpeeds(c)
		penTimes := r.formatPenaltyTimes(c)
		penSpeeds := r.formatPenaltySpeeds(c)

		row := fmt.Sprintf(
			"%d\t%s\t%s\t%s\t%s\t%s\t%s\t%d/%d",
			c.ID,
			status,
			timeStr,
			mainTimes,
			mainSpeeds,
			penTimes,
			penSpeeds,
			c.Hits,
			c.Shots,
		)
		if _, err := fmt.Fprintln(w, row); err != nil {
			r.logger.Error("failed to write row", "error", err)
		}
	}

	if err := w.Flush(); err != nil {
		r.logger.Error("failed to flush tabwriter", "error", err)
	}
	return sb.String()
}

func (r *ReportService) getStatusString(c *models.Competitor) string {
	switch c.Status {
	//case c.DisqualificationReason == "NotStarted":
	//	return "NotStarted"
	case models.NotStarted:
		return "NotStarted"
	case models.NotFinished:
		return "NotFinished"
	case models.Finished:
		return "Finished"
	case models.Disqualified:
		return "Disqualified"
	default:
		return "InProgress"
	}
}

// Получения чистых основных кругов
/*
func (r *ReportService) formatMainLapsDirtyClean(c *models.Competitor) string {
	mainLaps := c.MainLaps()
	penaltyLaps := c.PenaltyLaps()

	// Создаем новый массив кругов с скорректированным временем
	adjustedLaps := make([]models.Lap, len(mainLaps))

	for i, mainLap := range mainLaps {
		// Копируем основной круг
		adjustedLap := mainLap

		// Если круг завершен, вычисляем время с учетом штрафов
		if !mainLap.Finish.IsZero() {
			lapDuration := mainLap.Finish.Sub(mainLap.Start)

			// Находим все штрафные круги, которые были выполнены во время этого основного круга
			penaltyTime := time.Duration(0)
			for _, penaltyLap := range penaltyLaps {
				if !penaltyLap.Finish.IsZero() &&
					penaltyLap.Start.After(mainLap.Start) &&
					penaltyLap.Finish.Before(mainLap.Finish) {
					penaltyTime += penaltyLap.Finish.Sub(penaltyLap.Start)
				}
			}

			// Корректируем время круга
			adjustedDuration := lapDuration - penaltyTime
			if adjustedDuration < 0 {
				adjustedDuration = 0
			}
			adjustedLap.Finish = adjustedLap.Start.Add(adjustedDuration)
		}

		adjustedLaps[i] = adjustedLap
	}

	return r.formatLaps(adjustedLaps)
}
*/

// Получение грязных основных кругов
func (r *ReportService) formatMainLapsDirty(c *models.Competitor) string {
	return r.formatLaps(c.MainLaps())
}

func (r *ReportService) formatMainSpeeds(c *models.Competitor) string {
	return r.formatSpeeds(c.MainLaps(), r.config.LapLen)
}

func (r *ReportService) formatPenaltyTimes(c *models.Competitor) string {
	return r.formatLaps(c.PenaltyLaps())
}

func (r *ReportService) formatPenaltySpeeds(c *models.Competitor) string {
	return r.formatSpeedsPenalty(c.PenaltyLaps(), c.PenaltyMissedShots())
}

// Вспомогательные функции.
func (r *ReportService) formatLaps(laps []models.Lap) string {
	var times []string
	for _, lap := range laps {
		if !lap.Finish.IsZero() {
			times = append(times, utils.FormatDuration(lap.Finish.Sub(lap.Start)))
		}
	}
	return strings.Join(times, ", ")
}

func (r *ReportService) formatSpeeds(laps []models.Lap, distance int) string {
	var speeds []string
	for _, lap := range laps {
		if !lap.Finish.IsZero() {
			dur := lap.Finish.Sub(lap.Start).Seconds()
			speed := float64(distance) / dur
			speeds = append(speeds, fmt.Sprintf("%.3f", speed))
		}
	}
	return strings.Join(speeds, ", ")
}

func (r *ReportService) formatSpeedsPenalty(laps []models.Lap, missedShots []int) string {
	var speeds []string
	for i, lap := range laps {
		if missedShots[i] <= 0 {
			continue
		}
		if !lap.Finish.IsZero() {
			dur := lap.Finish.Sub(lap.Start).Seconds()
			distance := missedShots[i] * r.config.PenaltyLen
			speed := float64(distance) / dur
			speeds = append(speeds, fmt.Sprintf("%.3f", speed))
		}
	}
	return strings.Join(speeds, ", ")
}
