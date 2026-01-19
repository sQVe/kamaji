package output

import (
	"bytes"
	"io"
)

// Writer wraps an io.Writer with styled output.
type Writer struct {
	w       io.Writer
	msgType MessageType
	buf     []byte
}

// NewWriter creates a styled writer that prefixes each line.
func NewWriter(w io.Writer, t MessageType) *Writer {
	return &Writer{
		w:       w,
		msgType: t,
	}
}

// NewInfoWriter creates a writer with Info styling.
func NewInfoWriter(w io.Writer) *Writer {
	return NewWriter(w, Info)
}

// NewErrorWriter creates a writer with Error styling.
func NewErrorWriter(w io.Writer) *Writer {
	return NewWriter(w, Error)
}

// Write implements io.Writer, prefixing each complete line with styled prefix.
func (w *Writer) Write(p []byte) (n int, err error) {
	n = len(p)
	w.buf = append(w.buf, p...)

	for {
		idx := bytes.IndexByte(w.buf, '\n')
		if idx < 0 {
			break
		}

		line := w.buf[:idx]
		w.buf = w.buf[idx+1:]

		prefix := Prefix(w.msgType)
		if _, err := w.w.Write([]byte(prefix)); err != nil {
			return n, err
		}
		if _, err := w.w.Write(line); err != nil {
			return n, err
		}
		if _, err := w.w.Write([]byte{'\n'}); err != nil {
			return n, err
		}
	}

	return n, nil
}

// Flush writes any buffered partial line without a trailing newline.
func (w *Writer) Flush() error {
	if len(w.buf) == 0 {
		return nil
	}

	prefix := Prefix(w.msgType)
	if _, err := w.w.Write([]byte(prefix)); err != nil {
		return err
	}
	if _, err := w.w.Write(w.buf); err != nil {
		return err
	}
	w.buf = nil
	return nil
}
