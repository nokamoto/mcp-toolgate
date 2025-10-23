package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nokamoto/mcp-toolgate/internal/jsonrpc"
)

type replacer interface {
	Replace(input string) (string, error)
}

const (
	debug            = "DEBUG"
	allowedToolNames = "ALLOWED_TOOL_NAMES"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	errorChan := make(chan error, 1)

	var logger *slog.Logger
	var opts slog.HandlerOptions
	if os.Getenv(debug) != "" {
		opts.Level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &opts))

	logger.Info("Starting mcp-toolgate...")

	var replacer replacer
	replacer = jsonrpc.NewAllowedToolGate(strings.Split(os.Getenv(allowedToolNames), ","))

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				errorChan <- fmt.Errorf("failed to read line: %w", err)
				return
			}
			line = strings.TrimSuffix(line, "\n")
			replaced, err := replacer.Replace(line)
			if err != nil {
				errorChan <- fmt.Errorf("failed to replace line: %w", err)
				return
			}
			fmt.Println(replaced)
			logger.Debug("stdin content", "line", line)
		}
		close(errorChan)
	}()

	select {
	case <-ctx.Done():
		logger.Info("Shutting down mcp-toolgate...")
	case err := <-errorChan:
		if err != nil {
			logger.Error("Error occurred", "error", err)
			os.Exit(1)
		}
	}
}
