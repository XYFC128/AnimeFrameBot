package frame

import (
	"encoding/json"
	"net/http"
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

			queryStr := r.PathValue("query")
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
