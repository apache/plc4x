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

package abeth

import (
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// TagHandler parses ab-eth tag addresses of the form
//
//	N<fileNumber>:<elementNumber>[/<bitNumber>]:<fileType>[[<byteSize>]]
//
// e.g. N7:3:INTEGER[2], N7:3:WORD or N7:3/5:SINGLEBIT. Ported from plc4j's
// AbEthTag.ADDRESS_PATTERN.
type TagHandler struct {
	addressPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		addressPattern: regexp.MustCompile(`^N(?P<fileNumber>\d{1,7}):(?P<elementNumber>\d{1,7})(/(?P<bitNumber>\d{1,7}))?:(?P<dataType>[a-zA-Z_]+)(\[(?P<size>\d+)])?$`),
	}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	match := utils.GetSubgroupMatches(m.addressPattern, tagAddress)
	if match == nil {
		return nil, errors.Errorf("Unable to parse %s", tagAddress)
	}

	fileNumber, err := parseUint8(match["fileNumber"], "file number")
	if err != nil {
		return nil, err
	}
	elementNumber, err := parseUint8(match["elementNumber"], "element number")
	if err != nil {
		return nil, err
	}
	// The bit number is optional; plc4j defaults it to zero, which also means "this is not a bit
	// tag" as far as GetAddressString is concerned.
	bitNumber := uint8(0)
	if bitNumberString := match["bitNumber"]; bitNumberString != "" {
		if bitNumber, err = parseUint8(bitNumberString, "bit number"); err != nil {
			return nil, err
		}
	}

	fileType, ok := FileTypeByName(strings.ToUpper(match["dataType"]))
	if !ok {
		return nil, errors.Errorf("Unknown file type %s in %s", match["dataType"], tagAddress)
	}

	// The synthetic file types have a fixed width; every real file type takes its width from the
	// [<byteSize>] suffix, which is mandatory for them (plc4j throws a NumberFormatException when
	// it is missing, we say what is wrong).
	var byteSize uint8
	switch fileType {
	case FileTypeWord, FileTypeSinglebit:
		byteSize = 2
		if err := checkFixedWidthSize(fileType, byteSize, match["size"], tagAddress); err != nil {
			return nil, err
		}
	case FileTypeDword:
		byteSize = 4
		if err := checkFixedWidthSize(fileType, byteSize, match["size"], tagAddress); err != nil {
			return nil, err
		}
	default:
		sizeString := match["size"]
		if sizeString == "" {
			return nil, errors.Errorf("A %s tag needs an explicit byte size, e.g. %s[2]", fileType, tagAddress)
		}
		if byteSize, err = parseUint8(sizeString, "byte size"); err != nil {
			return nil, err
		}
		if byteSize == 0 {
			return nil, errors.New("byte size must be greater than zero")
		}
	}

	if fileType == FileTypeSinglebit && bitNumber > 15 {
		// A SINGLEBIT tag reads two bytes and picks one of their 16 bits.
		return nil, errors.Errorf("bit number %d is out of range, a SINGLEBIT tag covers bits 0 to 15", bitNumber)
	}

	return NewTag(byteSize, fileNumber, fileType, elementNumber, bitNumber), nil
}

// checkFixedWidthSize accepts a [<byteSize>] suffix on a synthetic file type only when it agrees
// with the width the type already fixes. A redundant suffix is harmless - GetAddressString never
// writes one, so the round trip is unaffected - but a suffix which contradicts the type is a
// mistake worth reporting: plc4j hard-codes the width and drops the suffix on the floor, which made
// N7:3:WORD[9] parse to a two byte tag without a word of complaint.
func checkFixedWidthSize(fileType FileType, byteSize uint8, sizeString string, tagAddress string) error {
	if sizeString == "" {
		return nil
	}
	declaredSize, err := parseUint8(sizeString, "byte size")
	if err != nil {
		return err
	}
	if declaredSize != byteSize {
		return errors.Errorf("a %s tag is always %d bytes wide, so the size %d in %s can't be honored",
			fileType, byteSize, declaredSize, tagAddress)
	}
	return nil
}

// ParseQuery is not supported: neither this driver nor plc4j's browses an ab-eth PLC.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("This driver doesn't support browsing")
}

// parseUint8 keeps the address numbers inside what the wire format can carry: every one of them is
// a single byte in the "protected typed logical read" command. plc4j parses them as a short and
// silently truncates on the wire.
func parseUint8(value string, what string) (uint8, error) {
	parsed, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, errors.Wrapf(err, "Error parsing %s %s (has to fit into a single byte)", what, value)
	}
	return uint8(parsed), nil
}
