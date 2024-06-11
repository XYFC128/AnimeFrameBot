package upload

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func HandleUpload(imageDir string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "Request too large", http.StatusBadRequest)
				return
			}

			file, handler, err := r.FormFile("image")
			if err != nil {
				http.Error(w, "Error retrieving file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			if !isImage(file) {
				http.Error(w, "File is not an image", http.StatusBadRequest)
				return
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				http.Error(w, "Error resetting file cursor", http.StatusInternalServerError)
				return
			}

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Error reading file", http.StatusInternalServerError)
				return
			}

			hash := sha256.Sum256(fileBytes)
			hashString := hex.EncodeToString(hash[:])

			fileName, err := url.QueryUnescape(handler.Filename)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			ext := filepath.Ext(fileName)
			baseName := strings.TrimSuffix(fileName, ext)
			newFileName := baseName + "_" + hashString + ext

			dst, err := os.Create(filepath.Join(imageDir, newFileName))
			if err != nil {
				http.Error(w, "Error creating file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				http.Error(w, "Error resetting file cursor", http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, file); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
		})
}
