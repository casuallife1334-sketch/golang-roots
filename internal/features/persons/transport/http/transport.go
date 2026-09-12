package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"io"
	"net/http"
)

type PersonsService interface {
	CreatePerson(context.Context, domain.CreatePersonInput) (domain.Person, error)
	GetPerson(context.Context, string) (domain.Person, error)
	GetPersons(context.Context) ([]domain.Person, error)
	PatchPerson(context.Context, string, domain.PatchPersonInput) (domain.Person, error)
	DeletePerson(context.Context, string) error
	UploadPersonPhoto(context.Context, string, string, io.Reader) (domain.Person, error)
	GetPersonPhoto(context.Context, string) (io.ReadCloser, error)
	DeletePersonPhoto(context.Context, string) error
}

type PersonsHTTPHandler struct {
	personsService PersonsService
}

func NewPersonsHTTPHandlers(personsService PersonsService) *PersonsHTTPHandler {
	return &PersonsHTTPHandler{personsService: personsService}
}

func (h *PersonsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/persons", Handler: h.CreatePerson},
		{Method: http.MethodGet, Path: "/persons", Handler: h.GetPersons},
		{Method: http.MethodGet, Path: "/persons/{id}", Handler: h.GetPerson},
		{Method: http.MethodPatch, Path: "/persons/{id}", Handler: h.PatchPerson},
		{Method: http.MethodDelete, Path: "/persons/{id}", Handler: h.DeletePerson},
		{Method: http.MethodPost, Path: "/persons/{id}/photo", Handler: h.UploadPersonPhoto},
		{Method: http.MethodGet, Path: "/persons/{id}/photo", Handler: h.GetPersonPhoto},
		{Method: http.MethodDelete, Path: "/persons/{id}/photo", Handler: h.DeletePersonPhoto},
	}
}
