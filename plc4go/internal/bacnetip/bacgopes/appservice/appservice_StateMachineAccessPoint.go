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
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comm"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/local/device"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

//go:generate go tool plc4xGenerator -type=StateMachineAccessPoint -prefix=appservice_
type StateMachineAccessPoint struct {
	ClientContract
	ServiceAccessPointContract
	*DefaultRFormatter `ignore:"true"`

	localDevice           LocalDeviceObject
	deviceInfoCache       *DeviceInfoCache
	nextInvokeId          uint8
	clientTransactions    []*ClientSSM
	serverTransactions    []*ServerSSM
	numberOfApduRetries   int
	apduTimeout           uint
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted `stringer:"true"`
	segmentationSupported readWriteModel.BACnetSegmentation    `stringer:"true"`
	segmentTimeout        uint
	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted `stringer:"true"`
	proposedWindowSize    uint8
	dccEnableDisable      readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable `stringer:"true"`
	applicationTimeout    uint

	log zerolog.Logger
}

func NewStateMachineAccessPoint(localLog zerolog.Logger, localDevice LocalDeviceObject, options ...Option) (*StateMachineAccessPoint, error) {
	s := &StateMachineAccessPoint{
		DefaultRFormatter: NewDefaultRFormatter(),
		// save a reference to the device information cache
		localDevice: localDevice,

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
	ApplyAppliers(options, s)
	optionsForParent := AddLeafTypeIfAbundant(options, s)
	// basic initialization
	var err error
	s.ClientContract, err = NewClient(s.log, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error building client")
	}
	s.ServiceAccessPointContract, err = NewServiceAccessPoint(s.log, optionsForParent...)
	if err != nil {
		return nil, errors.Wrap(err, "error building serviceAccessPoint")
	}
	if _debug != nil {
		_debug("__init__ localDevice=%r deviceInfoCache=%r sap=%r cid=%r", localDevice, s.deviceInfoCache, s.GetServiceID(), s.GetClientID())
	}
	s.log.Debug().
		Stringer("localDevice", localDevice).
		Stringer("deviceInfoCache", s.deviceInfoCache).
		Interface("sapID", s.GetServiceID()).
		Interface("cid", s.GetClientID()).
		Msg("NewStateMachineAccessPoint")
	return s, nil
}

func WithStateMachineAccessPointDeviceInfoCache(deviceInfoCache *DeviceInfoCache) GenericApplier[*StateMachineAccessPoint] {
	return WrapGenericApplier(func(s *StateMachineAccessPoint) { s.deviceInfoCache = deviceInfoCache })
}

// getNextInvokeId Called by clients to get an unused invoke ID
func (s *StateMachineAccessPoint) getNextInvokeId(address Address) (uint8, error) {
	if _debug != nil {
		_debug("get_next_invoke_id")
	}
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

func (s *StateMachineAccessPoint) SetDeviceInfoCache(cache *DeviceInfoCache) {
	s.deviceInfoCache = cache
}

// Confirmation Packets coming up the stack are APDU's
func (s *StateMachineAccessPoint) Confirmation(args Args, kwArgs KWArgs) error { // TODO: note we need a special method here as we don't contain src in the apdu
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")
	pdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("confirmation %r", pdu)
	}

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		if _debug != nil {
			_debug("    - communications enabled")
		}
		s.log.Debug().Msg("communications enabled")
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		apduType := pdu.GetRootMessage().(interface {
			GetApduType() readWriteModel.ApduType
		}).GetApduType()
		switch {
		case apduType == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			pdu.GetRootMessage().(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_DEVICE_COMMUNICATION_CONTROL:
			if _debug != nil {
				_debug("    - continue with DCC request")
			}
			s.log.Debug().Msg("continue with DCC request")
		case apduType == readWriteModel.ApduType_CONFIRMED_REQUEST_PDU &&
			pdu.GetRootMessage().(readWriteModel.APDUConfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetConfirmedServiceChoice_REINITIALIZE_DEVICE:
			if _debug != nil {
				_debug("    - continue with reinitialize device")
			}
			s.log.Debug().Msg("continue with reinitialize device")
		case apduType == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU &&
			pdu.GetRootMessage().(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_WHO_IS:
			if _debug != nil {
				_debug("    - continue with Who-Is")
			}
			s.log.Debug().Msg("continue with Who-Is")
		default:
			if _debug != nil {
				_debug("    - not a Who-Is, dropped")
			}
			s.log.Debug().Msg("not a Who-Is, dropped")
			return nil
		}
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		if _debug != nil {
			_debug("    - initiation disabled")
		}
		s.log.Debug().Msg("initiation disabled")
	}

	// make a more focused interpretation
	atype := APDUTypes[pdu.(APDU).GetApduType()] // TODO: why are we suddenly assuming apdu now here...
	if atype == nil {
		s.log.Warn().Msgf("    - unknown apduType: %s", pdu.(APDU).GetApduType())
		return nil
	}

	// decode it
	apdu := atype()
	if err := apdu.Decode(pdu); err != nil {
		return errors.Wrap(err, "apdu decode failed")
	}
	if _debug != nil {
		_debug("    - apdu: %r", apdu)
	}

	switch apdu := apdu.(type) {
	case *ConfirmedRequestPDU:
		// Find duplicates of this request
		var tr *ServerSSM
		for _, serverTransactionElement := range s.serverTransactions {
			if apdu.GetInvokeId() == serverTransactionElement.invokeId && apdu.GetPDUSource().Equals(serverTransactionElement.pduAddress) {
				tr = serverTransactionElement
				break
			}
		}
		if tr == nil {
			// build a server transaction
			var err error
			tr, err = NewServerSSM(s.log, s, apdu.GetPDUSource())
			if err != nil {
				return errors.Wrap(err, "Error building server ssm")
			}
			s.serverTransactions = append(s.serverTransactions, tr)
		}

		// let it run with the apdu
		if err := tr.Indication(NA(apdu), NoKWArgs()); err != nil {
			return errors.Wrap(err, "error runnning indication")
		}
	case *UnconfirmedRequestPDU:
		// deliver directly to the application
		if err := s.SapRequest(NA(apdu), NoKWArgs()); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUComplexAck, readWriteModel.APDUError, readWriteModel.APDUReject:
		// find the client transaction this is acking
		var tr *ClientSSM
		for _, _tr := range s.clientTransactions {
			if apdu.(interface{ GetInvokeId() uint8 }).GetInvokeId() == _tr.invokeId && apdu.(interface{ GetPDUSource() *Address }).GetPDUSource().Equals(_tr.pduAddress) {
				tr = _tr
				break
			}
		}
		if tr == nil {
			// TODO: log at least
			return nil
		}

		// send the packet on to the transaction
		if err := tr.Confirmation(NA(apdu), NoKWArgs()); err != nil {
			return errors.Wrap(err, "error running confirmation")
		}
	case *AbortPDU:
		// find the transaction being aborted
		if apdu.GetAPDUSrv() != nil {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if *apdu.GetApduInvokeID() == tr.invokeId && apdu.GetPDUSource().Equals(tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Confirmation(NA(apdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if *apdu.GetApduInvokeID() == serverTransactionElement.invokeId && apdu.GetPDUSource().Equals(serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Indication(NA(apdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	case *SegmentAckPDU:
		// find the transaction being aborted
		if apdu.GetServer() != nil {
			var tr *ClientSSM
			for _, tr := range s.clientTransactions {
				if *apdu.GetApduInvokeID() == tr.invokeId && apdu.GetPDUSource().Equals(tr.pduAddress) {
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Confirmation(NA(apdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error running confirmation")
			}
		} else {
			var tr *ServerSSM
			for _, serverTransactionElement := range s.serverTransactions {
				if *apdu.GetApduInvokeID() == serverTransactionElement.invokeId && apdu.GetPDUSource().Equals(serverTransactionElement.pduAddress) {
					tr = serverTransactionElement
					break
				}
			}
			if tr == nil {
				// TODO: log at least
				return nil
			}

			// send the packet on to the transaction
			if err := tr.Indication(NA(apdu), NoKWArgs()); err != nil {
				return errors.Wrap(err, "error running indication")
			}
		}
	default:
		return errors.Errorf("invalid APDU type %T", apdu)
	}
	return nil
}

// SapIndication This function is called when the application is requesting a new transaction as a client.
func (s *StateMachineAccessPoint) SapIndication(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapIndication")
	apdu := GA[APDU](args, 0)
	if _debug != nil {
		_debug("sap_indication %r", apdu)
	}

	pduDestination := apdu.GetPDUDestination()

	// check device communication control
	switch s.dccEnableDisable {
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_ENABLE:
		s.log.Debug().Msg("communications enabled")
		if _debug != nil {
			_debug("    - communications enabled")
		}
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE:
		s.log.Debug().Msg("communications disabled")
		if _debug != nil {
			_debug("    - communications disabled")
		}
		return nil
	case readWriteModel.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable_DISABLE_INITIATION:
		if _debug != nil {
			_debug("    - initiation disabled")
		}
		s.log.Debug().Msg("initiation disabled")
		// TODO: this should be quarded
		if apdu.GetRootMessage().(readWriteModel.APDU).GetApduType() == readWriteModel.ApduType_UNCONFIRMED_REQUEST_PDU && apdu.(readWriteModel.APDUUnconfirmedRequest).GetServiceRequest().GetServiceChoice() == readWriteModel.BACnetUnconfirmedServiceChoice_I_AM {
			s.log.Debug().Msg("continue with I-Am")
			if _debug != nil {
				_debug("    - continue with I-Am")
			}
		} else {
			s.log.Debug().Msg("not an I-Am")
			if _debug != nil {
				_debug("    - not an I-Am")
			}
			return nil
		}
	}

	switch _apdu := apdu.(type) {
	case readWriteModel.APDUUnconfirmedRequest:
		// deliver to the device
		if err := s.Request(NA(apdu), NoKWArgs()); err != nil {
			s.log.Debug().Err(err).Msg("error sending the request")
		}
	case readWriteModel.APDUConfirmedRequest:
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
		if _debug != nil {
			_debug("    - client segmentation state machine: %r", tr)
		}

		// add it to our transactions to track it
		s.clientTransactions = append(s.clientTransactions, tr)

		// let it run
		if err := tr.Indication(NA(apdu), NoKWArgs()); err != nil {
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
func (s *StateMachineAccessPoint) SapConfirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("SapConfirmation")
	apdu := GA[PDU](args, 0)
	if _debug != nil {
		_debug("sap_confirmation %r", apdu)
	}

	pduDestination := apdu.GetPDUDestination()
	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUComplexAck, readWriteModel.APDUError, readWriteModel.APDUReject:
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
		if err := tr.Confirmation(NA(apdu), NoKWArgs()); err != nil {
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

func (s *StateMachineAccessPoint) GetLocalDevice() LocalDeviceObject {
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

func (s *StateMachineAccessPoint) AlternateString() (string, bool) {
	if IsDebuggingActive() {
		return fmt.Sprintf("%r", s), true
	}
	return "", false
}
