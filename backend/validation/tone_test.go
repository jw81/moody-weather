package validation

import "testing"

func TestIsValidTone(t *testing.T) {
    tests := []struct {
        tone  string
        valid bool
    }{
        {"nice", true},
        {"normal", true},
        {"snarky", true},
        {"angry", false},
        {"", false},
    }

    for _, tt := range tests {
        t.Run(tt.tone, func(t *testing.T) {
            if got := IsValidTone(tt.tone); got != tt.valid {
                t.Errorf("IsValidTone(%q) = %v; want %v", tt.tone, got, tt.valid)
            }
        })
    }
}
