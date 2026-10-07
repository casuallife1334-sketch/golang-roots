package repository_test

import (
	"context"
	"errors"
	"fmt"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/repository/postgres"
	familyrepo "genealogy-tree/internal/features/families/repository"
	personrepo "genealogy-tree/internal/features/persons/repository"
	relrepo "genealogy-tree/internal/features/relationships/repository"
	"genealogy-tree/internal/usecases/family_lifecycle/repository"
	"os"
	"testing"

	"github.com/oklog/ulid/v2"
)

type fixture struct {
	db         *postgres.Pool
	lifecycle  *repository.Repository
	tree       string
	family     string
	parents    [2]string
	child      string
	spouse     string
	childLinks [2]string
}

func setupLifecycle(t *testing.T) fixture {
	t.Helper()
	url := os.Getenv("FAMILY_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set FAMILY_TEST_DATABASE_URL to run family lifecycle integration tests")
	}
	ctx := context.Background()
	db, err := postgres.NewPool(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	user := ulid.Make().String()
	t.Cleanup(func() {
		if _, err := db.Exec(ctx, `DELETE FROM users WHERE id = $1`, user); err != nil {
			t.Error(err)
		}
		db.Close()
	})
	f := fixture{
		db: db, tree: ulid.Make().String(), family: ulid.Make().String(),
		parents: [2]string{ulid.Make().String(), ulid.Make().String()},
		child:   ulid.Make().String(), spouse: ulid.Make().String(),
		childLinks: [2]string{ulid.Make().String(), ulid.Make().String()},
	}
	if _, err := db.Exec(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, 'test')`, user, fmt.Sprintf("%s@example.test", user)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(ctx, `INSERT INTO trees (id, owner_id, name) VALUES ($1, $2, 'test')`, f.tree, user); err != nil {
		t.Fatal(err)
	}
	for _, person := range []string{f.parents[0], f.parents[1], f.child} {
		if _, err := db.Exec(ctx, `INSERT INTO persons (id, tree_id, first_name, last_name) VALUES ($1, $2, 'Test', 'Person')`, person, f.tree); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(ctx, `INSERT INTO relationships (id, tree_id, person1_id, person2_id, type) VALUES ($1, $2, $3, $4, 'spouse')`, f.spouse, f.tree, f.parents[0], f.parents[1]); err != nil {
		t.Fatal(err)
	}
	for i, id := range f.childLinks {
		if _, err := db.Exec(ctx, `INSERT INTO relationships (id, tree_id, person1_id, person2_id, type, direction) VALUES ($1, $2, $3, $4, 'parent_child', 'parent')`, id, f.tree, f.parents[i], f.child); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(ctx, `INSERT INTO families (id, tree_id, auto_created) VALUES ($1, $2, true)`, f.family, f.tree); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{f.spouse, f.childLinks[0], f.childLinks[1]} {
		if _, err := db.Exec(ctx, `INSERT INTO family_relationships (tree_id, family_id, relationship_id) VALUES ($1, $2, $3)`, f.tree, f.family, id); err != nil {
			t.Fatal(err)
		}
	}
	families := familyrepo.NewFamiliesRepository(db)
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback(ctx)
	if err := families.ReconcileTx(ctx, tx, f.tree, []string{f.family}); err != nil {
		t.Fatal(err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	f.lifecycle = repository.NewRepository(db, relrepo.NewRelationshipsRepository(db), personrepo.NewPersonsRepository(db), families)
	return f
}

func (f fixture) checkMember(t *testing.T, person, role string, want bool) {
	t.Helper()
	var exists bool
	if err := f.db.QueryRow(context.Background(), `
		SELECT EXISTS(SELECT 1 FROM family_members WHERE family_id = $1 AND person_id = $2 AND role = $3::family_member_role)
	`, f.family, person, role).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists != want {
		t.Fatalf("member %s as %s: got %t, want %t", person, role, exists, want)
	}
}

func (f fixture) checkFamily(t *testing.T, want bool) {
	t.Helper()
	var exists bool
	if err := f.db.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM families WHERE id = $1)`, f.family).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists != want {
		t.Fatalf("family exists = %t, want %t", exists, want)
	}
}

func TestRelationshipDeletionReconcilesFamily(t *testing.T) {
	f := setupLifecycle(t)
	ctx := context.Background()
	if err := f.lifecycle.DeleteRelationship(ctx, f.tree, f.childLinks[0]); err != nil {
		t.Fatal(err)
	}
	f.checkMember(t, f.parents[0], "parent", false)
	f.checkMember(t, f.parents[0], "partner", true)
	f.checkMember(t, f.parents[1], "parent", true)
	f.checkMember(t, f.child, "child", true)
	if err := f.lifecycle.DeleteRelationship(ctx, f.tree, f.childLinks[1]); err != nil {
		t.Fatal(err)
	}
	f.checkFamily(t, true)
	f.checkMember(t, f.child, "child", false)
	if err := f.lifecycle.DeleteRelationship(ctx, f.tree, f.spouse); err != nil {
		t.Fatal(err)
	}
	f.checkFamily(t, false)
	if err := f.lifecycle.DeleteRelationship(ctx, f.tree, f.spouse); !errors.Is(err, coreerrors.ErrNotFound) {
		t.Fatalf("delete missing relationship: %v", err)
	}
}

func TestPersonDeletionReconcilesFamily(t *testing.T) {
	f := setupLifecycle(t)
	ctx := context.Background()
	if err := f.lifecycle.DeletePerson(ctx, f.tree, f.parents[0]); err != nil {
		t.Fatal(err)
	}
	f.checkFamily(t, true)
	f.checkMember(t, f.parents[0], "partner", false)
	f.checkMember(t, f.parents[1], "partner", false)
	f.checkMember(t, f.parents[1], "parent", true)
	f.checkMember(t, f.child, "child", true)
	if err := f.lifecycle.DeletePerson(ctx, f.tree, f.parents[1]); err != nil {
		t.Fatal(err)
	}
	f.checkFamily(t, false)
}

func TestChildDeletionKeepsSpouseFamily(t *testing.T) {
	f := setupLifecycle(t)
	if err := f.lifecycle.DeletePerson(context.Background(), f.tree, f.child); err != nil {
		t.Fatal(err)
	}
	f.checkFamily(t, true)
	f.checkMember(t, f.child, "child", false)
	f.checkMember(t, f.parents[0], "partner", true)
	f.checkMember(t, f.parents[1], "partner", true)
	f.checkMember(t, f.parents[0], "parent", false)
	f.checkMember(t, f.parents[1], "parent", false)
}

func TestManualFamilyPreservesMembership(t *testing.T) {
	f := setupLifecycle(t)
	ctx := context.Background()
	if _, err := f.db.Exec(ctx, `UPDATE families SET auto_created = false WHERE id = $1`, f.family); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{f.spouse, f.childLinks[0], f.childLinks[1]} {
		if err := f.lifecycle.DeleteRelationship(ctx, f.tree, id); err != nil {
			t.Fatal(err)
		}
	}
	f.checkFamily(t, true)
	f.checkMember(t, f.child, "child", true)
}
