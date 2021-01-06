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

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
)

type KnxInterfaceObjectType uint16

type IKnxInterfaceObjectType interface {
	Text() string
	Serialize(io utils.WriteBuffer) error
}

const (
	KnxInterfaceObjectType_OT_GENERAL                        KnxInterfaceObjectType = 65535
	KnxInterfaceObjectType_OT_DEVICE                         KnxInterfaceObjectType = 0
	KnxInterfaceObjectType_OT_ADDRESS_TABLE                  KnxInterfaceObjectType = 1
	KnxInterfaceObjectType_OT_ASSOCIATION_TABLE              KnxInterfaceObjectType = 2
	KnxInterfaceObjectType_OT_APPLICATION_PROGRAM            KnxInterfaceObjectType = 3
	KnxInterfaceObjectType_OT_INTERACE_PROGRAM               KnxInterfaceObjectType = 4
	KnxInterfaceObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE  KnxInterfaceObjectType = 5
	KnxInterfaceObjectType_OT_ROUTER                         KnxInterfaceObjectType = 6
	KnxInterfaceObjectType_OT_LTE_ADDRESS_ROUTING_TABLE      KnxInterfaceObjectType = 7
	KnxInterfaceObjectType_OT_CEMI_SERVER                    KnxInterfaceObjectType = 8
	KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE             KnxInterfaceObjectType = 9
	KnxInterfaceObjectType_OT_POLLING_MASTER                 KnxInterfaceObjectType = 10
	KnxInterfaceObjectType_OT_KNXIP_PARAMETER                KnxInterfaceObjectType = 11
	KnxInterfaceObjectType_OT_FILE_SERVER                    KnxInterfaceObjectType = 13
	KnxInterfaceObjectType_OT_SECURITY                       KnxInterfaceObjectType = 17
	KnxInterfaceObjectType_OT_RF_MEDIUM                      KnxInterfaceObjectType = 19
	KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR       KnxInterfaceObjectType = 409
	KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR        KnxInterfaceObjectType = 410
	KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC KnxInterfaceObjectType = 417
	KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC         KnxInterfaceObjectType = 418
	KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC           KnxInterfaceObjectType = 420
	KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC         KnxInterfaceObjectType = 421
	KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC        KnxInterfaceObjectType = 800
	KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC          KnxInterfaceObjectType = 801
)

func (e KnxInterfaceObjectType) Text() string {
	switch e {
	case 0:
		{ /* '0' */
			return "Device Object"
		}
	case 1:
		{ /* '1' */
			return "Addresstable Object"
		}
	case 10:
		{ /* '10' */
			return "Polling Master"
		}
	case 11:
		{ /* '11' */
			return "KNXnet/IP Parameter Object"
		}
	case 13:
		{ /* '13' */
			return "File Server Object"
		}
	case 17:
		{ /* '17' */
			return "Security Object"
		}
	case 19:
		{ /* '19' */
			return "RF Medium Object"
		}
	case 2:
		{ /* '2' */
			return "Associationtable Object"
		}
	case 3:
		{ /* '3' */
			return "Applicationprogram Object"
		}
	case 4:
		{ /* '4' */
			return "Interfaceprogram Object"
		}
	case 409:
		{ /* '409' */
			return "Indoor Brightness Sensor"
		}
	case 410:
		{ /* '410' */
			return "Indoor Luminance Sensor"
		}
	case 417:
		{ /* '417' */
			return "Light Switching Actuator Basic"
		}
	case 418:
		{ /* '418' */
			return "Dimming Actuator Basic"
		}
	case 420:
		{ /* '420' */
			return "Dimming   Sensor Basic"
		}
	case 421:
		{ /* '421' */
			return "Switching Sensor Basic"
		}
	case 5:
		{ /* '5' */
			return "KNX-Object Associationtable Object"
		}
	case 6:
		{ /* '6' */
			return "Router Object"
		}
	case 65535:
		{ /* '65535' */
			return "General Object"
		}
	case 7:
		{ /* '7' */
			return "LTE Address Routing Table Object"
		}
	case 8:
		{ /* '8' */
			return "cEMI Server Object"
		}
	case 800:
		{ /* '800' */
			return "Sunblind Actuator Basic"
		}
	case 801:
		{ /* '801' */
			return "Sunblind Sensor Basic"
		}
	case 9:
		{ /* '9' */
			return "Group Object Table Object"
		}
	default:
		{
			return ""
		}
	}
}
func KnxInterfaceObjectTypeByValue(value uint16) KnxInterfaceObjectType {
	switch value {
	case 0:
		return KnxInterfaceObjectType_OT_DEVICE
	case 1:
		return KnxInterfaceObjectType_OT_ADDRESS_TABLE
	case 10:
		return KnxInterfaceObjectType_OT_POLLING_MASTER
	case 11:
		return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
	case 13:
		return KnxInterfaceObjectType_OT_FILE_SERVER
	case 17:
		return KnxInterfaceObjectType_OT_SECURITY
	case 19:
		return KnxInterfaceObjectType_OT_RF_MEDIUM
	case 2:
		return KnxInterfaceObjectType_OT_ASSOCIATION_TABLE
	case 3:
		return KnxInterfaceObjectType_OT_APPLICATION_PROGRAM
	case 4:
		return KnxInterfaceObjectType_OT_INTERACE_PROGRAM
	case 409:
		return KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR
	case 410:
		return KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR
	case 417:
		return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
	case 418:
		return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
	case 420:
		return KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC
	case 421:
		return KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC
	case 5:
		return KnxInterfaceObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE
	case 6:
		return KnxInterfaceObjectType_OT_ROUTER
	case 65535:
		return KnxInterfaceObjectType_OT_GENERAL
	case 7:
		return KnxInterfaceObjectType_OT_LTE_ADDRESS_ROUTING_TABLE
	case 8:
		return KnxInterfaceObjectType_OT_CEMI_SERVER
	case 800:
		return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
	case 801:
		return KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC
	case 9:
		return KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE
	}
	return 0
}

func KnxInterfaceObjectTypeByName(value string) KnxInterfaceObjectType {
	switch value {
	case "OT_DEVICE":
		return KnxInterfaceObjectType_OT_DEVICE
	case "OT_ADDRESS_TABLE":
		return KnxInterfaceObjectType_OT_ADDRESS_TABLE
	case "OT_POLLING_MASTER":
		return KnxInterfaceObjectType_OT_POLLING_MASTER
	case "OT_KNXIP_PARAMETER":
		return KnxInterfaceObjectType_OT_KNXIP_PARAMETER
	case "OT_FILE_SERVER":
		return KnxInterfaceObjectType_OT_FILE_SERVER
	case "OT_SECURITY":
		return KnxInterfaceObjectType_OT_SECURITY
	case "OT_RF_MEDIUM":
		return KnxInterfaceObjectType_OT_RF_MEDIUM
	case "OT_ASSOCIATION_TABLE":
		return KnxInterfaceObjectType_OT_ASSOCIATION_TABLE
	case "OT_APPLICATION_PROGRAM":
		return KnxInterfaceObjectType_OT_APPLICATION_PROGRAM
	case "OT_INTERACE_PROGRAM":
		return KnxInterfaceObjectType_OT_INTERACE_PROGRAM
	case "OT_INDOOR_BRIGHTNESS_SENSOR":
		return KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR
	case "OT_INDOOR_LUMINANCE_SENSOR":
		return KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR
	case "OT_LIGHT_SWITCHING_ACTUATOR_BASIC":
		return KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
	case "OT_DIMMING_ACTUATOR_BASIC":
		return KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC
	case "OT_DIMMING_SENSOR_BASIC":
		return KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC
	case "OT_SWITCHING_SENSOR_BASIC":
		return KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC
	case "OT_EIBOBJECT_ASSOCIATATION_TABLE":
		return KnxInterfaceObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE
	case "OT_ROUTER":
		return KnxInterfaceObjectType_OT_ROUTER
	case "OT_GENERAL":
		return KnxInterfaceObjectType_OT_GENERAL
	case "OT_LTE_ADDRESS_ROUTING_TABLE":
		return KnxInterfaceObjectType_OT_LTE_ADDRESS_ROUTING_TABLE
	case "OT_CEMI_SERVER":
		return KnxInterfaceObjectType_OT_CEMI_SERVER
	case "OT_SUNBLIND_ACTUATOR_BASIC":
		return KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC
	case "OT_SUNBLIND_SENSOR_BASIC":
		return KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC
	case "OT_GROUP_OBJECT_TABLE":
		return KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE
	}
	return 0
}

func CastKnxInterfaceObjectType(structType interface{}) KnxInterfaceObjectType {
	castFunc := func(typ interface{}) KnxInterfaceObjectType {
		if sKnxInterfaceObjectType, ok := typ.(KnxInterfaceObjectType); ok {
			return sKnxInterfaceObjectType
		}
		return 0
	}
	return castFunc(structType)
}

func (m KnxInterfaceObjectType) LengthInBits() uint16 {
	return 16
}

func (m KnxInterfaceObjectType) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxInterfaceObjectTypeParse(io *utils.ReadBuffer) (KnxInterfaceObjectType, error) {
	val, err := io.ReadUint16(16)
	if err != nil {
		return 0, nil
	}
	return KnxInterfaceObjectTypeByValue(val), nil
}

func (e KnxInterfaceObjectType) Serialize(io utils.WriteBuffer) error {
	err := io.WriteUint16(16, uint16(e))
	return err
}

func (e KnxInterfaceObjectType) String() string {
	switch e {
	case KnxInterfaceObjectType_OT_DEVICE:
		return "OT_DEVICE"
	case KnxInterfaceObjectType_OT_ADDRESS_TABLE:
		return "OT_ADDRESS_TABLE"
	case KnxInterfaceObjectType_OT_POLLING_MASTER:
		return "OT_POLLING_MASTER"
	case KnxInterfaceObjectType_OT_KNXIP_PARAMETER:
		return "OT_KNXIP_PARAMETER"
	case KnxInterfaceObjectType_OT_FILE_SERVER:
		return "OT_FILE_SERVER"
	case KnxInterfaceObjectType_OT_SECURITY:
		return "OT_SECURITY"
	case KnxInterfaceObjectType_OT_RF_MEDIUM:
		return "OT_RF_MEDIUM"
	case KnxInterfaceObjectType_OT_ASSOCIATION_TABLE:
		return "OT_ASSOCIATION_TABLE"
	case KnxInterfaceObjectType_OT_APPLICATION_PROGRAM:
		return "OT_APPLICATION_PROGRAM"
	case KnxInterfaceObjectType_OT_INTERACE_PROGRAM:
		return "OT_INTERACE_PROGRAM"
	case KnxInterfaceObjectType_OT_INDOOR_BRIGHTNESS_SENSOR:
		return "OT_INDOOR_BRIGHTNESS_SENSOR"
	case KnxInterfaceObjectType_OT_INDOOR_LUMINANCE_SENSOR:
		return "OT_INDOOR_LUMINANCE_SENSOR"
	case KnxInterfaceObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC:
		return "OT_LIGHT_SWITCHING_ACTUATOR_BASIC"
	case KnxInterfaceObjectType_OT_DIMMING_ACTUATOR_BASIC:
		return "OT_DIMMING_ACTUATOR_BASIC"
	case KnxInterfaceObjectType_OT_DIMMING_SENSOR_BASIC:
		return "OT_DIMMING_SENSOR_BASIC"
	case KnxInterfaceObjectType_OT_SWITCHING_SENSOR_BASIC:
		return "OT_SWITCHING_SENSOR_BASIC"
	case KnxInterfaceObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE:
		return "OT_EIBOBJECT_ASSOCIATATION_TABLE"
	case KnxInterfaceObjectType_OT_ROUTER:
		return "OT_ROUTER"
	case KnxInterfaceObjectType_OT_GENERAL:
		return "OT_GENERAL"
	case KnxInterfaceObjectType_OT_LTE_ADDRESS_ROUTING_TABLE:
		return "OT_LTE_ADDRESS_ROUTING_TABLE"
	case KnxInterfaceObjectType_OT_CEMI_SERVER:
		return "OT_CEMI_SERVER"
	case KnxInterfaceObjectType_OT_SUNBLIND_ACTUATOR_BASIC:
		return "OT_SUNBLIND_ACTUATOR_BASIC"
	case KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC:
		return "OT_SUNBLIND_SENSOR_BASIC"
	case KnxInterfaceObjectType_OT_GROUP_OBJECT_TABLE:
		return "OT_GROUP_OBJECT_TABLE"
	}
	return ""
}
