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

type BACnetObjectType uint16

const (
	BACnetObjectType_ANALOG_INPUT           BACnetObjectType = 0x000
	BACnetObjectType_ANALOG_OUTPUT          BACnetObjectType = 0x001
	BACnetObjectType_ANALOG_VALUE           BACnetObjectType = 0x002
	BACnetObjectType_BINARY_INPUT           BACnetObjectType = 0x003
	BACnetObjectType_BINARY_OUTPUT          BACnetObjectType = 0x004
	BACnetObjectType_BINARY_VALUE           BACnetObjectType = 0x005
	BACnetObjectType_MULTISTATE_INPUT       BACnetObjectType = 0x00D
	BACnetObjectType_MULTISTATE_OUTPUT      BACnetObjectType = 0x00E
	BACnetObjectType_MULTISTATE_VALUE       BACnetObjectType = 0x013
	BACnetObjectType_SCHEDULE               BACnetObjectType = 0x011
	BACnetObjectType_DEVICE                 BACnetObjectType = 0x008
	BACnetObjectType_CALENDAR               BACnetObjectType = 0x006
	BACnetObjectType_COMMAND                BACnetObjectType = 0x007
	BACnetObjectType_EVENT_ENROLLMENT       BACnetObjectType = 0x009
	BACnetObjectType_FILE                   BACnetObjectType = 0x00A
	BACnetObjectType_GROUP                  BACnetObjectType = 0x00B
	BACnetObjectType_LOOP                   BACnetObjectType = 0x00C
	BACnetObjectType_NOTIFICATION_CLASS     BACnetObjectType = 0x00F
	BACnetObjectType_PROGRAM                BACnetObjectType = 0x010
	BACnetObjectType_AVERAGING              BACnetObjectType = 0x012
	BACnetObjectType_TREND_LOG              BACnetObjectType = 0x014
	BACnetObjectType_LIFE_SAFETY_POINT      BACnetObjectType = 0x015
	BACnetObjectType_LIFE_SAFETY_ZONE       BACnetObjectType = 0x016
	BACnetObjectType_ACCUMULATOR            BACnetObjectType = 0x017
	BACnetObjectType_PULSE_CONVERTER        BACnetObjectType = 0x018
	BACnetObjectType_EVENT_LOG              BACnetObjectType = 0x019
	BACnetObjectType_GLOBAL_GROUP           BACnetObjectType = 0x01A
	BACnetObjectType_TREND_LOG_MULTIPLE     BACnetObjectType = 0x01B
	BACnetObjectType_LOAD_CONTROL           BACnetObjectType = 0x01C
	BACnetObjectType_STRUCTURED_VIEW        BACnetObjectType = 0x01D
	BACnetObjectType_ACCESS_DOOR            BACnetObjectType = 0x01E
	BACnetObjectType_TIMER                  BACnetObjectType = 0x01F
	BACnetObjectType_ACCESS_CREDENTIAL      BACnetObjectType = 0x020
	BACnetObjectType_ACCESS_POINT           BACnetObjectType = 0x021
	BACnetObjectType_ACCESS_RIGHTS          BACnetObjectType = 0x022
	BACnetObjectType_ACCESS_USER            BACnetObjectType = 0x023
	BACnetObjectType_ACCESS_ZONE            BACnetObjectType = 0x024
	BACnetObjectType_CREDENTIAL_DATA_INPUT  BACnetObjectType = 0x025
	BACnetObjectType_NETWORK_SECURITY       BACnetObjectType = 0x026
	BACnetObjectType_BITSTRING_VALUE        BACnetObjectType = 0x027
	BACnetObjectType_CHARACTERSTRING_VALUE  BACnetObjectType = 0x028
	BACnetObjectType_DATEPATTERN_VALUE      BACnetObjectType = 0x029
	BACnetObjectType_DATE_VALUE             BACnetObjectType = 0x02A
	BACnetObjectType_DATETIMEPATTERN_VALUE  BACnetObjectType = 0x02B
	BACnetObjectType_DATETIME_VALUE         BACnetObjectType = 0x02C
	BACnetObjectType_INTEGER_VALUE          BACnetObjectType = 0x02D
	BACnetObjectType_LARGE_ANALOG_VALUE     BACnetObjectType = 0x02E
	BACnetObjectType_OCTETSTRING_VALUE      BACnetObjectType = 0x02F
	BACnetObjectType_POSITIVE_INTEGER_VALUE BACnetObjectType = 0x030
	BACnetObjectType_TIMEPATTERN_VALUE      BACnetObjectType = 0x031
	BACnetObjectType_TIME_VALUE             BACnetObjectType = 0x032
	BACnetObjectType_NOTIFICATION_FORWARDER BACnetObjectType = 0x033
	BACnetObjectType_ALERT_ENROLLMENT       BACnetObjectType = 0x034
	BACnetObjectType_CHANNEL                BACnetObjectType = 0x035
	BACnetObjectType_LIGHTING_OUTPUT        BACnetObjectType = 0x036
	BACnetObjectType_BINARY_LIGHTING_OUTPUT BACnetObjectType = 0x037
	BACnetObjectType_NETWORK_PORT           BACnetObjectType = 0x038
	BACnetObjectType_ELEVATOR_GROUP         BACnetObjectType = 0x039
	BACnetObjectType_ESCALATOR              BACnetObjectType = 0x03A
)

func CastBACnetObjectType(structType interface{}) BACnetObjectType {
	castFunc := func(typ interface{}) BACnetObjectType {
		if sBACnetObjectType, ok := typ.(BACnetObjectType); ok {
			return sBACnetObjectType
		}
		return 0
	}
	return castFunc(structType)
}

func BACnetObjectTypeParse(io spi.ReadBuffer) (BACnetObjectType, error) {
	// TODO: Implement ...
	return 0, nil
}

func (e BACnetObjectType) Serialize(io spi.WriteBuffer) {
	// TODO: Implement ...
}
