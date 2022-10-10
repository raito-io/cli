package role_based

import (
	"strings"

	"github.com/raito-io/cli/base"
)

var logger = base.Logger()

func find(s []string, q string) bool {
	for _, r := range s {
		if strings.EqualFold(r, q) {
			return true
		}
	}

	return false
}
