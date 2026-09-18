package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *RelationshipsService) GetRelationship(ctx context.Context, userID, treeID, id string) (domain.Relationship, error) {
	if err := s.treeAccess.CanReadTree(ctx, userID, treeID); err != nil {
		return domain.Relationship{}, err
	}
	return s.relationshipsRepository.GetRelationship(ctx, treeID, id)
}
