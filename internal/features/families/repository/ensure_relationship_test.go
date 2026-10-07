package repository

import (
	"context"
	"errors"
	"fmt"
	"genealogy-tree/internal/core/domain"
	coreerrors "genealogy-tree/internal/core/errors"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

func TestFindOrCreateFamily(t *testing.T) {
	url := os.Getenv("FAMILY_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set FAMILY_TEST_DATABASE_URL to run family selection integration tests")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	type member struct {
		person int
		role   domain.FamilyMemberRole
	}
	cases := []struct {
		name       string
		kind       domain.RelationshipType
		families   [][]member
		wantFamily int // -1: create a new family
		ambiguous  bool
		requested  *int
	}{
		{"spouses already adults", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberPartner}, {1, domain.FamilyMemberParent}}}, 0, false, nil},
		{"spouses adults in multiple families", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberPartner}, {1, domain.FamilyMemberParent}}, {{0, domain.FamilyMemberParent}, {1, domain.FamilyMemberPartner}}}, 0, true, nil},
		{"siblings do not share spouse family", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberChild}, {1, domain.FamilyMemberChild}}}, -1, false, nil},
		{"spouse joins single-parent family", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberParent}, {2, domain.FamilyMemberChild}}}, 0, false, nil},
		{"spouse does not join other partner", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberPartner}, {2, domain.FamilyMemberPartner}}}, -1, false, nil},
		{"spouse with two available families", domain.RelationshipSpouse, [][]member{{{0, domain.FamilyMemberParent}}, {{1, domain.FamilyMemberParent}}}, 0, true, nil},
		{"parent joins spouse family", domain.RelationshipParentChild, [][]member{{{0, domain.FamilyMemberPartner}, {2, domain.FamilyMemberPartner}}}, 0, false, nil},
		{"parent has two roles in one family", domain.RelationshipParentChild, [][]member{{{0, domain.FamilyMemberPartner}, {0, domain.FamilyMemberParent}}}, 0, false, nil},
		{"parent joins child's existing family", domain.RelationshipParentChild, [][]member{{{1, domain.FamilyMemberChild}, {2, domain.FamilyMemberParent}}}, 0, false, nil},
		{"child identifies one of parent's families", domain.RelationshipParentChild, [][]member{{{0, domain.FamilyMemberParent}, {2, domain.FamilyMemberPartner}}, {{0, domain.FamilyMemberParent}, {1, domain.FamilyMemberChild}}}, 1, false, nil},
		{"parent has multiple families", domain.RelationshipParentChild, [][]member{{{0, domain.FamilyMemberParent}}, {{0, domain.FamilyMemberPartner}}}, 0, true, nil},
		{"child has multiple families", domain.RelationshipParentChild, [][]member{{{1, domain.FamilyMemberChild}}, {{1, domain.FamilyMemberChild}}}, 0, true, nil},
		{"explicit family resolves ambiguity", domain.RelationshipParentChild, [][]member{{{0, domain.FamilyMemberParent}}, {{0, domain.FamilyMemberPartner}}}, 1, false, ptr(1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := pool.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer tx.Rollback(ctx)
			userID, treeID := ulid.Make().String(), ulid.Make().String()
			if _, err := tx.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, 'test')`, userID, fmt.Sprintf("%s@example.test", userID)); err != nil {
				t.Fatal(err)
			}
			if _, err := tx.Exec(ctx, `INSERT INTO trees (id, owner_id, name) VALUES ($1, $2, 'test')`, treeID, userID); err != nil {
				t.Fatal(err)
			}
			people := []string{ulid.Make().String(), ulid.Make().String(), ulid.Make().String()}
			for _, personID := range people {
				if _, err := tx.Exec(ctx, `INSERT INTO persons (id, tree_id, first_name, last_name) VALUES ($1, $2, 'Test', 'Person')`, personID, treeID); err != nil {
					t.Fatal(err)
				}
			}
			var familyIDs []string
			for _, members := range tc.families {
				id := ulid.Make().String()
				familyIDs = append(familyIDs, id)
				if _, err := tx.Exec(ctx, `INSERT INTO families (id, tree_id) VALUES ($1, $2)`, id, treeID); err != nil {
					t.Fatal(err)
				}
				for _, m := range members {
					if _, err := tx.Exec(ctx, `INSERT INTO family_members (tree_id, family_id, person_id, role) VALUES ($1, $2, $3, $4)`, treeID, id, people[m.person], m.role); err != nil {
						t.Fatal(err)
					}
				}
			}
			var requested *string
			if tc.requested != nil {
				requested = &familyIDs[*tc.requested]
			}
			got, err := (&FamiliesRepository{}).findOrCreateFamily(ctx, tx, treeID, domain.Relationship{
				Person1ID: people[0], Person2ID: people[1], Type: tc.kind,
			}, requested)
			if tc.ambiguous {
				if !errors.Is(err, coreerrors.ErrAmbiguousFamily) {
					t.Fatalf("error = %v, want ambiguous family", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tc.wantFamily >= 0 && got != familyIDs[tc.wantFamily] {
				t.Fatalf("family = %s, want %s", got, familyIDs[tc.wantFamily])
			}
			if tc.wantFamily < 0 {
				for _, id := range familyIDs {
					if got == id {
						t.Fatalf("reused unrelated family %s", got)
					}
				}
			}
		})
	}
}

func ptr(n int) *int { return &n }
