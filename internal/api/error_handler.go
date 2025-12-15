package api

import (
	"context"
	"net/http"

	"github.com/siyoga/rollstory/internal/generated/api"
	"github.com/siyoga/rollstory/pkg/logger"
)

// ErrorHandler centralizes error handling for RPC layer
type ErrorHandler struct {
	log *logger.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(log *logger.Logger) *ErrorHandler {
	return &ErrorHandler{
		log: log,
	}
}

// Handle converts application errors to appropriate HTTP responses
func (h *ErrorHandler) Handle(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// Log the error
	h.log.Error(ctx, "RPC error", err)

	// For now, return generic 500 error
	// In future, you can add specific error type handling here
	return err
}

// ToInternalServerError creates a 500 response
func (h *ErrorHandler) ToInternalServerError(message string) api.GetPing500JSONResponse {
	return api.GetPing500JSONResponse{
		InternalServerErrorJSONResponse: api.InternalServerErrorJSONResponse{
			Error: message,
		},
	}
}
