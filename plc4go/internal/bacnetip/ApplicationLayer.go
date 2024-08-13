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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"time"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/pkg/errors"
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
	GetLocalDevice() *LocalDeviceObject
	GetProposedWindowSize() uint8
	GetClientTransactions() []*ClientSSM
	RemoveClientTransaction(*ClientSSM)
	GetServerTransactions() []*ServerSSM
	RemoveServerTransaction(*ServerSSM)
	GetApplicationTimeout() uint
	GetDefaultAPDUTimeout() uint
	GetDefaultSegmentationSupported() readWriteModel.BACnetSegmentation
	GetDefaultAPDUSegmentTimeout() uint
	GetDefaultMaxSegmentsAccepted() readWriteModel.MaxSegmentsAccepted
	GetDefaultMaximumApduLengthAccepted() readWriteModel.MaxApduLengthAccepted
}

type SSMProcessingRequirements interface {
	processTask() error
}

// SSM - Segmentation State Machine
type SSM struct {
	*OneShotTask
	SSMProcessingRequirements

	ssmSAP SSMSAPRequirements

	pduAddress *Address
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
	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted

	log zerolog.Logger
}

func NewSSM(localLog zerolog.Logger, sap interface {
	SSMSAPRequirements
	SSMProcessingRequirements
}, pduAddress *Address) (*SSM, error) {
	localLog.Debug().Interface("sap", sap).Interface("pdu_address", pduAddress).Msg("init")
	var deviceInfo *DeviceInfo
	deviceInfoTemp, ok := sap.GetDeviceInfoCache().GetDeviceInfo(DeviceInfoCacheKey{PduSource: pduAddress})
	if ok {
		deviceInfo = &deviceInfoTemp
	}
	localDevice := sap.GetLocalDevice()
	var numberOfApduRetries uint
	if localDevice.NumberOfAPDURetries != nil {
		numberOfApduRetries = *localDevice.NumberOfAPDURetries
	}
	var apduTimeout uint
	if localDevice.APDUTimeout != nil {
		apduTimeout = *localDevice.APDUTimeout
	} else {
		apduTimeout = sap.GetDefaultAPDUTimeout()
	}
	var segmentationSupported readWriteModel.BACnetSegmentation
	if localDevice.SegmentationSupported != nil {
		segmentationSupported = *localDevice.SegmentationSupported
	} else {
		segmentationSupported = sap.GetDefaultSegmentationSupported()
	}
	var segmentTimeout uint
	if localDevice.APDUSegmentTimeout != nil {
		segmentTimeout = *localDevice.APDUSegmentTimeout
	} else {
		segmentTimeout = sap.GetDefaultAPDUSegmentTimeout()
	}
	var maxSegmentsAccepted readWriteModel.MaxSegmentsAccepted
	if localDevice.MaxSegmentsAccepted != nil {
		maxSegmentsAccepted = *localDevice.MaxSegmentsAccepted
	} else {
		maxSegmentsAccepted = sap.GetDefaultMaxSegmentsAccepted()
	}
	var maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted
	if localDevice.MaximumApduLengthAccepted != nil {
		maxApduLengthAccepted = *localDevice.MaximumApduLengthAccepted
	} else {
		maxApduLengthAccepted = sap.GetDefaultMaximumApduLengthAccepted()
	}
	ssm := &SSM{
		ssmSAP:                sap,
		pduAddress:            pduAddress,
		deviceInfo:            deviceInfo,
		state:                 SSMState_IDLE,
		numberOfApduRetries:   numberOfApduRetries,
		apduTimeout:           apduTimeout,
		segmentationSupported: segmentationSupported,
		segmentTimeout:        segmentTimeout,
		maxSegmentsAccepted:   maxSegmentsAccepted,
		maxApduLengthAccepted: maxApduLengthAccepted,
		log:                   localLog,
	}
	ssm.OneShotTask = NewOneShotTask(ssm, nil)
	ssm.SSMProcessingRequirements = sap
	return ssm, nil
}

func (s *SSM) StartTimer(millis uint) {
	s.log.Debug().Uint("millis", millis).Msg("Start timer")
	s.RestartTimer(millis)
}

func (s *SSM) StopTimer() {
	s.log.Debug().Msg("Stop Timer")
	if s.isScheduled {
		s.log.Debug().Msg("is scheduled")
		s.SuspendTask()
	}
}

func (s *SSM) RestartTimer(millis uint) {
	s.log.Debug().Uint("millis", millis).Msg("restart timer")
	if s.isScheduled {
		s.log.Debug().Msg("is scheduled")
		s.SuspendTask()
	}

	delta := time.Millisecond * time.Duration(millis)
	s.InstallTask(InstallTaskOptions{Delta: &delta})
}

// setState This function is called when the derived class wants to change state
func (s *SSM) setState(newState SSMState, timer *uint) error {
	s.log.Debug().Stringer("state", newState).Interface("timer", timer).Msg("setState")
	if s.state == SSMState_COMPLETED || s.state == SSMState_ABORTED {
		return errors.Errorf("Invalid state transition from %s to %s", s.state, newState)
	}

	s.StopTimer()

	s.state = newState

	if timer != nil {
		s.StartTimer(*timer)
	}
	return nil
}

// setSegmentationContext This function is called to set the segmentation context
func (s *SSM) setSegmentationContext(apdu readWriteModel.APDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("setSegmentationContext")
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
//
//	The segmentAPDU is the context
func (s *SSM) getSegment(index uint8) (segmentAPDU PDU, moreFollows bool, err error) {
	s.log.Debug().Uint8("index", index).Msg("Get segment")
	if s.segmentAPDU == nil {
		return nil, false, errors.New("No segment apdu set")
	}

	if index > s.segmentCount {
		return nil, false, errors.Errorf("Invalid segment number %d, APDU has %d segments", index, s.segmentCount)
	}

	// TODO: the original code does here something funky but it seems it is best to just return the original apdu
	if s.segmentCount == 1 {
		return NewPDU(s.segmentAPDU.originalApdu, WithPDUDestination(s.pduAddress)), false, nil
	}

	moreFollows = index < s.segmentCount-1
	sequenceNumber := index % 255
	proposedWindowSize := s.actualWindowSize
	if index == 0 {
		getProposedWindowSize := s.ssmSAP.GetProposedWindowSize()
		proposedWindowSize = &getProposedWindowSize
	}
	s.log.Debug().Interface("proposedWindowSize", proposedWindowSize).Msg("working with proposedWindowSize")
	serviceChoice := &s.segmentAPDU.serviceChoice
	offset := uint(index) * s.segmentSize
	segmentBytes := s.segmentAPDU.serviceBytes[offset : offset+s.segmentSize]
	if !s.segmentAPDU.isAck {
		s.log.Debug().Msg("confirmed request context")
		segmentedResponseAccepted := s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE || s.segmentationSupported == readWriteModel.BACnetSegmentation_SEGMENTED_BOTH
		s.log.Debug().Bool("segmentedResponseAccepted", segmentedResponseAccepted).Msg("segmentedResponseAccepted")
		segmentAPDU = NewPDU(readWriteModel.NewAPDUConfirmedRequest(
			true,
			moreFollows,
			segmentedResponseAccepted,
			s.maxSegmentsAccepted,
			s.maxApduLengthAccepted,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		), WithPDUDestination(s.pduAddress))
	} else {
		s.log.Debug().Msg("complex ack context")
		segmentAPDU = NewPDU(readWriteModel.NewAPDUComplexAck(
			true,
			moreFollows,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		), WithPDUDestination(s.pduAddress))
	}
	return segmentAPDU, moreFollows, nil
}

// TODO: check that function. looks a bit wonky to just append the payloads like that
// appendSegment This function appends the apdu content to the end of the current APDU being built.  The segmentAPDU is
//
//	the context
func (s *SSM) appendSegment(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("appendSegment")
	switch apdu := apdu.GetMessage().(type) {
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
	s.log.Debug().Uint8("sequenceA", sequenceA).Uint8("sequenceB", sequenceB).Msg("inWindow %d-%d")
	return (uint(sequenceA)-uint(sequenceB)-256)%256 < uint(*s.actualWindowSize)
}

func (s *SSM) fillWindow(sequenceNumber uint8) error {
	s.log.Debug().Uint8("sequenceNumber", sequenceNumber).Msg("fillWindow")
	for i := uint8(0); i < *s.actualWindowSize; i++ {
		apdu, moreFollows, err := s.getSegment(sequenceNumber + i)
		if err != nil {
			return errors.Wrapf(err, "Error sending out segment %d", i)
		}
		if err := s.ssmSAP.Request(NewArgs(NewPDU(apdu.GetMessage(), WithPDUDestination(s.pduAddress))), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if moreFollows {
			s.sentAllSegments = true
		}
	}
	return nil
}

type ClientSSM struct {
	*SSM

	passLogToModel bool
	log            zerolog.Logger
}

func NewClientSSM(localLog zerolog.Logger, sap SSMSAPRequirements, pduAddress *Address) (*ClientSSM, error) {
	localLog.Debug().Interface("sap", sap).Interface("pduAddress", pduAddress).Msg("init")
	c := &ClientSSM{
		log: localLog,
	}
	ssm, err := NewSSM(localLog, struct {
		SSMSAPRequirements
		SSMProcessingRequirements
	}{sap, c}, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if ssm.deviceInfo == nil {
		// TODO: get entry for device, store it in inventory
		localLog.Debug().Msg("Accquire device information")
	}
	c.SSM = ssm
	return c, nil
}

// setState This function is called when the client wants to change state
func (c *ClientSSM) setState(newState SSMState, timer *uint) error {
	c.log.Debug().Stringer("state", newState).Interface("timer", timer).Msg("setState")
	// do the regular state change
	if err := c.SSM.setState(newState, timer); err != nil {
		return errors.Wrap(err, "error during SSM state transition")
	}

	if c.state == SSMState_COMPLETED || c.state == SSMState_ABORTED {
		c.log.Debug().Msg("remove from active transaction")
		c.ssmSAP.RemoveClientTransaction(c)
		if c.deviceInfo == nil {
			// TODO: release device entry
			c.log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// Request This function is called by client transaction functions when it wants to send a message to the device
func (c *ClientSSM) Request(args Args, kwargs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("request")
	apdu := args.Get0PDU()

	// make sure it has a good source and destination
	apdu = NewPDUFromPDU(apdu, WithPDUSource(nil), WithPDUDestination(c.pduAddress))

	// send it via the device
	return c.ssmSAP.Request(NewArgs(apdu), kwargs)
}

// Indication This function is called after the device has bound a new transaction and wants to start the process
//
//	rolling
func (c *ClientSSM) Indication(args Args, kwargs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("indication")
	apdu := args.Get0PDU()
	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apduCasted, ok := apdu.GetMessage().(readWriteModel.APDUConfirmedRequestExactly); !ok {
		return errors.Errorf("Invalid APDU type %T", apduCasted)
	} else {
		apduConfirmedRequest = apduCasted
	}

	// save the request and set the segmentation context
	if err := c.setSegmentationContext(apduConfirmedRequest); err != nil {
		return errors.Wrap(err, "error setting context")
	}

	// if the max apdu length of the server isn't known, assume that it is the same size as our own and will be the segment
	//        size
	if c.deviceInfo == nil || c.deviceInfo.MaximumApduLengthAccepted != nil {
		c.segmentSize = uint(c.maxApduLengthAccepted.NumberOfOctets())
	} else if c.deviceInfo.MaximumNpduLength == nil {
		//      if the max npdu length of the server isn't known, assume that it is the same as the max apdu length accepted
		c.segmentSize = uint(c.maxApduLengthAccepted.NumberOfOctets())
	} else {
		c.segmentSize = min(*c.deviceInfo.MaximumNpduLength, uint(c.maxApduLengthAccepted.NumberOfOctets()))
	}
	c.log.Debug().Uint("segmentSize", c.segmentSize).Msg("segment size")

	c.invokeId = apduConfirmedRequest.GetInvokeId()
	c.log.Debug().Uint8("invokeId", c.invokeId).Msg("invoke ID")

	var segmentCount, more int
	segmentCount, more = len(c.segmentAPDU.serviceBytes)/int(c.segmentSize), len(c.segmentAPDU.serviceBytes)%int(c.segmentSize)
	c.segmentCount = uint8(segmentCount)
	if more > 0 {
		c.segmentCount += 1
	}
	c.log.Debug().Int("segmentCount", segmentCount).Msg("segment count")

	if c.segmentCount > 1 {
		if c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			c.log.Debug().Msg("local device can't send segmented requests")
			abort, err := c.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return c.Response(NewArgs(abort), kwargs)
		}

		if c.deviceInfo == nil {
			c.log.Debug().Msg("no server info for segmentation support")
		} else if *c.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && *c.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			c.log.Debug().Msg("server can't receive segmented requests")
			abort, err := c.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return c.Response(NewArgs(abort), kwargs)
		}

		// make sure we don't exceed the number of segments in our request that the server said it was willing to accept
		if c.deviceInfo == nil {
			c.log.Debug().Msg("no server info for maximum number of segments")
		} else if c.deviceInfo.MaxSegmentsAccepted == nil {
			c.log.Debug().Msg("server doesn't say maximum number of segments")
		} else if c.segmentCount > c.deviceInfo.MaxSegmentsAccepted.MaxSegments() {
			c.log.Debug().Msg("server can't receive enough segments")
			abort, err := c.abort(readWriteModel.BACnetAbortReason_APDU_TOO_LONG)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return c.Response(NewArgs(abort), kwargs)
		}
	}

	// send out the first segment (or the whole thing)
	if c.segmentCount == 1 {
		// unsegmented
		c.sentAllSegments = true
		c.retryCount = 0
		if err := c.setState(SSMState_AWAIT_CONFIRMATION, &c.apduTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	} else {
		// segmented
		c.sentAllSegments = false
		c.retryCount = 0
		c.segmentRetryCount = 0
		c.initialSequenceNumber = 0
		c.actualWindowSize = nil
		if err := c.setState(SSMState_SEGMENTED_REQUEST, &c.segmentTimeout); err != nil {
			return errors.Wrap(err, "error switching state")
		}
	}

	// deliver to the device
	segment, _, err := c.getSegment(0)
	if err != nil {
		return errors.Wrap(err, "error getting segment")
	}
	return c.Request(NewArgs(segment), kwargs)
}

// Response This function is called by client transaction functions when they want to send a message to the application.
func (c *ClientSSM) Response(args Args, kwargs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("response")

	apdu := args.Get0PDU()

	// make sure it has a good source and destination
	apdu = NewPDUFromPDU(apdu, WithPDUSource(c.pduAddress), WithPDUDestination(nil))

	// send it to the application
	return c.ssmSAP.SapResponse(NewArgs(apdu), kwargs)
}

// Confirmation This function is called by the device for all upstream messages related to the transaction.
func (c *ClientSSM) Confirmation(args Args, kwargs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("confirmation")
	apdu := args.Get0PDU()

	switch c.state {
	case SSMState_SEGMENTED_REQUEST:
		return c.segmentedRequest(apdu)
	case SSMState_AWAIT_CONFIRMATION:
		return c.awaitConfirmation(apdu)
	case SSMState_SEGMENTED_CONFIRMATION:
		return c.segmentedConfirmation(apdu)
	default:
		return errors.Errorf("Invalid state %c", c.state)
	}
}

// ProcessTask This function is called when something has taken too long
func (c *ClientSSM) ProcessTask() error {
	c.log.Debug().Stringer("currentState", c.state).Msg("ProcessTask")
	switch c.state {
	case SSMState_SEGMENTED_REQUEST:
		return c.segmentedRequestTimeout()
	case SSMState_AWAIT_CONFIRMATION:
		return c.awaitConfirmationTimeout()
	case SSMState_SEGMENTED_CONFIRMATION:
		return c.segmentedConfirmationTimeout()
	case SSMState_COMPLETED, SSMState_ABORTED:
		return nil
	default:
		return errors.Errorf("Invalid state %c", c.state)
	}
}

// abort This function is called when the transaction should be aborted
func (c *ClientSSM) abort(reason readWriteModel.BACnetAbortReason) (PDU, error) {
	c.log.Debug().Stringer("reason", reason).Msg("abort")

	// change the state to aborted
	if err := c.setState(SSMState_ABORTED, nil); err != nil {
		return nil, errors.Wrap(err, "Error setting state to aborted")
	}

	// build an abort _PDU to return
	abortApdu := readWriteModel.NewAPDUAbort(false, c.invokeId, readWriteModel.NewBACnetAbortReasonTagged(reason, uint32(reason), 0), 0)
	// return it
	return NewPDU(abortApdu), nil
}

// segmentedRequest This function is called when the client is sending a segmented request and receives an apdu
func (c *ClientSSM) segmentedRequest(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("segmentedRequest")

	switch _apdu := apdu.GetMessage().(type) {
	// server is ready for the next segment
	case readWriteModel.APDUSegmentAckExactly:
		c.log.Debug().Msg("segment ack")
		getActualWindowSize := _apdu.GetActualWindowSize()
		c.actualWindowSize = &getActualWindowSize

		// duplicate ack received?
		if !c.inWindow(_apdu.GetSequenceNumber(), c.initialSequenceNumber) {
			c.log.Debug().Msg("not in window")
			c.RestartTimer(c.segmentTimeout)
		} else if c.sentAllSegments {
			c.log.Debug().Msg("all done sending request")

			if err := c.setState(SSMState_AWAIT_CONFIRMATION, &c.apduTimeout); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		} else {
			c.log.Debug().Msg("More segments to send")

			c.initialSequenceNumber = _apdu.GetSequenceNumber() + 1
			c.retryCount = 0
			if err := c.fillWindow(c.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
			c.RestartTimer(c.segmentTimeout)
		}
	// simple ack
	case readWriteModel.APDUSimpleAckExactly:
		c.log.Debug().Msg("simple ack")

		if !c.sentAllSegments {
			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else {
			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	// complex ack
	case readWriteModel.APDUComplexAckExactly:
		c.log.Debug().Msg("complex ack")
		if !c.sentAllSegments {
			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else if !_apdu.GetSegmentedMessage() {
			// ack is not segmented
			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := c.Response(NewArgs(apdu), NoKWArgs); err != nil {
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else {
			// set the segmented response context
			if err := c.setSegmentationContext(_apdu); err != nil {
				return errors.Wrap(err, "error setting context")
			}

			// minimum of what the server is proposing and this client proposes
			minWindowSize := min(*_apdu.GetProposedWindowSize(), c.ssmSAP.GetProposedWindowSize())
			c.actualWindowSize = &minWindowSize
			c.lastSequenceNumber = 0
			c.initialSequenceNumber = 0
			if err := c.setState(SSMState_SEGMENTED_CONFIRMATION, &c.segmentTimeout); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	case readWriteModel.APDUErrorExactly:
		c.log.Debug().Msg("error/reject/abort")
		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NewArgs(apdu), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	default:
		return errors.Errorf("Invalid APDU type %T", apdu)
	}
	return nil
}

func (c *ClientSSM) segmentedRequestTimeout() error {
	c.log.Debug().Msg("segmentedRequestTimeout")

	// Try again
	if c.segmentRetryCount < c.numberOfApduRetries {
		c.log.Debug().Msg("retry segmented request")
		c.segmentRetryCount++
		c.StartTimer(c.segmentTimeout)

		if c.initialSequenceNumber == 0 {
			apdu, _, err := c.getSegment(0)
			if err != nil {
				return errors.Wrap(err, "error getting first segment")
			}
			if err := c.Request(NewArgs(apdu), NoKWArgs); err != nil {
				c.log.Debug().Err(err).Msg("error sending request")
			}
		} else {
			if err := c.fillWindow(c.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
		}
	} else {
		c.log.Debug().Msg("abort, no response from the device")

		abort, err := c.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Response(NewArgs(abort), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}
	return nil
}

func (c *ClientSSM) awaitConfirmation(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("awaitConfirmation")

	switch _apdu := apdu.GetMessage().(type) {
	case readWriteModel.APDUAbortExactly:
		c.log.Debug().Msg("Server aborted")

		if err := c.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NewArgs(apdu), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		c.log.Debug().Msg("simple ack, error or reject")

		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NewArgs(apdu), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUComplexAckExactly:
		c.log.Debug().Msg("complex ack")

		if !_apdu.GetSegmentedMessage() {
			c.log.Debug().Msg("unsegmented")

			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := c.Response(NewArgs(apdu), NoKWArgs); err != nil {
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else if c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			c.log.Debug().Msg("local device can't receive segmented messages")

			abort, err := c.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Response(NewArgs(abort), NoKWArgs); err != nil {
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else if *_apdu.GetSequenceNumber() == 0 {
			c.log.Debug().Msg("segmented response")

			// set the segmented response context
			if err := c.setSegmentationContext(_apdu); err != nil {
				return errors.Wrap(err, "error set segmentation context")
			}

			c.actualWindowSize = _apdu.GetProposedWindowSize()
			c.lastSequenceNumber = 0
			c.initialSequenceNumber = 0
			if err := c.setState(SSMState_SEGMENTED_CONFIRMATION, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}

			// send back a segment ack
			segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, c.invokeId, c.initialSequenceNumber, *c.actualWindowSize, 0)
			if err := c.Request(NewArgs(NewPDU(segmentAck)), NoKWArgs); err != nil {
				c.log.Debug().Err(err).Msg("error sending request")
			}
		} else {
			c.log.Debug().Msg("Invalid apdu in this state")

			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		}
	case readWriteModel.APDUSegmentAckExactly:
		c.log.Debug().Msg("segment ack(!?)")
		c.RestartTimer(c.segmentTimeout)
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

func (c *ClientSSM) awaitConfirmationTimeout() error {
	c.log.Debug().Msg("awaitConfirmationTimeout")

	if c.retryCount < c.numberOfApduRetries {
		c.log.Debug().
			Uint("retryCount", c.retryCount).
			Uint("numberOfApduRetries", c.numberOfApduRetries).
			Msg("no response, try again (retryCount < numberOfApduRetries)")
		c.retryCount++

		// save the retry count, indication acts like the request is coming from the application so the retryCount gets
		//            re-initialized.
		saveCount := c.retryCount
		if err := c.Indication(NewArgs(NewPDU(c.segmentAPDU.originalApdu, WithPDUDestination(c.pduAddress))), NoKWArgs); err != nil { // TODO: check that it is really the intention to re-send the original apdu here
			return err
		}
		c.retryCount = saveCount
	} else {
		c.log.Debug().
			Uint("retryCount", c.retryCount).
			Uint("numberOfApduRetries", c.numberOfApduRetries).
			Msg("retry count exceeded: retryCount >= numberOfApduRetries")

		abort, err := c.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Response(NewArgs(abort), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}
	return nil
}

func (c *ClientSSM) segmentedConfirmation(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("segmentedConfirmation")

	// the only messages we should be getting are complex acks
	apduComplexAck, ok := apdu.(readWriteModel.APDUComplexAckExactly)
	if !ok {
		c.log.Debug().Msg("complex ack required")

		abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
			c.log.Debug().Err(err).Msg("error sending request")
		}
		if err := c.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}

	// it must be segmented
	if !apduComplexAck.GetSegmentedMessage() {
		abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
			c.log.Debug().Err(err).Msg("error sending request")
		}
		if err := c.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}

	// proper segment number
	if sequenceNumber := *apduComplexAck.GetSequenceNumber(); sequenceNumber != c.lastSequenceNumber+1 {
		c.log.Debug().
			Uint8("sequenceNumber", sequenceNumber).
			Uint8("lastSequenceNumber", c.lastSequenceNumber+1).
			Msg("segment sequenceNumber received out of order, should be lastSequenceNumber")

		// segment received out of order
		c.RestartTimer(c.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(true, false, c.invokeId, c.initialSequenceNumber, *c.actualWindowSize, 0)
		if err := c.Request(NewArgs(NewPDU(segmentAck)), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending request")
		}
		return nil
	}

	// add the data
	if err := c.appendSegment(apdu); err != nil {
		return errors.Wrap(err, "error appending the segment")
	}

	// update the sequence number
	c.lastSequenceNumber = c.lastSequenceNumber + 1

	// last segment received
	if !apduComplexAck.GetMoreFollows() {
		c.log.Debug().Msg("No more follows")

		// send final ack
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, c.invokeId, c.lastSequenceNumber, *c.actualWindowSize, 0)
		if err := c.Request(NewArgs(NewPDU(segmentAck)), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending request")
		}

		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now...
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), c.log, options.WithPassLoggerToModel(c.passLogToModel))
		parse, err := readWriteModel.APDUParse(ctxForModel, c.segmentAPDU.serviceBytes, uint16(len(c.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := c.Response(NewArgs(NewPDU(parse)), NoKWArgs); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	} else if *apduComplexAck.GetSequenceNumber() == c.initialSequenceNumber+*c.actualWindowSize {
		c.log.Debug().Msg("last segment in the group")

		c.initialSequenceNumber = c.lastSequenceNumber
		c.RestartTimer(c.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, c.invokeId, c.lastSequenceNumber, *c.actualWindowSize, 0)
		if err := c.Request(NewArgs(NewPDU(segmentAck)), NoKWArgs); err != nil { // send it ot the device
			c.log.Debug().Err(err).Msg("error sending request")
		}
	} else {
		c.log.Debug().Msg("Wait for more segments")

		c.RestartTimer(c.segmentTimeout)
	}

	return nil
}

func (c *ClientSSM) segmentedConfirmationTimeout() error {
	c.log.Debug().Msg("segmentedConfirmationTimeout")

	abort, err := c.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
	if err != nil {
		return errors.Wrap(err, "error creating abort")
	}
	return c.Response(NewArgs(abort), NoKWArgs)
}

type ServerSSM struct {
	*SSM
	segmentedResponseAccepted bool

	passLogToModel bool
	log            zerolog.Logger
}

func NewServerSSM(localLog zerolog.Logger, sap SSMSAPRequirements, pduAddress *Address) (*ServerSSM, error) {
	localLog.Debug().Interface("sap", sap).Interface("pduAddress", pduAddress).Msg("init")
	s := &ServerSSM{
		segmentedResponseAccepted: true,
		log:                       localLog,
	}
	ssm, err := NewSSM(localLog, struct {
		SSMSAPRequirements
		SSMProcessingRequirements
	}{sap, s}, pduAddress)
	if err != nil {
		return nil, err
	}
	// TODO: if deviceEntry is not there get it now...
	if &ssm.deviceInfo == nil {
		// TODO: get entry for device, store it in inventory
		localLog.Debug().Msg("Accquire device information")
	}
	s.SSM = ssm
	return s, nil
}

// setState This function is called when the client wants to change state
func (s *ServerSSM) setState(newState SSMState, timer *uint) error {
	s.log.Debug().Stringer("state", newState).Interface("timer", timer).Msg("setState")
	// do the regular state change
	if err := s.SSM.setState(newState, timer); err != nil {
		return errors.Wrap(err, "error during SSM state transition")
	}

	if s.state == SSMState_COMPLETED || s.state == SSMState_ABORTED {
		s.log.Debug().Msg("remove from active transaction")
		s.ssmSAP.RemoveServerTransaction(s)
		if s.deviceInfo != nil {
			// TODO: release device entry
			s.log.Debug().Msg("release device entry")
		}
	}
	return nil
}

// Request This function is called by transaction functions to send to the application
func (s *ServerSSM) Request(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Request")
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	return s.ssmSAP.SapRequest(args, kwargs)
}

// Indication This function is called for each downstream packet related to
//
//	the transaction
func (s *ServerSSM) Indication(args Args, kwargs KWArgs) error { // TODO: maybe use another name for that
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	apdu := args.Get0PDU()
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
func (s *ServerSSM) Response(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Response")
	// make sure it has a good source and destination
	// TODO: check if source == none
	// TODO: check if destnation = s.pduAddress

	// send it via the device
	return s.ssmSAP.Request(args, kwargs)
}

// Confirmation This function is called when the application has provided a response and needs it to be sent to the
//
//	client.
func (s *ServerSSM) Confirmation(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	// check to see we are in the correct state
	if s.state != SSMState_AWAIT_RESPONSE {
		s.log.Debug().Msg("warning: no expecting a response")
	}

	apdu := args.Get0PDU()

	switch _apdu := apdu.GetMessage().(type) {
	// abort response
	case readWriteModel.APDUAbortExactly:
		s.log.Debug().Msg("abort")

		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// end the response to the device
		return s.Response(args, kwargs)
	// simple response
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		s.log.Debug().Msg("simple ack, error or reject")

		// transaction completed
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// send the response to the device
		return s.Response(args, kwargs)
	// complex ack
	case readWriteModel.APDUComplexAckExactly:
		s.log.Debug().Msg("complex ack")

		// save the response and set the segmentation context
		if err := s.setSegmentationContext(_apdu); err != nil {
			return errors.Wrap(err, "error settings segmentation context")
		}

		// the segment size is the minimum of the size of the largest packet that can be delivered to the client and the
		//            largest it can accept
		if s.deviceInfo == nil || s.deviceInfo.MaximumNpduLength == nil {
			s.segmentSize = uint(s.maxApduLengthAccepted.NumberOfOctets())
		} else {
			s.segmentSize = min(*s.deviceInfo.MaximumNpduLength, uint(s.maxApduLengthAccepted.NumberOfOctets()))
		}

		// compute the segment count
		if len(_apdu.GetSegment()) == 0 {
			// always at least one segment
			s.segmentCount = 1
		} else {
			// split into chunks, maybe need one more
			var segmentCount, more int
			segmentCount, more = len(s.segmentAPDU.serviceBytes)/int(s.segmentSize), len(s.segmentAPDU.serviceBytes)%int(s.segmentSize)
			if more > 0 {
				s.segmentCount += 1
			}
			s.log.Debug().Int("segmentCount", segmentCount).Msg("segment count")

			// make sure we support segmented transmit if we need to
			if s.segmentCount > 1 {
				s.log.Debug().Uint8("currentSegmentCount", s.segmentCount).Msg("segmentation required, currentSegmentCount segments")

				// make sure we support segmented transmit
				if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
					s.log.Debug().Msg("server can't send segmented requests")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(NewArgs(abort), NoKWArgs)
				}

				// make sure client supports segmented receive
				if !s.segmentedResponseAccepted {
					s.log.Debug().Msg("client can't receive segmented responses")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(NewArgs(abort), NoKWArgs)
				}

				// make sure we don't exceed the number of segments in our response that the client said it was willing to accept
				//                in the request
				if s.segmentCount > s.maxSegmentsAccepted.MaxSegments() {
					s.log.Debug().Msg("client can't receive enough segments")
					abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(NewArgs(abort), NoKWArgs)
				}
			}

			// initialize the state
			s.segmentRetryCount = 0
			s.initialSequenceNumber = 0
			s.actualWindowSize = nil

			// send out the first segment (or the whole thing)
			if s.segmentCount == 1 {
				if err := s.Response(args, NoKWArgs); err != nil {
					s.log.Debug().Err(err).Msg("error sending response")
				}
				if err := s.setState(SSMState_COMPLETED, nil); err != nil {
					return errors.Wrap(err, "Error setting state to aborted")
				}
			} else {
				segment, _, err := s.getSegment(0)
				if err != nil {
					return errors.Wrap(err, "error getting first segment")
				}
				if err := s.Response(NewArgs(segment), NoKWArgs); err != nil {
					s.log.Debug().Err(err).Msg("error sending response")
				}
				if err := s.setState(SSMState_SEGMENTED_RESPONSE, nil); err != nil {
					return errors.Wrap(err, "Error setting state to aborted")
				}
			}
		}
	default:
		return errors.Errorf("Invalid APDU type %T", apdu)
	}
	return nil
}

// ProcessTask This function is called when the client has failed to send all the segments of a segmented request,
//
//	the application has taken too long to complete the request, or the client failed to ack the segments of a
//	segmented response
func (s *ServerSSM) ProcessTask() error {
	s.log.Debug().Msg("ProcessTask")
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
func (s *ServerSSM) abort(reason readWriteModel.BACnetAbortReason) (PDU, error) {
	s.log.Debug().Stringer("apdu", reason).Msg("abort")

	// change the state to aborted
	if err := s.setState(SSMState_ABORTED, nil); err != nil {
		return nil, errors.Wrap(err, "Error setting state to aborted")
	}

	// build an abort _PDU to return
	abortApdu := readWriteModel.NewAPDUAbort(true, s.invokeId, readWriteModel.NewBACnetAbortReasonTagged(reason, uint32(reason), 0), 0)
	// return it
	return NewPDU(abortApdu), nil
}

func (s *ServerSSM) idle(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("idle")

	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		return errors.Errorf("Invalid APDU type %T", apdu)
	} else {
		apduConfirmedRequest = apdu
	}

	// save the invoke ID
	s.invokeId = apduConfirmedRequest.GetInvokeId()
	s.log.Debug().Uint8("invokeId", s.invokeId).Msg("invoke ID")

	// remember if the client accepts segmented responses
	s.segmentedResponseAccepted = apduConfirmedRequest.GetSegmentedResponseAccepted()

	// if there is a cache record, check to see if it needs to be updated
	if apduConfirmedRequest.GetSegmentedResponseAccepted() && s.deviceInfo != nil {
		switch *s.deviceInfo.SegmentationSupported {
		case readWriteModel.BACnetSegmentation_NO_SEGMENTATION:
			s.log.Debug().Msg("client actually supports segmented receive")
			segmentedReceive := readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE
			s.deviceInfo.SegmentationSupported = &segmentedReceive

		// TODO: bacpypes updates the cache here but as we have a pointer  to the entry we should need that. Maybe we should because concurrency... lets see later
		case readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT:
			s.log.Debug().Msg("client actually supports both segmented transmit and receive")
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
	s.maxApduLengthAccepted = getMaxApduLengthAccepted
	if s.deviceInfo != nil && s.deviceInfo.MaximumApduLengthAccepted != nil {
		if *s.deviceInfo.MaximumApduLengthAccepted < s.maxApduLengthAccepted {
			s.log.Debug().Msg("apdu max reponse encoding error")
		} else {
			s.maxApduLengthAccepted = *s.deviceInfo.MaximumApduLengthAccepted
		}
	}
	s.log.Debug().Stringer("maxApduLengthAccepted", s.maxApduLengthAccepted).Msg("maxApduLengthAccepted")

	// save the number of segments the client is willing to accept in the ack, if this is None then the value is unknown or more than 64
	getMaxSegmentsAccepted := apduConfirmedRequest.GetMaxSegmentsAccepted()
	s.maxSegmentsAccepted = getMaxSegmentsAccepted

	// unsegmented request
	if len(apduConfirmedRequest.GetSegment()) <= 0 {
		if err := s.setState(SSMState_AWAIT_RESPONSE, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		return s.Request(NewArgs(apdu), NoKWArgs)
	}

	// make sure we support segmented requests
	if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		return s.Response(NewArgs(abort), NoKWArgs)
	}

	// save the response and set the segmentation context
	if err := s.setSegmentationContext(apduConfirmedRequest); err != nil {
		return errors.Wrap(err, "error settings segmentation context")
	}

	// the window size is the minimum of what I would propose and what the device has proposed
	proposedWindowSize := *apduConfirmedRequest.GetProposedWindowSize()
	configuredWindowSize := s.ssmSAP.GetProposedWindowSize()
	minWindowSize := min(proposedWindowSize, configuredWindowSize)
	s.actualWindowSize = &minWindowSize
	s.log.Debug().
		Uint8("proposedWindowSize", proposedWindowSize).
		Uint8("configuredWindowSize", configuredWindowSize).
		Uint8("minWindowSize", minWindowSize).
		Msg("actualWindowSize? min(proposedWindowSize, configuredWindowSize) -> minWindowSize")

	// initialize the state
	s.lastSequenceNumber = 0
	s.initialSequenceNumber = 0
	if err := s.setState(SSMState_SEGMENTED_REQUEST, &s.segmentTimeout); err != nil {
		return errors.Wrap(err, "Error setting state to aborted")
	}

	// send back a segment ack
	segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
	s.log.Debug().Stringer("segack", segack).Msg("segAck")
	return s.Response(NewArgs(NewPDU(segack)), NoKWArgs)
}

func (s *ServerSSM) segmentedRequest(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("segmentedRequest")

	// some kind of problem
	if _, ok := apdu.(readWriteModel.APDUAbortExactly); ok {
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		return s.Response(NewArgs(apdu), NoKWArgs)
	}

	// the only messages we should be getting are confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if castedApdu, ok := apdu.(readWriteModel.APDUConfirmedRequestExactly); !ok {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
			s.log.Debug().Err(err).Msg("error sending response")
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
		if err := s.Request(NewArgs(abort), NoKWArgs); err != nil { // send it ot the device
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(NewArgs(abort), NoKWArgs); err != nil { // send it ot the application
			s.log.Debug().Err(err).Msg("error sending response")
		}
	}

	// proper segment number
	if actualSequenceNumber, expectedSequenceNumber := *apduConfirmedRequest.GetSequenceNumber(), s.lastSequenceNumber+1; actualSequenceNumber != expectedSequenceNumber {
		s.log.Debug().
			Uint8("actualSequenceNumber", actualSequenceNumber).
			Uint8("expectedSequenceNumber", expectedSequenceNumber).
			Msg("segment actualSequenceNumber received out of order, should be expectedSequenceNumber")

		// segment received out of order
		s.RestartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(true, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		return s.Response(NewArgs(NewPDU(segack)), NoKWArgs)
	}

	// add the data
	if err := s.appendSegment(apdu); err != nil {
		return errors.Wrap(err, "error appending segment")
	}

	// update the sequence number
	s.lastSequenceNumber++

	// last segment?
	if !apduConfirmedRequest.GetMoreFollows() {
		s.log.Debug().Msg("No more follows")

		// send back the final segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.lastSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Response(NewArgs(NewPDU(segack)), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending response")
		}

		// forward the whole thing to the application
		applicationTimeout := s.ssmSAP.GetApplicationTimeout()
		if err := s.setState(SSMState_AWAIT_RESPONSE, &applicationTimeout); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		// TODO: here we need to rebuild again yada yada
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now..
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), s.log, options.WithPassLoggerToModel(s.passLogToModel))
		parse, err := readWriteModel.APDUParse(ctxForModel, s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := s.Request(NewArgs(NewPDU(parse)), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
	} else if *apduConfirmedRequest.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		s.log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.RestartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Response(NewArgs(NewPDU(segack)), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending response")
		}
	} else {
		// wait for more segments
		s.RestartTimer(s.segmentTimeout)
	}

	return nil
}

func (s *ServerSSM) segmentedRequestTimeout() error {
	s.log.Debug().Msg("segmentedRequestTimeout")

	// give up
	if err := s.setState(SSMState_ABORTED, nil); err != nil {
		return errors.Wrap(err, "Error setting state to aborted")
	}
	return nil
}

func (s *ServerSSM) awaitResponse(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("awaitResponse")

	switch apdu.GetMessage().(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		s.log.Debug().Msg("client is trying this request again")
	case readWriteModel.APDUAbortExactly:
		s.log.Debug().Msg("client aborting this request")

		// forward to the application
		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Request(NewArgs(apdu), NoKWArgs); err != nil { // send it ot the device
			s.log.Debug().Err(err).Msg("error sending request")
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

// awaitResponseTimeout This function is called when the application has taken too long to respond to a clients request.
//
//	The client has probably long since given up
func (s *ServerSSM) awaitResponseTimeout() error {
	s.log.Debug().Msg("awaitResponseTimeout")

	abort, err := s.abort(readWriteModel.BACnetAbortReason(64)) // Note: this is a proprietary code used by bacpypes for server timeout. We just use that here too to keep consistent
	if err != nil {
		return errors.Wrap(err, "error creating abort")
	}
	if err := s.Request(NewArgs(abort), NoKWArgs); err != nil {
		s.log.Debug().Err(err).Msg("error sending request")
	}
	return nil
}

func (s *ServerSSM) segmentedResponse(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("segmentedResponse")

	// client is ready for the next segment
	switch _apdu := apdu.GetMessage().(type) {
	case readWriteModel.APDUSegmentAckExactly:
		s.log.Debug().Msg("segment ack")

		// actual window size is provided by client
		getActualWindowSize := _apdu.GetActualWindowSize()
		s.actualWindowSize = &getActualWindowSize

		// duplicate ack received?
		if !s.inWindow(_apdu.GetSequenceNumber(), s.initialSequenceNumber) {
			s.log.Debug().Msg("not in window")
			s.RestartTimer(s.segmentTimeout)
		} else if s.sentAllSegments {
			// final ack received?
			s.log.Debug().Msg("all done sending response")
			if err := s.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "Error setting state to aborted")
			}
		} else {
			s.log.Debug().Msg("more segments to send")

			s.initialSequenceNumber = _apdu.GetSequenceNumber() + 1
			actualWindowSize := _apdu.GetActualWindowSize()
			s.actualWindowSize = &actualWindowSize
			s.segmentRetryCount = 0
			if err := s.fillWindow(s.initialSequenceNumber); err != nil {
				return errors.Wrap(err, "error filling window")
			}
			s.RestartTimer(s.segmentRetryCount)
		}
	// some kind of problem
	case readWriteModel.APDUAbortExactly:
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Response(NewArgs(apdu), NoKWArgs); err != nil { // send it ot the application
			s.log.Debug().Err(err).Msg("error sending response")
		}
	default:
		return errors.Errorf("Invalid APDU type %T", apdu)
	}
	return nil
}

func (s *ServerSSM) segmentedResponseTimeout() error {
	s.log.Debug().Msg("segmentedResponseTimeout")

	// try again
	if s.segmentRetryCount < s.numberOfApduRetries {
		s.segmentRetryCount++
		s.StartTimer(s.segmentTimeout)
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

	localDevice           *LocalDeviceObject
	deviceInfoCache       *DeviceInfoCache
	nextInvokeId          uint8
	clientTransactions    []*ClientSSM
	serverTransactions    []*ServerSSM
	numberOfApduRetries   int
	apduTimeout           uint
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted
	segmentationSupported readWriteModel.BACnetSegmentation
	segmentTimeout        uint
	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted
	proposedWindowSize    uint8
	dccEnableDisable      readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable
	applicationTimeout    uint

	log zerolog.Logger
}

func NewStateMachineAccessPoint(localLog zerolog.Logger, localDevice *LocalDeviceObject, deviceInfoCache *DeviceInfoCache, sapID *int, cid *int) (*StateMachineAccessPoint, error) {
	localLog.Debug().
		Stringer("localDevice", localDevice).
		Stringer("deviceInfoCache", deviceInfoCache).
		Interface("sapID", sapID).
		Interface("cid", cid).
		Msg("NewStateMachineAccessPoint")

	s := &StateMachineAccessPoint{
		// save a reference to the device information cache
		localDevice:     localDevice,
		deviceInfoCache: deviceInfoCache,

		// client settings
		nextInvokeId:       1,
		clientTransactions: nil,

		// server settings
		serverTransactions: nil,

		// confirmed request defaults
		numberOfApduRetries:   3,
		apduTimeout:           3000,
		maxApduLengthAccepted: readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1024,

		// segmentation defaults
		segmentationSupported: readWriteModel.BACnetSegmentation_NO_SEGMENTATION,
		segmentTimeout:        1500,
		maxSegmentsAccepted:   readWriteModel.MaxSegmentsAccepted_NUM_SEGMENTS_02,
		proposedWindowSize:    2,

		// device communication control
		dccEnableDisable: readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE,

		// how long the state machine is willing to wait for the application
		// layer to form a response and send it
		applicationTimeout: 3000,

		log: localLog,
	}
	// basic initialization
	client, err := NewClient(localLog, cid, s)
	if err != nil {
		return nil, errors.Wrapf(err, "error building client for %d", cid)
	}
	s.Client = client
	serviceAccessPoint, err := NewServiceAccessPoint(localLog, sapID, s)
	if err != nil {
		return nil, errors.Wrapf(err, "error building serviceAccessPoint for %d", sapID)
	}
	s.ServiceAccessPoint = serviceAccessPoint
	return s, nil
}

func (s *StateMachineAccessPoint) String() string {
	return fmt.Sprintf("StateMachineAccessPoint(TBD...)") // TODO: fill some info here
}

// getNextInvokeId Called by clients to get an unused invoke ID
func (s *StateMachineAccessPoint) getNextInvokeId(address Address) (uint8, error) {
	s.log.Debug().Msg("getNextInvokeId")

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
			// TODO: replace deep equal
			if invokeId == tr.invokeId && address.Equals(tr.pduAddress) {
				return invokeId, nil
			}
		}
	}
}

func (s *StateMachineAccessPoint) GetDefaultAPDUTimeout() uint {
	return s.apduTimeout
}

func (s *StateMachineAccessPoint) GetDefaultSegmentationSupported() readWriteModel.BACnetSegmentation {
	return s.segmentationSupported
}

func (s *StateMachineAccessPoint) GetDefaultAPDUSegmentTimeout() uint {
	return s.segmentTimeout
}

func (s *StateMachineAccessPoint) GetDefaultMaxSegmentsAccepted() readWriteModel.MaxSegmentsAccepted {
	return s.maxSegmentsAccepted
}

func (s *StateMachineAccessPoint) GetDefaultMaximumApduLengthAccepted() readWriteModel.MaxApduLengthAccepted {
	return s.maxApduLengthAccepted
}

// Confirmation Packets coming up the stack are APDU's
func (s *StateMachineAccessPoint) Confirmation(args Args, kwargs KWArgs) error { // TODO: note we need a special method here as we don't contain src in the apdu
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")
	apdu := args.Get0PDU()

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		s.log.Debug().Msg("communications enabled")
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		apduType := apdu.GetMessage().(interface {
			GetApduType() readWriteModel.ApduType
		}).GetApduType()
		switch {
		case apduType == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			apdu.GetMessage().(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_DEVICE_COMMUNICATION_CONTROL:
			s.log.Debug().Msg("continue with DCC request")
		case apduType == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			apdu.GetMessage().(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_REINITIALIZE_DEVICE:
			s.log.Debug().Msg("continue with reinitialize device")
		case apduType == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU &&
			apdu.GetMessage().(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS:
			s.log.Debug().Msg("continue with Who-Is")
		default:
			s.log.Debug().Msg("not a Who-Is, dropped")
			return nil
		}
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		s.log.Debug().Msg("initiation disabled")
	}

	var pduSource = apdu.GetPDUSource()

	switch _apdu := apdu.GetMessage().(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		// Find duplicates of this request
		var tr *ServerSSM
		for _, serverTransactionElement := range s.serverTransactions {
			if _apdu.GetInvokeId() == serverTransactionElement.invokeId && pduSource.Equals(serverTransactionElement.pduAddress) {
				tr = serverTransactionElement
				break
			}
		}
		if tr == nil {
			// build a server transaction
			var err error
			tr, err = NewServerSSM(s.log, s, pduSource)
			if err != nil {
				return errors.Wrap(err, "Error building server ssm")
			}
			s.serverTransactions = append(s.serverTransactions, tr)
		}

		// let it run with the apdu
		if err := tr.Indication(NewArgs(apdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error runnning indication")
		}
	case readWriteModel.APDUUnconfirmedRequestExactly:
		// deliver directly to the application
		if err := s.SapRequest(NewArgs(apdu), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUComplexAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		// find the client transaction this is acking
		var tr *ClientSSM
		for _, _tr := range s.clientTransactions {
			if _apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == _tr.invokeId && pduSource.Equals(_tr.pduAddress) {
				tr = _tr
				break
			}
		}
		if tr == nil {
			// TODO: log at least
			return nil
		}

		// send the packet on to the transaction
		if err := tr.Confirmation(NewArgs(apdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	case readWriteModel.APDUAbortExactly:
		// find the transaction being aborted
		if _apdu.GetServer() {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && pduSource.Equals(tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Confirmation(NewArgs(apdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if _apdu.GetOriginalInvokeId() == serverTransactionElement.invokeId && pduSource.Equals(serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Indication(NewArgs(apdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	case readWriteModel.APDUSegmentAckExactly:
		// find the transaction being aborted
		if _apdu.GetServer() {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && pduSource.Equals(tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Confirmation(NewArgs(apdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if _apdu.GetOriginalInvokeId() == serverTransactionElement.invokeId && pduSource.Equals(serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Indication(NewArgs(apdu), NoKWArgs); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

// SapIndication This function is called when the application is requesting a new transaction as a client.
func (s *StateMachineAccessPoint) SapIndication(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")
	apdu := args.Get0PDU()

	pduDestination := apdu.GetPDUDestination()

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		s.log.Debug().Msg("communications enabled")
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		s.log.Debug().Msg("communications disabled")
		return nil
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		s.log.Debug().Msg("initiation disabled")
		// TODO: this should be quarded
		if apdu.GetMessage().(readWriteModel.APDU).GetApduType() == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU && apdu.(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_I_AM {
			s.log.Debug().Msg("continue with I-Am")
		} else {
			s.log.Debug().Msg("not an I-Am")
			return nil
		}
	}

	switch _apdu := apdu.GetMessage().(type) {
	case readWriteModel.APDUUnconfirmedRequestExactly:
		// deliver to the device
		if err := s.Request(NewArgs(apdu), NoKWArgs); err != nil {
			s.log.Debug().Err(err).Msg("error sending the request")
		}
	case readWriteModel.APDUConfirmedRequestExactly:
		// make sure it has an invoke ID
		// TODO: here it is getting slightly different: usually we give the invoke id from the outside as it is build already. So maybe we need to adjust that (we never create it, we need to check for collisions but maybe we should change that so we move the creation down here)
		// s.getNextInvokeId()...
		for _, tr := range s.clientTransactions {
			if _apdu.GetInvokeId() == tr.invokeId && pduDestination.Equals(tr.pduAddress) {
				return errors.New("invoke ID in use")
			}
		}

		// warning for bogus requests
		// TODO: not sure if we have that or if it is relvant (localstationaddr)

		// create a client transaction state machine
		tr, err := NewClientSSM(s.log, s, pduDestination)
		if err != nil {
			return errors.Wrap(err, "error creating client ssm")
		}

		// add it to our transactions to track it
		s.clientTransactions = append(s.clientTransactions, tr)

		// let it run
		if err := tr.Indication(NewArgs(apdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error doing indication")
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}

	return nil
}

// SapConfirmation This function is called when the application is responding to a request, the apdu may be a simple
//
//	ack, complex ack, error, reject or abort
func (s *StateMachineAccessPoint) SapConfirmation(args Args, kwargs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")
	apdu := args.Get0PDU()
	pduDestination := apdu.GetPDUDestination()
	switch apdu.GetMessage().(type) {
	case readWriteModel.APDUSimpleAckExactly, readWriteModel.APDUComplexAckExactly, readWriteModel.APDUErrorExactly, readWriteModel.APDURejectExactly:
		// find the client transaction this is acking
		var tr *ServerSSM
		for _, tr := range s.serverTransactions {
			if apdu.(interface{ GetOriginalInvokeId() uint8 }).GetOriginalInvokeId() == tr.invokeId && pduDestination.Equals(tr.pduAddress) {
				break
			}
		}
		if tr == nil {
			// TODO: log at least
			return nil
		}

		// pass control to the transaction
		if err := tr.Confirmation(NewArgs(apdu), NoKWArgs); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

func (s *StateMachineAccessPoint) GetDeviceInfoCache() *DeviceInfoCache {
	return s.deviceInfoCache
}

func (s *StateMachineAccessPoint) GetLocalDevice() *LocalDeviceObject {
	return s.localDevice
}

func (s *StateMachineAccessPoint) GetProposedWindowSize() uint8 {
	return s.proposedWindowSize
}

func (s *StateMachineAccessPoint) GetClientTransactions() []*ClientSSM {
	return s.clientTransactions
}

func (s *StateMachineAccessPoint) RemoveClientTransaction(c *ClientSSM) {
	indexFound := -1
	for i, tr := range s.clientTransactions {
		if tr == c {
			indexFound = i
			break
		}
	}
	if indexFound >= 0 {
		s.clientTransactions = append(s.clientTransactions[:indexFound], s.clientTransactions[indexFound+1:]...)
	}
}

func (s *StateMachineAccessPoint) GetServerTransactions() []*ServerSSM {
	return s.serverTransactions
}

func (s *StateMachineAccessPoint) RemoveServerTransaction(sssm *ServerSSM) {
	indexFound := -1
	for i, tr := range s.serverTransactions {
		if tr == sssm {
			indexFound = i
			break
		}
	}
	if indexFound >= 0 {
		s.serverTransactions = append(s.serverTransactions[:indexFound], s.serverTransactions[indexFound+1:]...)
	}
}

func (s *StateMachineAccessPoint) GetApplicationTimeout() uint {
	return s.applicationTimeout
}

type ApplicationServiceAccessPoint struct {
	*ApplicationServiceElement
	*ServiceAccessPoint

	log zerolog.Logger
}

func NewApplicationServiceAccessPoint(localLog zerolog.Logger, aseID *int, sapID *int) (*ApplicationServiceAccessPoint, error) {
	a := &ApplicationServiceAccessPoint{
		log: localLog,
	}
	applicationServiceElement, err := NewApplicationServiceElement(localLog, aseID, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating application service element")
	}
	a.ApplicationServiceElement = applicationServiceElement
	serviceAccessPoint, err := NewServiceAccessPoint(localLog, sapID, a)
	if err != nil {
		return nil, errors.Wrap(err, "error creating service access point")
	}
	a.ServiceAccessPoint = serviceAccessPoint
	return a, nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Indication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Indication")
	apdu := args.Get0PDU()

	switch _apdu := apdu.GetMessage().(type) {
	case readWriteModel.APDUConfirmedRequestExactly:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetConfirmedServiceChoiceKnows(uint8(_apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(NewArgs(apdu), NoKWArgs)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(NewArgs(apdu), NoKWArgs); err != nil {
				return err
			}
		} else {
			a.log.Debug().Err(errorFound).Msg("got error")

			// TODO: map it to a error... code temporary placeholder
			return a.Response(NewArgs(NewPDU(readWriteModel.NewAPDUReject(_apdu.GetInvokeId(), nil, 0))), NoKWArgs)
		}
	case readWriteModel.APDUUnconfirmedRequestExactly:
		//assume no errors found
		var errorFound error
		if !readWriteModel.BACnetUnconfirmedServiceChoiceKnows(uint8(_apdu.GetServiceRequest().GetServiceChoice())) {
			errorFound = errors.New("unrecognized service")
		}

		if errorFound == nil {
			errorFound = a.SapRequest(NewArgs(apdu), NoKWArgs)
		}
		// TODO: the handling here gets a bit different now... need to wrap the head around how to do this (error handling etc)

		if errorFound == nil {
			if err := a.SapRequest(NewArgs(apdu), NoKWArgs); err != nil {
				return err
			}
		} else {
			a.log.Debug().Err(errorFound).Msg("got error")
		}

	default:
		return errors.Errorf("unknown _PDU type %T", apdu)
	}
	return nil
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapIndication(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapIndication")

	// TODO: check if we need to check apdu here

	return a.Request(args, kwargs)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) Confirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("Confirmation")

	// TODO: check if we need to check apdu here

	// forward the decoded packet
	return a.SapResponse(args, kwargs)
}

// TODO: big WIP
func (a *ApplicationServiceAccessPoint) SapConfirmation(args Args, kwargs KWArgs) error {
	a.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwargs).Msg("SapConfirmation")

	// TODO: check if we need to check apdu here

	return a.Response(args, kwargs)
}
