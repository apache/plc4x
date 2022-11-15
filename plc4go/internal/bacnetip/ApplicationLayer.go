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
	"github.com/apache/plc4x/plc4go/internal/bacnetip/local"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type SSMState uint8

const (
	SSMState_IDLE SSMState = iota
	SSMState_SEGMENTED_REQUEST
	SSMState_AWAIT_CONFIRMATION
	SSMState_AWAIT_RESPONSE
	SSMState_SEGMENTED_RESPONSE
	SSMState_SEGMENTED_CONFIRMATION
	SSMState_COMPLETED
	SSMState_ABORTED
)

func (s SSMState) String() string {
	switch s {
	case SSMState_IDLE:
		return "IDLE"
	case SSMState_SEGMENTED_REQUEST:
		return "SEGMENTED_REQUEST"
	case SSMState_AWAIT_CONFIRMATION:
		return "AWAIT_CONFIRMATION"
	case SSMState_AWAIT_RESPONSE:
		return "AWAIT_RESPONSE"
	case SSMState_SEGMENTED_RESPONSE:
		return "SEGMENTED_RESPONSE"
	case SSMState_SEGMENTED_CONFIRMATION:
		return "SEGMENTED_CONFIRMATION"
	case SSMState_COMPLETED:
		return "COMPLETED"
	case SSMState_ABORTED:
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

type SSMSAPRequirements interface {
	_ServiceAccessPoint
	_Client
	GetDeviceInfoCache() *DeviceInfoCache
	GetLocalDevice() *local.LocalDeviceObject
	GetProposedWindowSize() uint8
	GetClientTransactions() []*ClientSSM
	GetServerTransactions() []*ServerSSM
	GetApplicationTimeout() uint
}

// SSM - Segmentation State Machine
type SSM struct {
	OneShotTask

	ssmSAP SSMSAPRequirements

	pduAddress []byte
	deviceInfo *DeviceInfo

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

func NewSSM(sap SSMSAPRequirements, pduAddress []byte) (SSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pdu_address", pduAddress).Msg("init")
	var deviceInfo *DeviceInfo
	deviceInfoTemp, ok := sap.GetDeviceInfoCache().GetDeviceInfo(DeviceInfoCacheKey{PduSource: pduAddress})
	if ok {
		deviceInfo = &deviceInfoTemp
	}
	localDevice := sap.GetLocalDevice()
	return SSM{
		ssmSAP:                sap,
		pduAddress:            pduAddress,
		deviceInfo:            deviceInfo,
		state:                 SSMState_IDLE,
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
	if s.state == SSMState_COMPLETED || s.state == SSMState_ABORTED {
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
		serializedBytes, err := apdu.GetServiceRequest().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		segmentAPDU := segmentAPDU{
			originalApdu:     apdu,
			originalInvokeId: apdu.GetInvokeId(),
			serviceBytes:     serializedBytes,
			serviceChoice:    apdu.GetServiceRequest().GetServiceChoice(),
		}
		s.segmentAPDU = &segmentAPDU
	case readWriteModel.APDUComplexAckExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		serializedBytes, err := apdu.GetServiceAck().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		segmentAPDU := segmentAPDU{
			originalApdu:  apdu,
			serviceBytes:  serializedBytes,
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
		serializedBytes, err := apdu.GetServiceRequest().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		s.segmentAPDU.serviceBytes = append(s.segmentAPDU.serviceBytes, serializedBytes...)
	case readWriteModel.APDUComplexAckExactly:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		serializedBytes, err := apdu.GetServiceAck().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		s.segmentAPDU.serviceBytes = append(s.segmentAPDU.serviceBytes, serializedBytes...)
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
		if err := s.ssmSAP.Request(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending request")
		}
		if moreFollows {
			s.sentAllSegments = true
		}
	}
	return nil
}

type ClientSSM struct {
	SSM
}

func NewClientSSM(sap SSMSAPRequirements, pduAddress []byte) (*ClientSSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pduAddress", pduAddress).Msg("init")
	ssm, err := NewSSM(sap, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if ssm.deviceInfo == nil {
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

	if s.state == SSMState_COMPLETED || s.state == SSMState_ABORTED {
		log.Debug().Msg("remove from active transaction")
		s.ssmSAP.GetClientTransactions() // TODO remove "this" transaction from the list
		if s.deviceInfo == nil {
			// TODO: release device entry
			log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// Request This function is called by client transaction functions when it wants to send a message to the device
func (s *ClientSSM) Request(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("request\n%s", apdu)
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	return s.ssmSAP.Request(apdu)
}

// Indication This function is called after the device has bound a new transaction and wants to start the process
//        rolling
func (s *ClientSSM) Indication(apdu readWriteModel.APDU) error { // TODO: maybe use another name for that
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
	if s.deviceInfo == nil || s.deviceInfo.MaximumApduLengthAccepted != nil {
		s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
	} else if s.deviceInfo.MaximumNpduLength == nil {
		//      if the max npdu length of the server isn't known, assume that it is the same as the max apdu length accepted
		s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
	} else {
		s.segmentSize = utils.Min(*s.deviceInfo.MaximumNpduLength, uint(s.maxApduLengthAccepted.NumberOfOctets()))
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
			return s.Response(abort)
		}

		if s.deviceInfo == nil {
			log.Debug().Msg("no server info for segmentation support")
		} else if *s.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && *s.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			log.Debug().Msg("server can't receive segmented requests")
			abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return s.Response(abort)
		}

		// make sure we don't exceed the number of segments in our request that the server said it was willing to accept
		if s.deviceInfo == nil {
			log.Debug().Msg("no server info for maximum number of segments")
		} else if s.deviceInfo.MaxSegmentsAccepted == nil {
			log.Debug().Msgf("server doesn't say maximum number of segments")
		} else if s.segmentCount > s.deviceInfo.MaxSegmentsAccepted.MaxSegments() {
			log.Debug().Msg("server can't receive enough segments")
			abort, err := s.abort(readWriteModel.BACnetAbortReason_APDU_TOO_LONG)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return s.Response(abort)
		}
	}

	// send out the first segment (or the whole thing)
	if s.segmentCount == 1 {
		// unsegmented
		s.sentAllSegments = true
		s.retryCount = 0
		if err := s.setState(SSMState_AWAIT_CONFIRMATION, &s.apduTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	} else {
		// segmented
		s.sentAllSegments = false
		s.retryCount = 0
		s.segmentRetryCount = 0
		s.initialSequenceNumber = 0
		s.actualWindowSize = nil
		if err := s.setState(SSMState_SEGMENTED_REQUEST, &s.segmentTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	}

	// deliver to the device
	segment, _, err := s.getSegment(0)
	if err != nil {
		return errors.Wrap(err, "error getting segment")
	}
	return s.Request(segment)
}

// Response This function is called by client transaction functions when they want to send a message to the application.
func (s *ClientSSM) Response(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("response\n%s", apdu)
	// make sure it has a good source and destination
	// TODO: check if source == s.pduAddress
	// TODO: check if

	// send it to the application
	return s.ssmSAP.SapResponse(apdu)
}

// Confirmation This function is called by the device for all upstream messages related to the transaction.
func (s *ClientSSM) Confirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("confirmation\n%s", apdu)

	switch s.state {
	case SSMState_SEGMENTED_REQUEST:
		return s.segmentedRequest(apdu)
	case SSMState_AWAIT_CONFIRMATION:
		return s.awaitConfirmation(apdu)
	case SSMState_SEGMENTED_CONFIRMATION:
		return s.segmentedConfirmation(apdu)
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// processTask This function is called when something has taken too long
func (s *ClientSSM) processTask() error {
	log.Debug().Msg("processTask")
	switch s.state {
	case SSMState_SEGMENTED_REQUEST:
		return s.segmentedRequestTimeout()
	case SSMState_AWAIT_CONFIRMATION:
		return s.awaitConfirmationTimeout()
	case SSMState_SEGMENTED_CONFIRMATION:
		return s.segmentedConfirmationTimeout()
	case SSMState_COMPLETED, SSMState_ABORTED:
		return nil
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// abort This function is called when the transaction should be aborted
func (s *ClientSSM) abort(reason readWriteModel.BACnetAbortReason) (readWriteModel.APDU, error) {
	log.Debug().Msgf("abort\n%s", reason)

	// change the state to aborted
	if err := s.setState(SSMState_ABORTED, nil); err != nil {
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

			if err := s.setState(SSMState_AWAIT_CONFIRMATION, &s.apduTimeout); err != nil {
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
			if err := s.Request(abort); err != nil { // send it ot the device
				log.Debug().Err(err).Msg("error sending request")
			}
			if err := s.Response(abort); err != nil { // send it ot the application
				log.Debug().Err(err).Msg("error sending response")
			}
		} else {
			if err := s.setState(SSMState_COMPLETED, nil); err != nil {
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
			if err := s.Request(abort); err != nil { // send it ot the device
				log.Debug().Err(err).Msg("error sending request")
			}
			if err := s.Response(abort); err != nil { // send it ot the application
				log.Debug().Err(err).Msg("error sending response")
			}
		} else if !apdu.GetSegmentedMessage() {
			// ack is not segmented
			if err := s.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := s.Response(apdu); err != nil {
				log.Debug().Err(err).Msg("error sending response")
			}
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
			if err := s.setState(SSMState_SEGMENTED_CONFIRMATION, &s.segmentTimeout); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	case readWriteModel.APDUErrorExactly:
		log.Debug().Msg("error/reject/abort")
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := s.Response(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
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
			if err := s.Request(apdu); err != nil {
				log.Debug().Err(err).Msg("error sending request")
			}
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
		if err := s.Response(abort); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
	}
	return nil
}

func (s *ClientSSM) awaitConfirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("awaitConfirmation\n%s", apdu)

	switch apdu := apdu.(type) {
	case readWriteModel.APDUAbortExactly:
		log.Debug().Msg("Server aborted")

		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := s.Response(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		log.Debug().Msg("simple ack, error or reject")

		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := s.Response(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUComplexAckExactly:
		log.Debug().Msg("complex ack")

		if !apdu.GetSegmentedMessage() {
			log.Debug().Msg("unsegmented")

			if err := s.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := s.Response(apdu); err != nil {
				log.Debug().Err(err).Msg("error sending response")
			}
		} else if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			log.Debug().Msg("local device can't receive segmented messages")

			abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := s.Response(abort); err != nil {
				log.Debug().Err(err).Msg("error sending response")
			}
		} else if *apdu.GetSequenceNumber() == 0 {
			log.Debug().Msg("segmented response")

			// set the segmented response context
			if err := s.setSegmentationContext(apdu); err != nil {
				return errors.Wrap(err, "error set segmentation context")
			}

			s.actualWindowSize = apdu.GetProposedWindowSize()
			s.lastSequenceNumber = 0
			s.initialSequenceNumber = 0
			if err := s.setState(SSMState_SEGMENTED_CONFIRMATION, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}

			// send back a segment ack
			segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
			if err := s.Request(segmentAck); err != nil {
				log.Debug().Err(err).Msg("error sending request")
			}
		} else {
			log.Debug().Msg("Invalid apdu in this state")

			abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := s.Request(abort); err != nil { // send it ot the device
				log.Debug().Err(err).Msg("error sending request")
			}
			if err := s.Response(abort); err != nil { // send it ot the application
				log.Debug().Err(err).Msg("error sending response")
			}
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
		if err := s.Indication(s.segmentAPDU.originalApdu); err != nil { // TODO: check that it is really the intention to re-send the original apdu here
			return err
		}
		s.retryCount = saveCount
	} else {
		log.Debug().Msg("retry count exceeded")

		abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Response(abort); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
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
		if err := s.Request(abort); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(abort); err != nil { // send it ot the application
			log.Debug().Err(err).Msg("error sending response")
		}
	}

	// it must be segmented
	if !apduComplexAck.GetSegmentedMessage() {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Request(abort); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(abort); err != nil { // send it ot the application
			log.Debug().Err(err).Msg("error sending response")
		}
	}

	// proper segment number
	if *apduComplexAck.GetSequenceNumber() != s.lastSequenceNumber+1 {
		log.Debug().Msgf("segment %d received out of order, should be %d", apduComplexAck.GetSequenceNumber(), s.lastSequenceNumber+1)

		// segment received out of order
		s.restartTimer(s.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(true, false, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Request(segmentAck); err != nil {
			log.Debug().Err(err).Msg("error sending request")
		}
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
		if err := s.Request(segmentAck); err != nil {
			log.Debug().Err(err).Msg("error sending request")
		}

		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now...
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		parse, err := readWriteModel.APDUParse(s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := s.Response(parse); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
	} else if *apduComplexAck.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.restartTimer(s.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, s.invokeId, s.lastSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Request(segmentAck); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
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
	return s.Response(abort)
}

type ServerSSM struct {
	SSM
	segmentedResponseAccepted bool
}

func NewServerSSM(sap SSMSAPRequirements, pduAddress []byte) (*ServerSSM, error) {
	log.Debug().Interface("sap", sap).Bytes("pduAddress", pduAddress).Msg("init")
	ssm, err := NewSSM(sap, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if &ssm.deviceInfo == nil {
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

	if s.state == SSMState_COMPLETED || s.state == SSMState_ABORTED {
		log.Debug().Msg("remove from active transaction")
		s.ssmSAP.GetServerTransactions() // TODO remove "this" transaction from the list
		if s.deviceInfo != nil {
			// TODO: release device entry
			log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// Request This function is called by transaction functions to send to the application
func (s *ServerSSM) Request(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("request\n%s", apdu)
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	return s.ssmSAP.SapRequest(apdu)
}

// Indication This function is called for each downstream packet related to
//        the transaction
func (s *ServerSSM) Indication(apdu readWriteModel.APDU) error { // TODO: maybe use another name for that
	log.Debug().Msgf("indication\n%s", apdu)
	// make sure we're getting confirmed requests

	switch s.state {
	case SSMState_IDLE:
		return s.idle(apdu)
	case SSMState_SEGMENTED_REQUEST:
		return s.segmentedRequest(apdu)
	case SSMState_AWAIT_RESPONSE:
		return s.awaitResponse(apdu)
	case SSMState_SEGMENTED_RESPONSE:
		return s.segmentedResponse(apdu)
	default:
		return errors.Errorf("invalid state %s", s.state)
	}
}

// Response This function is called by client transaction functions when they want to send a message to the application.
func (s *ServerSSM) Response(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("response\n%s", apdu)
	// make sure it has a good source and destination
	// TODO: check if source == none
	// TODO: check if destnation = s.pduAddress

	// send it via the device
	return s.ssmSAP.Request(apdu)
}

// Confirmation This function is called when the application has provided a response and needs it to be sent to the
//        client.
func (s *ServerSSM) Confirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("confirmation\n%s", apdu)

	// check to see we are in the correct state
	if s.state != SSMState_AWAIT_RESPONSE {
		log.Debug().Msg("warning: no expecting a response")
	}

	switch apdu := apdu.(type) {
	// abort response
	case readWriteModel.APDUAbortExactly:
		log.Debug().Msg("abort")

		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// end the response to the device
		return s.Response(apdu)
	// simple response
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		log.Debug().Msg("simple ack, error or reject")

		// transaction completed
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// send the response to the device
		return s.Response(apdu)
	// complex ack
	case readWriteModel.APDUComplexAckExactly:
		log.Debug().Msg("complex ack")

		// save the response and set the segmentation context
		if err := s.setSegmentationContext(apdu); err != nil {
			return errors.Wrap(err, "error settings segmentation context")
		}

		// the segment size is the minimum of the size of the largest packet that can be delivered to the client and the
		//            largest it can accept
		if s.deviceInfo == nil || s.deviceInfo.MaximumNpduLength == nil {
			s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
		} else {
			s.segmentSize = utils.Min(*s.deviceInfo.MaximumNpduLength, uint(s.maxApduLengthAccepted.NumberOfOctets()))
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
					return s.Response(abort)
				}

				// make sure client supports segmented receive
				if !s.segmentedResponseAccepted {
					log.Debug().Msg("client can't receive segmented responses")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(abort)
				}

				// make sure we don't exceed the number of segments in our response that the client said it was willing to accept
				//                in the request
				if s.maxSegmentsAccepted != nil && s.segmentCount > s.maxSegmentsAccepted.MaxSegments() {
					log.Debug().Msg("client can't receive enough segments")
					abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(abort)
				}
			}

			// initialize the state
			s.segmentRetryCount = 0
			s.initialSequenceNumber = 0
			s.actualWindowSize = nil

			// send out the first segment (or the whole thing)
			if s.segmentCount == 1 {
				if err := s.Response(apdu); err != nil {
					log.Debug().Err(err).Msg("error sending response")
				}
				if err := s.setState(SSMState_COMPLETED, nil); err != nil {
					return errors.Wrap(err, "Error setting state to aborted")
				}
			} else {
				segment, _, err := s.getSegment(0)
				if err != nil {
					return errors.Wrap(err, "error getting first segment")
				}
				if err := s.Response(segment); err != nil {
					log.Debug().Err(err).Msg("error sending response")
				}
				if err := s.setState(SSMState_SEGMENTED_RESPONSE, nil); err != nil {
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
	case SSMState_SEGMENTED_REQUEST:
		return s.segmentedRequestTimeout()
	case SSMState_AWAIT_CONFIRMATION:
		return s.awaitResponseTimeout()
	case SSMState_SEGMENTED_CONFIRMATION:
		return s.segmentedResponseTimeout()
	case SSMState_COMPLETED, SSMState_ABORTED:
		return nil
	default:
		return errors.Errorf("Invalid state %s", s.state)
	}
}

// abort This function is called when the transaction should be aborted
func (s *ServerSSM) abort(reason readWriteModel.BACnetAbortReason) (readWriteModel.APDU, error) {
	log.Debug().Msgf("abort\n%s", reason)

	// change the state to aborted
	if err := s.setState(SSMState_ABORTED, nil); err != nil {
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
	if apduConfirmedRequest.GetSegmentedResponseAccepted() && s.deviceInfo != nil {
		switch *s.deviceInfo.SegmentationSupported {
		case readWriteModel.BACnetSegmentation_NO_SEGMENTATION:
			log.Debug().Msg("client actually supports segmented receive")
			segmentedReceive := readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE
			s.deviceInfo.SegmentationSupported = &segmentedReceive

		// TODO: bacpypes updates the cache here but as we have a pointer  to the entry we should need that. Maybe we should because concurrency... lets see later
		case readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT:
			log.Debug().Msg("client actually supports both segmented transmit and receive")
			segmentedBoth := readWriteModel.BACnetSegmentation_SEGMENTED_BOTH
			s.deviceInfo.SegmentationSupported = &segmentedBoth

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
	if s.deviceInfo != nil && s.deviceInfo.MaximumApduLengthAccepted != nil {
		if *s.deviceInfo.MaximumApduLengthAccepted < *s.maxApduLengthAccepted {
			log.Debug().Msg("apdu max reponse encoding error")
		} else {
			s.maxApduLengthAccepted = s.deviceInfo.MaximumApduLengthAccepted
		}
	}
	log.Debug().Msgf("maxApduLengthAccepted %s", *s.maxApduLengthAccepted)

	// save the number of segments the client is willing to accept in the ack, if this is None then the value is unknown or more than 64
	getMaxSegmentsAccepted := apduConfirmedRequest.GetMaxSegmentsAccepted()
	s.maxSegmentsAccepted = &getMaxSegmentsAccepted

	// unsegmented request
	if len(apduConfirmedRequest.GetSegment()) <= 0 {
		if err := s.setState(SSMState_AWAIT_RESPONSE, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		return s.Request(apdu)
	}

	// make sure we support segmented requests
	if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		return s.Response(abort)
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
	if err := s.setState(SSMState_SEGMENTED_REQUEST, &s.segmentTimeout); err != nil {
		return errors.Wrap(err, "Error setting state to aborted")
	}

	// send back a segment ack
	segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
	log.Debug().Msgf("segAck: %s", segack)
	return s.Response(segack)
}

func (s *ServerSSM) segmentedRequest(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("segmentedRequest\n%s", apdu)

	// some kind of problem
	if _, ok := apdu.(readWriteModel.APDUAbortExactly); ok {
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		return s.Response(apdu)
	}

	// the only messages we should be getting are confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if castedApdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Request(abort); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(abort); err != nil { // send it ot the application
			log.Debug().Err(err).Msg("error sending response")
		}
	} else {
		apduConfirmedRequest = castedApdu
	}

	// it must be segmented
	if !apduConfirmedRequest.GetSegmentedMessage() {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Request(abort); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(abort); err != nil { // send it ot the application
			log.Debug().Err(err).Msg("error sending response")
		}
	}

	// proper segment number
	if *apduConfirmedRequest.GetSequenceNumber() != s.lastSequenceNumber+1 {
		log.Debug().Msgf("segment %d received out of order, should be %d", *apduConfirmedRequest.GetSequenceNumber(), s.lastSequenceNumber+1)

		// segment received out of order
		s.restartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(true, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		return s.Response(segack)
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
		if err := s.Response(segack); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}

		// forward the whole thing to the application
		applicationTimeout := s.ssmSAP.GetApplicationTimeout()
		if err := s.setState(SSMState_AWAIT_RESPONSE, &applicationTimeout); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		// TODO: here we need to rebuild again yada yada
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now..
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		parse, err := readWriteModel.APDUParse(s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := s.Request(parse); err != nil {
			log.Debug().Err(err).Msg("error sending request")
		}
	} else if *apduConfirmedRequest.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.restartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Response(segack); err != nil {
			log.Debug().Err(err).Msg("error sending response")
		}
	} else {
		// wait for more segments
		s.restartTimer(s.segmentTimeout)
	}

	return nil
}

func (s *ServerSSM) segmentedRequestTimeout() error {
	log.Debug().Msg("segmentedRequestTimeout")

	// give up
	if err := s.setState(SSMState_ABORTED, nil); err != nil {
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
		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Request(apdu); err != nil { // send it ot the device
			log.Debug().Err(err).Msg("error sending request")
		}
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
	if err := s.Request(abort); err != nil {
		log.Debug().Err(err).Msg("error sending request")
	}
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
			if err := s.setState(SSMState_COMPLETED, nil); err != nil {
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
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Response(apdu); err != nil { // send it ot the application
			log.Debug().Err(err).Msg("error sending response")
		}
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
		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
	}
	return nil
}

type StateMachineAccessPoint struct {
	*Client
	*ServiceAccessPoint

	localDevice           *local.LocalDeviceObject
	deviceInventory       *DeviceInfoCache
	nextInvokeId          uint8
	clientTransactions    []*ClientSSM
	serverTransactions    []*ServerSSM
	numberOfApduRetries   int
	apduTimeout           int
	maxApduLengthAccepted int
	segmentationSupported readWriteModel.BACnetSegmentation
	segmentTimeout        int
	maxSegmentsAccepted   int
	proposedWindowSize    uint8
	dccEnableDisable      readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable
	applicationTimeout    uint
}

func NewStateMachineAccessPoint(localDevice *local.LocalDeviceObject, deviceInventory *DeviceInfoCache, sapID *int, cid *int) (*StateMachineAccessPoint, error) {
	log.Debug().Msgf("NewStateMachineAccessPoint localDevice=%v deviceInventory=%v sap=%v cid=%v", localDevice, deviceInventory, sapID, cid)

	s := &StateMachineAccessPoint{
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
	// basic initialization
	client, err := NewClient(cid, s)
	if err != nil {
		return nil, errors.Wrapf(err, "error building client for %d", cid)
	}
	s.Client = client
	serviceAccessPoint, err := NewServiceAccessPoint(sapID, s)
	if err != nil {
		return nil, errors.Wrapf(err, "error building serviceAccessPoint for %d", sapID)
	}
	s.ServiceAccessPoint = serviceAccessPoint
	return s, nil
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

// ConfirmationFromSource Packets coming up the stack are APDU's
func (s *StateMachineAccessPoint) ConfirmationFromSource(apdu readWriteModel.APDU, pduSource []byte) error { // TODO: note we need a special method here as we don't contain src in the apdu
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
		if err := tr.Indication(apdu); err != nil {
			return errors.Wrap(err, "error runnning indication")
		}
	case readWriteModel.APDUUnconfirmedRequestExactly:
		// deliver directly to the application
		if err := s.SapRequest(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending request")
		}
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
		if err := tr.Confirmation(apdu); err != nil {
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
			if err := tr.Confirmation(apdu); err != nil {
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
			if err := tr.Indication(apdu); err != nil {
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
			if err := tr.Confirmation(apdu); err != nil {
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
			if err := tr.Indication(apdu); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}
	return nil
}

// SapIndication This function is called when the application is requesting a new transaction as a client.
func (s *StateMachineAccessPoint) SapIndication(apdu readWriteModel.APDU, pduDestination []byte) error {
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
		if err := s.Request(apdu); err != nil {
			log.Debug().Err(err).Msg("error sending the request")
		}
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
		if err := tr.Indication(apdu); err != nil {
			return errors.Wrap(err, "error doing indication")
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}

	return nil
}

// SapConfirmation This function is called when the application is responding to a request, the apdu may be a simple
//        ack, complex ack, error, reject or abort
func (s *StateMachineAccessPoint) SapConfirmation(apdu readWriteModel.APDU, pduDestination []byte) error {
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
		if err := tr.Confirmation(apdu); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	default:
		return errors.Errorf("invalid APDU %T", apdu)
	}
	return nil
}

func (s *StateMachineAccessPoint) GetDeviceInfoCache() *DeviceInfoCache {
	return s.deviceInventory
}

func (s *StateMachineAccessPoint) GetLocalDevice() *local.LocalDeviceObject {
	return s.localDevice
}

func (s *StateMachineAccessPoint) GetProposedWindowSize() uint8 {
	return s.proposedWindowSize
}

func (s *StateMachineAccessPoint) GetClientTransactions() []*ClientSSM {
	return s.clientTransactions
}

func (s *StateMachineAccessPoint) GetServerTransactions() []*ServerSSM {
	return s.serverTransactions
}

func (s *StateMachineAccessPoint) GetApplicationTimeout() uint {
	return s.applicationTimeout
}

type ApplicationServiceAccessPoint struct {
	*ApplicationServiceElement
	*ServiceAccessPoint
}

func NewApplicationServiceAccessPoint(aseID *int, sapID *int) (*ApplicationServiceAccessPoint, error) {
	a := &ApplicationServiceAccessPoint{}
	applicationServiceElement, err := NewApplicationServiceElement(aseID, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}
	a.ApplicationServiceElement = applicationServiceElement
	serviceAccessPoint, err := NewServiceAccessPoint(sapID, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}
	a.ServiceAccessPoint = serviceAccessPoint
	return a, nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Indication(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("Indication\n%s", apdu)

	switch apdu := apdu.(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetConfirmedServiceChoiceKnows(uint8(apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(apdu)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(apdu); err != nil {
				return err
			}
		} else {
			log.Debug().Err(errorFound).Msg("got error")

			// TODO: map it to a error... code temporary placeholder
			return a.Response(readWriteModel.NewAPDUReject(apdu.GetInvokeId(), nil, 0))
		}
	case readWriteModel.APDUUnconfirmedRequestExactly:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetUnconfirmedServiceChoiceKnows(uint8(apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(apdu)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(apdu); err != nil {
				return err
			}
		} else {
			log.Debug().Err(errorFound).Msg("got error")
		}

	default:
		return errors.Errorf("unknown PDU type %T", apdu)
	}
	return nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapIndication(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("SapIndication\n%s", apdu)

	// TODO: check if we need to check apdu here

	return a.Request(apdu)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Confirmation(apdu readWriteModel.APDU) error {
	log.Debug().Msgf("Confirmation\n%s", apdu)

	// TODO: check if we need to check apdu here

	return a.SapResponse(apdu)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapConfirmation(apdu readWriteModel.APDU, pduDestination []byte) error {
	log.Debug().Msgf("SapConfirmation\n%s", apdu)

	// TODO: check if we need to check apdu here

	return a.Response(apdu)
}
