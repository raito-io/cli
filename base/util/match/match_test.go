package match

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchesAny(t *testing.T) {
	singleCheck(t, "test", nil, false, false)
	singleCheck(t, "test", []string{}, false, false)
	singleCheck(t, "test", []string{"test"}, true, false)
	singleCheck(t, "test-a", []string{"test-a"}, true, false)
	singleCheck(t, "test", []string{"blah"}, false, false)
	singleCheck(t, "test", []string{"blah", "test"}, true, false)
	singleCheck(t, "test", []string{"blah", "t(.)+st"}, true, false)
	singleCheck(t, "pretestpost", []string{"blah", "t(.)+st"}, false, false)
	singleCheck(t, "teeeeeeest", []string{"blah", "t(.)+st"}, true, false)
	singleCheck(t, "test", []string{"(("}, false, true)

	// From the documentation:
	singleCheck(t, "test-prod", []string{".+-prod", ".+-dev"}, true, false)
	singleCheck(t, "test-dev", []string{".+-prod", ".+-dev"}, true, false)
	singleCheck(t, "-dev", []string{".+-prod", ".+-dev"}, false, false)
	singleCheck(t, "test-staging", []string{".+-prod", ".+-dev"}, false, false)
	singleCheck(t, "test_dev", []string{".+-prod", ".+-dev"}, false, false)
}

func singleCheck(t *testing.T, s string, p []string, expected bool, errorExpected bool) {
	m, e := MatchesAny(s, p)
	if errorExpected {
		assert.Error(t, e)
	} else {
		assert.NoError(t, e)
		assert.Equal(t, expected, m)
	}
}
