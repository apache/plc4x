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
	"fmt"
	"strconv"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

type DeviceQuery struct {
	apiModel.PlcQuery

	MainGroup   string // 5 Bits: Values 0-31
	MiddleGroup string // 3 Bits: values 0-7
	SubGroup    string // 8 Bits
}

func NewDeviceQuery(mainGroup string, middleGroup string, subGroup string) DeviceQuery {
	return DeviceQuery{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
	}
}

func (k DeviceQuery) GetQueryString() string {
	return k.GetAddressString()
}

func (k DeviceQuery) GetAddressString() string {
	return fmt.Sprintf("%s.%s.%s", k.MainGroup, k.MiddleGroup, k.SubGroup)
}

func (k DeviceQuery) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k DeviceQuery) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

// toKnxAddress converts the query into a single KnxAddress. As a query can contain
// wildcards, ranges and lists, this is only possible if all three segments are plain
// numbers. Otherwise nil is returned and the caller has to expand the query first
// (see Browser.calculateAddresses).
func (k DeviceQuery) toKnxAddress() driverModel.KnxAddress {
	mainGroup, err := strconv.ParseUint(k.MainGroup, 10, 8)
	if err != nil {
		return nil
	}
	middleGroup, err := strconv.ParseUint(k.MiddleGroup, 10, 8)
	if err != nil {
		return nil
	}
	subGroup, err := strconv.ParseUint(k.SubGroup, 10, 8)
	if err != nil {
		return nil
	}
	return driverModel.NewKnxAddress(uint8(mainGroup), uint8(middleGroup), uint8(subGroup))
}

func (k DeviceQuery) String() string {
	return fmt.Sprintf("knx.DeviceQuery{%s}", k.GetAddressString())
}

type CommunicationObjectQuery struct {
	apiModel.PlcQuery

	MainGroup   uint8 // 5 Bits: Values 0-31
	MiddleGroup uint8 // 3 Bits: values 0-7
	SubGroup    uint8 // 8 Bits
	DeviceTag
}

func NewCommunicationObjectQuery(mainGroup uint8, middleGroup uint8, subGroup uint8) CommunicationObjectQuery {
	return CommunicationObjectQuery{
		MainGroup:   mainGroup,
		MiddleGroup: middleGroup,
		SubGroup:    subGroup,
	}
}

func (k CommunicationObjectQuery) GetQueryString() string {
	return k.GetAddressString()
}

func (k CommunicationObjectQuery) GetAddressString() string {
	return fmt.Sprintf("%d.%d.%d#com-obj",
		k.MainGroup, k.MiddleGroup, k.SubGroup)
}

func (k CommunicationObjectQuery) GetValueType() values.PlcValueType {
	return values.Struct
}

func (k CommunicationObjectQuery) GetArrayInfo() []apiModel.ArrayInfo {
	return []apiModel.ArrayInfo{}
}

func (k CommunicationObjectQuery) toKnxAddress() driverModel.KnxAddress {
	individualAddress := driverModel.NewKnxAddress(
		k.MainGroup,
		k.MiddleGroup,
		k.SubGroup,
	)
	return individualAddress
}

// String is defined explicitly as CommunicationObjectQuery embeds both apiModel.PlcQuery
// and DeviceTag, which would otherwise make the promoted fmt.Stringer ambiguous.
func (k CommunicationObjectQuery) String() string {
	return fmt.Sprintf("knx.CommunicationObjectQuery{%s}", k.GetAddressString())
}
