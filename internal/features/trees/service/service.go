package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
)

var ErrInvalid = fmt.Errorf("%w: invalid tree input", coreerrors.ErrInvalidArgument)

type TreesRepository interface {
	CreateTree(context.Context, string, domain.CreateTreeInput) (domain.Tree, error)
	GetTrees(context.Context, string) ([]domain.Tree, error)
	GetTree(context.Context, string, string) (domain.Tree, error)
	GetMemberRole(context.Context, string, string) (domain.TreeRole, error)
	PatchTree(context.Context, string, string, domain.PatchTreeInput) (domain.Tree, error)
	DeleteTree(context.Context, string, string) error
	GetPhotoURLs(context.Context, string) ([]string, error)
}

type FileDeleter interface {
	Delete(context.Context, string) error
}

type TreesService struct {
	treesRepository TreesRepository
	fileStorage     FileDeleter
}

func NewTreesService(treesRepository TreesRepository, fileStorage FileDeleter) *TreesService {
	return &TreesService{treesRepository: treesRepository, fileStorage: fileStorage}
}
