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
	"math"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Reader struct {
	invokeIdGenerator *InvokeIdGenerator
	messageCodec      spi.MessageCodec
	tm                *spi.RequestTransactionManager

	// TODO make them configurable
	protocolVersion       uint8
	hopCount              uint8
	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted
	srcAddress            *struct {
		NetworkAddress uint16
		Address        []byte
	}
	dstAddress struct {
		NetworkAddress uint16
		Address        []byte
	}
}

func NewReader(invokeIdGenerator *InvokeIdGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) *Reader {
	return &Reader{
		invokeIdGenerator: invokeIdGenerator,
		messageCodec:      messageCodec,
		tm:                tm,

		protocolVersion:       1,
		hopCount:              255,
		maxSegmentsAccepted:   readWriteModel.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS,
		maxApduLengthAccepted: readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult)
	go func() {
		if len(readRequest.GetTagNames()) == 0 {
			result <- &spiModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("at least one field required"),
			}
			return
		}
		// create the service request
		var serviceRequest readWriteModel.BACnetConfirmedServiceRequest
		quantity := uint32(1)
		if len(readRequest.GetTag(readRequest.GetTagNames()[0]).GetArrayInfo()) > 0 {
			quantity = readRequest.GetTag(readRequest.GetTagNames()[0]).GetArrayInfo()[0].GetUpperBound() - readRequest.GetTag(readRequest.GetTagNames()[0]).GetArrayInfo()[0].GetLowerBound()
		}
		if isMultiRequest := len(readRequest.GetTagNames()) > 1 || quantity > 1; !isMultiRequest {
			// Single request
			singleTag := readRequest.GetTag(readRequest.GetTagNames()[0]).(BacNetPlcTag)
			objectIdentifier := readWriteModel.CreateBACnetContextTagObjectIdentifier(0, singleTag.GetObjectId().getId(), singleTag.GetObjectId().ObjectIdInstance)
			propertyIdentifier := readWriteModel.CreateBACnetPropertyIdentifierTagged(1, singleTag.GetProperties()[0].getId())
			var arrayIndex readWriteModel.BACnetContextTagUnsignedInteger
			if value := singleTag.GetProperties()[0].ArrayIndex; value != nil {
				arrayIndex = readWriteModel.CreateBACnetContextTagUnsignedInteger(2, *value)
			}
			serviceRequest = readWriteModel.NewBACnetConfirmedServiceRequestReadProperty(objectIdentifier, propertyIdentifier, arrayIndex, 0)
		} else {
			// Multi request
			var data []readWriteModel.BACnetReadAccessSpecification
			for _, tagName := range readRequest.GetTagNames() {
				tag := readRequest.GetTag(tagName).(BacNetPlcTag)
				objectIdentifier := readWriteModel.CreateBACnetContextTagObjectIdentifier(0, tag.GetObjectId().getId(), tag.GetObjectId().ObjectIdInstance)
				var listOfPropertyReferences []readWriteModel.BACnetPropertyReference
				for _, _property := range tag.GetProperties() {
					propertyIdentifier := readWriteModel.CreateBACnetPropertyIdentifierTagged(0, _property.getId())
					var arrayIndex readWriteModel.BACnetContextTagUnsignedInteger
					if value := _property.ArrayIndex; value != nil {
						arrayIndex = readWriteModel.CreateBACnetContextTagUnsignedInteger(1, *value)
					}
					listOfPropertyReferences = append(listOfPropertyReferences, readWriteModel.NewBACnetPropertyReference(propertyIdentifier, arrayIndex))
				}
				specification := readWriteModel.NewBACnetReadAccessSpecification(
					objectIdentifier,
					readWriteModel.CreateBACnetOpeningTag(1),
					listOfPropertyReferences,
					readWriteModel.CreateBACnetClosingTag(1),
				)
				data = append(data, specification)
			}

			serviceRequest = readWriteModel.NewBACnetConfirmedServiceRequestReadPropertyMultiple(data, 0, 0)
		}

		invokeId := m.invokeIdGenerator.getAndIncrement()

		// build apdu
		apdu := readWriteModel.NewAPDUConfirmedRequest(
			false,
			false,
			true,
			m.maxSegmentsAccepted,
			m.maxApduLengthAccepted,
			invokeId,
			nil,
			nil,
			serviceRequest,
			nil,
			nil,
			serviceRequest.GetLengthInBytes(),
		)

		// build npdu
		sourceSpecified := m.srcAddress != nil
		var sourceNetworkAddress *uint16
		var sourceLength *uint8
		var sourceAddress []uint8
		if sourceSpecified {
			sourceSpecified = true
			sourceNetworkAddress = &m.srcAddress.NetworkAddress
			sourceLengthValue := len(m.srcAddress.Address)
			if sourceLengthValue > math.MaxUint8 {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.New("source address length overflows"),
				}
				return
			}
			sourceLengthValueUint8 := uint8(sourceLengthValue)
			sourceLength = &sourceLengthValueUint8
			sourceAddress = m.srcAddress.Address
			if sourceLengthValueUint8 == 0 {
				// If we define the len 0 we must not send the array
				sourceAddress = nil
			}
		}
		control := readWriteModel.NewNPDUControl(false, true, sourceSpecified, true, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
		destinationLengthValue := len(m.dstAddress.Address)
		if destinationLengthValue > math.MaxUint8 {
			result <- &spiModel.DefaultPlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("destination address length overflows"),
			}
			return
		}
		destinationNetworkAddress := &m.dstAddress.NetworkAddress
		destinationLengthValueUint8 := uint8(destinationLengthValue)
		destinationtLength := &destinationLengthValueUint8
		destinationAddress := m.dstAddress.Address
		if len(m.dstAddress.Address) == 0 {
			// If we define the len 0 we must not send the array
			destinationAddress = nil
		}
		npdu := readWriteModel.NewNPDU(m.protocolVersion, control, destinationNetworkAddress, destinationtLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, &m.hopCount, nil, apdu, 0)
		bvlc := readWriteModel.NewBVLCOriginalUnicastNPDU(npdu, 0)
		// Start a new request-transaction (Is ended in the response-handler)
		transaction := m.tm.StartTransaction()
		transaction.Submit(func() {

			// Send the  over the wire
			log.Trace().Msg("Send ")
			if err := m.messageCodec.SendRequest(ctx, bvlc, func(message spi.Message) bool {
				bvlc, ok := message.(readWriteModel.BVLCExactly)
				if !ok {
					log.Debug().Msgf("Received strange type %T", bvlc)
					return false
				}
				var npdu readWriteModel.NPDU
				if npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU }); ok {
					npdu = npduRetriever.GetNpdu()
				} else {
					log.Debug().Msgf("bvlc has no way to give a npdu %T", bvlc)
					return false
				}
				if npdu.GetControl().GetMessageTypeFieldPresent() {
					return false
				}
				if invokeIdFromApdu, err := getInvokeIdFromApdu(npdu.GetApdu()); err != nil {
					log.Debug().Err(err).Msg("Error getting invoke id")
					return false
				} else {
					return invokeIdFromApdu == invokeId
				}
			}, func(message spi.Message) error {
				// Convert the response into an
				log.Trace().Msg("convert response to ")
				apdu := message.(readWriteModel.BVLC).(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu().GetApdu()

				// TODO: implement segment handling

				// Convert the bacnet response into a PLC4X response
				log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(apdu, readRequest)

				if err != nil {
					result <- &spiModel.DefaultPlcReadRequestResult{
						Request: readRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
					return transaction.EndRequest()
				}
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: readResponse,
				}
				return transaction.EndRequest()
			}, func(err error) error {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request: readRequest,
					Err:     errors.Wrap(err, "got timeout while waiting for response"),
				}
				return transaction.EndRequest()
			}, time.Second*1); err != nil {
				result <- &spiModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrap(err, "error sending message"),
				}
				_ = transaction.EndRequest()
			}
		})
	}()
	return result
}

func (m *Reader) ToPlc4xReadResponse(apdu readWriteModel.APDU, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	var complexAck readWriteModel.APDUComplexAck
	var errorClass *readWriteModel.ErrorClass
	var errorCode *readWriteModel.ErrorCode
	var rejectReason *readWriteModel.BACnetRejectReason
	var abortReason *readWriteModel.BACnetAbortReason
	switch apdu := apdu.(type) {
	case readWriteModel.APDUComplexAck:
		complexAck = apdu
	case readWriteModel.APDUError:
		apduError := apdu.GetError()
		var bacError readWriteModel.Error
		switch concreteError := apduError.(type) {
		case readWriteModel.BACnetErrorGeneral:
			bacError = concreteError.GetError()
		default:
			bacError = concreteError.(interface {
				GetErrorType() readWriteModel.ErrorEnclosed
			}).GetErrorType().GetError()
		}
		errorClassValue := bacError.GetErrorClass().GetValue()
		errorClass = &errorClassValue
		errorCodeValue := bacError.GetErrorCode().GetValue()
		errorCode = &errorCodeValue
	case readWriteModel.APDUReject:
		rejectReasonValue := apdu.GetRejectReason().GetValue()
		rejectReason = &rejectReasonValue
	case readWriteModel.APDUAbort:
		abortReasonValue := apdu.GetAbortReason().GetValue()
		abortReason = &abortReasonValue
	default:
		return nil, errors.Errorf("unsupported response type %T", apdu)
	}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}

	// If the result contains any form of non-null error code, handle this instead.
	if errorClass != nil {
		log.Warn().Msgf("Got an unknown error response from the PLC. Error Class: %d, Error Code %d. "+
			"We probably need to implement explicit handling for this, so please file a bug-report "+
			"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
			"containing a capture of the communication.",
			errorClass, errorCode)
		for _, tagName := range readRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			plcValues[tagName] = spiValues.NewPlcNULL()
		}
		return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
	}
	if rejectReason != nil {
		log.Warn().Msgf("Got an unknown error response from the PLC. Reject Reason %d. "+
			"We probably need to implement explicit handling for this, so please file a bug-report "+
			"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
			"containing a capture of the communication.",
			rejectReason)
		for _, tagName := range readRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			plcValues[tagName] = spiValues.NewPlcNULL()
		}
		return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
	}
	if abortReason != nil {
		log.Warn().Msgf("Got an unknown error response from the PLC. Abort Reason %d. "+
			"We probably need to implement explicit handling for this, so please file a bug-report "+
			"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
			"containing a capture of the communication.",
			abortReason)
		for _, tagName := range readRequest.GetTagNames() {
			responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			plcValues[tagName] = spiValues.NewPlcNULL()
		}
		return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
	}

	switch complexAck := complexAck.(type) {
	case readWriteModel.BACnetServiceAckReadPropertyExactly:
		// TODO: super lazy implementation for now
		responseCodes[readRequest.GetTagNames()[0]] = apiModel.PlcResponseCode_OK
		plcValues[readRequest.GetTagNames()[0]] = spiValues.NewPlcSTRING(complexAck.GetValues().(fmt.Stringer).String())
	case readWriteModel.BACnetServiceAckReadPropertyMultipleExactly:

		// way to know how to interpret the responses is by aligning them with the
		// items from the request as this information is not returned by the PLC.
		if len(readRequest.GetTagNames()) != len(complexAck.GetData()) {
			return nil, errors.New("The number of requested items doesn't match the number of returned items")
		}
		for i, tagName := range readRequest.GetTagNames() {
			// TODO: super lazy implementation for now
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
			plcValues[tagName] = spiValues.NewPlcSTRING(complexAck.GetData()[i].GetListOfResults().(fmt.Stringer).String())
		}
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

func getInvokeIdFromApdu(apdu readWriteModel.APDU) (uint8, error) {
	var invokeId uint8
	if originalInvokeIdRetriever, ok := apdu.(interface {
		GetOriginalInvokeId() uint8
	}); ok {
		invokeId = originalInvokeIdRetriever.GetOriginalInvokeId()
	} else if invokeIdRetriever, ok := apdu.(interface {
		GetInvokeId() uint8
	}); ok {
		invokeId = invokeIdRetriever.GetInvokeId()
	} else {
		return 0, errors.Errorf("No way to get invoke id from %T", apdu)
	}
	return invokeId, nil
}
