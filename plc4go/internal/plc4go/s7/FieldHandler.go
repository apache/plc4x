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

package s7

import (
	"encoding/hex"
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

type FieldType uint8

//go:generate stringer -type FieldType
const (
	S7Field       FieldType = 0x00
	S7StringField FieldType = 0x01
)

func (i FieldType) GetName() string {
	return i.String()
}

type FieldHandler struct {
	addressPattern                *regexp.Regexp
	dataBlockAddressPattern       *regexp.Regexp
	dataBlockShortPattern         *regexp.Regexp
	dataBlockStringAddressPattern *regexp.Regexp
	dataBlockStringShortPattern   *regexp.Regexp
	plcProxyAddressPattern        *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		addressPattern: regexp.MustCompile(`^%(?P<memoryArea>.)(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		//blockNumber usually has its max hat around 64000 --> 5digits
		dataBlockAddressPattern:       regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}).DB(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		dataBlockShortPattern:         regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}):(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		dataBlockStringAddressPattern: regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}).DB(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(\[(?P<numElements>\d+)])?`),
		dataBlockStringShortPattern:   regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}):(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(\[(?P<numElements>\d+)])?`),
		plcProxyAddressPattern:        regexp.MustCompile(`[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}`),
	}
}

const (
	DATA_TYPE          = "dataType"
	STRING_LENGTH      = "stringLength"
	TRANSFER_SIZE_CODE = "transferSizeCode"
	BLOCK_NUMBER       = "blockNumber"
	BYTE_OFFSET        = "byteOffset"
	BIT_OFFSET         = "bitOffset"
	NUM_ELEMENTS       = "numElements"
	MEMORY_AREA        = "memoryArea"
)

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.dataBlockStringAddressPattern, query); match != nil {
		dataType := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		atoi, err := strconv.Atoi(match[STRING_LENGTH])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting stringlength")
		}
		stringLength := uint16(atoi)
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		atoi, err = strconv.Atoi(match[BYTE_OFFSET])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			atoi, err := strconv.Atoi(match[BIT_OFFSET])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(atoi)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			atoi, err := strconv.Atoi(match[NUM_ELEMENTS])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(atoi)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}

		return NewStringField(memoryArea, 0, byteOffset, bitOffset, numElements, stringLength, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockStringShortPattern, query); match != nil {
		dataType := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		atoi, err := strconv.Atoi(match[STRING_LENGTH])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting stringlength")
		}
		stringLength := uint16(atoi)
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		atoi, err = strconv.Atoi(match[BLOCK_NUMBER])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		atoi, err = strconv.Atoi(match[BYTE_OFFSET])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			atoi, err := strconv.Atoi(match[NUM_ELEMENTS])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(atoi)
		}

		return NewStringField(memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockAddressPattern, query); match != nil {
		dataType := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		atoi, err := strconv.Atoi(match[BLOCK_NUMBER])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		atoi, err = strconv.Atoi(match[BYTE_OFFSET])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			atoi, err := strconv.Atoi(match[BIT_OFFSET])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(atoi)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			atoi, err := strconv.Atoi(match[NUM_ELEMENTS])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(atoi)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}

		return NewField(memoryArea, blockNumber, byteOffset, bitOffset, numElements, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockShortPattern, query); match != nil {
		dataType := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		atoi, err := strconv.Atoi(match[BLOCK_NUMBER])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		atoi, err = strconv.Atoi(match[BYTE_OFFSET])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			atoi, err := strconv.Atoi(match[BIT_OFFSET])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(atoi)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			atoi, err := strconv.Atoi(match[NUM_ELEMENTS])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(atoi)
		}

		return NewField(memoryArea, blockNumber, byteOffset, bitOffset, numElements, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.plcProxyAddressPattern, query); match != nil {
		addressData, err := hex.DecodeString(strings.ReplaceAll(query, "[-]", ""))
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse address: %s", query)
		}
		rb := utils.NewReadBufferByteBased(addressData)
		s7Address, err := readWriteModel.S7AddressAnyParse(rb)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse address: %s", query)
		}
		s7AddressAny := s7Address.Child.(*readWriteModel.S7AddressAny)
		if (s7AddressAny.TransportSize != readWriteModel.TransportSize_BOOL) && s7AddressAny.BitAddress != 0 {
			return nil, errors.New("A bit offset other than 0 is only supported for type BOOL")
		}

		return NewField(
			s7AddressAny.Area,
			s7AddressAny.DbNumber,
			s7AddressAny.ByteAddress,
			s7AddressAny.BitAddress,
			s7AddressAny.NumberOfElements,
			s7AddressAny.TransportSize,
		), nil
	} else if match := utils.GetSubgroupMatches(m.addressPattern, query); match != nil {
		dataType := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		memoryArea, err := getMemoryAreaForShortName(match[MEMORY_AREA])
		if err != nil {
			return nil, errors.Wrap(err, "Error getting memory area")
		}
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		atoi, err := strconv.Atoi(match[BYTE_OFFSET])
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(atoi)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			atoi, err := strconv.Atoi(match[BIT_OFFSET])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(atoi)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			atoi, err := strconv.Atoi(match[NUM_ELEMENTS])
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(atoi)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}
		if (dataType != readWriteModel.TransportSize_BOOL) && bitOffset != 0 {
			errors.New("A bit offset other than 0 is only supported for type BOOL")
		}

		return NewField(memoryArea, 0, byteOffset, bitOffset, numElements, dataType), nil
	}
	return nil, errors.Errorf("Unable to parse %s", query)
}

func checkDatablockNumber(blockNumber int) (uint16, error) {
	//ToDo check the value or add reference - limit eventually depending on active S7 --> make a case selection
	if blockNumber > 64000 || blockNumber < 1 {
		return 0, errors.New("Datablock numbers larger than 64000 or smaller than 1 are not supported.")
	}
	return uint16(blockNumber), nil
}

func checkByteOffset(byteOffset int) (uint16, error) {
	//ToDo check the value or add reference
	if byteOffset > 2097151 || byteOffset < 0 {
		return 0, errors.New("ByteOffset must be smaller than 2097151 and positive.")
	}
	return uint16(byteOffset), nil
}

func getSizeCode(value string) uint8 {
	if value == "" {
		return 0
	}
	if len(value) > 1 {
		return 0
	}
	chars := []rune(value)
	atoi, err := strconv.Atoi(string(chars[0]))
	if err != nil {
		return 0
	}
	return uint8(atoi)
}

func getMemoryAreaForShortName(shortName string) (readWriteModel.MemoryArea, error) {
	for _, memoryArea := range readWriteModel.MemoryAreaValues {
		if memoryArea.ShortName() == shortName {
			return memoryArea, nil
		}
	}
	return 0, errors.Errorf("Unknown memory area for short name: '%s'", shortName)
}
