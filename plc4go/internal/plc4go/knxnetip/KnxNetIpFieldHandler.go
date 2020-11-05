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
package knxnetip

import (
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "plc4x.apache.org/plc4go-modbus-driver/v0/pkg/plc4go/model"
    "regexp"
)

type KnxNetIpFieldType uint8

const (
    KNX_NET_IP_FIELD_COIL KnxNetIpFieldType = 0x00
)

func (m KnxNetIpFieldType) GetName() string {
    switch m {
    case KNX_NET_IP_FIELD_COIL:
        return "ModbusFieldHoldingRegister"
    }
    return ""
}

type FieldHandler struct {
    knxNetIpGroupAddress3Level *regexp.Regexp
    knxNetIpGroupAddress2Level *regexp.Regexp
    knxNetIpGroupAddress1Level *regexp.Regexp
    spi.PlcFieldHandler
}

func NewFieldHandler() FieldHandler {
    return FieldHandler{
        knxNetIpGroupAddress3Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*))/(?P<middleGroup>(\d|\*))/(?P<subGroup>(\d{1,3}|\*))`),
        knxNetIpGroupAddress2Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*))/(?P<subGroup>(\d{1,4}|\*))`),
        knxNetIpGroupAddress1Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,5}|\*))`),
    }
}

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
    if match := utils.GetSubgropMatches(m.knxNetIpGroupAddress3Level, query); match != nil {
        return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
    } else if match := utils.GetSubgropMatches(m.knxNetIpGroupAddress2Level, query); match != nil {
        return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
    } else if match := utils.GetSubgropMatches(m.knxNetIpGroupAddress1Level, query); match != nil {
        return NewModbusPlcFieldFromStrings(MODBUS_FIELD_DISCRETE_INPUT, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
    }
    return nil, errors.New("Invalid address format for address '" + query + "'")
}
