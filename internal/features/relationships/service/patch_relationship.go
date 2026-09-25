package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"strings"
	"unicode/utf8"
)

const maxRelationshipCommentLength = 5000

func (s *RelationshipsService) PatchRelationship(ctx context.Context, userID, treeID, id string, input domain.PatchRelationshipInput) (domain.Relationship, error) {
	if input.Metadata.Set {
		if input.Metadata.Value == nil {
			return domain.Relationship{}, ErrInvalid
		}
		metadata, err := normalizeMetadata(*input.Metadata.Value)
		if err != nil {
			return domain.Relationship{}, err
		}
		input.Metadata.Value = &metadata
	}
	if err := s.treeAccess.CanWriteTree(ctx, userID, treeID); err != nil {
		return domain.Relationship{}, err
	}
	return s.relationshipsRepository.PatchRelationship(ctx, treeID, id, input)
}

func normalizeMetadata(metadata map[string]any) (map[string]any, error) {
	normalized := make(map[string]any, len(metadata))
	for key, value := range metadata {
		normalized[key] = value
	}

	comment, exists := normalized["comment"]
	if !exists {
		return normalized, nil
	}
	commentString, ok := comment.(string)
	if !ok {
		return nil, ErrInvalid
	}
	commentString = strings.TrimSpace(commentString)
	if utf8.RuneCountInString(commentString) > maxRelationshipCommentLength {
		return nil, ErrInvalid
	}
	if commentString == "" {
		delete(normalized, "comment")
	} else {
		normalized["comment"] = commentString
	}
	return normalized, nil
}
