package gedcom

import (
	"bytes"
	"genealogy-tree/internal/core/domain"
	"strings"
	"testing"
	"time"
)

func TestExportWritesIndividualsAndFamiliesFromRelationships(t *testing.T) {
	birthDate := time.Date(1980, time.January, 2, 0, 0, 0, 0, time.UTC)
	male := domain.GenderMale
	female := domain.GenderFemale
	persons := []domain.Person{
		{ID: "01", FirstName: "Ivan", LastName: "Petrov", Gender: &male, BirthDate: &birthDate},
		{ID: "02", FirstName: "Anna", LastName: "Sidorova", Gender: &female},
		{ID: "03", FirstName: "Petr", LastName: "Petrov"},
	}
	spouse := domain.Relationship{Person1ID: "01", Person2ID: "02", Type: domain.RelationshipSpouse, Metadata: map[string]any{}}
	parent1 := domain.Relationship{Person1ID: "01", Person2ID: "03", Type: domain.RelationshipParentChild, Direction: direction(domain.DirectionParent), Metadata: map[string]any{}}
	parent2 := domain.Relationship{Person1ID: "02", Person2ID: "03", Type: domain.RelationshipParentChild, Direction: direction(domain.DirectionParent), Metadata: map[string]any{}}

	var output bytes.Buffer
	if err := Export(persons, []domain.Relationship{spouse, parent1, parent2}, &output); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{
		"0 HEAD",
		"1 CHAR UTF-8",
		"0 @I1@ INDI",
		"1 NAME Ivan /Petrov/",
		"1 SEX M",
		"1 BIRT",
		"2 DATE 2 JAN 1980",
		"1 FAMS @F1@",
		"1 FAMC @F1@",
		"0 @F1@ FAM",
		"1 HUSB",
		"1 WIFE",
		"1 CHIL",
		"0 TRLR",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("export does not contain %q:\n%s", expected, text)
		}
	}
	parsed, err := Parse(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.Errors) != 0 {
		t.Fatalf("round-trip parse errors: %+v", parsed.Errors)
	}
	if len(parsed.Persons) != len(persons) || len(parsed.Relationships) != 3 {
		t.Fatalf("round-trip counts = %d persons, %d relationships", len(parsed.Persons), len(parsed.Relationships))
	}
}
