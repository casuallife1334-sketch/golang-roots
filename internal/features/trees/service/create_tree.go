package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
)

func (s *TreesService) CreateTree(ctx context.Context, ownerID string, input domain.CreateTreeInput) (domain.Tree, error) {
	if strings.TrimSpace(input.Name) == "" {
		return domain.Tree{}, ErrInvalid
	}
	input.Name = strings.TrimSpace(input.Name)
	return s.treesRepository.CreateTree(ctx, ownerID, input)
}
