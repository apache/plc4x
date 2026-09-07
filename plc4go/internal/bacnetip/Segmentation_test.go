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

package bacnetip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

// makeSegment returns an APDUComplexAck carrying the given bytes as a single
// segment for the reassembler.
func makeSegment(invokeId, seq uint8, more bool, payload []byte) model.APDUComplexAck {
	return model.NewAPDUComplexAck(true, more, invokeId, &seq, nil, nil, nil, payload)
}

// ── inboundReassembler ─────────────────────────────────────────────────────

func TestInboundReassembler_ThreeSegments(t *testing.T) {
	r := NewInboundReassembler(7, 1)
	expected := []byte("hello world this is a long payload that needs three chunks!")
	chunks := [][]byte{expected[:20], expected[20:40], expected[40:]}

	for i, chunk := range chunks {
		ack, err := r.AcceptSegment(makeSegment(7, uint8(i), i != len(chunks)-1, chunk))
		require.NoError(t, err)
		assert.False(t, ack.GetNegativeAck(), "segment %d: positive ack expected", i)
		assert.Equal(t, uint8(i), ack.GetSequenceNumber())
	}
	require.True(t, r.Complete())
	assert.Equal(t, expected, r.Bytes())
}

func TestInboundReassembler_RejectsOutOfOrder(t *testing.T) {
	r := NewInboundReassembler(3, 1)
	// First send seq 1 instead of seq 0 — should fail with a NAK for seq 0.
	ack, err := r.AcceptSegment(makeSegment(3, 1, true, []byte("x")))
	require.Error(t, err)
	require.NotNil(t, ack)
	assert.True(t, ack.GetNegativeAck())
	assert.Equal(t, uint8(0), ack.GetSequenceNumber(), "NAK should request the missing seq")
}

func TestInboundReassembler_RejectsInvokeIdMismatch(t *testing.T) {
	r := NewInboundReassembler(7, 1)
	_, err := r.AcceptSegment(makeSegment(8, 0, false, []byte("hi")))
	require.Error(t, err)
}

func TestInboundReassembler_RejectsNonSegmented(t *testing.T) {
	r := NewInboundReassembler(7, 1)
	// segmentedMessage=false
	seq := uint8(0)
	apdu := model.NewAPDUComplexAck(false, false, 7, &seq, nil, nil, nil, []byte("hi"))
	_, err := r.AcceptSegment(apdu)
	require.Error(t, err)
}

func TestInboundReassembler_SinglePositiveAckAfterFinalSegment(t *testing.T) {
	r := NewInboundReassembler(11, 1)
	ack, err := r.AcceptSegment(makeSegment(11, 0, false, []byte("done")))
	require.NoError(t, err)
	assert.False(t, ack.GetNegativeAck())
	assert.True(t, r.Complete())
}

// ── outboundSegmenter ─────────────────────────────────────────────────────

func TestOutboundSegmenter_SplitsPayload(t *testing.T) {
	payload := make([]byte, 500)
	for i := range payload {
		payload[i] = byte(i)
	}
	// maxApdu=128 → 128-20 (headerOverhead) = 108 bytes/segment → ceil(500/108) = 5 segments.
	s := NewOutboundSegmenter(9, payload, 128, 1)
	var collected []byte
	segments := 0
	for s.HasMore() {
		_, seg, more := s.NextSegment()
		collected = append(collected, seg...)
		segments++
		assert.Equal(t, more, s.HasMore())
	}
	assert.Equal(t, payload, collected)
	assert.GreaterOrEqual(t, segments, 4, "expected at least 4 segments for 500B payload at ~108B each")
	assert.LessOrEqual(t, segments, 6)
}

func TestOutboundSegmenter_SingleSegmentWhenFits(t *testing.T) {
	payload := []byte("short")
	s := NewOutboundSegmenter(1, payload, 1476, 1)
	_, seg, more := s.NextSegment()
	assert.Equal(t, payload, seg)
	assert.False(t, more, "single-segment payload should not have moreFollows=true")
	assert.False(t, s.HasMore())
}

func TestOutboundSegmenter_RewindResetsCursor(t *testing.T) {
	payload := make([]byte, 300)
	for i := range payload {
		payload[i] = byte(i)
	}
	s := NewOutboundSegmenter(2, payload, 128, 1)
	_, _, _ = s.NextSegment() // seq 0
	_, _, _ = s.NextSegment() // seq 1
	s.Rewind(1)
	seq, seg, _ := s.NextSegment()
	assert.Equal(t, uint8(1), seq)
	// Same bytes as we got the first time we asked for seq 1.
	assert.Len(t, seg, 108)
}

func TestOutboundSegmenter_AcknowledgeNAK(t *testing.T) {
	s := NewOutboundSegmenter(5, make([]byte, 50), 128, 1)
	ack := model.NewAPDUSegmentAck(true, false, 5, 0, 1) // negativeAck=true
	nak, err := s.AcknowledgeSegment(ack)
	require.NoError(t, err)
	assert.True(t, nak)
}

func TestOutboundSegmenter_AcknowledgeWrongInvokeIdErrors(t *testing.T) {
	s := NewOutboundSegmenter(5, make([]byte, 50), 128, 1)
	ack := model.NewAPDUSegmentAck(false, false, 6, 0, 1) // wrong invoke id
	_, err := s.AcknowledgeSegment(ack)
	require.Error(t, err)
}
