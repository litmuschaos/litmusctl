package utils

import (
	"strings"
	"testing"
)

func isValidString(s string) bool {
	// Define the set of valid characters
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"

	// Iterate over each character in the string
	for _, char := range s {
		// Check if the character is not in the set of valid characters
		if !strings.ContainsRune(validChars, char) {
			return false
		}
	}
	return true
}

func FuzzGenerateRandomString(f *testing.F) {
	f.Add(10)
	f.Fuzz(func(t *testing.T, n int) {
		randomString, err := GenerateRandomString(n)
		if err != nil && n > 0 {
			t.Errorf("Error generating random string: %v", err)
		}
		// Perform checks on the generated string
		// Check if the length matches the expected length
		if n >= 0 && len(randomString) != n {
			t.Errorf("Generated string length doesn't match expected length")
		}

		// Check if the string contains only valid characters
		if !isValidString(randomString) {
			t.Errorf("Generated string contains invalid characters")
		}
	})

}
