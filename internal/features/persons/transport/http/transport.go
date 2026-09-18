package http

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/transport/http/server"
	"io"
	"net/http"
)

type PersonsService interface {
	CreatePerson(context.Context, string, string, domain.CreatePersonInput) (domain.Person, error)
	GetPerson(context.Context, string, string, string) (domain.Person, error)
	GetPersons(context.Context, string, string) ([]domain.Person, error)
	PatchPerson(context.Context, string, string, string, domain.PatchPersonInput) (domain.Person, error)
	DeletePerson(context.Context, string, string, string) error
	UploadPersonPhoto(context.Context, string, string, string, string, io.Reader) (domain.Person, error)
	GetPersonPhoto(context.Context, string, string, string) (io.ReadCloser, error)
	DeletePersonPhoto(context.Context, string, string, string) error
}

type PersonsHTTPHandler struct {
	personsService PersonsService
}

func NewPersonsHTTPHandlers(personsService PersonsService) *PersonsHTTPHandler {
	return &PersonsHTTPHandler{personsService: personsService}
}

func (h *PersonsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{Method: http.MethodPost, Path: "/trees/{tree_id}/persons", Handler: h.CreatePerson},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/persons", Handler: h.GetPersons},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/persons/{id}", Handler: h.GetPerson},
		{Method: http.MethodPatch, Path: "/trees/{tree_id}/persons/{id}", Handler: h.PatchPerson},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/persons/{id}", Handler: h.DeletePerson},
		{Method: http.MethodPost, Path: "/trees/{tree_id}/persons/{id}/photo", Handler: h.UploadPersonPhoto},
		{Method: http.MethodGet, Path: "/trees/{tree_id}/persons/{id}/photo", Handler: h.GetPersonPhoto},
		{Method: http.MethodDelete, Path: "/trees/{tree_id}/persons/{id}/photo", Handler: h.DeletePersonPhoto},
	}
}
