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

package bacnetip

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
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
		deviceIdentifierString := match[DEVICE_IDENTIFIER]
		var deviceIdentifier uint32
		if deviceIdentifierString == "*" {
			// TODO: find a way to express a wildcard. -1 not an option here
			deviceIdentifier = 0
		} else {
			if parsedDeviceIdentifier, err := strconv.ParseUint(deviceIdentifierString, 10, 32); err != nil {
				return nil, err
			} else {
				deviceIdentifier = uint32(parsedDeviceIdentifier)
			}
		}
		objectTypeString := match[OBJECT_TYPE]
		var objectType uint16
		if objectTypeString == "*" {
			// TODO: find a way to express a wildcard. -1 not an option here
			deviceIdentifier = 0
		} else {
			if parsedObjectType, err := strconv.ParseUint(objectTypeString, 10, 16); err != nil {
				return nil, err
			} else {
				objectType = uint16(parsedObjectType)
			}
		}
		objectInstanceString := match[OBJECT_INSTANCE]
		var objectInstance uint32
		if objectInstanceString == "*" {
			// TODO: find a way to express a wildcard. -1 not an option here
			objectInstance = 0
		} else {
			if parsedObjectInstance, err := strconv.ParseUint(objectInstanceString, 10, 32); err != nil {
				return nil, err
			} else {
				objectInstance = uint32(parsedObjectInstance)
			}
		}

		return NewField(deviceIdentifier, objectType, objectInstance), nil
	}
	return nil, errors.Errorf("Unable to parse %s", query)
}
