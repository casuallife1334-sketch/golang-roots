package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
)

func (s *TreesService) PatchTree(ctx context.Context, userID, treeID string, input domain.PatchTreeInput) (domain.Tree, error) {
	if input.Name != nil && strings.TrimSpace(*input.Name) == "" {
		return domain.Tree{}, ErrInvalid
	}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		input.Name = &name
	}
	if err := s.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Tree{}, err
	}
	return s.treesRepository.PatchTree(ctx, userID, treeID, input)
}
