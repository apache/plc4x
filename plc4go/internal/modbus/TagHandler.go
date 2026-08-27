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
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TagType uint8

//go:generate go tool stringer -type TagType
//go:generate go tool plc4xLicencer -type=TagType
const (
	Coil             TagType = 0x00
	DiscreteInput    TagType = 0x01
	InputRegister    TagType = 0x03
	HoldingRegister  TagType = 0x04
	ExtendedRegister TagType = 0x06
)

func (i TagType) GetName() string {
	return fmt.Sprintf("ModbusTag%s", i.String())
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

	options []options.WithOption

	log zerolog.Logger
}

func NewTagHandler(_options ...options.WithOption) TagHandler {
	// The selection sits between the address and the type, as it does in plc4j's
	// ModbusTag.ADDRESS_PATTERN: "holding-register:1[0..2]:STRING(20)" is three 20-character
	// strings. STRING and WSTRING carry the length of one string in parentheses. Behind all of
	// that, curly braces may carry per-tag settings (plc4j TagConfigParser.TAG_CONFIG_PATTERN).
	generalAddressPattern := `(?P<address>\d+)` + spiModel.ArrayGroupPattern + `(:(?P<datatype>[a-zA-Z_]+)(\((?P<stringLength>\d+)\))?)?` + tagConfigPattern + `$`
	generalFixedDigitAddressPattern := `(?P<address>\d{4,5})?` + spiModel.ArrayGroupPattern + `(:(?P<datatype>[a-zA-Z_]+)(\((?P<stringLength>\d+)\))?)?` + tagConfigPattern + `$`
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
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
		options:                        _options,
		log:                            customLogger,
	}
}

// defaultDataTypeFor is the type an address without an explicit one is read and written as. Ported
// from plc4j, where the per-area factories default to BOOL for the bit areas
// (ModbusTagCoil.java:95, and likewise ModbusTagDiscreteInput) and to INT for the register areas
// (ModbusTagHoldingRegister.java:94, and likewise ModbusTagInputRegister/ExtendedRegister).
func defaultDataTypeFor(tagType TagType) readWriteModel.ModbusDataType {
	switch tagType {
	case Coil, DiscreteInput:
		return readWriteModel.ModbusDataType_BOOL
	default:
		return readWriteModel.ModbusDataType_INT
	}
}

func dataTypeFor(tagType TagType, dataTypeString string) (readWriteModel.ModbusDataType, error) {
	if dataTypeString == "" {
		return defaultDataTypeFor(tagType), nil
	}
	typeByName, ok := readWriteModel.ModbusDataTypeByName(dataTypeString)
	if !ok {
		return 0, errors.Errorf("Unknown type %s", dataTypeString)
	}
	return typeByName, nil
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	for _, candidate := range []struct {
		pattern *regexp.Regexp
		tagType TagType
	}{
		{m.plc4xCoilPattern, Coil},
		{m.numericCoilPattern, Coil},
		{m.plc4xDiscreteInputPattern, DiscreteInput},
		{m.numericDiscreteInputPattern, DiscreteInput},
		{m.plc4xInputRegisterPattern, InputRegister},
		{m.numericInputRegisterPattern, InputRegister},
		{m.plc4xHoldingRegisterPattern, HoldingRegister},
		{m.numericHoldingRegisterPattern, HoldingRegister},
		{m.plc4xExtendedRegisterPattern, ExtendedRegister},
		{m.numericExtendedRegisterPattern, ExtendedRegister},
	} {
		match := utils.GetSubgroupMatches(candidate.pattern, tagAddress)
		if match == nil {
			continue
		}
		dataType, err := dataTypeFor(candidate.tagType, match["datatype"])
		if err != nil {
			return nil, err
		}
		config, err := parseTagConfig(match["config"])
		if err != nil {
			return nil, err
		}
		return NewModbusPlcTagFromStrings(candidate.tagType, match["address"], match["array"], match["stringLength"], dataType, config, m.options...)
	}
	// "holding-register:1:INT[4]" - the count after the type - no longer parses, so name the
	// form to write rather than reporting only that nothing matched.
	return nil, spiModel.InvalidAddressError(tagAddress,
		"{area}:{address}[selection]:{TYPE} - for example holding-register:1[0..3]:INT")
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}

// tagConfigPattern is the optional trailing block of per-tag settings, written in curly braces
// behind the address. Ported from plc4j's TagConfigParser.TAG_CONFIG_PATTERN.
const tagConfigPattern = `(\{(?P<config>[^}]*)})?`

// tagConfigKeyValuePattern picks the individual settings out of that block. It follows JSON syntax
// with unquoted keys, as plc4j's TagConfigParser.KEY_VALUE_PATTERN does. A bare word is accepted on
// top of what plc4j allows, so that an unquoted enum name such as "byte-order: BIG_ENDIAN" is read
// rather than silently skipped.
var tagConfigKeyValuePattern = regexp.MustCompile(`(?P<parameter>[\w\-_]+):\s*(?P<value>-?\d+\.\d+|-?\d+|"[^"]*"|'[^']*'|true|false|[A-Za-z_][\w\-]*),?`)

// tagConfig is what a tag may say about itself on top of its address.
type tagConfig struct {
	unitId    *uint8
	byteOrder *ByteOrder
}

// parseTagConfig reads the settings out of the curly-brace block. An empty block, or none at all,
// simply leaves everything at the connection's defaults.
func parseTagConfig(config string) (tagConfig, error) {
	var parsed tagConfig
	if config == "" {
		return parsed, nil
	}
	for _, match := range tagConfigKeyValuePattern.FindAllStringSubmatch(config, -1) {
		parameter, value := match[1], unquoteTagConfigValue(match[2])
		switch parameter {
		case "unit-id":
			unitId, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return tagConfig{}, errors.Errorf("Couldn't parse unit-id '%s' into an int between 0 and 255", value)
			}
			asUint8 := uint8(unitId)
			parsed.unitId = &asUint8
		case "byte-order":
			byteOrder, ok := ByteOrderByName(value)
			if !ok {
				return tagConfig{}, errors.Errorf("Unknown byte-order %s", value)
			}
			parsed.byteOrder = &byteOrder
		default:
			// plc4j's TagConfigParser collects every key into a map and leaves it to each tag to
			// pick out the ones it knows, so an address carrying a setting this driver has no use
			// for still parses.
		}
	}
	return parsed, nil
}

// unquoteTagConfigValue strips the quotes a string value may be written with.
func unquoteTagConfigValue(value string) string {
	for _, quote := range []string{"'", `"`} {
		if len(value) >= 2 && strings.HasPrefix(value, quote) && strings.HasSuffix(value, quote) {
			return value[1 : len(value)-1]
		}
	}
	return value
}
