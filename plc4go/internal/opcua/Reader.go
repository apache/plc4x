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
	"fmt"
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

type Reader struct {
	messageCodec *MessageCodec

	log zerolog.Logger
}

func NewReader(messageCodec *MessageCodec, _options ...options.WithOption) *Reader {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Reader{
		messageCodec: messageCodec,

		log: customLogger,
	}
}

func (m *Reader) Read(ctx context.Context, readRequest apiModel.PlcReadRequest) <-chan apiModel.PlcReadRequestResult {
	m.log.Trace().Msg("Reading")
	result := make(chan apiModel.PlcReadRequestResult, 1)
	go m.readSync(ctx, readRequest, result)
	return result
}

func (m *Reader) readSync(ctx context.Context, readRequest apiModel.PlcReadRequest, result chan apiModel.PlcReadRequestResult) {
	defer func() {
		if err := recover(); err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
		}
	}()

	requestHeader := readWriteModel.NewRequestHeader(
		m.messageCodec.channel.getAuthenticationToken(),
		m.messageCodec.channel.getCurrentDateTime(),
		0,
		0,
		NULL_STRING,
		REQUEST_TIMEOUT_LONG,
		NULL_EXTENSION_OBJECT,
	)
	readValueArray := make([]readWriteModel.ExtensionObjectDefinition, len(readRequest.GetTagNames()))
	for i, tagName := range readRequest.GetTagNames() {
		tag := readRequest.GetTag(tagName).(Tag)

		nodeId, err := generateNodeId(tag)
		if err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "error generating node id from tag %s", tag))
			return
		}

		readValueArray[i] = readWriteModel.NewReadValueId(nodeId,
			0xD,
			NULL_STRING,
			readWriteModel.NewQualifiedName(0, NULL_STRING),
		)
	}

	opcuaReadRequest := readWriteModel.NewReadRequest(
		requestHeader,
		0.0,
		readWriteModel.TimestampsToReturn_timestampsToReturnNeither,
		int32(len(readValueArray)),
		readValueArray)
	identifier, err := strconv.ParseUint(opcuaReadRequest.GetIdentifier(), 10, 16)
	if err != nil {
		result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "error parsing identifier"))
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
		opcuaReadRequest,
		false)

	buffer := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err = extObject.SerializeWithWriteBuffer(ctx, buffer); err != nil {
		result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "Unable to serialise the ReadRequest"))
		return
	}

	consumer := func(opcuaResponse []byte) {
		reply, err := readWriteModel.ExtensionObjectParseWithBuffer(ctx, utils.NewReadBufferByteBased(opcuaResponse, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)), false)
		if err != nil {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, errors.Wrapf(err, "Unable to read the reply"))
			return
		}
		if readResponse, ok := reply.(readWriteModel.ReadResponseExactly); ok {
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, spiModel.NewDefaultPlcReadResponse(m.readResponse(readRequest, readResponse.GetResults())), nil)
			return
		} else {
			if serviceFault, ok := reply.(readWriteModel.ServiceFaultExactly); ok {
				header := serviceFault.GetResponseHeader()
				m.log.Error().Msgf("Read request ended up with ServiceFault: %s", header)
			} else {
				m.log.Error().Msgf("Remote party returned an error '%s'", reply)
			}

			responseCodes := map[string]apiModel.PlcResponseCode{}
			for _, tagName := range readRequest.GetTagNames() {
				responseCodes[tagName] = apiModel.PlcResponseCode_INTERNAL_ERROR
			}
			result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, spiModel.NewDefaultPlcReadResponse(readRequest, responseCodes, nil), nil)
		}
	}

	errorDispatcher := func(err error) {
		result <- spiModel.NewDefaultPlcReadRequestResult(readRequest, nil, err)
	}

	m.messageCodec.channel.submit(ctx, m.messageCodec, errorDispatcher, result, consumer, buffer)
}

func (m *Reader) readResponse(readRequestIn apiModel.PlcReadRequest, results []readWriteModel.DataValue) (readRequest apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]apiValues.PlcValue) {
	readRequest = readRequestIn
	responseCodes = map[string]apiModel.PlcResponseCode{}
	values = map[string]apiValues.PlcValue{}
	count := 0
	for _, tagName := range readRequest.GetTagNames() {
		responseCode := apiModel.PlcResponseCode_OK
		var value apiValues.PlcValue
		if results[count].GetValueSpecified() {
			variant := results[count].GetValue()
			m.log.Trace().Msgf("Response of type %T", variant)
			switch variant := variant.(type) {
			case readWriteModel.VariantBooleanExactly:
				array := variant.GetValue()
				boolValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					boolValues[i] = spiValues.NewPlcBOOL(t != 0)
				}
				value = spiValues.NewPlcList(boolValues)
			case readWriteModel.VariantSByteExactly:
				array := variant.GetValue()
				value = spiValues.NewPlcRawByteArray(array)
			case readWriteModel.VariantByteExactly:
				array := variant.GetValue()
				value = spiValues.NewPlcRawByteArray(array)
			case readWriteModel.VariantInt16Exactly:
				array := variant.GetValue()
				int16Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int16Values[i] = spiValues.NewPlcINT(t)
				}
				value = spiValues.NewPlcList(int16Values)
			case readWriteModel.VariantUInt16Exactly:
				array := variant.GetValue()
				uint16Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint16Values[i] = spiValues.NewPlcUINT(t)
				}
				value = spiValues.NewPlcList(uint16Values)
			case readWriteModel.VariantInt32Exactly:
				array := variant.GetValue()
				int32Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int32Values[i] = spiValues.NewPlcDINT(t)
				}
				value = spiValues.NewPlcList(int32Values)
			case readWriteModel.VariantUInt32Exactly:
				array := variant.GetValue()
				uint32Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint32Values[i] = spiValues.NewPlcUDINT(t)
				}
				value = spiValues.NewPlcList(uint32Values)
			case readWriteModel.VariantInt64Exactly:
				array := variant.GetValue()
				int64Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int64Values[i] = spiValues.NewPlcLINT(t)
				}
				value = spiValues.NewPlcList(int64Values)
			case readWriteModel.VariantUInt64Exactly:
				array := variant.GetValue()
				uint64Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint64Values[i] = spiValues.NewPlcULINT(t)
				}
				value = spiValues.NewPlcList(uint64Values)
			case readWriteModel.VariantFloatExactly:
				array := variant.GetValue()
				floatValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					floatValues[i] = spiValues.NewPlcREAL(t)
				}
				value = spiValues.NewPlcList(floatValues)
			case readWriteModel.VariantDoubleExactly:
				array := variant.GetValue()
				doubleValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					doubleValues[i] = spiValues.NewPlcLREAL(t)
				}
				value = spiValues.NewPlcList(doubleValues)
			case readWriteModel.VariantStringExactly:
				array := variant.GetValue()
				stringValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					stringValues[i] = spiValues.NewPlcSTRING(t.GetStringValue())
				}
				value = spiValues.NewPlcList(stringValues)
			case readWriteModel.VariantDateTimeExactly:
				array := variant.GetValue()
				dateTimeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					dateTimeValues[i] = spiValues.NewPlcDATE_AND_TIME(getDateTime(t))
				}
				value = spiValues.NewPlcList(dateTimeValues)
			case readWriteModel.VariantGuidExactly:
				array := variant.GetValue()
				guidValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					//These two data section aren't little endian like the rest.
					data4Bytes := t.GetData4()
					data4 := 0
					for _, data4Byte := range data4Bytes {
						data4 = (data4 << 8) + (int(data4Byte) & 0xff)
					}
					data5Bytes := t.GetData5()
					data5 := 0
					for _, data5Byte := range data5Bytes {
						data5 = (data5 << 8) + (int(data5Byte) & 0xff)
					}
					guidValues[i] = spiValues.NewPlcSTRING(fmt.Sprintf("%x-%x-%x-%x-%x", t.GetData1(), t.GetData2(), t.GetData3(), data4, data5))
				}
				value = spiValues.NewPlcList(guidValues)
			case readWriteModel.VariantXmlElementExactly:
				array := variant.GetValue()
				xmlElementValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					xmlElementValues[i] = spiValues.NewPlcSTRING(t.GetStringValue())
				}
				value = spiValues.NewPlcList(xmlElementValues)
			case readWriteModel.VariantLocalizedTextExactly:
				array := variant.GetValue()
				localizedTextValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					v := ""
					if t.GetLocaleSpecified() {
						v += t.GetLocale().GetStringValue() + "|"
					}
					if t.GetTextSpecified() {
						v += t.GetText().GetStringValue()
					}
					localizedTextValues[i] = spiValues.NewPlcSTRING(v)
				}
				value = spiValues.NewPlcList(localizedTextValues)
			case readWriteModel.VariantQualifiedNameExactly:
				array := variant.GetValue()
				qualifiedNameValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					qualifiedNameValues[i] = spiValues.NewPlcSTRING(fmt.Sprintf("ns=%d;s=%s", t.GetNamespaceIndex(), t.GetName().GetStringValue()))
				}
				value = spiValues.NewPlcList(qualifiedNameValues)
			case readWriteModel.VariantExtensionObjectExactly:
				array := variant.GetValue()
				extensionObjectValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					extensionObjectValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(extensionObjectValues)
			case readWriteModel.VariantNodeIdExactly:
				array := variant.GetValue()
				nodeIdValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					nodeIdValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(nodeIdValues)
			case readWriteModel.VariantStatusCodeExactly:
				array := variant.GetValue()
				statusCodeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					statusCodeValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(statusCodeValues)
			case readWriteModel.VariantByteStringExactly:
				array := variant.GetValue()
				statusCodeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					statusCodeValues[i] = spiValues.NewPlcRawByteArray(t.GetValue())
				}
				value = spiValues.NewPlcList(statusCodeValues)
			default:
				responseCode = apiModel.PlcResponseCode_UNSUPPORTED
				m.log.Error().Msgf("Data type - %T is not supported ", variant)
			}
		} else {
			if results[count].GetStatusCode().GetStatusCode() == uint32(readWriteModel.OpcuaStatusCode_BadNodeIdUnknown) {
				responseCode = apiModel.PlcResponseCode_NOT_FOUND
			} else {
				responseCode = apiModel.PlcResponseCode_UNSUPPORTED
			}
			m.log.Error().Msgf("Error while reading value from OPC UA server error code:- %s", results[count].GetStatusCode())
		}
		count++
		responseCodes[tagName] = responseCode
		values[tagName] = value
	}
	return
}
