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
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"strconv"
)

func GroupAddressToString(groupAddress *driverModel.KnxGroupAddress) string {
	if groupAddress != nil {
		switch groupAddress.Child.(type) {
		case *driverModel.KnxGroupAddress3Level:
			level3 := driverModel.CastKnxGroupAddress3Level(groupAddress)
			return strconv.Itoa(int(level3.MainGroup)) + "/" + strconv.Itoa(int(level3.MiddleGroup)) + "/" + strconv.Itoa(int(level3.SubGroup))
		case *driverModel.KnxGroupAddress2Level:
			level2 := driverModel.CastKnxGroupAddress2Level(groupAddress)
			return strconv.Itoa(int(level2.MainGroup)) + "/" + strconv.Itoa(int(level2.SubGroup))
		case *driverModel.KnxGroupAddressFreeLevel:
			level1 := driverModel.CastKnxGroupAddressFreeLevel(groupAddress)
			return strconv.Itoa(int(level1.SubGroup))
		}
	}
	return ""
}

func FieldToKnxAddress(field KnxNetIpField) *driverModel.KnxAddress {
	if field.IsPatternField() {
		return nil
	}
	var mainAddress int
	var middleAddress int
	var subAddress int
	switch field.(type) {
	case KnxNetIpDevicePropertyAddressPlcField:
		plcField := field.(KnxNetIpDevicePropertyAddressPlcField)
		mainAddress, _ = strconv.Atoi(plcField.MainGroup)
		middleAddress, _ = strconv.Atoi(plcField.MiddleGroup)
		subAddress, _ = strconv.Atoi(plcField.SubGroup)
	case KnxNetIpGroupAddress3LevelPlcField:
		plcField := field.(KnxNetIpGroupAddress3LevelPlcField)
		mainAddress, _ = strconv.Atoi(plcField.MainGroup)
		middleAddress, _ = strconv.Atoi(plcField.MiddleGroup)
		subAddress, _ = strconv.Atoi(plcField.SubGroup)
	default:
		return nil
	}

	return driverModel.NewKnxAddress(uint8(mainAddress), uint8(middleAddress), uint8(subAddress))
}

func Int8ArrayToKnxAddress(data []int8) *driverModel.KnxAddress {
	readBuffer := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(data))
	knxAddress, err := driverModel.KnxAddressParse(readBuffer)
	if err != nil {
		return nil
	}
	return knxAddress
}

func KnxAddressToInt8Array(knxAddress driverModel.KnxAddress) []int8 {
	targetAddress := make([]int8, 2)
	targetAddress[0] = int8((knxAddress.MainGroup&0xF)<<4 | (knxAddress.MiddleGroup & 0xF))
	targetAddress[1] = int8(knxAddress.SubGroup)
	return targetAddress
}
