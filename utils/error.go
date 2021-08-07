package utils

import "fmt"

type CustomError struct {
	Process string
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error in %s. %s: %s", e.Process, e.Message, e.Err)
}