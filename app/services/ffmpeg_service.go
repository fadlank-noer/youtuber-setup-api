package services

import (
	"bufio"
	"fmt"

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
	args := []string{
		"-i", body.VideoInput.TmpFileloc,
		"-c:v", "libx265",
	}

	// Compose Dynamic CRF Args
	if utils.Contains([]string{"18", "19", "20", "21", "22"}, body.CRFCode) {
		args = append(args, "-preset", "slow", "-crf", body.CRFCode)
	} else if utils.Contains([]string{"23", "24"}, body.CRFCode) {
		args = append(args, "-preset", "medium", "-crf", body.CRFCode)
	} else {
		args = append(args, "-preset", "fast", "-crf", body.CRFCode)
	}

	// Compose Final Args
	args = append(args,
		"-c:a", "copy",
		"-write_tmcd", "0",
		"-f", "mp4",
		"-movflags", "frag_keyframe+empty_moov",
		"pipe:1",
	)

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ffmpeg.WriteTmcd(w, args); err != nil {
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		}
	}))

	return nil
}
