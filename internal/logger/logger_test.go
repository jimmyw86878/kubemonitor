package logger

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	NewLogger()
	msg := fmt.Sprintf("Error message example")
	Error.Println(msg)
}

func TestWarning(t *testing.T) {
	NewLogger()
	msg := fmt.Sprintf("Warning message example")
	Warning.Println(msg)
}

func TestInfo(t *testing.T) {
	NewLogger()
	msg := fmt.Sprintf("Info message example")
	Info.Println(msg)
}
