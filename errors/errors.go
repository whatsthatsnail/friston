package errors

import "fmt"

// TODO: Report column of error, print line of error, etc.

// Prints error message and sets error flag
func ThrowError(line int, message string) {
	Report(line, "Error: " + message)
}

// Print any line dependant message (error, warning, etc.)
func Report(line int, message string) {
	fmt.Printf("[Line %d] %s\n", line, message)
}