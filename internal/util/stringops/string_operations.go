package stringops

import "strings"

func TrimSpaceInCommaSeparatedList(list string) string {
	trimmed := strings.Split(list, ",")
	result := make([]string, 0, len(trimmed))

	for _, s := range trimmed {
		v := strings.TrimSpace(s)

		if v == "" {
			continue
		}

		result = append(result, v)
	}

	return strings.Join(result, ",")
}
