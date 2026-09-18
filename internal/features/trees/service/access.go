package service

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
)

func (s *TreesService) CanReadTree(ctx context.Context, userID, treeID string) error {
	_, err := s.treesRepository.GetMemberRole(ctx, userID, treeID)
	return err
}

func (s *TreesService) CanWriteTree(ctx context.Context, userID, treeID string) error {
	role, err := s.treesRepository.GetMemberRole(ctx, userID, treeID)
	if err != nil {
		return err
	}
	if role != domain.TreeRoleOwner {
		return ErrAccessDenied
	}
	return nil
}

var ErrAccessDenied = fmt.Errorf("%w: tree access denied", coreerrors.ErrForbidden)
