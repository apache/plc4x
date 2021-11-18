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

[discriminatedType FirmataMessage(bit response) byteOrder='BIG_ENDIAN'
    [discriminator uint 4 messageType]
    [typeSwitch 'messageType'
        // Reading operations
        // Data-Format is: in both bytes only the least significant 7 bits
        // count (unsigned 14 bit integer number)
        // The first byte contains the least significant part
        // The second byte contains the most significant part
        ['0xE' FirmataMessageAnalogIO
            [simple uint 4 pin]
            [array int 8 data count '2']
        ]
        // Bitmask containing the first 7 bits in the least significant
        // bits of the first byte and bit 8 in the second byte
        // The 'pinBlock' refers to the block of bytes (0 refers to the
        // first 8 pins, 1 to the second and so on.
        ['0x9' FirmataMessageDigitalIO
            [simple uint 4 pinBlock]
            [array int 8 data count '2']
        ]

        ['0xC' FirmataMessageSubscribeAnalogPinValue
            [simple uint 4 pin]
            [reserved uint 7 '0x00']
            [simple bit enable]
        ]
        ['0xD' FirmataMessageSubscribeDigitalPinValue
            [simple uint 4 pin]
            [reserved uint 7 '0x00']
            [simple bit enable]
        ]

        // Command
        ['0xF' FirmataMessageCommand
            [simple FirmataCommand('response') command]
        ]
    ]
]

[discriminatedType FirmataCommand(bit response)
    [discriminator uint 4 commandCode]
    [typeSwitch 'commandCode'
        ['0x0' FirmataCommandSysex
            [simple SysexCommand('response') command]
            [reserved uint 8 '0xF7']
        ]
        ['0x4' FirmataCommandSetPinMode
            [simple uint 8 pin]
            [simple PinMode mode]
        ]
        ['0x5' FirmataCommandSetDigitalPinValue
            [simple uint 8 pin]
            [reserved uint 7 '0x00']
            [simple bit on]
        ]
        ['0x9' FirmataCommandProtocolVersion
            [simple uint 8 majorVersion]
            [simple uint 8 minorVersion]
        ]
        ['0xF' FirmataCommandSystemReset
        ]
    ]
]

[discriminatedType SysexCommand(bit response)
    [discriminator uint 8 commandType]
    [typeSwitch 'commandType','response'
        ['0x00' SysexCommandExtendedId
            [array int 8 id count '2']
        ]
        ['0x69','false' SysexCommandAnalogMappingQueryRequest
        ]
        ['0x69','true' SysexCommandAnalogMappingQueryResponse
            [simple uint 8 pin]
        ]
        ['0x6A' SysexCommandAnalogMappingResponse
        ]
        ['0x6B' SysexCommandCapabilityQuery
        ]
        ['0x6C' SysexCommandCapabilityResponse
        ]
        ['0x6D' SysexCommandPinStateQuery
            [simple uint 8 pin]
        ]
        ['0x6E' SysexCommandPinStateResponse
            [simple uint 8 pin]
            [simple uint 8 pinMode]
            [simple uint 8 pinState]
        ]
        ['0x6F' SysexCommandExtendedAnalog
        ]
        ['0x71' SysexCommandStringData
        ]
        ['0x79','false' SysexCommandReportFirmwareRequest
        ]
        ['0x79','true' SysexCommandReportFirmwareResponse
            [simple uint 8 majorVersion]
            [simple uint 8 minorVersion]
            [manualArray byte 'fileName' terminated 'STATIC_CALL("isSysexEnd", readBuffer)' 'STATIC_CALL("parseSysexString", readBuffer)' 'STATIC_CALL("serializeSysexString", writeBuffer, _value)' 'STATIC_CALL("lengthSysexString", fileName)']
        ]
        ['0x7A' SysexCommandSamplingInterval
        ]
        ['0x7E' SysexCommandSysexNonRealtime
        ]
        ['0x7F' SysexCommandSysexRealtime
        ]
    ]
]

[enum uint 8 PinMode
    ['0x0' PinModeInput]
    ['0x1' PinModeOutput]
    ['0x2' PinModeAnalog]
    ['0x3' PinModePwm]
    ['0x4' PinModeServo]
    ['0x5' PinModeShift]
    ['0x6' PinModeI2C]
    ['0x7' PinModeOneWire]
    ['0x8' PinModeStepper]
    ['0x9' PinModeEncoder]
    ['0xA' PinModeSerial]
    ['0xB' PinModePullup]
]
