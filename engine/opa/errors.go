package opa

import (
	"fmt"
	"github.com/open-policy-agent/opa/rego"
)

type UnexpectedResultExpressionError struct {
	exps []*rego.ExpressionValue
}

func (e *UnexpectedResultExpressionError) Error() string {
	return fmt.Sprintf("unexpected result expressions: %v", e.exps)
}

type UnexpectedResultSetError struct {
	set rego.ResultSet
}

func (e *UnexpectedResultSetError) Error() string {
	return fmt.Sprintf("unexpected result set: %v", e.set)
}

type EvaluationError struct {
	e error
}

func (e *EvaluationError) Error() string {
	return fmt.Sprintf("error in query evaluation: %s", e.e.Error())
}
