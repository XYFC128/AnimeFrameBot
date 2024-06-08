package upload

import (
	"mime/multipart"
	"path/filepath"
	"strings"
)

func isImage(handler *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}
