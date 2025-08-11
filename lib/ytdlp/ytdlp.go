package ytdlp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
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

func GetYtDlpPath() string {
	switch runtime.GOOS {
	case "windows":
		return "lib/ytdlp/yt-dlp.exe"
	case "linux":
		return "lib/ytdlp/yt-dlp"
	case "darwin":
		return "lib/ytdlp/yt-dlp"
	default:
		fmt.Println("OS tidak didukung:", runtime.GOOS)
		os.Exit(1)
		return ""
	}
}

func ytdlpStringDataExtractor(args []string) (string, error) {
	// Get Yt-Dlp Path by OS
	ytdlpExec := GetYtDlpPath()

	// Run Execution
	lib_zerolog.Logger().Info().Msg("Run exec ytdlpStringDataExtractor()")
	cmd := exec.Command(ytdlpExec, args...)
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
		output_string, err := ytdlpStringDataExtractor(args)

		// Cmd error
		if err != nil {
			loop_break = false
		}

		// Youtube Has Warning
		if strings.Contains(output_string, "WARNING") || !strings.Contains(output_string, "mp4a") {
			loop_break = false
		}

		// Video Only Not Meet Length Criteria
		count := strings.Count(output_string, "https")
		if count < 7 {
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
		lib_zerolog.Logger().Error().Msg("GetListResolution() cmd exec failed!")
		return metadata, errors.New("cmd execution failed")
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

func (c *Ytdlp) GetVideoTitle() (string, error) {
	// Get Video Title
	var video_title string
	args := []string{"-e", c.Url}
	output_string, err := ytdlpStringDataExtractor(args)
	if err != nil {
		lib_zerolog.Logger().Error().Msg("GetVideoTitle() cmd exec failed!")
		return "", errors.New("cmd execution failed")
	}

	// Split by new line
	lines := strings.Split(output_string, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "error:") {
			continue
		}
		video_title = line
	}

	return video_title, nil
}

func (c *Ytdlp) DownloadVideo(w *bufio.Writer, args []string) error {
	// Get Yt-Dlp Path by OS
	ytdlpExec := GetYtDlpPath()
	args = append(args, c.Url)

	// Run Execution
	lib_zerolog.Logger().Info().Msg("Run exec ytdlp->DownloadVideo()")
	cmd := exec.Command(ytdlpExec, args...)

	// Get Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan stdout pipe: %w", err)
	}

	// Get Stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("gagal menjalankan yt-dlp: %w", err)
	}

	// Log output dari stderr secara async
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Println("yt-dlp log:", scanner.Text())
		}
	}()

	// Kirim stdout secara bertahap (chunk)
	buf := make([]byte, 5*1024*1024) // 5MB chunk
	for {
		n, readErr := stdout.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("gagal menulis ke client: %w", writeErr)
			}
			w.Flush() // flush setiap chunk
		}
		if readErr != nil {
			if readErr != io.EOF {
				log.Println("Error membaca stdout:", readErr)
			}
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp selesai dengan error: %w", err)
	}

	log.Println("Stream selesai.")
	return nil
}
