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

package eip

import (
	"context"
	"encoding/binary"
	"strconv"
	"strings"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	plc4goModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func (m *Reader) Read(ctx context.Context, readRequest model.PlcReadRequest) <-chan model.PlcReadRequestResult {
	// TODO: handle ctx
	log.Trace().Msg("Reading")
	result := make(chan model.PlcReadRequestResult)
	go func() {
		classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6))
		instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 1))
		for _, tagName := range readRequest.GetTagNames() {
			plcTag := readRequest.GetTag(tagName).(EIPPlcTag)
			tag := plcTag.GetTag()
			elementsNb := uint16(1)
			if plcTag.GetElementNb() > 1 {
				elementsNb = plcTag.GetElementNb()
			}
			ansi, err := toAnsi(tag)
			if err != nil {
				result <- &plc4goModel.DefaultPlcReadRequestResult{
					Request:  readRequest,
					Response: nil,
					Err:      errors.Wrapf(err, "Error encoding eip ansi for tag %s", tagName),
				}
				return
			}
			request := readWriteModel.NewCipReadRequest(ansi, elementsNb, uint16(0))
			requestItem := readWriteModel.NewCipUnconnectedRequest(classSegment, instanceSegment, request,
				m.configuration.backplane, m.configuration.slot, uint16(0))
			typeIds := []readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(requestItem),
			}
			pkt := readWriteModel.NewCipRRData(0, 0, typeIds, *m.sessionHandle,
				uint32(readWriteModel.CIPStatus_Success), []byte(DefaultSenderContext), 0)

			transaction := m.tm.StartTransaction()
			transaction.Submit(func() {
				if err := m.messageCodec.SendRequest(ctx, pkt,
					func(message spi.Message) bool {
						eipPacket := message.(readWriteModel.EipPacket)
						if eipPacket == nil {
							return false
						}
						cipRRData := eipPacket.(readWriteModel.CipRRData)
						if cipRRData == nil {
							return false
						}
						return cipRRData.GetSessionHandle() == *m.sessionHandle
					},
					func(message spi.Message) error {
						cipRRData := message.(readWriteModel.CipRRData)
						unconnectedDataItem := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
						// Convert the eip response into a PLC4X response
						log.Trace().Msg("convert response to PLC4X response")
						readResponse, err := m.ToPlc4xReadResponse(unconnectedDataItem.GetService(), readRequest)
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
					}, time.Second*1); err != nil {
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

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))

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

func (m *Reader) ToPlc4xReadResponse(response readWriteModel.CipService, readRequest model.PlcReadRequest) (model.PlcReadResponse, error) {
	plcValues := map[string]values.PlcValue{}
	responseCodes := map[string]model.PlcResponseCode{}
	switch response := response.(type) {
	case readWriteModel.CipReadResponse: // only 1 tag
		cipReadResponse := response
		tagName := readRequest.GetTagNames()[0]
		tag := readRequest.GetTag(tagName).(EIPPlcTag)
		code := decodeResponseCode(cipReadResponse.GetStatus())
		var plcValue values.PlcValue
		_type := cipReadResponse.GetData().GetDataType()
		data := utils.NewReadBufferByteBased(cipReadResponse.GetData().GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		if code == model.PlcResponseCode_OK {
			var err error
			plcValue, err = parsePlcValue(tag, data, _type)
			if err != nil {
				return nil, err
			}
		}
		plcValues[tagName] = plcValue
		responseCodes[tagName] = code
	case readWriteModel.MultipleServiceResponse: //Multiple response
		multipleServiceResponse := response
		nb := multipleServiceResponse.GetServiceNb()
		arr := make([]readWriteModel.CipService, nb)
		read := utils.NewReadBufferByteBased(multipleServiceResponse.GetServicesData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		total := read.GetTotalBytes()
		for i := uint16(0); i < nb; i++ {
			length := uint16(0)
			offset := multipleServiceResponse.GetOffsets()[i] - multipleServiceResponse.GetOffsets()[0] //Substract first offset as we only have the service in the buffer (not servicesNb and offsets)
			if i == nb-1 {
				length = uint16(total) - offset //Get the rest if last
			} else {
				length = multipleServiceResponse.GetOffsets()[i+1] - offset - multipleServiceResponse.GetOffsets()[0] //Calculate length with offsets (substracting first offset)
			}
			serviceBuf := utils.NewReadBufferByteBased(read.GetBytes()[offset:offset+length], utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
			var err error
			// TODO: If we're using a connected connection, do this differently
			arr[i], err = readWriteModel.CipServiceParseWithBuffer(context.Background(), serviceBuf, false, length)
			if err != nil {
				return nil, err
			}
		}
		services := readWriteModel.NewServices(multipleServiceResponse.GetOffsets(), arr, uint16(0))
		for i, tagName := range readRequest.GetTagNames() {
			tag := readRequest.GetTag(tagName).(EIPPlcTag)
			if cipReadResponse, ok := services.Services[i].(readWriteModel.CipReadResponse); ok {
				code := decodeResponseCode(cipReadResponse.GetStatus())
				_type := cipReadResponse.GetData().GetDataType()
				data := utils.NewReadBufferByteBased(cipReadResponse.GetData().GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				var plcValue values.PlcValue
				if code == model.PlcResponseCode_OK {
					var err error
					plcValue, err = parsePlcValue(tag, data, _type)
					if err != nil {
						return nil, err
					}
				}

				plcValues[tagName] = plcValue
				responseCodes[tagName] = code
			} else {
				responseCodes[tagName] = model.PlcResponseCode_INTERNAL_ERROR
			}
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", response)
	}

	// Return the response
	log.Trace().Msg("Returning the response")
	return plc4goModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

func parsePlcValue(tag EIPPlcTag, data utils.ReadBufferByteBased, _type readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
	nb := tag.GetElementNb()
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
