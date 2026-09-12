package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *RelationshipsService) GetRelationship(ctx context.Context, id string) (domain.Relationship, error) {
	return s.relationshipsRepository.GetRelationship(ctx, id)
}
