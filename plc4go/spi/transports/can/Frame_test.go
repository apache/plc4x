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

package can

import (
	"bytes"
	"testing"
)

func TestFrameRoundTrip(t *testing.T) {
	for _, frame := range []Frame{
		{ID: 0x601, Data: []byte{0x40, 0x18, 0x10, 0x01, 0, 0, 0, 0}},
		{ID: 0x000, Data: []byte{0x80, 0x00}},
		{ID: 0x701, Data: []byte{0x05}},
		{ID: 0x80},
		{ID: 0x1ABCDEF0, Extended: true, Data: []byte{1, 2, 3}},
		{ID: 0x123, Remote: true},
	} {
		record, err := frame.Encode()
		if err != nil {
			t.Fatalf("%s: %v", frame, err)
		}
		if len(record) != FrameSize {
			t.Fatalf("%s: record is %d bytes", frame, len(record))
		}
		decoded, err := Decode(record)
		if err != nil {
			t.Fatalf("%s: %v", frame, err)
		}
		if decoded.ID != frame.ID || decoded.Extended != frame.Extended || decoded.Remote != frame.Remote || !bytes.Equal(decoded.Data, frame.Data) {
			t.Fatalf("round trip changed %s into %s", frame, decoded)
		}
	}
}

func TestFrameLayout(t *testing.T) {
	record, err := Frame{ID: 0x581, Data: []byte{0x43, 0x00, 0x10, 0x00, 0x91, 0x01, 0x0F, 0x00}}.Encode()
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00, 0x00, 0x05, 0x81, 8, 0, 0, 0, 0x43, 0x00, 0x10, 0x00, 0x91, 0x01, 0x0F, 0x00}
	if !bytes.Equal(record, want) {
		t.Fatalf("layout\n got % X\nwant % X", record, want)
	}
}

func TestFrameRejects(t *testing.T) {
	if _, err := (Frame{ID: 0x800}).Encode(); err == nil {
		t.Fatal("an 12 bit standard identifier was accepted")
	}
	if _, err := (Frame{ID: 1, Data: make([]byte, 9)}).Encode(); err == nil {
		t.Fatal("9 data bytes were accepted")
	}
	record := make([]byte, FrameSize)
	record[4] = 9
	if _, err := Decode(record); err == nil {
		t.Fatal("a data length code of 9 was accepted")
	}
	if _, err := Decode(record[:8]); err == nil {
		t.Fatal("a short record was accepted")
	}
}
