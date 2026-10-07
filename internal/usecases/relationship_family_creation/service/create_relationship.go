package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/validation"
	relationshipservice "genealogy-tree/internal/features/relationships/service"
)

func (s *RelationshipFamilyCreationService) CreateRelationship(ctx context.Context, userID, treeID string, command domain.CreateRelationshipCommand) (domain.Relationship, error) {
	input, err := relationshipservice.PrepareCreateRelationship(command.Relationship)
	if err != nil {
		return domain.Relationship{}, err
	}
	if command.FamilyID != nil {
		if err := validation.ValidateULID(*command.FamilyID); err != nil {
			return domain.Relationship{}, relationshipservice.ErrInvalid
		}
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Relationship{}, err
	}
	return s.repository.CreateRelationshipWithFamily(ctx, treeID, input, command.FamilyID)
}
