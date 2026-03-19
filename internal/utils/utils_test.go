package utils

import "testing"

func TestNormalizeDurationString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "integer seconds", input: "3600", want: "01:00:00"},
		{name: "decimal seconds", input: "3661.2", want: "01:01:01"},
		{name: "surrounding spaces", input: " 65 ", want: "00:01:05"},
		{name: "non numeric", input: "01:02:03", want: "01:02:03"},
		{name: "negative numeric", input: "-15", want: "-15"},
		{name: "empty", input: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeDurationString(tt.input); got != tt.want {
				t.Fatalf("NormalizeDurationString(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
