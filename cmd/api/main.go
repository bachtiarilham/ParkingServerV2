package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"modulegue/config"
	"modulegue/core/queue"
	"modulegue/database"
	"modulegue/internal/middleware"
	"modulegue/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := database.OpenMySQL(cfg)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("mysql ping warning: %v; write requests will stay buffered in sqlite queue until mysql recovers", err)
	}

	queueOpenConns := cfg.QueueWorkerCount + 1
	if queueOpenConns < 2 {
		queueOpenConns = 2
	}

	queue, err := queue.Open(
		cfg.QueueSQLitePath,
		time.Duration(cfg.QueueLeaseSeconds)*time.Second,
		time.Duration(cfg.QueueRetrySeconds)*time.Second,
		queueOpenConns,
	)
	if err != nil {
		log.Fatalf("open sqlite queue: %v", err)
	}
	defer queue.Close()

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()

	handler, workerPool := server.NewRouter(&cfg, db, queue, log.Default())
	log.Printf("job queue worker count: %d (sqlite max open conns: %d)", cfg.QueueWorkerCount, queueOpenConns)
	workerWG := workerPool.Start(workerCtx)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           middleware.WithRequestLogger(handler, log.Default()),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("park-server listening on :%s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	workerCancel()
	done := make(chan struct{})
	go func() {
		workerWG.Wait()
		close(done)
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	select {
	case <-done:
	case <-ctx.Done():
		log.Printf("worker shutdown timeout: %v", ctx.Err())
	}
}
