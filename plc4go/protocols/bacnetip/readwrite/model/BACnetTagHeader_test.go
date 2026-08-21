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
	"fmt"
	"testing"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// A tag header whose length-value-type is 5 only carries an extended length when
// the tag is primitive and not boolean. For a boolean or constructed tag the
// extended length fields are absent, so actualLength must fall back to the
// length-value-type instead of reading them.
func TestBACnetTagHeaderActualLengthWithoutExtendedLength(t *testing.T) {
	// 0x15 -> tagNumber 1, APPLICATION_TAGS, lengthValueType 5: a boolean tag,
	// so no extLength follows on the wire.
	header, err := BACnetTagHeaderParseWithBuffer(context.Background(), utils.NewReadBufferByteBased([]byte{0x15}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if header.GetExtLength() != nil {
		t.Fatalf("expected no extLength, got %d", *header.GetExtLength())
	}
	if actual := header.GetActualLength(); actual != 5 {
		t.Errorf("expected actualLength 5, got %d", actual)
	}
}

// Every possible tag header must either parse or return an error, never panic.
func TestBACnetTagHeaderParseAllFirstBytesWithoutPanic(t *testing.T) {
	// Enough trailing bytes for extTagNumber plus the widest extended length.
	trailer := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	for first := 0; first <= 0xFF; first++ {
		for trailerLen := 0; trailerLen <= len(trailer); trailerLen++ {
			data := append([]byte{byte(first)}, trailer[:trailerLen]...)
			t.Run(fmt.Sprintf("0x%02X/%d", first, trailerLen), func(t *testing.T) {
				header, err := BACnetTagHeaderParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(data))
				if err == nil {
					// Touching the virtual fields must be safe too.
					_ = header.GetActualLength()
					_ = header.GetActualTagNumber()
				}
			})
		}
	}
}
