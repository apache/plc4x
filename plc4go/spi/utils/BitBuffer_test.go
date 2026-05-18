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
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── ReadBitBuffer ──────────────────────────────────────────────────────────

func TestReadBitBuffer_BitsRemaining(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xFF, 0x00})
	assert.Equal(t, uint64(16), b.BitsRemaining())

	_, _ = b.ReadBits(4)
	assert.Equal(t, uint64(12), b.BitsRemaining())

	_, _ = b.ReadBits(12)
	assert.Equal(t, uint64(0), b.BitsRemaining())
}

func TestReadBitBuffer_ReadBits_ByteAligned(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xAB, 0xCD})

	v, err := b.ReadBits(8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xAB), v)

	v, err = b.ReadBits(8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xCD), v)
}

func TestReadBitBuffer_ReadBits_CrossByte(t *testing.T) {
	// 0xAB = 1010_1011, 0xCD = 1100_1101
	// Reading 4 bits then 8 bits: first 4 = 0xA (1010), next 8 cross the byte boundary
	// bits 4-11: 1011_1100 = 0xBC
	b := NewReadBitBuffer([]byte{0xAB, 0xCD})

	v, err := b.ReadBits(4)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xA), v)

	v, err = b.ReadBits(8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xBC), v)
}

func TestReadBitBuffer_ReadBits_Partial(t *testing.T) {
	// 0b10110000 — read top 3 bits = 0b101 = 5
	b := NewReadBitBuffer([]byte{0b10110000})

	v, err := b.ReadBits(3)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), v)
}

func TestReadBitBuffer_ReadBool(t *testing.T) {
	// 0b10000000: first bit = 1 (true), second bit = 0 (false)
	b := NewReadBitBuffer([]byte{0b10000000})

	v, err := b.ReadBool()
	require.NoError(t, err)
	assert.True(t, v)

	v, err = b.ReadBool()
	require.NoError(t, err)
	assert.False(t, v)
}

func TestReadBitBuffer_ReadByte(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xDE, 0xAD})

	byt, err := b.ReadByte()
	require.NoError(t, err)
	assert.Equal(t, byte(0xDE), byt)
}

func TestReadBitBuffer_Read(t *testing.T) {
	b := NewReadBitBuffer([]byte{0x01, 0x02, 0x03})
	p := make([]byte, 3)
	n, err := b.Read(p)
	require.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, p)
}

func TestReadBitBuffer_ReadBits_EOF(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xFF})
	_, err := b.ReadBits(9)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestReadBitBuffer_ReadBits_Zero(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xFF})
	v, err := b.ReadBits(0)
	require.NoError(t, err)
	assert.Equal(t, uint64(0), v)
}

func TestReadBitBuffer_ReadBits_TooMany(t *testing.T) {
	b := NewReadBitBuffer(make([]byte, 9))
	_, err := b.ReadBits(65)
	require.Error(t, err)
	assert.Equal(t, uint64(72), b.BitsRemaining(), "cursor must not advance on error")
}

func TestReadBitBuffer_ReadBits_64(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	b := NewReadBitBuffer(data)
	v, err := b.ReadBits(64)
	require.NoError(t, err)
	assert.Equal(t, uint64(0x0102030405060708), v)
	assert.Equal(t, uint64(0), b.BitsRemaining())
}

func TestReadBitBuffer_EmptyBuffer(t *testing.T) {
	b := NewReadBitBuffer([]byte{})
	assert.Equal(t, uint64(0), b.BitsRemaining())
	_, err := b.ReadBits(1)
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestReadBitBuffer_Read_Partial(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xAA})
	p := make([]byte, 3)
	n, err := b.Read(p)
	assert.Equal(t, 1, n)
	assert.Error(t, err)
	assert.Equal(t, byte(0xAA), p[0])
}

func TestReadBitBuffer_ResetTo(t *testing.T) {
	b := NewReadBitBuffer([]byte{0xAB, 0xCD})

	v, err := b.ReadBits(8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xAB), v)

	b.ResetTo(0)
	v, err = b.ReadBits(8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xAB), v)
}

func TestReadBitBuffer_ResetTo_NonZero(t *testing.T) {
	// Reset to bit position 4 (mid-byte) and read the lower nibble of 0xAB.
	// 0xAB = 1010_1011 → lower 4 bits = 0xB
	b := NewReadBitBuffer([]byte{0xAB, 0xCD})
	b.ResetTo(4)
	v, err := b.ReadBits(4)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xB), v)
}

// ── WriteBitBuffer ─────────────────────────────────────────────────────────

func TestWriteBitBuffer_WriteByteAligned(t *testing.T) {
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xAB, 8))
	require.NoError(t, w.WriteBits(0xCD, 8))
	assert.Equal(t, []byte{0xAB, 0xCD}, w.Bytes())
}

func TestWriteBitBuffer_WriteBool(t *testing.T) {
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBool(true))
	require.NoError(t, w.WriteBool(false))
	require.NoError(t, w.WriteBool(true))
	require.NoError(t, w.WriteBool(true))
	require.NoError(t, w.WriteBool(false))
	require.NoError(t, w.WriteBool(false))
	require.NoError(t, w.WriteBool(false))
	require.NoError(t, w.WriteBool(false))
	// 1011_0000 = 0xB0
	assert.Equal(t, []byte{0xB0}, w.Bytes())
}

func TestWriteBitBuffer_WritePartialThenAlign(t *testing.T) {
	// Write 4 bits of 0xA (1010), then 4 bits of 0xB (1011) → 0xAB
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xA, 4))
	require.NoError(t, w.WriteBits(0xB, 4))
	assert.Equal(t, []byte{0xAB}, w.Bytes())
}

func TestWriteBitBuffer_WriteCrossByte(t *testing.T) {
	// Write 4 bits 0xA, then 8 bits 0xBC, then 4 bits 0xD
	// Result bytes: 0xA<<4 | 0xB = 0xAB, 0xC<<4 | 0xD = 0xCD
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xA, 4))
	require.NoError(t, w.WriteBits(0xBC, 8))
	require.NoError(t, w.WriteBits(0xD, 4))
	assert.Equal(t, []byte{0xAB, 0xCD}, w.Bytes())
}

func TestWriteBitBuffer_PendingFlushedWithZeroPad(t *testing.T) {
	// Write 3 bits: 0b101 → should be padded to 0b10100000 = 0xA0
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0b101, 3))
	assert.Equal(t, []byte{0xA0}, w.Bytes())
}

func TestWriteBitBuffer_ByteLen(t *testing.T) {
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xAB, 8))
	require.NoError(t, w.WriteBits(0b101, 3)) // pending, not yet flushed
	assert.Equal(t, 1, w.ByteLen())           // only the complete byte
}

func TestWriteBitBuffer_TryWriteByte(t *testing.T) {
	w := NewWriteBitBuffer(0)
	w.TryWriteByte(0xDE)
	w.TryWriteByte(0xAD)
	assert.NoError(t, w.GetTryError())
	assert.Equal(t, []byte{0xDE, 0xAD}, w.Bytes())
}

func TestWriteBitBuffer_WriteBits_Zero(t *testing.T) {
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xFF, 0))
	assert.Equal(t, []byte{}, w.Bytes())
}

func TestWriteBitBuffer_WriteBits_TooMany(t *testing.T) {
	w := NewWriteBitBuffer(0)
	err := w.WriteBits(0, 65)
	require.Error(t, err)
	assert.Equal(t, []byte{}, w.Bytes(), "buffer must not change on error")
}

func TestWriteBitBuffer_WriteBits_64(t *testing.T) {
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0x0102030405060708, 64))
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, w.Bytes())
}

func TestWriteBitBuffer_SmallWriteIntoPartialPending(t *testing.T) {
	// Write 3 bits (1=101), then 2 bits (1=11) — neither fills the pending byte.
	// After both: pending = 101_11_000 → but pending only holds one byte (8 bits),
	// so 5 bits used: 1011_1000 wait — it's MSB first:
	// bit 7: 1, bit 6: 0, bit 5: 1 → after first write, pending=1010_0000
	// then write 11: bit 4: 1, bit 3: 1 → pending=1011_1000 = 0xB8
	// Then flush via Bytes()
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0b101, 3))
	require.NoError(t, w.WriteBits(0b11, 2))
	assert.Equal(t, 0, w.ByteLen(), "no complete bytes yet")
	assert.Equal(t, []byte{0b10111000}, w.Bytes()) // 101_11_000 zero-padded = 0xB8
}

func TestWriteBitBuffer_HighBitsInValIgnored(t *testing.T) {
	// Only the lower n bits of val should be written.
	// WriteBits(0xFF, 4) should write 0b1111, not the upper 4 bits.
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(0xFF, 4))
	require.NoError(t, w.WriteBits(0x00, 4))
	assert.Equal(t, []byte{0xF0}, w.Bytes())
}

// ── Round-trip ─────────────────────────────────────────────────────────────

func TestBitBuffer_RoundTrip(t *testing.T) {
	// Write a sequence of non-byte-aligned fields and read them back.
	w := NewWriteBitBuffer(0)
	require.NoError(t, w.WriteBits(1, 1))       // 1 bit
	require.NoError(t, w.WriteBits(0xA, 4))     // 4 bits
	require.NoError(t, w.WriteBits(0x1F, 5))    // 5 bits
	require.NoError(t, w.WriteBits(0xABCD, 16)) // 16 bits (2 full bytes)

	r := NewReadBitBuffer(w.Bytes())

	v, err := r.ReadBits(1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), v)

	v, err = r.ReadBits(4)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xA), v)

	v, err = r.ReadBits(5)
	require.NoError(t, err)
	assert.Equal(t, uint64(0x1F), v)

	v, err = r.ReadBits(16)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xABCD), v)
}
