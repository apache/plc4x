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

package socketcan

import (
	"context"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/spi/transports/can"
)

// TestErrorFrames checks socketcan.error-mask on a real or virtual CAN interface named by
// SOCKETCAN_TEST_INTERFACE (for example vcan0): an error frame reaches a socket whose mask has its
// class, and not one whose mask is 0. A virtual interface generates no errors of its own, so the
// test sends one; the kernel filters error frames on receive by the same mask either way.
func TestErrorFrames(t *testing.T) {
	interfaceName := os.Getenv("SOCKETCAN_TEST_INTERFACE")
	if interfaceName == "" {
		t.Skip("SOCKETCAN_TEST_INTERFACE is not set")
	}
	const busOff = 0x40 // CAN_ERR_BUSOFF
	open := func(mask []string) *TransportInstance {
		options := map[string][]string{}
		if mask != nil {
			options[ErrorMaskOption] = mask
		}
		instance, err := NewTransport().CreateTransportInstance(url.URL{Scheme: TransportCode, Host: interfaceName}, options)
		if err != nil {
			t.Fatal(err)
		}
		if err := instance.Connect(t.Context()); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = instance.Close() })
		return instance.(*TransportInstance)
	}
	sender := open(nil)
	withoutMask := open(nil)
	busOffOnly := open([]string{"0x40"})
	allClasses := open([]string{"0x1FFFFFFF"})
	otherClass := open([]string{"0x04"}) // CAN_ERR_CRTL only

	errorFrame, err := can.Frame{ID: busOff, Error: true, Data: make([]byte, 8)}.Encode()
	if err != nil {
		t.Fatal(err)
	}
	// A data frame after the error frame marks the end: a socket that sees it first got no error frame.
	marker, err := can.Frame{ID: 0x7E5, Data: []byte{0xA5}}.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if err := sender.Write(t.Context(), append(errorFrame, marker...)); err != nil {
		t.Fatal(err)
	}

	next := func(instance *TransportInstance) can.Frame {
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		record, err := instance.Read(ctx, can.FrameSize)
		if err != nil {
			t.Fatalf("%s: %v", instance, err)
		}
		frame, err := can.Decode(record)
		if err != nil {
			t.Fatal(err)
		}
		return frame
	}
	for _, receiver := range []struct {
		name      string
		instance  *TransportInstance
		wantError bool
	}{
		{"no mask", withoutMask, false},
		{"mask 0x40", busOffOnly, true},
		{"mask 0x1FFFFFFF", allClasses, true},
		{"mask 0x04", otherClass, false},
	} {
		first := next(receiver.instance)
		if receiver.wantError {
			if !first.Error || first.ID != busOff {
				t.Errorf("%s: first frame %s, want the bus-off error frame", receiver.name, first)
				continue
			}
			first = next(receiver.instance)
		}
		if first.Error || first.ID != 0x7E5 {
			t.Errorf("%s: got %s, want the marker frame", receiver.name, first)
			continue
		}
		t.Logf("%s: error frame delivered: %v", receiver.name, receiver.wantError)
	}
}
