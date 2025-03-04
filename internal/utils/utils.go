package utils

import "strings"

func SanitizeReleaseName(name string) string {
	chars := []string{" ", "?", "\n", "\r", "\t", "=", "*", "(", ")", "&", "%",
		"$", "#", "@", "+", "!"}
	result := name

	for _, char := range chars {
		result = strings.ReplaceAll(result, char, "")
	}

	return result
}
