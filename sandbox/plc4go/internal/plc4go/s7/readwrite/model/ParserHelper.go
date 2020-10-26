//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

type S7ParserHelper struct {
}

func (m S7ParserHelper) Parse(typeName string, arguments []string, io *utils.ReadBuffer) (spi.Message, error) {
    switch typeName {
    case "SzlId":
        return SzlIdParse(io)
    case "S7Message":
        return S7MessageParse(io)
    case "S7VarPayloadStatusItem":
        return S7VarPayloadStatusItemParse(io)
    case "S7Parameter":
        messageType, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        return S7ParameterParse(io, messageType)
    case "SzlDataTreeItem":
        return SzlDataTreeItemParse(io)
    case "COTPPacket":
        cotpLen, err := utils.StrToUint16(arguments[0])
        if err != nil {
            return nil, err
        }
        return COTPPacketParse(io, cotpLen)
    case "S7PayloadUserDataItem":
        cpuFunctionType, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        return S7PayloadUserDataItemParse(io, cpuFunctionType)
    case "COTPParameter":
        rest, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        return COTPParameterParse(io, rest)
    case "TPKTPacket":
        return TPKTPacketParse(io)
    case "S7Payload":
        messageType, err := utils.StrToUint8(arguments[0])
        if err != nil {
            return nil, err
        }
        var parameter IS7Parameter
        return S7PayloadParse(io, messageType, parameter)
    case "S7VarRequestParameterItem":
        return S7VarRequestParameterItemParse(io)
    case "S7VarPayloadDataItem":
        lastItem, err := utils.StrToBool(arguments[0])
        if err != nil {
            return nil, err
        }
        return S7VarPayloadDataItemParse(io, lastItem)
    case "S7Address":
        return S7AddressParse(io)
    case "S7ParameterUserDataItem":
        return S7ParameterUserDataItemParse(io)
    }
    return nil, errors.New("Unsupported type " + typeName)
}
