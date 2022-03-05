// Package url contains some helper functions to deal with URLs.
package url

import (
	"strings"
)

// CutOffPrefix simply cuts off a prefix from a string. It returns the original input string when it doesn't start with the prefix.
func CutOffPrefix(input string, prefix string) string {
	if strings.HasPrefix(input, prefix) {
		return input[len(prefix):]
	}
	return input
}

// CutOffSuffix simply cuts off a suffix from a string. It returns the original input string when it doesn't start with the suffix.
func CutOffSuffix(input string, suffix string) string {
	if strings.HasSuffix(input, suffix) {
		return input[:len(input)-len(suffix)]
	}
	return input
}

// CutOffSchema cuts off the schema part of a URL.
// For example: https://www.google.com would become www.google.com
func CutOffSchema(input string) string {
	if strings.Contains(input, "://") {
		return input[strings.Index(input, "://")+3:]
	}
	return input
}

// GetRelativePath returns the relative path part of a URL. It can handle both full URLs (with a schema) or just absolute paths.
// For example:
//  - https://www.google.com/my/path would become my/path
//  - /a/cool/path would become a/cool/path
func GetRelativePath(path string) string {
	relPath := CutOffSchema(path)
	if relPath != path && strings.Contains(relPath, "/") {
		relPath = relPath[strings.Index(relPath, "/")+1:]
	}
	return CutOffPrefix(relPath, "/")
}
