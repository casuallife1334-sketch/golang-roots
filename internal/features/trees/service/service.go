package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/storage"
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

type TreesService struct {
	treesRepository TreesRepository
	fileStorage     storage.FileStorage
}

func NewTreesService(treesRepository TreesRepository, fileStorage storage.FileStorage) *TreesService {
	return &TreesService{treesRepository: treesRepository, fileStorage: fileStorage}
}
