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
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/opcua/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func generateNodeId(tag Tag) (readWriteModel.NodeId, error) {
	var nodeId readWriteModel.NodeId
	if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_BINARY_IDENTIFIER {
		parsedIdentifier, err := strconv.ParseUint(tag.GetIdentifier(), 10, 8)
		if err != nil {
			return nil, errors.New("Error parsing identifier")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdTwoByte(uint8(parsedIdentifier)))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_NUMBER_IDENTIFIER {
		parsedIdentifier, err := strconv.ParseUint(tag.GetIdentifier(), 10, 32)
		if err != nil {
			return nil, errors.New("Error parsing identifier")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdNumeric( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), uint32(parsedIdentifier)))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_GUID_IDENTIFIER {
		guid, err := uuid.Parse(tag.GetIdentifier())
		if err != nil {
			return nil, errors.Wrap(err, "error parsing guid")
		}
		guidBytes, err := guid.MarshalBinary() // TODO: do we need to do flip it here?
		if err != nil {
			return nil, errors.Wrap(err, "error marshaling guid")
		}
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdGuid( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), guidBytes))
	} else if tag.GetIdentifierType() == readWriteModel.OpcuaIdentifierType_STRING_IDENTIFIER {
		nodeId = readWriteModel.NewNodeId(readWriteModel.NewNodeIdString( /*TODO: do we want to check for overflow?*/ uint16(tag.GetNamespace()), readWriteModel.NewPascalString(utils.ToPtr(tag.GetIdentifier()))))
	}
	return nodeId, nil
}

func getDateTime(dateTime int64) time.Time {
	return time.UnixMilli((dateTime - EPOCH_OFFSET) / 10000)
}

func readResponse(localLog zerolog.Logger, readRequestIn apiModel.PlcReadRequest, tagNames []string, results []readWriteModel.DataValue) (readRequest apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]apiValues.PlcValue) {
	readRequest = readRequestIn
	responseCodes = map[string]apiModel.PlcResponseCode{}
	values = map[string]apiValues.PlcValue{}
	count := 0
	for _, tagName := range tagNames {
		responseCode := apiModel.PlcResponseCode_OK
		var value apiValues.PlcValue
		if results[count].GetValueSpecified() {
			variant := results[count].GetValue()
			localLog.Trace().Type("type", variant).Msg("Response of type")
			switch variant := variant.(type) {
			case readWriteModel.VariantBoolean:
				array := variant.GetValue()
				boolValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					boolValues[i] = spiValues.NewPlcBOOL(t != 0)
				}
				value = spiValues.NewPlcList(boolValues)
			case readWriteModel.VariantSByte:
				array := variant.GetValue()
				value = spiValues.NewPlcRawByteArray(array)
			case readWriteModel.VariantByte:
				array := variant.GetValue()
				value = spiValues.NewPlcRawByteArray(array)
			case readWriteModel.VariantInt16:
				array := variant.GetValue()
				int16Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int16Values[i] = spiValues.NewPlcINT(t)
				}
				value = spiValues.NewPlcList(int16Values)
			case readWriteModel.VariantUInt16:
				array := variant.GetValue()
				uint16Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint16Values[i] = spiValues.NewPlcUINT(t)
				}
				value = spiValues.NewPlcList(uint16Values)
			case readWriteModel.VariantInt32:
				array := variant.GetValue()
				int32Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int32Values[i] = spiValues.NewPlcDINT(t)
				}
				value = spiValues.NewPlcList(int32Values)
			case readWriteModel.VariantUInt32:
				array := variant.GetValue()
				uint32Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint32Values[i] = spiValues.NewPlcUDINT(t)
				}
				value = spiValues.NewPlcList(uint32Values)
			case readWriteModel.VariantInt64:
				array := variant.GetValue()
				int64Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					int64Values[i] = spiValues.NewPlcLINT(t)
				}
				value = spiValues.NewPlcList(int64Values)
			case readWriteModel.VariantUInt64:
				array := variant.GetValue()
				uint64Values := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					uint64Values[i] = spiValues.NewPlcULINT(t)
				}
				value = spiValues.NewPlcList(uint64Values)
			case readWriteModel.VariantFloat:
				array := variant.GetValue()
				floatValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					floatValues[i] = spiValues.NewPlcREAL(t)
				}
				value = spiValues.NewPlcList(floatValues)
			case readWriteModel.VariantDouble:
				array := variant.GetValue()
				doubleValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					doubleValues[i] = spiValues.NewPlcLREAL(t)
				}
				value = spiValues.NewPlcList(doubleValues)
			case readWriteModel.VariantString:
				array := variant.GetValue()
				stringValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					stringValues[i] = spiValues.NewPlcSTRING(*t.GetStringValue())
				}
				value = spiValues.NewPlcList(stringValues)
			case readWriteModel.VariantDateTime:
				array := variant.GetValue()
				dateTimeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					dateTimeValues[i] = spiValues.NewPlcDATE_AND_TIME(getDateTime(t))
				}
				value = spiValues.NewPlcList(dateTimeValues)
			case readWriteModel.VariantGuid:
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
			case readWriteModel.VariantXmlElement:
				array := variant.GetValue()
				xmlElementValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					xmlElementValues[i] = spiValues.NewPlcSTRING(*t.GetStringValue())
				}
				value = spiValues.NewPlcList(xmlElementValues)
			case readWriteModel.VariantLocalizedText:
				array := variant.GetValue()
				localizedTextValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					v := ""
					if t.GetLocaleSpecified() {
						v += *t.GetLocale().GetStringValue() + "|"
					}
					if t.GetTextSpecified() {
						v += *t.GetText().GetStringValue()
					}
					localizedTextValues[i] = spiValues.NewPlcSTRING(v)
				}
				value = spiValues.NewPlcList(localizedTextValues)
			case readWriteModel.VariantQualifiedName:
				array := variant.GetValue()
				qualifiedNameValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					qualifiedNameValues[i] = spiValues.NewPlcSTRING(fmt.Sprintf("ns=%d;s=%s", t.GetNamespaceIndex(), *t.GetName().GetStringValue()))
				}
				value = spiValues.NewPlcList(qualifiedNameValues)
			case readWriteModel.VariantExtensionObject:
				array := variant.GetValue()
				extensionObjectValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					extensionObjectValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(extensionObjectValues)
			case readWriteModel.VariantNodeId:
				array := variant.GetValue()
				nodeIdValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					nodeIdValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(nodeIdValues)
			case readWriteModel.VariantStatusCode:
				array := variant.GetValue()
				statusCodeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					statusCodeValues[i] = spiValues.NewPlcSTRING(t.String())
				}
				value = spiValues.NewPlcList(statusCodeValues)
			case readWriteModel.VariantByteString:
				array := variant.GetValue()
				statusCodeValues := make([]apiValues.PlcValue, len(array))
				for i, t := range array {
					statusCodeValues[i] = spiValues.NewPlcRawByteArray(t.GetValue())
				}
				value = spiValues.NewPlcList(statusCodeValues)
			default:
				responseCode = apiModel.PlcResponseCode_UNSUPPORTED
				localLog.Error().Type("variant", variant).Msg("Data type is not supported ")
			}
		} else {
			if results[count].GetStatusCode().GetStatusCode() == uint32(readWriteModel.OpcuaStatusCode_BadNodeIdUnknown) {
				responseCode = apiModel.PlcResponseCode_NOT_FOUND
			} else {
				responseCode = apiModel.PlcResponseCode_UNSUPPORTED
			}
			localLog.Error().Stringer("statusCode", results[count].GetStatusCode()).Msg("Error while reading value from OPC UA server error code")
		}
		count++
		responseCodes[tagName] = responseCode
		values[tagName] = value
	}
	return
}
