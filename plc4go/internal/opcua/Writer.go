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

package opcua

import (
	"context"
	"encoding/binary"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"runtime/debug"
	"strconv"
)

type Writer struct {
	messageCodec *MessageCodec

	log zerolog.Logger
}

func NewWriter(messageCodec *MessageCodec, _options ...options.WithOption) *Writer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Writer{
		messageCodec: messageCodec,

		log: customLogger,
	}
}

func (m *Writer) Write(ctx context.Context, writeRequest apiModel.PlcWriteRequest) <-chan apiModel.PlcWriteRequestResult {
	m.log.Trace().Msg("Writing")
	result := make(chan apiModel.PlcWriteRequestResult, 1)
	go m.WriteSync(ctx, writeRequest, result)
	return result
}

func (m *Writer) WriteSync(ctx context.Context, writeRequest apiModel.PlcWriteRequest, result chan apiModel.PlcWriteRequestResult) {
	defer func() {
		if err := recover(); err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
		}
	}()

	requestHeader := readWriteModel.NewRequestHeader(
		m.messageCodec.channel.getAuthenticationToken(),
		m.messageCodec.channel.getCurrentDateTime(),
		m.messageCodec.channel.getRequestHandle(),
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT,
	)
	writeValueArray := make([]readWriteModel.ExtensionObjectDefinition, len(writeRequest.GetTagNames()))
	for i, tagName := range writeRequest.GetTagNames() {
		tag := writeRequest.GetTag(tagName).(Tag)

		nodeId, err := generateNodeId(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "error generating node id from tag %s", tag))
			return
		}

		plcValue, err := m.fromPlcValue(tagName, tag, writeRequest)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Error getting plcValue"))
			return
		}
		writeValueArray[i] = readWriteModel.NewWriteValue(nodeId,
			0xD,
			NULL_STRING,
			readWriteModel.NewDataValue(
				false,
				false,
				false,
				false,
				false,
				true,
				plcValue,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
		)
	}

	opcuaWriteRequest := readWriteModel.NewWriteRequest(
		requestHeader,
		int32(len(writeValueArray)),
		writeValueArray,
	)

	identifier, err := strconv.ParseUint(opcuaWriteRequest.GetIdentifier(), 10, 16)
	if err != nil {
		result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "error parsing identifier"))
		return
	}
	expandedNodeId := readWriteModel.NewExpandedNodeId(false, //Namespace Uri Specified
		false, //Server Index Specified
		readWriteModel.NewNodeIdFourByte(0, uint16(identifier)),
		nil,
		nil)

	extObject := readWriteModel.NewExtensionObject(
		expandedNodeId,
		nil,
		opcuaWriteRequest,
		false)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err = extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Unable to serialise the ReadRequest"))
		return
	}

	consumer := func(opcuaResponse []byte) {
		reply, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, errors.Wrapf(err, "Unable to read the reply"))
			return
		}
		if writeResponse, ok := reply.(readWriteModel.WriteResponseExactly); ok {
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, spiModel.NewDefaultPlcWriteResponse(m.writeResponse(writeRequest, writeResponse.GetResults())), nil)
			return
		} else {
			if serviceFault, ok := reply.(readWriteModel.ServiceFaultExactly); ok {
				header := serviceFault.GetResponseHeader()
				m.log.Error().Msgf("Read request ended up with ServiceFault: %s", header)
			} else {
				m.log.Error().Msgf("Remote party returned an error '%s'", reply)
			}

			responseCodes := map[string]apiModel.PlcResponseCode{}
			for _, tagName := range writeRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
			result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, spiModel.NewDefaultPlcWriteResponse(writeRequest, responseCodes), nil)
		}
	}

	errorDispatcher := func(err error) {
		result <- spiModel.NewDefaultPlcWriteRequestResult(writeRequest, nil, err)
	}

	m.messageCodec.channel.submit(ctx, m.messageCodec, errorDispatcher, consumer, buffer)
}

func (m *Writer) writeResponse(requestIn apiModel.PlcWriteRequest, results []readWriteModel.StatusCode) (request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) {
	request = requestIn
	responseCodes = map[string]apiModel.PlcResponseCode{}
	for i, tagName := range request.GetTagNames() {
		statusCode := results[i].GetStatusCode()
		switch readWriteModel.OpcuaStatusCode(statusCode) {
		case readWriteModel.OpcuaStatusCode_Good:
			responseCodes[tagName] = apiModel.PlcResponseCode_OK
		case readWriteModel.OpcuaStatusCode_BadNodeIdUnknown:
			responseCodes[tagName] = apiModel.PlcResponseCode_NOT_FOUND
		default:
			responseCodes[tagName] = apiModel.PlcResponseCode_REMOTE_ERROR
		}
	}
	return
}

func (m *Writer) fromPlcValue(tagName string, tag Tag, request apiModel.PlcWriteRequest) (readWriteModel.Variant, error) {
	var valueObject spiValues.PlcList
	if value := request.GetValue(tagName); value.IsList() {
		valueObject = value.(spiValues.PlcList)
	} else {
		valueObject = spiValues.NewPlcList([]apiValues.PlcValue{value})
	}

	plcValueList := valueObject.GetList()
	dataType := tag.GetValueType()
	if dataType == (apiValues.NULL) {
		dataType = plcValueList[0].GetPlcValueType()
	}
	length := valueObject.GetLength()
	switch dataType {
	// Simple boolean values
	case apiValues.BOOL:
		tmpBOOL := make([]byte, length)
		for i := uint32(0); i < length; i++ {
			tmpBOOL[i] = valueObject.GetIndex(i).GetByte()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantBoolean(arrayLength, tmpBOOL, false, false, nil, nil), nil

	// 8-Bit Bit-Strings (Groups of Boolean Values)
	case apiValues.BYTE:
		tmpBYTE := make([]byte, length)
		for i := uint32(0); i < length; i++ {
			tmpBYTE[i] = valueObject.GetIndex(i).GetByte()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantByte(arrayLength, tmpBYTE, false, false, nil, nil), nil

	// 16-Bit Bit-Strings (Groups of Boolean Values)
	case apiValues.WORD:
		tmpWORD := make([]uint16, length)
		for i := uint32(0); i < length; i++ {
			tmpWORD[i] = valueObject.GetIndex(i).GetUint16()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt16(arrayLength, tmpWORD, false, false, nil, nil), nil

	// 32-Bit Bit-Strings (Groups of Boolean Values)
	case apiValues.DWORD:
		tmpDWORD := make([]uint32, length)
		for i := uint32(0); i < length; i++ {
			tmpDWORD[i] = valueObject.GetIndex(i).GetUint32()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt32(arrayLength, tmpDWORD, false, false, nil, nil), nil

	// 64-Bit Bit-Strings (Groups of Boolean Values)
	case apiValues.LWORD:
		tmpLWORD := make([]uint64, length)
		for i := uint32(0); i < length; i++ {
			tmpLWORD[i] = valueObject.GetIndex(i).GetUint64()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt64(arrayLength, tmpLWORD, false, false, nil, nil), nil

	// 8-Bit Unsigned Integers
	case apiValues.USINT:
		tmpUSINT := make([]byte, length)
		for i := uint32(0); i < length; i++ {
			tmpUSINT[i] = valueObject.GetIndex(i).GetByte()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantByte(arrayLength, tmpUSINT, false, false, nil, nil), nil

	// 8-Bit Signed Integers
	case apiValues.SINT:
		tmpSINT := make([]byte, length)
		for i := uint32(0); i < length; i++ {
			tmpSINT[i] = valueObject.GetIndex(i).GetByte()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantSByte(arrayLength, tmpSINT, false, false, nil, nil), nil

	// 16-Bit Unsigned Integers
	case apiValues.UINT:
		tmpUINT := make([]uint16, length)
		for i := uint32(0); i < length; i++ {
			tmpUINT[i] = valueObject.GetIndex(i).GetUint16()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt16(arrayLength, tmpUINT, false, false, nil, nil), nil

	// 16-Bit Signed Integers
	case apiValues.INT:
		tmpINT := make([]int16, length)
		for i := uint32(0); i < length; i++ {
			tmpINT[i] = valueObject.GetIndex(i).GetInt16()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantInt16(arrayLength, tmpINT, false, false, nil, nil), nil

	// 32-Bit Unsigned Integers
	case apiValues.UDINT:
		tmpUDINT := make([]uint32, length)
		for i := uint32(0); i < length; i++ {
			tmpUDINT[i] = valueObject.GetIndex(i).GetUint32()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt32(arrayLength, tmpUDINT, false, false, nil, nil), nil

	// 32-Bit Signed Integers
	case apiValues.DINT:
		tmpDINT := make([]int32, length)
		for i := uint32(0); i < length; i++ {
			tmpDINT[i] = valueObject.GetIndex(i).GetInt32()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantInt32(arrayLength, tmpDINT, false, false, nil, nil), nil

	// 64-Bit Unsigned Integers
	case apiValues.ULINT:
		tmpULINT := make([]uint64, length)
		for i := uint32(0); i < length; i++ {
			tmpULINT[i] = valueObject.GetIndex(i).GetUint64()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantUInt64(arrayLength, tmpULINT, false, false, nil, nil), nil

	// 64-Bit Signed Integers
	case apiValues.LINT:
		tmpUINT := make([]int64, length)
		for i := uint32(0); i < length; i++ {
			tmpUINT[i] = valueObject.GetIndex(i).GetInt64()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantInt64(arrayLength, tmpUINT, false, false, nil, nil), nil

	// 32-Bit Floating Point Values
	case apiValues.REAL:
		tmpREAL := make([]float32, length)
		for i := uint32(0); i < length; i++ {
			tmpREAL[i] = valueObject.GetIndex(i).GetFloat32()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantFloat(arrayLength, tmpREAL, false, false, nil, nil), nil

	// 64-Bit Floating Point Values
	case apiValues.LREAL:
		tmpLREAL := make([]float64, length)
		for i := uint32(0); i < length; i++ {
			tmpLREAL[i] = valueObject.GetIndex(i).GetFloat64()
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantDouble(arrayLength, tmpLREAL, false, false, nil, nil), nil

	// UTF-8 Characters and Strings
	case apiValues.CHAR:
		fallthrough
	case apiValues.STRING:
		fallthrough
		// UTF-16 Characters and Strings
	case apiValues.WCHAR:
		fallthrough
	case apiValues.WSTRING:
		tmpString := make([]readWriteModel.PascalString, length)
		for i := uint32(0); i < length; i++ {
			tmpString[i] = readWriteModel.NewPascalString(valueObject.GetIndex(i).GetString())
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantString(arrayLength, tmpString, false, false, nil, nil), nil

	case apiValues.DATE_AND_TIME:
		tmpDateTime := make([]int64, length)
		for i := uint32(0); i < length; i++ {
			tmpDateTime[i] = valueObject.GetIndex(i).GetDateTime().UnixMilli() / 1000
		}
		var arrayLength *int32
		if length != 1 {
			int32Length := int32(length)
			arrayLength = &int32Length
		}
		return readWriteModel.NewVariantDateTime(arrayLength, tmpDateTime, false, false, nil, nil), nil
	default:
		return nil, errors.Errorf("Unsupported write tag type %s", dataType)
	}
}
