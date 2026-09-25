package service

import (
	"genealogy-tree/internal/core/domain"
	"testing"
)

func TestCanWriteTree(t *testing.T) {
	for _, test := range []struct {
		role domain.TreeRole
		want bool
	}{
		{role: domain.TreeRoleOwner, want: true},
		{role: domain.TreeRoleEditor, want: true},
		{role: domain.TreeRoleViewer, want: false},
	} {
		if got := canWriteTree(test.role); got != test.want {
			t.Fatalf("canWriteTree(%q) = %v, want %v", test.role, got, test.want)
		}
	}
}
