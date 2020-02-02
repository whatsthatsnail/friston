package errors

import "fmt"

// TODO: Report column of error, print line of error, etc.
// TODO: Differentiate between syntax and runtime errors, stop execution with runtime errors, etc.

// Prints error message and sets error flag
func ThrowError(line int, message string) {
	report(line, "Error: " + message)
}

// Print any line dependant message (error, warning, etc.)
func report(line int, message string) {
	if line == 0 {
		fmt.Printf("%s\n", message)
	} else {
		fmt.Printf("[Line %d] %s\n", line, message)
	}
}