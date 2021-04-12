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
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Writer struct {
	messageCodec spi.MessageCodec
}

func NewWriter(messageCodec spi.MessageCodec) Writer {
	return Writer{
		messageCodec: messageCodec,
	}
}

func (m Writer) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	result := make(chan model.PlcWriteRequestResult)
	go func() {
		// If we are requesting only one field, use a
		if len(writeRequest.GetFieldNames()) != 1 {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("s7 only supports single-item requests"),
			}
			return
		}
		fieldName := writeRequest.GetFieldNames()[0]

		// Get the s7 field instance from the request
		field := writeRequest.GetField(fieldName)
		s7Field, err := CastTos7FieldFromPlcField(field)
		if err != nil {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			return
		}
		stringLength := int32(0)
		if s7StringField, ok := field.(*PlcStringField); ok {
			stringLength = s7StringField.stringLength
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(fieldName)
		io := utils.NewWriteBuffer()
		if err := readWriteModel.DataItemSerialize(io, value, s7Field.Datatype.DataProtocolId(), stringLength); err != nil {
			result <- model.PlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error serializing value"),
			}
			return
		}
		utils.Uint8ArrayToInt8Array(io.GetBytes())

		// TODO disect child type
		/*
			var pdu *readWriteModel.s7PDU
			switch s7Field.FieldType {
			case Coil:
				pdu = readWriteModel.News7PDUWriteMultipleCoilsRequest(
					s7Field.Address,
					s7Field.Quantity,
					data)
			case HoldingRegister:
				pdu = readWriteModel.News7PDUWriteMultipleHoldingRegistersRequest(
					s7Field.Address,
					numWords,
					data)
			case ExtendedRegister:
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.New("s7 currently doesn't support extended register requests"),
				}
				return
			default:
				result <- model.PlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.New("unsupported field type"),
				}
				return
			}
		*/

		// Assemble the paket
		request := readWriteModel.S7MessageRequest{}

		// Send the  over the wire
		err = m.messageCodec.SendRequest(
			request,
			func(message interface{}) bool {
				readWriteModel.CastS7MessageResponse(message)
				// TODO: match to message
				return false
			},
			func(message interface{}) error {
				// Convert the response into an
				response := readWriteModel.CastS7MessageResponse(message)
				// Convert the s7 response into a PLC4X response
				readResponse, err := m.ToPlc4xWriteResponse(request, *response, writeRequest)

				if err != nil {
					result <- model.PlcWriteRequestResult{
						Request: writeRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
				} else {
					result <- model.PlcWriteRequestResult{
						Request:  writeRequest,
						Response: readResponse,
					}
				}
				return nil
			},
			func(err error) error {
				result <- model.PlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.New("got timeout while waiting for response"),
				}
				return nil
			},
			time.Second*1)
	}()
	return result
}

func (m Writer) ToPlc4xWriteResponse(request readWriteModel.S7MessageRequest, response readWriteModel.S7MessageResponse, writeRequest model.PlcWriteRequest) (model.PlcWriteResponse, error) {
	responseCodes := map[string]model.PlcResponseCode{}
	fieldName := writeRequest.GetFieldNames()[0]

	// we default to an error until its proven wrong
	responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
	// TODO disect child type
	/*
		switch response.Pdu.Child.(type) {
		case *readWriteModel.s7PDUWriteMultipleCoilsResponse:
			req := readWriteModel.Casts7PDUWriteMultipleCoilsRequest(request.Pdu)
			resp := readWriteModel.Casts7PDUWriteMultipleCoilsResponse(response.Pdu)
			if req.Quantity == resp.Quantity {
				responseCodes[fieldName] = model.PlcResponseCode_OK
			}
		case *readWriteModel.s7PDUWriteMultipleHoldingRegistersResponse:
			req := readWriteModel.Casts7PDUWriteMultipleHoldingRegistersRequest(request.Pdu)
			resp := readWriteModel.Casts7PDUWriteMultipleHoldingRegistersResponse(response.Pdu)
			if req.Quantity == resp.Quantity {
				responseCodes[fieldName] = model.PlcResponseCode_OK
			}
		case *readWriteModel.s7PDUError:
			resp := readWriteModel.Casts7PDUError(&response.Pdu)
			switch resp.ExceptionCode {
			case readWriteModel.s7ErrorCode_ILLEGAL_FUNCTION:
				responseCodes[fieldName] = model.PlcResponseCode_UNSUPPORTED
			case readWriteModel.s7ErrorCode_ILLEGAL_DATA_ADDRESS:
				responseCodes[fieldName] = model.PlcResponseCode_INVALID_ADDRESS
			case readWriteModel.s7ErrorCode_ILLEGAL_DATA_VALUE:
				responseCodes[fieldName] = model.PlcResponseCode_INVALID_DATA
			case readWriteModel.s7ErrorCode_SLAVE_DEVICE_FAILURE:
				responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
			case readWriteModel.s7ErrorCode_ACKNOWLEDGE:
				responseCodes[fieldName] = model.PlcResponseCode_OK
			case readWriteModel.s7ErrorCode_SLAVE_DEVICE_BUSY:
				responseCodes[fieldName] = model.PlcResponseCode_REMOTE_BUSY
			case readWriteModel.s7ErrorCode_NEGATIVE_ACKNOWLEDGE:
				responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
			case readWriteModel.s7ErrorCode_MEMORY_PARITY_ERROR:
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			case readWriteModel.s7ErrorCode_GATEWAY_PATH_UNAVAILABLE:
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			case readWriteModel.s7ErrorCode_GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND:
				responseCodes[fieldName] = model.PlcResponseCode_REMOTE_ERROR
			default:
				log.Debug().Msgf("Unmapped exception code %x", resp.ExceptionCode)
			}
		default:
			return nil, errors.Errorf("unsupported response type %T", response.Pdu.Child)
		}
	*/

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
