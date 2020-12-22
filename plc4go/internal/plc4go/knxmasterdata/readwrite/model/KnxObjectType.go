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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

type KnxObjectType uint8

type IKnxObjectType interface {
    Text() string
    Serialize(io utils.WriteBuffer) error
}

const(
    KnxObjectType_OT_DEVICE KnxObjectType = 0
    KnxObjectType_OT_ADDRESS_TABLE KnxObjectType = 1
    KnxObjectType_OT_ASSOCIATION_TABLE KnxObjectType = 2
    KnxObjectType_OT_APPLICATION_PROGRAM KnxObjectType = 3
    KnxObjectType_OT_INTERACE_PROGRAM KnxObjectType = 4
    KnxObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE KnxObjectType = 5
    KnxObjectType_OT_ROUTER KnxObjectType = 6
    KnxObjectType_OT_LTE_ADDRESS_ROUTING_TABLE KnxObjectType = 7
    KnxObjectType_OT_CEMI_SERVER KnxObjectType = 8
    KnxObjectType_OT_GROUP_OBJECT_TABLE KnxObjectType = 9
    KnxObjectType_OT_POLLING_MASTER KnxObjectType = 10
    KnxObjectType_OT_KNXIP_PARAMETER KnxObjectType = 11
    KnxObjectType_OT_FILE_SERVER KnxObjectType = 13
    KnxObjectType_OT_SECURITY KnxObjectType = 17
    KnxObjectType_OT_RF_MEDIUM KnxObjectType = 19
    KnxObjectType_OT_INDOOR_BRIGHTNESS_SENSOR KnxObjectType = 409
    KnxObjectType_OT_INDOOR_LUMINANCE_SENSOR KnxObjectType = 410
    KnxObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC KnxObjectType = 417
    KnxObjectType_OT_DIMMING_ACTUATOR_BASIC KnxObjectType = 418
    KnxObjectType_OT_DIMMING_SENSOR_BASIC KnxObjectType = 420
    KnxObjectType_OT_SWITCHING_SENSOR_BASIC KnxObjectType = 421
    KnxObjectType_OT_SUNBLIND_ACTUATOR_BASIC KnxObjectType = 800
    KnxObjectType_OT_SUNBLIND_SENSOR_BASIC KnxObjectType = 801
)


func (e KnxObjectType) Text() string {
    switch e  {
        case 0: { /* '0' */
            return "Device Object"
        }
        case 1: { /* '1' */
            return "Addresstable Object"
        }
        case 10: { /* '10' */
            return "Polling Master"
        }
        case 11: { /* '11' */
            return "KNXnet/IP Parameter Object"
        }
        case 13: { /* '13' */
            return "File Server Object"
        }
        case 17: { /* '17' */
            return "Security Object"
        }
        case 19: { /* '19' */
            return "RF Medium Object"
        }
        case 2: { /* '2' */
            return "Associationtable Object"
        }
        case 3: { /* '3' */
            return "Applicationprogram Object"
        }
        case 4: { /* '4' */
            return "Interfaceprogram Object"
        }
        case 409: { /* '409' */
            return "Indoor Brightness Sensor"
        }
        case 410: { /* '410' */
            return "Indoor Luminance Sensor"
        }
        case 417: { /* '417' */
            return "Light Switching Actuator Basic"
        }
        case 418: { /* '418' */
            return "Dimming Actuator Basic"
        }
        case 420: { /* '420' */
            return "Dimming   Sensor Basic"
        }
        case 421: { /* '421' */
            return "Switching Sensor Basic"
        }
        case 5: { /* '5' */
            return "KNX-Object Associationtable Object"
        }
        case 6: { /* '6' */
            return "Router Object"
        }
        case 7: { /* '7' */
            return "LTE Address Routing Table Object"
        }
        case 8: { /* '8' */
            return "cEMI Server Object"
        }
        case 800: { /* '800' */
            return "Sunblind Actuator Basic"
        }
        case 801: { /* '801' */
            return "Sunblind Sensor Basic"
        }
        case 9: { /* '9' */
            return "Group Object Table Object"
        }
        default: {
            return ""
        }
    }
}
func KnxObjectTypeValueOf(value uint8) KnxObjectType {
    switch value {
        case 0:
            return KnxObjectType_OT_DEVICE
        case 1:
            return KnxObjectType_OT_ADDRESS_TABLE
        case 10:
            return KnxObjectType_OT_POLLING_MASTER
        case 11:
            return KnxObjectType_OT_KNXIP_PARAMETER
        case 13:
            return KnxObjectType_OT_FILE_SERVER
        case 17:
            return KnxObjectType_OT_SECURITY
        case 19:
            return KnxObjectType_OT_RF_MEDIUM
        case 2:
            return KnxObjectType_OT_ASSOCIATION_TABLE
        case 3:
            return KnxObjectType_OT_APPLICATION_PROGRAM
        case 4:
            return KnxObjectType_OT_INTERACE_PROGRAM
        case 409:
            return KnxObjectType_OT_INDOOR_BRIGHTNESS_SENSOR
        case 410:
            return KnxObjectType_OT_INDOOR_LUMINANCE_SENSOR
        case 417:
            return KnxObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC
        case 418:
            return KnxObjectType_OT_DIMMING_ACTUATOR_BASIC
        case 420:
            return KnxObjectType_OT_DIMMING_SENSOR_BASIC
        case 421:
            return KnxObjectType_OT_SWITCHING_SENSOR_BASIC
        case 5:
            return KnxObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE
        case 6:
            return KnxObjectType_OT_ROUTER
        case 7:
            return KnxObjectType_OT_LTE_ADDRESS_ROUTING_TABLE
        case 8:
            return KnxObjectType_OT_CEMI_SERVER
        case 800:
            return KnxObjectType_OT_SUNBLIND_ACTUATOR_BASIC
        case 801:
            return KnxObjectType_OT_SUNBLIND_SENSOR_BASIC
        case 9:
            return KnxObjectType_OT_GROUP_OBJECT_TABLE
    }
    return 0
}

func CastKnxObjectType(structType interface{}) KnxObjectType {
    castFunc := func(typ interface{}) KnxObjectType {
        if sKnxObjectType, ok := typ.(KnxObjectType); ok {
            return sKnxObjectType
        }
        return 0
    }
    return castFunc(structType)
}

func (m KnxObjectType) LengthInBits() uint16 {
    return 8
}

func (m KnxObjectType) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxObjectTypeParse(io *utils.ReadBuffer) (KnxObjectType, error) {
    val, err := io.ReadUint8(8)
    if err != nil {
        return 0, nil
    }
    return KnxObjectTypeValueOf(val), nil
}

func (e KnxObjectType) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(8, uint8(e))
    return err
}

func (e KnxObjectType) String() string {
    switch e {
    case KnxObjectType_OT_DEVICE:
        return "OT_DEVICE"
    case KnxObjectType_OT_ADDRESS_TABLE:
        return "OT_ADDRESS_TABLE"
    case KnxObjectType_OT_POLLING_MASTER:
        return "OT_POLLING_MASTER"
    case KnxObjectType_OT_KNXIP_PARAMETER:
        return "OT_KNXIP_PARAMETER"
    case KnxObjectType_OT_FILE_SERVER:
        return "OT_FILE_SERVER"
    case KnxObjectType_OT_SECURITY:
        return "OT_SECURITY"
    case KnxObjectType_OT_RF_MEDIUM:
        return "OT_RF_MEDIUM"
    case KnxObjectType_OT_ASSOCIATION_TABLE:
        return "OT_ASSOCIATION_TABLE"
    case KnxObjectType_OT_APPLICATION_PROGRAM:
        return "OT_APPLICATION_PROGRAM"
    case KnxObjectType_OT_INTERACE_PROGRAM:
        return "OT_INTERACE_PROGRAM"
    case KnxObjectType_OT_INDOOR_BRIGHTNESS_SENSOR:
        return "OT_INDOOR_BRIGHTNESS_SENSOR"
    case KnxObjectType_OT_INDOOR_LUMINANCE_SENSOR:
        return "OT_INDOOR_LUMINANCE_SENSOR"
    case KnxObjectType_OT_LIGHT_SWITCHING_ACTUATOR_BASIC:
        return "OT_LIGHT_SWITCHING_ACTUATOR_BASIC"
    case KnxObjectType_OT_DIMMING_ACTUATOR_BASIC:
        return "OT_DIMMING_ACTUATOR_BASIC"
    case KnxObjectType_OT_DIMMING_SENSOR_BASIC:
        return "OT_DIMMING_SENSOR_BASIC"
    case KnxObjectType_OT_SWITCHING_SENSOR_BASIC:
        return "OT_SWITCHING_SENSOR_BASIC"
    case KnxObjectType_OT_EIBOBJECT_ASSOCIATATION_TABLE:
        return "OT_EIBOBJECT_ASSOCIATATION_TABLE"
    case KnxObjectType_OT_ROUTER:
        return "OT_ROUTER"
    case KnxObjectType_OT_LTE_ADDRESS_ROUTING_TABLE:
        return "OT_LTE_ADDRESS_ROUTING_TABLE"
    case KnxObjectType_OT_CEMI_SERVER:
        return "OT_CEMI_SERVER"
    case KnxObjectType_OT_SUNBLIND_ACTUATOR_BASIC:
        return "OT_SUNBLIND_ACTUATOR_BASIC"
    case KnxObjectType_OT_SUNBLIND_SENSOR_BASIC:
        return "OT_SUNBLIND_SENSOR_BASIC"
    case KnxObjectType_OT_GROUP_OBJECT_TABLE:
        return "OT_GROUP_OBJECT_TABLE"
    }
    return ""
}
