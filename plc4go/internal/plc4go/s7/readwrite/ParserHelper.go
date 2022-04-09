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

package readwrite

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

type S7ParserHelper struct {
}

func (m S7ParserHelper) Parse(typeName string, arguments []string, io utils.ReadBuffer) (interface{}, error) {
	switch typeName {
	case "DataItem":
		dataProtocolId, err := utils.StrToString(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		stringLength, err := utils.StrToInt32(arguments[1])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.DataItemParse(io, dataProtocolId, stringLength)
	case "SzlId":
		return model.SzlIdParse(io)
	case "AlarmMessageObjectAckType":
		return model.AlarmMessageObjectAckTypeParse(io)
	case "AlarmMessageAckPushType":
		return model.AlarmMessageAckPushTypeParse(io)
	case "S7Message":
		return model.S7MessageParse(io)
	case "S7VarPayloadStatusItem":
		return model.S7VarPayloadStatusItemParse(io)
	case "S7Parameter":
		messageType, err := utils.StrToUint8(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.S7ParameterParse(io, messageType)
	case "S7DataAlarmMessage":
		cpuFunctionType, err := utils.StrToUint8(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.S7DataAlarmMessageParse(io, cpuFunctionType)
	case "SzlDataTreeItem":
		return model.SzlDataTreeItemParse(io)
	case "COTPPacket":
		cotpLen, err := utils.StrToUint16(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.COTPPacketParse(io, cotpLen)
	case "S7PayloadUserDataItem":
		cpuFunctionType, err := utils.StrToUint8(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		cpuSubfunction, err := utils.StrToUint8(arguments[1])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.S7PayloadUserDataItemParse(io, cpuFunctionType, cpuSubfunction)
	case "DateAndTime":
		return model.DateAndTimeParse(io)
	case "COTPParameter":
		rest, err := utils.StrToUint8(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		return model.COTPParameterParse(io, rest)
	case "AlarmMessageObjectPushType":
		return model.AlarmMessageObjectPushTypeParse(io)
	case "State":
		return model.StateParse(io)
	case "AlarmMessagePushType":
		return model.AlarmMessagePushTypeParse(io)
	case "TPKTPacket":
		return model.TPKTPacketParse(io)
	case "AlarmMessageAckType":
		return model.AlarmMessageAckTypeParse(io)
	case "AssociatedValueType":
		return model.AssociatedValueTypeParse(io)
	case "AlarmMessageAckObjectPushType":
		return model.AlarmMessageAckObjectPushTypeParse(io)
	case "S7Payload":
		messageType, err := utils.StrToUint8(arguments[0])
		if err != nil {
			return nil, errors.Wrap(err, "Error parsing")
		}
		var parameter model.S7Parameter
		return model.S7PayloadParse(io, messageType, &parameter)
	case "S7VarRequestParameterItem":
		return model.S7VarRequestParameterItemParse(io)
	case "S7VarPayloadDataItem":
		return model.S7VarPayloadDataItemParse(io)
	case "AlarmMessageQueryType":
		return model.AlarmMessageQueryTypeParse(io)
	case "AlarmMessageAckResponseType":
		return model.AlarmMessageAckResponseTypeParse(io)
	case "AlarmMessageObjectQueryType":
		return model.AlarmMessageObjectQueryTypeParse(io)
	case "S7Address":
		return model.S7AddressParse(io)
	case "S7ParameterUserDataItem":
		return model.S7ParameterUserDataItemParse(io)
	}
	return nil, errors.Errorf("Unsupported type %s", typeName)
}
