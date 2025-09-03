package services

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ffmpeg"
	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/utils"
)

func WriteTmcdCompressService(c *fiber.Ctx, body *types.FfmpegTmcdCompressRequest) error {
	// Set Client Header
	filename := "video.mp4"
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Set("Connection", "Keep-Alive")

	// Compose Args
	args1 := []string{
		"-i", body.VideoInput.TmpFileloc,
		"-c:v", "libx264",
	}

	// Compose Dynamic CRF Args
	if utils.Contains([]string{"18", "19", "20", "21", "22"}, body.CRFCode) {
		args1 = append(args1, "-preset", "slow", "-crf", body.CRFCode)
	} else if utils.Contains([]string{"23", "24"}, body.CRFCode) {
		args1 = append(args1, "-preset", "medium", "-crf", body.CRFCode)
	} else {
		args1 = append(args1, "-preset", "fast", "-crf", body.CRFCode)
	}

	// Compose Final Args
	args1 = append(args1,
		"-pix_fmt", "yuv420p",
		"-c:a", "aac", "-b:a", "128k",
		"-write_tmcd", "0",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov+faststart",
		"pipe:1",
	)

	// Compose Second Pipe
	args2 := []string{
		"-i", "pipe:0",
		"-c", "copy",
		"-movflags", "+faststart",
		"-f", "mp4",
		"pipe:1",
	}

	fmt.Printf("Args1: \n%s", strings.Join(args1, " "))
	fmt.Printf("Args2: \n%s", strings.Join(args2, " "))

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ffmpeg.WriteTmcdChained(w, args1, args2); err != nil {
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		}
	}))

	return nil
}
