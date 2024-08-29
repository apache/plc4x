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

package s7

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

type Reader struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager

	passLogToModel bool
	log            zerolog.Logger
}

func NewReader(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, _options ...options.WithOption) *Reader {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		tpduGenerator:  tpduGenerator,
		messageCodec:   messageCodec,
		tm:             tm,
		passLogToModel: passLoggerToModel,
		log:            customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()

		requestItems := make([]readWriteModel.S7VarRequestParameterItem, len(readRequest.GetTagNames()))
		for i, tagName := range readRequest.GetTagNames() {
			tag := readRequest.GetTag(tagName)
			address, err := encodeS7Address(tag)
			if err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrapf(err, "Error encoding s7 address for tag %s", tagName),
				)
				return
			}
			requestItems[i] = readWriteModel.NewS7VarRequestParameterItemAddress(address)
		}

		// Create a read request template.
		// tpuId will be inserted before sending in #readInternal so we insert 0 as dummy here
		s7MessageRequest := readWriteModel.NewS7MessageRequest(
			0,
			readWriteModel.NewS7ParameterReadVarRequest(requestItems),
			nil,
		)

		tpduId := m.tpduGenerator.getAndIncrement()

		request := s7MessageRequest
		// Create a new Request with correct tpuId (is not known before)
		s7MessageRequest = readWriteModel.NewS7MessageRequest(tpduId, request.Parameter, request.Payload)

		// Assemble the finished paket
		m.log.Trace().Msg("Assemble paket")
		// TODO: why do we use a uint16 above and the cotp a uint8?
		tpktPacket := readWriteModel.NewTPKTPacket(
			readWriteModel.NewCOTPPacketData(true,
				uint8(tpduId),
				nil,
				s7MessageRequest,
				0,
			),
		)
		// Start a new request-transaction (Is ended in the response-handler)
		transaction := m.tm.StartTransaction()
		transaction.Submit(func(transaction transactions.RequestTransaction) {

			// Send the  over the wire
			m.log.Trace().Msg("Send ")
			if err := m.messageCodec.SendRequest(ctx, tpktPacket, func(message spi.Message) bool {
				tpktPacket, ok := message.(readWriteModel.TPKTPacketExactly)
				if !ok {
					return false
				}
				cotpPacketData, ok := tpktPacket.GetPayload().(readWriteModel.COTPPacketDataExactly)
				if !ok {
					return false
				}
				payload := cotpPacketData.GetPayload()
				if payload == nil {
					return false
				}
				return payload.GetTpduReference() == tpduId
			}, func(message spi.Message) error {
				// Convert the response into an
				m.log.Trace().Msg("convert response to ")
				tpktPacket := message.(readWriteModel.TPKTPacket)
				cotpPacketData := tpktPacket.GetPayload().(readWriteModel.COTPPacketData)
				payload := cotpPacketData.GetPayload()
				// Convert the s7 response into a PLC4X response
				m.log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(payload, readRequest)

				if err != nil {
					result <- spiModel.NewDefaultPlcReadRequestResult(
						readRequest,
						nil,
						errors.Wrap(err, "Error decoding response"),
					)
					return transaction.EndRequest()
				}
				result <- spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					readResponse,
					nil,
				)
				return transaction.EndRequest()
			}, func(err error) error {
				result <- spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrap(err, "got timeout while waiting for response"),
				)
				return transaction.EndRequest()
			}, time.Second*1); err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(
					readRequest,
					nil,
					errors.Wrap(err, "error sending message"),
				)
				if err := transaction.FailRequest(errors.Errorf("timeout after %s", 1*time.Second)); err != nil {
					m.log.Debug().Err(err).Msg("Error failing request")
				}
			}
		})
	}()
	return result
}

func (m *Reader) ToPlc4xReadResponse(response readWriteModel.S7Message, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	var errorClass uint8
	var errorCode uint8
	switch messageResponseData := response.(type) {
	case readWriteModel.S7MessageResponseData:
		errorClass = messageResponseData.GetErrorClass()
		errorCode = messageResponseData.GetErrorCode()
	case readWriteModel.S7MessageResponse:
		errorClass = messageResponseData.GetErrorClass()
		errorCode = messageResponseData.GetErrorCode()
	default:
		return nil, errors.Errorf("unsupported response type %T", response)
	}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	plcValues := map[string]apiValues.PlcValue{}

	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		// This is usually the case if PUT/GET wasn't enabled on the PLC
		if (errorClass == 129) && (errorCode == 4) {
			m.log.Warn().Msg("Got an error response from the PLC. This particular response code usually indicates " +
				"that PUT/GET is not enabled on the PLC.")
			for _, tagName := range readRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_ACCESS_DENIED
				plcValues[tagName] = spiValues.NewPlcNULL()
			}
			m.log.Trace().Msg("Returning the response")
			return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
		} else {
			m.log.Warn().
				Uint8("errorClass", errorClass).
				Uint8("errorCode", errorCode).
				Msg("Got an unknown error response from the PLC. Error Class: %d, Error Code %d. " +
					"We probably need to implement explicit handling for this, so please file a bug-report " +
					"on https://github.com/apache/plc4x/issues and ideally attach a WireShark dump " +
					"containing a capture of the communication.")
			for _, tagName := range readRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
				plcValues[tagName] = spiValues.NewPlcNULL()
			}
			return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
		}
	}

	// In all other cases all went well.
	payload := response.GetPayload().(readWriteModel.S7PayloadReadVarResponse)

	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(readRequest.GetTagNames()) != len(payload.GetItems()) {
		return nil, errors.New("The number of requested items doesn't match the number of returned items")
	}

	payloadItems := payload.GetItems()
	for i, tagName := range readRequest.GetTagNames() {
		tag := readRequest.GetTag(tagName).(PlcTag)
		payloadItem := payloadItems[i]

		responseCode := decodeResponseCode(payloadItem.GetReturnCode())
		// Decode the data according to the information from the request
		m.log.Trace().Msg("decode data")
		responseCodes[tagName] = responseCode
		if responseCode == apiModel.PlcResponseCode_OK {
			ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
			plcValue, err := parsePlcValue(ctxForModel, tag, payloadItem.GetData())
			if err != nil {
				return nil, errors.Wrap(err, "Error parsing data item")
			}
			plcValues[tagName] = plcValue
		}
	}

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

// Currently we only support the S7 Any type of addresses. This helper simply converts the S7Tag from PLC4X into
// S7Address objects.
func encodeS7Address(tag apiModel.PlcTag) (readWriteModel.S7Address, error) {
	s7Tag, ok := tag.(PlcTag)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", tag)
	}
	transportSize := s7Tag.GetDataType()
	numElements := s7Tag.GetNumElements()
	// For these date-types we have to convert the requests to simple byte-array requests
	// As otherwise the S7 will deny them with "Data type not supported" replies.
	if (transportSize == readWriteModel.TransportSize_TIME) /*|| (transportSize == TransportSize.S7_S5TIME)*/ ||
		(transportSize == readWriteModel.TransportSize_LTIME) || (transportSize == readWriteModel.TransportSize_DATE) ||
		(transportSize == readWriteModel.TransportSize_TIME_OF_DAY) || (transportSize == readWriteModel.TransportSize_DATE_AND_TIME) {
		numElements = numElements * uint16(transportSize.SizeInBytes())
		transportSize = readWriteModel.TransportSize_BYTE
	}
	if transportSize == readWriteModel.TransportSize_STRING {
		transportSize = readWriteModel.TransportSize_CHAR
		stringLength := uint16(254)
		if s7StringTag, ok := tag.(PlcStringTag); ok {
			stringLength = s7StringTag.stringLength
		}
		numElements = numElements * (stringLength + 2)
	} else if transportSize == readWriteModel.TransportSize_WSTRING {
		transportSize = readWriteModel.TransportSize_CHAR
		stringLength := uint16(254)
		if s7StringTag, ok := tag.(PlcStringTag); ok {
			stringLength = s7StringTag.stringLength
		}
		numElements = numElements * (stringLength + 2) * 2
	}
	return readWriteModel.NewS7AddressAny(
		transportSize,
		numElements,
		s7Tag.GetBlockNumber(),
		s7Tag.GetMemoryArea(),
		s7Tag.GetByteOffset(),
		s7Tag.GetBitOffset(),
	), nil
}

func parsePlcValue(ctx context.Context, tag PlcTag, data []byte) (apiValues.PlcValue, error) {
	// TODO: port over
	panic("not implemented yet")
}

// Helper to convert the return codes returned from the S7 into one of our standard
func decodeResponseCode(dataTransportErrorCode readWriteModel.DataTransportErrorCode) apiModel.PlcResponseCode {
	switch dataTransportErrorCode {
	case readWriteModel.DataTransportErrorCode_OK:
		return apiModel.PlcResponseCode_OK
	case readWriteModel.DataTransportErrorCode_NOT_FOUND:
		return apiModel.PlcResponseCode_NOT_FOUND
	case readWriteModel.DataTransportErrorCode_INVALID_ADDRESS:
		return apiModel.PlcResponseCode_INVALID_ADDRESS
	case readWriteModel.DataTransportErrorCode_DATA_TYPE_NOT_SUPPORTED:
		return apiModel.PlcResponseCode_INVALID_DATATYPE
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
}
