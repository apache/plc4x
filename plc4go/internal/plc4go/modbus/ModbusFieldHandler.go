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
package modbus

import (
	"errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"regexp"
)

type ModbusFieldType uint8

const (
	MODBUS_FIELD_COIL              ModbusFieldType = 0x00
	MODBUS_FIELD_DISCRETE_INPUT    ModbusFieldType = 0x01
	MODBUS_FIELD_INPUT_REGISTER    ModbusFieldType = 0x03
	MODBUS_FIELD_HOLDING_REGISTER  ModbusFieldType = 0x04
	MODBUS_FIELD_EXTENDED_REGISTER ModbusFieldType = 0x06
)

func (m ModbusFieldType) GetName() string {
	switch m {
	case MODBUS_FIELD_COIL:
		return "ModbusFieldHoldingRegister"
	case MODBUS_FIELD_DISCRETE_INPUT:
		return "ModbusFieldDiscreteInput"
	case MODBUS_FIELD_INPUT_REGISTER:
		return "ModbusFieldInputRegister"
	case MODBUS_FIELD_HOLDING_REGISTER:
		return "ModbusFieldHoldingRegister"
	case MODBUS_FIELD_EXTENDED_REGISTER:
		return "ModbusFieldExtendedRegister"
	}
	return ""
}

type FieldHandler struct {
	plc4xCoilPattern               *regexp.Regexp
	numericCoilPattern             *regexp.Regexp
	plc4xDiscreteInputPattern      *regexp.Regexp
	numericDiscreteInputPattern    *regexp.Regexp
	plc4xInputRegisterPattern      *regexp.Regexp
	numericInputRegisterPattern    *regexp.Regexp
	plc4xHoldingRegisterPattern    *regexp.Regexp
	numericHoldingRegisterPattern  *regexp.Regexp
	plc4xExtendedRegisterPattern   *regexp.Regexp
	numericExtendedRegisterPattern *regexp.Regexp
	spi.PlcFieldHandler
}

func NewFieldHandler() FieldHandler {
	generalAddressPattern := `(?P<address>\d+)(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)])?$`
	generalFixedDigitAddressPattern := `(?P<address>\d{4,5})?(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)])?$`
	return FieldHandler{
		plc4xCoilPattern:               regexp.MustCompile("^coil:" + generalAddressPattern),
		numericCoilPattern:             regexp.MustCompile("^0[xX]?" + generalFixedDigitAddressPattern),
		plc4xDiscreteInputPattern:      regexp.MustCompile("^discrete-input:" + generalAddressPattern),
		numericDiscreteInputPattern:    regexp.MustCompile("^1[xX]?" + generalFixedDigitAddressPattern),
		plc4xInputRegisterPattern:      regexp.MustCompile("^input-register:" + generalAddressPattern),
		numericInputRegisterPattern:    regexp.MustCompile("^3[xX]?" + generalFixedDigitAddressPattern),
		plc4xHoldingRegisterPattern:    regexp.MustCompile("^holding-register:" + generalAddressPattern),
		numericHoldingRegisterPattern:  regexp.MustCompile("^4[xX]?" + generalFixedDigitAddressPattern),
		plc4xExtendedRegisterPattern:   regexp.MustCompile("^extended-register:" + generalAddressPattern),
		numericExtendedRegisterPattern: regexp.MustCompile("^6[xX]?" + generalFixedDigitAddressPattern),
	}
}

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgropMatches(m.plc4xCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.numericCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.plc4xDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_DISCRETE_INPUT, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.numericDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_DISCRETE_INPUT, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.plc4xInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_INPUT_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.numericInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_INPUT_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.plc4xHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_HOLDING_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.numericHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_HOLDING_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.plc4xExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_EXTENDED_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	} else if match := utils.GetSubgropMatches(m.numericExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_EXTENDED_REGISTER, match["address"], match["quantity"], "IEC61131_"+match["datatype"])
	}
	return nil, errors.New("Invalid address format for address '" + query + "'")
}
