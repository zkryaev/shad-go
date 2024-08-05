//go:build !change

package externalsort

import (
	"io"
)

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}

type LinesReader struct {
	r io.Reader
}

type LinesWriter struct {
	w io.Writer
}

func (lrw *LinesReader) ReadLine() (string, error) {
	buf := make([]byte, 0, 10)
	oneByte := make([]byte, 1)
	var err error
	var n int
	n, err = lrw.r.Read(oneByte)
	for rune(oneByte[0]) != '\n' {
		if err != nil && err != io.EOF { // invalid read
			return string(buf), err
		}
		if n == 0 && err == io.EOF { // EOF
			return string(buf), io.EOF
		}
		buf = append(buf, oneByte[0])
		n, err = lrw.r.Read(oneByte)
	}
	return string(buf), err
}

func (lrw *LinesWriter) Write(l string) (err error) {
	_, err = lrw.w.Write([]byte(l))
	if err != nil {
		return err
	}
	lrw.w.Write([]byte("\n"))
	return nil
}
