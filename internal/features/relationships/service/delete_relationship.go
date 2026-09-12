package service

import "context"

func (s *RelationshipsService) DeleteRelationship(ctx context.Context, id string) error {
	return s.relationshipsRepository.DeleteRelationship(ctx, id)
}
