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
	CreatePerson(context.Context, string, domain.CreatePersonInput) (domain.Person, error)
	GetPerson(context.Context, string, string) (domain.Person, error)
	GetPersons(context.Context, string) ([]domain.Person, error)
	PatchPerson(context.Context, string, string, domain.PatchPersonInput) (domain.Person, error)
	DeletePerson(context.Context, string, string) error
	UpdatePersonPhoto(context.Context, string, string, *string) (domain.Person, error)
}

type TreeAccess interface {
	CanReadTree(context.Context, string, string) error
	CanWriteTree(context.Context, string, string) error
}

type PersonsService struct {
	personsRepository PersonsRepository
	fileStorage       storage.FileStorage
	treeAccess        TreeAccess
}

func NewPersonsService(personsRepository PersonsRepository, fileStorage storage.FileStorage, treeAccess TreeAccess) *PersonsService {
	return &PersonsService{personsRepository: personsRepository, fileStorage: fileStorage, treeAccess: treeAccess}
}
