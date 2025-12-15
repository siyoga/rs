package ping

import (
	"context"

	"github.com/siyoga/rollstory/pkg/logger"
)

// Response represents the ping response from service layer
type Response struct {
	Message string
}

// Handler contains business logic for ping operations
type Handler struct {
	log *logger.Logger
}

// NewHandler creates a new ping service handler
func NewHandler(log *logger.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

// Handle executes the ping business logic
func (h *Handler) Handle(ctx context.Context) (*Response, error) {
	// Business logic goes here
	// For now, it's simple, but in real app this could involve:
	// - checking database connectivity
	// - verifying external service availability
	// - running health checks

	h.log.Debug(ctx, "Ping request processed")

	return &Response{
		Message: "pong",
	}, nil
}