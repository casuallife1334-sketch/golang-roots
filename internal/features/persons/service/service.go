package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/storage"
)

var ErrInvalid = fmt.Errorf("%w: invalid person input", coreerrors.ErrInvalidArgument)

type PersonsRepository interface {
	CreatePerson(context.Context, domain.CreatePersonInput) (domain.Person, error)
	GetPerson(context.Context, string) (domain.Person, error)
	GetPersons(context.Context) ([]domain.Person, error)
	PatchPerson(context.Context, string, domain.PatchPersonInput) (domain.Person, error)
	DeletePerson(context.Context, string) error
	UpdatePersonPhoto(context.Context, string, *string) (domain.Person, error)
}

type PersonsService struct {
	personsRepository PersonsRepository
	fileStorage       storage.FileStorage
}

func NewPersonsService(personsRepository PersonsRepository, fileStorage storage.FileStorage) *PersonsService {
	return &PersonsService{personsRepository: personsRepository, fileStorage: fileStorage}
}
