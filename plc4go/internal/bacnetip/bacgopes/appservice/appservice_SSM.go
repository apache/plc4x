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

package appservice

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
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
	ServiceAccessPointContract
	Client
	GetDeviceInfoCache() *DeviceInfoCache
	GetLocalDevice() LocalDeviceObject
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
	if localDevice.GetNumberOfAPDURetries() != nil {
		numberOfApduRetries = *localDevice.GetNumberOfAPDURetries()
	}
	var apduTimeout uint
	if localDevice.GetAPDUTimeout() != nil {
		apduTimeout = *localDevice.GetAPDUTimeout()
	} else {
		apduTimeout = sap.GetDefaultAPDUTimeout()
	}
	var segmentationSupported readWriteModel.BACnetSegmentation
	if localDevice.GetSegmentationSupported() != nil {
		segmentationSupported = *localDevice.GetSegmentationSupported()
	} else {
		segmentationSupported = sap.GetDefaultSegmentationSupported()
	}
	var segmentTimeout uint
	if localDevice.GetAPDUSegmentTimeout() != nil {
		segmentTimeout = *localDevice.GetAPDUSegmentTimeout()
	} else {
		segmentTimeout = sap.GetDefaultAPDUSegmentTimeout()
	}
	var maxSegmentsAccepted readWriteModel.MaxSegmentsAccepted
	if localDevice.GetMaxSegmentsAccepted() != nil {
		maxSegmentsAccepted = *localDevice.GetMaxSegmentsAccepted()
	} else {
		maxSegmentsAccepted = sap.GetDefaultMaxSegmentsAccepted()
	}
	var maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted
	if localDevice.GetMaximumApduLengthAccepted() != nil {
		maxApduLengthAccepted = *localDevice.GetMaximumApduLengthAccepted()
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
	ssm.OneShotTask = NewOneShotTask(ssm)
	ssm.SSMProcessingRequirements = sap
	return ssm, nil
}

func (s *SSM) StartTimer(millis uint) {
	s.log.Debug().Uint("millis", millis).Msg("Start timer")
	s.RestartTimer(millis)
}

func (s *SSM) StopTimer() {
	s.log.Debug().Msg("Stop Timer")
	if s.GetIsScheduled() {
		s.log.Debug().Msg("is scheduled")
		s.SuspendTask()
	}
}

func (s *SSM) RestartTimer(millis uint) {
	s.log.Debug().Uint("millis", millis).Msg("restart timer")
	if s.GetIsScheduled() {
		s.log.Debug().Msg("is scheduled")
		s.SuspendTask()
	}

	s.InstallTask(WithInstallTaskOptionsDelta(time.Duration(millis) * time.Millisecond))
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
	case readWriteModel.APDUConfirmedRequest:
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
	case readWriteModel.APDUComplexAck:
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
	s.log.Debug().Uint8("index", index).Msg("GA segment")
	if s.segmentAPDU == nil {
		return nil, false, errors.New("No segment apdu set")
	}

	if index > s.segmentCount {
		return nil, false, errors.Errorf("Invalid segment number %d, APDU has %d segments", index, s.segmentCount)
	}

	// TODO: the original code does here something funky but it seems it is best to just return the original apdu
	if s.segmentCount == 1 {
		return NewPDU(NoArgs, NKW(KWCPCIDestination, s.pduAddress), WithRootMessage(s.segmentAPDU.originalApdu)), false, nil
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
		segmentAPDU = NewPDU(NoArgs, NKW(KWCPCIDestination, s.pduAddress), WithRootMessage(readWriteModel.NewAPDUConfirmedRequest(
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
		)))
	} else {
		s.log.Debug().Msg("complex ack context")
		segmentAPDU = NewPDU(NoArgs, NKW(KWCPCIDestination, s.pduAddress), WithRootMessage(readWriteModel.NewAPDUComplexAck(
			true,
			moreFollows,
			s.segmentAPDU.originalInvokeId,
			&sequenceNumber,
			proposedWindowSize,
			nil,
			serviceChoice,
			segmentBytes,
			0,
		)))
	}
	return segmentAPDU, moreFollows, nil
}

// TODO: check that function. looks a bit wonky to just append the payloads like that
// appendSegment This function appends the apdu content to the end of the current APDU being built.  The segmentAPDU is
//
//	the context
func (s *SSM) appendSegment(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("appendSegment")
	switch apdu := apdu.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		if apdu.GetSegmentedMessage() || apdu.GetMoreFollows() {
			return errors.New("Can't handle already segmented message")
		}
		serializedBytes, err := apdu.GetServiceRequest().Serialize()
		if err != nil {
			return errors.Wrap(err, "Can serialize service request")
		}
		s.segmentAPDU.serviceBytes = append(s.segmentAPDU.serviceBytes, serializedBytes...)
	case readWriteModel.APDUComplexAck:
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
		if err := s.ssmSAP.Request(NA(NewPDU(NoArgs, NKW(KWCPCIDestination, s.pduAddress)), WithRootMessage(apdu.GetRootMessage())), NoKWArgs()); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if moreFollows {
			s.sentAllSegments = true
		}
	}
	return nil
}
