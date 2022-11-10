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

func NewTransactionStateMachine(messageCodec *MessageCodec, deviceInventory *DeviceInventory) TransactionStateMachine {
	return TransactionStateMachine{
		MessageCodec:          messageCodec,
		deviceInventory:       deviceInventory,
		retryCount:            3,
		segmentRetryCount:     3,
		duplicateCount:        0,
		sentAllSegments:       false,
		lastSequenceNumber:    0,
		initialSequenceNumber: 0,
		actualWindowSize:      0,
		proposeWindowSize:     2,
		segmentTimer:          1500,
		RequestTimer:          3000,
	}
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
			log.Trace().Msg("Message type field present")
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
			// TODO: this is a "client" request
			// TODO: check if adpu length is the magic number (it should be "unencoded")
			return false, nil
		case readWriteModel.APDUComplexAckExactly:
			// TODO: this is a "server" response
			// TODO: check if adpu length is the magic number (it should be "unencoded")
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

// TODO: this is a placeholder for a tasking framework
type _Task struct {
	taskTime    time.Time
	isScheduled bool
}

func (t *_Task) installTask(when *time.Time, delta *time.Duration) {
	// TODO: schedule task
}

func (t *_Task) suspendTask() {
	// TODO: suspend task
}

func (t *_Task) resume() {
	// TODO: resume task
}

type OneShotTask struct {
	_Task
}

// TODO: this is the interface to the outside for the SSM
type ServiceAccessPoint interface {
	GetDeviceInventory() *DeviceInventory
	GetLocalDevice() DeviceEntry
	GetProposedWindowSize() uint8
	Request(apdu readWriteModel.APDU)
	// TODO: wrap that properly
	GetClientTransactions() []interface{}
}

type SSMState uint8

const (
	IDLE SSMState = iota
	SEGMENTED_REQUEST
	AWAIT_CONFIRMATION
	AWAIT_RESPONSE
	SEGMENTED_RESPONSE
	SEGMENTED_CONFIRMATION
	COMPLETED
	ABORTED
)

func (s SSMState) String() string {
	switch s {
	case IDLE:
		return "IDLE"
	case SEGMENTED_REQUEST:
		return "SEGMENTED_REQUEST"
	case AWAIT_CONFIRMATION:
		return "AWAIT_CONFIRMATION"
	case AWAIT_RESPONSE:
		return "AWAIT_RESPONSE"
	case SEGMENTED_RESPONSE:
		return "SEGMENTED_RESPONSE"
	case SEGMENTED_CONFIRMATION:
		return "SEGMENTED_CONFIRMATION"
	case COMPLETED:
		return "COMPLETED"
	case ABORTED:
		return "ABORTED"
	default:
		return "Unknown"
	}
}

type segmentAPDU struct {
	originalApdu     readWriteModel.APDU
	originalInvokeId uint8
	serviceBytes     []byte
	serviceChoice    readWriteModel.BACnetConfirmedServiceChoice
	isAck            bool
}

// SSM - Segmentation State Machine
type SSM struct {
	OneShotTask

	ssmSAP ServiceAccessPoint

	pduAddress  []byte
	deviceEntry DeviceEntry

	invokeId uint8

	state        SSMState
	segmentAPDU  *segmentAPDU // TODO: rename that to segmentAPDUSource or something
	segmentSize  uint
	segmentCount uint8

	retryCount            uint
	segmentRetryCount     uint
	sentAllSegments       bool
	lastSequenceNumber    uint8
	initialSequenceNumber uint8
	actualWindowSize      uint8

	numberOfApduRetries   uint
	apduTimeout           uint
	segmentationSupported readWriteModel.BACnetSegmentation
	segmentTimeout        uint
	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted
}

func NewSSM(sap ServiceAccessPoint, pduAddress []byte) (SSM, error) {
	deviceEntry, err := sap.GetDeviceInventory().getEntryForDestination(pduAddress)
	if err != nil {
		return SSM{}, errors.Wrap(err, "Can't create SSM")
	}
	localDevice := sap.GetLocalDevice()
	return SSM{
		ssmSAP:                sap,
		pduAddress:            pduAddress,
		deviceEntry:           deviceEntry,
		state:                 IDLE,
		numberOfApduRetries:   localDevice.NumberOfAPDURetries,
		apduTimeout:           localDevice.APDUTimeout,
		segmentationSupported: localDevice.SegmentationSupported,
		segmentTimeout:        localDevice.APDUSegmentTimeout,
		maxSegmentsAccepted:   localDevice.MaxSegmentsAccepted,
		maxApduLengthAccepted: localDevice.MaximumApduLengthAccepted,
	}, nil
}

func (s *SSM) startTimer(millis int64) {
	s.restartTimer(millis)
}

func (s *SSM) stopTimer() {
	if s.isScheduled {
		s.suspendTask()
	}
}

func (s *SSM) restartTimer(millis int64) {
	if s.isScheduled {
		s.suspendTask()
	}

	delta := time.Millisecond * time.Duration(millis)
	s.installTask(nil, &delta)
}

func (s *SSM) setState(newState SSMState, timer *int64) error {
	if s.state == COMPLETED || s.state == ABORTED {
		return errors.Errorf("Invalid state transition from %s to %s", s.state, newState)
	}

	s.stopTimer()

	s.state = newState

	if timer != nil {
		s.startTimer(*timer)
	}
	return nil
}

func (s *SSM) setSegmentationContext(apdu readWriteModel.APDU) error {
	switch apdu := apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		bytes, err := apdu.GetServiceRequest().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		segmentAPDU := segmentAPDU{
			originalApdu:     apdu,
			originalInvokeId: apdu.GetInvokeId(),
			serviceBytes:     bytes,
			serviceChoice:    apdu.GetServiceRequest().GetServiceChoice(),
		}
		s.segmentAPDU = &segmentAPDU
	case readWriteModel.APDUComplexAckExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		bytes, err := apdu.GetServiceAck().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		segmentAPDU := segmentAPDU{
			originalApdu:  apdu,
			serviceBytes:  bytes,
			serviceChoice: apdu.GetServiceAck().GetServiceChoice(),
			isAck:         true,
		}
		s.segmentAPDU = &segmentAPDU
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

func (s *SSM) getSegment(index uint8) (segmentAPDU readWriteModel.APDU, moreFollows bool, err error) {
	if s.segmentAPDU == nil {
		return nil, false, errors.New("No segment apdu set")
	}

	if index > s.segmentCount {
		return nil, false, errors.Errorf("Invalid segment number %d, APDU has %d segments", index, s.segmentCount)
	}

	// TODO: the original code does here something funky but it seems it is best to just return the original apdu
	if s.segmentCount == 1 {
		return s.segmentAPDU.originalApdu, false, nil
	}

	moreFollows = index < s.segmentCount-1
	sequenceNumber := index % 255
	proposedWindowSize := s.actualWindowSize
	if index == 0 {
		proposedWindowSize = s.ssmSAP.GetProposedWindowSize()
	}
	serviceChoice := &s.segmentAPDU.serviceChoice
	offset := uint(index) * s.segmentSize
	segmentBytes := s.segmentAPDU.serviceBytes[offset : offset+s.segmentSize]
	if !s.segmentAPDU.isAck {
		segmentAPDU = readWriteModel.NewAPDUConfirmedRequest(
			true,
			moreFollows,
			s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE || s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_BOTH,
			s.maxSegmentsAccepted,
			s.maxApduLengthAccepted,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			&proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		)
	} else {
		segmentAPDU = readWriteModel.NewAPDUComplexAck(
			true,
			moreFollows,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			&proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		)
	}
	return segmentAPDU, moreFollows, nil
}

// TODO: check that function. looks a bit wonky to just append the payloads like that
func (s *SSM) appendSegment(apdu readWriteModel.APDU) error {
	switch apdu := apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		bytes, err := apdu.GetServiceRequest().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		s.segmentAPDU.serviceBytes = append(s.segmentAPDU.serviceBytes, bytes...)
	case readWriteModel.APDUComplexAckExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		bytes, err := apdu.GetServiceAck().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		s.segmentAPDU.serviceBytes = append(s.segmentAPDU.serviceBytes, bytes...)
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

func (s *SSM) inWindow(sequenceA, sequenceB uint) bool {
	return (sequenceA-sequenceB-256)%256 < uint(s.actualWindowSize)
}

func (s *SSM) fillWindow(sequenceNumber uint8) error {
	for i := uint8(0); i < s.actualWindowSize; i++ {
		apdu, moreFollows, err := s.getSegment(sequenceNumber + i)
		if err != nil {
			return errors.Wrapf(err, "Error sending out segment %d", i)
		}
		s.ssmSAP.Request(apdu)
		if moreFollows {
			s.sentAllSegments = true
		}
	}
	return nil
}

type ClientSSM struct {
	SSM
}

func NewClientSSM(sap ServiceAccessPoint, pduAddress []byte) (ClientSSM, error) {
	ssm, err := NewSSM(sap, pduAddress)
	if err != nil {
		return ClientSSM{}, err
	}
	// TODO: if deviceEntry is not there get it now...
	if &ssm.deviceEntry == &NoDeviceEntry {
		// TODO: get entry for device, store it in inventory
	}
	return ClientSSM{
		SSM: ssm,
	}, nil
}

func (s *ClientSSM) setState(newState SSMState, timer *int64) error {
	// do the regular state change
	if err := s.SSM.setState(newState, timer); err != nil {
		return errors.Wrap(err, "error during SSM state transition")
	}

	if s.state == COMPLETED || s.state == ABORTED {
		s.ssmSAP.GetClientTransactions() // TODO remove this
		if &s.deviceEntry != &NoDeviceEntry {
			// TODO: release device entry
		}
	}
	return nil
}

func (s *ClientSSM) request(apdu readWriteModel.APDU) {
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	s.ssmSAP.Request(apdu)
}
