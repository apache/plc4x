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

package bacnetip

import (
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
		addressPattern: regexp.MustCompile(`^(?P<deviceIdentifier>(\d|\*))/(?P<objectType>(\d|\*))/(?P<objectInstance>(\d|\*))`),
	}
}

const (
	DEVICE_IDENTIFIER = "deviceIdentifier"
	OBJECT_TYPE       = "objectType"
	OBJECT_INSTANCE   = "objectInstance"
)

func (m FieldHandler) ParseQuery(query string) (model.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.addressPattern, query); match != nil {
		deviceIdentifier, err := strconv.ParseUint(match[DEVICE_IDENTIFIER], 10, 32)
		if err != nil {
			return nil, err
		}
		objectType, err := strconv.ParseUint(match[OBJECT_TYPE], 10, 16)
		if err != nil {
			return nil, err
		}
		objectInstance, err := strconv.ParseUint(match[OBJECT_INSTANCE], 10, 32)
		if err != nil {
			return nil, err
		}

		return NewField(uint32(deviceIdentifier), uint16(objectType), uint32(objectInstance)), nil
	}
	return nil, errors.Errorf("Unable to parse %s", query)
}
