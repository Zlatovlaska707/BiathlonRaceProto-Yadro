package utils

import (
	"fmt"
	"strings"
	"time"
)

const timeLayout = "15:04:05.000"

// Преобразует строку формата HH:MM:SS.sss в time.Time
func ParseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(timeLayout, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %w", err)
	}
	return t, nil
}

// Преобразует строку формата HH:MM:SS в time.Duration
func ParseDuration(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration format")
	}

	var h, m, s time.Duration
	_, err := fmt.Sscanf(parts[0], "%d", &h)
	_, err = fmt.Sscanf(parts[1], "%d", &m)
	_, err = fmt.Sscanf(parts[2], "%d", &s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return h*time.Hour + m*time.Minute + s*time.Second, nil
}

// Форматирует duration в HH:MM:SS.sss
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000

	return fmt.Sprintf(
		"%02d:%02d:%02d.%03d",
		hours,
		minutes,
		seconds,
		milliseconds,
	)
}

func FormatTimestamp(t time.Time) string {
	return t.Format(timeLayout)
}
