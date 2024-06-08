package upload

import (
	"mime/multipart"
	"net/http"
)

func isImage(file multipart.File) bool {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		return true
	}
	return false
}
