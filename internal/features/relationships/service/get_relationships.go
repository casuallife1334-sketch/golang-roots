package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *RelationshipsService) GetRelationships(ctx context.Context, userID, treeID, personID string) ([]domain.Relationship, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return nil, err
	}
	return s.relationshipsRepository.GetRelationships(ctx, treeID, personID)
}
