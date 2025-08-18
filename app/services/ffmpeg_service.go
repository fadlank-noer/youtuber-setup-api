package services

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ffmpeg"
	"github.com/youtuber-setup-api/lib/zerolog"
)

func WriteTmcdService(c *fiber.Ctx, body *types.FfmpegWriteTmcdRequest) error {
	// Set Client Header
	filename := "video.mp4"
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Set("Connecton", "Keep-Alive")

	// Compose Args
	args := []string{
		"-i", body.VideoInput.TmpFileloc,
		"-c:a", "copy",
		"-c:v", "copy",
		"-write_tmcd", "0",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"pipe:1",
	}

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ffmpeg.WriteTmcd(w, args); err != nil {
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		}
	}))

	return nil
}
