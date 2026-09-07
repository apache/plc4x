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
	"time"

	"github.com/rs/zerolog"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// segmentAckWaitTimeout bounds how long we wait for the peer's SegmentAck after
// transmitting a window of request segments (ASHRAE 135 T-seg, default 5s).
const segmentAckWaitTimeout = 5 * time.Second

// proposedSegmentWindowSize is the window size we propose in every request
// segment. The peer's SegmentAck answers with the ACTUAL window size, which
// then governs how many segments we put on the wire per ack.
const proposedSegmentWindowSize = 16

// maxOutboundRequestSegments caps a segmented request at the uint8 sequence
// space so we never have to disambiguate wrapped sequence numbers. A request
// this large (~250 * peer-max-APDU bytes) exceeds any real device's
// max-segments-accepted anyway.
const maxOutboundRequestSegments = 255

// unsegmentedConfirmedRequestHeaderBytes is the fixed APDU header of an
// UNsegmented confirmed request: PDU-type/flags octet, max-segments/max-APDU
// octet, and the invoke id. Everything after it is the serialized service
// request, so `header + len(serviceRequest)` is the wire APDU size used to
// decide whether a request fits the peer's MaxApduLengthAccepted.
const unsegmentedConfirmedRequestHeaderBytes = 3

// segmentedRequestSender transmits an oversized confirmed request as a
// sequence of APDUConfirmedRequest segments per ASHRAE 135 clause 5.4
// (segment 0 alone, then windows of segments each answered by a SegmentAck,
// with NAK-driven rewind and bounded retransmission on ack timeout).
//
// The caller must register its expectation for the FINAL service response
// (excluding APDUSegmentAck — see responseMatcherExcludingSegmentAcks) BEFORE
// calling send: the peer may answer immediately after acking the last segment.
type segmentedRequestSender struct {
	messageCodec  spi.MessageCodec
	routedDest    *routedDestination
	driverContext DriverContext
	log           zerolog.Logger
}

// needsSegmentedRequest reports whether the serialized service-request payload
// exceeds the peer's declared APDU ceiling. With an unknown ceiling (0) the
// answer is always false — we then send unsegmented exactly as before.
func (d DriverContext) needsSegmentedRequest(payloadBytes int) bool {
	return d.peerMaxApduBytes > 0 && payloadBytes+unsegmentedConfirmedRequestHeaderBytes > int(d.peerMaxApduBytes)
}

// send drives the segmented transmission of payload (the serialized
// BACnetConfirmedServiceRequest, starting with its service-choice byte) under
// the given invoke id. It returns once the peer has acknowledged the final
// segment; the service response itself is delivered through the caller's
// pre-registered expectation.
func (s *segmentedRequestSender) send(ctx context.Context, invokeId uint8, payload []byte) error {
	if len(payload) == 0 {
		return errors.New("empty confirmed-request payload")
	}
	segmenter := NewOutboundSegmenter(invokeId, payload, s.driverContext.peerMaxApduBytes, proposedSegmentWindowSize)
	total := segmenter.TotalSegments()
	if total > maxOutboundRequestSegments {
		return errors.Errorf("confirmed request of %d bytes needs %d segments, exceeding the %d-segment ceiling", len(payload), total, maxOutboundRequestSegments)
	}
	serviceChoice := readWriteModel.BACnetConfirmedServiceChoice(payload[0])

	s.log.Debug().
		Uint8("invokeId", invokeId).
		Int("payloadBytes", len(payload)).
		Int("segments", total).
		Uint16("peerMaxApdu", s.driverContext.peerMaxApduBytes).
		Msg("sending segmented confirmed request")

	// Window is 1 until the first SegmentAck reveals the peer's actual window
	// size (clause 5.4.4: segment 0 is sent alone).
	window := uint8(1)
	lastAcked := -1
	retries := uint8(0)

	for {
		// Register the ack expectation BEFORE putting segments on the wire,
		// otherwise the peer's ack could race ahead of the expectation.
		ackCh, errCh := s.expectSegmentAck(ctx, invokeId)

		sent := 0
		for sent < int(window) && segmenter.HasMore() {
			seq, chunk, moreFollows := segmenter.NextSegment()
			apdu := s.segmentApdu(invokeId, seq, moreFollows, chunk, serviceChoice)
			if err := s.messageCodec.Send(ctx, "requestSegment", wrapAPDU(apdu, true, s.routedDest)); err != nil {
				return errors.Wrapf(err, "error sending request segment %d", seq)
			}
			sent++
		}

		ack, err := s.awaitMeaningfulAck(ctx, invokeId, lastAcked, ackCh, errCh)
		if err != nil {
			if ctx.Err() != nil {
				return errors.Wrap(ctx.Err(), "context done awaiting segment ack")
			}
			// Ack timeout: retransmit from the last acknowledged segment,
			// bounded by the configured APDU retry budget.
			retries++
			if retries > s.driverContext.configuration.ApduRetries {
				return errors.Wrapf(err, "no segment ack after %d retries", s.driverContext.configuration.ApduRetries)
			}
			segmenter.Rewind(uint8(lastAcked + 1))
			continue
		}
		retries = 0

		// The actual window size in the peer's ack is authoritative for the
		// remainder of the transmission (capped at our own proposal).
		if actual := ack.GetActualWindowSize(); actual > 0 {
			window = min(actual, proposedSegmentWindowSize)
		}

		if ackSeq := int(ack.GetSequenceNumber()); ackSeq > lastAcked {
			lastAcked = ackSeq
		}

		if lastAcked >= total-1 {
			if !ack.GetNegativeAck() {
				return nil
			}
			// A NAK naming the final segment is a protocol oddity; resend just
			// the final segment rather than looping forever.
			segmenter.Rewind(uint8(total - 1))
			lastAcked = total - 2
			continue
		}

		// Resume right after the last segment the peer confirmed in order.
		// After an in-order positive ack this is a no-op rewind; after a NAK
		// (or an ack for a partial window) it retransmits the missing tail.
		segmenter.Rewind(uint8(lastAcked + 1))
	}
}

// awaitMeaningfulAck waits on the registered expectation and swallows stale
// duplicate acks (positive acks for segments at or below lastAcked), re-arming
// the expectation until a meaningful ack, an error, or ctx cancellation.
func (s *segmentedRequestSender) awaitMeaningfulAck(ctx context.Context, invokeId uint8, lastAcked int, ackCh <-chan readWriteModel.APDUSegmentAck, errCh <-chan error) (readWriteModel.APDUSegmentAck, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			return nil, err
		case ack := <-ackCh:
			if !ack.GetNegativeAck() && int(ack.GetSequenceNumber()) <= lastAcked {
				s.log.Debug().Uint8("ackSeq", ack.GetSequenceNumber()).Msg("ignoring stale duplicate segment ack")
				ackCh, errCh = s.expectSegmentAck(ctx, invokeId)
				continue
			}
			return ack, nil
		}
	}
}

// segmentApdu builds one APDUConfirmedRequest segment. Segment 0 carries the
// service-choice byte as the first payload byte (the generated model only
// parses/serializes a discrete segmentServiceChoice for sequence numbers > 0),
// so chunks slice the serialized service request contiguously and the peer's
// reassembly concatenates them back to the original payload.
func (s *segmentedRequestSender) segmentApdu(invokeId uint8, seq uint8, moreFollows bool, chunk []byte, serviceChoice readWriteModel.BACnetConfirmedServiceChoice) readWriteModel.APDUConfirmedRequest {
	seqV := seq
	windowV := uint8(proposedSegmentWindowSize)
	var choicePtr *readWriteModel.BACnetConfirmedServiceChoice
	if seq != 0 {
		choicePtr = &serviceChoice
	}
	return readWriteModel.NewAPDUConfirmedRequest(
		true,
		moreFollows,
		true,
		s.driverContext.maxSegmentsAccepted,
		s.driverContext.maxApduLengthAccepted,
		invokeId,
		&seqV,
		&windowV,
		nil,
		choicePtr,
		chunk,
	)
}

// expectSegmentAck registers a one-shot expectation for the next SegmentAck
// with the given invoke id, returning channels for the ack or an error/timeout.
func (s *segmentedRequestSender) expectSegmentAck(ctx context.Context, invokeId uint8) (<-chan readWriteModel.APDUSegmentAck, <-chan error) {
	ackCh := make(chan readWriteModel.APDUSegmentAck, 1)
	errCh := make(chan error, 1)

	expectCtx, cancel := utils.WithNamedTimeout(ctx, "segment ack wait timeout", segmentAckWaitTimeout)

	m := s.messageCodec
	m.Expect(expectCtx, "requestSegmentAck",
		func(message spi.Message) bool {
			apdu, ok := apduFromMessage(message)
			if !ok {
				return false
			}
			ack, ok := apdu.(readWriteModel.APDUSegmentAck)
			return ok && ack.GetOriginalInvokeId() == invokeId
		},
		func(message spi.Message) error {
			cancel()
			apdu, _ := apduFromMessage(message)
			ackCh <- apdu.(readWriteModel.APDUSegmentAck)
			return nil
		},
		func(err error) error {
			cancel()
			errCh <- err
			return nil
		},
	)
	return ackCh, errCh
}

// responseMatcherExcludingSegmentAcks wraps a response matcher so it never
// consumes the peer's SegmentAcks for our own request segments — those belong
// to the segmentedRequestSender's expectations, which are armed concurrently
// on the same invoke id.
func responseMatcherExcludingSegmentAcks(matcher func(spi.Message) bool) func(spi.Message) bool {
	return func(message spi.Message) bool {
		if apdu, ok := apduFromMessage(message); ok {
			if _, isSegmentAck := apdu.(readWriteModel.APDUSegmentAck); isSegmentAck {
				return false
			}
		}
		return matcher(message)
	}
}
