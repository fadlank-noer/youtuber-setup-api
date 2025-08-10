package ytdlp

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	lib_zerolog "github.com/youtuber-setup-api/lib/zerolog"
)

type Ytdlp struct {
	Url string
}

type GetResolutionMediaItem struct {
	ID         string
	Size       string
	Resolution string
	Codec      string
}

type GetResolutionMediaFormats struct {
	AudioOnly  []GetResolutionMediaItem
	VideoOnly  []GetResolutionMediaItem
	VideoAudio []GetResolutionMediaItem
}

func execYtdlpCmd(args []string) (string, error) {
	// Run Execution
	lib_zerolog.Logger().Info().Msg("Run exec execYtdlpCmd()")
	cmd := exec.Command("lib/ytdlp/yt-dlp.exe", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		lib_zerolog.Logger().Error().Msg("yt-dlp exec -> " + strings.Join(args, " ") + " Failed")
		fmt.Println(err)
		return "", err
	}

	// Convert output to string
	output := out.String()
	return output, nil
}

func (c *Ytdlp) GetListResolution() (GetResolutionMediaFormats, error) {
	// Define Initial Value
	max_retries := 3
	args := []string{"-F", c.Url}
	metadata := GetResolutionMediaFormats{
		AudioOnly:  []GetResolutionMediaItem{},
		VideoOnly:  []GetResolutionMediaItem{},
		VideoAudio: []GetResolutionMediaItem{},
	}

	// Run exec yt-dlp
	cmd_success := false
	var output string
	for i := 0; i < max_retries; i++ {
		// Define Loop Break
		loop_break := true
		output_string, err := execYtdlpCmd(args)

		// Cmd error
		if err != nil {
			loop_break = false
		}

		// Youtube Has Warning
		if strings.Contains(output_string, "have been skipped") || !strings.Contains(output_string, "mp4a") {
			loop_break = false
		}

		// Break Loop
		if loop_break {
			output = output_string
			cmd_success = true
			break
		}
	}

	// Check Success
	if !cmd_success {
		// kalau mau isi default data dari JSON
		return metadata, errors.New("cmd execution failed!")
	}

	// Split by new line
	lines := strings.Split(output, "\n")

	// Loop string stdout data
	for _, line := range lines {
		// Line Filter
		line_continue := true
		if strings.Contains(line, "https") && !strings.HasPrefix(line, "[") {
			if strings.Contains(line, "mp4") || strings.Contains(line, "m4a") {
				line_continue = false
			}
		}

		// Skip with this prefix
		if line_continue {
			continue
		}

		// Get Columns
		cols := strings.Fields(line)
		print("Cols test: ", strings.Join(cols, ", "), "\n")
		if len(cols) >= 3 {
			// Get video_audio
			if strings.Contains(line, "mp4a") {
				metadata.VideoAudio = append(metadata.VideoAudio, GetResolutionMediaItem{
					ID:         cols[0],
					Size:       cols[6],
					Resolution: cols[2],
					Codec:      cols[10] + "+" + cols[11],
				})
			}

			// Get mp4 video only
			if strings.Contains(line, "mp4") && strings.Contains(line, "video only") {
				metadata.VideoOnly = append(metadata.VideoOnly, GetResolutionMediaItem{
					ID:         cols[0],
					Size:       cols[5],
					Resolution: cols[2],
					Codec:      cols[9],
				})
			}

			// Get m4a audio only
			if strings.Contains(line, "m4a") && strings.Contains(line, "audio only") {
				metadata.AudioOnly = append(metadata.AudioOnly, GetResolutionMediaItem{
					ID:         cols[0],
					Size:       cols[6],
					Resolution: "audio only",
					Codec:      cols[12],
				})
			}
		}
	}

	// Merge Video Only and Audio Only
	for _, video_only := range metadata.VideoOnly {
		// Audio Loop
		for _, audio_only := range metadata.AudioOnly {
			metadata.VideoAudio = append(metadata.VideoAudio, GetResolutionMediaItem{
				ID:         video_only.ID + "+" + audio_only.ID,
				Size:       video_only.Size + "+" + audio_only.Size,
				Resolution: video_only.Resolution,
				Codec:      video_only.Codec + "+" + audio_only.Codec,
			})
		}
	}

	return metadata, nil
}
