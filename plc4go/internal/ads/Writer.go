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
	"context"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(m.GetPlcTagHandler(), m.GetPlcValueHandler(), m)
}

func (m *Connection) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	/*	// TODO: handle context
		result := make(chan model.PlcWriteRequestResult)
		go func() {
			// If we are requesting only one field, use a
			if len(writeRequest.GetTagNames()) != 1 {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.New("ads only supports single-item requests"),
				}
				return
			}
			fieldName := writeRequest.GetTagNames()[0]

			// Get the ads field instance from the request
			field := writeRequest.GetTag(fieldName)
			if needsResolving(field) {
				adsField, err := castToSymbolicPlcTagFromPlcTag(field)
				if err != nil {
					result <- &plc4goModel.DefaultPlcWriteRequestResult{
						Request:  writeRequest,
						Response: nil,
						Err:      errors.Wrap(err, "invalid field item type"),
					}
					log.Debug().Msgf("Invalid field item type %T", field)
					return
				}
				field, err = m.reader.resolveTag(ctx, adsField)
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
			adsField, err := castToDirectAdsTagFromPlcTag(field)
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
			io := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
			if err := readWriteModel.DataItemSerializeWithWriteBuffer(io, value, adsField.Datatype.PlcValueType(), adsField.StringLength); err != nil {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request:  writeRequest,
					Response: nil,
					Err:      errors.Wrap(err, "error serializing value"),
				}
				return
			}
			/data := io.GetBytes()

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
			)/
			switch adsField.TagType {
			case DirectAdsStringField:
				//userdata.Data = readWriteModel.NewAdsWriteRequest(adsField.IndexGroup, adsField.IndexOffset, data)
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
			/userdata.InvokeId = m.getInvokeId()

			// Assemble the finished amsTcpPaket
			log.Trace().Msg("Assemble amsTcpPaket")
			amsTcpPaket := readWriteModel.NewAmsTCPPacket(userdata)

			// Send the TCP Paket over the wire
			err = m.messageCodec.SendRequest(ctx, amsTcpPaket, func(message spi.Message) bool {
				paket := readWriteModel.CastAmsTCPPacket(message)
				return paket.GetUserdata().GetInvokeId() == transactionIdentifier
			}, func(message spi.Message) error {
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
			}, func(err error) error {
				result <- &plc4goModel.DefaultPlcWriteRequestResult{
					Request: writeRequest,
					Err:     errors.New("got timeout while waiting for response"),
				}
				return nil
			}, time.Second*1)/
		}()
		return result
	*/
	return nil
}

func (m *Connection) ToPlc4xWriteResponse(requestTcpPaket driverModel.AmsTCPPacket, responseTcpPaket driverModel.AmsTCPPacket, writeRequest apiModel.PlcWriteRequest) (apiModel.PlcWriteResponse, error) {
	responseCodes := map[string]apiModel.PlcResponseCode{}
	tagName := writeRequest.GetTagNames()[0]

	// we default to an error until its proven wrong
	responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
	switch writeResponse := responseTcpPaket.GetUserdata().(type) {
	case driverModel.AdsWriteResponseExactly:
		responseCodes[tagName] = apiModel.PlcResponseCode(writeResponse.GetResult())
	default:
		return nil, errors.Errorf("unsupported response type %T", responseTcpPaket.GetUserdata())
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return internalModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil
}
