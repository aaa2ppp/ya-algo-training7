package main

import "io"

type BitWriter struct {
	w         io.ByteWriter
	stickyErr error
	cnt       int
	outByte   byte
}

func NewBitWriter(w io.ByteWriter) *BitWriter {
	return &BitWriter{w: w}
}

func (w *BitWriter) Cnt() int {
	return w.cnt
}

func (w *BitWriter) WriteBit(v uint) error {
	if w.stickyErr != nil {
		return w.stickyErr
	}

	w.outByte |= byte(v&1) << (w.cnt & 7)
	w.cnt++

	if w.cnt&7 == 0 {
		if err := w.w.WriteByte(w.outByte); err != nil {
			w.stickyErr = err
			return err
		}
		w.outByte = 0
	}

	return nil
}

func (w *BitWriter) WriteBits(v uint, n int) error {
	for i := 0; i < n; i++ {
		if err := w.WriteBit(v); err != nil {
			return err
		}
		v >>= 1
	}
	return nil
}

type BitReader struct {
	r         io.ByteReader
	stickyErr error
	cnt       int
	maxBitCnt int
	inByte    byte
}

func NewBitReader(r io.ByteReader) *BitReader {
	return &BitReader{r: r, maxBitCnt: -1}
}

func (r *BitReader) Cnt() int {
	return r.cnt
}

func (r *BitReader) MaxBitCnt(n int) {
	r.maxBitCnt = n
}

func (r *BitReader) ReadBit() (uint, error) {
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

func (r *BitReader) ReadBits(n int) (uint, error) {
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
