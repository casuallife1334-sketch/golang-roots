package repository

import (
	"context"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/repository/postgres"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oklog/ulid/v2"
)

func TestFamilyTreeIntegrity(t *testing.T) {
	url := os.Getenv("FAMILY_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set FAMILY_TEST_DATABASE_URL to run family integrity integration tests")
	}
	ctx := context.Background()
	db, err := postgres.NewPool(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	user := ulid.Make().String()
	defer func() {
		if _, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, user); err != nil {
			t.Error(err)
		}
	}()
	if _, err := db.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, 'test')`, user, fmt.Sprintf("%s@example.test", user)); err != nil {
		t.Fatal(err)
	}
	trees := [2]string{ulid.Make().String(), ulid.Make().String()}
	people := [2]string{ulid.Make().String(), ulid.Make().String()}
	for i := range trees {
		if _, err := db.Exec(ctx, `INSERT INTO trees (id, owner_id, name) VALUES ($1, $2, 'test')`, trees[i], user); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(ctx, `INSERT INTO persons (id, tree_id, first_name, last_name) VALUES ($1, $2, 'Test', 'Person')`, people[i], trees[i]); err != nil {
			t.Fatal(err)
		}
	}
	family := ulid.Make().String()
	if _, err := db.Exec(ctx, `INSERT INTO families (id, tree_id) VALUES ($1, $2)`, family, trees[0]); err != nil {
		t.Fatal(err)
	}
	otherPerson := ulid.Make().String()
	if _, err := db.Exec(ctx, `INSERT INTO persons (id, tree_id, first_name, last_name) VALUES ($1, $2, 'Other', 'Person')`, otherPerson, trees[1]); err != nil {
		t.Fatal(err)
	}
	relationship := ulid.Make().String()
	if _, err := db.Exec(ctx, `INSERT INTO relationships (id, tree_id, person1_id, person2_id, type) VALUES ($1, $2, $3, $4, 'spouse')`, relationship, trees[1], people[1], otherPerson); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		query string
		args  []any
	}{
		{"member from another tree", `INSERT INTO family_members (tree_id, family_id, person_id, role) VALUES ($1, $2, $3, 'child')`, []any{trees[0], family, people[1]}},
		{"relationship from another tree", `INSERT INTO family_relationships (tree_id, family_id, relationship_id) VALUES ($1, $2, $3)`, []any{trees[0], family, relationship}},
		{"wrong association tree", `INSERT INTO family_members (tree_id, family_id, person_id, role) VALUES ($1, $2, $3, 'child')`, []any{trees[1], family, people[1]}},
		{"relationship with foreign person", `INSERT INTO relationships (id, tree_id, person1_id, person2_id, type) VALUES ($1, $2, $3, $4, 'spouse')`, []any{ulid.Make().String(), trees[0], people[0], people[1]}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := db.Exec(ctx, tc.query, tc.args...)
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) || pgErr.Code != "23503" {
				t.Fatalf("error = %v, want foreign key violation", err)
			}
		})
	}
	if _, err := db.Exec(ctx, `INSERT INTO family_members (tree_id, family_id, person_id, role) VALUES ($1, $2, $3, 'parent')`, trees[0], family, people[0]); err != nil {
		t.Fatalf("same-tree member rejected: %v", err)
	}
}
