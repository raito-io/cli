package role_based

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	//Given
	l1 := []string{"a", "b", "c"}
	l2 := []string{"a", "e", "c"}

	//When
	t1 := find(l1, "b")
	t2 := find(l2, "b")

	//Then
	assert.True(t, t1)
	assert.False(t, t2)
}
