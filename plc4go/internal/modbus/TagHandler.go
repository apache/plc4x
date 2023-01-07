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

package modbus

import (
	"fmt"
	"regexp"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	model2 "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type TagType uint8

//go:generate stringer -type TagType
const (
	Coil             TagType = 0x00
	DiscreteInput    TagType = 0x01
	InputRegister    TagType = 0x03
	HoldingRegister  TagType = 0x04
	ExtendedRegister TagType = 0x06
)

func (i TagType) GetName() string {
	return fmt.Sprintf("ModbusTags%s", i.String())
}

type TagHandler struct {
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

func NewTagHandler() TagHandler {
	generalAddressPattern := `(?P<address>\d+)\.?(?P<offset>\d+)?(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)])?$`
	generalFixedDigitAddressPattern := `(?P<address>\d{4,5})?\.?(?P<offset>\d+)?(:(?P<datatype>[a-zA-Z_]+))?(\[(?P<quantity>\d+)])?$`
	return TagHandler{
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

func (m TagHandler) ParseTag(tagAddress string) (model.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.plc4xCoilPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(Coil, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.numericCoilPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(Coil, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.plc4xDiscreteInputPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(DiscreteInput, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.numericDiscreteInputPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(DiscreteInput, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.plc4xInputRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(InputRegister, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.numericInputRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(InputRegister, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.plc4xHoldingRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(HoldingRegister, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.numericHoldingRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(HoldingRegister, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.plc4xExtendedRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(ExtendedRegister, match["address"], match["offset"], match["quantity"], typeByName)
	} else if match := utils.GetSubgroupMatches(m.numericExtendedRegisterPattern, tagAddress); match != nil {
		typeByName, ok := model2.ModbusDataTypeByName(match["datatype"])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match["datatype"])
		}
		return NewModbusPlcTagFromStrings(ExtendedRegister, match["address"], match["offset"], match["quantity"], typeByName)
	}
	return nil, errors.Errorf("Invalid address format for address '%s'", tagAddress)
}

func (m TagHandler) ParseQuery(query string) (model.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
