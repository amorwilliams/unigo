package utils

import (
	l4g "code.google.com/p/log4go"
	"crypto/rand"
	"fmt"
	"io"
)

// NewUUID() provides unique identifier strings.
func NewUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		l4g.Error(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
