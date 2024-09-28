package utils

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) Print(a ...any) {
	fmt.Fprint(w.writer, a...)
}

func (w *Writer) Println(a ...any) {
	fmt.Fprintln(w.writer, a...)
}

func (w *Writer) Printf(format string, a ...any) {
	fmt.Fprintf(w.writer, format, a...)
}
