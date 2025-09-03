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

func WriteTmcdChained(w *bufio.Writer, args1, args2 []string) error {
	ffmpegExec := GetPath()
	zerolog.Logger().Info().Msg("Run exec ffmpeg->WriteTmcdChained()")

	// Define Commands
	_r, _w := io.Pipe()
	cmd1 := exec.Command(ffmpegExec, args1...)
	cmd2 := exec.Command(ffmpegExec, args2...)

	// Pipe Chaining
	cmd1.Stdout = _w
	cmd2.Stdin = _r

	// Write Buffer
	var b2 bytes.Buffer
	cmd2.Stdout = &b2

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	_w.Close()
	cmd2.Wait()
	io.Copy(os.Stdout, &b2)

	// Proses kedua (remux → +faststart)

	// stdin2, err := cmd2.StdinPipe()
	// if err != nil {
	// 	return fmt.Errorf("Error getting stdin pipe cmd2: %w", err)
	// }

	// stdout2, err := cmd2.StdoutPipe()
	// if err != nil {
	// 	return fmt.Errorf("Error getting stdout pipe cmd2: %w", err)
	// }

	// Buffer stderr biar bisa log error lengkap
	// var stderr1, stderr2 bytes.Buffer
	// cmd1.Stderr = &stderr1
	// cmd2.Stderr = &stderr2

	// // Start kedua proses
	// if err := cmd1.Start(); err != nil {
	// 	return fmt.Errorf("Error starting cmd1: %w", err)
	// }
	// if err := cmd2.Start(); err != nil {
	// 	return fmt.Errorf("Error starting cmd2: %w", err)
	// }

	// // Goroutine: pipe output cmd1 → input cmd2
	// go func() {
	// 	defer stdin2.Close()
	// 	io.Copy(stdin2, stdout1)
	// }()

	// Stream output cmd2 → writer ke client
	buf := make([]byte, 5*1024*1024)
	for {
		n, readErr := b2.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("WriteTmcdChained() fail writing io: %w", writeErr)
			}
			w.Flush()
		}
		if readErr != nil {
			if readErr != io.EOF {
				zerolog.Logger().Error().Msgf("WriteTmcdChained() error reading io stdout: %v", readErr)
			}
			break
		}
	}

	// Tunggu kedua proses selesai
	// if err := cmd1.Wait(); err != nil {
	// 	zerolog.Logger().Error().Msgf("FFmpeg cmd1 stderr:\n%s", stderr1.String())
	// 	return fmt.Errorf("WriteTmcdChained() cmd1 error: %w", err)
	// }
	// if err := cmd2.Wait(); err != nil {
	// 	zerolog.Logger().Error().Msgf("FFmpeg cmd2 stderr:\n%s", stderr2.String())
	// 	return fmt.Errorf("WriteTmcdChained() cmd2 error: %w", err)
	// }

	zerolog.Logger().Info().Msg("WriteTmcdChained() complete.")
	return nil
}
