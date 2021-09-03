// utils package implements objects commonly used in other packages.
package utils

import "fmt"

// CustomError wraps an error with additional information
type CustomError struct {
	Process string // In which process the error took place
	Message string // Detail of the error
	Err     error  // Wrapping an original error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error in %s. %s: %s", e.Process, e.Message, e.Err)
}

// AsMessage method creates a message to notify the error to Slack
func (e *CustomError) AsMessage() string {
	return fmt.Sprintf("Process Failed at %s!", e.Process)
}
