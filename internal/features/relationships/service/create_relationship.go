package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/validation"
)

func (s *RelationshipsService) CreateRelationship(ctx context.Context, userID, treeID string, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	prepared, err := PrepareCreateRelationship(in)
	if err != nil {
		return domain.Relationship{}, err
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Relationship{}, err
	}
	return s.relationshipsRepository.CreateRelationship(ctx, treeID, prepared)
}

func PrepareCreateRelationship(in domain.CreateRelationshipInput) (domain.CreateRelationshipInput, error) {
	metadata, err := normalizeMetadata(in.Metadata)
	if err != nil {
		return domain.CreateRelationshipInput{}, err
	}
	in.Metadata = metadata
	if err := validation.ValidateULID(in.Person1ID); err != nil {
		return domain.CreateRelationshipInput{}, ErrInvalid
	}
	if err := validation.ValidateULID(in.Person2ID); err != nil || in.Person1ID == in.Person2ID {
		return domain.CreateRelationshipInput{}, ErrInvalid
	}
	if in.Type != domain.RelationshipParentChild && in.Type != domain.RelationshipSpouse {
		return domain.CreateRelationshipInput{}, ErrInvalid
	}
	if in.Type == domain.RelationshipParentChild {
		if in.Direction == nil || (*in.Direction != domain.DirectionParent && *in.Direction != domain.DirectionChild) {
			return domain.CreateRelationshipInput{}, ErrInvalid
		}
		if *in.Direction == domain.DirectionChild {
			in.Person1ID, in.Person2ID = in.Person2ID, in.Person1ID
			d := domain.DirectionParent
			in.Direction = &d
		}
	} else {
		in.Direction = nil
		if in.Person1ID > in.Person2ID {
			in.Person1ID, in.Person2ID = in.Person2ID, in.Person1ID
		}
	}
	return in, nil
}
