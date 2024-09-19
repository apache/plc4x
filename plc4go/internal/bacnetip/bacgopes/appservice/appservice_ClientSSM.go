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

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/apdu"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/pdu"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

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
	if _debug != nil {
		_debug("__init__ %s %r", sap, pduAddress)
	}
	ssm, err := NewSSM(localLog, struct {
		SSMSAPRequirements
		SSMProcessingRequirements
	}{sap, c}, pduAddress)
	if err != nil {
		return nil, err
	}
	// acquire the device info
	if ssm.deviceInfo != nil {
		if _debug != nil {
			_debug("    - acquire device information")
		}
		localLog.Debug().Msg("Accquire device information")
		c.ssmSAP.GetDeviceInfoCache().Acquire(c.deviceInfo._cacheKey)
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
func (c *ClientSSM) Request(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("request")
	apdu := GA[PDU](args, 0)

	// make sure it has a good source and destination
	apdu.SetPDUSource(nil)
	apdu.SetPDUDestination(c.pduAddress)

	// send it via the device
	return c.ssmSAP.Request(NA(apdu), kwArgs)
}

// Indication This function is called after the device has bound a new transaction and wants to start the process
//
//	rolling
func (c *ClientSSM) Indication(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("indication")
	apdu := GA[APDU](args, 0)
	// make sure we're getting confirmed requests
	var apduConfirmedRequest readWriteModel.APDUConfirmedRequest
	if apduCasted, ok := apdu.GetRootMessage().(readWriteModel.APDUConfirmedRequest); !ok {
		return errors.Errorf("Invalid APDU type %T", apduCasted)
	} else {
		apduConfirmedRequest = apduCasted
	}

	// save the request and set the segmentation context
	if err := c.setSegmentationContext(apduConfirmedRequest); err != nil {
		return errors.Wrap(err, "error setting context")
	}

	// if the max apdu length of the server isn't known, assume that it is the same size as our own and will be the segment
	// size
	if c.deviceInfo == nil || c.deviceInfo.MaximumApduLengthAccepted != nil {
		c.segmentSize = uint(c.maxApduLengthAccepted.NumberOfOctets())
	} else if c.deviceInfo.MaximumNpduLength == nil {
		// if the max npdu length of the server isn't known, assume that it is the same as the max apdu length accepted
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
			return c.Response(NA(abort), kwArgs)
		}

		if c.deviceInfo == nil {
			c.log.Debug().Msg("no server info for segmentation support")
		} else if *c.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_TRANSMIT && *c.deviceInfo.SegmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			c.log.Debug().Msg("server can't receive segmented requests")
			abort, err := c.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "Error creating abort")
			}
			return c.Response(NA(abort), kwArgs)
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
			return c.Response(NA(abort), kwArgs)
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
	return c.Request(NA(segment), kwArgs)
}

// Response This function is called by client transaction functions when they want to send a message to the application.
func (c *ClientSSM) Response(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("response")

	apdu := GA[PDU](args, 0)

	// make sure it has a good source and destination
	apdu.SetPDUSource(nil)
	apdu.SetPDUDestination(c.pduAddress)

	// send it to the application
	return c.ssmSAP.SapResponse(NA(apdu), kwArgs)
}

// Confirmation This function is called by the device for all upstream messages related to the transaction.
func (c *ClientSSM) Confirmation(args Args, kwArgs KWArgs) error {
	c.log.Debug().Stringer("Args", args).Stringer("KWArgs", kwArgs).Msg("confirmation")
	apdu := GA[PDU](args, 0)

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
	return NewPDU(NoArgs, NoKWArgs(), WithRootMessage(abortApdu)), nil
}

// segmentedRequest This function is called when the client is sending a segmented request and receives an apdu
func (c *ClientSSM) segmentedRequest(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("segmentedRequest")

	switch _apdu := apdu.GetRootMessage().(type) {
	// server is ready for the next segment
	case readWriteModel.APDUSegmentAck:
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
	case readWriteModel.APDUSimpleAck:
		c.log.Debug().Msg("simple ack")

		if !c.sentAllSegments {
			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else {
			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
		}
	// complex ack
	case readWriteModel.APDUComplexAck:
		c.log.Debug().Msg("complex ack")
		if !c.sentAllSegments {
			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else if !_apdu.GetSegmentedMessage() {
			// ack is not segmented
			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := c.Response(NA(apdu), NoKWArgs()); err != nil {
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
	case readWriteModel.APDUError:
		c.log.Debug().Msg("error/reject/abort")
		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NA(apdu), NoKWArgs()); err != nil {
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
			if err := c.Request(NA(apdu), NoKWArgs()); err != nil {
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
		if err := c.Response(NA(abort), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}
	return nil
}

func (c *ClientSSM) awaitConfirmation(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("awaitConfirmation")

	switch _apdu := apdu.GetRootMessage().(type) {
	case readWriteModel.APDUAbort:
		c.log.Debug().Msg("Server aborted")

		if err := c.setState(SSMState_ABORTED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NA(apdu), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUSimpleAck, readWriteModel.APDUError, readWriteModel.APDUReject:
		c.log.Debug().Msg("simple ack, error or reject")

		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		if err := c.Response(NA(apdu), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	case readWriteModel.APDUComplexAck:
		c.log.Debug().Msg("complex ack")

		if !_apdu.GetSegmentedMessage() {
			c.log.Debug().Msg("unsegmented")

			if err := c.setState(SSMState_COMPLETED, nil); err != nil {
				return errors.Wrap(err, "error switching state")
			}
			if err := c.Response(NA(apdu), NoKWArgs()); err != nil {
				c.log.Debug().Err(err).Msg("error sending response")
			}
		} else if c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_RECEIVE && c.segmentationSupported != readWriteModel.BACnetSegmentation_SEGMENTED_BOTH {
			c.log.Debug().Msg("local device can't receive segmented messages")

			abort, err := c.abort(readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Response(NA(abort), NoKWArgs()); err != nil {
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
			if err := c.Request(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segmentAck))), NoKWArgs()); err != nil {
				c.log.Debug().Err(err).Msg("error sending request")
			}
		} else {
			c.log.Debug().Msg("Invalid apdu in this state")

			abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
			if err != nil {
				return errors.Wrap(err, "error creating abort")
			}
			if err := c.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
				c.log.Debug().Err(err).Msg("error sending request")
			}
			if err := c.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
				c.log.Debug().Err(err).Msg("error sending response")
			}
		}
	case readWriteModel.APDUSegmentAck:
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
		if err := c.Indication(NA(NewPDU(NoArgs, NKW(KWCPCIDestination, c.pduAddress)), WithRootMessage(c.segmentAPDU.originalApdu)), NoKWArgs()); err != nil { // TODO: check that it is really the intention to re-send the original apdu here
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
		if err := c.Response(NA(abort), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}
	return nil
}

func (c *ClientSSM) segmentedConfirmation(apdu PDU) error {
	c.log.Debug().Stringer("apdu", apdu).Msg("segmentedConfirmation")

	// the only messages we should be getting are complex acks
	apduComplexAck, ok := apdu.(readWriteModel.APDUComplexAck)
	if !ok {
		c.log.Debug().Msg("complex ack required")

		abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
			c.log.Debug().Err(err).Msg("error sending request")
		}
		if err := c.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
			c.log.Debug().Err(err).Msg("error sending response")
		}
	}

	// it must be segmented
	if !apduComplexAck.GetSegmentedMessage() {
		abort, err := c.abort(readWriteModel.BACnetAbortReason_INVALID_APDU_IN_THIS_STATE)
		if err != nil {
			return errors.Wrap(err, "error creating abort")
		}
		if err := c.Request(NA(abort), NoKWArgs()); err != nil { // send it ot the device
			c.log.Debug().Err(err).Msg("error sending request")
		}
		if err := c.Response(NA(abort), NoKWArgs()); err != nil { // send it ot the application
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
		if err := c.Request(NA(NewPDU(NoArgs, NoKWArgs()), WithRootMessage(segmentAck)), NoKWArgs()); err != nil {
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
		if err := c.Request(NA(NewPDU(NoArgs, NoKWArgs()), WithRootMessage(segmentAck)), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending request")
		}

		if err := c.setState(SSMState_COMPLETED, nil); err != nil {
			return errors.Wrap(err, "error switching state")
		}
		// TODO: this is nonsense... We need to parse the service and the apdu not sure where to get it from now...
		// TODO: it should be the original apdu, we might just need to use that as base and forward it as non segmented
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), c.log, options.WithPassLoggerToModel(c.passLogToModel))
		parse, err := readWriteModel.APDUParse[readWriteModel.APDU](ctxForModel, c.segmentAPDU.serviceBytes, uint16(len(c.segmentAPDU.serviceBytes)))
		if err != nil {
			return errors.Wrap(err, "error parsing apdu")
		}
		if err := c.Request(NA(NewPDU(NoArgs, NoKWArgs()), WithRootMessage(parse)), NoKWArgs()); err != nil {
			c.log.Debug().Err(err).Msg("error sending response")
		}
	} else if *apduComplexAck.GetSequenceNumber() == c.initialSequenceNumber+*c.actualWindowSize {
		c.log.Debug().Msg("last segment in the group")

		c.initialSequenceNumber = c.lastSequenceNumber
		c.RestartTimer(c.segmentTimeout)
		segmentAck := readWriteModel.NewAPDUSegmentAck(false, false, c.invokeId, c.lastSequenceNumber, *c.actualWindowSize, 0)
		if err := c.Request(NA(NewPDU(NoArgs, NoKWArgs(), WithRootMessage(segmentAck))), NoKWArgs()); err != nil { // send it ot the device
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
	return c.Response(NA(abort), NoKWArgs())
}
