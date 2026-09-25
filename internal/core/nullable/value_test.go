package nullable

import (
	"encoding/json"
	"testing"
)

func TestValue(t *testing.T) {
	tests := []struct {
		name      string
		payload   string
		set       bool
		value     string
		hasValue  bool
		wantError bool
	}{
		{name: "omitted", payload: `{}`, set: false},
		{name: "null", payload: `{"field":null}`, set: true},
		{name: "value", payload: `{"field":"Петрович"}`, set: true, value: "Петрович", hasValue: true},
		{name: "wrong type", payload: `{"field":42}`, wantError: true},
	}

	type payload struct {
		Field Value[string] `json:"field"`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var decoded payload
			err := json.Unmarshal([]byte(tt.payload), &decoded)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if decoded.Field.Set != tt.set {
				t.Fatalf("Set = %v, want %v", decoded.Field.Set, tt.set)
			}
			if (decoded.Field.Value != nil) != tt.hasValue {
				t.Fatalf("Value presence = %v, want %v", decoded.Field.Value != nil, tt.hasValue)
			}
			if tt.hasValue && *decoded.Field.Value != tt.value {
				t.Fatalf("Value = %q, want %q", *decoded.Field.Value, tt.value)
			}
		})
	}
}
