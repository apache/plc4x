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

[discriminatedType ModbusPDU(uint 8 umasRequestFunctionKey, uint 16 byteLength)
    [discriminator bit         errorFlag]
    [discriminator uint 7      functionFlag]
    [typeSwitch errorFlag,functionFlag
        ['true'                     ModbusPDUError
            [simple ModbusErrorCode  exceptionCode]
        ]
        ['false','0x5a'     UmasPDU
            [simple     UmasPDUItem('umasRequestFunctionKey','byteLength - 1')    item]
        ]
    ]
]

[type UmasPDUItem(uint 8 umasRequestFunctionKey, uint 16 byteLength) byteOrder='LITTLE_ENDIAN'
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
        ['0x22'      UmasPDUReadVariableRequest
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableReadRequestReference variables count 'variableCount']
        ]
        ['0x23'      UmasPDUWriteVariableRequest
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableWriteRequestReference variables count 'variableCount']
        ]
        ['0x26'     UmasPDUReadUnlocatedVariableNamesRequest
            [simple     uint 16         recordType]
            [simple     uint 8          index]
            [simple     uint 32         hardwareId]
            [simple     uint 16         blockNo]
            [simple     uint 16         offset]
            [const      uint 16         blank 0x00]
        ]
        ['0xFD'     UmasPDUErrorResponse
            [array      uint 8          block count 'byteLength - 2']
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
            [simple     uint 16         range]
            [simple     uint 32         ident]
            [simple     uint 16         model]
            [simple     uint 16         comVersion]
            [simple     uint 16         comPatch]
            [simple     uint 16         intVersion]
            [simple     uint 16         hardwareVersion]
            [simple     uint 32         crashCode]
            [simple     uint 16         unknown1]
            [simple     uint 8          hostnameLength]
            [simple     vstring         'hostnameLength*8' hostname]
            [simple     uint 8          numberOfMemoryBanks]
            [array      PlcMemoryBlockIdent memoryIdents count 'numberOfMemoryBanks']
        ]
        ['0xFE', '0x04'     UmasPDUPlcStatusResponse
            [simple     uint 8          notUsed1]
            [simple     uint 16         notUsed2]
            [simple     uint 8          numberOfBlocks]
            [array      uint 32         blocks count 'numberOfBlocks']
        ]
        ['0xFE', '0x20'     UmasPDUReadMemoryBlockResponse
            [simple     uint 8          range]
            [simple     uint 16         numberOfBytes]
            [array      uint 8          block count 'numberOfBytes']
        ]
        ['0xFE', '0x22'     UmasPDUReadVariableResponse
            [array      uint 8          block count 'byteLength - 2']
        ]
        ['0xFE', '0x23'     UmasPDUWriteVariableResponse
            [array      uint 8          block count 'byteLength - 2']
        ]
        ['0xFE', '0x26'     UmasPDUReadUnlocatedVariableResponse
            [array      uint 8          block count 'byteLength - 2']
        ]
    ]
]

[type UmasMemoryBlock(uint 16 blockNumber, uint 16 offset)
    [typeSwitch blockNumber, offset
        ['0x30', '0x00' UmasMemoryBlockBasicInfo
            [simple     uint 16          range]
            [simple uint 16 notSure]
            [simple uint 8  index]
            [simple uint 32 hardwareId]
        ]
    ]
]

[type UmasPDUReadUnlocatedVariableNamesResponse
    [simple     uint 8          range]
    [simple     uint 16         nextAddress]
    [simple     uint 16         unknown1]
    [simple     uint 16         noOfRecords]
    [array      UmasUnlocatedVariableReference         records count 'noOfRecords']
]

[type UmasPDUReadUmasUDTDefinitionResponse
    [simple     uint 8          range]
    [simple     uint 32         unknown1]
    [simple     uint 16         noOfRecords]
    [array      UmasUDTDefinition         records count 'noOfRecords']
]

[type UmasPDUReadDatatypeNamesResponse
    [simple     uint 8         range]
    [simple     uint 16        nextAddress]
    [simple     uint 8         unknown1]
    [simple     uint 16        noOfRecords]
    [array      UmasDatatypeReference         records count 'noOfRecords']
]

[type VariableReadRequestReference
    [simple     uint 4           isArray]
    [simple     uint 4           dataSizeIndex]
    [simple     uint 16          block]
    [const      uint 8           unknown1 0x01]
    [simple     uint 16          baseOffset]
    [simple     uint 8           offset]
    [optional   uint 16          arrayLength 'isArray']
]

[type VariableWriteRequestReference
    [simple     uint 4           isArray]
    [simple     uint 4           dataSizeIndex]
    [simple     uint 16          block]
    [simple     uint 16          baseOffset]
    [simple     uint 16          offset]
    [optional   uint 16          arrayLength 'isArray']
    [array      byte             recordData     length  'isArray == 1 ? dataSizeIndex * arrayLength : dataSizeIndex']
]

[type UmasUnlocatedVariableReference
    [simple     uint 16          dataType]
    [simple     uint 16          block]
    [simple     uint 16          offset]
    [simple     uint 16          unknown4]
    [simple     uint 16          stringLength]
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, stringLength)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, stringLength)' '(stringLength * 8)']
]

[type UmasUDTDefinition
    [simple     uint 16          dataType]
    [simple     uint 16          offset]
    [simple     uint 16          unknown5]
    [simple     uint 16          unknown4]
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, -1)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, -1)' '(stringLength * 8)']
]

[type UmasDatatypeReference
    [simple     uint 16          dataSize]
    [simple     uint 16          unknown1]
    [simple     uint 8           classIdentifier]
    [simple     uint 8           dataType]
    [simple     uint 8           stringLength]
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, stringLength)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, stringLength)' '(stringLength * 8)']
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
            [reserved uint 7 '0x0000'                         ]
            [simple   bit     value                           ]
        ]
        ['EBOOL','1'  BOOL
            // TODO: Possibly change the order of the bit and the reserved part.
            [reserved uint 7 '0x0000'                         ]
            [simple   bit     value                            ]
        ]
        ['BYTE','1'  BYTE
            [simple byte value]
        ]
        ['BYTE' List
            [array    byte     value count 'numberOfValues' ]
        ]
        ['WORD'      WORD
            [simple   uint 16 value]
        ]
        ['DWORD'     DWORD
            [simple   uint 32 value]
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
        ['REAL','1' REAL
            [simple float 32  value]
        ]
        ['REAL' List
            [array float 32 value count 'numberOfValues']
        ]
        ['STRING','1' STRING
            [manual vstring value  'STATIC_CALL("parseTerminatedStringBytes", readBuffer, numberOfValues)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, _value, numberOfValues)' '(numberOfValues * 8)']
        ]
        ['STRING' List
            [array float 32 value count 'numberOfValues']
        ]
        ['TIME','1' TIME
            [simple uint 32 value]
        ]
        ['TIME' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['DATE','1' DATE
            [simple uint 8 day encoding='BCD']
            [simple uint 8 month encoding='BCD']
            [simple uint 16 year encoding='BCD']
        ]
        ['TOD','1' TIME_OF_DAY
            [simple uint 32 value]
        ]
        ['TOD' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['DATE_AND_TIME','1' DATE_AND_TIME
            [simple uint 8 unused]
            [simple uint 8 seconds encoding='BCD']
            [simple uint 8 minutes encoding='BCD']
            [simple uint 8 hour encoding='BCD']
            [simple uint 8 day encoding='BCD']
            [simple uint 8 month encoding='BCD']
            [simple uint 16 year encoding='BCD']
        ]
    ]
]

[enum uint 8 UmasDataType(uint 8 dataTypeSize, uint 8 requestSize, vstring data_type_conversion )
    ['1' BOOL ['1','1','"BOOL"']]
    ['2' UNKNOWN2 ['1','1','"BOOL"']]
    ['3' UNKNOWN3 ['1','1','"BOOL"']]
    ['4' INT ['2', '2','"INT"']]
    ['5' UINT ['2','2','"UINT"']]
    ['6' DINT ['4','3','"DINT"']]
    ['7' UDINT ['4','3','"UDINT"']]
    ['8' REAL ['4','3','"REAL"']]
    ['9' STRING ['1','17','"STRING"']]
    ['10' TIME ['4','3','"TIME"']]
    ['11' UNKNOWN11 ['1','1','"BYTE"']]
    ['12' UNKNOWN12 ['1','1','"BYTE"']]
    ['13' UNKNOWN13 ['1','1','"BYTE"']]
    ['14' DATE ['4','3','"DATE"']]
    ['15' TOD ['4','3','"TIME_OF_DAY"']]
    ['16' DATE_AND_TIME ['8','4','"DATE_AND_TIME"']]
    ['17' UNKNOWN17 ['1','1','"BYTE"']]
    ['18' UNKNOWN18 ['1','1','"BYTE"']]
    ['19' UNKNOWN19 ['1','1','"BYTE"']]
    ['20' UNKNOWN20 ['1','1','"BYTE"']]
    ['21' BYTE ['1','1','"BYTE"']]
    ['22' WORD ['2','2','"WORD"']]
    ['23' DWORD ['4','3','"DWORD"']]
    ['24' UNKNOWN24 ['1','1','"BYTE"']]
    ['25' EBOOL ['1','1','"BOOL"']]
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
