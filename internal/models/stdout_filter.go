package models

import (
	"io"
	"strings"
)

type RequiredWordsWriter struct {
	Writer        io.Writer
	RequiredWords []string
}

func (wfw *RequiredWordsWriter) Write(p []byte) (n int, err error) {
	input := string(p)

	for _, word := range wfw.RequiredWords {
		if !strings.Contains(input, word) {
			return len(p), nil
		}
	}

	return wfw.Writer.Write(p)
}
