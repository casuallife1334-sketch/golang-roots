package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *TreesService) GetTrees(ctx context.Context, userID string) ([]domain.Tree, error) {
	return s.treesRepository.GetTrees(ctx, userID)
}

func (s *TreesService) GetTree(ctx context.Context, userID, treeID string) (domain.Tree, error) {
	return s.treesRepository.GetTree(ctx, userID, treeID)
}
