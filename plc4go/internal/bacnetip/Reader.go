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
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Reader struct {
	invokeIdGenerator *InvokeIdGenerator
	messageCodec      spi.MessageCodec
	tm                transactions.RequestTransactionManager
	driverContext     DriverContext
	routedDest        *routedDestination // nil for local-segment connections

	maxSegmentsAccepted   readWriteModel.MaxSegmentsAccepted
	maxApduLengthAccepted readWriteModel.MaxApduLengthAccepted

	wg sync.WaitGroup // use to track spawned go routines

	log zerolog.Logger
}

func NewReader(invokeIdGenerator *InvokeIdGenerator, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, driverContext DriverContext, routedDest *routedDestination, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		invokeIdGenerator: invokeIdGenerator,
		messageCodec:      messageCodec,
		tm:                tm,
		driverContext:     driverContext,
		routedDest:        routedDest,

		maxSegmentsAccepted:   readWriteModel.MaxSegmentsAccepted_MORE_THAN_64_SEGMENTS,
		maxApduLengthAccepted: readWriteModel.MaxApduLengthAccepted_NUM_OCTETS_1476,

		log: customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	m.wg.Go(func() {
		if len(readRequest.GetTagNames()) == 0 {
			utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.New("at least one field required")))
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
			serviceRequest = readWriteModel.NewBACnetConfirmedServiceRequestReadProperty(0, objectIdentifier, propertyIdentifier, arrayIndex)
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

			serviceRequest = readWriteModel.NewBACnetConfirmedServiceRequestReadPropertyMultiple(0, data)
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
		)

		// If the request exceeds the peer's declared APDU ceiling (large
		// ReadPropertyMultiple access lists), it has to go out as a segmented
		// request (ASHRAE 135 clause 5.4) — or fail fast when the peer can't
		// receive segments, instead of provoking an abort.
		var segmentedPayload []byte
		if m.driverContext.peerMaxApduBytes > 0 {
			payload, serErr := serviceRequest.Serialize()
			if serErr != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrap(serErr, "Error serializing read request")))
				return
			}
			if m.driverContext.needsSegmentedRequest(len(payload)) {
				if !m.driverContext.peerAcceptsSegmentedRequests {
					utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil,
						errors.Errorf("read request of %d bytes exceeds the peer's max APDU of %d and the peer does not support segmented requests", len(payload), m.driverContext.peerMaxApduBytes)))
					return
				}
				segmentedPayload = payload
			}
		}

		// Start a new request-transaction (Is ended in the response-handler)
		transaction := m.tm.StartTransaction("read")
		transaction.Submit("readOperation", func(transactionContext context.Context, transaction transactions.RequestTransaction) {
			// readCtx is the caller's context; it outlives the transaction so a
			// multi-segment reassembly (driven after we end the transaction) can
			// keep awaiting segments without being cancelled by transaction teardown.
			readCtx := ctx
			ctx, cancel := context.WithCancel(ctx)
			context.AfterFunc(transactionContext, cancel)

			acceptsReadResponse := func(message spi.Message) bool {
				bvlc, ok := message.(readWriteModel.BVLC)
				if !ok {
					m.log.Debug().Type("bvlc", bvlc).Msg("Received strange type")
					return false
				}
				var npdu readWriteModel.NPDU
				if npduRetriever, ok := bvlc.(interface{ GetNpdu() readWriteModel.NPDU }); ok {
					npdu = npduRetriever.GetNpdu()
				} else {
					m.log.Debug().Type("bvlc", bvlc).Msg("bvlc has no way to give a npdu")
					return false
				}
				if npdu.GetControl().GetMessageTypeFieldPresent() {
					return false
				}
				if invokeIdFromApdu, err := getInvokeIdFromApdu(npdu.GetApdu()); err != nil {
					m.log.Debug().Err(err).Msg("Error getting invoke id")
					return false
				} else {
					return invokeIdFromApdu == invokeId
				}
			}
			handleReadResponse := func(message spi.Message) error {
				// Convert the response into an
				m.log.Trace().Msg("convert response to ")
				apdu := message.(readWriteModel.BVLC).(interface{ GetNpdu() readWriteModel.NPDU }).GetNpdu().GetApdu()

				// Segmented response: the device split the service ack across
				// multiple APDUs. We can't drive the segment-ack/await loop here
				// because this callback runs while the codec holds its expectation
				// lock — instead hand off to a goroutine that uses the caller's
				// context, and free this transaction slot immediately.
				if complexAck, ok := apdu.(readWriteModel.APDUComplexAck); ok && complexAck.GetSegmentedMessage() {
					m.log.Trace().Uint8("invokeId", complexAck.GetOriginalInvokeId()).Msg("segmented read response — starting reassembly")
					m.wg.Go(func() {
						// Reassembly works on wire-controlled sizes; a panic in this
						// goroutine would otherwise kill the whole process, so recover
						// and fail the read instead.
						defer func() {
							if r := recover(); r != nil {
								utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", r, debug.Stack())))
							}
						}()
						m.reassembleSegmentedRead(readCtx, readRequest, complexAck, result)
					})
					return transaction.EndRequest()
				}

				// Convert the bacnet response into a PLC4X response
				m.log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(apdu, readRequest)

				if err != nil {
					utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
						readRequest,
						nil,
						errors.Wrap(err, "Error decoding response"),
					))
					return transaction.EndRequest()
				}
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					readResponse,
					nil,
				))
				return transaction.EndRequest()
			}
			handleReadError := func(err error) error {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrap(err, "got timeout while waiting for response"),
				))
				return transaction.EndRequest()
			}

			if segmentedPayload != nil {
				// Register the response expectation BEFORE driving the segments:
				// the peer may answer right after acking the final segment. The
				// matcher must not consume the peer's SegmentAcks — those belong
				// to the sender's own expectations.
				m.messageCodec.Expect(ctx, "segmentedReadResponse",
					responseMatcherExcludingSegmentAcks(acceptsReadResponse),
					handleReadResponse,
					handleReadError,
				)

				sender := &segmentedRequestSender{
					messageCodec:  m.messageCodec,
					routedDest:    m.routedDest,
					driverContext: m.driverContext,
					log:           m.log,
				}
				if err := sender.send(ctx, invokeId, segmentedPayload); err != nil {
					// Deliver the real cause first (the buffered result channel
					// takes it), then cancel so the armed response expectation
					// resolves without delivering a second, less specific error.
					utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
						readRequest,
						nil,
						errors.Wrap(err, "error sending segmented read request"),
					))
					cancel()
				}
				return
			}

			// Send the  over the wire
			m.log.Trace().Msg("Send ")
			if err := m.messageCodec.SendRequest(ctx, "read", wrapAPDU(apdu, true, m.routedDest), acceptsReadResponse, handleReadResponse, handleReadError); err != nil {
				utils.DeliverResult(m.log, result, spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrap(err, "error sending message"),
				))
				if err := transaction.FailRequest(errors.Errorf("timeout after %s", time.Second*1)); err != nil {
					m.log.Debug().Err(err).Msg("Error failing request")
				}
			}
		})
	})
	return result
}

func (m *Reader) ToPlc4xReadResponse(apdu readWriteModel.APDU, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	tagNames := readRequest.GetTagNames()
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}

	// Dispatch on APDU subtype. Errors/Reject/Abort fan out to every requested
	// tag with the same response code so consumers can read a sane PlcResponse
	// even on protocol-level failures.
	switch apdu := apdu.(type) {
	case readWriteModel.APDUComplexAck:
		return m.decodeComplexAck(apdu, readRequest)
	case readWriteModel.APDUError:
		responseCode := mapErrorAPDU(apdu, m.log)
		return broadcastResponseCode(readRequest, responseCode, responseCodes, plcValues), nil
	case readWriteModel.APDUReject:
		m.log.Warn().Stringer("rejectReason", reasonOfReject(apdu)).Msg("BACnet REJECT received")
		return broadcastResponseCode(readRequest, apiModel.PlcResponseCode_INVALID_DATA, responseCodes, plcValues), nil
	case readWriteModel.APDUAbort:
		reason := apdu.GetAbortReason().GetValue()
		code := apiModel.PlcResponseCode_INTERNAL_ERROR
		if reason == readWriteModel.BACnetAbortReason_SEGMENTATION_NOT_SUPPORTED {
			code = apiModel.PlcResponseCode_UNSUPPORTED
		}
		m.log.Warn().Stringer("abortReason", &reason).Msg("BACnet ABORT received")
		return broadcastResponseCode(readRequest, code, responseCodes, plcValues), nil
	default:
		// Reject any other APDU type — but mark all tags as remote-error so the
		// caller still receives a well-formed PlcReadResponse instead of nil.
		_ = tagNames
		m.log.Warn().Type("apdu", apdu).Msg("unsupported APDU type on read response")
		return broadcastResponseCode(readRequest, apiModel.PlcResponseCode_REMOTE_ERROR, responseCodes, plcValues), nil
	}
}

// decodeComplexAck handles a non-segmented APDUComplexAck for ReadProperty and
// ReadPropertyMultiple. Segmented responses are reassembled separately (see
// reassembleSegmentedRead) and then handed to decodeServiceAck directly.
func (m *Reader) decodeComplexAck(apdu readWriteModel.APDUComplexAck, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	if apdu.GetSegmentedMessage() {
		// Should not happen: segmented responses are intercepted before this
		// point. Surface UNSUPPORTED rather than panicking on a nil ServiceAck.
		m.log.Warn().Uint8("invokeId", apdu.GetOriginalInvokeId()).Msg("segmented APDU reached decodeComplexAck unexpectedly")
		return broadcastResponseCode(readRequest, apiModel.PlcResponseCode_UNSUPPORTED, map[string]apiModel.PlcResponseCode{}, map[string]values.PlcValue{}), nil
	}
	return m.decodeServiceAck(apdu.GetServiceAck(), readRequest)
}

// decodeServiceAck converts a decoded BACnetServiceAck (whether it arrived in a
// single APDU or was reassembled from segments) into a PLC4X read response.
func (m *Reader) decodeServiceAck(serviceAckMsg readWriteModel.BACnetServiceAck, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	tagNames := readRequest.GetTagNames()
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}

	switch serviceAck := serviceAckMsg.(type) {
	case readWriteModel.BACnetServiceAckReadProperty:
		if len(tagNames) == 0 {
			return nil, errors.New("ReadProperty response without a corresponding requested tag")
		}
		responseCodes[tagNames[0]] = apiModel.PlcResponseCode_OK
		plcValues[tagNames[0]] = constructedDataToPlcValue(serviceAck.GetValues())
		// Any extra requested tags (shouldn't happen with single ReadProperty,
		// but be defensive) get NOT_FOUND.
		for _, n := range tagNames[1:] {
			responseCodes[n] = apiModel.PlcResponseCode_NOT_FOUND
			plcValues[n] = spiValues.NewPlcNULL()
		}
	case readWriteModel.BACnetServiceAckReadPropertyMultiple:
		data := serviceAck.GetData()
		if len(tagNames) != len(data) {
			return nil, errors.Errorf("ReadPropertyMultiple expected %d results, got %d", len(tagNames), len(data))
		}
		for i, tagName := range tagNames {
			list := data[i].GetListOfResults()
			if list == nil {
				responseCodes[tagName] = apiModel.PlcResponseCode_NOT_FOUND
				plcValues[tagName] = spiValues.NewPlcNULL()
				continue
			}
			// Per BACnet ReadPropertyMultiple semantics each requested property
			// in the access spec produces one read result. We collapse to a
			// single value if there's just one, otherwise wrap as PlcList.
			results := list.GetListOfReadAccessProperty()
			code, value := readResultsToPlcValue(results)
			responseCodes[tagName] = code
			plcValues[tagName] = value
		}
	default:
		m.log.Warn().Type("serviceAck", serviceAck).Msg("unsupported ServiceAck variant")
		return broadcastResponseCode(readRequest, apiModel.PlcResponseCode_INTERNAL_ERROR, responseCodes, plcValues), nil
	}

	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

// readResultsToPlcValue collapses a ReadPropertyMultiple result list to a
// (responseCode, value) pair. Errors in any one element propagate to the whole
// tag — BACnet doesn't promise atomicity but we fail closed for clarity.
func readResultsToPlcValue(results []readWriteModel.BACnetReadAccessProperty) (apiModel.PlcResponseCode, values.PlcValue) {
	if len(results) == 0 {
		return apiModel.PlcResponseCode_NOT_FOUND, spiValues.NewPlcNULL()
	}
	collected := make([]values.PlcValue, 0, len(results))
	for _, rap := range results {
		rr := rap.GetReadResult()
		if rr == nil {
			collected = append(collected, spiValues.NewPlcNULL())
			continue
		}
		if e := rr.GetPropertyAccessError(); e != nil {
			code := bacnetErrorEnclosedToResponseCode(e)
			return code, spiValues.NewPlcNULL()
		}
		collected = append(collected, constructedDataToPlcValue(rr.GetPropertyValue()))
	}
	if len(collected) == 1 {
		return apiModel.PlcResponseCode_OK, collected[0]
	}
	return apiModel.PlcResponseCode_OK, spiValues.NewPlcList(collected)
}

// mapErrorAPDU translates an APDUError into the appropriate plc4go response code
// based on BACnet ErrorClass + ErrorCode. Unknown combinations default to
// REMOTE_ERROR so callers always get a meaningful non-OK result.
func mapErrorAPDU(apdu readWriteModel.APDUError, log zerolog.Logger) apiModel.PlcResponseCode {
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
	errorClass := bacError.GetErrorClass().GetValue()
	errorCode := bacError.GetErrorCode().GetValue()
	log.Warn().
		Stringer("errorClass", &errorClass).
		Stringer("errorCode", &errorCode).
		Msg("BACnet error response received")
	return mapErrorClassCodeToResponseCode(errorClass, errorCode)
}

// mapErrorClassCodeToResponseCode is the table the plan calls out: class+code
// to plc4go's flat response codes.
func mapErrorClassCodeToResponseCode(class readWriteModel.ErrorClass, code readWriteModel.ErrorCode) apiModel.PlcResponseCode {
	switch class {
	case readWriteModel.ErrorClass_OBJECT:
		if code == readWriteModel.ErrorCode_UNKNOWN_OBJECT {
			return apiModel.PlcResponseCode_NOT_FOUND
		}
	case readWriteModel.ErrorClass_PROPERTY:
		switch code {
		case readWriteModel.ErrorCode_UNKNOWN_PROPERTY:
			return apiModel.PlcResponseCode_INVALID_ADDRESS
		case readWriteModel.ErrorCode_WRITE_ACCESS_DENIED:
			return apiModel.PlcResponseCode_ACCESS_DENIED
		}
	case readWriteModel.ErrorClass_DEVICE:
		if code == readWriteModel.ErrorCode_OPERATIONAL_PROBLEM {
			return apiModel.PlcResponseCode_REMOTE_ERROR
		}
	}
	return apiModel.PlcResponseCode_REMOTE_ERROR
}

// bacnetErrorEnclosedToResponseCode extracts the inner Error from a
// per-property ErrorEnclosed (used in ReadPropertyMultiple results) and maps it
// to a plc4go response code.
func bacnetErrorEnclosedToResponseCode(enclosed readWriteModel.ErrorEnclosed) apiModel.PlcResponseCode {
	err := enclosed.GetError()
	return mapErrorClassCodeToResponseCode(err.GetErrorClass().GetValue(), err.GetErrorCode().GetValue())
}

// broadcastResponseCode populates the same code and a PlcNULL value for every
// tag in the request, then wraps in a DefaultPlcReadResponse.
func broadcastResponseCode(req apiModel.PlcReadRequest, code apiModel.PlcResponseCode, codes map[string]apiModel.PlcResponseCode, vals map[string]values.PlcValue) apiModel.PlcReadResponse {
	for _, name := range req.GetTagNames() {
		codes[name] = code
		vals[name] = spiValues.NewPlcNULL()
	}
	return spiModel.NewDefaultPlcReadResponse(req, codes, vals)
}

// reasonOfReject returns a pointer to the reject reason for Stringer formatting.
func reasonOfReject(apdu readWriteModel.APDUReject) *readWriteModel.BACnetRejectReason {
	r := apdu.GetRejectReason().GetValue()
	return &r
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
