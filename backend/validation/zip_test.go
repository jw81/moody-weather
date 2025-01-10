package validation

import "testing"

func TestIsValidZipCode(t *testing.T) {
	tests := []struct {
		zipCode string
		valid   bool
	}{
		{"12345", true},  // valid
		{"abcde", false}, // invalid (non-numeric)
		{"1234", false},  // invalid (too short)
		{"123456", false}, // invalid (too long)
	}

	for _, tt := range tests {
		t.Run(tt.zipCode, func(t *testing.T) {
			if got := IsValidZipCode(tt.zipCode); got != tt.valid {
				t.Errorf("IsValidZipCode(%q) = %v; want %v", tt.zipCode, got, tt.valid)
			}
		})
	}
}
