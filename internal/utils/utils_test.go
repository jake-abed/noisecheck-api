package utils

import (
	"testing"
)

func TestSanitizeReleaseName(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"DreebyMan3", "DreebyMan3"},
		{"", ""},
		{"420**Doggy", "420Doggy"},
		{"Why? Who? Argh!", "WhyWhoArgh"},
		{"@$yee()", "yee"},
		{"1-@2&-=+dog-4", "1-2-dog-4"},
		{"yo!-b1rd@up-#hee!", "yo-b1rdup-hee"},
		{"311 is for me! 420 4 u?", "311isforme4204u"},
	}

	for _, tt := range tests {
		res := SanitizeReleaseName(tt.Input)

		if res != tt.Expected {
			t.Errorf("Did not sanitize: result=%s != expected=%s", res, tt.Expected)
		}
	}
}
