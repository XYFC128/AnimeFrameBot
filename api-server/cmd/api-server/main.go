package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Join(filepath.Dir(b), "../..")
)

func run(ctx context.Context, addr string) error {
	runCtx, runCancel := signal.NotifyContext(ctx, os.Interrupt)
	defer runCancel()

	serverHandler := NewServer(filepath.Join(basepath, "images"))
	httpServer := &http.Server{
		Addr:    addr,
		Handler: serverHandler,
	}

	go func() {
		log.Printf("API server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening and serving: %s", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-runCtx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
		defer shutdownCancel()
		log.Println("Shutting down API server")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("error shutting down: %s", err)
		}
	}()
	wg.Wait()
	log.Println("API server closed")
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, ":8763"); err != nil {
		log.Fatalf("error running server: %s", err)
	}
}
