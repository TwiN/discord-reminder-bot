package format

import (
	"strconv"
	"strings"
	"time"
)

func PrettyDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	var parts []string
	if days != 0 {
		if days > 1 {
			parts = append(parts, strconv.Itoa(days)+" days")
		} else {
			parts = append(parts, strconv.Itoa(days)+" day")
		}
	}
	if hours != 0 {
		if hours > 1 {
			parts = append(parts, strconv.Itoa(hours)+" hours")
		} else {
			parts = append(parts, strconv.Itoa(hours)+" hour")
		}
	}
	if minutes != 0 {
		if minutes > 1 {
			parts = append(parts, strconv.Itoa(minutes)+" minutes")
		} else {
			parts = append(parts, strconv.Itoa(minutes)+" minute")
		}
	}
	if seconds != 0 {
		if seconds > 1 {
			parts = append(parts, strconv.Itoa(seconds)+" seconds")
		} else {
			parts = append(parts, strconv.Itoa(seconds)+" second")
		}
	}
	var pretty string // 1h 2m 3s
	for i := range parts {
		if i == 0 {
			pretty = parts[i]
		} else if i+1 == len(parts) {
			pretty = pretty + " and " + parts[i]
		} else {
			pretty = pretty + ", " + parts[i]
		}
	}
	return pretty
}

// ParseDuration is a more lenient version of time.ParseDuration
func ParseDuration(s string) (time.Duration, error) {
	s = strings.ToLower(s)
	if s == "tomorrow" {
		s = "24h"
	}
	var carryOver time.Duration
	duration, err := time.ParseDuration(s)
	if err != nil {
		if strings.Contains(s, "min") {
			s = strings.ReplaceAll(s, "minutes", "m")
			s = strings.ReplaceAll(s, "minute", "m")
			s = strings.ReplaceAll(s, "mins", "m")
			s = strings.ReplaceAll(s, "min", "m")
		}
		if strings.Contains(s, "hour") {
			s = strings.ReplaceAll(s, "hours", "h")
			s = strings.ReplaceAll(s, "hour", "h")
		}
		if strings.Contains(s, "hr") {
			s = strings.ReplaceAll(s, "hrs", "h")
			s = strings.ReplaceAll(s, "hr", "h")
		}
		if strings.Contains(s, "day") {
			s = strings.ReplaceAll(s, "days", "d")
			s = strings.ReplaceAll(s, "day", "d")
		}
		if strings.Contains(s, "d") {
			numberStart, numberEnd := -1, -1
			for i := range s {
				if s[i] >= '0' && s[i] <= '9' {
					if numberStart == -1 {
						numberStart = i
					}
				} else {
					if s[i] == 'd' {
						numberEnd = i
						break
					} else {
						numberStart = -1
					}
				}
			}
			if numberStart > -1 && numberEnd >= numberStart {
				numberOfDays, _ := strconv.Atoi(s[numberStart:numberEnd])
				carryOver = time.Duration(numberOfDays) * 24 * time.Hour
				s = s[:numberStart] + s[numberEnd+1:]
				// If the only duration in s was a number of days, we'll return the carryover right now, otherwise
				// time.ParseDuration(s) will fail because s == ""
				if len(s) == 0 {
					return carryOver, nil
				}
			}
		}
	}
	// Now that s has been sanitized, we can try parsing it again
	duration, err = time.ParseDuration(s)
	if err == nil {
		duration += carryOver
	}
	return duration, err
}
