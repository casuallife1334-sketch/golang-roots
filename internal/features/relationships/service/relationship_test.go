package service

import (
	"context"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/nullable"
	"strings"
	"testing"
)

func TestNormalizeMetadata(t *testing.T) {
	t.Run("turns nil into an empty object", func(t *testing.T) {
		metadata, err := normalizeMetadata(nil)
		if err != nil {
			t.Fatal(err)
		}
		if metadata == nil || len(metadata) != 0 {
			t.Fatalf("metadata = %#v", metadata)
		}
	})

	t.Run("preserves unknown fields and trims comment", func(t *testing.T) {
		metadata, err := normalizeMetadata(map[string]any{
			"comment": "  Семейная заметка  ",
			"source":  "archive",
		})
		if err != nil {
			t.Fatal(err)
		}
		if metadata["comment"] != "Семейная заметка" || metadata["source"] != "archive" {
			t.Fatalf("metadata = %#v", metadata)
		}
	})

	t.Run("removes an empty comment", func(t *testing.T) {
		metadata, err := normalizeMetadata(map[string]any{"comment": " \n\t ", "other": true})
		if err != nil {
			t.Fatal(err)
		}
		if _, exists := metadata["comment"]; exists || metadata["other"] != true {
			t.Fatalf("metadata = %#v", metadata)
		}
	})

	for name, metadata := range map[string]map[string]any{
		"non-string comment": {"comment": 42},
		"comment too long":   {"comment": strings.Repeat("я", maxRelationshipCommentLength+1)},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := normalizeMetadata(metadata); err != ErrInvalid {
				t.Fatalf("error = %v, want ErrInvalid", err)
			}
		})
	}
}

func TestPatchRelationshipMetadataSemantics(t *testing.T) {
	t.Run("absent metadata is a repository no-op", func(t *testing.T) {
		repository := &relationshipRepositoryStub{}
		access := &treeAccessStub{}
		service := NewRelationshipsService(repository, access)

		if _, err := service.PatchRelationship(context.Background(), "user", "tree", "relationship", domain.PatchRelationshipInput{}); err != nil {
			t.Fatal(err)
		}
		if !repository.patchCalled || repository.patchInput.Metadata.Set {
			t.Fatalf("patch input = %#v", repository.patchInput)
		}
		if !access.writeCalled {
			t.Fatal("write access was not checked")
		}
	})

	t.Run("null metadata is invalid", func(t *testing.T) {
		repository := &relationshipRepositoryStub{}
		access := &treeAccessStub{}
		service := NewRelationshipsService(repository, access)

		_, err := service.PatchRelationship(context.Background(), "user", "tree", "relationship", domain.PatchRelationshipInput{
			Metadata: nullable.Value[map[string]any]{Set: true},
		})
		if err != ErrInvalid {
			t.Fatalf("error = %v, want ErrInvalid", err)
		}
		if repository.patchCalled || access.writeCalled {
			t.Fatal("invalid metadata reached access or repository")
		}
	})

	t.Run("object metadata is normalized and updated", func(t *testing.T) {
		metadata := map[string]any{"comment": "  note  ", "other": float64(1)}
		repository := &relationshipRepositoryStub{}
		service := NewRelationshipsService(repository, &treeAccessStub{})

		if _, err := service.PatchRelationship(context.Background(), "user", "tree", "relationship", domain.PatchRelationshipInput{
			Metadata: nullable.Value[map[string]any]{Set: true, Value: &metadata},
		}); err != nil {
			t.Fatal(err)
		}
		got := *repository.patchInput.Metadata.Value
		if got["comment"] != "note" || got["other"] != float64(1) {
			t.Fatalf("metadata = %#v", got)
		}
	})
}

type relationshipRepositoryStub struct {
	patchCalled bool
	patchInput  domain.PatchRelationshipInput
}

func (r *relationshipRepositoryStub) CreateRelationship(context.Context, string, domain.CreateRelationshipInput) (domain.Relationship, error) {
	return domain.Relationship{}, nil
}

func (r *relationshipRepositoryStub) GetRelationship(context.Context, string, string) (domain.Relationship, error) {
	return domain.Relationship{}, nil
}

func (r *relationshipRepositoryStub) GetRelationships(context.Context, string, string) ([]domain.Relationship, error) {
	return nil, nil
}

func (r *relationshipRepositoryStub) PatchRelationship(_ context.Context, _, _ string, input domain.PatchRelationshipInput) (domain.Relationship, error) {
	r.patchCalled = true
	r.patchInput = input
	return domain.Relationship{}, nil
}

func (r *relationshipRepositoryStub) DeleteRelationship(context.Context, string, string) error {
	return nil
}

type treeAccessStub struct {
	writeCalled bool
}

func (a *treeAccessStub) CanReadTree(context.Context, string, string) error {
	return nil
}

func (a *treeAccessStub) CanWriteTree(context.Context, string, string) error {
	a.writeCalled = true
	return nil
}
