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
package ads

import (
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"sync/atomic"
	"time"
)

type Reader struct {
	transactionIdentifier uint32
	targetAmsNetId        readWriteModel.AmsNetId
	targetAmsPort         uint16
	sourceAmsNetId        readWriteModel.AmsNetId
	sourceAmsPort         uint16
	messageCodec          spi.MessageCodec
}

func NewReader(messageCodec spi.MessageCodec, targetAmsNetId readWriteModel.AmsNetId, targetAmsPort uint16, sourceAmsNetId readWriteModel.AmsNetId, sourceAmsPort uint16) *Reader {
	return &Reader{
		transactionIdentifier: 0,
		targetAmsNetId:        targetAmsNetId,
		targetAmsPort:         targetAmsPort,
		sourceAmsNetId:        sourceAmsNetId,
		sourceAmsPort:         sourceAmsPort,
		messageCodec:          messageCodec,
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
				Err:      errors.New("ads only supports single-item requests"),
			}
			log.Debug().Msgf("ads only supports single-item requests. Got %d fields", len(readRequest.GetFieldNames()))
			return
		}
		// If we are requesting only one field, use a
		fieldName := readRequest.GetFieldNames()[0]
		field := readRequest.GetField(fieldName)
		adsField, err := CastToAdsFieldFromPlcField(field)
		if err != nil {
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			log.Debug().Msgf("Invalid field item type %T", field)
			return
		}
		userdata := readWriteModel.AmsPacket{
			TargetAmsNetId: &m.targetAmsNetId,
			TargetAmsPort:  m.targetAmsPort,
			SourceAmsNetId: &m.sourceAmsNetId,
			SourceAmsPort:  m.sourceAmsPort,
			CommandId:      readWriteModel.CommandId_ADS_READ,
			State:          readWriteModel.NewState(false, false, false, false, false, true, false, false, false),
			ErrorCode:      0,
			InvokeId:       0,
			Data:           nil,
		}
		switch adsField.FieldType {
		case StringField:
			// TODO: what is our read length?
			userdata.Data = readWriteModel.NewAdsReadRequest(adsField.IndexGroup, adsField.IndexOffset, 1)
		case Field:
			// TODO: what is our read length?
			userdata.Data = readWriteModel.NewAdsReadRequest(adsField.IndexGroup, adsField.IndexOffset, 1)
		case SymbolicStringField:
			panic("implement me")
		case SymbolicField:
			panic("implement me")
		default:
			result <- model.PlcReadRequestResult{
				Request:  readRequest,
				Response: nil,
				Err:      errors.Errorf("unsupported field type %x", adsField.FieldType),
			}
			log.Debug().Msgf("Unsupported field type %x", adsField.FieldType)
			return
		}

		// Calculate a new transaction identifier
		transactionIdentifier := atomic.AddUint32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 1
			atomic.StoreUint32(&m.transactionIdentifier, 1)
		}
		log.Debug().Msgf("Calculated transaction identifier %x", transactionIdentifier)
		userdata.InvokeId = transactionIdentifier

		// Assemble the finished tcp paket
		log.Trace().Msg("Assemble tcp paket")
		amsTcpPaket := readWriteModel.AmsTCPPacket{
			Userdata: &userdata,
		}

		// Send the TCP Paket over the wire
		log.Trace().Msg("Send TCP Paket")
		if err = m.messageCodec.SendRequest(
			amsTcpPaket,
			func(message interface{}) bool {
				paket := readWriteModel.CastAmsTCPPacket(message)
				return paket.Userdata.InvokeId == transactionIdentifier
			},
			func(message interface{}) error {
				// Convert the response into an amsTcpPaket
				log.Trace().Msg("convert response to amsTcpPaket")
				amsTcpPaket := readWriteModel.CastAmsTCPPacket(message)
				// Convert the ads response into a PLC4X response
				log.Trace().Msg("convert response to PLC4X response")
				readResponse, err := m.ToPlc4xReadResponse(*amsTcpPaket, readRequest)

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

func (m *Reader) ToPlc4xReadResponse(amsTcpPaket readWriteModel.AmsTCPPacket, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	var data []uint8
	switch amsTcpPaket.Userdata.Data.Child.(type) {
	case *readWriteModel.AdsReadResponse:
		readResponse := readWriteModel.CastAdsReadResponse(amsTcpPaket.Userdata.Data)
		data = utils.Int8ArrayToUint8Array(readResponse.Data)
		// Pure Boolean ...
	default:
		return nil, errors.Errorf("unsupported response type %T", amsTcpPaket.Userdata.Data.Child)
	}

	// Get the field from the request
	log.Trace().Msg("get a field from request")
	fieldName := readRequest.GetFieldNames()[0]
	field, err := CastToAdsFieldFromPlcField(readRequest.GetField(fieldName))
	if err != nil {
		return nil, errors.Wrap(err, "error casting to ads-field")
	}

	// Decode the data according to the information from the request
	log.Trace().Msg("decode data")
	rb := utils.NewReadBuffer(data)
	value, err := readWriteModel.DataItemParse(rb, field.Datatype.DataFormatName(), field.StringLength)
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
