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
package knxnetip

import (
	"encoding/hex"
	"errors"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"regexp"
	"strconv"
)

type KnxNetIpFieldType uint8

const (
	KNX_NET_IP_FIELD_COIL KnxNetIpFieldType = 0x00
)

func (m KnxNetIpFieldType) GetName() string {
	switch m {
	case KNX_NET_IP_FIELD_COIL:
		return "ModbusFieldHoldingRegister"
	}
	return ""
}

type FieldHandler struct {
	knxNetIpGroupAddress3Level             *regexp.Regexp
	knxNetIpGroupAddress2Level             *regexp.Regexp
	knxNetIpGroupAddress1Level             *regexp.Regexp
	knxNetIpDeviceQuery                    *regexp.Regexp
	knxNetIpDevicePropertyAddress          *regexp.Regexp
	knxNetIpDeviceMemoryAddress            *regexp.Regexp
	knxNetIpDeviceCommunicationObjectQuery *regexp.Regexp
	spi.PlcFieldHandler
}

func NewFieldHandler() FieldHandler {
	return FieldHandler{
		knxNetIpGroupAddress3Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\/(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		knxNetIpGroupAddress2Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))/(?P<subGroup>(\d{1,4}|\*|\[(\d{1,4}|\d{1,4}\-\d{1,4})(,(\d{1,4}|\d{1,4}\-\d{1,4}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),
		knxNetIpGroupAddress1Level: regexp.MustCompile(`^(?P<mainGroup>(\d{1,5}|\*|\[(\d{1,5}|\d{1,5}\-\d{1,5})(,(\d{1,5}|\d{1,5}\-\d{1,5}))*]))(:(?P<datatype>[a-zA-Z_]+))?$`),

		knxNetIpDeviceQuery:                    regexp.MustCompile(`^(?P<mainGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<middleGroup>(\d{1,2}|\*|\[(\d{1,2}|\d{1,2}\-\d{1,2})(,(\d{1,2}|\d{1,2}\-\d{1,2}))*]))\.(?P<subGroup>(\d{1,3}|\*|\[(\d{1,3}|\d{1,3}\-\d{1,3})(,(\d{1,3}|\d{1,3}\-\d{1,3}))*]))$`),
		knxNetIpDevicePropertyAddress:          regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<objectId>\d{1,3})\/(?P<propertyId>\d{1,3})(\/(?P<propertyIndex>\d{1,4}))?(\[(?P<numElements>\d{1,2})])?$`),
		knxNetIpDeviceMemoryAddress:            regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#(?P<address>[0-9a-fA-F]{1,8})(:(?P<datatype>[a-zA-Z_]+)(\[(?P<numElements>\d+)])?)?$`),
		knxNetIpDeviceCommunicationObjectQuery: regexp.MustCompile(`^(?P<mainGroup>\d{1,2})\.(?P<middleGroup>\d)\.(?P<subGroup>\d{1,3})#com-obj$`),
	}
}

func (m FieldHandler) ParseQuery(query string) (apiModel.PlcField, error) {
	if match := utils.GetSubgroupMatches(m.knxNetIpGroupAddress1Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewKnxNetIpGroupAddress1LevelPlcField(match["mainGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpGroupAddress2Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewKnxNetIpGroupAddress2LevelPlcField(match["mainGroup"], match["subGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpGroupAddress3Level, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		var fieldType driverModel.KnxDatapointType
		if ok {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		return NewKnxNetIpGroupAddress3LevelPlcField(match["mainGroup"], match["middleGroup"], match["subGroup"], &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpDeviceQuery, query); match != nil {
		return NewKnxNetIpDeviceQueryField(
			match["mainGroup"], match["middleGroup"], match["subGroup"]), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpDevicePropertyAddress, query); match != nil {
		mainGroup, _ := strconv.Atoi(match["mainGroup"])
		middleGroup, _ := strconv.Atoi(match["middleGroup"])
		subGroup, _ := strconv.Atoi(match["subGroup"])
		objectId, _ := strconv.Atoi(match["objectId"])
		propertyId, _ := strconv.Atoi(match["propertyId"])
		propertyIndex := 1
		propertyInd, ok := match["propertyIndex"]
		if ok && len(propertyInd) > 0 {
			propertyIndex, _ = strconv.Atoi(propertyInd)
		}
		numberOfElements := 1
		numElements, ok := match["numElements"]
		if ok && len(numElements) > 0 {
			numberOfElements, _ = strconv.Atoi(numElements)
		}
		return NewKnxNetIpDevicePropertyAddressPlcField(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup), uint8(objectId), uint8(propertyId),
			uint16(propertyIndex), uint8(numberOfElements)), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpDeviceMemoryAddress, query); match != nil {
		fieldTypeName, ok := match["datatype"]
		// This is a 0-255 valued 1-byte value.
		fieldType := driverModel.KnxDatapointType_DPT_DecimalFactor
		if ok && len(fieldTypeName) > 0 {
			fieldType = driverModel.KnxDatapointTypeByName(fieldTypeName)
		}
		mainGroup, _ := strconv.Atoi(match["mainGroup"])
		middleGroup, _ := strconv.Atoi(match["middleGroup"])
		subGroup, _ := strconv.Atoi(match["subGroup"])
		addressData, _ := hex.DecodeString(match["address"])
		var address uint16
		if len(addressData) == 2 {
			address = uint16(addressData[0])<<8 | uint16(addressData[1])
		} else if len(addressData) == 1 {
			address = uint16(addressData[0])
		} else {
			return nil, errors.New("invalid address: " + match["address"])
		}
		numberOfElements := 1
		numElements, ok := match["numElements"]
		if ok && len(numElements) > 0 {
			numberOfElements, _ = strconv.Atoi(numElements)
		}
		return NewKnxNetIpDeviceMemoryAddressPlcField(uint8(mainGroup), uint8(middleGroup), uint8(subGroup), address, uint8(numberOfElements), &fieldType), nil
	} else if match := utils.GetSubgroupMatches(m.knxNetIpDeviceCommunicationObjectQuery, query); match != nil {
		mainGroup, _ := strconv.Atoi(match["mainGroup"])
		middleGroup, _ := strconv.Atoi(match["middleGroup"])
		subGroup, _ := strconv.Atoi(match["subGroup"])
		return NewKnxNetIpCommunicationObjectQueryField(
			uint8(mainGroup), uint8(middleGroup), uint8(subGroup)), nil
	}
	return nil, errors.New("Invalid address format for address '" + query + "'")
}
