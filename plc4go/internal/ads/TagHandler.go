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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	model2 "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type TagHandler struct {
	directAdsStringTag *regexp.Regexp
	directAdsTag       *regexp.Regexp
	symbolicAdsTag     *regexp.Regexp
	arrayInfoSegment   *regexp.Regexp
	driverContext      *DriverContext
}

// NewTagHandler this constructor creates a version of the TagHandler that's detached from a connection and can't provide context-sensitive feedback.
func NewTagHandler() TagHandler {
	return TagHandler{
		directAdsStringTag: regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		directAdsTag:       regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>\w+)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		symbolicAdsTag:     regexp.MustCompile(`^(?P<symbolicAddress>[^\[]+)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		arrayInfoSegment:   regexp.MustCompile(`((^(?P<numElements>\d+)$)|(^((?P<startElement>\d+)\.\.(?P<endElement>\d+))$)|(^((?P<startElement2>\d+):(?P<numElements2>\d+)))$)`),
	}
}

// NewTagHandlerWithDriverContext this constructor creates a version of the TagHandler that is connected to a connection and can provide context-sensitive feedback.
func NewTagHandlerWithDriverContext(driverContext *DriverContext) TagHandler {
	return TagHandler{
		directAdsStringTag: regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		directAdsTag:       regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>\w+)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		symbolicAdsTag:     regexp.MustCompile(`^(?P<symbolicAddress>[^\[]+)(?P<arrayInfo>((\[(\d+)])|(\[(\d+)\.\.(\d+)])|(\[(\d+):(\d+)]))*)`),
		arrayInfoSegment:   regexp.MustCompile(`((^(?P<numElements>\d+)$)|(^((?P<startElement>\d+)\.\.(?P<endElement>\d+))$)|(^((?P<startElement2>\d+):(?P<numElements2>\d+)))$)`),
		driverContext:      driverContext,
	}
}

func (m TagHandler) ParseTag(query string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.directAdsStringTag, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			decodeString, err := hex.DecodeString(indexGroupHexString)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			// Fill up the array with missing bytes to get an array of size 4 bytes.
			for i := len(decodeString); i < 4; i++ {
				decodeString = append([]byte{0}, decodeString...)
			}
			indexGroup = binary.BigEndian.Uint32(decodeString)
		} else {
			parsedIndexGroup, err := strconv.ParseUint(match["indexGroup"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexGroup = uint32(parsedIndexGroup)
		}
		var indexOffset uint32
		if indexOffsetHexString := match["indexOffsetHex"]; indexOffsetHexString != "" {
			decodeString, err := hex.DecodeString(indexOffsetHexString)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			// Fill up the array with missing bytes to get an array of size 4 bytes.
			for i := len(decodeString); i < 4; i++ {
				decodeString = append([]byte{0}, decodeString...)
			}
			indexOffset = binary.BigEndian.Uint32(decodeString)
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
		plcValueType, ok := values.PlcValueByName(adsDataTypeName)
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

		if match["arrayInfo"] != "" {
			arrayInfoString := match["arrayInfo"]

			arrayInfo = []apiModel.ArrayInfo{}

			// Cut off the starting and ending bracket
			arrayInfoString = arrayInfoString[1:(len(arrayInfoString) - 1)]
			// Split the remaining string into separate segments.
			arrayInfoSegments := strings.Split(arrayInfoString, "][")
			for _, currentSegment := range arrayInfoSegments {
				if match := utils.GetSubgroupMatches(m.arrayInfoSegment, currentSegment); match != nil {
					if match["startElement"] != "" && match["endElement"] != "" {
						startElement, err := m.getUint32Value(match["startElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						endElement, err := m.getUint32Value(match["endElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: endElement,
						})
					} else if match["startElement2"] != "" && match["numElements2"] != "" {
						startElement, err := m.getUint32Value(match["startElement2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						numElements, err := m.getUint32Value(match["numElements2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: startElement + numElements,
						})
					} else if match["numElements"] != "" {
						numElements, err := m.getUint32Value(match["numElements"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: numElements,
						})
					}
				}
			}
		}

		return model.NewDirectAdsPlcTag(indexGroup, indexOffset, plcValueType, stringLength, arrayInfo)
	} else if match := utils.GetSubgroupMatches(m.directAdsTag, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			decodeString, err := hex.DecodeString(indexGroupHexString)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			// Fill up the array with missing bytes to get an array of size 4 bytes.
			for i := len(decodeString); i < 4; i++ {
				decodeString = append([]byte{0}, decodeString...)
			}
			indexGroup = binary.BigEndian.Uint32(decodeString)
		} else {
			parsedIndexGroup, err := strconv.ParseUint(match["indexGroup"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexGroup = uint32(parsedIndexGroup)
		}
		var indexOffset uint32
		if indexOffsetHexString := match["indexOffsetHex"]; indexOffsetHexString != "" {
			decodeString, err := hex.DecodeString(indexOffsetHexString)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			// Fill up the array with missing bytes to get an array of size 4 bytes.
			for i := len(decodeString); i < 4; i++ {
				decodeString = append([]byte{0}, decodeString...)
			}
			indexOffset = binary.BigEndian.Uint32(decodeString)
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
		plcValueType, ok := values.PlcValueByName(adsDataTypeName)
		if !ok {
			return nil, fmt.Errorf("invalid ads data type")
		}

		var arrayInfo []apiModel.ArrayInfo

		if match["arrayInfo"] != "" {
			arrayInfoString := match["arrayInfo"]

			arrayInfo = []apiModel.ArrayInfo{}

			// Cut off the starting and ending bracket
			arrayInfoString = arrayInfoString[1:(len(arrayInfoString) - 1)]
			// Split the remaining string into separate segments.
			arrayInfoSegments := strings.Split(arrayInfoString, "][")
			for _, currentSegment := range arrayInfoSegments {
				if match := utils.GetSubgroupMatches(m.arrayInfoSegment, currentSegment); match != nil {
					if match["startElement"] != "" && match["endElement"] != "" {
						startElement, err := m.getUint32Value(match["startElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						endElement, err := m.getUint32Value(match["endElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: endElement,
						})
					} else if match["startElement2"] != "" && match["numElements2"] != "" {
						startElement, err := m.getUint32Value(match["startElement2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						numElements, err := m.getUint32Value(match["numElements2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: startElement + numElements,
						})
					} else if match["numElements"] != "" {
						numElements, err := m.getUint32Value(match["numElements"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: numElements,
						})
					}
				}
			}
		}

		return model.NewDirectAdsPlcTag(indexGroup, indexOffset, plcValueType, model.NONE, arrayInfo)
	} else if match := utils.GetSubgroupMatches(m.symbolicAdsTag, query); match != nil {
		var arrayInfo []apiModel.ArrayInfo

		if match["arrayInfo"] != "" {
			arrayInfoString := match["arrayInfo"]

			arrayInfo = []apiModel.ArrayInfo{}

			// Cut off the starting and ending bracket
			arrayInfoString = arrayInfoString[1:(len(arrayInfoString) - 1)]
			// Split the remaining string into separate segments.
			arrayInfoSegments := strings.Split(arrayInfoString, "][")
			for _, currentSegment := range arrayInfoSegments {
				if match := utils.GetSubgroupMatches(m.arrayInfoSegment, currentSegment); match != nil {
					if match["startElement"] != "" && match["endElement"] != "" {
						startElement, err := m.getUint32Value(match["startElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						endElement, err := m.getUint32Value(match["endElement"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: endElement,
						})
					} else if match["startElement2"] != "" && match["numElements2"] != "" {
						startElement, err := m.getUint32Value(match["startElement2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						numElements, err := m.getUint32Value(match["numElements2"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: startElement,
							UpperBound: startElement + numElements,
						})
					} else if match["numElements"] != "" {
						numElements, err := m.getUint32Value(match["numElements"])
						if err != nil {
							return nil, fmt.Errorf("error parsing array info: %s, got error: %v", currentSegment, err)
						}
						arrayInfo = append(arrayInfo, model2.DefaultArrayInfo{
							LowerBound: 0,
							UpperBound: numElements,
						})
					}
				}
			}
		}

		return model.NewAdsSymbolicPlcTag(match["symbolicAddress"], arrayInfo)
	} else {
		return nil, errors.Errorf("Invalid address format for address '%s'", query)
	}
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
