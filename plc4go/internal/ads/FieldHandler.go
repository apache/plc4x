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

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type FieldHandler struct {
	directAdsStringField *regexp.Regexp
	directAdsField       *regexp.Regexp
	symbolicAdsField     *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		directAdsStringField: regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)((\[(?P<numElements>\d+)])|(\[(?P<startElement>\d+)\.\.(?P<endElement>\d+)])|(\[(?P<startElement2>\d+):(?P<numElements2>\d+)]))?`),
		directAdsField:       regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>\w+)((\[(?P<numElements>\d+)])|(\[(?P<startElement>\d+)\.\.(?P<endElement>\d+)])|(\[(?P<startElement2>\d+):(?P<numElements2>\d+)]))?`),
		symbolicAdsField:     regexp.MustCompile(`^(?P<symbolicAddress>[^\[]+)((\[(?P<numElements>\d+)])|(\[(?P<startElement>\d+)\.\.(?P<endElement>\d+)])|(\[(?P<startElement2>\d+):(?P<numElements2>\d+)]))?`),
	}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.directAdsStringField, query); match != nil {
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
		adsDataType, ok := model.AdsDataTypeByName(adsDataTypeName)
		if !ok {
			return nil, fmt.Errorf("invalid ads data type")
		}

		stringLength := NONE
		numElements := NONE
		startElement := NONE
		endElement := NONE

		tmpStringLength, err := strconv.ParseInt(match["stringLength"], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Error decoding string length")
		}
		stringLength = int32(tmpStringLength)

		if match["numElements"] != "" {
			tmpNumElements, err := strconv.ParseUint(match["numElements"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements': %v", err)
			} else {
				numElements = int32(tmpNumElements)
			}
		} else if match["startElement"] != "" && match["endElement"] != "" {
			tmpStartElement, err := strconv.ParseUint(match["startElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement': %v", err)
			} else {
				startElement = int32(tmpStartElement)
			}
			tmpEndElement, err := strconv.ParseUint(match["endElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'endElement': %v", err)
			} else {
				endElement = int32(tmpEndElement)
			}
		} else if match["startElement2"] != "" && match["numElements2"] != "" {
			tmpStartElement2, err := strconv.ParseUint(match["startElement2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement2': %v", err)
			} else {
				startElement = int32(tmpStartElement2)
			}
			tmpNumElements2, err := strconv.ParseUint(match["numElements2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements2': %v", err)
			} else {
				numElements = int32(tmpNumElements2)
			}
		}

		return newDirectAdsPlcField(indexGroup, indexOffset, adsDataType, stringLength, numElements, startElement, endElement)
	} else if match := utils.GetSubgroupMatches(m.directAdsField, query); match != nil {
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
		adsDataType, ok := model.AdsDataTypeByName(adsDataTypeName)
		if !ok {
			return nil, fmt.Errorf("invalid ads data type")
		}

		numElements := NONE
		startElement := NONE
		endElement := NONE

		if match["numElements"] != "" {
			tmpNumElements, err := strconv.ParseUint(match["numElements"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements': %v", err)
			} else {
				numElements = int32(tmpNumElements)
			}
		} else if match["startElement"] != "" && match["endElement"] != "" {
			tmpStartElement, err := strconv.ParseUint(match["startElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement': %v", err)
			} else {
				startElement = int32(tmpStartElement)
			}
			tmpEndElement, err := strconv.ParseUint(match["endElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'endElement': %v", err)
			} else {
				endElement = int32(tmpEndElement)
			}
		} else if match["startElement2"] != "" && match["numElements2"] != "" {
			tmpStartElement2, err := strconv.ParseUint(match["startElement2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement2': %v", err)
			} else {
				startElement = int32(tmpStartElement2)
			}
			tmpNumElements2, err := strconv.ParseUint(match["numElements2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements2': %v", err)
			} else {
				numElements = int32(tmpNumElements2)
			}
		}

		return newDirectAdsPlcField(indexGroup, indexOffset, adsDataType, NONE, numElements, startElement, endElement)
	} else if match := utils.GetSubgroupMatches(m.symbolicAdsField, query); match != nil {
		numElements := NONE
		startElement := NONE
		endElement := NONE

		if match["numElements"] != "" {
			tmpNumElements, err := strconv.ParseUint(match["numElements"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements': %v", err)
			} else {
				numElements = int32(tmpNumElements)
			}
		} else if match["startElement"] != "" && match["endElement"] != "" {
			tmpStartElement, err := strconv.ParseUint(match["startElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement': %v", err)
			} else {
				startElement = int32(tmpStartElement)
			}
			tmpEndElement, err := strconv.ParseUint(match["endElement"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'endElement': %v", err)
			} else {
				endElement = int32(tmpEndElement)
			}
		} else if match["startElement2"] != "" && match["numElements2"] != "" {
			tmpStartElement2, err := strconv.ParseUint(match["startElement2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'startElement2': %v", err)
			} else {
				startElement = int32(tmpStartElement2)
			}
			tmpNumElements2, err := strconv.ParseUint(match["numElements2"], 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid address format parsing 'numElements2': %v", err)
			} else {
				numElements = int32(tmpNumElements2)
			}
		}

		return newAdsSymbolicPlcField(match["symbolicAddress"], numElements, startElement, endElement)
	} else {
		return nil, errors.Errorf("Invalid address format for address '%s'", query)
	}
}
