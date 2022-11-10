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
	"bytes"
	"context"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
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
		var entryForDestination = DeviceEntryDefault
		if npdu.GetControl().GetDestinationSpecified() {
			if retrievedEntry, err := t.deviceInventory.getEntryForDestination(npdu.GetDestinationAddress()); err != nil {
				// Get information from the device first
				// TODO: get information with who-has maybe or directed... not sure now
				// TODO: set entry once received
				_ = retrievedEntry
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

// TODO: this is the interface to the outside for the SSM // TODO: maybe we should port that as non interface first
type ServiceAccessPoint interface {
	GetDeviceInventory() *DeviceInventory
	GetLocalDevice() DeviceEntry
	GetProposedWindowSize() uint8
	Request(apdu readWriteModel.APDU)
	SapRequest(apdu readWriteModel.APDU)
	SapResponse(apdu readWriteModel.APDU)
	// TODO: wrap that properly
	GetClientTransactions() []interface{}
	// TODO: wrap that properly
	GetServerTransactions() []interface{}
	GetApplicationTimeout() *uint
}

// TODO: interface to client // TODO: maybe we should port that as non interface first
type Client interface {
	request(apdu readWriteModel.APDU)
	confirmation(apdu readWriteModel.APDU)
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
	deviceEntry *DeviceEntry

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
	actualWindowSize      *uint8

	numberOfApduRetries   uint
	apduTimeout           uint
	segmentationSupported readWriteModel.BACnetSegmentation
	segmentTimeout        uint
	maxSegmentsAccepted   *readWriteModel.MaxSegmentsAccepted
	maxApduLengthAccepted *readWriteModel.MaxApduLengthAccepted
}

func NewSSM(sap ServiceAccessPoint, pduAddress []byte) (SSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pdu_address", pduAddress).Msg("init")
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

func (s *SSM) startTimer(millis uint) {
	log.Debug().Msgf("Start timer %d", millis)
	s.restartTimer(millis)
}

func (s *SSM) stopTimer() {
	log.Debug().Msg("Stop Timer")
	if s.isScheduled {
		log.Debug().Msg("is scheduled")
		s.suspendTask()
	}
}

func (s *SSM) restartTimer(millis uint) {
	log.Debug().Msgf("restartTimer %d", millis)
	if s.isScheduled {
		log.Debug().Msg("is scheduled")
		s.suspendTask()
	}

	delta := time.Millisecond * time.Duration(millis)
	s.installTask(nil, &delta)
}

// setState This function is called when the derived class wants to change state
func (s *SSM) setState(newState SSMState, timer *uint) error {
	log.Debug().Msgf("setState %s timer=%d", newState, timer)
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

// setSegmentationContext This function is called to set the segmentation context
func (s *SSM) setSegmentationContext(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("setSegmentationContext\n%s", apdu)
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

// getSegment This function returns an APDU coorisponding to a particular segment of a confirmed request or complex ack.
//         The segmentAPDU is the context
func (s *SSM) getSegment(index uint8) (segmentAPDU readWriteModel.APDU, moreFollows bool, err error) {
	log.Debug().Msgf("Get segment %d", index)
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
		getProposedWindowSize := s.ssmSAP.GetProposedWindowSize()
		proposedWindowSize = &getProposedWindowSize
	}
	log.Debug().Msgf("proposedWindowSize %d", proposedWindowSize)
	serviceChoice := &s.segmentAPDU.serviceChoice
	offset := uint(index) * s.segmentSize
	segmentBytes := s.segmentAPDU.serviceBytes[offset : offset+s.segmentSize]
	if !s.segmentAPDU.isAck {
		log.Debug().Msg("confirmed request context")
		segmentedResponseAccepted := s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE || s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_BOTH
		log.Debug().Msgf("segmentedResponseAccepted %t", segmentedResponseAccepted)
		segmentAPDU = readWriteModel.NewAPDUConfirmedRequest(
			true,
			moreFollows,
			segmentedResponseAccepted,
			*s.maxSegmentsAccepted,
			*s.maxApduLengthAccepted,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		)
	} else {
		log.Debug().Msg("complex ack context")
		segmentAPDU = readWriteModel.NewAPDUComplexAck(
			true,
			moreFollows,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		)
	}
	return segmentAPDU, moreFollows, nil
}

// TODO: check that function. looks a bit wonky to just append the payloads like that
// appendSegment This function appends the apdu content to the end of the current APDU being built.  The segmentAPDU is
//        the context
func (s *SSM) appendSegment(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("appendSegment\n%s", apdu)
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

func (s *SSM) inWindow(sequenceA, sequenceB uint8) bool {
	log.Debug().Msgf("inWindow %d-%d", sequenceA, sequenceB)
	return (uint(sequenceA)-uint(sequenceB)-256)%256 < uint(*s.actualWindowSize)
}

func (s *SSM) fillWindow(sequenceNumber uint8) error {
	log.Debug().Msgf("fillWindow %d", sequenceNumber)
	for i := uint8(0); i < *s.actualWindowSize; i++ {
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

func NewClientSSM(sap ServiceAccessPoint, pduAddress []byte) (*ClientSSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pduAddress", pduAddress).Msg("init")
	ssm, err := NewSSM(sap, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if ssm.deviceEntry == nil {
		// TODO: get entry for device, store it in inventory
		log.Debug().Msg("Accquire device information")
	}
	return &ClientSSM{
		SSM: ssm,
	}, nil
}

// setState This function is called when the client wants to change state
func (s *ClientSSM) setState(newState SSMState, timer *uint) error {
	log.Debug().Msgf("setState %s timer=%d", newState, timer)
	// do the regular state change
	if err := s.SSM.setState(newState, timer); err != nil {
		return errors.Wrap(err, "error during SSM state transition")
	}

	if s.state == COMPLETED || s.state == ABORTED {
		log.Debug().Msg("remove from active transaction")
		s.ssmSAP.GetClientTransactions() // TODO remove "this" transaction from the list
		if s.deviceEntry == nil {
			// TODO: release device entry
			log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// request This function is called by client transaction functions when it wants to send a message to the device
func (s *ClientSSM) request(apdu readWriteModel.APDU) {
	log.Debug().Msgf("request\n%s", apdu)
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	s.ssmSAP.Request(apdu)
}

// TODO: maybe use another name for that
// indication This function is called after the device has bound a new transaction and wants to start the process
//        rolling
func (s *ClientSSM) indication(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("indication\n%s", apdu)
	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		return errors.Errorf("Invalid APDU type %T", apdu)
	}

	// save the request and set the segmentation context
	if err := s.setSegmentationContext(apdu); err != nil {
		return errors.Wrap(err, "error setting context")
	}

	// if the max apdu length of the server isn't known, assume that it is the same size as our own and will be the segment
	//        size
	if s.deviceEntry == nil || s.deviceEntry.MaximumApduLengthAccepted != nil {
		s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
	} else if s.deviceEntry.MaximumNpduLength == nil {
		//      if the max npdu length of the server isn't known, assume that it is the same as the max apdu length accepted
		s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
	} else {
		s.segmentSize = utils.Min(*s.deviceEntry.MaximumNpduLength, uint(s.maxApduLengthAccepted.NumberOfOctets()))
	}
	log.Debug().Msgf("segment size %d", s.segmentSize)

	s.invokeId = apduConfirmedRequest.GetInvokeId()
	log.Debug().Msgf("invoke ID: %d", s.invokeId)

	var segmentCount, more int
	segmentCount, more = len(s.segmentAPDU.serviceBytes)/int(s.segmentSize), len(s.segmentAPDU.serviceBytes)%int(s.segmentSize)
	s.segmentCount = uint8(segmentCount)
	if more > 0 {
		s.segmentCount += 1
	}
	log.Debug().Msgf("segment count %d", segmentCount)

	if s.segmentCount > 1 {
		if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			log.Debug().Msg("local device can't send segmented requests")
			abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			s.response(abort)
			return nil
		}

		if s.deviceEntry == nil {
			log.Debug().Msg("no server info for segmentation support")
		} else if s.deviceEntry.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && s.deviceEntry.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			log.Debug().Msg("server can't receive segmented requests")
			abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			s.response(abort)
			return nil
		}

		// make sure we don't exceed the number of segments in our request that the server said it was willing to accept
		if s.deviceEntry == nil {
			log.Debug().Msg("no server info for maximum number of segments")
		} else if s.deviceEntry.MaxSegmentsAccepted == nil {
			log.Debug().Msgf("server doesn't say maximum number of segments")
		} else if s.segmentCount > s.deviceEntry.MaxSegmentsAccepted.MaxSegments() {
			log.Debug().Msg("server can't receive enough segments")
			abort, err := s.abort(readWriteModel.BACnetAbortReason_APDU_TOO_LONG)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			s.response(abort)
			return nil
		}
	}

	// send out the first segment (or the whole thing)
	if s.segmentCount == 1 {
		// unsegmented
		s.sentAllSegments = true
		s.retryCount = 0
		if err := s.setState(AWAIT_CONFIRMATION, &s.apduTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	} else {
		// segmented
		s.sentAllSegments = false
		s.retryCount = 0
		s.segmentRetryCount = 0
		s.initialSequenceNumber = 0
		s.actualWindowSize = nil
		if err := s.setState(SEGMENTED_REQUEST, &s.segmentTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	}

	// deliver to the device
	segment, _, err := s.getSegment(0)
	if err != nil {
		return errors.Wrap(err, "error getting segment")
	}
	s.request(segment)
	return nil
}

// response This function is called by client transaction functions when they want to send a message to the application.
func (s *ClientSSM) response(apdu readWriteModel.APDU) {
	log.Debug().Msgf("response\n%s", apdu)
	// make sure it has a good source and destination
	// TODO: check if source == s.pduAddress
	// TODO: check if

	// send it to the application
	s.ssmSAP.SapResponse(apdu)
}

// confirmation This function is called by the device for all upstream messages related to the transaction.
func (s *ClientSSM) confirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("confirmation\n%s", apdu)

	switch s.state {
	case SEGMENTED_REQUEST:
		return s.segmentedRequest(apdu)
	case AWAIT_CONFIRMATION:
		return s.awaitConfirmation(apdu)
	case SEGMENTED_CONFIRMATION:
		return s.segmentedConfirmation(apdu)
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// processTask This function is called when something has taken too long
func (s *ClientSSM) processTask() error {
	log.Debug().Msg("processTask")
	switch s.state {
	case SEGMENTED_REQUEST:
		return s.segmentedRequestTimeout()
	case AWAIT_CONFIRMATION:
		return s.awaitConfirmationTimeout()
	case SEGMENTED_CONFIRMATION:
		return s.segmentedConfirmationTimeout()
	case COMPLETED, ABORTED:
		return nil
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// abort This function is called when the transaction should be aborted
func (s *ClientSSM) abort(reason readWriteModel.BACnetAbortReason) (readWriteModel.APDU, error) {
	log.Debug().Msgf("abort\n%s", reason)

	// change the state to aborted
	if err := s.setState(ABORTED, nil); err != nil {
		return nil, errors.Wrap(err, "Error setting state to aborted")
	}

	// build an abort PDU to return
	abortApdu := readWriteModel.NewAPDUAbort(false, s.invokeId, readWriteModel.NewBACnetAbortReasonTagged(reason, uint32(reason), 0), 0)
	// return it
	return abortApdu, nil
}

// segmentedRequest This function is called when the client is sending a segmented request and receives an apdu
func (s *ClientSSM) segmentedRequest(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("segmentedRequest\n%s", apdu)

	switch apdu := apdu.(type) {
	// server is ready for the next segment
	case readWriteModel.APDUSegmentAckExactly:
		log.Debug().Msg("segment ack")
		getActualWindowSize := apdu.GetActualWindowSize()
		s.actualWindowSize = &getActualWindowSize

		// duplicate ack received?
		if !s.inWindow(apdu.GetSequenceNumber(), s.initialSequenceNumber) {
			log.Debug().Msg("not in window")
			s.restartTimer(s.segmentTimeout)
		} else if s.sentAllSegments {
			log.Debug().Msg("all done sending request")

			if err := s.setState(AWAIT_CONFIRMATION, &s.apduTimeout); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		} else {
			log.Debug().Msg("More segments to send")

			s.initialSequenceNumber = apdu.GetSequenceNumber() + 1
			s.retryCount = 0
			if err := s.fillWindow(s.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
			s.restartTimer(s.segmentTimeout)
		}
	// simple ack
	case readWriteModel.APDUSimpleAckExactly:
		log.Debug().Msg("simple ack")

		if !s.sentAllSegments {
			abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			s.request(abort)  // send it ot the device
			s.response(abort) // send it ot the application
		} else {
			if err := s.setState(COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	// complex ack
	case readWriteModel.APDUComplexAckExactly:
		log.Debug().Msg("complex ack")
		if !s.sentAllSegments {
			abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			s.request(abort)  // send it ot the device
			s.response(abort) // send it ot the application
		} else if !apdu.GetSegmentedMessage() {
			// ack is not segmented
			if err := s.setState(COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			s.response(apdu)
		} else {
			// set the segmented response context
			if err := s.setSegmentationContext(apdu); err != nil {
				return errors.Wrap(err, "error setting context")
			}

			// minimum of what the server is proposing and this client proposes
			minWindowSize := utils.Min(*apdu.GetProposedWindowSize(), s.ssmSAP.GetProposedWindowSize())
			s.actualWindowSize = &minWindowSize
			s.lastSequenceNumber = 0
			s.initialSequenceNumber = 0
			if err := s.setState(SEGMENTED_CONFIRMATION, &s.segmentTimeout); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	case readWriteModel.APDUErrorExactly:
		log.Debug().Msg("error/reject/abort")
		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		s.response(apdu)
	default:
		return errors.Errorf("Invalid apdu %T", apdu)
	}
	return nil
}

func (s *ClientSSM) segmentedRequestTimeout() error {
	log.Debug().Msg("segmentedRequestTimeout")

	// Try again
	if s.segmentRetryCount < s.numberOfApduRetries {
		log.Debug().Msg("retry segmented request")
		s.segmentRetryCount++
		s.startTimer(s.segmentTimeout)

		if s.initialSequenceNumber == 0 {
			apdu, _, err := s.getSegment(0)
			if err != nil {
				return errors.Wrap(err, "error getting first segment")
			}
			s.request(apdu)
		} else {
			if err := s.fillWindow(s.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
		}
	} else {
		log.Debug().Msg("abort, no response from the device")

		abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.response(abort)
	}
	return nil
}

func (s *ClientSSM) awaitConfirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("awaitConfirmation\n%s", apdu)

	switch apdu := apdu.(type) {
	case readWriteModel.APDUAbortExactly:
		log.Debug().Msg("Server aborted")

		if err := s.setState(ABORTED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		s.response(apdu)
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		log.Debug().Msg("simple ack, error or reject")

		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		s.response(apdu)
	case readWriteModel.APDUComplexAckExactly:
		log.Debug().Msg("complex ack")

		if !apdu.GetSegmentedMessage() {
			log.Debug().Msg("unsegmented")

			if err := s.setState(COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			s.response(apdu)
		} else if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			log.Debug().Msg("local device can't receive segmented messages")

			abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			s.response(abort)
		} else if *apdu.GetSequenceNumber() == 0 {
			log.Debug().Msg("segmented response")

			// set the segmented response context
			if err := s.setSegmentationContext(apdu); err != nil {
				return errors.Wrap(err, "error set segmentation context")
			}

			s.actualWindowSize = apdu.GetProposedWindowSize()
			s.lastSequenceNumber = 0
			s.initialSequenceNumber = 0
			if err := s.setState(SEGMENTED_CONFIRMATION, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}

			// send back a segment ack
			segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
			s.request(segmentAck)
		} else {
			log.Debug().Msg("Invalid apdu in this state")

			abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			s.request(abort)  // send it to the device
			s.response(abort) // send it to the application
		}
	case readWriteModel.APDUSegmentAckExactly:
		log.Debug().Msg("segment ack(!?)")
		s.restartTimer(s.segmentTimeout)
	default:
		return errors.Errorf("invalid apdu %T", apdu)
	}
	return nil
}

func (s *ClientSSM) awaitConfirmationTimeout() error {
	log.Debug().Msg("awaitConfirmationTimeout")

	if s.retryCount < s.numberOfApduRetries {
		log.Debug().Msgf("no response, try again (%d < %d)", s.retryCount, s.numberOfApduRetries)
		s.retryCount++

		// save the retry count, indication acts like the request is coming from the application so the retryCount gets
		//            re-initialized.
		saveCount := s.retryCount
		if err := s.indication(s.segmentAPDU.originalApdu); err != nil { // TODO: check that it is really the intention to re-send the original apdu here
			return err
		}
		s.retryCount = saveCount
	} else {
		log.Debug().Msg("retry count exceeded")

		abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.response(abort)
	}
	return nil
}

func (s *ClientSSM) segmentedConfirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("segmentedConfirmation\n%s", apdu)

	// the only messages we should be getting are complex acks
	apduComplexAck, ok := apdu.(readWriteModel.APDUComplexAckExactly)
	if !ok {
		log.Debug().Msg("complex ack required")

		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.request(abort)  // send it to the device
		s.response(abort) // send it to the application
	}

	// it must be segmented
	if !apduComplexAck.GetSegmentedMessage() {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.request(abort)  // send it to the device
		s.response(abort) // send it to the application
	}

	// proper segment number
	if *apduComplexAck.GetSequenceNumber() != s.lastSequenceNumber+1 {
		log.Debug().Msgf("segment %d received out of order, should be %d", apduComplexAck.GetSequenceNumber(), s.lastSequenceNumber+1)

		// segment received out of order
		s.restartTimer(s.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(true, false, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		s.request(segmentAck)
		return nil
	}

	// add the data
	if err := s.appendSegment(apdu); err != nil {
		return errors.Wrap(err, "error appending the segment")
	}

	// update the sequence number
	s.lastSequenceNumber = s.lastSequenceNumber + 1

	// last segment received
	if !apduComplexAck.GetMoreFollows() {
		log.Debug().Msg("No more follows")

		// send final ack
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, s.invokeId, s.lastSequenceNumber, *s.actualWindowSize, 0)
		s.request(segmentAck)

		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now...
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		parse, err := readWriteModel.APDUParse(s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		s.response(parse)
	} else if *apduComplexAck.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.restartTimer(s.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, s.invokeId, s.lastSequenceNumber, *s.actualWindowSize, 0)
		s.request(segmentAck)
	} else {
		log.Debug().Msg("Wait for more segments")

		s.restartTimer(s.segmentTimeout)
	}

	return nil
}

func (s *ClientSSM) segmentedConfirmationTimeout() error {
	log.Debug().Msg("segmentedConfirmationTimeout")

	abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
	if err != nil {
		return errors.Wrap(err, "error creating abort")
	}
	s.response(abort)
	return nil
}

type ServerSSM struct {
	SSM
	segmentedResponseAccepted bool
}

func NewServerSSM(sap ServiceAccessPoint, pduAddress []byte) (*ServerSSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pduAddress", pduAddress).Msg("init")
	ssm, err := NewSSM(sap, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if &ssm.deviceEntry == nil {
		// TODO: get entry for device, store it in inventory
		log.Debug().Msg("Accquire device information")
	}
	return &ServerSSM{
		SSM:                       ssm,
		segmentedResponseAccepted: true,
	}, nil
}

// setState This function is called when the client wants to change state
func (s *ServerSSM) setState(newState SSMState, timer *uint) error {
	log.Debug().Msgf("setState %s timer=%d", newState, timer)
	// do the regular state change
	if err := s.SSM.setState(newState, timer); err != nil {
		return errors.Wrap(err, "error during SSM state transition")
	}

	if s.state == COMPLETED || s.state == ABORTED {
		log.Debug().Msg("remove from active transaction")
		s.ssmSAP.GetServerTransactions() // TODO remove "this" transaction from the list
		if s.deviceEntry != nil {
			// TODO: release device entry
			log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// request This function is called by transaction functions to send to the application
func (s *ServerSSM) request(apdu readWriteModel.APDU) {
	log.Debug().Msgf("request\n%s", apdu)
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	s.ssmSAP.SapRequest(apdu)
}

// TODO: maybe use another name for that
// indication This function is called for each downstream packet related to
//        the transaction
func (s *ServerSSM) indication(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("indication\n%s", apdu)
	// make sure we're getting confirmed requests

	switch s.state {
	case IDLE:
		return s.idle(apdu)
	case SEGMENTED_REQUEST:
		return s.segmentedRequest(apdu)
	case AWAIT_RESPONSE:
		return s.awaitResponse(apdu)
	case SEGMENTED_RESPONSE:
		return s.segmentedResponse(apdu)
	default:
		return errors.Errorf("invalid state %s", s.state)
	}
}

// response This function is called by client transaction functions when they want to send a message to the application.
func (s *ServerSSM) response(apdu readWriteModel.APDU) {
	log.Debug().Msgf("response\n%s", apdu)
	// make sure it has a good source and destination
	// TODO: check if source == none
	// TODO: check if destnation = s.pduAddress

	// send it via the device
	s.ssmSAP.Request(apdu)
}

// confirmation This function is called when the application has provided a response and needs it to be sent to the
//        client.
func (s *ServerSSM) confirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("confirmation\n%s", apdu)

	// check to see we are in the correct state
	if s.state != AWAIT_RESPONSE {
		log.Debug().Msg("warning: no expecting a response")
	}

	switch apdu := apdu.(type) {
	// abort response
	case readWriteModel.APDUAbortExactly:
		log.Debug().Msg("abort")

		if err := s.setState(ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// end the response to the device
		s.response(apdu)
		return nil
	// simple response
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		log.Debug().Msg("simple ack, error or reject")

		// transaction completed
		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// send the response to the device
		s.response(apdu)
		return nil
	// complex ack
	case readWriteModel.APDUComplexAckExactly:
		log.Debug().Msg("complex ack")

		// save the response and set the segmentation context
		if err := s.setSegmentationContext(apdu); err != nil {
			return errors.Wrap(err, "error settings segmentation context")
		}

		// the segment size is the minimum of the size of the largest packet that can be delivered to the client and the
		//            largest it can accept
		if s.deviceEntry == nil || s.deviceEntry.MaximumNpduLength == nil {
			s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
		} else {
			s.segmentSize = utils.Min(*s.deviceEntry.MaximumNpduLength, uint(s.maxApduLengthAccepted.NumberOfOctets()))
		}

		// compute the segment count
		if len(apdu.GetSegment()) == 0 {
			// always at least one segment
			s.segmentCount = 1
		} else {
			// split into chunks, maybe need one more
			var segmentCount, more int
			segmentCount, more = len(s.segmentAPDU.serviceBytes)/int(s.segmentSize), len(s.segmentAPDU.serviceBytes)%int(s.segmentSize)
			if more > 0 {
				s.segmentCount += 1
			}
			log.Debug().Msgf("segment count: %d", segmentCount)

			// make sure we support segmented transmit if we need to
			if s.segmentCount > 1 {
				log.Debug().Msgf("segmentation required, %d segments", s.segmentCount)

				// make sure we support segmented transmit
				if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
					log.Debug().Msg("server can't send segmented requests")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					s.response(abort)
					return nil
				}

				// make sure client supports segmented receive
				if !s.segmentedResponseAccepted {
					log.Debug().Msg("client can't receive segmented responses")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					s.response(abort)
					return nil
				}

				// make sure we don't exceed the number of segments in our response that the client said it was willing to accept
				//                in the request
				if s.maxSegmentsAccepted != nil && s.segmentCount > s.maxSegmentsAccepted.MaxSegments() {
					log.Debug().Msg("client can't receive enough segments")
					abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					s.response(abort)
					return nil
				}
			}

			// initialize the state
			s.segmentRetryCount = 0
			s.initialSequenceNumber = 0
			s.actualWindowSize = nil

			// send out the first segment (or the whole thing)
			if s.segmentCount == 1 {
				s.response(apdu)
				if err := s.setState(COMPLETED, nil); err != nil {
					return errors.Wrap(err, "Error setting state to aborted")
				}
			} else {
				segment, _, err := s.getSegment(0)
				if err != nil {
					return errors.Wrap(err, "error getting first segment")
				}
				s.response(segment)
				if err := s.setState(SEGMENTED_RESPONSE, nil); err != nil {
					return errors.Wrap(err, "Error setting state to aborted")
				}
			}
		}
	default:
		return errors.Errorf("Invalid APDU %T", apdu)
	}
	return nil
}

// processTask This function is called when the client has failed to send all the segments of a segmented request,
//        the application has taken too long to complete the request, or the client failed to ack the segments of a
//        segmented response
func (s *ServerSSM) processTask() error {
	log.Debug().Msg("processTask")
	switch s.state {
	case SEGMENTED_REQUEST:
		return s.segmentedRequestTimeout()
	case AWAIT_CONFIRMATION:
		return s.awaitResponseTimeout()
	case SEGMENTED_CONFIRMATION:
		return s.segmentedResponseTimeout()
	case COMPLETED, ABORTED:
		return nil
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// abort This function is called when the transaction should be aborted
func (s *ServerSSM) abort(reason readWriteModel.BACnetAbortReason) (readWriteModel.APDU, error) {
	log.Debug().Msgf("abort\n%s", reason)

	// change the state to aborted
	if err := s.setState(ABORTED, nil); err != nil {
		return nil, errors.Wrap(err, "Error setting state to aborted")
	}

	// build an abort PDU to return
	abortApdu := readWriteModel.NewAPDUAbort(true, s.invokeId, readWriteModel.NewBACnetAbortReasonTagged(reason, uint32(reason), 0), 0)
	// return it
	return abortApdu, nil
}

func (s *ServerSSM) idle(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("idle %s", apdu)

	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		return errors.Errorf("Invalid APDU type %T", apdu)
	}

	// save the invoke ID
	s.invokeId = apduConfirmedRequest.GetInvokeId()
	log.Debug().Msgf("invoke ID: %d", s.invokeId)

	// remember if the client accepts segmented responses
	s.segmentedResponseAccepted = apduConfirmedRequest.GetSegmentedResponseAccepted()

	// if there is a cache record, check to see if it needs to be updated
	if apduConfirmedRequest.GetSegmentedResponseAccepted() && s.deviceEntry != nil {
		switch s.deviceEntry.SegmentationSupported {
		case readWriteModel.BACnetSegmentation_NO_SEGMENTATION:
			log.Debug().Msg("client actually supports segmented receive")
			s.deviceEntry.SegmentationSupported = readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE

		// TODO: bacpypes updates the cache here but as we have a pointer  to the entry we should need that. Maybe we should because concurrency... lets see later
		case readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT:
			log.Debug().Msg("client actually supports both segmented transmit and receive")
			s.deviceEntry.SegmentationSupported = readWriteModel.BACnetSegmentation_SEGMENTED_BOTH

			// TODO: bacpypes updates the cache here but as we have a pointer  to the entry we should need that. Maybe we should because concurrency... lets see later
		case readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE, readWriteModel.BACnetSegmentation_SEGMENTED_BOTH:
		// all good
		default:
			return errors.New("invalid segmentation supported in device info")
		}
	}

	// decode the maximum that the client can receive in one APDU, and if  there is a value in the device information then
	//        use that one because  it came from reading device object property value or from an I-Am  message that was
	//        received
	getMaxApduLengthAccepted := apduConfirmedRequest.GetMaxApduLengthAccepted()
	s.maxApduLengthAccepted = &getMaxApduLengthAccepted
	if s.deviceEntry != nil && s.deviceEntry.MaximumApduLengthAccepted != nil {
		if *s.deviceEntry.MaximumApduLengthAccepted < *s.maxApduLengthAccepted {
			log.Debug().Msg("apdu max reponse encoding error")
		} else {
			s.maxApduLengthAccepted = s.deviceEntry.MaximumApduLengthAccepted
		}
	}
	log.Debug().Msgf("maxApduLengthAccepted %s", *s.maxApduLengthAccepted)

	// save the number of segments the client is willing to accept in the ack, if this is None then the value is unknown or more than 64
	getMaxSegmentsAccepted := apduConfirmedRequest.GetMaxSegmentsAccepted()
	s.maxSegmentsAccepted = &getMaxSegmentsAccepted

	// unsegmented request
	if len(apduConfirmedRequest.GetSegment()) <= 0 {
		if err := s.setState(AWAIT_RESPONSE, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		s.request(apdu)
		return nil
	}

	// make sure we support segmented requests
	if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.response(abort)
		return nil
	}

	// save the response and set the segmentation context
	if err := s.setSegmentationContext(apdu); err != nil {
		return errors.Wrap(err, "error settings segmentation context")
	}

	// the window size is the minimum of what I would propose and what the device has proposed
	minWindowSize := utils.Min(*apduConfirmedRequest.GetProposedWindowSize(), s.ssmSAP.GetProposedWindowSize())
	s.actualWindowSize = &minWindowSize
	log.Debug().Msgf("actualWindowSize? min(%d, %d) -> %d", apduConfirmedRequest.GetProposedWindowSize(), s.ssmSAP.GetProposedWindowSize(), s.actualWindowSize)

	// initialize the state
	s.lastSequenceNumber = 0
	s.initialSequenceNumber = 0
	if err := s.setState(SEGMENTED_REQUEST, &s.segmentTimeout); err != nil {
		return errors.Wrap(err, "Error setting state to aborted")
	}

	// send back a segment ack
	segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
	log.Debug().Msgf("segAck: %s", segack)
	s.response(segack)
	return nil
}

func (s *ServerSSM) segmentedRequest(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("segmentedRequest\n%s", apdu)

	// some kind of problem
	if _, ok := apdu.(readWriteModel.APDUAbortExactly); ok {
		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		s.response(apdu)
		return nil
	}

	// the only messages we should be getting are confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if castedApdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.request(abort)  // send it ot the device
		s.response(abort) // send it ot the application
	} else {
		apduConfirmedRequest = castedApdu
	}

	// it must be segmented
	if !apduConfirmedRequest.GetSegmentedMessage() {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		s.request(abort)  // send it ot the device
		s.response(abort) // send it ot the application
	}

	// proper segment number
	if *apduConfirmedRequest.GetSequenceNumber() != s.lastSequenceNumber+1 {
		log.Debug().Msgf("segment %d received out of order, should be %d", *apduConfirmedRequest.GetSequenceNumber(), s.lastSequenceNumber+1)

		// segment received out of order
		s.restartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(true, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		s.response(segack)
		return nil
	}

	// add the data
	if err := s.appendSegment(apdu); err != nil {
		return errors.Wrap(err, "error appending segment")
	}

	// update the sequence number
	s.lastSequenceNumber++

	// last segment?
	if !apduConfirmedRequest.GetMoreFollows() {
		log.Debug().Msg("No more follows")

		// send back the final segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.lastSequenceNumber, *s.actualWindowSize, 0)
		s.response(segack)

		// forward the whole thing to the application
		if err := s.setState(AWAIT_RESPONSE, s.ssmSAP.GetApplicationTimeout()); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		// TODO: here we need to rebuild again yada yada
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now..
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		parse, err := readWriteModel.APDUParse(s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		s.request(parse)
	} else if *apduConfirmedRequest.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.restartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		s.response(segack)
	} else {
		// wait for more segments
		s.restartTimer(s.segmentTimeout)
	}

	return nil
}

func (s *ServerSSM) segmentedRequestTimeout() error {
	log.Debug().Msg("segmentedRequestTimeout")

	// give up
	if err := s.setState(ABORTED, nil); err != nil {
		return errors.Wrap(err, "Error setting state to aborted")
	}
	return nil
}

func (s *ServerSSM) awaitResponse(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("awaitResponse\n%s", apdu)

	switch apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		log.Debug().Msg("client is trying this request again")
	case readWriteModel.APDUAbortExactly:
		log.Debug().Msg("client aborting this request")

		// forward to the application
		if err := s.setState(ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		s.request(apdu)
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}
	return nil
}

// awaitResponseTimeout This function is called when the application has taken too long to respond to a clients request.
//         The client has probably long since given up
func (s *ServerSSM) awaitResponseTimeout() error {
	log.Debug().Msg("awaitResponseTimeout")

	abort, err := s.abort(readWriteModel.BACnetAbortReason(64)) // Note: this is a proprietary code used by bacpypes for server timeout. We just use that here too to keep consistent
	if err != nil {
		return errors.Wrap(err, "error creating abort")
	}
	s.request(abort)
	return nil
}

func (s *ServerSSM) segmentedResponse(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("segmentedResponse\n%s", apdu)

	// client is ready for the next segment
	switch apdu := apdu.(type) {
	case readWriteModel.APDUSegmentAckExactly:
		log.Debug().Msg("segment ack")

		// actual window size is provided by client
		getActualWindowSize := apdu.GetActualWindowSize()
		s.actualWindowSize = &getActualWindowSize

		// duplicate ack received?
		if !s.inWindow(apdu.GetSequenceNumber(), s.initialSequenceNumber) {
			log.Debug().Msg("not in window")
			s.restartTimer(s.segmentTimeout)
		} else if s.sentAllSegments {
			// final ack received?
			log.Debug().Msg("all done sending response")
			if err := s.setState(COMPLETED, nil); err != nil {
				return errors.Wrap(err, "Error setting state to aborted")
			}
		} else {
			log.Debug().Msg("more segments to send")

			s.initialSequenceNumber = apdu.GetSequenceNumber() + 1
			actualWindowSize := apdu.GetActualWindowSize()
			s.actualWindowSize = &actualWindowSize
			s.segmentRetryCount = 0
			if err := s.fillWindow(s.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
			s.restartTimer(s.segmentRetryCount)
		}
	// some kind of problem
	case readWriteModel.APDUAbortExactly:
		if err := s.setState(COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		s.response(apdu)
	default:
		return errors.Errorf("Invalid APDU %T", apdu)
	}
	return nil
}

func (s *ServerSSM) segmentedResponseTimeout() error {
	log.Debug().Msg("segmentedResponseTimeout")

	// try again
	if s.segmentRetryCount < s.numberOfApduRetries {
		s.segmentRetryCount++
		s.startTimer(s.segmentTimeout)
		if err := s.fillWindow(s.initialSequenceNumber); err != nil {
			return errors.Wrap(err, "error filling window")
		}
	} else {
		// five up
		if err := s.setState(ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
	}
	return nil
}

type StateMachineAccessPoint struct {
	Client
	ServiceAccessPoint

	localDevice           DeviceEntry
	deviceInventory       *DeviceInventory
	nextInvokeId          uint8
	clientTransactions    []*ClientSSM
	serverTransactions    []*ServerSSM
	numberOfApduRetries   int
	apduTimeout           int
	maxApduLengthAccepted int
	segmentationSupported readWriteModel.BACnetSegmentation
	segmentTimeout        int
	maxSegmentsAccepted   int
	proposedWindowSize    int
	dccEnableDisable      readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable
	applicationTimeout    int
}

func NewStateMachineAccessPoint(localDevice DeviceEntry, deviceInventory *DeviceInventory, sapID int, cid int) StateMachineAccessPoint {
	log.Debug().Msgf("NewStateMachineAccessPoint localDevice=%v deviceInventory=%v sap=%v cid=%v", localDevice, deviceInventory, sapID, cid)

	// basic initialization
	// TODO: init client
	// TODO: init sap
	return StateMachineAccessPoint{
		// save a reference to the device information cache
		localDevice:     localDevice,
		deviceInventory: deviceInventory,

		// client settings
		nextInvokeId:       1,
		clientTransactions: nil,

		// server settings
		serverTransactions: nil,

		// confirmed request defaults
		numberOfApduRetries:   3,
		apduTimeout:           3000,
		maxApduLengthAccepted: 1024,

		// segmentation defaults
		segmentationSupported: readWriteModel.BACnetSegmentation_NO_SEGMENTATION,
		segmentTimeout:        1500,
		maxSegmentsAccepted:   2,
		proposedWindowSize:    2,

		// device communication control
		dccEnableDisable: readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE,

		// how long the state machine is willing to wait for the application
		// layer to form a response and send it
		applicationTimeout: 3000,
	}
}

// getNextInvokeId Called by clients to get an unused invoke ID
func (s *StateMachineAccessPoint) getNextInvokeId(address []byte) (uint8, error) {
	log.Debug().Msg("getNextInvokeId")

	initialID := s.nextInvokeId
	for {
		invokeId := s.nextInvokeId
		s.nextInvokeId++

		// see if we've checked for them all
		if initialID == s.nextInvokeId {
			return 0, errors.New("No available invoke ID")
		}

		if len(s.clientTransactions) == 0 {
			return invokeId, nil
		}

		// TODO: double check that the logic here is right
		for _, tr := range s.clientTransactions {
			if invokeId == tr.invokeId && bytes.Equal(address, tr.pduAddress) {
				return invokeId, nil
			}
		}
	}
}

// confirmation Packets coming up the stack are APDU's
func (s *StateMachineAccessPoint) confirmation(apdu readWriteModel.APDU, pduSource []byte) error {
	log.Debug().Msgf("confirmation\n%s", apdu)

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		log.Debug().Msg("communications enabled")
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		switch {
		case apdu.GetApduType() == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			apdu.(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_DEVICE_COMMUNICATION_CONTROL:
			log.Debug().Msg("continue with DCC request")
		case apdu.GetApduType() == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			apdu.(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_REINITIALIZE_DEVICE:
			log.Debug().Msg("continue with reinitialize device")
		case apdu.GetApduType() == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU &&
			apdu.(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS:
			log.Debug().Msg("continue with Who-Is")
		default:
			log.Debug().Msg("not a Who-Is, dropped")
			return nil
		}
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		log.Debug().Msg("initiation disabled")
	}

	switch apdu := apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		// Find duplicates of this request
		var tr *ServerSSM
		for _, serverTransactionElement := range s.serverTransactions {
			if apdu.GetInvokeId() == serverTransactionElement.invokeId && bytes.Equal(pduSource, serverTransactionElement.pduAddress) {
				tr = serverTransactionElement
				break
			}
		}
		if tr == nil {
			// build a server transaction
			var err error
			tr, err = NewServerSSM(s, pduSource)
			if err != nil {
				return errors.Wrap(err, "Error building server ssm")
			}
			s.serverTransactions = append(s.serverTransactions, tr)
		}

		// let it run with the apdu
		if err := tr.indication(apdu); err != nil {
			return errors.Wrap(err, "error runnning indication")
		}
	case readWriteModel.APDUUnconfirmedRequestExactly:
		// deliver directly to the application
		s.SapRequest(apdu)
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUComplexAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		// find the client transaction this is acking
		var tr *ClientSSM
		for _, tr := range s.clientTransactions {
			if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && bytes.Equal(pduSource, tr.pduAddress) {
				break
			}
		}
		if tr == nil {
			// TODO: log at least
			return nil
		}

		// send the packet on to the transaction
		if err := tr.confirmation(apdu); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	case readWriteModel.APDUAbortExactly:
		// find the transaction being aborted
		if apdu.GetServer() {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && bytes.Equal(pduSource, tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.confirmation(apdu); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if apdu.GetOriginalInvokeId() == serverTransactionElement.invokeId && bytes.Equal(pduSource, serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.indication(apdu); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	case readWriteModel.APDUSegmentAckExactly:
		// find the transaction being aborted
		if apdu.GetServer() {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && bytes.Equal(pduSource, tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.confirmation(apdu); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if apdu.GetOriginalInvokeId() == serverTransactionElement.invokeId && bytes.Equal(pduSource, serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.indication(apdu); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}
	return nil
}

// sapIndication This function is called when the application is requesting a new transaction as a client.
func (s *StateMachineAccessPoint) sapIndication(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("sapIndication\n%s", apdu)

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		log.Debug().Msg("communications enabled")
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		log.Debug().Msg("communications disabled")
		return nil
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		log.Debug().Msg("initiation disabled")
		if apdu.GetApduType() == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU && apdu.(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_I_AM {
			log.Debug().Msg("continue with I-Am")
		} else {
			log.Debug().Msg("not an I-Am")
			return nil
		}
	}

	switch apdu := apdu.(type) {
	case readWriteModel.APDUUnconfirmedRequestExactly:
		// deliver to the device
		s.request(apdu)
	case readWriteModel.APDUConfirmedRequestExactly:
		// make sure it has an invoke ID
		// TODO: here it is getting slightly different: usually we give the invoke id from the outside as it is build already. So maybe we need to adjust that (we never create it, we need to check for collisions but maybe we should change that so we move the creation down here)
		// s.getNextInvokeId()...
		for _, tr := range s.clientTransactions {
			if apdu.GetInvokeId() == tr.invokeId && bytes.Equal(pduDestination, tr.pduAddress) {
				return errors.New("invoke ID in use")
			}
		}

		// warning for bogus requests
		// TODO: not sure if we have that or if it is relvant (localstationaddr)

		// create a client transaction state machine
		tr, err := NewClientSSM(s, pduDestination)
		if err != nil {
			return errors.Wrap(err, "error creating client ssm")
		}

		// add it to our transactions to track it
		s.clientTransactions = append(s.clientTransactions, tr)

		// let it run
		if err := tr.indication(apdu); err != nil {
			return errors.Wrap(err, "error doing indication")
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}

	return nil
}

// sapConfirmation This function is called when the application is responding to a request, the apdu may be a simple
//        ack, complex ack, error, reject or abort
func (s *StateMachineAccessPoint) sapConfirmation(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("sapConfirmation\n%s", apdu)
	switch apdu.(type) {
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUComplexAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		// find the client transaction this is acking
		var tr *ServerSSM
		for _, tr := range s.serverTransactions {
			if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && bytes.Equal(pduDestination, tr.pduAddress) {
				break
			}
		}
		if tr == nil {
			// TODO: log at least
			return nil
		}

		// pass control to the transaction
		if err := tr.confirmation(apdu); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}
	return nil
}
