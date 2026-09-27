package validation

import "testing"

func TestValidateULID(t *testing.T) {
	if err := ValidateULID("01ARZ3NDEKTSV4RRFFQ69G5FAV"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateULID("invalid"); err != ErrInvalidULID {
		t.Fatalf("ValidateULID() error = %v, want ErrInvalidULID", err)
	}
}
