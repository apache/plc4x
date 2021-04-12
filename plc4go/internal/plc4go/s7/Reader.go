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

type Reader struct {
	messageCodec spi.MessageCodec
}

func NewReader(messageCodec spi.MessageCodec) *Reader {
	return &Reader{
		messageCodec: messageCodec,
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		if len(readRequest.GetFieldNames()) != 1 {
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.New("s7 only supports single-item requests"),
			}
			log.Debug().Msgf("s7 only supports single-item requests. Got %d fields", len(readRequest.GetFieldNames()))
			return
		}
		// If we are requesting only one field, use a
		fieldName := readRequest.GetFieldNames()[0]
		field := readRequest.GetField(fieldName)
		s7Field, err := CastTos7FieldFromPlcField(field)
		if err != nil {
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
		switch s7Field.FieldType {
		case FIELD:
			panic("Not implemented")
		case STRING_FIELD:
			panic("Not implemented")
		default:
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Errorf("unsupported field type %x", s7Field.FieldType),
			}
			log.Debug().Msgf("Unsupported field type %x", s7Field.FieldType)
			return
		}

		// Assemble the finished paket
		log.Trace().Msg("Assemble paket")
		request := readWriteModel.S7MessageRequest{}

		// Send the  over the wire
		log.Trace().Msg("Send ")
		if err = m.messageCodec.SendRequest(
			request,
			func(message interface{}) bool {
				readWriteModel.CastS7MessageResponse(message)
				// TODO: match to message
				return false
			},
			func(message interface{}) error {
				// Convert the response into an
				log.Trace().Msg("convert response to ")
				response := readWriteModel.CastS7MessageResponse(message)
				// Convert the s7 response into a PLC4X response
				log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(*response, readRequest)

				if err != nil {
					result <- model.PlcReadRequestResult{
						Request: readRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
					// TODO: should we return the error here?
					return nil
				}
				result <- model.PlcReadRequestResult{
					Request:  readRequest,
					Response: readResponse,
				}
				return nil
			},
			func(err error) error {
				result <- model.PlcReadRequestResult{
					Request: readRequest,
					Err:     errors.Wrap(err, "got timeout while waiting for response"),
				}
				return nil
			},
			time.Second*1); err != nil {
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error sending message"),
			}
		}
	}()
	return result
}

func (m *Reader) ToPlc4xReadResponse(response readWriteModel.S7MessageResponse, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	var data []uint8
	// TODO disect child type
	/*
		switch response.Pdu.Child.(type) {
		case *readWriteModel.s7PDUReadDiscreteInputsResponse:
			pdu := readWriteModel.Casts7PDUReadDiscreteInputsResponse(response.Pdu)
			data = utils.Int8ArrayToUint8Array(pdu.Value)
			// Pure Boolean ...
		case *readWriteModel.s7PDUReadCoilsResponse:
			pdu := readWriteModel.Casts7PDUReadCoilsResponse(&response.Pdu)
			data = utils.Int8ArrayToUint8Array(pdu.Value)
			// Pure Boolean ...
		case *readWriteModel.s7PDUReadInputRegistersResponse:
			pdu := readWriteModel.Casts7PDUReadInputRegistersResponse(response.Pdu)
			data = utils.Int8ArrayToUint8Array(pdu.Value)
			// DataIo ...
		case *readWriteModel.s7PDUReadHoldingRegistersResponse:
			pdu := readWriteModel.Casts7PDUReadHoldingRegistersResponse(response.Pdu)
			data = utils.Int8ArrayToUint8Array(pdu.Value)
		case *readWriteModel.s7PDUError:
			return nil, errors.Errorf("got an error from remote. Errorcode %x", response.Pdu.Child.(*readWriteModel.s7PDUError).ExceptionCode)
		default:
			return nil, errors.Errorf("unsupported response type %T", response.Pdu.Child)
		}
	*/

	// Get the field from the request
	log.Trace().Msg("get a field from request")
	fieldName := readRequest.GetFieldNames()[0]
	field, err := CastTos7FieldFromPlcField(readRequest.GetField(fieldName))
	if err != nil {
		return nil, errors.Wrap(err, "error casting to s7-field")
	}

	// Decode the data according to the information from the request
	log.Trace().Msg("decode data")
	rb := utils.NewReadBuffer(data)
	value, err := readWriteModel.DataItemParse(rb, field.Datatype.DataProtocolId(), int32(field.GetQuantity()))
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing data item")
	}
	responseCodes := map[string]model.PlcResponseCode{}
	plcValues := map[string]values.PlcValue{}
	plcValues[fieldName] = value
	responseCodes[fieldName] = model.PlcResponseCode_OK

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}
