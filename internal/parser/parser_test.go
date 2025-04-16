package parser

import (
	"testing"
)

func TestParseDiffOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name: "single line change 1",
			input: `--- Original
					+++ New
					@@ -20,7 +20,7 @@
					}
 
					func doo() {
					-       ddd := 6
					+       ddd := 5
							slog.Info(strconv.Itoa(ddd))
							fmt.Println("doo")
					 }`,
			expected: 23,
		},
		{
			name: "single line change 2",
			input: `--- Original
					+++ New
					@@ -14,7 +14,7 @@
					 func foo() {
							jjj := 6
					-       slog.Info(strconv.Itoa(jjj))
					+       _, _, _ = slog.Info, strconv.Itoa, jjj
					 
							fmt.Println("foo")
					 }`,
			expected: 17,
		},
		{
			name: "multiple changes should return fallback",
			input: `--- Original
					+++ New
					@@ -14,7 +14,7 @@
					 func foo() {
							jjj := 6
					-       slog.Info(strconv.Itoa(jjj))
					+       _, _, _ = slog.Info, strconv.Itoa, jjj
					 
							fmt.Println("foo")
					 }
					@@ -20,7 +20,7 @@
					}
 
					func doo() {
					-       ddd := 6
					+       ddd := 5
							slog.Info(strconv.Itoa(ddd))
							fmt.Println("doo")
					 }`,
			expected: 0,
		},
		{
			name:     "empty input",
			input:    "",
			expected: 0,
		},
		{
			name: "invalid line numbers",
			input: `@@ -abc +def @@
					-garbage
					+garbage`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindOriginalStartLine([]byte(tt.input))
			if got != tt.expected {
				t.Errorf("FindOriginalStartLine() = %v, want %v", got, tt.expected)
			}
		})
	}
}
