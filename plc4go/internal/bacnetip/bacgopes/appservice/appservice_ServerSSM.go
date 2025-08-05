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
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

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
func (s *ServerSSM) Request(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Request")
	// TODO: ensure apdu has destination, otherwise
	// TODO: we would need a BVLC to send something or not... maybe the todo above is nonsense, as we are in a connection context
	return s.ssmSAP.SapRequest(args, kwArgs)
}

// Indication This function is called for each downstream packet related to
//
//	the transaction
func (s *ServerSSM) Indication(args Args, kwArgs KWArgs) error { // TODO: maybe use another name for that
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Indication")
	apdu := GA[PDU](args, 0)
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
func (s *ServerSSM) Response(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Response")
	// make sure it has a good source and destination
	// TODO: check if source == none
	// TODO: check if destnation = s.pduAddress

	// send it via the device
	return s.ssmSAP.Request(args, kwArgs)
}

// Confirmation This function is called when the application has provided a response and needs it to be sent to the
//
//	client.
func (s *ServerSSM) Confirmation(args Args, kwArgs KWArgs) error {
	s.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("Confirmation")

	// check to see we are in the correct state
	if s.state != SSMState_AWAIT_RESPONSE {
		s.log.Debug().Msg("warning: no expecting a response")
	}

	apdu := GA[PDU](args, 0)

	switch _apdu := apdu.GetRootMessage().(type) {
	// abort response
	case readWriteModel.APDUAbort:
		s.log.Debug().Msg("abort")

		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// end the response to the device
		return s.Response(args, kwArgs)
	// simple response
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUError, readWriteModel.APDUReject:
		s.log.Debug().Msg("simple ack, error or reject")

		// transaction completed
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}

		// send the response to the device
		return s.Response(args, kwArgs)
	// complex ack
	case readWriteModel.APDUComplexAck:
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
					return s.Response(NA(abort), NoKWArgs())
				}

				// make sure client supports segmented receive
				if !s.segmentedResponseAccepted {
					s.log.Debug().Msg("client can't receive segmented responses")
					abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(NA(abort), NoKWArgs())
				}

				// make sure we don't exceed the number of segments in our response that the client said it was willing to accept
				//                in the request
				if s.segmentCount > s.maxSegmentsAccepted.MaxSegments() {
					s.log.Debug().Msg("client can't receive enough segments")
					abort, err := s.abort(readWriteModel.BACnetAbortReason(65)) // Note: this is a proprietary code used by bacpypes for no response. We just use that here too to keep consistent
					if err != nil {
						return errors.Wrap(err, "Error creating abort")
					}
					return s.Response(NA(abort), NoKWArgs())
				}
			}

			// initialize the state
			s.segmentRetryCount = 0
			s.initialSequenceNumber = 0
			s.actualWindowSize = nil

			// send out the first segment (or the whole thing)
			if s.segmentCount == 1 {
				if err := s.Response(args, NoKWArgs()); err != nil {
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
				if err := s.Response(NA(segment), NoKWArgs()); err != nil {
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
	return NewPDU(NoArgs, NoKWArgs(), WithRootMessage(abortApdu)), nil
}

func (s *ServerSSM) idle(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("idle")

	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apdu, ok := apdu.(readWriteModel.APDUConfirmedRequest); !ok {
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
		return s.Request(NA(apdu), NoKWArgs())
	}

	// make sure we support segmented requests
	if s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && s.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		return s.Response(NA(abort), NoKWArgs())
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
	return s.Response(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segack))), NoKWArgs())
}

func (s *ServerSSM) segmentedRequest(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("segmentedRequest")

	// some kind of problem
	if _, ok := apdu.(readWriteModel.APDUAbort); ok {
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		return s.Response(NA(apdu), NoKWArgs())
	}

	// the only messages we should be getting are confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if castedApdu, ok := apdu.(readWriteModel.APDUConfirmedRequest); !ok {
		abort, err := s.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := s.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
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
		if err := s.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
			s.log.Debug().Err(err).Msg("error sending request")
		}
		if err := s.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
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
		return s.Response(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segack))), NoKWArgs())
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
		if err := s.Response(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segack))), NoKWArgs()); err != nil {
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
		parse, err := readWriteModel.APDUParse[readWriteModel.APDU](ctxForModel, s.segmentAPDU.serviceBytes, uint16(len(s.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := s.Request(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(parse))), NoKWArgs()); err != nil {
			s.log.Debug().Err(err).Msg("error sending request")
		}
	} else if *apduConfirmedRequest.GetSequenceNumber() == s.initialSequenceNumber+*s.actualWindowSize {
		s.log.Debug().Msg("last segment in the group")

		s.initialSequenceNumber = s.lastSequenceNumber
		s.RestartTimer(s.segmentTimeout)

		// send back a segment ack
		segack := readWriteModel.NewAPDUSegmentAck(false, true, s.invokeId, s.initialSequenceNumber, *s.actualWindowSize, 0)
		if err := s.Response(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segack))), NoKWArgs()); err != nil {
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

	switch apdu.GetRootMessage().(type) {
	case readWriteModel.APDUConfirmedRequest:
		s.log.Debug().Msg("client is trying this request again")
	case readWriteModel.APDUAbort:
		s.log.Debug().Msg("client aborting this request")

		// forward to the application
		if err := s.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Request(NA(apdu), NoKWArgs()); err != nil { // send it ot the device
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
	if err := s.Request(NA(abort), NoKWArgs()); err != nil {
		s.log.Debug().Err(err).Msg("error sending request")
	}
	return nil
}

func (s *ServerSSM) segmentedResponse(apdu PDU) error {
	s.log.Debug().Stringer("apdu", apdu).Msg("segmentedResponse")

	// client is ready for the next segment
	switch _apdu := apdu.GetRootMessage().(type) {
	case readWriteModel.APDUSegmentAck:
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
	case readWriteModel.APDUAbort:
		if err := s.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "Error setting state to aborted")
		}
		if err := s.Response(NA(apdu), NoKWArgs()); err != nil { // send it ot the application
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
