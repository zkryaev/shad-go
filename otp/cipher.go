//go:build !solution

package otp

import (
	"errors"
	"io"
)

type DataReadWriter struct {
	InputCipherText  io.Reader
	CipherKeys       io.Reader
	OutputCipherText io.Writer
}

func (d *DataReadWriter) Read(p []byte) (n int, err error) {
	if len(p) == 0 { // nothing to do
		return 0, nil
	}
	n, err = d.InputCipherText.Read(p)
	switch {
	case n == 0 && err == io.EOF:
		return n, err
	case err != nil && err != io.EOF:
		return n, err
	}
	keys := make([]byte, n)
	d.CipherKeys.Read(keys)
	// doesn't return error
	for i := 0; i < n; i++ {
		p[i] ^= keys[i]
	}
	return n, err
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &DataReadWriter{
		InputCipherText:  r,
		CipherKeys:       prng,
		OutputCipherText: nil,
	}
}

func (d *DataReadWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 { // nothing to write
		return 0, nil
	}
	keys := make([]byte, len(p))
	n, _ = d.CipherKeys.Read(keys)
	if n < len(p) {
		return n, errors.New("Not enough keys")
	}
	encryptData := make([]byte, len(p))
	for i := 0; i < len(p); i++ {
		encryptData[i] = p[i] ^ keys[i]
	}
	n, err = d.OutputCipherText.Write(encryptData)
	return n, err
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &DataReadWriter{
		InputCipherText:  nil,
		CipherKeys:       prng,
		OutputCipherText: w,
	}
}
