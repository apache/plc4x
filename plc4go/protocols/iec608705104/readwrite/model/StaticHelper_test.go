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

package model

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TestFinished pins the termination predicate of the terminated "apdus" array.
// It mirrors plc4j's StaticHelper.finished(ReadBuffer), which stops as soon as
// fewer than eight bits are left, i.e. the moment no further full byte can be
// read.
func TestFinished(t *testing.T) {
	readBuffer := utils.NewReadBufferByteBased([]byte{0x01, 0x02})
	finished := Finished(context.Background(), readBuffer)

	assert.False(t, finished(), "two bytes left")
	_, err := readBuffer.ReadByte("")
	require.NoError(t, err)
	assert.False(t, finished(), "one byte left")
	_, err = readBuffer.ReadByte("")
	require.NoError(t, err)
	assert.True(t, finished(), "buffer drained")
}

// TestFinishedStopsShortOfAPartialByte pins that a trailing fragment of less
// than a byte terminates the array too, instead of provoking a parse of a stub
// APDU that cannot possibly be complete.
func TestFinishedStopsShortOfAPartialByte(t *testing.T) {
	readBuffer := utils.NewReadBufferByteBased([]byte{0xFF})
	finished := Finished(context.Background(), readBuffer)

	require.False(t, finished())
	_, err := readBuffer.ReadUint8("", 1)
	require.NoError(t, err)
	assert.True(t, finished(), "seven bits are not a byte")
}

// TestFinishedDoesNotConsume guards that evaluating the predicate leaves the
// buffer position alone - ReadTerminatedArrayField calls it before every element,
// so a predicate that consumed anything would corrupt the following APDU.
func TestFinishedDoesNotConsume(t *testing.T) {
	readBuffer := utils.NewReadBufferByteBased([]byte{0x01, 0x02})
	finished := Finished(context.Background(), readBuffer)

	for range 3 {
		require.False(t, finished())
	}
	assert.Equal(t, uint32(0), readBuffer.GetPos())
}

// TestAPDUsParseTerminates is the end-to-end check that the predicate actually
// wires up: the raw frames are the ones from the IEC-60870-5-104
// parser-serializer testsuite, and the point of the APDUs wrapper is that
// devices commonly send several of them in a single packet.
func TestAPDUsParseTerminates(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want int
	}{
		// "4 <- U (TESTFR act)" from the testsuite.
		{"single frame", "680443000000", 1},
		// TESTFR act followed by TESTFR con, as a device would coalesce them.
		{"two frames", "680443000000680483000000", 2},
		{"three frames", "680443000000680483000000680407000000", 3},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			raw, err := hex.DecodeString(test.raw)
			require.NoError(t, err)

			apdus, err := APDUsParse(context.Background(), raw)
			require.NoError(t, err)
			assert.Len(t, apdus.GetApdus(), test.want)

			serialized, err := apdus.Serialize()
			require.NoError(t, err)
			assert.Equal(t, raw, serialized, "must round trip")
		})
	}
}
