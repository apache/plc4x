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
	"context"
	"strconv"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"

	"github.com/pkg/errors"
)

func NumericGroupAddressToString(numericAddress uint16, groupAddress GroupAddressTag) (string, error) {
	if groupAddress == nil {
		return "", nil
	}
	switch groupAddress.(type) {
	case GroupAddress3LevelPlcTag:
		main := numericAddress >> 11
		middle := (numericAddress >> 8) & 0x07
		sub := numericAddress & 0xFF
		return strconv.Itoa(int(main)) + "/" + strconv.Itoa(int(middle)) + "/" + strconv.Itoa(int(sub)), nil
	case GroupAddress2LevelPlcTag:
		main := numericAddress >> 11
		sub := numericAddress & 0x07FF
		return strconv.Itoa(int(main)) + "/" + strconv.Itoa(int(sub)), nil
	case GroupAddress1LevelPlcTag:
		return strconv.Itoa(int(numericAddress)), nil
	default:
		return "", errors.Errorf("Unmapped %T", groupAddress)
	}
}

func GroupAddressToString(groupAddress driverModel.KnxGroupAddress) (string, error) {
	if groupAddress == nil {
		return "", nil
	}
	switch groupAddress := groupAddress.(type) {
	case driverModel.KnxGroupAddress3Level:
		level3 := groupAddress
		return strconv.Itoa(int(level3.GetMainGroup())) + "/" + strconv.Itoa(int(level3.GetMiddleGroup())) + "/" + strconv.Itoa(int(level3.GetSubGroup())), nil
	case driverModel.KnxGroupAddress2Level:
		level2 := groupAddress
		return strconv.Itoa(int(level2.GetMainGroup())) + "/" + strconv.Itoa(int(level2.GetSubGroup())), nil
	case driverModel.KnxGroupAddressFreeLevel:
		level1 := groupAddress
		return strconv.Itoa(int(level1.GetSubGroup())), nil
	default:
		return "", errors.Errorf("Unmapped %T", groupAddress)
	}
}

func ByteArrayToKnxAddress(ctx context.Context, data []byte) driverModel.KnxAddress {
	knxAddress, err := driverModel.KnxAddressParse(ctx, data)
	if err != nil {
		return nil
	}
	return knxAddress
}

func KnxAddressToByteArray(knxAddress driverModel.KnxAddress) []byte {
	targetAddress := make([]byte, 2)
	targetAddress[0] = (knxAddress.GetMainGroup()&0xF)<<4 | (knxAddress.GetMiddleGroup() & 0xF)
	targetAddress[1] = knxAddress.GetSubGroup()
	return targetAddress
}

func Uint16ToKnxAddress(data uint16) driverModel.KnxAddress {
	main := uint8(data >> 12)
	middle := uint8(data>>8) & 0xF
	sub := uint8(data & 0xFF)
	knxAddress := driverModel.NewKnxAddress(
		main,
		middle,
		sub,
	)
	return knxAddress
}

func Uint16ToKnxGroupAddress(ctx context.Context, data uint16, numLevels uint8) driverModel.KnxGroupAddress {
	rawData := make([]uint8, 2)
	rawData[0] = uint8(data >> 8)
	rawData[1] = uint8(data & 0xFF)
	knxGroupAddress, err := driverModel.KnxGroupAddressParse(ctx, rawData, numLevels)
	if err != nil {
		return nil
	}
	return knxGroupAddress
}
