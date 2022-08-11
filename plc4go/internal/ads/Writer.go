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

package ads

import (
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"sync/atomic"
	"time"
)

type Writer struct {
	transactionIdentifier uint32
	targetAmsNetId        readWriteModel.AmsNetId
	targetAmsPort         uint16
	sourceAmsNetId        readWriteModel.AmsNetId
	sourceAmsPort         uint16
	messageCodec          spi.MessageCodec
	reader                *Reader
}

func NewWriter(messageCodec spi.MessageCodec, targetAmsNetId readWriteModel.AmsNetId, targetAmsPort uint16, sourceAmsNetId readWriteModel.AmsNetId, sourceAmsPort uint16, reader *Reader) *Writer {
	return &Writer{
		transactionIdentifier: 0,
		targetAmsNetId:        targetAmsNetId,
		targetAmsPort:         targetAmsPort,
		sourceAmsNetId:        sourceAmsNetId,
		sourceAmsPort:         sourceAmsPort,
		messageCodec:          messageCodec,
		reader:                reader,
	}
}

func (m *Writer) Write(writeRequest model.PlcWriteRequest) <-chan model.PlcWriteRequestResult {
	result := make(chan model.PlcWriteRequestResult)
	go func() {
		// If we are requesting only one field, use a
		if len(writeRequest.GetFieldNames()) != 1 {
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("ads only supports single-item requests"),
			}
			return
		}
		fieldName := writeRequest.GetFieldNames()[0]

		// Get the ads field instance from the request
		field := writeRequest.GetField(fieldName)
		if needsResolving(field) {
			adsField, err := castToSymbolicPlcFieldFromPlcField(field)
			if err != nil {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid field item type"),
				}
				log.Debug().Msgf("Invalid field item type %T", field)
				return
			}
			field, err = m.reader.resolveField(adsField)
			if err != nil {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "invalid field item type"),
				}
				log.Debug().Msgf("Invalid field item type %T", field)
				return
			}
		}
		adsField, err := castToDirectAdsFieldFromPlcField(field)
		if err != nil {
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "invalid field item type"),
			}
			return
		}

		// Get the value from the request and serialize it to a byte array
		value := writeRequest.GetValue(fieldName)
		io := utils.NewLittleEndianWriteBufferByteBased()
		if err := readWriteModel.DataItemSerialize(io, value, adsField.Datatype.DataFormatName(), adsField.StringLength); err != nil {
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.Wrap(err, "error serializing value"),
			}
			return
		}
		data := io.GetBytes()

		userdata := readWriteModel.NewAmsPacket(
			m.targetAmsNetId,
			m.targetAmsPort,
			m.sourceAmsNetId,
			m.sourceAmsPort,
			readWriteModel.CommandId_ADS_READ,
			readWriteModel.NewState(false, false, false, false, false, true, false, false, false),
			0,
			0,
			nil,
		)
		switch adsField.FieldType {
		case DirectAdsStringField:
			userdata.Data = readWriteModel.NewAdsWriteRequest(adsField.IndexGroup, adsField.IndexOffset, data)
			panic("implement me")
		case DirectAdsField:
			panic("implement me")
		case SymbolicAdsStringField, SymbolicAdsField:
			panic("we should never reach this point as symbols are resolved before")
		default:
			result <- &plc4goModel.DefaultPlcWriteRequestResult{
				Request:  writeRequest,
				Response: nil,
				Err:      errors.New("unsupported field type"),
			}
			return
		}

		// Calculate a new unit identifier
		// TODO: this is not threadsafe as the whole operation is not atomic
		transactionIdentifier := atomic.AddUint32(&m.transactionIdentifier, 1)
		if transactionIdentifier > math.MaxUint8 {
			transactionIdentifier = 0
			atomic.StoreUint32(&m.transactionIdentifier, 0)
		}
		userdata.InvokeId = transactionIdentifier

		// Assemble the finished amsTcpPaket
		log.Trace().Msg("Assemble amsTcpPaket")
		amsTcpPaket := readWriteModel.NewAmsTCPPacket(userdata)

		// Send the TCP Paket over the wire
		err = m.messageCodec.SendRequest(
			amsTcpPaket,
			func(message spi.Message) bool {
				paket := readWriteModel.CastAmsTCPPacket(message)
				return paket.GetUserdata().GetInvokeId() == transactionIdentifier
			},
			func(message spi.Message) error {
				// Convert the response into an responseAmsTcpPaket
				responseAmsTcpPaket := readWriteModel.CastAmsTCPPacket(message)
				// Convert the ads response into a PLC4X response
				readResponse, err := m.ToPlc4xWriteResponse(amsTcpPaket, responseAmsTcpPaket, writeRequest)

				if err != nil {
					result <- &plc4goModel.DefaultPlcWriteRequestResult{
						Request: writeRequest,
						Err:     errors.Wrap(err, "Error decoding response"),
					}
				} else {
					result <- &plc4goModel.DefaultPlcWriteRequestResult{
						Request:  writeRequest,
						Response: readResponse,
					}
				}
				return nil
			},
			func(err error) error {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.New("got timeout while waiting for response"),
				}
				return nil
			},
			time.Second*1)
	}()
	return result
}

func (m *Writer) ToPlc4xWriteResponse(requestTcpPaket readWriteModel.AmsTCPPacket, responseTcpPaket readWriteModel.AmsTCPPacket, writeRequest model.PlcWriteRequest) (model.PlcWriteResponse, error) {
	responseCodes := map[string]model.PlcResponseCode{}
	fieldName := writeRequest.GetFieldNames()[0]

	// we default to an error until its proven wrong
	responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
	switch responseTcpPaket.GetUserdata().GetData().(type) {
	case readWriteModel.AdsWriteResponse:
		resp := readWriteModel.CastAdsWriteResponse(responseTcpPaket.GetUserdata().GetData())
		responseCodes[fieldName] = model.PlcResponseCode(resp.GetResult())
	default:
		return nil, errors.Errorf("unsupported response type %T", responseTcpPaket.GetUserdata().GetData())
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
