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

type ExtendedAPCI uint8

type IExtendedAPCI interface {
    Serialize(io utils.WriteBuffer) error
}

const(
    ExtendedAPCI_OPEN_ROUTING_TABLE_REQUEST_PDU ExtendedAPCI = 0x00
    ExtendedAPCI_READ_ROUTING_TABLE_REQUEST_PDU ExtendedAPCI = 0x01
    ExtendedAPCI_READ_ROUTING_TABLE_RESPONSE_PDU ExtendedAPCI = 0x02
    ExtendedAPCI_WRITE_ROUTING_TABLE_REQUEST_PDU ExtendedAPCI = 0x03
    ExtendedAPCI_READ_ROUTER_MEMORY_REQUEST_PDU ExtendedAPCI = 0x08
    ExtendedAPCI_READ_ROUTER_MEMORY_RESPONSE_PDU ExtendedAPCI = 0x09
    ExtendedAPCI_WRITE_ROUTER_MEMORY_REQUEST_PDU ExtendedAPCI = 0x0A
    ExtendedAPCI_READ_ROUTER_STATUS_REQUEST_PDU ExtendedAPCI = 0x0D
    ExtendedAPCI_READ_ROUTER_STATUS_RESPONSE_PDU ExtendedAPCI = 0x0E
    ExtendedAPCI_WRITE_ROUTER_STATUS_REQUEST_PDU ExtendedAPCI = 0x0F
    ExtendedAPCI_MEMORY_BIT_WRITE_PDU ExtendedAPCI = 0x10
    ExtendedAPCI_AUTHORIZE_REQUEST_PDU ExtendedAPCI = 0x11
    ExtendedAPCI_AUTHORIZE_RESPONSE_PDU ExtendedAPCI = 0x12
    ExtendedAPCI_KEY_WRITE_PDU ExtendedAPCI = 0x13
    ExtendedAPCI_KEY_RESPONSE_PDU ExtendedAPCI = 0x14
    ExtendedAPCI_PROPERTY_VALUE_READ_PDU ExtendedAPCI = 0x15
    ExtendedAPCI_PROPERTY_VALUE_RESPONSE_PDU ExtendedAPCI = 0x16
    ExtendedAPCI_PROPERTY_VALUE_WRITE_PDU ExtendedAPCI = 0x17
    ExtendedAPCI_PROPERTY_DESCRIPTION_READ_PDU ExtendedAPCI = 0x18
    ExtendedAPCI_PROPERTY_DESCRIPTION_RESPONSE_PDU ExtendedAPCI = 0x19
    ExtendedAPCI_NETWORK_PARAMETER_READ_PDU ExtendedAPCI = 0x1A
    ExtendedAPCI_NETWORK_PARAMETER_RESPONSE_PDU ExtendedAPCI = 0x1B
    ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_READ_PDU ExtendedAPCI = 0x1C
    ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_RESPONSE_PDU ExtendedAPCI = 0x1D
    ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_WRITE_PDU ExtendedAPCI = 0x1E
    ExtendedAPCI_DOMAIN_ADDRESS_WRITE ExtendedAPCI = 0x20
    ExtendedAPCI_DOMAIN_ADDRESS_READ ExtendedAPCI = 0x21
    ExtendedAPCI_DOMAIN_ADDRESS_RESPONSE ExtendedAPCI = 0x22
    ExtendedAPCI_DOMAIN_ADDRESS_SELECTIVE_READ ExtendedAPCI = 0x23
    ExtendedAPCI_NETWORK_PARAMETER_WRITE ExtendedAPCI = 0x24
    ExtendedAPCI_LINK_READ ExtendedAPCI = 0x25
    ExtendedAPCI_LINK_RESPONSE ExtendedAPCI = 0x26
    ExtendedAPCI_LINK_WRITE ExtendedAPCI = 0x27
    ExtendedAPCI_GROUP_PROPERTY_VALUE_READ ExtendedAPCI = 0x28
    ExtendedAPCI_GROUP_PROPERTY_VALUE_RESPONSE ExtendedAPCI = 0x29
    ExtendedAPCI_GROUP_PROPERTY_VALUE_WRITE ExtendedAPCI = 0x2A
    ExtendedAPCI_GROUP_PROPERTY_VALUE_INFO_REPORT ExtendedAPCI = 0x2B
    ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_READ ExtendedAPCI = 0x2C
    ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_RESPONSE ExtendedAPCI = 0x2D
    ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_WRITE ExtendedAPCI = 0x2E
    ExtendedAPCI_FILE_STREAM_INFO_REPORT ExtendedAPCI = 0x30
)

func ExtendedAPCIValueOf(value uint8) ExtendedAPCI {
    switch value {
        case 0x00:
            return ExtendedAPCI_OPEN_ROUTING_TABLE_REQUEST_PDU
        case 0x01:
            return ExtendedAPCI_READ_ROUTING_TABLE_REQUEST_PDU
        case 0x02:
            return ExtendedAPCI_READ_ROUTING_TABLE_RESPONSE_PDU
        case 0x03:
            return ExtendedAPCI_WRITE_ROUTING_TABLE_REQUEST_PDU
        case 0x08:
            return ExtendedAPCI_READ_ROUTER_MEMORY_REQUEST_PDU
        case 0x09:
            return ExtendedAPCI_READ_ROUTER_MEMORY_RESPONSE_PDU
        case 0x0A:
            return ExtendedAPCI_WRITE_ROUTER_MEMORY_REQUEST_PDU
        case 0x0D:
            return ExtendedAPCI_READ_ROUTER_STATUS_REQUEST_PDU
        case 0x0E:
            return ExtendedAPCI_READ_ROUTER_STATUS_RESPONSE_PDU
        case 0x0F:
            return ExtendedAPCI_WRITE_ROUTER_STATUS_REQUEST_PDU
        case 0x10:
            return ExtendedAPCI_MEMORY_BIT_WRITE_PDU
        case 0x11:
            return ExtendedAPCI_AUTHORIZE_REQUEST_PDU
        case 0x12:
            return ExtendedAPCI_AUTHORIZE_RESPONSE_PDU
        case 0x13:
            return ExtendedAPCI_KEY_WRITE_PDU
        case 0x14:
            return ExtendedAPCI_KEY_RESPONSE_PDU
        case 0x15:
            return ExtendedAPCI_PROPERTY_VALUE_READ_PDU
        case 0x16:
            return ExtendedAPCI_PROPERTY_VALUE_RESPONSE_PDU
        case 0x17:
            return ExtendedAPCI_PROPERTY_VALUE_WRITE_PDU
        case 0x18:
            return ExtendedAPCI_PROPERTY_DESCRIPTION_READ_PDU
        case 0x19:
            return ExtendedAPCI_PROPERTY_DESCRIPTION_RESPONSE_PDU
        case 0x1A:
            return ExtendedAPCI_NETWORK_PARAMETER_READ_PDU
        case 0x1B:
            return ExtendedAPCI_NETWORK_PARAMETER_RESPONSE_PDU
        case 0x1C:
            return ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_READ_PDU
        case 0x1D:
            return ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_RESPONSE_PDU
        case 0x1E:
            return ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_WRITE_PDU
        case 0x20:
            return ExtendedAPCI_DOMAIN_ADDRESS_WRITE
        case 0x21:
            return ExtendedAPCI_DOMAIN_ADDRESS_READ
        case 0x22:
            return ExtendedAPCI_DOMAIN_ADDRESS_RESPONSE
        case 0x23:
            return ExtendedAPCI_DOMAIN_ADDRESS_SELECTIVE_READ
        case 0x24:
            return ExtendedAPCI_NETWORK_PARAMETER_WRITE
        case 0x25:
            return ExtendedAPCI_LINK_READ
        case 0x26:
            return ExtendedAPCI_LINK_RESPONSE
        case 0x27:
            return ExtendedAPCI_LINK_WRITE
        case 0x28:
            return ExtendedAPCI_GROUP_PROPERTY_VALUE_READ
        case 0x29:
            return ExtendedAPCI_GROUP_PROPERTY_VALUE_RESPONSE
        case 0x2A:
            return ExtendedAPCI_GROUP_PROPERTY_VALUE_WRITE
        case 0x2B:
            return ExtendedAPCI_GROUP_PROPERTY_VALUE_INFO_REPORT
        case 0x2C:
            return ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_READ
        case 0x2D:
            return ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_RESPONSE
        case 0x2E:
            return ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_WRITE
        case 0x30:
            return ExtendedAPCI_FILE_STREAM_INFO_REPORT
    }
    return 0
}

func CastExtendedAPCI(structType interface{}) ExtendedAPCI {
    castFunc := func(typ interface{}) ExtendedAPCI {
        if sExtendedAPCI, ok := typ.(ExtendedAPCI); ok {
            return sExtendedAPCI
        }
        return 0
    }
    return castFunc(structType)
}

func (m ExtendedAPCI) LengthInBits() uint16 {
    return 6
}

func (m ExtendedAPCI) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ExtendedAPCIParse(io *utils.ReadBuffer) (ExtendedAPCI, error) {
    val, err := io.ReadUint8(6)
    if err != nil {
        return 0, nil
    }
    return ExtendedAPCIValueOf(val), nil
}

func (e ExtendedAPCI) Serialize(io utils.WriteBuffer) error {
    err := io.WriteUint8(6, uint8(e))
    return err
}

func (e ExtendedAPCI) String() string {
    switch e {
    case ExtendedAPCI_OPEN_ROUTING_TABLE_REQUEST_PDU:
        return "OPEN_ROUTING_TABLE_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTING_TABLE_REQUEST_PDU:
        return "READ_ROUTING_TABLE_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTING_TABLE_RESPONSE_PDU:
        return "READ_ROUTING_TABLE_RESPONSE_PDU"
    case ExtendedAPCI_WRITE_ROUTING_TABLE_REQUEST_PDU:
        return "WRITE_ROUTING_TABLE_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTER_MEMORY_REQUEST_PDU:
        return "READ_ROUTER_MEMORY_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTER_MEMORY_RESPONSE_PDU:
        return "READ_ROUTER_MEMORY_RESPONSE_PDU"
    case ExtendedAPCI_WRITE_ROUTER_MEMORY_REQUEST_PDU:
        return "WRITE_ROUTER_MEMORY_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTER_STATUS_REQUEST_PDU:
        return "READ_ROUTER_STATUS_REQUEST_PDU"
    case ExtendedAPCI_READ_ROUTER_STATUS_RESPONSE_PDU:
        return "READ_ROUTER_STATUS_RESPONSE_PDU"
    case ExtendedAPCI_WRITE_ROUTER_STATUS_REQUEST_PDU:
        return "WRITE_ROUTER_STATUS_REQUEST_PDU"
    case ExtendedAPCI_MEMORY_BIT_WRITE_PDU:
        return "MEMORY_BIT_WRITE_PDU"
    case ExtendedAPCI_AUTHORIZE_REQUEST_PDU:
        return "AUTHORIZE_REQUEST_PDU"
    case ExtendedAPCI_AUTHORIZE_RESPONSE_PDU:
        return "AUTHORIZE_RESPONSE_PDU"
    case ExtendedAPCI_KEY_WRITE_PDU:
        return "KEY_WRITE_PDU"
    case ExtendedAPCI_KEY_RESPONSE_PDU:
        return "KEY_RESPONSE_PDU"
    case ExtendedAPCI_PROPERTY_VALUE_READ_PDU:
        return "PROPERTY_VALUE_READ_PDU"
    case ExtendedAPCI_PROPERTY_VALUE_RESPONSE_PDU:
        return "PROPERTY_VALUE_RESPONSE_PDU"
    case ExtendedAPCI_PROPERTY_VALUE_WRITE_PDU:
        return "PROPERTY_VALUE_WRITE_PDU"
    case ExtendedAPCI_PROPERTY_DESCRIPTION_READ_PDU:
        return "PROPERTY_DESCRIPTION_READ_PDU"
    case ExtendedAPCI_PROPERTY_DESCRIPTION_RESPONSE_PDU:
        return "PROPERTY_DESCRIPTION_RESPONSE_PDU"
    case ExtendedAPCI_NETWORK_PARAMETER_READ_PDU:
        return "NETWORK_PARAMETER_READ_PDU"
    case ExtendedAPCI_NETWORK_PARAMETER_RESPONSE_PDU:
        return "NETWORK_PARAMETER_RESPONSE_PDU"
    case ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_READ_PDU:
        return "INDIVIDUAL_ADDRESS_SERIAL_NUMBER_READ_PDU"
    case ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_RESPONSE_PDU:
        return "INDIVIDUAL_ADDRESS_SERIAL_NUMBER_RESPONSE_PDU"
    case ExtendedAPCI_INDIVIDUAL_ADDRESS_SERIAL_NUMBER_WRITE_PDU:
        return "INDIVIDUAL_ADDRESS_SERIAL_NUMBER_WRITE_PDU"
    case ExtendedAPCI_DOMAIN_ADDRESS_WRITE:
        return "DOMAIN_ADDRESS_WRITE"
    case ExtendedAPCI_DOMAIN_ADDRESS_READ:
        return "DOMAIN_ADDRESS_READ"
    case ExtendedAPCI_DOMAIN_ADDRESS_RESPONSE:
        return "DOMAIN_ADDRESS_RESPONSE"
    case ExtendedAPCI_DOMAIN_ADDRESS_SELECTIVE_READ:
        return "DOMAIN_ADDRESS_SELECTIVE_READ"
    case ExtendedAPCI_NETWORK_PARAMETER_WRITE:
        return "NETWORK_PARAMETER_WRITE"
    case ExtendedAPCI_LINK_READ:
        return "LINK_READ"
    case ExtendedAPCI_LINK_RESPONSE:
        return "LINK_RESPONSE"
    case ExtendedAPCI_LINK_WRITE:
        return "LINK_WRITE"
    case ExtendedAPCI_GROUP_PROPERTY_VALUE_READ:
        return "GROUP_PROPERTY_VALUE_READ"
    case ExtendedAPCI_GROUP_PROPERTY_VALUE_RESPONSE:
        return "GROUP_PROPERTY_VALUE_RESPONSE"
    case ExtendedAPCI_GROUP_PROPERTY_VALUE_WRITE:
        return "GROUP_PROPERTY_VALUE_WRITE"
    case ExtendedAPCI_GROUP_PROPERTY_VALUE_INFO_REPORT:
        return "GROUP_PROPERTY_VALUE_INFO_REPORT"
    case ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_READ:
        return "DOMAIN_ADDRESS_SERIAL_NUMBER_READ"
    case ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_RESPONSE:
        return "DOMAIN_ADDRESS_SERIAL_NUMBER_RESPONSE"
    case ExtendedAPCI_DOMAIN_ADDRESS_SERIAL_NUMBER_WRITE:
        return "DOMAIN_ADDRESS_SERIAL_NUMBER_WRITE"
    case ExtendedAPCI_FILE_STREAM_INFO_REPORT:
        return "FILE_STREAM_INFO_REPORT"
    }
    return ""
}
