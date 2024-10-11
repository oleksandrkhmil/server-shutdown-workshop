package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	engine.GET("/v1/resource", func(c *gin.Context) {
		slog.Info("Received request")
		time.Sleep(10 * time.Second) // Long processing...
		slog.Info("Processed request")
		c.JSON(http.StatusOK, gin.H{"message": "Request processed"})
	})

	server := &http.Server{
		Addr:    ":4400",
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Unexpected server error", "error", err.Error())
		}
	}()

	slog.Info("Waiting signal")
	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownChannel
	slog.Info("Signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err.Error())
	}
	slog.Info("App closed")
}
