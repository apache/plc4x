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

package eip

import (
	"fmt"
	"regexp"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type TagHandler struct {
	addressPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		// The selection sits before the type, as it does everywhere else. The trailing
		// ":elementNb" count is gone: "%rate:DINT:4" is now written "%rate[0..3]:DINT".
		addressPattern: regexp.MustCompile(`^%(?P<tag>[%a-zA-Z_.0-9]+)` + spiModel.ArrayGroupPattern + `(?::(?P<dataType>[A-Z]+))?$`),
	}
}

const (
	TAG       = "tag"
	DATA_TYPE = "dataType"
	ARRAY     = "array"
)

// eipConstraints is what a CIP request can encode of a selection: one dimension, starting no
// later than 255 - the member segment that carries the offset is a uint 8.
var eipConstraints = spiModel.SingleDimension.
	WithMaxIndex(255).
	WithOnlyTrailingDimensionMayBeRange(true)

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	matches := m.addressPattern.FindStringSubmatch(tagAddress)
	if matches == nil {
		// "%rate:DINT:4" and "%rate:DINT[4]" both used to parse. Neither does now, so say what
		// to write instead rather than reporting only that the address did not match.
		return nil, spiModel.InvalidAddressError(tagAddress, "%tag[selection]:TYPE - for example %rate[0..3]:DINT")
	}

	tagName := matches[m.addressPattern.SubexpIndex(TAG)]
	dataTypeStr := matches[m.addressPattern.SubexpIndex(DATA_TYPE)]
	arrayExpression := matches[m.addressPattern.SubexpIndex(ARRAY)]

	var dataType model.CIPDataTypeCode
	if dataTypeStr == "" {
		dataType = model.CIPDataTypeCode_DINT
	} else {
		var found bool
		dataType, found = model.CIPDataTypeCodeByName(dataTypeStr)
		if !found {
			return nil, fmt.Errorf("unknown data type: %s", dataTypeStr)
		}
	}

	selection, err := spiModel.ParseArrayExpression(arrayExpression, tagAddress, eipConstraints)
	if err != nil {
		return nil, err
	}

	return NewTagWithSelection(tagName, dataType, selection), nil
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
