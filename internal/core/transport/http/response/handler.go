package response

import (
	"encoding/json"
	"errors"
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"net/http"
)

type HTTPResponseHandler struct{ rw http.ResponseWriter }

func NewHTTPResponseHandler(rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{rw: rw}
}

func (h *HTTPResponseHandler) JSONResponse(body any, status int) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(status)
	_ = json.NewEncoder(h.rw).Encode(body)
}

func (h *HTTPResponseHandler) NoContentResponse() { h.rw.WriteHeader(http.StatusNoContent) }

func (h *HTTPResponseHandler) ErrorResponse(err error, message string) {
	status := http.StatusInternalServerError
	errorText := "internal server error"
	switch {
	case errors.Is(err, coreerrors.ErrInvalidArgument):
		status, errorText = http.StatusBadRequest, coreerrors.ErrInvalidArgument.Error()
	case errors.Is(err, coreerrors.ErrConflict):
		status, errorText = http.StatusConflict, coreerrors.ErrConflict.Error()
	case errors.Is(err, coreerrors.ErrNotFound):
		status, errorText = http.StatusNotFound, coreerrors.ErrNotFound.Error()
	}
	h.JSONResponse(ErrorResponse{Error: errorText, Message: message}, status)
}

func (h *HTTPResponseHandler) PanicResponse(value any, message string) {
	h.ErrorResponse(fmt.Errorf("unexpected panic: %v", value), message)
}
