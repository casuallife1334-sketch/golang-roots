package request

import (
	"errors"
	coreerrors "genealogy-tree/internal/core/errors"
	"genealogy-tree/internal/core/security"
	"net/http/httptest"
	"testing"
)

const testULID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"

func TestGetTreeContext(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("tree_id", testULID)
	r = r.WithContext(security.WithUserID(r.Context(), "user-id"))

	userID, treeID, err := GetTreeContext(r)
	if err != nil {
		t.Fatalf("GetTreeContext() error = %v", err)
	}
	if userID != "user-id" || treeID != testULID {
		t.Fatalf("GetTreeContext() = (%q, %q), want (%q, %q)", userID, treeID, "user-id", testULID)
	}
}

func TestGetTreeContextRequiresUser(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("tree_id", testULID)

	_, _, err := GetTreeContext(r)
	if !errors.Is(err, coreerrors.ErrUnauthorized) {
		t.Fatalf("GetTreeContext() error = %v, want ErrUnauthorized", err)
	}
}

func TestGetTreeIDValidatesPathValue(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("tree_id", "invalid")

	_, err := GetTreeID(r)
	if !errors.Is(err, coreerrors.ErrInvalidArgument) {
		t.Fatalf("GetTreeID() error = %v, want ErrInvalidArgument", err)
	}
}
