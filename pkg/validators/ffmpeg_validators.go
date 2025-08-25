package validators

import (
	"fmt"

	"github.com/youtuber-setup-api/pkg/utils"
)

func CRFCodeValidator(crf_code string) error {
	// Define Valid Parameter
	valid_crf := []string{
		"18", "19", "20", "21",
		"22", "23", "24", "25",
		"26", "27", "28", "29",
	}
	if utils.Contains(valid_crf, crf_code) {
		return nil
	}

	return fmt.Errorf("Invalid crf_code!")
}
