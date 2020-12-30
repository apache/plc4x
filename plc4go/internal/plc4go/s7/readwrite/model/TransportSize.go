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

type TransportSize int8

type ITransportSize interface {
    Supported_S7_300() bool
    Supported_LOGO() bool
    SizeInBytes() uint8
    Supported_S7_400() bool
    Supported_S7_1200() bool
    SizeCode() uint8
    Supported_S7_1500() bool
    DataTransportSize() DataTransportSize
    BaseType() TransportSize
    DataProtocolId() string
    Serialize(io utils.WriteBuffer) error
}

const(
    TransportSize_BOOL TransportSize = 0x01
    TransportSize_BYTE TransportSize = 0x02
    TransportSize_WORD TransportSize = 0x04
    TransportSize_DWORD TransportSize = 0x06
    TransportSize_LWORD TransportSize = 0x00
    TransportSize_INT TransportSize = 0x05
    TransportSize_UINT TransportSize = 0x05
    TransportSize_SINT TransportSize = 0x02
    TransportSize_USINT TransportSize = 0x02
    TransportSize_DINT TransportSize = 0x07
    TransportSize_UDINT TransportSize = 0x07
    TransportSize_LINT TransportSize = 0x00
    TransportSize_ULINT TransportSize = 0x00
    TransportSize_REAL TransportSize = 0x08
    TransportSize_LREAL TransportSize = 0x30
    TransportSize_CHAR TransportSize = 0x03
    TransportSize_WCHAR TransportSize = 0x13
    TransportSize_STRING TransportSize = 0x03
    TransportSize_WSTRING TransportSize = 0x00
    TransportSize_TIME TransportSize = 0x0B
    TransportSize_LTIME TransportSize = 0x00
    TransportSize_DATE TransportSize = 0x09
    TransportSize_TIME_OF_DAY TransportSize = 0x06
    TransportSize_TOD TransportSize = 0x06
    TransportSize_DATE_AND_TIME TransportSize = 0x0F
    TransportSize_DT TransportSize = 0x0F
)


func (e TransportSize) Supported_S7_300() bool {
    switch e  {
        case 0x00: { /* '0x00' */
            return false
        }
        case 0x01: { /* '0x01' */
            return true
        }
        case 0x02: { /* '0x02' */
            return true
        }
        case 0x03: { /* '0x03' */
            return true
        }
        case 0x04: { /* '0x04' */
            return true
        }
        case 0x05: { /* '0x05' */
            return true
        }
        case 0x06: { /* '0x06' */
            return true
        }
        case 0x07: { /* '0x07' */
            return true
        }
        case 0x08: { /* '0x08' */
            return true
        }
        case 0x09: { /* '0x09' */
            return true
        }
        case 0x0B: { /* '0x0B' */
            return true
        }
        case 0x0F: { /* '0x0F' */
            return true
        }
        case 0x13: { /* '0x13' */
            return false
        }
        case 0x30: { /* '0x30' */
            return false
        }
        default: {
            return false
        }
    }
}

func (e TransportSize) Supported_LOGO() bool {
    switch e  {
        case 0x00: { /* '0x00' */
            return false
        }
        case 0x01: { /* '0x01' */
            return true
        }
        case 0x02: { /* '0x02' */
            return true
        }
        case 0x03: { /* '0x03' */
            return true
        }
        case 0x04: { /* '0x04' */
            return true
        }
        case 0x05: { /* '0x05' */
            return true
        }
        case 0x06: { /* '0x06' */
            return true
        }
        case 0x07: { /* '0x07' */
            return true
        }
        case 0x08: { /* '0x08' */
            return true
        }
        case 0x09: { /* '0x09' */
            return true
        }
        case 0x0B: { /* '0x0B' */
            return true
        }
        case 0x0F: { /* '0x0F' */
            return false
        }
        case 0x13: { /* '0x13' */
            return true
        }
        case 0x30: { /* '0x30' */
            return false
        }
        default: {
            return false
        }
    }
}

func (e TransportSize) SizeInBytes() uint8 {
    switch e  {
        case 0x00: { /* '0x00' */
            return 8
        }
        case 0x01: { /* '0x01' */
            return 1
        }
        case 0x02: { /* '0x02' */
            return 1
        }
        case 0x03: { /* '0x03' */
            return 1
        }
        case 0x04: { /* '0x04' */
            return 2
        }
        case 0x05: { /* '0x05' */
            return 2
        }
        case 0x06: { /* '0x06' */
            return 4
        }
        case 0x07: { /* '0x07' */
            return 4
        }
        case 0x08: { /* '0x08' */
            return 4
        }
        case 0x09: { /* '0x09' */
            return 2
        }
        case 0x0B: { /* '0x0B' */
            return 4
        }
        case 0x0F: { /* '0x0F' */
            return 12
        }
        case 0x13: { /* '0x13' */
            return 2
        }
        case 0x30: { /* '0x30' */
            return 8
        }
        default: {
            return 0
        }
    }
}

func (e TransportSize) Supported_S7_400() bool {
    switch e  {
        case 0x00: { /* '0x00' */
            return false
        }
        case 0x01: { /* '0x01' */
            return true
        }
        case 0x02: { /* '0x02' */
            return true
        }
        case 0x03: { /* '0x03' */
            return true
        }
        case 0x04: { /* '0x04' */
            return true
        }
        case 0x05: { /* '0x05' */
            return true
        }
        case 0x06: { /* '0x06' */
            return true
        }
        case 0x07: { /* '0x07' */
            return true
        }
        case 0x08: { /* '0x08' */
            return true
        }
        case 0x09: { /* '0x09' */
            return true
        }
        case 0x0B: { /* '0x0B' */
            return true
        }
        case 0x0F: { /* '0x0F' */
            return true
        }
        case 0x13: { /* '0x13' */
            return false
        }
        case 0x30: { /* '0x30' */
            return false
        }
        default: {
            return false
        }
    }
}

func (e TransportSize) Supported_S7_1200() bool {
    switch e  {
        case 0x00: { /* '0x00' */
            return false
        }
        case 0x01: { /* '0x01' */
            return true
        }
        case 0x02: { /* '0x02' */
            return true
        }
        case 0x03: { /* '0x03' */
            return true
        }
        case 0x04: { /* '0x04' */
            return true
        }
        case 0x05: { /* '0x05' */
            return true
        }
        case 0x06: { /* '0x06' */
            return true
        }
        case 0x07: { /* '0x07' */
            return true
        }
        case 0x08: { /* '0x08' */
            return true
        }
        case 0x09: { /* '0x09' */
            return true
        }
        case 0x0B: { /* '0x0B' */
            return true
        }
        case 0x0F: { /* '0x0F' */
            return false
        }
        case 0x13: { /* '0x13' */
            return true
        }
        case 0x30: { /* '0x30' */
            return true
        }
        default: {
            return false
        }
    }
}

func (e TransportSize) SizeCode() uint8 {
    switch e  {
        case 0x00: { /* '0x00' */
            return 'X'
        }
        case 0x01: { /* '0x01' */
            return 'X'
        }
        case 0x02: { /* '0x02' */
            return 'B'
        }
        case 0x03: { /* '0x03' */
            return 'B'
        }
        case 0x04: { /* '0x04' */
            return 'W'
        }
        case 0x05: { /* '0x05' */
            return 'W'
        }
        case 0x06: { /* '0x06' */
            return 'D'
        }
        case 0x07: { /* '0x07' */
            return 'D'
        }
        case 0x08: { /* '0x08' */
            return 'D'
        }
        case 0x09: { /* '0x09' */
            return 'X'
        }
        case 0x0B: { /* '0x0B' */
            return 'X'
        }
        case 0x0F: { /* '0x0F' */
            return 'X'
        }
        case 0x13: { /* '0x13' */
            return 'X'
        }
        case 0x30: { /* '0x30' */
            return 'X'
        }
        default: {
            return 0
        }
    }
}

func (e TransportSize) Supported_S7_1500() bool {
    switch e  {
        case 0x00: { /* '0x00' */
            return true
        }
        case 0x01: { /* '0x01' */
            return true
        }
        case 0x02: { /* '0x02' */
            return true
        }
        case 0x03: { /* '0x03' */
            return true
        }
        case 0x04: { /* '0x04' */
            return true
        }
        case 0x05: { /* '0x05' */
            return true
        }
        case 0x06: { /* '0x06' */
            return true
        }
        case 0x07: { /* '0x07' */
            return true
        }
        case 0x08: { /* '0x08' */
            return true
        }
        case 0x09: { /* '0x09' */
            return true
        }
        case 0x0B: { /* '0x0B' */
            return true
        }
        case 0x0F: { /* '0x0F' */
            return true
        }
        case 0x13: { /* '0x13' */
            return true
        }
        case 0x30: { /* '0x30' */
            return true
        }
        default: {
            return false
        }
    }
}

func (e TransportSize) DataTransportSize() DataTransportSize {
    switch e  {
        case 0x00: { /* '0x00' */
            return 0
        }
        case 0x01: { /* '0x01' */
            return DataTransportSize_BIT
        }
        case 0x02: { /* '0x02' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x03: { /* '0x03' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x04: { /* '0x04' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x05: { /* '0x05' */
            return DataTransportSize_INTEGER
        }
        case 0x06: { /* '0x06' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x07: { /* '0x07' */
            return DataTransportSize_INTEGER
        }
        case 0x08: { /* '0x08' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x09: { /* '0x09' */
            return DataTransportSize_BYTE_WORD_DWORD
        }
        case 0x0B: { /* '0x0B' */
            return 0
        }
        case 0x0F: { /* '0x0F' */
            return 0
        }
        case 0x13: { /* '0x13' */
            return 0
        }
        case 0x30: { /* '0x30' */
            return 0
        }
        default: {
            return 0
        }
    }
}

func (e TransportSize) BaseType() TransportSize {
    switch e  {
        case 0x00: { /* '0x00' */
            return 0
        }
        case 0x01: { /* '0x01' */
            return 0
        }
        case 0x02: { /* '0x02' */
            return 0
        }
        case 0x03: { /* '0x03' */
            return 0
        }
        case 0x04: { /* '0x04' */
            return 0
        }
        case 0x05: { /* '0x05' */
            return 0
        }
        case 0x06: { /* '0x06' */
            return TransportSize_WORD
        }
        case 0x07: { /* '0x07' */
            return TransportSize_INT
        }
        case 0x08: { /* '0x08' */
            return 0
        }
        case 0x09: { /* '0x09' */
            return 0
        }
        case 0x0B: { /* '0x0B' */
            return 0
        }
        case 0x0F: { /* '0x0F' */
            return 0
        }
        case 0x13: { /* '0x13' */
            return 0
        }
        case 0x30: { /* '0x30' */
            return TransportSize_REAL
        }
        default: {
            return 0
        }
    }
}

func (e TransportSize) DataProtocolId() string {
    switch e  {
        case 0x00: { /* '0x00' */
            return "IEC61131_LWORD"
        }
        case 0x01: { /* '0x01' */
            return "IEC61131_BOOL"
        }
        case 0x02: { /* '0x02' */
            return "IEC61131_BYTE"
        }
        case 0x03: { /* '0x03' */
            return "IEC61131_CHAR"
        }
        case 0x04: { /* '0x04' */
            return "IEC61131_WORD"
        }
        case 0x05: { /* '0x05' */
            return "IEC61131_INT"
        }
        case 0x06: { /* '0x06' */
            return "IEC61131_DWORD"
        }
        case 0x07: { /* '0x07' */
            return "IEC61131_DINT"
        }
        case 0x08: { /* '0x08' */
            return "IEC61131_REAL"
        }
        case 0x09: { /* '0x09' */
            return "IEC61131_DATE"
        }
        case 0x0B: { /* '0x0B' */
            return "IEC61131_TIME"
        }
        case 0x0F: { /* '0x0F' */
            return "IEC61131_DATE_AND_TIME"
        }
        case 0x13: { /* '0x13' */
            return "IEC61131_WCHAR"
        }
        case 0x30: { /* '0x30' */
            return "IEC61131_LREAL"
        }
        default: {
            return ""
        }
    }
}
func TransportSizeValueOf(value int8) TransportSize {
    switch value {
        case 0x00:
            return TransportSize_LWORD
        case 0x01:
            return TransportSize_BOOL
        case 0x02:
            return TransportSize_BYTE
        case 0x03:
            return TransportSize_CHAR
        case 0x04:
            return TransportSize_WORD
        case 0x05:
            return TransportSize_INT
        case 0x06:
            return TransportSize_DWORD
        case 0x07:
            return TransportSize_DINT
        case 0x08:
            return TransportSize_REAL
        case 0x09:
            return TransportSize_DATE
        case 0x0B:
            return TransportSize_TIME
        case 0x0F:
            return TransportSize_DATE_AND_TIME
        case 0x13:
            return TransportSize_WCHAR
        case 0x30:
            return TransportSize_LREAL
    }
    return 0
}

func CastTransportSize(structType interface{}) TransportSize {
    castFunc := func(typ interface{}) TransportSize {
        if sTransportSize, ok := typ.(TransportSize); ok {
            return sTransportSize
        }
        return 0
    }
    return castFunc(structType)
}

func (m TransportSize) LengthInBits() uint16 {
    return 8
}

func (m TransportSize) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func TransportSizeParse(io *utils.ReadBuffer) (TransportSize, error) {
    val, err := io.ReadInt8(8)
    if err != nil {
        return 0, nil
    }
    return TransportSizeValueOf(val), nil
}

func (e TransportSize) Serialize(io utils.WriteBuffer) error {
    err := io.WriteInt8(8, int8(e))
    return err
}

func (e TransportSize) String() string {
    switch e {
    case TransportSize_LWORD:
        return "LWORD"
    case TransportSize_BOOL:
        return "BOOL"
    case TransportSize_BYTE:
        return "BYTE"
    case TransportSize_CHAR:
        return "CHAR"
    case TransportSize_WORD:
        return "WORD"
    case TransportSize_INT:
        return "INT"
    case TransportSize_DWORD:
        return "DWORD"
    case TransportSize_DINT:
        return "DINT"
    case TransportSize_REAL:
        return "REAL"
    case TransportSize_DATE:
        return "DATE"
    case TransportSize_TIME:
        return "TIME"
    case TransportSize_DATE_AND_TIME:
        return "DATE_AND_TIME"
    case TransportSize_WCHAR:
        return "WCHAR"
    case TransportSize_LREAL:
        return "LREAL"
    }
    return ""
}
