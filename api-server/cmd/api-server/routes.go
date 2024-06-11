package main

import (
	"AnimeFrameBot/internal/frame"
	"AnimeFrameBot/internal/upload"
	"net/http"
)

func addRoutes(mux *http.ServeMux, imageDir string) {
	mux.HandleFunc("GET /frame/random/{count}", frame.HandleRandom(imageDir))
	mux.HandleFunc("GET /frame/fuzzy/{query}/{count}", frame.HandleFuzzy(imageDir))
	mux.HandleFunc("GET /frame/exact/{query}/{count}", frame.HandleExact(imageDir))
	mux.HandleFunc("POST /frame", upload.HandleUpload(imageDir))
	mux.HandleFunc("GET /frame/{image}", frame.HandleDownload(imageDir))
}
