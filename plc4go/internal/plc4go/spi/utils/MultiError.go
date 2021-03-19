package utils

import (
	"fmt"
	"strings"
)

// MultiError is a Wrapper for multiple Errors
type MultiError struct {
	// MainError denotes the error which summarize the error
	MainError error
	// Errors are the child errors
	Errors []error
}

func (m MultiError) Error() string {
	mainErrorText := "Child errors:\n"
	if m.MainError != nil {
		mainErrorText = fmt.Sprintf("Main Error: %v\nChild errors:\n", m.MainError)
	}
	return mainErrorText + strings.Join(func(errors []error) []string {
		result := make([]string, len(errors))
		for i, errorElement := range errors {
			result[i] = errorElement.Error()
		}
		return result
	}(m.Errors), "\n")
}
