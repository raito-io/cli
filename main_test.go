package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "v1.2.3 (yyyy-mm-dd)", buildVersion("v1.2.3", "yyyy-mm-dd"))
}
