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

package eip

import (
	readWriteModel "github.com/apache/plc4x/plc4go/internal/plc4go/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

type FieldHandler struct {
	addressPattern *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		addressPattern: regexp.MustCompile(`^%(?P<tag>[a-zA-Z_.0-9]+\[?[0-9]*]?):?(?P<dataType>[A-Z]*):?(?P<elementNb>[0-9]*)`),
	}
}

const (
	TAG        = "tag"
	DATA_TYPE  = "dataType"
	ELEMENT_NB = "elementNb"
)

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.addressPattern, query); match != nil {
		tag := match[TAG]
		_type := readWriteModel.CIPDataTypeCodeByName(match[DATA_TYPE])
		atoi, _ := strconv.Atoi(match[ELEMENT_NB])
		elementNb := uint16(atoi)
		return NewField(tag, _type, elementNb), nil
	}
	return nil, errors.Errorf("Unable to parse %s", query)
}
