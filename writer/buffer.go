package writer

import (
	"bytes"
	"log"
)

// WriteToBuffer writes to bytes.Buffer and logs any error.
func WriteToBuffer(buffer *bytes.Buffer, data []byte) {
	_, err := buffer.Write(data)
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
}
