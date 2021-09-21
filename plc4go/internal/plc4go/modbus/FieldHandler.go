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

package modbus

import (
	"fmt"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"regexp"
)

type FieldType uint8

//go:generate stringer -type FieldType
const (
	Coil             FieldType = 0x00
	DiscreteInput    FieldType = 0x01
	InputRegister    FieldType = 0x03
	HoldingRegister  FieldType = 0x04
	ExtendedRegister FieldType = 0x06
)

func (i FieldType) GetName() string {
	return fmt.Sprintf("ModbusField%s", i.String())
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
	if match := utils.GetSubgroupMatches(m.plc4xCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(Coil, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.numericCoilPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(Coil, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.plc4xDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(DiscreteInput, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.numericDiscreteInputPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(DiscreteInput, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.plc4xInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(InputRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.numericInputRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(InputRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.plc4xHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(HoldingRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.numericHoldingRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(HoldingRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.plc4xExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(ExtendedRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	} else if match := utils.GetSubgroupMatches(m.numericExtendedRegisterPattern, query); match != nil {
		return NewModbusPlcFieldFromStrings(ExtendedRegister, match["address"], match["quantity"], model2.ModbusDataTypeByName(match["datatype"]))
	}
	return nil, errors.Errorf("Invalid address format for address '%s'", query)
}
