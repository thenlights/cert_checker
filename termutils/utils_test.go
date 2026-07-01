package termutils

import "testing"

func TestYesNo(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"yes", Green},
		{"no", Red},
		{"maybe", Magenta},
		{"", Magenta},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := YesNo(tt.input); got != tt.want {
				t.Errorf("YesNo(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpirationColor(t *testing.T) {
	tests := []struct {
		name      string
		expiresIn int
		want      string
	}{
		{"far from expiry", 60, Green},
		{"just above green threshold", 31, Green},
		{"yellow lower bound", 15, Yellow},
		{"yellow upper bound", 30, Yellow},
		{"red below yellow", 14, Red},
		{"expired", -3, Red},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpirationColor(tt.expiresIn); got != tt.want {
				t.Errorf("ExpirationColor(%d) = %q, want %q", tt.expiresIn, got, tt.want)
			}
		})
	}
}
