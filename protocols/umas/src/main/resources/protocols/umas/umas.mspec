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

// https://Umas.org/docs/Umas_Application_Protocol_V1_1b.pdf

// Remark: The different fields are encoded in Big-endian.

[type UmasConstants
    [const          uint 16     UmasTcpDefaultPort 502]
]

[discriminatedType ModbusTcpADU byteOrder='BIG_ENDIAN'
    // It is used for transaction pairing, the Umas server copies in the response the transaction
    // identifier of the request.
    [simple         uint 16     transactionIdentifier]

    // It is used for intra-system multiplexing. The Umas protocol is identified by the value 0.
    [const          uint 16     protocolIdentifier    0x0000]

    // The length field is a byte count of the following fields, including the Unit Identifier and
    // data fields.
    [implicit       uint 16     length                'COUNT(pduArray) + 1']

    // This field is used for intra-system routing purpose. It is typically used to communicate to
    // a Umas+ or a Umas serial line slave through a gateway between an Ethernet TCP-IP network
    // and a Umas serial line. This field is set by the Umas Client in the request and must be
    // returned with the same value in the response by the server.
    [simple         uint 8      unitIdentifier]

    // The actual Modbus payload
    [array      byte        pduArray         count   'length - 1']
]

[discriminatedType ModbusPDU(uint 8 umasRequestFunctionKey)
    [discriminator bit         errorFlag]
    [discriminator uint 7      functionFlag]
    [typeSwitch errorFlag,functionFlag
        ['true'                     ModbusPDUError
            [simple ModbusErrorCode  exceptionCode]
        ]

        ['false','0x5a'     UmasPDU
            [simple     UmasPDUItem('umasRequestFunctionKey')    item]
        ]
    ]
]

[type UmasPDUItem(uint 8 umasRequestFunctionKey) byteOrder='LITTLE_ENDIAN'
    [simple     uint 8     pairingKey]
    [discriminator     uint 8     umasFunctionKey]
    [typeSwitch umasFunctionKey, umasRequestFunctionKey
        ['0x01'      UmasInitCommsRequest
            [simple     uint 8         unknownObject]
        ]
        ['0x02'      UmasPDUPlcIdentRequest
        ]
        ['0x03'      UmasPDUProjectInfoRequest
            [simple uint 8 subcode]
        ]
        ['0x04'      UmasPDUPlcStatusRequest
        ]
        ['0x20'      UmasPDUReadMemoryBlockRequest
            [simple     uint 8         range]
            [simple     uint 16        blockNumber]
            [simple     uint 16        offset]
            [simple     uint 16        unknownObject1]
            [simple     uint 16        numberOfBytes]
        ]
        ['0x26'     UmasPDUReadUnlocatedVariableNamesRequest
            [const      uint 16         range 0xdd02]
            [simple     uint 32         hardwareId]
            [simple     uint 8          hardwareIdIndex]
            [simple     uint 16         blockNo]
            [const      uint 32         blank 0x0000]
        ]

        ['0xFE', '0x01'     UmasInitCommsResponse
            [simple     uint 16         maxFrameSize]
            [simple     uint 16         firmwareVersion]
            [simple     uint 32         notSure]
            [simple     uint 32         internalCode]
            [simple     uint 8          hostnameLength]
            [simple     vstring         'hostnameLength*8' hostname]
        ]
        ['0xFE', '0x02'     UmasPDUPlcIdentResponse
            [simple     uint 8          range]
            [simple     uint 32         ident]
            [simple     uint 16         model]
            [simple     uint 16         comVersion]
            [simple     uint 16         comPatch]
            [simple     uint 16         intVersion]
            [simple     uint 16         hardwareVersion]
            [simple     uint 32         crashCode]
            [simple     uint 32          hostnameLength]
            [simple     vstring         'hostnameLength*8' hostname]
            [simple     uint 8          numberOfMemoryBanks
            [array      PlcMemoryBlockIdent memoryIdents count 'numberOfMemoryBanks']
        ]
        ['0xFE', '0x04'     UmasPDUPlcStatusResponse
            [simple     uint 24         notUsed]
            [simple     uint 8          numberOfBlocks]
            [array      uint 32         blocks count 'numberOfBlocks']
        ]
        ['0xFE', '0x20'     UmasPDUReadMemoryBlockResponse
            [simple     uint 16         numberOfBytes]
            [array      uint 8          block count 'numberOfBytes']
        ]
        ['0xFE', '0x26'     UmasPDUReadUnlocatedVariableNamesResponse
            [simple     uint 32         range]
            [simple     uint 16         noOfRecords]
            [array      UmasUnlocatedVariableReference         records count 'noOfRecords']
        ]
    ]
]

[type UmasMemoryBlock(uint 16 blockNumber, uint 16 offset)
    [typeSwitch blockNumber, offset
        ['0x30', '0x00' UmasMemoryBlockBasicInfo
            [simple     uint 8          range]
            [simple uint 32 notSure]
            [simple uint 32 hardwareId]
        ]
    ]
]

[type UmasUnlocatedVariableReference
    [simple     uint 16          unknown1]
    [simple     uint 16          unknown2]
    [simple     uint 16          unknown3]
    [simple     uint 16          unknown4]
    [simple     uint 8           stringLength]
    [simple     uint 8           stringType]
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, stringLength)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, stringLength)' '(stringLength * 8)'']
]

[type PlcMemoryBlockIdent
    [simple uint 8 blockType]
    [simple uint 8 folio]
    [simple uint 16 status]
    [simple uint 32 memoryLength]
]


[dataIo DataItem(UmasDataType dataType, uint 16 numberOfValues)
    [typeSwitch dataType,numberOfValues
        ['BOOL','1'  BOOL
            // TODO: Possibly change the order of the bit and the reserved part.
            [reserved uint 15 '0x0000'                         ]
            [simple   bit     value                            ]
        ]
        ['BOOL'      List
            // TODO: Handle adding some reserved bits at the end to fill up the last word.
            [array    bit     value count 'numberOfValues'     ]
        ]
        ['BYTE','1'  BYTE
            [reserved uint 8 '0x00']
            [simple uint 8 value]
        ]
        ['BYTE' List
            // TODO: If the number of values is odd, add a reserved byte
            [array    bit     value count 'numberOfValues * 8' ]
        ]
        ['WORD'      WORD
            [simple   uint 16 value]
        ]
        ['DWORD'     DWORD
            [simple   uint 32 value]
        ]
        ['LWORD'     LWORD
            [simple   uint 64 value]
        ]
        ['SINT','1' SINT
            [reserved uint 8  '0x00']
            [simple   int 8   value ]
        ]
        ['SINT' List
            [array int 8 value count 'numberOfValues']
        ]
        ['INT','1' INT
            [simple int 16 value]
        ]
        ['INT' List
            [array int 16 value count 'numberOfValues']
        ]
        ['DINT','1' DINT
            [simple int 32 value]
        ]
        ['DINT' List
            [array int 32 value count 'numberOfValues']
        ]
        ['LINT','1' LINT
            [simple int 64 value]
        ]
        ['LINT' List
            [array int 64 value count 'numberOfValues']
        ]
        ['USINT','1' USINT
            [reserved uint 8 '0x00']
            [simple   uint 8 value ]
        ]
        ['USINT' List
            [array uint 8 value count 'numberOfValues']
        ]
        ['UINT','1' UINT
            [simple uint 16 value]
        ]
        ['UINT' List
            [array uint 16 value count 'numberOfValues']
        ]
        ['UDINT','1' UDINT
            [simple uint 32 value]
        ]
        ['UDINT' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['ULINT','1' ULINT
            [simple uint 64 value]
        ]
        ['ULINT' List
            [array uint 64 value count 'numberOfValues']
        ]
        ['REAL','1' REAL
            [simple float 32  value]
        ]
        ['REAL' List
            [array float 32 value count 'numberOfValues']
        ]
        ['LREAL','1' LREAL
            [simple float 64  value]
        ]
        ['LREAL' List
            [array float 64 value count 'numberOfValues']
        ]
        ['CHAR','1' CHAR
            [simple string 8 value encoding='"UTF-8"']
        ]
        ['CHAR' List
            [array string 8 value count 'numberOfValues' encoding='"UTF-8"']
        ]
        ['WCHAR','1' WCHAR
            [simple string 16 value encoding='"UTF-16"']
        ]
        ['WCHAR' List
            [array string 16 value count 'numberOfValues' encoding='"UTF-16"']
        ]
    ]
]

[enum uint 8 UmasDataType(uint 8 dataTypeSize)
    ['1' BOOL ['2']]
    ['2' BYTE ['2']]
    ['3' WORD ['2']]
    ['4' DWORD ['4']]
    ['5' LWORD ['8']]
    ['6' SINT ['2']]
    ['7' INT ['2']]
    ['8' DINT ['4']]
    ['9' LINT ['8']]
    ['10' USINT ['2']]
    ['11' UINT ['2']]
    ['12' UDINT ['4']]
    ['13' ULINT ['8']]
    ['14' REAL ['4']]
    ['15' LREAL ['8']]
    ['16' TIME ['8']]
    ['17' LTIME ['8']]
    ['18' DATE ['8']]
    ['19' LDATE ['8']]
    ['20' TIME_OF_DAY ['8']]
    ['21' LTIME_OF_DAY ['8']]
    ['22' DATE_AND_TIME ['8']]
    ['23' LDATE_AND_TIME ['8']]
    ['24' CHAR ['1']]
    ['25' WCHAR ['2']]
    ['26' STRING ['1']]
    ['27' WSTRING ['2']]
]

[enum uint 8 ModbusErrorCode
    ['1'    ILLEGAL_FUNCTION]
    ['2'    ILLEGAL_DATA_ADDRESS]
    ['3'    ILLEGAL_DATA_VALUE]
    ['4'    SLAVE_DEVICE_FAILURE]
    ['5'    ACKNOWLEDGE]
    ['6'    SLAVE_DEVICE_BUSY]
    ['7'    NEGATIVE_ACKNOWLEDGE]
    ['8'    MEMORY_PARITY_ERROR]
    ['10'   GATEWAY_PATH_UNAVAILABLE]
    ['11'   GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
]

[enum uint 8 UmasDeviceInformationLevel
    ['0x01' BASIC     ]
    ['0x02' REGULAR   ]
    ['0x03' EXTENDED  ]
    ['0x04' INDIVIDUAL]
]

[enum uint 7 UmasDeviceInformationConformityLevel
    ['0x01' BASIC_STREAM_ONLY   ]
    ['0x02' REGULAR_STREAM_ONLY ]
    ['0x03' EXTENDED_STREAM_ONLY]
]

[enum uint 8 UmasDeviceInformationMoreFollows
    ['0x00' NO_MORE_OBJECTS_AVAILABLE]
    ['0xFF' MORE_OBJECTS_AVAILABLE   ]
]
