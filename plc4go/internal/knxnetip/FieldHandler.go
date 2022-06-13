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

package knxnetip

import (
	"encoding/hex"
	"errors"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"regexp"
	"strconv"
)

type FieldHandler struct {
	groupAddress3Level             *regexp.Regexp
	groupAddress2Level             *regexp.Regexp
	groupAddress1Level             *regexp.Regexp
	deviceQuery                    *regexp.Regexp
	devicePropertyAddress          *regexp.Regexp
	deviceMemoryAddress            *regexp.Regexp
	deviceCommunicationObjectQuery *regexp.Regexp
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		groupAddress3Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		groupAddress2Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))/(?P<subGroup>(\d{1,4}|\*|\[(\d{1,4}|\d{1,4}\-\d{1,4})(,(\d{1,4}|\d{1,4}\-\d{1,4}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		groupAddress1Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,5}|\*|\[(\d{1,5}|\d{1,5}\-\d{1,5})(,(\d{1,5}|\d{1,5}\-\d{1,5}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),

		deviceQuery:                    regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))$`),
		devicePropertyAddress:          regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<objectId>\d{1,3})\/(?P<propertyId>\d{1,3})(\/(?P<propertyIndex>\d{1,4}))?(\[(?P<numElements>\d{1,2})])?$`),
		deviceMemoryAddress:            regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<address>[0-9a-fA-F]{1,8})(:(?P<datatype>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?)?$`),
		deviceCommunicationObjectQuery: regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#com-obj$`),
	}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.groupAddress1Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewGroupAddress1LevelPlcField(match["mainGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.groupAddress2Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewGroupAddress2LevelPlcField(match["mainGroup"], match["subGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.groupAddress3Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewGroupAddress3LevelPlcField(match["mainGroup"], match["middleGroup"], match["subGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.deviceQuery, query); match != nil {
		return NewDeviceQueryField(
			match["mainGroup"], match["middleGroup"], match["subGroup"]), nil
	} else if match := utils.GetSubgroupMatches(m.devicePropertyAddress, query); match != nil {
		mainGroup, _ := strconv.ParseUint(match["mainGroup"], 10, 8)
		middleGroup, _ := strconv.ParseUint(match["middleGroup"], 10, 8)
		subGroup, _ := strconv.ParseUint(match["subGroup"], 10, 8)
		objectId, _ := strconv.ParseUint(match["objectId"], 10, 8)
		propertyId, _ := strconv.ParseUint(match["propertyId"], 10, 8)
		propertyIndex := uint64(1)
		propertyInd, ok := match["propertyIndex"]
		if ok && len(propertyInd) > 0 {
			propertyIndex, _ = strconv.ParseUint(propertyInd, 10, 16)
		}
		numberOfElements := uint64(1)
		numElements, ok := match["numElements"]
		if ok && len(numElements) > 0 {
			numberOfElements, _ = strconv.ParseUint(numElements, 10, 8)
		}
		return NewDevicePropertyAddressPlcField(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup), uint8(objectId), uint8(propertyId),
			uint16(propertyIndex), uint8(numberOfElements)), nil
	} else if match := utils.GetSubgroupMatches(m.deviceMemoryAddress, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		// This is a 0-255 valued 1-byte value.
		fieldType := driverModel.KnxDatapointType_DPT_DecimalFactor
		if ok && len(fieldTypeName) > 0 {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		mainGroup, _ := strconv.ParseUint(match["mainGroup"], 10, 8)
		middleGroup, _ := strconv.ParseUint(match["middleGroup"], 10, 8)
		subGroup, _ := strconv.ParseUint(match["subGroup"], 10, 8)
		addressData, _ := hex.DecodeString(match["address"])
		var address uint16
		if len(addressData) == 2 {
			address = uint16(addressData[0])<<8 | uint16(addressData[1])
		} else if len(addressData) == 1 {
			address = uint16(addressData[0])
		} else {
			return nil, errors.New("invalid address: " + match["address"])
		}
		numberOfElements := uint64(1)
		numElements, ok := match["numElements"]
		if ok && len(numElements) > 0 {
			numberOfElements, _ = strconv.ParseUint(numElements, 10, 8)
		}
		return NewDeviceMemoryAddressPlcField(uint8(mainGroup), uint8(middleGroup), uint8(subGroup), address, uint8(numberOfElements), &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.deviceCommunicationObjectQuery, query); match != nil {
		mainGroup, _ := strconv.ParseUint(match["mainGroup"], 10, 8)
		middleGroup, _ := strconv.ParseUint(match["middleGroup"], 10, 8)
		subGroup, _ := strconv.ParseUint(match["subGroup"], 10, 8)
		return NewCommunicationObjectQueryField(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup)), nil
	}
	return nil, errors.New("Invalid address format for address '" + query + "'")
}
