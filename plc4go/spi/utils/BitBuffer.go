/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

import (
	stdErrors "errors"
	"io"

	"github.com/pkg/errors"
)

// ReadBitBuffer reads bits from a fixed byte slice in big-endian bit order.
// Bit 0 is the MSB of byte 0; bit 7 is the LSB of byte 0; bit 8 is the MSB of byte 1, etc.
type ReadBitBuffer struct {
	data   []byte
	bitPos uint64
}

// NewReadBitBuffer creates a ReadBitBuffer backed by data.
func NewReadBitBuffer(data []byte) *ReadBitBuffer {
	return &ReadBitBuffer{data: data}
}

// ResetTo repositions the read cursor to the given bit offset.
func (b *ReadBitBuffer) ResetTo(bitPos uint64) {
	b.bitPos = bitPos
}

// BitsRemaining returns the number of bits that can still be read.
func (b *ReadBitBuffer) BitsRemaining() uint64 {
	total := uint64(len(b.data)) * 8
	if b.bitPos >= total {
		return 0
	}
	return total - b.bitPos
}

// ReadBool reads 1 bit.
func (b *ReadBitBuffer) ReadBool() (bool, error) {
	v, err := b.ReadBits(1)
	return v != 0, err
}

// ReadByte reads 8 bits as a byte.
func (b *ReadBitBuffer) ReadByte() (byte, error) {
	v, err := b.ReadBits(8)
	return byte(v), err
}

// Read reads len(p) bytes, implementing io.Reader.
func (b *ReadBitBuffer) Read(p []byte) (int, error) {
	for i := range p {
		v, err := b.ReadBits(8)
		if err != nil {
			return i, err
		}
		p[i] = byte(v)
	}
	return len(p), nil
}

// ReadBits reads n bits (0 ≤ n ≤ 64) in big-endian order, returning them right-aligned in a uint64.
func (b *ReadBitBuffer) ReadBits(n uint8) (uint64, error) {
	if n == 0 {
		return 0, nil
	}
	if n > 64 {
		return 0, errors.New("cannot read more than 64 bits at once")
	}
	if uint64(n) > b.BitsRemaining() {
		return 0, io.ErrUnexpectedEOF
	}

	var result uint64
	bitsLeft := n

	// Partial first byte when the cursor is not byte-aligned
	if bitOffset := b.bitPos % 8; bitOffset != 0 {
		bitsAvail := uint8(8 - bitOffset)
		bitsToRead := bitsAvail
		if bitsToRead > bitsLeft {
			bitsToRead = bitsLeft
		}
		byteVal := b.data[b.bitPos/8]
		shift := bitsAvail - bitsToRead
		mask := uint8((1 << bitsToRead) - 1)
		result = uint64((byteVal >> shift) & mask)
		b.bitPos += uint64(bitsToRead)
		bitsLeft -= bitsToRead
	}

	// Full bytes
	for bitsLeft >= 8 {
		result = (result << 8) | uint64(b.data[b.bitPos/8])
		b.bitPos += 8
		bitsLeft -= 8
	}

	// Partial last byte
	if bitsLeft > 0 {
		byteVal := b.data[b.bitPos/8]
		result = (result << bitsLeft) | uint64(byteVal>>(8-bitsLeft))
		b.bitPos += uint64(bitsLeft)
	}

	return result, nil
}

// WriteBitBuffer accumulates bits in big-endian order into an internal byte slice.
// Bits are placed MSB-first within each byte.
type WriteBitBuffer struct {
	buf         []byte
	pending     byte  // partial byte being assembled from the MSB side
	pendingBits uint8 // number of bits used in pending (0–7)
	tryError    error // accumulated error from TryWrite methods
}

// NewWriteBitBuffer creates a WriteBitBuffer with the given initial capacity hint.
func NewWriteBitBuffer(initialCap int) *WriteBitBuffer {
	return &WriteBitBuffer{buf: make([]byte, 0, initialCap)}
}

// WriteBool writes a single bit.
func (w *WriteBitBuffer) WriteBool(v bool) error {
	var bit uint64
	if v {
		bit = 1
	}
	return w.WriteBits(bit, 1)
}

// WriteBits writes n bits (0 ≤ n ≤ 64) from the LSBs of val, MSB first.
func (w *WriteBitBuffer) WriteBits(val uint64, n uint8) error {
	if n == 0 {
		return nil
	}
	if n > 64 {
		return errors.New("cannot write more than 64 bits at once")
	}

	bitsLeft := n

	// Fill the in-progress pending byte first
	if w.pendingBits > 0 {
		bitsToWrite := uint8(8 - w.pendingBits)
		if bitsToWrite > bitsLeft {
			bitsToWrite = bitsLeft
		}
		shift := bitsLeft - bitsToWrite
		bits := uint8((val >> shift) & ((1 << bitsToWrite) - 1))
		w.pending |= bits << (8 - w.pendingBits - bitsToWrite)
		w.pendingBits += bitsToWrite
		bitsLeft -= bitsToWrite
		if w.pendingBits == 8 {
			w.buf = append(w.buf, w.pending)
			w.pending = 0
			w.pendingBits = 0
		}
	}

	// Write full bytes directly
	for bitsLeft >= 8 {
		bitsLeft -= 8
		w.buf = append(w.buf, byte(val>>bitsLeft))
	}

	// Store remaining bits in pending
	if bitsLeft > 0 {
		w.pending = uint8(val&((1<<bitsLeft)-1)) << (8 - bitsLeft)
		w.pendingBits = bitsLeft
	}

	return nil
}

// TryWriteByte writes a byte, accumulating any error into GetTryError.
func (w *WriteBitBuffer) TryWriteByte(b byte) {
	if err := w.WriteBits(uint64(b), 8); err != nil && w.tryError == nil {
		w.tryError = stdErrors.Join(w.tryError, err)
	}
}

// GetTryError returns the accumulated error from TryWriteByte calls.
func (w *WriteBitBuffer) GetTryError() error {
	return w.tryError
}

// flush pads the pending partial byte with zero bits and appends it.
func (w *WriteBitBuffer) flush() {
	if w.pendingBits > 0 {
		w.buf = append(w.buf, w.pending)
		w.pending = 0
		w.pendingBits = 0
	}
}

// Bytes flushes any pending partial byte (zero-padded) and returns all written bytes.
func (w *WriteBitBuffer) Bytes() []byte {
	w.flush()
	return w.buf
}

// ByteLen returns the number of complete bytes written, excluding any pending partial byte.
func (w *WriteBitBuffer) ByteLen() int {
	return len(w.buf)
}
