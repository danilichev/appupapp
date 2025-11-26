package errors

import (
	z "github.com/Oudwins/zog"
)

type ValidationError struct {
	Issues  *z.ZogIssueMap `json:"-"`
	Message string         `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(issues *z.ZogIssueMap) *ValidationError {
	return &ValidationError{
		Issues:  issues,
		Message: "Validation failed",
	}
}
