package ffmpeg

import (
	"bufio"
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
	// Get Exec Path by OS
	ffmpegExec := GetPath()

	// Run Execution
	zerolog.Logger().Info().Msg("Run exec ffmpeg->WriteTmcd()")
	cmd := exec.Command(ffmpegExec, args...)

	// Get Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting WriteTmcd() stdout pipe: %w", err)
	}

	// Get Stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error getting WriteTmcd() stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error exec ffmpeg WriteTmcd(): %w", err)
	}

	// Stderr logging async
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			zerolog.Logger().Error().Msg(fmt.Sprintf("WriteTmcd() ffmpeg log: %s", scanner.Text()))
		}
	}()

	// Kirim stdout secara bertahap (chunk)
	buf := make([]byte, 5*1024*1024) // 5MB chunk
	for {
		n, readErr := stdout.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("WriteTmcd() fail writing io: %w", writeErr)
			}
			w.Flush() // flush setiap chunk
		}
		if readErr != nil {
			if readErr != io.EOF {
				zerolog.Logger().Error().Msg(fmt.Sprintf("WriteTmcd() error reading io stdout: %w", readErr))
			}
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("WriteTmcd() complete with error: %w", err)
	}

	zerolog.Logger().Info().Msg(fmt.Sprintf("WriteTmcd() complete."))
	return nil
}
