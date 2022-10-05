package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/raito-io/cli/base/data_usage"
)

type SimpleDataUsageStatementHandler struct {
	*DataUsageStatementHandler
	Statements []data_usage.Statement
}

func NewSimpleDataUsageStatementHandler(t mockConstructorTestingTNewDataUsageStatementHandler) *SimpleDataUsageStatementHandler {
	result := &SimpleDataUsageStatementHandler{
		DataUsageStatementHandler: NewDataUsageStatementHandler(t),
		Statements:                make([]data_usage.Statement, 0),
	}

	result.EXPECT().AddStatements(mock.AnythingOfType("[]data_usage.Statement")).Run(func(statements []data_usage.Statement) {
		result.Statements = append(result.Statements, statements...)
	}).Return(nil)

	return result
}
