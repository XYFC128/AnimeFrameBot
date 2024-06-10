package frame

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

func HandleRandom(imageDir string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			frames, err := initFrames(imageDir)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			imageCountStr := r.PathValue("count")
			imageCount, err := strconv.Atoi(imageCountStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			randomFrames, err := getRandomFrames(frames, imageCount)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			bytes, err := json.Marshal(randomFrames)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if _, err := w.Write(bytes); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
}

func HandleFuzzy(imageDir string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			frames, err := initFrames(imageDir)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			queryStrRaw := r.PathValue("query")
			queryStr, err := url.QueryUnescape(queryStrRaw)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			
			imageCountStr := r.PathValue("count")
			imageCount, err := strconv.Atoi(imageCountStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			randomFrames, err := matchSubtitles(frames, queryStr, imageCount)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			bytes, err := json.Marshal(randomFrames)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if _, err := w.Write(bytes); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
}

func HandleDownload(imageDir string) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fileNameRaw := r.PathValue("image")
			fileName, err := url.QueryUnescape(fileNameRaw)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			http.ServeFile(w, r, filepath.Join(imageDir, fileName))
		})
}
