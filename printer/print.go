// Package printer is replacement for fmt.* functions
package printer

import (
	"fmt"
	"io"
	"log"
)

// Println wraps fmt.Println and handles errors
func Println(v ...any) {
	_, err := fmt.Println(v...)
	if err != nil {
		log.Printf("Println error: %v", err)
	}
}

// Printf wraps fmt.Printf and handles errors
func Printf(format string, v ...any) {
	_, err := fmt.Printf(format, v...)
	if err != nil {
		log.Printf("Printf error: %v", err)
	}
}

// Print wraps fmt.Print and handles errors
func Print(v ...any) {
	_, err := fmt.Print(v...)
	if err != nil {
		log.Printf("Print error: %v", err)
	}
}

// Errorf wraps fmt.Errorf and returns a formatted error message
func Errorf(format string, v ...any) error {
	return fmt.Errorf(format, v...)
}

// Fprint wraps fmt.Fprint and writes to the provided io.Writer, handling errors
func Fprint(w io.Writer, v ...any) {
	_, err := fmt.Fprint(w, v...)
	if err != nil {
		log.Printf("Fprint error: %v", err)
	}
}

// Fprintf wraps fmt.Fprintf and writes formatted output to the provided io.Writer, handling errors
func Fprintf(w io.Writer, format string, v ...any) {
	_, err := fmt.Fprintf(w, format, v...)
	if err != nil {
		log.Printf("Fprintf error: %v", err)
	}
}

// Fprintln wraps fmt.Fprintln and writes to the provided io.Writer, handling errors
func Fprintln(w io.Writer, v ...any) {
	_, err := fmt.Fprintln(w, v...)
	if err != nil {
		log.Printf("Fprintln error: %v", err)
	}
}
