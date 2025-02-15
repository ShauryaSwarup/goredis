package resp

import (
	"bufio"
	"fmt"
)

type Writer struct {
	bufw *bufio.Writer
}

func NewWriter(w *bufio.Writer) *Writer {
	return &Writer{bufw: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	nn, err := w.bufw.Write(bytes)
	// FLUSHING IS IMPORTANT TO WRITE 😭
	w.bufw.Flush()
	if err != nil {
		return err
	}
	fmt.Println("Bytes written: ", nn)
	return nil
}
