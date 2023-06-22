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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/rs/zerolog"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/pkg/errors"
)

type Reader struct {
	messageCodec  spi.MessageCodec
	tm            transactions.RequestTransactionManager
	configuration Configuration
	sessionHandle *uint32

	log zerolog.Logger
}

func NewReader(messageCodec spi.MessageCodec, tm transactions.RequestTransactionManager, configuration Configuration, sessionHandle *uint32, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		messageCodec:  messageCodec,
		tm:            tm,
		configuration: configuration,
		sessionHandle: sessionHandle,

		log: customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	// TODO: handle ctx
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		classSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewClassID(0, 6))
		instanceSegment := readWriteModel.NewLogicalSegment(readWriteModel.NewInstanceID(0, 1))
		for _, tagName := range readRequest.GetTagNames() {
			plcTag := readRequest.GetTag(tagName).(PlcTag)
			tag := plcTag.GetTag()
			elementsNb := uint16(1)
			if plcTag.GetElementNb() > 1 {
				elementsNb = plcTag.GetElementNb()
			}
			ansi, err := toAnsi(tag)
			if err != nil {
				result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "Error encoding eip ansi for tag %s", tagName))
				return
			}
			requestItem := readWriteModel.NewCipUnconnectedRequest(classSegment, instanceSegment,
				readWriteModel.NewCipReadRequest(ansi, elementsNb, 0),
				m.configuration.backplane, m.configuration.slot, uint16(0))
			typeIds := []readWriteModel.TypeId{
				readWriteModel.NewNullAddressItem(),
				readWriteModel.NewUnConnectedDataItem(requestItem),
			}
			request := readWriteModel.NewCipRRData(0, 0, typeIds, *m.sessionHandle, uint32(readWriteModel.CIPStatus_Success), []byte(DefaultSenderContext), 0)
			transaction := m.tm.StartTransaction()
			transaction.Submit(func(transaction transactions.RequestTransaction) {
				if err := m.messageCodec.SendRequest(
					ctx,
					request,
					func(message spi.Message) bool {
						eipPacket := message.(readWriteModel.EipPacketExactly)
						if eipPacket == nil {
							return false
						}
						cipRRData := eipPacket.(readWriteModel.CipRRDataExactly)
						if cipRRData == nil {
							return false
						}
						return cipRRData.GetSessionHandle() == *m.sessionHandle
					},
					func(message spi.Message) error {
						cipRRData := message.(readWriteModel.CipRRData)
						m.log.Trace().Msgf("handling:\n%s", cipRRData)
						unconnectedDataItem := cipRRData.GetTypeIds()[1].(readWriteModel.UnConnectedDataItem)
						// Convert the eip response into a PLC4X response
						m.log.Trace().Msg("convert response to PLC4X response")
						readResponse, err := m.ToPlc4xReadResponse(unconnectedDataItem.GetService(), readRequest)
						if err != nil {
							result <- spiModel.NewDefaultPlcReadRequestResult(
								readRequest,
								nil,
								errors.Wrap(err, "Error decoding response"),
							)
							return transaction.EndRequest()
						}
						result <- spiModel.NewDefaultPlcReadRequestResult(
							readRequest,
							readResponse,
							nil,
						)
						return transaction.EndRequest()
					},
					func(err error) error {
						result <- spiModel.NewDefaultPlcReadRequestResult(
							readRequest,
							nil,
							errors.Wrap(err, "got timeout while waiting for response"),
						)
						return transaction.EndRequest()
					},
					time.Second*1,
				); err != nil {
					result <- spiModel.NewDefaultPlcReadRequestResult(
						readRequest,
						nil,
						errors.Wrap(err, "error sending message"),
					)
					if err := transaction.FailRequest(errors.Errorf("timeout after %s", time.Second*1)); err != nil {
						m.log.Debug().Err(err).Msg("Error failing request")
					}
				}
			})
		}
	}()
	return result
}

func toAnsi(tag string) ([]byte, error) {
	resourceAddressPattern := regexp.MustCompile("([.\\[\\]])*([A-Za-z_0-9]+){1}")

	segments := make([]readWriteModel.PathSegment, 0)
	lengthInBytes := uint16(0)
	submatch := resourceAddressPattern.FindAllStringSubmatch(tag, -1)
	for _, match := range submatch {
		identifier := match[2]
		qualifier := match[1]

		var newSegment readWriteModel.PathSegment
		if len(qualifier) > 0 {
			if qualifier == "[" {
				numericIdentifier, err := strconv.Atoi(identifier)
				if err != nil {
					return nil, fmt.Errorf("error parsing address %s, identifier %s couldn't be parsed to an integer", tag, identifier)
				}
				newSegment = readWriteModel.NewLogicalSegment(readWriteModel.NewMemberID(0, uint8(numericIdentifier)))
			} else {
				newSegment = readWriteModel.NewDataSegment(readWriteModel.NewAnsiExtendedSymbolSegment(identifier, nil))
			}
		} else {
			var pad *uint8
			if len(identifier)%2 != 0 {
				paddingValue := uint8(0)
				pad = &paddingValue
			}
			newSegment = readWriteModel.NewDataSegment(readWriteModel.NewAnsiExtendedSymbolSegment(identifier, pad))
		}
		lengthInBytes += newSegment.GetLengthInBytes(context.Background())
		segments = append(segments, newSegment)
	}
	buffer := utils.NewWriteBufferByteBased(
		utils.WithInitialSizeForByteBasedBuffer(int(lengthInBytes)),
		utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	for _, segment := range segments {
		if err := segment.SerializeWithWriteBuffer(context.Background(), buffer); err != nil {
			return nil, errors.Wrap(err, "error converting tag to ansi")
		}
	}
	return buffer.GetBytes(), nil
}

func (m *Reader) ToPlc4xReadResponse(response readWriteModel.CipService, readRequest apiModel.PlcReadRequest) (apiModel.PlcReadResponse, error) {
	plcValues := map[string]values.PlcValue{}
	responseCodes := map[string]apiModel.PlcResponseCode{}
	switch response := response.(type) {
	case readWriteModel.CipReadResponseExactly: // only 1 tag
		cipReadResponse := response
		tagName := readRequest.GetTagNames()[0]
		tag := readRequest.GetTag(tagName).(PlcTag)
		code := decodeResponseCode(cipReadResponse.GetStatus())
		var plcValue values.PlcValue
		_type := cipReadResponse.GetData().GetDataType()
		data := utils.NewReadBufferByteBased(cipReadResponse.GetData().GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		if code == apiModel.PlcResponseCode_OK {
			var err error
			plcValue, err = parsePlcValue(tag, data, _type)
			if err != nil {
				return nil, err
			}
		}
		plcValues[tagName] = plcValue
		responseCodes[tagName] = code
	case readWriteModel.MultipleServiceResponseExactly: //Multiple response
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
			tag := readRequest.GetTag(tagName).(PlcTag)
			if cipReadResponse, ok := services.Services[i].(readWriteModel.CipReadResponse); ok {
				code := decodeResponseCode(cipReadResponse.GetStatus())
				_type := cipReadResponse.GetData().GetDataType()
				data := utils.NewReadBufferByteBased(cipReadResponse.GetData().GetData(), utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
				var plcValue values.PlcValue
				if code == apiModel.PlcResponseCode_OK {
					var err error
					plcValue, err = parsePlcValue(tag, data, _type)
					if err != nil {
						return nil, err
					}
				}

				plcValues[tagName] = plcValue
				responseCodes[tagName] = code
			} else {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
		}
	default:
		return nil, errors.Errorf("unsupported response type %T", response)
	}

	// Return the response
	m.log.Trace().Msg("Returning the response")
	return spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, plcValues), nil
}

func parsePlcValue(tag PlcTag, data utils.ReadBufferByteBased, _type readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
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
					return nil, errors.New("Unexpected size")
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
				return nil, errors.New("Unexpected size")
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
func decodeResponseCode(status uint8) apiModel.PlcResponseCode {
	//TODO other status
	switch status {
	case 0:
		return apiModel.PlcResponseCode_OK
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
}
