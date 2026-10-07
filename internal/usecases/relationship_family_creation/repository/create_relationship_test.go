package repository_test

import (
	"context"
	"fmt"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/core/repository/postgres"
	familyrepo "genealogy-tree/internal/features/families/repository"
	relrepo "genealogy-tree/internal/features/relationships/repository"
	"genealogy-tree/internal/usecases/relationship_family_creation/repository"
	"os"
	"testing"
	"time"

	"github.com/oklog/ulid/v2"
)

func TestConcurrentParentChildRelationshipsShareFamily(t *testing.T) {
	url := os.Getenv("FAMILY_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set FAMILY_TEST_DATABASE_URL to run family creation integration tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	db, err := postgres.NewPool(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for i := range 5 {
		t.Run(fmt.Sprintf("round %d", i+1), func(t *testing.T) {
			userID, treeID := ulid.Make().String(), ulid.Make().String()
			if _, err := db.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, 'test')`, userID, fmt.Sprintf("%s@example.test", userID)); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if _, err := db.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, userID); err != nil {
					t.Error(err)
				}
			})
			if _, err := db.Exec(ctx, `INSERT INTO trees (id, owner_id, name) VALUES ($1, $2, 'test')`, treeID, userID); err != nil {
				t.Fatal(err)
			}
			parents := [2]string{ulid.Make().String(), ulid.Make().String()}
			child := ulid.Make().String()
			for _, personID := range []string{parents[0], parents[1], child} {
				if _, err := db.Exec(ctx, `INSERT INTO persons (id, tree_id, first_name, last_name) VALUES ($1, $2, 'Test', 'Person')`, personID, treeID); err != nil {
					t.Fatal(err)
				}
			}
			families := familyrepo.NewFamiliesRepository(db)
			creation := repository.NewRelationshipFamilyRepository(db, relrepo.NewRelationshipsRepository(db), families)
			direction := domain.DirectionParent
			start := make(chan struct{})
			results := make(chan error, len(parents))
			for _, parentID := range parents {
				go func(parentID string) {
					<-start
					_, err := creation.CreateRelationshipWithFamily(ctx, treeID, domain.CreateRelationshipInput{
						Person1ID: parentID,
						Person2ID: child,
						Type:      domain.RelationshipParentChild,
						Direction: &direction,
						Metadata:  map[string]any{},
					}, nil)
					results <- err
				}(parentID)
			}
			close(start)
			for range parents {
				if err := <-results; err != nil {
					t.Fatal(err)
				}
			}
			items, err := families.GetFamilies(ctx, treeID)
			if err != nil {
				t.Fatal(err)
			}
			if len(items) != 1 || len(items[0].RelationshipIDs) != 2 || len(items[0].Members) != 3 {
				t.Fatalf("families = %+v, want one family with two links and three members", items)
			}
		})
	}
}
