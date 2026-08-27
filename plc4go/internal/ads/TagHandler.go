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

package ads

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TagHandler struct {
	directAdsStringTag *regexp.Regexp
	directAdsTag       *regexp.Regexp
	symbolicAdsTag     *regexp.Regexp
	driverContext      *DriverContext
}

// NewTagHandler this constructor creates a version of the TagHandler that's detached from a connection and can't provide context-sensitive feedback.
func NewTagHandler() TagHandler {
	return TagHandler{
		directAdsStringTag: regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+))` + spiModel.ArrayGroupPattern + `:(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)$`),
		directAdsTag:       regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+))` + spiModel.ArrayGroupPattern + `:(?P<adsDataType>\w+)$`),
		symbolicAdsTag:     regexp.MustCompile(`^(?P<symbolicAddress>[^\[]+)` + spiModel.ArrayGroupPattern + `$`),
	}
}

// NewTagHandlerWithDriverContext this constructor creates a version of the TagHandler that is connected to a connection and can provide context-sensitive feedback.
func NewTagHandlerWithDriverContext(driverContext *DriverContext) TagHandler {
	return TagHandler{
		directAdsStringTag: regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+))` + spiModel.ArrayGroupPattern + `:(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)$`),
		directAdsTag:       regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+))` + spiModel.ArrayGroupPattern + `:(?P<adsDataType>\w+)$`),
		symbolicAdsTag:     regexp.MustCompile(`^(?P<symbolicAddress>[^\[]+)` + spiModel.ArrayGroupPattern + `$`),
		driverContext:      driverContext,
	}
}

func (m TagHandler) ParseTag(query string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.directAdsStringTag, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			// ParseUint instead of hex.DecodeString: addresses like 0x8 have an odd
			// number of digits, which byte-wise decoding rejects.
			parsed, err := strconv.ParseUint(indexGroupHexString, 16, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding indexGroup")
			}
			indexGroup = uint32(parsed)
		} else {
			parsedIndexGroup, err := strconv.ParseUint(match["indexGroup"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexGroup = uint32(parsedIndexGroup)
		}
		var indexOffset uint32
		if indexOffsetHexString := match["indexOffsetHex"]; indexOffsetHexString != "" {
			// ParseUint instead of hex.DecodeString: addresses like 0x8 have an odd
			// number of digits, which byte-wise decoding rejects.
			parsed, err := strconv.ParseUint(indexOffsetHexString, 16, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding indexOffset")
			}
			indexOffset = uint32(parsed)
		} else {
			parsedIndexOffset, err := strconv.ParseUint(match["indexOffset"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = uint32(parsedIndexOffset)
		}
		adsDataTypeName := match["adsDataType"]
		if adsDataTypeName == "" {
			return nil, errors.Errorf("Missing ads data type")
		}
		plcValueType, ok := apiValues.PlcValueTypeByName(adsDataTypeName)
		if !ok {
			return nil, fmt.Errorf("invalid ads data type")
		}

		stringLength := model.NONE
		var arrayInfo []apiModel.ArrayInfo

		tmpStringLength, err := strconv.ParseInt(match["stringLength"], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Error decoding string length")
		}
		stringLength = int32(tmpStringLength)

		// The selection goes through the shared parser, so plc4go accepts exactly what plc4j
		// accepts. A direct address names a memory location, so it carries one dimension.
		arrayInfo, err = spiModel.ParseArrayExpression(match["array"], query, spiModel.SingleDimension)
		if err != nil {
			return nil, err
		}

		return model.NewDirectAdsPlcTag(indexGroup, indexOffset, plcValueType, stringLength, arrayInfo)
	} else if match := utils.GetSubgroupMatches(m.directAdsTag, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			// ParseUint instead of hex.DecodeString: addresses like 0x8 have an odd
			// number of digits, which byte-wise decoding rejects.
			parsed, err := strconv.ParseUint(indexGroupHexString, 16, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding indexGroup")
			}
			indexGroup = uint32(parsed)
		} else {
			parsedIndexGroup, err := strconv.ParseUint(match["indexGroup"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexGroup = uint32(parsedIndexGroup)
		}
		var indexOffset uint32
		if indexOffsetHexString := match["indexOffsetHex"]; indexOffsetHexString != "" {
			// ParseUint instead of hex.DecodeString: addresses like 0x8 have an odd
			// number of digits, which byte-wise decoding rejects.
			parsed, err := strconv.ParseUint(indexOffsetHexString, 16, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding indexOffset")
			}
			indexOffset = uint32(parsed)
		} else {
			parsedIndexOffset, err := strconv.ParseUint(match["indexOffset"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = uint32(parsedIndexOffset)
		}
		adsDataTypeName := match["adsDataType"]
		if adsDataTypeName == "" {
			return nil, errors.Errorf("Missing ads data type")
		}
		plcValueType, ok := apiValues.PlcValueTypeByName(adsDataTypeName)
		if !ok {
			return nil, fmt.Errorf("invalid ads data type")
		}

		// The selection goes through the shared parser, so plc4go accepts exactly what plc4j
		// accepts. A direct address names a memory location, so it carries one dimension.
		arrayInfo, err := spiModel.ParseArrayExpression(match["array"], query, spiModel.SingleDimension)
		if err != nil {
			return nil, err
		}

		return model.NewDirectAdsPlcTag(indexGroup, indexOffset, plcValueType, model.NONE, arrayInfo)
	} else if match := utils.GetSubgroupMatches(m.symbolicAdsTag, query); match != nil {
		// A symbolic address is anything that is not a direct one, so an address that looks
		// direct but does not parse would otherwise be accepted here as a symbol name of its
		// own - "1234/5678:BOOL[42]", written before the notation moved, would silently become
		// a symbol lookup rather than an error. Report it instead, naming the address to write.
		if looksLikeADirectAddress(query) {
			return nil, spiModel.InvalidAddressError(query,
				"{indexGroup}/{indexOffset}[selection]:{TYPE} - for example 0x4020/0[0..3]:DINT")
		}

		// The selection goes through the shared parser, so plc4go accepts exactly what plc4j
		// accepts. Only the last dimension may span more than one element: a range before it
		// would ask for a member of several elements at once, which is not one contiguous read.
		arrayInfo, err := spiModel.ParseArrayExpression(match["array"], query,
			spiModel.Unconstrained.WithOnlyTrailingDimensionMayBeRange(true))
		if err != nil {
			return nil, err
		}

		return model.NewAdsSymbolicPlcTag(match["symbolicAddress"], arrayInfo)
	} else {
		// The start-and-count form was a plc4go extension with no counterpart in plc4j, so it
		// is gone. Name the range that selects the same elements rather than reporting only
		// that nothing matched.
		if match := startAndCountForm.FindStringSubmatch(query); match != nil {
			start, _ := strconv.Atoi(match[2])
			count, _ := strconv.Atoi(match[3])
			if count > 0 {
				return nil, errors.Errorf("invalid address '%s': the start-and-count form '[a:b]' "+
					"is no longer supported, so this address is now written '%s[%d..%d]'",
					query, match[1], start, start+count-1)
			}
		}
		return nil, spiModel.InvalidAddressError(query,
			"{symbol}[selection] or {indexGroup}/{indexOffset}[selection]:{TYPE} - "+
				"for example MAIN.g_arr[0..3]")
	}
}

// looksLikeADirectAddress reports whether an address has the shape of a direct one - an index
// group and offset separated by a slash - regardless of whether it parses. A symbolic address
// never contains a slash, so this cannot mistake one for the other.
// startAndCountForm matches the removed "[a:b]" spelling - a start and a count - so a rejection
// can name the range that selects the same elements.
var startAndCountForm = regexp.MustCompile(`^(.*)\[(\d+):(\d+)]$`)

var directAddressShape = regexp.MustCompile(`^(?:0[xX][0-9a-fA-F]+|\d+)/(?:0[xX][0-9a-fA-F]+|\d+):`)

func looksLikeADirectAddress(query string) bool {
	return directAddressShape.MatchString(query)
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return symbolicPlcQuery{
		query: query,
	}, nil
}

func (m TagHandler) getUint32Value(stringValue string) (uint32, error) {
	intValue, err := strconv.ParseUint(stringValue, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid number format parsing '%s' as int32: %v", stringValue, err)
	} else {
		return uint32(intValue), nil
	}
}
