/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package eip

import (
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type Reader struct {
	messageCodec  spi.MessageCodec
	tm            *spi.RequestTransactionManager
	configuration Configuration
	sessionHandle *uint32
}

func NewReader(messageCodec spi.MessageCodec, tm *spi.RequestTransactionManager, configuration Configuration, sessionHandle *uint32) *Reader {
	return &Reader{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,
		sessionHandle: sessionHandle,
	}
}

func (m *Reader) Read(readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {

		requestItems := make([]*readWriteModel.CipService, len(readRequest.GetFieldNames()))
		for i, fieldName := range readRequest.GetFieldNames() {
			plcField := readRequest.GetField(fieldName).(EIPPlcField)
			tag := plcField.GetTag()
			elements := uint16(1)
			if plcField.GetElementNb() > 1 {
				elements = plcField.GetElementNb()
			}
			ansi, err := toAnsi(tag)
			if err != nil {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding eip ansi for field %s", fieldName),
				}
				return
			}
			request := readWriteModel.NewCipReadRequest(getRequestSize(tag), ansi, elements)
			requestItems[i] = request
		}
		if len(requestItems) > 1 {
			nb := uint16(len(requestItems))
			offsets := make([]uint16, nb)
			offset := 2 + nb*2
			for i := uint16(0); i < nb; i++ {
				offsets[i] = offset
				offset += requestItems[i].LengthInBytes()
			}

			serviceArr := make([]*readWriteModel.CipService, nb)
			for i := uint16(0); i < nb; i++ {
				serviceArr[i] = requestItems[i]
			}

			data := readWriteModel.NewServices(nb, offsets, serviceArr)
			//Encapsulate the data
			pkt := readWriteModel.NewCipRRData(
				readWriteModel.NewCipExchange(
					readWriteModel.NewCipUnconnectedRequest(
						readWriteModel.NewMultipleServiceRequest(data),
						m.configuration.backplane,
						m.configuration.slot,
					),
				),
				*m.sessionHandle,
				0,
				make([]byte, 8),
				0,
			)

			// Start a new request-transaction (Is ended in the response-handler)
			transaction := m.tm.StartTransaction()
			transaction.Submit(func() {
				// Send the  over the wire
				log.Trace().Msg("Send ")
				if err := m.messageCodec.SendRequest(
					pkt,
					func(message interface{}) bool {
						eipPacket := readWriteModel.CastEipPacket(message)
						if eipPacket == nil {
							return false
						}
						cipRRData := readWriteModel.CastCipRRData(eipPacket.Child)
						if cipRRData == nil {
							return false
						}
						if eipPacket.SessionHandle != *m.sessionHandle {
							return false
						}
						multipleServiceResponse := readWriteModel.CastMultipleServiceResponse(cipRRData.Exchange.Service)
						if multipleServiceResponse == nil {
							return false
						}
						if multipleServiceResponse.ServiceNb != nb {
							return false
						}
						return true
					},
					func(message interface{}) error {
						// Convert the response into an
						log.Trace().Msg("convert response to ")
						eipPacket := readWriteModel.CastEipPacket(message)
						cipRRData := readWriteModel.CastCipRRData(eipPacket.Child)
						multipleServiceResponse := readWriteModel.CastMultipleServiceResponse(cipRRData.Exchange.Service)
						// Convert the eip response into a PLC4X response
						log.Trace().Msg("convert response to PLC4X response")
						readResponse, err := m.ToPlc4xReadResponse(multipleServiceResponse.Parent, readRequest)

						if err != nil {
							result <- &plc4goModel.DefaultPlcReadRequestResult{
								Request: readRequest,
								Err:     errors.Wrap(err, "Error decoding response"),
							}
							return transaction.EndRequest()
						}
						result <- &plc4goModel.DefaultPlcReadRequestResult{
							Request:  readRequest,
							Response: readResponse,
						}
						return transaction.EndRequest()
					},
					func(err error) error {
						result <- &plc4goModel.DefaultPlcReadRequestResult{
							Request: readRequest,
							Err:     errors.Wrap(err, "got timeout while waiting for response"),
						}
						return transaction.EndRequest()
					},
					time.Second*1); err != nil {
					result <- &plc4goModel.DefaultPlcReadRequestResult{
						Request:  readRequest,
						Response: nil,
						Err:      errors.Wrap(err, "error sending message"),
					}
					_ = transaction.EndRequest()
				}
			})
		} else if len(requestItems) == 1 {
			//Encapsulate the data
			pkt := readWriteModel.NewCipRRData(
				readWriteModel.NewCipExchange(
					readWriteModel.NewCipUnconnectedRequest(
						requestItems[0],
						m.configuration.backplane,
						m.configuration.slot,
					),
				),
				*m.sessionHandle,
				0,
				make([]byte, 8),
				0,
			)

			// Start a new request-transaction (Is ended in the response-handler)
			transaction := m.tm.StartTransaction()
			transaction.Submit(func() {
				// Send the  over the wire
				log.Trace().Msg("Send ")
				if err := m.messageCodec.SendRequest(
					pkt,
					func(message interface{}) bool {
						eipPacket := readWriteModel.CastEipPacket(message)
						if eipPacket == nil {
							return false
						}
						cipRRData := readWriteModel.CastCipRRData(eipPacket.Child)
						if cipRRData == nil {
							return false
						}
						if eipPacket.SessionHandle != *m.sessionHandle {
							return false
						}
						cipReadResponse := readWriteModel.CastCipReadResponse(cipRRData.Exchange.Service)
						if cipReadResponse == nil {
							return false
						}
						return true
					},
					func(message interface{}) error {
						// Convert the response into an
						log.Trace().Msg("convert response to ")
						eipPacket := readWriteModel.CastEipPacket(message)
						cipRRData := readWriteModel.CastCipRRData(eipPacket.Child)
						cipReadResponse := readWriteModel.CastCipReadResponse(cipRRData.Exchange.Service)
						// Convert the eip response into a PLC4X response
						log.Trace().Msg("convert response to PLC4X response")
						readResponse, err := m.ToPlc4xReadResponse(cipReadResponse.Parent, readRequest)

						if err != nil {
							result <- &plc4goModel.DefaultPlcReadRequestResult{
								Request: readRequest,
								Err:     errors.Wrap(err, "Error decoding response"),
							}
							return transaction.EndRequest()
						}
						result <- &plc4goModel.DefaultPlcReadRequestResult{
							Request:  readRequest,
							Response: readResponse,
						}
						return transaction.EndRequest()
					},
					func(err error) error {
						result <- &plc4goModel.DefaultPlcReadRequestResult{
							Request: readRequest,
							Err:     errors.Wrap(err, "got timeout while waiting for response"),
						}
						return transaction.EndRequest()
					},
					time.Second*1); err != nil {
					result <- &plc4goModel.DefaultPlcReadRequestResult{
						Request:  readRequest,
						Response: nil,
						Err:      errors.Wrap(err, "error sending message"),
					}
					_ = transaction.EndRequest()
				}
			})
		}
	}()
	return result
}

func getRequestSize(tag string) int8 {
	//We need the size of the request in words (0x91, tagLength, ... tag + possible pad)
	// Taking half to get word size
	isArray := false
	isStruct := false
	tagIsolated := tag

	if strings.Contains(tag, "[") {
		isArray = true
		tagIsolated = tag[0:strings.Index(tag, "[")]
	}

	if strings.Contains(tag, ".") {
		isStruct = true
		tagIsolated = strings.Replace(tagIsolated, ".", "", -1)
	}
	dataLength := (len(tagIsolated) + 2) + (len(tagIsolated) % 2)
	if isArray {
		dataLength += 2
	}
	if isStruct {
		dataLength += 2
	}
	requestPathSize := (int8)(dataLength / 2)
	return requestPathSize
}

func toAnsi(tag string) ([]byte, error) {
	arrayIndex := byte(0)
	isArray := false
	isStruct := false
	tagFinal := tag
	if strings.Contains(tag, "[") {
		isArray = true
		index := tag[strings.Index(tag, "[")+1 : strings.Index(tag, "]")]
		parsedArrayIndex, err := strconv.ParseUint(index, 10, 8)
		if err != nil {
			return nil, err
		}
		arrayIndex = byte(parsedArrayIndex)
		tagFinal = tag[0:strings.Index(tag, "[")]
	}
	if strings.Contains(tag, ".") {
		tagFinal = tag[0:strings.Index(tag, ".")]
		isStruct = true
	}
	isPadded := len(tagFinal)%2 != 0
	dataSegLength := 2 + len(tagFinal)
	if isPadded {
		dataSegLength += 1
	}
	if isArray {
		dataSegLength += 2
	}

	if isStruct {
		for _, subStr := range strings.Split(tag[strings.Index(tag, ".")+1:], ".") {
			dataSegLength += 2 + len(subStr) + len(subStr)%2
		}
	}

	buffer := utils.NewLittleEndianWriteBufferByteBased()

	err := buffer.WriteByte("", 0x91)
	if err != nil {
		return nil, err
	}
	err = buffer.WriteByte("", byte(len(tagFinal)))
	if err != nil {
		return nil, err
	}

	quoteToASCII := strconv.QuoteToASCII(tagFinal)
	err = buffer.WriteByteArray("", []byte(quoteToASCII)[1:len(quoteToASCII)-1])
	if err != nil {
		return nil, err
	}

	if isPadded {
		err = buffer.WriteByte("", 0x00)
		if err != nil {
			return nil, err
		}
	}

	if isArray {
		err = buffer.WriteByte("", 0x28)
		if err != nil {
			return nil, err
		}
		err = buffer.WriteByte("", arrayIndex)
		if err != nil {
			return nil, err
		}
	}
	if isStruct {
		ansi, err := toAnsi(tag[strings.Index(tag, ".")+1:])
		if err != nil {
			return nil, err
		}
		err = buffer.WriteByteArray("", ansi)
		if err != nil {
			return nil, err
		}
	}
	return buffer.GetBytes(), nil
}

func (m *Reader) ToPlc4xReadResponse(response *readWriteModel.CipService, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	plcValues := map[string]values.PlcValue{}
	responseCodes := map[string]model.PlcResponseCode{}
	switch response.Child.(type) {
	case *readWriteModel.CipReadResponse: // only 1 field
		cipReadResponse := response.Child.(*readWriteModel.CipReadResponse)
		fieldName := readRequest.GetFieldNames()[0]
		field := readRequest.GetField(fieldName).(EIPPlcField)
		code := decodeResponseCode(cipReadResponse.Status)
		var plcValue values.PlcValue
		_type := cipReadResponse.DataType
		data := utils.NewLittleEndianReadBufferByteBased(cipReadResponse.Data)
		if code == model.PlcResponseCode_OK {
			var err error
			plcValue, err = parsePlcValue(field, data, _type)
			if err != nil {
				return nil, err
			}
		}
		plcValues[fieldName] = plcValue
		responseCodes[fieldName] = code
	case *readWriteModel.MultipleServiceResponse: //Multiple response
		multipleServiceResponse := response.Child.(*readWriteModel.MultipleServiceResponse)
		nb := multipleServiceResponse.ServiceNb
		arr := make([]*readWriteModel.CipService, nb)
		read := utils.NewLittleEndianReadBufferByteBased(multipleServiceResponse.ServicesData)
		total := read.GetTotalBytes()
		for i := uint16(0); i < nb; i++ {
			length := uint16(0)
			offset := multipleServiceResponse.Offsets[i] - multipleServiceResponse.Offsets[0] //Substract first offset as we only have the service in the buffer (not servicesNb and offsets)
			if i == nb-1 {
				length = uint16(total) - offset //Get the rest if last
			} else {
				length = multipleServiceResponse.Offsets[i+1] - offset - multipleServiceResponse.Offsets[0] //Calculate length with offsets (substracting first offset)
			}
			serviceBuf := utils.NewLittleEndianReadBufferByteBased(read.GetBytes()[offset : offset+length])
			var err error
			arr[i], err = readWriteModel.CipServiceParse(serviceBuf, length)
			if err != nil {
				return nil, err
			}
		}
		services := readWriteModel.NewServices(nb, multipleServiceResponse.Offsets, arr)
		for i, fieldName := range readRequest.GetFieldNames() {
			field := readRequest.GetField(fieldName).(EIPPlcField)
			if cipReadResponse, ok := services.Services[i].Child.(*readWriteModel.CipReadResponse); ok {
				code := decodeResponseCode(cipReadResponse.Status)
				_type := cipReadResponse.DataType
				data := utils.NewLittleEndianReadBufferByteBased(cipReadResponse.Data)
				var plcValue values.PlcValue
				if code == model.PlcResponseCode_OK {
					var err error
					plcValue, err = parsePlcValue(field, data, _type)
					if err != nil {
						return nil, err
					}
				}

				plcValues[fieldName] = plcValue
				responseCodes[fieldName] = code
			} else {
				responseCodes[fieldName] = model.PlcResponseCode_INTERNAL_ERROR
			}
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", response.Child)
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

func parsePlcValue(field EIPPlcField, data utils.ReadBufferByteBased, _type readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
	nb := field.GetElementNb()
	if nb > 1 {
		list := make([]values.PlcValue, 0)
		for i := uint16(0); i < nb; i++ {
			switch _type {
			case readWriteModel.CIPDataTypeCode_DINT:
				readInt32, err := data.ReadInt32("", _type.Size()*8)
				if err != nil {
					return nil, err
				}
				list = append(list, spiValues.NewPlcDINT(readInt32))
			case readWriteModel.CIPDataTypeCode_INT:
				readInt16, err := data.ReadInt16("", _type.Size()*8)
				if err != nil {
					return nil, err
				}
				list = append(list, spiValues.NewPlcINT(readInt16))
			case readWriteModel.CIPDataTypeCode_SINT:
				readInt8, err := data.ReadInt8("", _type.Size()*8)
				if err != nil {
					return nil, err
				}
				list = append(list, spiValues.NewPlcSINT(readInt8))
			case readWriteModel.CIPDataTypeCode_REAL:
				if _type.Size()*8 != 64 {
					panic("Unexpected size")
				}
				readFloat64, err := data.ReadFloat64("", 64)
				if err != nil {
					return nil, err
				}
				list = append(list, spiValues.NewPlcLREAL(readFloat64))
			case readWriteModel.CIPDataTypeCode_BOOL:
				bit, err := data.ReadBit("")
				if err != nil {
					return nil, err
				}
				list = append(list, spiValues.NewPlcBOOL(bit))
			default:
				return nil, errors.Errorf("Unknown type %v", _type)
			}
		}
		return spiValues.NewPlcList(list), nil
	} else {
		switch _type {
		case readWriteModel.CIPDataTypeCode_SINT:
			readByte, err := data.ReadInt8("", _type.Size()*8)
			if err != nil {
				return nil, err
			}
			return spiValues.NewPlcSINT(readByte), nil
		case readWriteModel.CIPDataTypeCode_INT:
			readInt16, err := data.ReadInt16("", _type.Size()*8)
			if err != nil {
				return nil, err
			}
			return spiValues.NewPlcINT(readInt16), nil
		case readWriteModel.CIPDataTypeCode_DINT:
			readInt32, err := data.ReadInt32("", _type.Size()*8)
			if err != nil {
				return nil, err
			}
			return spiValues.NewPlcDINT(readInt32), nil
		case readWriteModel.CIPDataTypeCode_REAL:
			if _type.Size()*8 != 64 {
				panic("Unexpected size")
			}
			readFloat32, err := data.ReadFloat32("", 64)
			if err != nil {
				return nil, err
			}
			return spiValues.NewPlcREAL(readFloat32), nil
		case readWriteModel.CIPDataTypeCode_BOOL:
			readBit, err := data.ReadBit("")
			if err != nil {
				return nil, err
			}
			return spiValues.NewPlcBOOL(readBit), nil
		default:
			return nil, errors.Errorf("Unknown type %v", _type)
		}
	}
}

// Helper to convert the return codes returned from the eip into one of our standard
func decodeResponseCode(status uint8) model.PlcResponseCode {
	//TODO other status
	switch status {
	case 0:
		return model.PlcResponseCode_OK
	default:
		return model.PlcResponseCode_INTERNAL_ERROR
	}
}
