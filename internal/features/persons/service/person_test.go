package service

import (
	"strings"
	"testing"
)

func TestNormalizePatronymic(t *testing.T) {
	t.Run("allows null", func(t *testing.T) {
		value, err := normalizePatronymic(nil)
		if err != nil || value != nil {
			t.Fatalf("value = %v, error = %v", value, err)
		}
	})

	t.Run("trims a value", func(t *testing.T) {
		input := "  Петрович  "
		value, err := normalizePatronymic(&input)
		if err != nil || value == nil || *value != "Петрович" {
			t.Fatalf("value = %v, error = %v", value, err)
		}
	})

	for _, input := range []string{"   ", strings.Repeat("я", 201)} {
		if _, err := normalizePatronymic(&input); err != ErrInvalid {
			t.Fatalf("normalizePatronymic(%q) error = %v, want ErrInvalid", input, err)
		}
	}
}
