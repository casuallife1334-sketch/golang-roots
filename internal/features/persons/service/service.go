package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"io"
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

type PhotoStorage interface {
	Put(context.Context, string, io.Reader, string) (string, error)
	Get(context.Context, string) (io.ReadCloser, error)
	Delete(context.Context, string) error
}

type PersonsService struct {
	personsRepository PersonsRepository
	fileStorage       PhotoStorage
	treeAccess        TreeAccess
}

func NewPersonsService(personsRepository PersonsRepository, fileStorage PhotoStorage, treeAccess TreeAccess) *PersonsService {
	return &PersonsService{personsRepository: personsRepository, fileStorage: fileStorage, treeAccess: treeAccess}
}
