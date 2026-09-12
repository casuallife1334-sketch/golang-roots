package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"github.com/oklog/ulid/v2"
)

func (s *RelationshipsService) CreateRelationship(ctx context.Context, in domain.CreateRelationshipInput) (domain.Relationship, error) {
	if _, err := ulid.Parse(in.Person1ID); err != nil {
		return domain.Relationship{}, ErrInvalid
	}
	if _, err := ulid.Parse(in.Person2ID); err != nil || in.Person1ID == in.Person2ID {
		return domain.Relationship{}, ErrInvalid
	}
	if in.Type != domain.RelationshipParentChild && in.Type != domain.RelationshipSpouse {
		return domain.Relationship{}, ErrInvalid
	}
	if in.Type == domain.RelationshipParentChild {
		if in.Direction == nil || (*in.Direction != domain.DirectionParent && *in.Direction != domain.DirectionChild) {
			return domain.Relationship{}, ErrInvalid
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
	return s.relationshipsRepository.CreateRelationship(ctx, in)
}
