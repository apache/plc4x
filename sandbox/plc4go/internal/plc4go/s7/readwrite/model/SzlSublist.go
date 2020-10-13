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
package model

import "plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"

type SzlSublist uint8

type ISzlSublist interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

const (
	SzlSublist_MODULE_IDENTIFICATION                                         SzlSublist = 0x11
	SzlSublist_CPU_FEATURES                                                  SzlSublist = 0x12
	SzlSublist_USER_MEMORY_AREA                                              SzlSublist = 0x13
	SzlSublist_SYSTEM_AREAS                                                  SzlSublist = 0x14
	SzlSublist_BLOCK_TYPES                                                   SzlSublist = 0x15
	SzlSublist_STATUS_MODULE_LEDS                                            SzlSublist = 0x19
	SzlSublist_COMPONENT_IDENTIFICATION                                      SzlSublist = 0x1C
	SzlSublist_INTERRUPT_STATUS                                              SzlSublist = 0x22
	SzlSublist_ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS           SzlSublist = 0x25
	SzlSublist_COMMUNICATION_STATUS_DATA                                     SzlSublist = 0x32
	SzlSublist_STATUS_SINGLE_MODULE_LED                                      SzlSublist = 0x74
	SzlSublist_DP_MASTER_SYSTEM_INFORMATION                                  SzlSublist = 0x90
	SzlSublist_MODULE_STATUS_INFORMATION                                     SzlSublist = 0x91
	SzlSublist_RACK_OR_STATION_STATUS_INFORMATION                            SzlSublist = 0x92
	SzlSublist_RACK_OR_STATION_STATUS_INFORMATION_2                          SzlSublist = 0x94
	SzlSublist_ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION SzlSublist = 0x95
	SzlSublist_MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP         SzlSublist = 0x96
	SzlSublist_DIAGNOSTIC_BUFFER                                             SzlSublist = 0xA0
	SzlSublist_MODULE_DIAGNOSTIC_DATA                                        SzlSublist = 0xB1
)

func CastSzlSublist(structType interface{}) SzlSublist {
	castFunc := func(typ interface{}) SzlSublist {
		if sSzlSublist, ok := typ.(SzlSublist); ok {
			return sSzlSublist
		}
		return 0
	}
	return castFunc(structType)
}

func (m SzlSublist) LengthInBits() uint16 {
	return 8
}

func (m SzlSublist) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func SzlSublistParse(io *spi.ReadBuffer) (SzlSublist, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e SzlSublist) Serialize(io spi.WriteBuffer) error {
	// TODO: Implement ...
	return nil
}
