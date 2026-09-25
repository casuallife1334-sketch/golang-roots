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
	if !canWriteTree(role) {
		return ErrAccessDenied
	}
	return nil
}

func canWriteTree(role domain.TreeRole) bool {
	return role == domain.TreeRoleOwner || role == domain.TreeRoleEditor
}

var ErrAccessDenied = fmt.Errorf("%w: tree access denied", coreerrors.ErrForbidden)
