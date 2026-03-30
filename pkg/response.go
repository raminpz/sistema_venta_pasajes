package pkg

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    any    `json:"meta"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, status int, message string, data any, meta any) {
	response := Response{
		Code:    status,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	WriteJSON(w, status, response)
}

func WriteError(w http.ResponseWriter, _ *http.Request, err error) {
	appErr := AsAppError(err)
	response := ErrorResponse{
		Code:    appErr.Status,
		Message: appErr.Message,
		Error:   appErr.Code,
		Details: appErr.Details,
	}

	WriteJSON(w, appErr.Status, response)
}

// Alias para compatibilidad con handlers existentes
func Error(w http.ResponseWriter, err error) {
	WriteError(w, nil, err)
}

func JSON(w http.ResponseWriter, status int, data any) {
	WriteJSON(w, status, data)
}
