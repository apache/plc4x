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

package ads

import (
	"encoding/binary"
	"encoding/hex"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
)

type FieldType uint8

//go:generate stringer -type FieldType
const (
	DirectAdsStringField   FieldType = 0x00
	DirectAdsField         FieldType = 0x01
	SymbolicAdsStringField FieldType = 0x03
	SymbolicAdsField       FieldType = 0x04
)

func (i FieldType) GetName() string {
	return i.String()
}

type FieldHandler struct {
	directAdsStringField   *regexp.Regexp
	directAdsField         *regexp.Regexp
	symbolicAdsStringField *regexp.Regexp
	symbolicAdsField       *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		directAdsStringField:   regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>STRING|WSTRING)\((?P<stringLength>\d{1,3})\)(\[(?P<numberOfElements>\d+)])?`),
		directAdsField:         regexp.MustCompile(`^((0[xX](?P<indexGroupHex>[0-9a-fA-F]+))|(?P<indexGroup>\d+))/((0[xX](?P<indexOffsetHex>[0-9a-fA-F]+))|(?P<indexOffset>\d+)):(?P<adsDataType>\w+)(\[(?P<numberOfElements>\d+)])?`),
		symbolicAdsStringField: regexp.MustCompile(`^(?P<symbolicAddress>.+):(?P<adsDataType>'STRING'|'WSTRING')\((?P<stringLength>\d{1,3})\)(\[(?P<numberOfElements>\d+)])?`),
		symbolicAdsField:       regexp.MustCompile(`^(?P<symbolicAddress>.+):(?P<adsDataType>\w+)(\[(?P<numberOfElements>\d+)])?`),
	}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.directAdsStringField, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			decodeString, err := hex.DecodeString(indexGroupHexString[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
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
			decodeString, err := hex.DecodeString(indexOffsetHexString[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = binary.BigEndian.Uint32(decodeString)
		} else {
			parsedIndexOffset, err := strconv.ParseUint(match["indexOffset"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = uint32(parsedIndexOffset)
		}
		stringLength, err := strconv.ParseInt(match["stringLength"], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Error decoding string length")
		}
		numberOfElements, err := strconv.ParseUint(match["numberOfElements"], 10, 32)
		if err != nil {
			log.Trace().Msg("Falling back to number of elements 1")
			numberOfElements = 1
		}

		return newDirectAdsPlcField(indexGroup, indexOffset, model2.AdsDataTypeByName(match["adsDataType"]), int32(stringLength), uint32(numberOfElements))
	} else if match := utils.GetSubgroupMatches(m.directAdsField, query); match != nil {
		var indexGroup uint32
		if indexGroupHexString := match["indexGroupHex"]; indexGroupHexString != "" {
			decodeString, err := hex.DecodeString(indexGroupHexString[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
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
			decodeString, err := hex.DecodeString(indexOffsetHexString[2:])
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = binary.BigEndian.Uint32(decodeString)
		} else {
			parsedIndexOffset, err := strconv.ParseUint(match["indexOffset"], 10, 32)
			if err != nil {
				return nil, errors.Wrap(err, "Error decoding index group")
			}
			indexOffset = uint32(parsedIndexOffset)
		}

		adsDataType := model2.AdsDataTypeByName(match["adsDataType"])
		numberOfElements, err := strconv.ParseUint(match["numberOfElements"], 10, 32)
		if err != nil {
			log.Trace().Msg("Falling back to number of elements 1")
			numberOfElements = 1
		}
		return newDirectAdsPlcField(indexGroup, indexOffset, adsDataType, int32(0), uint32(numberOfElements))
	} else if match := utils.GetSubgroupMatches(m.symbolicAdsStringField, query); match != nil {
		stringLength, err := strconv.ParseInt(match["stringLength"], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Error decoding string length")
		}
		numberOfElements, err := strconv.ParseUint(match["numberOfElements"], 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, "Error decoding number of elements")
		}
		return newAdsSymbolicPlcField(match["symbolicAddress"], model2.AdsDataTypeByName(match["adsDataType"]), int32(stringLength), uint32(numberOfElements))
	} else if match := utils.GetSubgroupMatches(m.symbolicAdsField, query); match != nil {
		numberOfElements, err := strconv.ParseUint(match["numberOfElements"], 10, 32)
		if err != nil {
			log.Trace().Msg("Falling back to number of elements 1")
			numberOfElements = 1
		}
		return newAdsSymbolicPlcField(match["symbolicAddress"], model2.AdsDataTypeByName(match["adsDataType"]), int32(0), uint32(numberOfElements))
	} else {
		return nil, errors.Errorf("Invalid address format for address '%s'", query)
	}
}
