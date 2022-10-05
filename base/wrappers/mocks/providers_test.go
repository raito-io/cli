package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raito-io/cli/base/data_usage"
)

func TestNewSimpleDataUsageStatementHandler(t *testing.T) {
	//Given
	statements1 := []data_usage.Statement{{User: "user1", Credits: 123456}, {User: "user2", Credits: 234567}}
	statements2 := []data_usage.Statement{{User: "user3", Credits: 3141592}}

	mock := NewSimpleDataUsageStatementHandler(t)
	err := mock.AddStatements(statements1)

	assert.NoError(t, err)
	assert.Len(t, mock.Statements, 2)
	assert.Equal(t, statements1, mock.Statements)

	err = mock.AddStatements(statements2)

	assert.NoError(t, err)
	assert.Len(t, mock.Statements, 3)
	assert.Equal(t, statements1, mock.Statements[0:2])
	assert.Equal(t, statements2, mock.Statements[2:])
}
