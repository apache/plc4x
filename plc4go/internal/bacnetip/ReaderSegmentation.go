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

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// segmentWaitTimeout bounds how long we wait for each follow-up segment before
// giving up on a segmented response.
const segmentWaitTimeout = 5 * time.Second

// reassembleSegmentedRead drives the BACnet segmented-response protocol for a
// read: it acknowledges each segment and waits for the next until the device
// signals no more follow, then reparses the concatenated payload into a service
// ack and produces the final PLC4X read result.
//
// It runs in its own goroutine (NOT under the codec expectation lock) so it is
// free to register fresh expectations and send segment acks.
func (m *Reader) reassembleSegmentedRead(ctx context.Context, readRequest apiModel.PlcReadRequest, first readWriteModel.APDUComplexAck, result chan apiModel.PlcReadRequestResult) {
	invokeId := first.GetOriginalInvokeId()

	// Advertise a window size of 1 (ack every segment) — simplest and most
	// broadly compatible; the reassembler echoes this in each SegmentAck.
	reassembler := NewInboundReassembler(invokeId, 1)

	ack, err := reassembler.AcceptSegment(first)
	if err != nil {
		m.failSegmentedRead(readRequest, result, errors.Wrap(err, "error accepting first segment"))
		return
	}

	for {
		if reassembler.Complete() {
			// Acknowledge the final segment (best-effort).
			if sendErr := m.sendSegmentAck(ctx, ack); sendErr != nil {
				m.log.Debug().Err(sendErr).Msg("error sending final segment ack")
			}
			break
		}

		// Register the expectation for the next segment BEFORE acking the
		// current one, otherwise the device's next segment could race ahead of
		// our expectation and be dropped.
		segCh, errCh := m.expectSegment(ctx, invokeId)
		if sendErr := m.sendSegmentAck(ctx, ack); sendErr != nil {
			m.failSegmentedRead(readRequest, result, errors.Wrap(sendErr, "error sending segment ack"))
			return
		}

		select {
		case <-ctx.Done():
			m.failSegmentedRead(readRequest, result, errors.Wrap(ctx.Err(), "context cancelled awaiting segment"))
			return
		case waitErr := <-errCh:
			m.failSegmentedRead(readRequest, result, errors.Wrap(waitErr, "error awaiting segment"))
			return
		case seg := <-segCh:
			ack, err = reassembler.AcceptSegment(seg)
			if err != nil {
				m.failSegmentedRead(readRequest, result, errors.Wrap(err, "error accepting segment"))
				return
			}
		}
	}

	// Reparse the concatenated payload (service-choice byte + data) into a
	// service ack, then decode it the same way as a single-APDU response.
	payload := reassembler.Bytes()
	serviceAck, err := readWriteModel.BACnetServiceAckParse[readWriteModel.BACnetServiceAck](ctx, payload, uint32(len(payload)))
	if err != nil {
		m.failSegmentedRead(readRequest, result, errors.Wrap(err, "error parsing reassembled service ack"))
		return
	}

	readResponse, err := m.decodeServiceAck(serviceAck, readRequest)
	if err != nil {
		m.failSegmentedRead(readRequest, result, errors.Wrap(err, "error decoding reassembled response"))
		return
	}
	utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, readResponse, nil))
}

func (m *Reader) failSegmentedRead(readRequest apiModel.PlcReadRequest, result chan apiModel.PlcReadRequestResult, err error) {
	m.log.Debug().Err(err).Msg("segmented read failed")
	utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err))
}

// expectSegment registers a one-shot expectation for the next segmented
// APDUComplexAck with the given invoke id, returning channels for the segment or
// an error/timeout.
func (m *Reader) expectSegment(ctx context.Context, invokeId uint8) (<-chan readWriteModel.APDUComplexAck, <-chan error) {
	segCh := make(chan readWriteModel.APDUComplexAck, 1)
	errCh := make(chan error, 1)

	expectCtx, cancel := utils.WithNamedTimeout(ctx, "segment wait timeout", segmentWaitTimeout)

	m.messageCodec.Expect(expectCtx, "readSegment",
		func(message spi.Message) bool {
			apdu, ok := apduFromMessage(message)
			if !ok {
				return false
			}
			complexAck, ok := apdu.(readWriteModel.APDUComplexAck)
			return ok && complexAck.GetSegmentedMessage() && complexAck.GetOriginalInvokeId() == invokeId
		},
		func(message spi.Message) error {
			cancel()
			apdu, _ := apduFromMessage(message)
			segCh <- apdu.(readWriteModel.APDUComplexAck)
			return nil
		},
		func(err error) error {
			cancel()
			errCh <- err
			return nil
		},
	)
	return segCh, errCh
}

// sendSegmentAck transmits a BACnet SegmentAck wrapped in NPDU+BVLC.
func (m *Reader) sendSegmentAck(ctx context.Context, ack readWriteModel.APDUSegmentAck) error {
	if ack == nil {
		return nil
	}
	return m.messageCodec.Send(ctx, "segmentAck", wrapAPDU(ack, false, m.routedDest))
}

// apduFromMessage safely extracts the APDU from a received BVLC message,
// returning ok=false rather than panicking on unexpected shapes.
func apduFromMessage(message spi.Message) (readWriteModel.APDU, bool) {
	bvlc, ok := message.(readWriteModel.BVLC)
	if !ok {
		return nil, false
	}
	npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU })
	if !ok {
		return nil, false
	}
	npdu := npduRetriever.GetNpdu()
	if npdu == nil || npdu.GetControl().GetMessageTypeFieldPresent() {
		return nil, false
	}
	apdu := npdu.GetApdu()
	if apdu == nil {
		return nil, false
	}
	return apdu, true
}
