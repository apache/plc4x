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

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

type ApplicationId uint8

type IApplicationId interface {
	Serialize(writeBuffer utils.WriteBuffer) error
}

const (
	ApplicationId_RESERVED                           ApplicationId = 0x00
	ApplicationId_FREE_USAGE                         ApplicationId = 0x01
	ApplicationId_TEMPERATURE_BROADCAST              ApplicationId = 0x02
	ApplicationId_ROOM_CONTROL_SYSTEM                ApplicationId = 0x03
	ApplicationId_LIGHTING                           ApplicationId = 0x04
	ApplicationId_VENTILATION                        ApplicationId = 0x05
	ApplicationId_IRRIGATION_CONTROL                 ApplicationId = 0x06
	ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL ApplicationId = 0x07
	ApplicationId_HEATING                            ApplicationId = 0x08
	ApplicationId_AIR_CONDITIONING                   ApplicationId = 0x09
	ApplicationId_TRIGGER_CONTROL                    ApplicationId = 0x0A
	ApplicationId_ENABLE_CONTROL                     ApplicationId = 0x0B
	ApplicationId_AUDIO_AND_VIDEO                    ApplicationId = 0x0C
	ApplicationId_SECURITY                           ApplicationId = 0x0D
	ApplicationId_METERING                           ApplicationId = 0x0E
	ApplicationId_ACCESS_CONTROL                     ApplicationId = 0x0F
	ApplicationId_CLOCK_AND_TIMEKEEPING              ApplicationId = 0x10
	ApplicationId_TELEPHONY_STATUS_AND_CONTROL       ApplicationId = 0x11
	ApplicationId_MEASUREMENT                        ApplicationId = 0x12
	ApplicationId_TESTING                            ApplicationId = 0x13
)

var ApplicationIdValues []ApplicationId

func init() {
	_ = errors.New
	ApplicationIdValues = []ApplicationId{
		ApplicationId_RESERVED,
		ApplicationId_FREE_USAGE,
		ApplicationId_TEMPERATURE_BROADCAST,
		ApplicationId_ROOM_CONTROL_SYSTEM,
		ApplicationId_LIGHTING,
		ApplicationId_VENTILATION,
		ApplicationId_IRRIGATION_CONTROL,
		ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL,
		ApplicationId_HEATING,
		ApplicationId_AIR_CONDITIONING,
		ApplicationId_TRIGGER_CONTROL,
		ApplicationId_ENABLE_CONTROL,
		ApplicationId_AUDIO_AND_VIDEO,
		ApplicationId_SECURITY,
		ApplicationId_METERING,
		ApplicationId_ACCESS_CONTROL,
		ApplicationId_CLOCK_AND_TIMEKEEPING,
		ApplicationId_TELEPHONY_STATUS_AND_CONTROL,
		ApplicationId_MEASUREMENT,
		ApplicationId_TESTING,
	}
}

func ApplicationIdByValue(value uint8) ApplicationId {
	switch value {
	case 0x00:
		return ApplicationId_RESERVED
	case 0x01:
		return ApplicationId_FREE_USAGE
	case 0x02:
		return ApplicationId_TEMPERATURE_BROADCAST
	case 0x03:
		return ApplicationId_ROOM_CONTROL_SYSTEM
	case 0x04:
		return ApplicationId_LIGHTING
	case 0x05:
		return ApplicationId_VENTILATION
	case 0x06:
		return ApplicationId_IRRIGATION_CONTROL
	case 0x07:
		return ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL
	case 0x08:
		return ApplicationId_HEATING
	case 0x09:
		return ApplicationId_AIR_CONDITIONING
	case 0x0A:
		return ApplicationId_TRIGGER_CONTROL
	case 0x0B:
		return ApplicationId_ENABLE_CONTROL
	case 0x0C:
		return ApplicationId_AUDIO_AND_VIDEO
	case 0x0D:
		return ApplicationId_SECURITY
	case 0x0E:
		return ApplicationId_METERING
	case 0x0F:
		return ApplicationId_ACCESS_CONTROL
	case 0x10:
		return ApplicationId_CLOCK_AND_TIMEKEEPING
	case 0x11:
		return ApplicationId_TELEPHONY_STATUS_AND_CONTROL
	case 0x12:
		return ApplicationId_MEASUREMENT
	case 0x13:
		return ApplicationId_TESTING
	}
	return 0
}

func ApplicationIdByName(value string) ApplicationId {
	switch value {
	case "RESERVED":
		return ApplicationId_RESERVED
	case "FREE_USAGE":
		return ApplicationId_FREE_USAGE
	case "TEMPERATURE_BROADCAST":
		return ApplicationId_TEMPERATURE_BROADCAST
	case "ROOM_CONTROL_SYSTEM":
		return ApplicationId_ROOM_CONTROL_SYSTEM
	case "LIGHTING":
		return ApplicationId_LIGHTING
	case "VENTILATION":
		return ApplicationId_VENTILATION
	case "IRRIGATION_CONTROL":
		return ApplicationId_IRRIGATION_CONTROL
	case "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL":
		return ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL
	case "HEATING":
		return ApplicationId_HEATING
	case "AIR_CONDITIONING":
		return ApplicationId_AIR_CONDITIONING
	case "TRIGGER_CONTROL":
		return ApplicationId_TRIGGER_CONTROL
	case "ENABLE_CONTROL":
		return ApplicationId_ENABLE_CONTROL
	case "AUDIO_AND_VIDEO":
		return ApplicationId_AUDIO_AND_VIDEO
	case "SECURITY":
		return ApplicationId_SECURITY
	case "METERING":
		return ApplicationId_METERING
	case "ACCESS_CONTROL":
		return ApplicationId_ACCESS_CONTROL
	case "CLOCK_AND_TIMEKEEPING":
		return ApplicationId_CLOCK_AND_TIMEKEEPING
	case "TELEPHONY_STATUS_AND_CONTROL":
		return ApplicationId_TELEPHONY_STATUS_AND_CONTROL
	case "MEASUREMENT":
		return ApplicationId_MEASUREMENT
	case "TESTING":
		return ApplicationId_TESTING
	}
	return 0
}

func ApplicationIdKnows(value uint8) bool {
	for _, typeValue := range ApplicationIdValues {
		if uint8(typeValue) == value {
			return true
		}
	}
	return false
}

func CastApplicationId(structType interface{}) ApplicationId {
	castFunc := func(typ interface{}) ApplicationId {
		if sApplicationId, ok := typ.(ApplicationId); ok {
			return sApplicationId
		}
		return 0
	}
	return castFunc(structType)
}

func (m ApplicationId) GetLengthInBits() uint16 {
	return 8
}

func (m ApplicationId) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ApplicationIdParse(readBuffer utils.ReadBuffer) (ApplicationId, error) {
	val, err := readBuffer.ReadUint8("ApplicationId", 8)
	if err != nil {
		return 0, nil
	}
	return ApplicationIdByValue(val), nil
}

func (e ApplicationId) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteUint8("ApplicationId", 8, uint8(e), utils.WithAdditionalStringRepresentation(e.name()))
}

func (e ApplicationId) name() string {
	switch e {
	case ApplicationId_RESERVED:
		return "RESERVED"
	case ApplicationId_FREE_USAGE:
		return "FREE_USAGE"
	case ApplicationId_TEMPERATURE_BROADCAST:
		return "TEMPERATURE_BROADCAST"
	case ApplicationId_ROOM_CONTROL_SYSTEM:
		return "ROOM_CONTROL_SYSTEM"
	case ApplicationId_LIGHTING:
		return "LIGHTING"
	case ApplicationId_VENTILATION:
		return "VENTILATION"
	case ApplicationId_IRRIGATION_CONTROL:
		return "IRRIGATION_CONTROL"
	case ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL:
		return "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL"
	case ApplicationId_HEATING:
		return "HEATING"
	case ApplicationId_AIR_CONDITIONING:
		return "AIR_CONDITIONING"
	case ApplicationId_TRIGGER_CONTROL:
		return "TRIGGER_CONTROL"
	case ApplicationId_ENABLE_CONTROL:
		return "ENABLE_CONTROL"
	case ApplicationId_AUDIO_AND_VIDEO:
		return "AUDIO_AND_VIDEO"
	case ApplicationId_SECURITY:
		return "SECURITY"
	case ApplicationId_METERING:
		return "METERING"
	case ApplicationId_ACCESS_CONTROL:
		return "ACCESS_CONTROL"
	case ApplicationId_CLOCK_AND_TIMEKEEPING:
		return "CLOCK_AND_TIMEKEEPING"
	case ApplicationId_TELEPHONY_STATUS_AND_CONTROL:
		return "TELEPHONY_STATUS_AND_CONTROL"
	case ApplicationId_MEASUREMENT:
		return "MEASUREMENT"
	case ApplicationId_TESTING:
		return "TESTING"
	}
	return ""
}

func (e ApplicationId) String() string {
	return e.name()
}
