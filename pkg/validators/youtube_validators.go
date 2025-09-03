package validators

import (
	"fmt"
	"strings"
	"time"
)

func YoutubeURLValidator(url string) error {
	// Youtube Link are either youtube.com or youtu.be
	default_url := false
	if strings.Contains(url, "youtube.com/") {
		default_url = true
	}
	if !default_url || strings.Contains(url, "youtu.be/") {
		return nil
	}

	// Default URL validator
	if default_url {
		if strings.Contains(url, "/watch") || strings.Contains(url, "/live") || strings.Contains(url, "/clip") || strings.Contains(url, "/shorts") {
			return nil
		}
	}

	return fmt.Errorf("invalid youtube link")
}

func ValidateSectionTimes(start, end string) error {
	// Valid format section: HH:MM:SS.s
	layout := "15:04:05.0"

	startTime, err := time.Parse(layout, start)
	if err != nil {
		return fmt.Errorf("wrong start_time format: HH:MM:SS.s")
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return fmt.Errorf("wrong end_time format: HH:MM:SS.s")
	}

	// Validate end >= start
	if endTime.Before(startTime) {
		return fmt.Errorf("end_time cannot be less than start_time")
	}

	return nil
}
