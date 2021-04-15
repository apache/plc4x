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
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Writer struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            *spi.RequestTransactionManager
}

func NewWriter(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager) Writer {
	return Writer{
		tpduGenerator: tpduGenerator,
		messageCodec:  messageCodec,
		tm:            tm,
	}
}

func (m Writer) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	result := make(chan model.PlcWriteRequestResult)
	go func() {
		parameterItems := make([]*readWriteModel.S7VarRequestParameterItem, len(writeRequest.GetFieldNames()))
		payloadItems := make([]*readWriteModel.S7VarPayloadDataItem, len(writeRequest.GetFieldNames()))
		for i, fieldName := range writeRequest.GetFieldNames() {
			field := writeRequest.GetField(fieldName)
			plcValue := writeRequest.GetValue(fieldName)
			s7Address, err := encodeS7Address(field)
			if err != nil {
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding s7 address for field %s", fieldName),
				}
				return
			}
			parameterItems[i] = readWriteModel.NewS7VarRequestParameterItemAddress(s7Address)
			value, err := serializePlcValue(field, plcValue)
			if err != nil {
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding value for field %s", fieldName),
				}
				return
			}
			payloadItems[i] = value
		}
		tpduId := m.tpduGenerator.getAndIncrement()

		// Create a new Request with correct tpuId (is not known before)
		s7MessageRequest := readWriteModel.NewS7MessageRequest(
			tpduId,
			readWriteModel.NewS7ParameterWriteVarRequest(parameterItems),
			readWriteModel.NewS7PayloadWriteVarRequest(payloadItems),
		)

		// Assemble the finished paket
		log.Trace().Msg("Assemble paket")
		// TODO: why do we use a uint16 above and the cotp a uint8?
		tpktPacket := readWriteModel.NewTPKTPacket(
			readWriteModel.NewCOTPPacketData(
				true,
				uint8(tpduId),
				nil,
				s7MessageRequest,
			),
		)

		// Start a new request-transaction (Is ended in the response-handler)
		transaction := m.tm.StartTransaction()
		transaction.Submit(func() {
			// Send the  over the wire
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
					readResponse, err := m.ToPlc4xWriteResponse(*payload, writeRequest)

					if err != nil {
						result <- model.PlcWriteRequestResult{
							Request: writeRequest,
							Err:     errors.Wrap(err, "Error decoding response"),
						}
						return transaction.EndRequest()
					}
					result <- model.PlcWriteRequestResult{
						Request:  writeRequest,
						Response: readResponse,
					}
					return transaction.EndRequest()
				},
				func(err error) error {
					result <- model.PlcWriteRequestResult{
						Request: writeRequest,
						Err:     errors.New("got timeout while waiting for response"),
					}
					return transaction.EndRequest()
				},
				time.Second*1); err != nil {
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "error sending message"),
				}
				_ = transaction.EndRequest()
			}
		})
	}()
	return result
}

func (m Writer) ToPlc4xWriteResponse(response readWriteModel.S7Message, writeRequest model.PlcWriteRequest) (model.PlcWriteResponse, error) {
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

	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		// This is usually the case if PUT/GET wasn't enabled on the PLC
		if (errorClass == 129) && (errorCode == 4) {
			log.Warn().Msg("Got an error response from the PLC. This particular response code usually indicates " +
				"that PUT/GET is not enabled on the PLC.")
			for _, fieldName := range writeRequest.GetFieldNames() {
				responseCodes[fieldName] = model.PlcResponseCode_ACCESS_DENIED
			}
			log.Trace().Msg("Returning the response")
			return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
		} else {
			log.Warn().Msgf("Got an unknown error response from the PLC. Error Class: %d, Error Code %d. "+
				"We probably need to implement explicit handling for this, so please file a bug-report "+
				"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
				"containing a capture of the communication.",
				errorClass, errorCode)
			for _, fieldName := range writeRequest.GetFieldNames() {
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			}
			return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
		}
	}

	// In all other cases all went well.
	payload := response.Payload.Child.(*readWriteModel.S7PayloadWriteVarResponse)

	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(writeRequest.GetFieldNames()) != len(payload.Items) {
		return nil, errors.New("The number of requested items doesn't match the number of returned items")
	}

	payloadItems := payload.Items
	for i, fieldName := range writeRequest.GetFieldNames() {
		payloadItem := payloadItems[i]

		responseCode := decodeResponseCode(payloadItem.ReturnCode)
		// Decode the data according to the information from the request
		log.Trace().Msg("decode data")
		responseCodes[fieldName] = responseCode
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}

func serializePlcValue(field model.PlcField, plcValue values.PlcValue) (*readWriteModel.S7VarPayloadDataItem, error) {
	s7Field, ok := field.(S7PlcField)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", field)
	}
	transportSize := s7Field.GetDataType().DataTransportSize()
	stringLength := uint16(254)
	if s7StringField, ok := field.(*PlcStringField); ok {
		stringLength = s7StringField.stringLength
	}
	io := utils.NewWriteBuffer()
	err := readWriteModel.DataItemSerialize(io, plcValue, s7Field.GetDataType().DataProtocolId(), int32(stringLength))
	if err != nil {
		return nil, errors.Wrapf(err, "Error serializing field item of type: '%v'", s7Field.GetDataType())
	}
	data := io.GetBytes()
	return readWriteModel.NewS7VarPayloadDataItem(
		readWriteModel.DataTransportErrorCode_OK,
		transportSize, utils.Uint8ArrayToInt8Array(data),
	), nil
}
