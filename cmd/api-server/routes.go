package main

import (
	"net/http"

	"AnimeFrameBot/internal/frame"
	"AnimeFrameBot/internal/upload"
)

func addRoutes(mux *http.ServeMux, imageDir string) {
	mux.HandleFunc("GET /frame/random/{count}", frame.HandleRandom(imageDir))
	mux.HandleFunc("GET /frame/fuzzy/{query}/{count}", frame.HandleFuzzy(imageDir))
	mux.HandleFunc("POST /frame/", upload.HandleUpload(imageDir))
}
