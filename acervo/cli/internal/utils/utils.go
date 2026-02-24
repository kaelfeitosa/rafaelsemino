package utils

import "strings"

// CleanWikilink removes [[ and ]] from a string
func CleanWikilink(s string) string {
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")
	return s
}
