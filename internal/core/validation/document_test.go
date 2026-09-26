package validation

import "testing"

const validDocumentOwnerID = "01ARZ3NDEKTSV4RRFFQ69G5FAV"

func TestDocumentValidation(t *testing.T) {
	if err := ValidateDocumentOwnerType("person"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDocumentOwnerType("unknown"); err != ErrInvalidDocumentOwner {
		t.Fatalf("ValidateDocumentOwnerType() error = %v", err)
	}
	if err := ValidateDocumentOwnerID(validDocumentOwnerID); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDocumentOwnerID("invalid"); err != ErrInvalidDocumentOwner {
		t.Fatalf("ValidateDocumentOwnerID() error = %v", err)
	}
	if err := ValidateDocument("application/pdf", MaxDocumentSize); err != nil {
		t.Fatal(err)
	}
	if err := ValidateDocument("application/pdf", MaxDocumentSize+1); err != ErrDocumentTooLarge {
		t.Fatalf("ValidateDocument() error = %v", err)
	}
	if err := ValidateDocument("text/plain", 10); err != ErrInvalidDocument {
		t.Fatalf("ValidateDocument() error = %v", err)
	}
}

func TestNormalizeDocumentFileName(t *testing.T) {
	name, err := NormalizeDocumentFileName("  archive.pdf  ")
	if err != nil || name != "archive.pdf" {
		t.Fatalf("NormalizeDocumentFileName() = %q, %v", name, err)
	}
}
