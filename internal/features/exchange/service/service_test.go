package service

import (
	"context"
	"errors"
	"genealogy-tree/internal/core/domain"
	"genealogy-tree/internal/features/exchange/gedcom"
	"strings"
	"testing"
)

type fakeRepository struct {
	importCalled bool
}

func (r *fakeRepository) GetTreeData(context.Context, string) ([]domain.Person, []domain.Relationship, error) {
	return nil, nil, nil
}

func (r *fakeRepository) ImportSnapshot(context.Context, string, gedcom.Snapshot) (int, int, error) {
	r.importCalled = true
	return 1, 1, nil
}

type fakeTreeAccess struct {
	writeErr error
}

func (a fakeTreeAccess) CanReadTree(context.Context, string, string) error { return nil }

func (a fakeTreeAccess) CanWriteTree(context.Context, string, string) error { return a.writeErr }

func TestImportRejectsParentCycleBeforeRepository(t *testing.T) {
	repository := &fakeRepository{}
	service := NewExchangeService(repository, fakeTreeAccess{})
	input := strings.NewReader(`0 HEAD
0 @I1@ INDI
1 NAME A /A/
0 @I2@ INDI
1 NAME B /B/
0 @F1@ FAM
1 HUSB @I1@
1 CHIL @I2@
0 @F2@ FAM
1 HUSB @I2@
1 CHIL @I1@
0 TRLR
`)
	_, _, err := service.Import(context.Background(), "user", "tree", input)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("error = %v, want ErrInvalid", err)
	}
	if repository.importCalled {
		t.Fatal("repository import was called for invalid snapshot")
	}
}

func TestPreviewDoesNotWrite(t *testing.T) {
	repository := &fakeRepository{}
	service := NewExchangeService(repository, fakeTreeAccess{})
	snapshot, err := service.Preview(strings.NewReader(`0 HEAD
0 @I1@ INDI
1 NAME Ivan /Petrov/
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Persons) != 1 || repository.importCalled {
		t.Fatalf("unexpected preview result: %+v", snapshot)
	}
}

func TestPreviewReportsParentCycle(t *testing.T) {
	service := NewExchangeService(&fakeRepository{}, fakeTreeAccess{})
	snapshot, err := service.Preview(strings.NewReader(`0 HEAD
0 @I1@ INDI
1 NAME A /A/
0 @I2@ INDI
1 NAME B /B/
0 @F1@ FAM
1 HUSB @I1@
1 CHIL @I2@
0 @F2@ FAM
1 HUSB @I2@
1 CHIL @I1@
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Errors) == 0 {
		t.Fatal("expected cycle error in preview")
	}
}
