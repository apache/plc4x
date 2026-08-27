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
	"regexp"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type TagHandler struct {
	groupAddress3Level             *regexp.Regexp
	groupAddress2Level             *regexp.Regexp
	groupAddress1Level             *regexp.Regexp
	deviceQuery                    *regexp.Regexp
	devicePropertyAddress          *regexp.Regexp
	deviceMemoryAddress            *regexp.Regexp
	deviceCommunicationObjectQuery *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		groupAddress3Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		groupAddress2Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))/(?P<subGroup>(\d{1,4}|\*|\[(\d{1,4}|\d{1,4}\-\d{1,4})(,(\d{1,4}|\d{1,4}\-\d{1,4}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		groupAddress1Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,5}|\*|\[(\d{1,5}|\d{1,5}\-\d{1,5})(,(\d{1,5}|\d{1,5}\-\d{1,5}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),

		deviceQuery: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))$`),
		// The two device forms carry a real element count, so they take the shared notation.
		// A property address has no type suffix, so the selection ends it.
		devicePropertyAddress:          regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<objectId>\d{1,3})\/(?P<propertyId>\d{1,3})(\/(?P<propertyIndex>\d{1,4}))?` + spiModel.ArrayGroupPattern + `$`),
		deviceMemoryAddress:            regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<address>[0-9a-fA-F]{1,8})` + spiModel.ArrayGroupPattern + `(:(?P<datatype>[a-zA-Z_]+))?$`),
		deviceCommunicationObjectQuery: regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#com-obj$`),
	}
}

// selectionOf reads the array expression a device address carries and returns how far past the
// written start the selection begins and how many elements it spans. An address with no
// expression reads one element where it says.
//
// The group-address forms are deliberately not routed through here: their brackets hold a set of
// group addresses to match - "[1-3,5]" - not an array selection, and they never described a
// count.
func selectionOf(expression string, tagAddress string) (uint64, uint64, error) {
	if expression == "" {
		return 0, 1, nil
	}
	dimensions, err := spiModel.ParseArrayExpression(expression, tagAddress, spiModel.SingleDimension)
	if err != nil {
		return 0, 0, err
	}
	dimension := dimensions[0]
	elements := uint64(dimension.GetSize())
	if elements > 0xFF {
		return 0, 0, errors.Errorf("A selection of %d elements in '%s' is more than the count "+
			"field can carry", elements, tagAddress)
	}
	return uint64(dimension.GetLowerBound() - dimension.GetBase()), elements, nil
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.groupAddress1Level, tagAddress); match != nil {
		tagTypeName, ok := match["datatype"]
		var tagType driverModel.KnxDatapointType
		if ok {
			tagType, ok = driverModel.KnxDatapointTypeByName(tagTypeName)
			if !ok {
				return nil, errors.Errorf("Unknown type %s", tagTypeName)
			}
		}
		return NewGroupAddress1LevelPlcTag(match["mainGroup"], &tagType), nil
	} else if match := utils.GetSubgroupMatches(m.groupAddress2Level, tagAddress); match != nil {
		tagTypeName, ok := match["datatype"]
		var tagType driverModel.KnxDatapointType
		if ok {
			tagType, ok = driverModel.KnxDatapointTypeByName(tagTypeName)
			if !ok {
				return nil, errors.Errorf("Unknown type %s", tagTypeName)
			}
		}
		return NewGroupAddress2LevelPlcTag(match["mainGroup"], match["subGroup"], &tagType), nil
	} else if match := utils.GetSubgroupMatches(m.groupAddress3Level, tagAddress); match != nil {
		tagTypeName, ok := match["datatype"]
		var tagType driverModel.KnxDatapointType
		if ok {
			tagType, ok = driverModel.KnxDatapointTypeByName(tagTypeName)
			if !ok {
				return nil, errors.Errorf("Unknown type %s", tagTypeName)
			}
		}
		return NewGroupAddress3LevelPlcTag(match["mainGroup"], match["middleGroup"], match["subGroup"], &tagType), nil
	} else if match := utils.GetSubgroupMatches(m.devicePropertyAddress, tagAddress); match != nil {
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
		// A property is read with a start index and a count, and the property index written in
		// the address is that start index - so a selection that starts past the first element
		// moves it, exactly as an offset moves the address of a memory-addressed driver.
		offset, numberOfElements, err := selectionOf(match["array"], tagAddress)
		if err != nil {
			return nil, err
		}
		propertyIndex += offset
		return NewDevicePropertyAddressPlcTag(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup), uint8(objectId), uint8(propertyId),
			uint16(propertyIndex), uint8(numberOfElements)), nil
	} else if match := utils.GetSubgroupMatches(m.deviceMemoryAddress, tagAddress); match != nil {
		tagTypeName, ok := match["datatype"]
		// This is a 0-255 valued 1-byte value.
		tagType := driverModel.KnxDatapointType_DPT_DecimalFactor
		if ok && len(tagTypeName) > 0 {
			tagType, _ = driverModel.KnxDatapointTypeByName(tagTypeName)
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
		offset, numberOfElements, err := selectionOf(match["array"], tagAddress)
		if err != nil {
			return nil, err
		}
		if offset != 0 {
			// Unlike the property form there is nothing exact to move here: what one element
			// occupies depends on the datapoint type, which is measured in bits and need not be
			// a whole number of bytes. Guessing a byte size would move the address to somewhere
			// no one asked for, so this is reported instead.
			return nil, errors.Errorf("Array selection in tag '%s' must start at the first "+
				"element: a memory address is moved in bytes and a datapoint type is measured "+
				"in bits, so there is no offset to apply", tagAddress)
		}
		return NewDeviceMemoryAddressPlcTag(uint8(mainGroup), uint8(middleGroup), uint8(subGroup), address, uint8(numberOfElements), &tagType), nil
	} else if match := utils.GetSubgroupMatches(m.deviceCommunicationObjectQuery, tagAddress); match != nil {
		mainGroup, _ := strconv.ParseUint(match["mainGroup"], 10, 8)
		middleGroup, _ := strconv.ParseUint(match["middleGroup"], 10, 8)
		subGroup, _ := strconv.ParseUint(match["subGroup"], 10, 8)
		return NewCommunicationObjectQuery(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup)), nil
	}
	return nil, spiModel.InvalidAddressError(tagAddress,
		"a group address, or {area}.{line}.{device}#{object}/{property}[selection] - "+
			"for example 1.2.3#11/1[0..3]")
}

func (m TagHandler) ParseQuery(query string) (apiModel.PlcQuery, error) {
	// The "#com-obj" form has to be checked first as the plain device-query pattern
	// would otherwise never see it (it doesn't allow any suffix, so it simply wouldn't
	// match, leaving communication-object browsing unreachable).
	if match := utils.GetSubgroupMatches(m.deviceCommunicationObjectQuery, query); match != nil {
		mainGroup, _ := strconv.ParseUint(match["mainGroup"], 10, 8)
		middleGroup, _ := strconv.ParseUint(match["middleGroup"], 10, 8)
		subGroup, _ := strconv.ParseUint(match["subGroup"], 10, 8)
		return NewCommunicationObjectQuery(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup)), nil
	}
	if match := utils.GetSubgroupMatches(m.deviceQuery, query); match != nil {
		return NewDeviceQuery(
			match["mainGroup"], match["middleGroup"], match["subGroup"]), nil
	}
	return nil, errors.New("Invalid address format for query '" + query + "'")
}
