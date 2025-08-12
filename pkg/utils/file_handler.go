package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RequestFileResult struct {
	multipart.FileHeader
	TmpFilename string
	TmpFileloc  string
}

func RequestBodyFileHandler(c *fiber.Ctx, file_attr []string) (map[string]*RequestFileResult, error) {
	// Map File
	files := make(map[string]*RequestFileResult)

	// Get Unix Timestamp for TMP Filename
	unix_timestamp := time.Now().UnixMilli() // millisecond

	// Loop File Attr
	file_length := len(file_attr)
	for _index_file, _fileattr := range file_attr {
		// Define Temporary Name
		var tmp_name = fmt.Sprintf("%d", unix_timestamp)
		if file_length > 1 {
			tmp_name = fmt.Sprintf("%s_%d", tmp_name, (_index_file + 1))
		}

		// Get File Header
		fileHeader, err := c.FormFile(_fileattr)
		if err != nil {
			return nil, fmt.Errorf("No \"%s\" file found in Request Body!", _fileattr)
		}
		extension := path.Ext(fileHeader.Filename)

		// Compose Tmp Filename and Tmp Fileloc
		tmp_name = fmt.Sprintf("%s%s", tmp_name, extension)
		tmp_fileloc := fmt.Sprintf("./tmp_file/%s", tmp_name)

		// Open File for Stream
		src, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("Failed to open tmp \"%s\" file!", _fileattr)
		}
		defer src.Close()

		// Make File Destination
		dst, err := os.Create(tmp_fileloc)
		if err != nil {
			return nil, fmt.Errorf("Failed to make tmp \"%s\" file destination!", _fileattr)
		}
		defer dst.Close()

		// Copy Stdin Streaming
		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("Failed to save \"%s\" file!", _fileattr)
		}

		// Save in file mapping
		files[_fileattr] = &RequestFileResult{
			TmpFilename: tmp_name,
			TmpFileloc:  tmp_fileloc,
			FileHeader:  *fileHeader,
		}
	}

	return files, nil
}
