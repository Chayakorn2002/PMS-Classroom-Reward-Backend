package main

import (
	"context"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chayakorn2002/pms-classroom-backend/config"
	"github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/server"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/logger"
	"github.com/Chayakorn2002/pms-classroom-backend/utils/network"
)

func init() {
	logger.InitLogger()
}

func main() {
	ctx := context.Background()

	server, err := server.NewHttpServer()
	if err != nil {
		panic(err)
	}

	config := config.ProvideConfig()
	port := config.RestServer.Port

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	local, _ := network.GetLocalIP()
	go func() {
		slog.InfoContext(ctx, fmt.Sprintf("Server started on port %d", port))
		slog.InfoContext(ctx, fmt.Sprintf("Local: http://localhost:%d", port))
		slog.InfoContext(ctx, fmt.Sprintf("Network: http://%s:%d", local, port))
		slog.InfoContext(ctx, "Waiting for incoming requests... (Ctrl+C to quit)")

		if err := server.ListenAndServe(); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("failed to serve: %v\n", err))
		}
	}()

	<-ctx.Done()

	stop()

	slog.InfoContext(ctx, "Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to shutdown: %v\n", err))
	}

	slog.InfoContext(ctx, "Server shutdown successfully.")
}
