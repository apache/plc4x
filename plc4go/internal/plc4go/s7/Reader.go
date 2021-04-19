//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package s7

import (
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Reader struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            *spi.RequestTransactionManager
}

func NewReader(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) *Reader {
	return &Reader{
		tpduGenerator: tpduGenerator,
		messageCodec:  messageCodec,
		tm:            tm,
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {

		requestItems := make([]*readWriteModel.S7VarRequestParameterItem, len(readRequest.GetFieldNames()))
		for i, fieldName := range readRequest.GetFieldNames() {
			field := readRequest.GetField(fieldName)
			address, err := encodeS7Address(field)
			if err != nil {
				result <- model.PlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding s7 address for field %s", fieldName),
				}
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
		log.Trace().Msg("Assemble paket")
		// TODO: why do we use a uint16 above and the cotp a uint8?
		tpktPacket := readWriteModel.NewTPKTPacket(
			readWriteModel.NewCOTPPacketData(true,
				uint8(tpduId),
				nil,
				s7MessageRequest,
			),
		)
		// Start a new request-transaction (Is ended in the response-handler)
		transaction := m.tm.StartTransaction()
		transaction.Submit(func() {

			// Send the  over the wire
			log.Trace().Msg("Send ")
			if err := m.messageCodec.SendRequest(
				tpktPacket,
				func(message interface{}) bool {
					tpktPacket := readWriteModel.CastTPKTPacket(message)
					if tpktPacket == nil {
						return false
					}
					cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
					if cotpPacketData == nil {
						return false
					}
					payload := cotpPacketData.Parent.Payload
					if payload == nil {
						return false
					}
					return payload.TpduReference == tpduId
				},
				func(message interface{}) error {
					// Convert the response into an
					log.Trace().Msg("convert response to ")
					tpktPacket := readWriteModel.CastTPKTPacket(message)
					cotpPacketData := readWriteModel.CastCOTPPacketData(tpktPacket.Payload)
					payload := cotpPacketData.Parent.Payload
					// Convert the s7 response into a PLC4X response
					log.Trace().Msg("convert response to PLC4X response")
					readResponse, err := m.ToPlc4xReadResponse(*payload, readRequest)

					if err != nil {
						result <- model.PlcReadRequestResult{
							Request: readRequest,
							Err:     errors.Wrap(err, "Error decoding response"),
						}
						return transaction.EndRequest()
					}
					result <- model.PlcReadRequestResult{
						Request:  readRequest,
						Response: readResponse,
					}
					return transaction.EndRequest()
				},
				func(err error) error {
					result <- model.PlcReadRequestResult{
						Request: readRequest,
						Err:     errors.Wrap(err, "got timeout while waiting for response"),
					}
					return transaction.EndRequest()
				},
				time.Second*1); err != nil {
				result <- model.PlcReadRequestResult{
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

func (m *Reader) ToPlc4xReadResponse(response readWriteModel.S7Message, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	var errorClass uint8
	var errorCode uint8
	switch response.Child.(type) {
	case *readWriteModel.S7MessageResponseData:
		messageResponseData := response.Child.(*readWriteModel.S7MessageResponseData)
		errorClass = messageResponseData.ErrorClass
		errorCode = messageResponseData.ErrorCode
	case *readWriteModel.S7MessageResponse:
		messageResponseData := response.Child.(*readWriteModel.S7MessageResponse)
		errorClass = messageResponseData.ErrorClass
		errorCode = messageResponseData.ErrorCode
	default:
		return nil, errors.Errorf("unsupported response type %T", response.Child)
	}
	responseCodes := map[string]model.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}

	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		// This is usually the case if PUT/GET wasn't enabled on the PLC
		if (errorClass == 129) && (errorCode == 4) {
			log.Warn().Msg("Got an error response from the PLC. This particular response code usually indicates " +
				"that PUT/GET is not enabled on the PLC.")
			for _, fieldName := range readRequest.GetFieldNames() {
				responseCodes[fieldName] = model.PlcResponseCode_ACCESS_DENIED
				plcValues[fieldName] = spiValues.NewPlcNULL()
			}
			log.Trace().Msg("Returning the response")
			return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
		} else {
			log.Warn().Msgf("Got an unknown error response from the PLC. Error Class: %d, Error Code %d. "+
				"We probably need to implement explicit handling for this, so please file a bug-report "+
				"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
				"containing a capture of the communication.",
				errorClass, errorCode)
			for _, fieldName := range readRequest.GetFieldNames() {
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
				plcValues[fieldName] = spiValues.NewPlcNULL()
			}
			return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
		}
	}

	// In all other cases all went well.
	payload := response.Payload.Child.(*readWriteModel.S7PayloadReadVarResponse)

	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(readRequest.GetFieldNames()) != len(payload.Items) {
		return nil, errors.New("The number of requested items doesn't match the number of returned items")
	}

	payloadItems := payload.Items
	for i, fieldName := range readRequest.GetFieldNames() {
		field := readRequest.GetField(fieldName).(S7PlcField)
		payloadItem := payloadItems[i]

		responseCode := decodeResponseCode(payloadItem.ReturnCode)
		// Decode the data according to the information from the request
		log.Trace().Msg("decode data")
		rb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(payloadItem.Data))
		responseCodes[fieldName] = responseCode
		if responseCode == model.PlcResponseCode_OK {
			plcValue, err := readWriteModel.DataItemParse(rb, field.GetDataType().DataProtocolId(), int32(field.GetNumElements()))
			if err != nil {
				return nil, errors.Wrap(err, "Error parsing data item")
			}
			plcValues[fieldName] = plcValue
		}
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

// Currently we only support the S7 Any type of addresses. This helper simply converts the S7Field from PLC4X into
// S7Address objects.
func encodeS7Address(field model.PlcField) (*readWriteModel.S7Address, error) {
	s7Field, ok := field.(S7PlcField)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", field)
	}
	transportSize := s7Field.GetDataType()
	numElements := s7Field.GetNumElements()
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
		if s7StringField, ok := field.(PlcStringField); ok {
			stringLength = s7StringField.stringLength
		}
		numElements = numElements * (stringLength + 2)
	} else if transportSize == readWriteModel.TransportSize_WSTRING {
		transportSize = readWriteModel.TransportSize_CHAR
		stringLength := uint16(254)
		if s7StringField, ok := field.(PlcStringField); ok {
			stringLength = s7StringField.stringLength
		}
		numElements = numElements * (stringLength + 2) * 2
	}
	return readWriteModel.NewS7AddressAny(
		transportSize,
		numElements,
		s7Field.GetBlockNumber(),
		s7Field.GetMemoryArea(),
		s7Field.GetByteOffset(),
		s7Field.GetBitOffset(),
	), nil
}

// Helper to convert the return codes returned from the S7 into one of our standard
func decodeResponseCode(dataTransportErrorCode readWriteModel.DataTransportErrorCode) model.PlcResponseCode {
	switch dataTransportErrorCode {
	case readWriteModel.DataTransportErrorCode_OK:
		return model.PlcResponseCode_OK
	case readWriteModel.DataTransportErrorCode_NOT_FOUND:
		return model.PlcResponseCode_NOT_FOUND
	case readWriteModel.DataTransportErrorCode_INVALID_ADDRESS:
		return model.PlcResponseCode_INVALID_ADDRESS
	case readWriteModel.DataTransportErrorCode_DATA_TYPE_NOT_SUPPORTED:
		return model.PlcResponseCode_INVALID_DATATYPE
	default:
		return model.PlcResponseCode_INTERNAL_ERROR
	}
}
