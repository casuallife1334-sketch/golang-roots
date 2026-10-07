package http

import (
	"context"
	"genealogy-tree/internal/core/transport/http/server"
	"genealogy-tree/internal/features/exchange/gedcom"
	"genealogy-tree/internal/features/exchange/service"
	"io"
	"net/http"
)

type ExchangeService interface {
	Preview(io.Reader) (gedcom.Snapshot, error)
	Import(context.Context, string, string, io.Reader) (service.ImportResult, gedcom.Snapshot, error)
	Export(context.Context, string, string) ([]byte, error)
}

type ExchangeHTTPHandler struct {
	exchangeService ExchangeService
}

func NewExchangeHTTPHandler(exchangeService ExchangeService) *ExchangeHTTPHandler {
	return &ExchangeHTTPHandler{exchangeService: exchangeService}
}

func (h *ExchangeHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/imports/preview", Handler: h.PreviewImport},
		{Method: http.MethodPost, Path: "/trees/{tree_id}/import", Handler: h.ImportTree},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/export", Handler: h.ExportTree},
	}
}
