package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Chayakorn2002/pms-classroom-backend/config"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/exceptions"
	"github.com/Chayakorn2002/pms-classroom-backend/domain/services"
	sqlite_repository "github.com/Chayakorn2002/pms-classroom-backend/internal/adapters/repositories/sqlite"
	"github.com/Chayakorn2002/pms-classroom-backend/internal/infrastructure/classroom"
	transport_middleware "github.com/Chayakorn2002/pms-classroom-backend/middlewares/transport"
	"github.com/rs/cors"
)

func NewHttpServer() (*http.Server, error) {
	// Initialize context
	ctx := context.Background()

	// Provide config
	cfg := config.ProvideConfig()

	// Initialize sqlite repository
	repo, err := sqlite_repository.NewSqliteRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite repository: %v", err)
	}

	// Setup exception
	applicationError := exceptions.NewApplicationError()

	// Setup Google Classroom service client
	classroomSvc, err := classroom.GetGoogleClassroomClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Classroom service: %v", err)
	}

	// Setup service
	service := services.NewService(
		repo,
		cfg,
		applicationError,
		classroomSvc,
	)

	var middlewares []transport_middleware.TransportMiddleware
	middlewares = append(middlewares, transport_middleware.RequestIdMiddleware())
	middlewares = append(middlewares, transport_middleware.ClaimMiddleware())
	middlewares = append(middlewares, transport_middleware.LoggingMiddleware())

	middlewareStack := transport_middleware.CreateStack(middlewares...)

	handler := registerRoute(service)
	handler = cors.Default().Handler(handler)
	wrappedHandler := middlewareStack(handler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.RestServer.Port),
		Handler: wrappedHandler,
	}

	return server, nil
}
