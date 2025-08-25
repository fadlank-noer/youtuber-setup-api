package ffmpeg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/youtuber-setup-api/lib/zerolog"
)

func GetPath() string {
	switch runtime.GOOS {
	case "windows":
		return "lib/ffmpeg/ffmpeg.exe"
	case "linux":
		return "lib/ffmpeg/ffmpeg"
	case "darwin":
		return "lib/ffmpeg/ffmpeg"
	default:
		zerolog.Logger().Error().Msg(fmt.Sprintf("%s is not supported", runtime.GOOS))
		os.Exit(1)
		return ""
	}
}

func WriteTmcd(w *bufio.Writer, args []string) error {
	ffmpegExec := GetPath()

	zerolog.Logger().Info().Msg("Run exec ffmpeg->WriteTmcd()")
	cmd := exec.Command(ffmpegExec, args...)

	// Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting WriteTmcd() stdout pipe: %w", err)
	}

	// Buffer untuk stderr (biar bisa di-log full)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error exec ffmpeg WriteTmcd(): %w", err)
	}

	// Stream stdout ke writer
	buf := make([]byte, 5*1024*1024)
	for {
		n, readErr := stdout.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("WriteTmcd() fail writing io: %w", writeErr)
			}
			w.Flush()
		}
		if readErr != nil {
			if readErr != io.EOF {
				zerolog.Logger().Error().Msgf("WriteTmcd() error reading io stdout: %v", readErr)
			}
			break
		}
	}

	// Tunggu ffmpeg selesai
	if err := cmd.Wait(); err != nil {
		// log error lengkap dari ffmpeg
		zerolog.Logger().Error().Msgf("FFmpeg stderr:\n%s", stderr.String())
		return fmt.Errorf("WriteTmcd() complete with error: %w", err)
	}

	zerolog.Logger().Info().Msg("WriteTmcd() complete.")
	return nil
}
