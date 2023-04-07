package util

import (
	"fmt"
	"strings"
	"time"
)

func FormatDuration(d time.Duration) string {
	nanoseconds := int(d.Nanoseconds()) % 1000
	microseconds := int(d.Microseconds()) % 1000
	seconds := int(d.Seconds()) % 60
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())

	format := ""
	if hours > 0 {
		format = fmt.Sprintf("%dh %dm %ds %dms %dns", hours, minutes, seconds, microseconds, nanoseconds)
	} else if minutes > 0 {
		format = fmt.Sprintf("%dm %ds %dms %dns", minutes, seconds, microseconds, nanoseconds)
	} else if seconds > 0 {
		format = fmt.Sprintf("%ds %dms %dns", seconds, microseconds, nanoseconds)
	} else if microseconds > 0 {
		format = fmt.Sprintf("%dms %dns", microseconds, nanoseconds)
	} else {
		format = fmt.Sprintf("%dns", nanoseconds)
	}

	return format
}

func FormatPercentageBar(percentage float64, width int) string {
	filledWidth := int(percentage * float64(width))
	filled := strings.Repeat("=", filledWidth)

	unfilledWidth := width - filledWidth
	unfilled := strings.Repeat(" ", unfilledWidth)

	return fmt.Sprintf("[%s%s] %.1f%%", filled, unfilled, percentage*100)
}
