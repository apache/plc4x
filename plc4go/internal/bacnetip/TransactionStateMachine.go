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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

// TransactionStateMachine is the implementation of the bacnet transaction state machine
type TransactionStateMachine struct {
	*MessageCodec
	deviceInventory       *DeviceInventory
	retryCount            int
	segmentRetryCount     int
	duplicateCount        int
	sentAllSegments       bool
	lastSequenceNumber    int
	initialSequenceNumber int
	actualWindowSize      int
	proposeWindowSize     int
	segmentTimer          int
	RequestTimer          int
}

func (t *TransactionStateMachine) GetCodec() spi.MessageCodec {
	return t
}

func (t *TransactionStateMachine) Send(message spi.Message) error {
	if handled, err := t.handleOutboundMessage(message); handled {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "Error handling message")
	} else {
		return t.MessageCodec.Send(message)
	}
}

func (t *TransactionStateMachine) Expect(ctx context.Context, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// TODO: detect overflow
	return t.MessageCodec.Expect(ctx, acceptsMessage, handleMessage, handleError, ttl)
}

func (t *TransactionStateMachine) SendRequest(ctx context.Context, message spi.Message, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// Note: this code is copied on purpose from default codec as we want to call "this" `Send` and `Expect`
	if err := ctx.Err(); err != nil {
		return errors.Wrap(err, "Not sending message as context is aborted")
	}
	log.Trace().Msg("Sending request")
	// Send the actual message
	err := t.Send(message)
	if err != nil {
		return errors.Wrap(err, "Error sending the request")
	}
	return t.Expect(ctx, acceptsMessage, handleMessage, handleError, ttl)
}

func (t *TransactionStateMachine) handleOutboundMessage(message spi.Message) (handled bool, err error) {
	switch message := message.(type) {
	case readWriteModel.BVLCExactly:
		bvlc := message
		var npdu readWriteModel.NPDU
		if npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU }); ok {
			npdu = npduRetriever.GetNpdu()
		} else {
			log.Debug().Msgf("bvlc has no way to give a npdu %T", bvlc)
			return false, nil
		}
		if npdu.GetControl().GetMessageTypeFieldPresent() {
			log.Trace().Msg("No message type field present")
			return false, nil
		}
		var entryForDestination = NoDeviceEntry
		if npdu.GetControl().GetDestinationSpecified() {
			var err error
			if entryForDestination, err = t.deviceInventory.getEntryForDestination(npdu.GetDestinationAddress()); err != nil {
				// Get information from the device first
				// TODO: get information with who-has maybe or directed... not sure now
				// TODO: set entry once received
			}
		}
		// TODO: should we continue if we don't have a destination
		_ = entryForDestination
		apdu := npdu.GetApdu()
		switch apdu := apdu.(type) {
		case readWriteModel.APDUConfirmedRequestExactly:
			// TODO: check if adpu length is the magic number (it should be "unencoded")
			maximumApduLengthForDevice := entryForDestination.MaximumApduLengthAcceptedLength.GetActualValue()
			apduLengthDoesOverflow := uint64(apdu.GetLengthInBytes()) > maximumApduLengthForDevice
			if apduLengthDoesOverflow && !entryForDestination.SegmentationSupported {
				return false, errors.Errorf("We have a overflow. We need %d bytes, but device only supports a max of %d and no segmentation", apdu.GetLengthInBytes(), maximumApduLengthForDevice)
			}
			// TODO: handle potential retry
			if apduLengthDoesOverflow {
				// TODO: handle potential segmentation
			}
			return false, nil
		case readWriteModel.APDUComplexAckExactly:
			// TODO: check if adpu length is the magic number (it should be "unencoded")
			maximumApduLengthForDevice := entryForDestination.MaximumApduLengthAcceptedLength.GetActualValue()
			apduLengthDoesOverflow := uint64(apdu.GetLengthInBytes()) > maximumApduLengthForDevice
			if apduLengthDoesOverflow && !entryForDestination.SegmentationSupported {
				return false, errors.Errorf("We have a overflow. We need %d bytes, but device only supports a max of %d and no segmentation", apdu.GetLengthInBytes(), maximumApduLengthForDevice)
			}
			if apduLengthDoesOverflow {
				// TODO: handle potential segmentation
			}
			return false, nil
		default:
			log.Trace().Msgf("APDU type not relevant %T present", apdu)
			return false, nil
		}
	default:
		log.Trace().Msgf("Message type not relevant %T present", message)
		return false, nil
	}
}
