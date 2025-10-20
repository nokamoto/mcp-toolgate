package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	scanner := bufio.NewScanner(os.Stdin)
	errorChan := make(chan error, 1)

	var logger *slog.Logger
	var opts slog.HandlerOptions
	if os.Getenv("DEBUG") != "" {
		opts.Level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &opts))

	logger.Info("Starting mcp-toolgate...")

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			logger.Debug("stdin content", "line", line)
		}
		if err := scanner.Err(); err != nil {
			errorChan <- err
		}
		close(errorChan)
	}()

	select {
	case <-ctx.Done():
		logger.Info("Shutting down mcp-toolgate...")
	case err := <-errorChan:
		if err != nil {
			logger.Error("Error reading from stdin", "error", err)
			os.Exit(1)
		}
	}
}
