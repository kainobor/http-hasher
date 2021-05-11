package src

import (
	"fmt"
	"strings"
)

type (
	Writer interface {
		Write(input ...string)
	}

	ioWriter struct{}
)

func NewIOWriter() *ioWriter {
	return &ioWriter{}
}

func (w *ioWriter) Write(input ...string) {
	output := strings.Join(input, " ")
	fmt.Println(output)
}
