package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Convert names to "safer" paths using camelCase
// It could be a bit confusing, but it is a simple way to handle
// spaces in paths.
// "This is a Test" -> "thisIsATest"
func ToCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(strings.ToLower(words[0]))

	caser := cases.Title(language.English)

	for _, word := range words[1:] {
		// result.WriteString(strings.Title(word))
		result.WriteString(caser.String(word))
	}

	return result.String()
}
