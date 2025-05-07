package bitio

import (
	"errors"
	"io"
)

var ErrWriterClosed = errors.New("writer is closed")

type Writer struct {
	w         io.ByteWriter
	stickyErr error
	cnt       int
	outByte   byte
}

func NewWriter(w io.ByteWriter) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Cnt() int {
	return w.cnt
}

func (w *Writer) flush() error {
	if err := w.w.WriteByte(w.outByte); err != nil {
		w.stickyErr = err
		return err
	}
	w.outByte = 0
	return nil
}

func (w *Writer) Close() error {
	w.stickyErr = ErrWriterClosed
	if w.cnt&7 != 0 {
		return w.flush()
	}
	return nil
}

func (w *Writer) WriteBit(v uint) error {
	if w.stickyErr != nil {
		return w.stickyErr
	}

	w.outByte |= byte(v&1) << (w.cnt & 7)
	w.cnt++

	if w.cnt&7 == 0 {
		return w.flush()
	}

	return nil
}

func (w *Writer) WriteBits(v uint, n int) error {
	for i := 0; i < n; i++ {
		if err := w.WriteBit(v); err != nil {
			return err
		}
		v >>= 1
	}
	return nil
}

type Reader struct {
	r         io.ByteReader
	stickyErr error
	cnt       int
	maxBitCnt int
	inByte    byte
}

func NewReader(r io.ByteReader) *Reader {
	return &Reader{r: r, maxBitCnt: -1}
}

func (r *Reader) Cnt() int {
	return r.cnt
}

func (r *Reader) MaxBitCnt(n int) {
	r.maxBitCnt = n
}

func (r *Reader) ReadBit() (uint, error) {
	if r.stickyErr != nil {
		return 0, r.stickyErr
	}

	if r.cnt == r.maxBitCnt {
		return 0, io.EOF
	}

	if r.cnt&7 == 0 {
		x, err := r.r.ReadByte()
		if err != nil {
			r.stickyErr = err
			return 0, err
		}
		r.inByte = x
	}

	v := uint(r.inByte>>(r.cnt&7)) & 1
	r.cnt++

	return v, nil
}

func (r *Reader) ReadBits(n int) (uint, error) {
	var v uint
	for i := 0; i < n; i++ {
		x, err := r.ReadBit()
		if err != nil {
			return v, err
		}
		v |= x << i
	}
	return v, nil
}
