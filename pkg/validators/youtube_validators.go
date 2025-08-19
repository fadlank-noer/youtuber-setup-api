package validators

import (
	"fmt"
	"strings"
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

	return fmt.Errorf("Invalid youtube link!")
}
