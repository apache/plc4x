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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/rs/zerolog"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"

	"github.com/pkg/errors"
)

type Writer struct {
	tpduGenerator *TpduGenerator
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager

	log zerolog.Logger
}

func NewWriter(tpduGenerator *TpduGenerator, messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, _options ...options.WithOption) Writer {
	return Writer{
		tpduGenerator: tpduGenerator,
		messageCodec:  messageCodec,
		tm:            tm,
		log:           options.ExtractCustomLogger(_options...),
	}
}

func (m Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	// TODO: handle context
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v", err))
			}
		}()
		parameterItems := make([]readWriteModel.S7VarRequestParameterItem, len(writeRequest.GetTagNames()))
		payloadItems := make([]readWriteModel.S7VarPayloadDataItem, len(writeRequest.GetTagNames()))
		for i, tagName := range writeRequest.GetTagNames() {
			tag := writeRequest.GetTag(tagName)
			plcValue := writeRequest.GetValue(tagName)
			s7Address, err := encodeS7Address(tag)
			if err != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error encoding s7 address for tag %s", tagName))
				return
			}
			parameterItems[i] = readWriteModel.NewS7VarRequestParameterItemAddress(s7Address)
			value, err := serializePlcValue(tag, plcValue)
			if err != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error encoding value for tag %s", tagName))
				return
			}
			payloadItems[i] = value
		}
		tpduId := m.tpduGenerator.getAndIncrement()

		// Create a new Request with correct tpuId (is not known before)
		s7MessageRequest := readWriteModel.NewS7MessageRequest(
			tpduId,
			readWriteModel.NewS7ParameterWriteVarRequest(parameterItems),
			readWriteModel.NewS7PayloadWriteVarRequest(payloadItems, nil),
		)

		// Assemble the finished paket
		m.log.Trace().Msg("Assemble paket")
		// TODO: why do we use a uint16 above and the cotp a uint8?
		tpktPacket := readWriteModel.NewTPKTPacket(
			readWriteModel.NewCOTPPacketData(
				true,
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
				readResponse, err := m.ToPlc4xWriteResponse(payload, writeRequest)

				if err != nil {
					result <- &spiModel.DefaultPlcWriteRequestResult{
						Request: writeRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
					return transaction.EndRequest()
				}
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: readResponse,
				}
				return transaction.EndRequest()
			}, func(err error) error {
				result <- &spiModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.New("got timeout while waiting for response"),
				}
				return transaction.EndRequest()
			}, time.Second*1); err != nil {
				result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrap(err, "error sending message"))
				_ = transaction.EndRequest()
			}
		})
	}()
	return result
}

func (m Writer) ToPlc4xWriteResponse(response readWriteModel.S7Message, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
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

	// If the result contains any form of non-null error code, handle this instead.
	if (errorClass != 0) || (errorCode != 0) {
		// This is usually the case if PUT/GET wasn't enabled on the PLC
		if (errorClass == 129) && (errorCode == 4) {
			m.log.Warn().Msg("Got an error response from the PLC. This particular response code usually indicates " +
				"that PUT/GET is not enabled on the PLC.")
			for _, tagName := range writeRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_ACCESS_DENIED
			}
			m.log.Trace().Msg("Returning the response")
			return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
		} else {
			m.log.Warn().Msgf("Got an unknown error response from the PLC. Error Class: %d, Error Code %d. "+
				"We probably need to implement explicit handling for this, so please file a bug-report "+
				"on https://issues.apache.org/jira/projects/PLC4X and ideally attach a WireShark dump "+
				"containing a capture of the communication.",
				errorClass, errorCode)
			for _, tagName := range writeRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
			return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
		}
	}

	// In all other cases all went well.
	payload := response.GetPayload().(readWriteModel.S7PayloadWriteVarResponse)

	// If the numbers of items don't match, we're in big trouble as the only
	// way to know how to interpret the responses is by aligning them with the
	// items from the request as this information is not returned by the PLC.
	if len(writeRequest.GetTagNames()) != len(payload.GetItems()) {
		return nil, errors.New("The number of requested items doesn't match the number of returned items")
	}

	payloadItems := payload.GetItems()
	for i, tagName := range writeRequest.GetTagNames() {
		payloadItem := payloadItems[i]

		responseCode := decodeResponseCode(payloadItem.GetReturnCode())
		// Decode the data according to the information from the request
		m.log.Trace().Msg("decode data")
		responseCodes[tagName] = responseCode
	}

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}

func serializePlcValue(tag apiModel.PlcTag, plcValue apiValues.PlcValue) (readWriteModel.S7VarPayloadDataItem, error) {
	s7Tag, ok := tag.(PlcTag)
	if !ok {
		return nil, errors.Errorf("Unsupported address type %t", tag)
	}
	transportSize := s7Tag.GetDataType().DataTransportSize()
	stringLength := uint16(254)
	if s7StringTag, ok := tag.(*PlcStringTag); ok {
		stringLength = s7StringTag.stringLength
	}
	data, err := readWriteModel.DataItemSerialize(plcValue, s7Tag.GetDataType().DataProtocolId(), int32(stringLength))
	if err != nil {
		return nil, errors.Wrapf(err, "Error serializing tag item of type: '%v'", s7Tag.GetDataType())
	}
	return readWriteModel.NewS7VarPayloadDataItem(
		readWriteModel.DataTransportErrorCode_OK,
		transportSize, data,
	), nil
}
