package logger

import "fmt"

// Logger defines an interface on what to do about logging messages.
// The user can decide whether to log debug messages or not.
// In the code there is just one logger which is either a verbose
// logger or a quite one.
type Logger interface {
	Log(message string, args ...any)
	Debug(message string, args ...any)
}

// VerboseLogger logs debug messages.
type VerboseLogger struct{}

// Log just logs normal messages.
func (*VerboseLogger) Log(message string, args ...any) {
	fmt.Printf(message, args...)
	fmt.Println()
}

// Debug is used for messages which can normally be ignored.
func (*VerboseLogger) Debug(message string, args ...any) {
	fmt.Printf(message, args...)
	fmt.Println()
}

// QuiteLogger 's LogDebug is ignored.
type QuiteLogger struct{}

// Log just logs normal messages.
func (*QuiteLogger) Log(message string, args ...any) {
	fmt.Printf(message, args...)
	fmt.Println()
}

// Debug is ignored.
func (*QuiteLogger) Debug(message string, args ...any) {
	// I'm quite.
}
