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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// knxFloatSamples are 2 byte KNX floats (DPT 9) with the value they encode: sign in bit 15,
// exponent in bits 14-11, a 12 bit two's complement mantissa of sign and bits 10-0, value =
// 0.01 * mantissa * 2^exponent.
var knxFloatSamples = []struct {
	raw   uint16
	value float32
}{
	{0x0000, 0},
	{0x0001, 0.01},
	{0x0064, 1},
	{0x0C33, 21.5}, // 1075 * 2^1
	{0x0C38, 21.6}, // 1080 * 2^1
	{0x8BCD, -21.5},
	{0x87FF, -0.01},
	{0x5DB8, 29982.72}, // 1464 * 2^11
}

func TestReadFloat32_KnxFloat(t *testing.T) {
	samples := append(knxFloatSamples, struct {
		raw   uint16
		value float32
	}{0x7FFF, 670760.96}, // the largest value: 2047 * 2^15
		struct {
			raw   uint16
			value float32
		}{0xF800, -671088.64}) // the smallest value: -2048 * 2^15
	for _, sample := range samples {
		rb := NewReadBufferByteBased([]byte{byte(sample.raw >> 8), byte(sample.raw)})
		value, err := rb.ReadFloat32("value", 16, WithEncoding("KNXFloat"))
		require.NoError(t, err, "0x%04X", sample.raw)
		assert.InDelta(t, sample.value, value, 0.001, "0x%04X", sample.raw)
		assert.False(t, rb.HasMore(1))
	}
}

func TestWriteFloat32_KnxFloat(t *testing.T) {
	for _, sample := range knxFloatSamples {
		wb := NewWriteBufferByteBased()
		require.NoError(t, wb.WriteFloat32("value", 16, sample.value, WithEncoding("KNXFloat")), "%v", sample.value)
		assert.Equal(t, []byte{byte(sample.raw >> 8), byte(sample.raw)}, wb.GetBytes(), "%v", sample.value)
	}
	// The mantissa is truncated to 11 bits, as the java implementation does: 30000 is 3000000
	// hundredths, which only fit after 11 halvings.
	wb := NewWriteBufferByteBased()
	require.NoError(t, wb.WriteFloat32("value", 16, 30000, WithEncoding("KNXFloat")))
	assert.Equal(t, []byte{0x5D, 0xB8}, wb.GetBytes())

	for _, value := range []float32{670761, -670761} {
		assert.Error(t, NewWriteBufferByteBased().WriteFloat32("value", 16, value, WithEncoding("KNXFloat")), "%v", value)
	}
	assert.Error(t, NewWriteBufferByteBased().WriteFloat32("value", 32, 1, WithEncoding("KNXFloat")))
	_, err := NewReadBufferByteBased([]byte{0, 0, 0, 0}).ReadFloat32("value", 32, WithEncoding("KNXFloat"))
	assert.Error(t, err)
}
