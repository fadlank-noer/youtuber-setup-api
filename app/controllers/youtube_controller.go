package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/services"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/pkg/utils"
)

func GetVideoResolutionList(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeGetResolutionRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request Body!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return services.GetVideoResolutionListService(c, body)
}

func DownloadVideo(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeDownloadRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request Body!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return services.DownloadVideoService(c, body)
}

// func TestGetDownload(c *fiber.Ctx) error {
// 	// Define Url and Filename
// 	url := "https://youtu.be/yIsPsyYJxRk?si=Su-F9kJE5wZgMCAb"
// 	filename := "video.mp4"

// 	// Set header untuk client
// 	// "inline" akan mencoba memutar video di browser, "attachment" akan langsung men-download
// 	c.Set("Content-Type", "video/mp4")
// 	c.Set("Content-Disposition", `inline; filename="`+filename+`"`)

// 	// SetBodyStreamWriter menjalankan semua proses di dalam satu goroutine
// 	// yang dikelola oleh Fiber, memastikan koneksi tidak terputus.
// 	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
// 		// 1. Argumen yt-dlp diperbaiki
// 		// -f "best[ext=mp4]/best": Memilih format MP4 terbaik yang sudah digabung.
// 		//      Ini lebih simpel daripada "bv+ba" yang mungkin memerlukan ffmpeg untuk merging.
// 		// -o -: Mengarahkan output video ke stdout (standard output), bukan ke file.
// 		args := []string{
// 			"-f", "best[ext=mp4]/best",
// 			url,
// 			"-o", "-",
// 		}
// 		cmd := exec.Command("lib/ytdlp/yt-dlp.exe", args...)

// 		// Ambil stdout dari command untuk di-stream
// 		stdout, err := cmd.StdoutPipe()
// 		if err != nil {
// 			log.Println("Error: Gagal mendapatkan stdout pipe:", err)
// 			return // Hentikan proses jika gagal
// 		}

// 		// Ambil stderr untuk melihat log/error dari yt-dlp saat berjalan
// 		stderr, err := cmd.StderrPipe()
// 		if err != nil {
// 			log.Println("Error: Gagal mendapatkan stderr pipe:", err)
// 			return
// 		}

// 		// 2. Jalankan command di dalam stream writer
// 		if err := cmd.Start(); err != nil {
// 			log.Println("Error: Gagal menjalankan command yt-dlp:", err)
// 			return
// 		}

// 		// Baca stderr di goroutine terpisah agar tidak memblokir
// 		go func() {
// 			scanner := bufio.NewScanner(stderr)
// 			for scanner.Scan() {
// 				log.Println("yt-dlp log:", scanner.Text())
// 			}
// 		}()

// 		// 3. Salin output dari stdout ke response writer
// 		// io.Copy akan menyalin data secara efisien sampai stream selesai.
// 		_, err = io.Copy(w, stdout)
// 		if err != nil {
// 			log.Println("Error: Gagal menyalin stream ke client:", err)
// 		}

// 		// Flush buffer untuk memastikan semua data terkirim
// 		if err := w.Flush(); err != nil {
// 			log.Println("Error: Gagal flush writer:", err)
// 		}

// 		// 4. Tunggu command selesai dan bersihkan resource
// 		// Ini penting untuk menangkap error dari proses yt-dlp itu sendiri.
// 		if err := cmd.Wait(); err != nil {
// 			log.Println("Error: Proses yt-dlp selesai dengan error:", err)
// 		}

// 		log.Println("Stream selesai.")
// 	}))

// 	return nil
// }
