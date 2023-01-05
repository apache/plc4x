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

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type TagHandler struct {
	addressPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		addressPattern: regexp.MustCompile(`^%(?P<tag>[a-zA-Z_.0-9]+\[?[0-9]*]?):?(?P<dataType>[A-Z]*):?(?P<elementNb>[0-9]*)`),
	}
}

const (
	TAG        = "tag"
	DATA_TYPE  = "dataType"
	ELEMENT_NB = "elementNb"
)

func (m TagHandler) ParseTag(query string) (model.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.addressPattern, query); match != nil {
		tag := match[TAG]
		_type, ok := readWriteModel.CIPDataTypeCodeByName(match[DATA_TYPE])
		if !ok {
			return nil, errors.Errorf("Unknown type %s", match[DATA_TYPE])
		}
		parsedUint, _ := strconv.ParseUint(match[ELEMENT_NB], 10, 16)
		elementNb := uint16(parsedUint)
		return NewTag(tag, _type, elementNb), nil
	}
	return nil, errors.Errorf("Unable to parse %s", query)
}

func (m TagHandler) ParseQuery(query string) (model.PlcQuery, error) {
	return nil, fmt.Errorf("queries not supported")
}
