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
	"testing"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// An optional field is read only when its condition holds, but the read can still come up empty on
// a datagram that ends early: ReadOptionalField reports "not enough bytes" as an absent field.
// Anything deriving a value from such a field therefore has to require it to be present, or the
// derivation dereferences a nil pointer. Each case below announces a field and then stops before it.
func TestTruncatedOptionalIsRejectedNotDereferenced(t *testing.T) {
	ctx := context.Background()
	buf := func(b ...byte) utils.ReadBufferByteBased { return utils.NewReadBufferByteBased(b) }

	tests := []struct {
		name  string
		parse func() (any, error)
	}{
		// An actualLength of 1 announces a one byte payload that is not there.
		{"unsigned integer payload", func() (any, error) {
			return BACnetTagPayloadUnsignedIntegerParseWithBuffer(ctx, buf(), uint32(1))
		}},
		{"signed integer payload", func() (any, error) {
			return BACnetTagPayloadSignedIntegerParseWithBuffer(ctx, buf(), uint32(1))
		}},
		// A segmented APDU announces a sequence number and a proposed window size after the
		// invoke id; both are missing here.
		{"segmented complex ack", func() (any, error) {
			return APDUParseWithBuffer[APDU](ctx, buf(0x38, 0x01), uint16(2))
		}},
		{"segmented confirmed request", func() (any, error) {
			return APDUParseWithBuffer[APDU](ctx, buf(0x08, 0x00, 0x01), uint16(3))
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("parsing panicked instead of returning an error: %v", r)
				}
			}()
			if _, err := test.parse(); err == nil {
				t.Error("expected a truncated message to be rejected")
			}
		})
	}
}
