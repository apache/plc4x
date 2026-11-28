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

package virtualcan

import (
	"bytes"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/spi/transports/can"
)

func connect(t *testing.T, busName string) *TransportInstance {
	t.Helper()
	instance, err := NewTransport().CreateTransportInstance(url.URL{Scheme: TransportCode, Host: busName}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := instance.Connect(t.Context()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = instance.Close() })
	return instance.(*TransportInstance)
}

func record(t *testing.T, frame can.Frame) []byte {
	t.Helper()
	encoded, err := frame.Encode()
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func TestBusDeliversToOthersOnly(t *testing.T) {
	a, b, c := connect(t, t.Name()), connect(t, t.Name()), connect(t, t.Name())
	other := connect(t, t.Name()+"-other")
	frame := record(t, can.Frame{ID: 0x181, Data: []byte{1, 2}})
	if err := a.Write(t.Context(), frame); err != nil {
		t.Fatal(err)
	}
	for _, receiver := range []*TransportInstance{b, c} {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		got, err := receiver.Read(ctx, can.FrameSize)
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, frame) {
			t.Fatalf("got % X, want % X", got, frame)
		}
	}
	for _, silent := range []*TransportInstance{a, other} {
		if n, _ := silent.GetNumBytesAvailableInBuffer(); n != 0 {
			t.Fatalf("%s received %d bytes it should not have", silent, n)
		}
	}
}

func TestReadTimesOutAndCloseWakes(t *testing.T) {
	a := connect(t, t.Name())
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Millisecond)
	defer cancel()
	if _, err := a.Read(ctx, can.FrameSize); err == nil {
		t.Fatal("read on an idle bus returned")
	}
	done := make(chan error, 1)
	go func() {
		_, err := a.Read(t.Context(), can.FrameSize)
		done <- err
	}()
	time.Sleep(10 * time.Millisecond)
	_ = a.Close()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("read on a closed instance succeeded")
		}
	case <-time.After(time.Second):
		t.Fatal("Close did not wake a pending read")
	}
}

func TestOverrunDropsWholeFrames(t *testing.T) {
	a, b := connect(t, t.Name()), connect(t, t.Name())
	for i := range maxBufferedFrames + 10 {
		if err := a.Write(t.Context(), record(t, can.Frame{ID: uint32(i % 0x7FF), Data: []byte{byte(i)}})); err != nil {
			t.Fatal(err)
		}
	}
	if b.DroppedFrames() != 10 {
		t.Fatalf("dropped %d frames, want 10", b.DroppedFrames())
	}
	first, err := b.Read(t.Context(), can.FrameSize)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := can.Decode(first)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Data[0] != 10 {
		t.Fatalf("oldest kept frame is #%d, want #10", decoded.Data[0])
	}
}

func TestWriteRejectsMalformedRecords(t *testing.T) {
	a := connect(t, t.Name())
	if err := a.Write(t.Context(), make([]byte, 10)); err == nil {
		t.Fatal("a partial record was accepted")
	}
	bad := make([]byte, can.FrameSize)
	bad[4] = 12
	if err := a.Write(t.Context(), bad); err == nil {
		t.Fatal("a record with a data length code of 12 was accepted")
	}
}
