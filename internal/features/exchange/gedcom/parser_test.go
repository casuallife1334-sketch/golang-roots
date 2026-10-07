package gedcom

import (
	"genealogy-tree/internal/core/domain"
	"strings"
	"testing"
)

func TestParseBuildsPersonsAndRelationships(t *testing.T) {
	snapshot, err := Parse(strings.NewReader(`0 HEAD
1 CHAR UTF-8
0 @I1@ INDI
1 NAME Ivan Petrovich /Petrov/
2 GIVN Ivan Petrovich
2 SURN Petrov
1 SEX M
1 BIRT
2 DATE 2 JAN 1980
0 @I2@ INDI
1 NAME Anna /Sidorova/
1 SEX F
0 @I3@ INDI
1 NAME Petr /Petrov/
1 BIRT
2 DATE 1985
0 @F1@ FAM
1 HUSB @I1@
1 WIFE @I2@
1 CHIL @I3@
1 MARR
2 DATE 1 MAY 1979
2 PLAC Moscow
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Errors) != 0 {
		t.Fatalf("unexpected parse errors: %+v", snapshot.Errors)
	}
	if len(snapshot.Persons) != 3 {
		t.Fatalf("persons = %d, want 3", len(snapshot.Persons))
	}
	if len(snapshot.Relationships) != 3 {
		t.Fatalf("relationships = %d, want 3", len(snapshot.Relationships))
	}
	if snapshot.Persons[0].Input.FirstName != "Ivan" || snapshot.Persons[0].Input.LastName != "Petrov" {
		t.Fatalf("unexpected first person: %+v", snapshot.Persons[0].Input)
	}
	if snapshot.Persons[0].Input.Patronymic == nil || *snapshot.Persons[0].Input.Patronymic != "Petrovich" {
		t.Fatalf("unexpected patronymic: %+v", snapshot.Persons[0].Input.Patronymic)
	}
	if snapshot.Persons[0].Input.Gender == nil || *snapshot.Persons[0].Input.Gender != domain.GenderMale {
		t.Fatalf("unexpected gender: %+v", snapshot.Persons[0].Input.Gender)
	}
	if snapshot.Persons[0].Input.BirthDate == nil || snapshot.Persons[0].Input.BirthDate.Year() != 1980 {
		t.Fatalf("unexpected birth date: %v", snapshot.Persons[0].Input.BirthDate)
	}
	if snapshot.Relationships[0].Input.Type != domain.RelationshipSpouse {
		t.Fatalf("first relationship type = %s, want spouse", snapshot.Relationships[0].Input.Type)
	}
	if snapshot.Relationships[1].Input.Type != domain.RelationshipParentChild {
		t.Fatalf("second relationship type = %s, want parent_child", snapshot.Relationships[1].Input.Type)
	}
	gedcomMetadata, ok := snapshot.Relationships[0].Input.Metadata["gedcom"].(map[string]any)
	if !ok || gedcomMetadata["marriage_place"] != "Moscow" {
		t.Fatalf("unexpected marriage metadata: %+v", snapshot.Relationships[0].Input.Metadata)
	}
}

func TestParseReportsIncompleteNamesAndReferencesAsWarnings(t *testing.T) {
	snapshot, err := Parse(strings.NewReader(`0 HEAD
0 @I1@ INDI
1 NAME /Petrov/
0 @F1@ FAM
1 HUSB @I1@
1 CHIL @I404@
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", snapshot.Errors)
	}
	if len(snapshot.Warnings) < 2 {
		t.Fatalf("warnings = %+v, want incomplete name and reference warnings", snapshot.Warnings)
	}
}

func TestParseReportsUnsupportedDateAsWarning(t *testing.T) {
	snapshot, err := Parse(strings.NewReader(`0 HEAD
0 @I1@ INDI
1 NAME Ivan /Petrov/
1 BIRT
2 DATE ABT 1980
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Warnings) == 0 {
		t.Fatal("expected unsupported date warning")
	}
	if len(snapshot.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", snapshot.Errors)
	}
}

func TestParseIgnoresStandardMetadataAndNormalizesWrappedSurname(t *testing.T) {
	snapshot, err := Parse(strings.NewReader(`0 HEAD
0 @I1@ INDI
1 _UID external-id
1 NAME Ivan / (Petrov)/
2 GIVN Ivan
2 SURN Petrov
1 SEX M
1 RESI
2 PLAC Moscow
1 OBJE
2 FILE photo.jpg
3 FORM jpg
2 TITL Ivan
2 _PRIM Y
0 @F1@ FAM
1 HUSB @I1@
1 NOTE imported note
2 CONT continued note
0 TRLR
`))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Errors) != 0 {
		t.Fatalf("unexpected errors: %+v", snapshot.Errors)
	}
	if len(snapshot.Warnings) != 0 {
		t.Fatalf("unexpected warnings: %+v", snapshot.Warnings)
	}
	if snapshot.Persons[0].Input.LastName != "Petrov" {
		t.Fatalf("last name = %q, want Petrov", snapshot.Persons[0].Input.LastName)
	}
}
