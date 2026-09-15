package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
)

func (s *RelationshipsService) GetRelationships(ctx context.Context, personID string) ([]domain.Relationship, error) {
	return s.relationshipsRepository.GetRelationships(ctx, personID)
}
