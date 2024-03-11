package match

import (
	"regexp"
	"strings"
)

// MatchesAny returns true if the string matches any of the regex patterns
func MatchesAny(s string, regexPatters []string) (bool, error) {
	for _, p := range regexPatters {
		if !strings.HasPrefix(s, "^") {
			p = "^" + p
		}

		if !strings.HasSuffix(s, "$") {
			p += "$"
		}

		match, err := regexp.Match(p, []byte(s))
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}

	return false, nil
}
