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

package readwrite

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// Code generated by code-generation. DO NOT EDIT.

type S7XmlParserHelper struct {
}

// Temporary imports to silent compiler warnings (TODO: migrate from static to emission based imports)
func init() {
	_ = strconv.ParseUint
	_ = strconv.ParseInt
	_ = strings.Join
	_ = utils.Dump
}

func (m S7XmlParserHelper) Parse(typeName string, xmlString string, parserArguments ...string) (interface{}, error) {
	switch typeName {
	case "DataItem":
		// TODO: find a way to parse the sub types
		var dataProtocolId string
		parsedInt1, err := strconv.ParseInt(parserArguments[1], 10, 32)
		if err != nil {
			return nil, err
		}
		stringLength := int32(parsedInt1)
		return model.DataItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), dataProtocolId, stringLength)
	case "SzlId":
		return model.SzlIdParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageObjectAckType":
		return model.AlarmMessageObjectAckTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageAckPushType":
		return model.AlarmMessageAckPushTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7Message":
		return model.S7MessageParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7VarPayloadStatusItem":
		return model.S7VarPayloadStatusItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7Parameter":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 8)
		if err != nil {
			return nil, err
		}
		messageType := uint8(parsedUint0)
		return model.S7ParameterParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), messageType)
	case "S7DataAlarmMessage":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 4)
		if err != nil {
			return nil, err
		}
		cpuFunctionType := uint8(parsedUint0)
		return model.S7DataAlarmMessageParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), cpuFunctionType)
	case "SzlDataTreeItem":
		return model.SzlDataTreeItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "COTPPacket":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 16)
		if err != nil {
			return nil, err
		}
		cotpLen := uint16(parsedUint0)
		return model.COTPPacketParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), cotpLen)
	case "S7PayloadUserDataItem":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 4)
		if err != nil {
			return nil, err
		}
		cpuFunctionType := uint8(parsedUint0)
		parsedUint1, err := strconv.ParseUint(parserArguments[1], 10, 8)
		if err != nil {
			return nil, err
		}
		cpuSubfunction := uint8(parsedUint1)
		return model.S7PayloadUserDataItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), cpuFunctionType, cpuSubfunction)
	case "DateAndTime":
		return model.DateAndTimeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "COTPParameter":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 8)
		if err != nil {
			return nil, err
		}
		rest := uint8(parsedUint0)
		return model.COTPParameterParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), rest)
	case "AlarmMessageObjectPushType":
		return model.AlarmMessageObjectPushTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "State":
		return model.StateParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessagePushType":
		return model.AlarmMessagePushTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "TPKTPacket":
		return model.TPKTPacketParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageAckType":
		return model.AlarmMessageAckTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AssociatedValueType":
		return model.AssociatedValueTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageAckObjectPushType":
		return model.AlarmMessageAckObjectPushTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7Payload":
		parsedUint0, err := strconv.ParseUint(parserArguments[0], 10, 8)
		if err != nil {
			return nil, err
		}
		messageType := uint8(parsedUint0)
		// TODO: find a way to parse the sub types
		var parameter model.S7Parameter
		return model.S7PayloadParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)), messageType, parameter)
	case "S7VarRequestParameterItem":
		return model.S7VarRequestParameterItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7VarPayloadDataItem":
		return model.S7VarPayloadDataItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageQueryType":
		return model.AlarmMessageQueryTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageAckResponseType":
		return model.AlarmMessageAckResponseTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "AlarmMessageObjectQueryType":
		return model.AlarmMessageObjectQueryTypeParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7Address":
		return model.S7AddressParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	case "S7ParameterUserDataItem":
		return model.S7ParameterUserDataItemParse(utils.NewXmlReadBuffer(strings.NewReader(xmlString)))
	}
	return nil, errors.Errorf("Unsupported type %s", typeName)
}
