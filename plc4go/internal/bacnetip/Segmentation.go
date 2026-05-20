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
	"context"

	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// inboundReassembler buffers a multi-segment confirmed-service response from a
// remote BACnet device. The caller feeds segments one at a time via
// AcceptSegment; once `MoreFollows` is false on the final segment, Complete
// returns the concatenated raw service-ack payload bytes ready for re-parse.
//
// The reassembler does not parse — it simply collects bytes and tracks the
// expected sequence number, returning a SegmentAck PDU the caller must transmit
// after every received segment.
type inboundReassembler struct {
	invokeId     uint8
	expectedSeq  uint8
	buffer       []byte
	windowSize   uint8 // actual window size we advertised in the original request
	moreFollows  bool
	complete     bool
}

// NewInboundReassembler creates a fresh reassembler tracking the given invoke
// id. windowSize is the actual window size we want to advertise back to the
// publisher in each SegmentAck (often capped to 1 for simplicity until a real
// device shows up that benefits from larger windows).
func NewInboundReassembler(invokeId, windowSize uint8) *inboundReassembler {
	if windowSize == 0 {
		windowSize = 1
	}
	return &inboundReassembler{
		invokeId:   invokeId,
		windowSize: windowSize,
	}
}

// AcceptSegment consumes one APDUComplexAck segment. Returns the SegmentAck the
// caller must send back, and an error if the segment's sequence number doesn't
// match what we expect.
//
// When the returned ack has Complete()=true the reassembler is finished and
// Bytes() yields the concatenated payload.
func (r *inboundReassembler) AcceptSegment(seg model.APDUComplexAck) (model.APDUSegmentAck, error) {
	if !seg.GetSegmentedMessage() {
		return nil, errors.New("non-segmented APDU passed to reassembler")
	}
	if seg.GetOriginalInvokeId() != r.invokeId {
		return nil, errors.Errorf("segment invoke-id mismatch: got %d, expected %d", seg.GetOriginalInvokeId(), r.invokeId)
	}
	seqPtr := seg.GetSequenceNumber()
	if seqPtr == nil {
		return nil, errors.New("segmented APDU is missing sequence number")
	}
	seq := *seqPtr
	if seq != r.expectedSeq {
		// Out-of-order segments are a protocol error; signal a negative ack so
		// the publisher restarts from the missing segment.
		return model.NewAPDUSegmentAck(true, true, r.invokeId, r.expectedSeq, r.windowSize), errors.Errorf("out-of-order segment: got seq %d, expected %d", seq, r.expectedSeq)
	}
	r.buffer = append(r.buffer, seg.GetSegment()...)
	r.expectedSeq = seq + 1
	r.moreFollows = seg.GetMoreFollows()
	r.complete = !r.moreFollows
	// server=true: we are the server (i.e. response receiver) acknowledging.
	return model.NewAPDUSegmentAck(false, true, r.invokeId, seq, r.windowSize), nil
}

// Complete reports whether the final segment has been received.
func (r *inboundReassembler) Complete() bool { return r.complete }

// Bytes returns the concatenated payload bytes accumulated across all
// accepted segments. Only meaningful once Complete() is true.
func (r *inboundReassembler) Bytes() []byte { return r.buffer }

// outboundSegmenter splits a large confirmed-request payload into APDU segments
// honoring the peer's MaxApduLengthAccepted. The caller drives the loop:
// SendNext returns the next segment to put on the wire, then waits for the
// remote SegmentAck via AcknowledgeSegment before calling SendNext again.
type outboundSegmenter struct {
	invokeId     uint8
	payload      []byte
	maxSegment   uint16 // bytes per segment (peer's max APDU minus header overhead)
	windowSize   uint8
	nextSeq      uint8
	cursor       int  // byte offset into payload
	done         bool
}

// NewOutboundSegmenter prepares a segmenter for the given serialized request
// payload. maxApdu is the peer's MaxApduLengthAccepted in bytes; we subtract a
// conservative 20-byte overhead allowance for APDU/NPDU/BVLC headers.
func NewOutboundSegmenter(invokeId uint8, payload []byte, maxApdu uint16, windowSize uint8) *outboundSegmenter {
	if windowSize == 0 {
		windowSize = 1
	}
	const headerOverhead = 20
	segmentBytes := uint16(0)
	if maxApdu > headerOverhead {
		segmentBytes = maxApdu - headerOverhead
	}
	if segmentBytes == 0 {
		segmentBytes = 1 // safety net; a maxApdu < overhead means we can't really segment, but try.
	}
	return &outboundSegmenter{
		invokeId:   invokeId,
		payload:    payload,
		maxSegment: segmentBytes,
		windowSize: windowSize,
	}
}

// HasMore reports whether more segments remain to be sent.
func (s *outboundSegmenter) HasMore() bool { return !s.done }

// NextSegment returns the next chunk of payload bytes along with a flag
// indicating whether more segments follow. Advances the internal cursor.
func (s *outboundSegmenter) NextSegment() (seq uint8, segment []byte, moreFollows bool) {
	if s.done {
		return 0, nil, false
	}
	end := s.cursor + int(s.maxSegment)
	if end >= len(s.payload) {
		end = len(s.payload)
		s.done = true
	}
	segment = s.payload[s.cursor:end]
	seq = s.nextSeq
	moreFollows = !s.done
	s.cursor = end
	s.nextSeq++
	return seq, segment, moreFollows
}

// AcknowledgeSegment consumes a SegmentAck from the remote. A negative ack
// (NAK) is reported back to the caller via the boolean; the caller may then
// rewind to the requested sequence by calling Rewind.
func (s *outboundSegmenter) AcknowledgeSegment(ack model.APDUSegmentAck) (negativeAck bool, err error) {
	if ack.GetOriginalInvokeId() != s.invokeId {
		return false, errors.Errorf("segment-ack invoke-id mismatch: got %d, expected %d", ack.GetOriginalInvokeId(), s.invokeId)
	}
	return ack.GetNegativeAck(), nil
}

// Rewind resets the cursor to the start of the segment numbered seq so the
// next NextSegment call resends from that point. Used in response to a NAK.
func (s *outboundSegmenter) Rewind(seq uint8) {
	s.cursor = int(seq) * int(s.maxSegment)
	if s.cursor > len(s.payload) {
		s.cursor = len(s.payload)
	}
	s.nextSeq = seq
	s.done = false
}

// AwaitableSegmenter is a thin context-aware wrapper for callers that want to
// honor ctx cancellation while waiting for segment acks. Reserved for future
// use; not consumed by the current Reader/Writer flow.
func AwaitableSegmenter(_ context.Context) {}
