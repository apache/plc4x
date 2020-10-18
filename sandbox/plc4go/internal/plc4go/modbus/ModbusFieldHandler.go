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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/pkg/plc4go/model"
	"regexp"
)

const (
	MODBUS_FIELD_COIL              = uint8(0)
	MODBUS_FIELD_DISCRETE_INPUT    = uint8(1)
	MODBUS_FIELD_INPUT_REGISTER    = uint8(3)
	MODBUS_FIELD_HOLDING_REGISTER  = uint8(4)
	MODBUS_FIELD_EXTENDED_REGISTER = uint8(6)
)

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
	if match := GetSubgropMatches(m.plc4xCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.numericCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_COIL, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.plc4xDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_DISCRETE_INPUT, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.numericDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_DISCRETE_INPUT, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.plc4xInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_INPUT_REGISTER, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.numericInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_INPUT_REGISTER, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.plc4xHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_HOLDING_REGISTER, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.numericHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_HOLDING_REGISTER, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.plc4xExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_EXTENDED_REGISTER, match["address"], match["quantity"], match["datatype"])
	} else if match := GetSubgropMatches(m.numericExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(MODBUS_FIELD_EXTENDED_REGISTER, match["address"], match["quantity"], match["datatype"])
	}
	return nil, errors.New("Invalid address format for address '" + query + "'")
}

func GetSubgropMatches(r *regexp.Regexp, query string) map[string]string {
	match := r.FindStringSubmatch(query)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}
