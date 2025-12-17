package ping

import (
	"context"
	"github.com/siyoga/rollstory/internal/generated/api"
	"github.com/siyoga/rollstory/pkg/logger"
)

// Handler is a thin RPC adapter for ping endpoint
// It transforms HTTP requests/responses and delegates business logic to app layer
type Handler struct {
	pingHandler  PingHandler
	errorHandler ErrorHandler
	log          *logger.Logger
}

// NewHandler creates a new RPC ping handler
func NewHandler(
	pingHandler PingHandler,
	errorHandler ErrorHandler,
	log *logger.Logger,
) *Handler {
	return &Handler{
		pingHandler:  pingHandler,
		errorHandler: errorHandler,
		log:          log,
	}
}

// GetPing handles the GET /ping endpoint
// This is the network layer - it only transforms data and delegates to service layer
func (h *Handler) GetPing(ctx context.Context, request api.GetPingRequestObject) (api.GetPingResponseObject, error) {
	// Delegate to app layer (business logic)
	response, err := h.pingHandler.Handle(ctx)
	if err != nil {
		// Delegate error handling
		if handledErr := h.errorHandler.Handle(ctx, err); handledErr != nil {
			return api.GetPing500JSONResponse{
				InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse{
					Error: "internal server error",
				},
			}, nil
		}
	}

	// Transform app response to API response
	return api.GetPing200JSONResponse{
		Message: response.Message,
	}, nil
}
