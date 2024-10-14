package string

import "strings"

func TrimSpaceInCommaSeparatedList(list string) string {
	trimmed := strings.Split(list, ",")
	for i, s := range trimmed {
		trimmed[i] = strings.TrimSpace(s)
	}

	return strings.Join(trimmed, ",")
}
