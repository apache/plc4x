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
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

type TagHandler struct {
	addressPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		addressPattern: regexp.MustCompile(`^%(?P<tag>[%a-zA-Z_.0-9]+\[?[0-9]*]?):?(?P<dataType>[A-Z]*):?(?P<elementNb>[0-9]*)`),
	}
}

const (
	TAG        = "tag"
	DATA_TYPE  = "dataType"
	ELEMENT_NB = "elementNb"
)

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	matches := m.addressPattern.FindStringSubmatch(tagAddress)
	if matches == nil {
		return nil, fmt.Errorf("invalid tag address: %s", tagAddress)
	}

	tagName := matches[m.addressPattern.SubexpIndex(TAG)]
	dataTypeStr := matches[m.addressPattern.SubexpIndex(DATA_TYPE)]
	elementNbStr := matches[m.addressPattern.SubexpIndex(ELEMENT_NB)]

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

	elementNb := uint16(1)
	if elementNbStr != "" {
		nb, err := strconv.ParseUint(elementNbStr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid element count: %s", elementNbStr)
		}
		elementNb = uint16(nb)
	}

	return NewTag(tagName, dataType, elementNb), nil
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
