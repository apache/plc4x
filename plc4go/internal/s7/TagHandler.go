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

package s7

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
)

type TagType uint8

//go:generate stringer -type TagType
//go:generate go run ../../tools/plc4xlicenser/gen.go -type=TagType
const (
	S7Tag       TagType = 0x00
	S7StringTag TagType = 0x01
)

func (i TagType) GetName() string {
	return i.String()
}

type TagHandler struct {
	addressPattern                *regexp.Regexp
	dataBlockAddressPattern       *regexp.Regexp
	dataBlockShortPattern         *regexp.Regexp
	dataBlockStringAddressPattern *regexp.Regexp
	dataBlockStringShortPattern   *regexp.Regexp
	plcProxyAddressPattern        *regexp.Regexp

	passLogToModel bool
	log            zerolog.Logger
}

func NewTagHandler(_options ...options.WithOption) TagHandler {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger, _ := options.ExtractCustomLogger(_options...)
	return TagHandler{
		addressPattern: regexp.MustCompile(`^%(?P<memoryArea>.)(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		//blockNumber usually has its max hat around 64000 --> 5digits
		dataBlockAddressPattern:       regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}).DB(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		dataBlockShortPattern:         regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}):(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?`),
		dataBlockStringAddressPattern: regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}).DB(?P<transferSizeCode>[XBWD]?)(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(\[(?P<numElements>\d+)])?`),
		dataBlockStringShortPattern:   regexp.MustCompile(`^%DB(?P<blockNumber>\d{1,5}):(?P<byteOffset>\d{1,7})(.(?P<bitOffset>[0-7]))?:(?P<dataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(\[(?P<numElements>\d+)])?`),
		plcProxyAddressPattern:        regexp.MustCompile(`[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}-[0-9A-F]{2}`),

		passLogToModel: passLoggerToModel,
		log:            customLogger,
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

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.dataBlockStringAddressPattern, tagAddress); match != nil {
		dataType, ok := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		parsedStringLength, err := strconv.ParseUint(match[STRING_LENGTH], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting stringlength")
		}
		stringLength := uint16(parsedStringLength)
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		parsedByteOffset, err := strconv.ParseUint(match[BYTE_OFFSET], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(parsedByteOffset)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			parsedBitOffset, err := strconv.ParseUint(match[BIT_OFFSET], 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(parsedBitOffset)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			parsedNumElements, err := strconv.ParseUint(match[NUM_ELEMENTS], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(parsedNumElements)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}

		return NewStringTag(memoryArea, 0, byteOffset, bitOffset, numElements, stringLength, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockStringShortPattern, tagAddress); match != nil {
		dataType, ok := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		parsedStringLength, err := strconv.ParseUint(match[STRING_LENGTH], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting stringlength")
		}
		stringLength := uint16(parsedStringLength)
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		parsedBlockNumber, err := strconv.ParseUint(match[BLOCK_NUMBER], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(parsedBlockNumber)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		parsedByteOffset, err := strconv.ParseUint(match[BYTE_OFFSET], 10, 8)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(parsedByteOffset)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			parsedNumElements, err := strconv.ParseUint(match[NUM_ELEMENTS], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(parsedNumElements)
		}

		return NewStringTag(memoryArea, blockNumber, byteOffset, bitOffset, numElements, stringLength, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockAddressPattern, tagAddress); match != nil {
		dataType, ok := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		parsedBlockNumber, err := strconv.ParseUint(match[BLOCK_NUMBER], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(parsedBlockNumber)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		parsedByteOffset, err := strconv.ParseUint(match[BYTE_OFFSET], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(parsedByteOffset)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			parsedBitOffset, err := strconv.ParseUint(match[BIT_OFFSET], 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(parsedBitOffset)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			parsedNumElements, err := strconv.ParseUint(match[NUM_ELEMENTS], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(parsedNumElements)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}

		return NewTag(memoryArea, blockNumber, byteOffset, bitOffset, numElements, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.dataBlockShortPattern, tagAddress); match != nil {
		dataType, ok := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		memoryArea := readWriteModel.MemoryArea_DATA_BLOCKS
		parsedBlockNumber, err := strconv.ParseUint(match[BLOCK_NUMBER], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting blocknumber")
		}
		blockNumber, err := checkDatablockNumber(parsedBlockNumber)
		if err != nil {
			return nil, errors.Wrap(err, "Error checking blocknumber")
		}
		parsedByteOffset, err := strconv.ParseUint(match[BYTE_OFFSET], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(parsedByteOffset)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			parsedBitOffset, err := strconv.ParseUint(match[BIT_OFFSET], 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(parsedBitOffset)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			parsedNumElements, err := strconv.ParseUint(match[NUM_ELEMENTS], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(parsedNumElements)
		}

		return NewTag(memoryArea, blockNumber, byteOffset, bitOffset, numElements, dataType), nil
	} else if match := utils.GetSubgroupMatches(m.plcProxyAddressPattern, tagAddress); match != nil {
		addressData, err := hex.DecodeString(strings.ReplaceAll(tagAddress, "[-]", ""))
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse address: %s", tagAddress)
		}
		ctxForModel := options.GetLoggerContextForModel(context.TODO(), m.log, options.WithPassLoggerToModel(m.passLogToModel))
		s7Address, err := readWriteModel.S7AddressAnyParse(ctxForModel, addressData)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to parse address: %s", tagAddress)
		}
		s7AddressAny := s7Address.(readWriteModel.S7AddressAny)
		if (s7AddressAny.GetTransportSize() != readWriteModel.TransportSize_BOOL) && s7AddressAny.GetBitAddress() != 0 {
			return nil, errors.New("A bit offset other than 0 is only supported for type BOOL")
		}

		return NewTag(
			s7AddressAny.GetArea(),
			s7AddressAny.GetDbNumber(),
			s7AddressAny.GetByteAddress(),
			s7AddressAny.GetBitAddress(),
			s7AddressAny.GetNumberOfElements(),
			s7AddressAny.GetTransportSize(),
		), nil
	} else if match := utils.GetSubgroupMatches(m.addressPattern, tagAddress); match != nil {
		dataType, ok := readWriteModel.TransportSizeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		memoryArea, err := getMemoryAreaForShortName(match[MEMORY_AREA])
		if err != nil {
			return nil, errors.Wrap(err, "Error getting memory area")
		}
		transferSizeCode := getSizeCode(match[TRANSFER_SIZE_CODE])
		parsedTransferSizeCode, err := strconv.ParseUint(match[BYTE_OFFSET], 10, 16)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		byteOffset, err := checkByteOffset(parsedTransferSizeCode)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting byteoffset")
		}
		bitOffset := uint8(0)
		if match[BIT_OFFSET] != "" {
			parsedBitOffset, err := strconv.ParseUint(match[BIT_OFFSET], 10, 8)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting byteoffset")
			}
			bitOffset = uint8(parsedBitOffset)
		} else if dataType == readWriteModel.TransportSize_BOOL {
			return nil, errors.New("Expected bit offset for BOOL parameters.")
		}
		numElements := uint16(1)
		if match[NUM_ELEMENTS] != "" {
			parsedNumElements, err := strconv.ParseUint(match[NUM_ELEMENTS], 10, 16)
			if err != nil {
				return nil, errors.Wrap(err, "Error converting numelements")
			}
			numElements = uint16(parsedNumElements)
		}

		if (transferSizeCode != 0) && (dataType.ShortName() != transferSizeCode) {
			return nil, errors.Errorf("Transfer size code '%d' doesn't match specified data type '%s'", transferSizeCode, dataType)
		}
		if (dataType != readWriteModel.TransportSize_BOOL) && bitOffset != 0 {
			return nil, errors.New("A bit offset other than 0 is only supported for type BOOL")
		}

		return NewTag(memoryArea, 0, byteOffset, bitOffset, numElements, dataType), nil
	}
	return nil, errors.Errorf("Unable to parse %s", tagAddress)
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}

func checkDatablockNumber(blockNumber uint64) (uint16, error) {
	//ToDo check the value or add reference - limit eventually depending on active S7 --> make a case selection
	if blockNumber > 64000 || blockNumber < 1 {
		return 0, errors.New("Datablock numbers larger than 64000 or smaller than 1 are not supported.")
	}
	return uint16(blockNumber), nil
}

func checkByteOffset(byteOffset uint64) (uint16, error) {
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
	parsedSizeCode, err := strconv.ParseUint(string(chars[0]), 10, 8)
	if err != nil {
		return 0
	}
	return uint8(parsedSizeCode)
}

func getMemoryAreaForShortName(shortName string) (readWriteModel.MemoryArea, error) {
	for _, memoryArea := range readWriteModel.MemoryAreaValues {
		if memoryArea.ShortName() == shortName {
			return memoryArea, nil
		}
	}
	return 0, errors.Errorf("Unknown memory area for short name: '%s'", shortName)
}
