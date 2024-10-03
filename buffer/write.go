package buffer

import (
	"bytes"
	"log"
)

// WriteToBuffer writes to bytes.Buffer and logs any error.
func Write(buffer *bytes.Buffer, value string) {
	_, err := buffer.Write([]byte(value))
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
}
