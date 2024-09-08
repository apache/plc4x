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

// https://modbus.org/docs/Modbus_Application_Protocol_V1_1b.pdf

// Remark: The different fields are encoded in Big-endian for Modbus. UMAS uses Little-Endian

[type ModbusConstants
    [const          uint 16     ModbusTcpDefaultPort 502]
]

[discriminatedType ModbusTcpADU byteOrder='BIG_ENDIAN'
    // It is used for transaction pairing, the Mobus server copies in the response the transaction
    // identifier of the request.
    [simple         uint 16     transactionIdentifier]

    // It is used for intra-system multiplexing. The Modbus protocol is identified by the value 0.
    [const          uint 16     protocolIdentifier    0x0000]

    // The length field is a byte count of the following fields, including the Unit Identifier and
    // data fields.
    [implicit       uint 16     length                'COUNT(pduArray) + 1']

    // This field is used for intra-system routing purpose. It is typically used to communicate to
    // a Modbus serial line slave through a gateway between an Ethernet TCP-IP network
    // and a Umas serial line. This field is set by the Modbus Client in the request and must be
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
    // Not too sure what this is. It is normally 0x00 or 0xFF
    [simple     uint 8     pairingKey]
    [discriminator     uint 8     umasFunctionKey]
    [typeSwitch umasFunctionKey, umasRequestFunctionKey
        // Opens a connection to the PLC to enable reading and writing. It use one connection from the PLC.
        ['0x01'      UmasInitCommsRequest
            [const     uint 8         unknownObject     0x00]
        ]
        ['0x02'      UmasPDUPlcIdentRequest
        ]
        // Requests information about the running project
        ['0x03'      UmasPDUProjectInfoRequest
            [simple uint 8 subcode]
        ]
        // Requests information about the PLC, it mainly contains useful information about each memory block
        ['0x04'      UmasPDUPlcStatusRequest
        ]
        // Request to read data from a memory block defined by the Block No. and Offset.
        ['0x20'      UmasPDUReadMemoryBlockRequest
            [const      uint 8         range        0x00]
            // The memory block number to be read.
            [simple     uint 16        blockNumber]
            // The offset from the start of the block to read
            [simple     uint 16        offset]
            // Not sure but probably some other offset
            [simple     uint 16        unknownObject1]
            // The number of bytes to read
            [simple     uint 16        numberOfBytes]
        ]
        // Used to read the values of variables. The variable information must firstly be looked up using a UmasPDUReadUnlocatedVariableNamesRequest
        ['0x22'      UmasPDUReadVariableRequest
            // This can be calculated from a UmasPDUPlcStatusResponse by adding blocks 3 and 4.
            // It is a unique value for the current project on the PLC. If there are any modifications then this changes.
            // It is used to detect a programming change to know when to re-read the tag information.
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableReadRequestReference variables count 'variableCount']
        ]
        ['0x23'      UmasPDUWriteVariableRequest
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableWriteRequestReference variables count 'variableCount']
        ]
        // Is used to read information from the data dictionary such as tag names, data types, etc...
        ['0x26'     UmasPDUReadDataDictionaryRequest
            // I've seen values 0xDD02 and 0xDD03 being used
            [simple     uint 16         recordType]
            // This is an index from the UmasMemoryBlockBasicInfo message
            [simple     uint 8          index]
            // This seems to be a unique number returned from the UmasMemoryBlockBasicInfo
            [simple     uint 32         hardwareId]
            [simple     uint 16         blockNo]
            [simple     uint 16         offset]
            [const      uint 16         blank 0x00]
        ]
        ['0xFD'     UmasPDUErrorResponse
            [array      uint 8          block count 'byteLength - 2']
        ]
        ['0xFE', '0x01'     UmasInitCommsResponse
            // The largest UMAS packet size. I don't think it includes the Modbus header
            [simple     uint 16         maxFrameSize]
            // PLC Firmware Version
            [simple     uint 16         firmwareVersion]
            [simple     uint 32         notSure]
            [simple     uint 32         internalCode]
            [simple     uint 8          hostnameLength]
            // PLC Name returned a string
            [simple     vstring         'hostnameLength*8' hostname]
        ]
        ['0xFE', '0x02'     UmasPDUPlcIdentResponse
            // I don't have a lot of information about these, but can be used to identify the PLC/Firmware/Hardware version if we ever have a need
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
            // List of memory blocks within the PLC.
            [array      PlcMemoryBlockIdent memoryIdents count 'numberOfMemoryBanks']
        ]
        ['0xFE', '0x04'     UmasPDUPlcStatusResponse
            // This has some information, presumably about the current project. We need blocks 3 and 4 to work out the CRC of the project.
            // We monitor these to check if the project has changed, if so we re-read the data dictionary
            [simple     uint 8          notUsed1]
            [simple     uint 16         notUsed2]
            [simple     uint 8          numberOfBlocks]
            [array      uint 32         blocks count 'numberOfBlocks']
        ]
        ['0xFE', '0x20'     UmasPDUReadMemoryBlockResponse
            // Simply reads a section of memory. A UmasMemoryBlock can be returned, I don't have much information on what other blocks mean though.
            [simple     uint 8          range]
            [simple     uint 16         numberOfBytes]
            [array      uint 8          block count 'numberOfBytes']
        ]
        ['0xFE', '0x22'     UmasPDUReadVariableResponse
            // It just returns the variable data which can then be decoded with the DataIO.
            [array      uint 8          block count 'byteLength - 2']
        ]
        ['0xFE', '0x23'     UmasPDUWriteVariableResponse
            [array      uint 8          block count 'byteLength - 2']
        ]
        // Reads information from the data dictionary.
        // we can read one of UmasPDUReadUnlocatedVariableNamesResponse, UmasPDUReadUmasUDTDefinitionResponse or UmasPDUReadDatatypeNamesResponse
        // Which depends on the request that was sent.
        ['0xFE', '0x26'     UmasPDUReadDataDictionaryResponse
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
    [const      uint 8           unknown1 0x01]
    [simple     uint 16          baseOffset]
    [simple     uint 8           offset]
    [optional   uint 16          arrayLength 'isArray']
    [array      byte       recordData     length  'isArray == 1 ? dataSizeIndex * arrayLength : dataSizeIndex']
]

[type UmasUnlocatedVariableReference
    [simple     uint 8           dataType]
    [simple     uint 8           unknown1]
    [simple     uint 16          block]
    [simple     uint 8           offset]
    [simple     uint 8           unknown5]
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
            // TODO: Possibly change the order of the bit and the reserved part.
            [reserved uint 7 '0x0000'                         ]
            [simple   bit     value                            ]
        ]
        ['BOOL'      List
            // TODO: Handle adding some reserved bits at the end to fill up the last word.
            [array    bit     value count 'numberOfValues'     ]
        ]
        ['BYTE','1'  BYTE
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
            [manual vstring value  'STATIC_CALL("parseTerminatedStringBytes", readBuffer, numberOfValues)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, numberOfValues)' '(numberOfValues * 8)']
        ]
        ['STRING' List
            [array float 32 value count 'numberOfValues']
        ]
    ]
]

[enum uint 8 UmasDataType(uint 8 dataTypeSize, uint 8 requestSize)
    ['1' BOOL ['1','1']]
    ['2' UNKNOWN2 ['1','1']]
    ['3' UNKNOWN3 ['1','1']]
    ['4' INT ['2', '2']]
    ['5' UINT ['2','2']]
    ['6' DINT ['4','3']]
    ['7' UDINT ['4','3']]
    ['8' REAL ['4','3']]
    ['9' STRING ['1','17']]
    ['10' TIME ['4','3']]
    ['11' UNKNOWN11 ['1','1']]
    ['12' UNKNOWN12 ['1','1']]
    ['13' UNKNOWN13 ['1','1']]
    ['14' DATE ['4','3']]
    ['15' TOD ['4','3']]
    ['16' DT ['4','3']]
    ['17' UNKNOWN17 ['1','1']]
    ['18' UNKNOWN18 ['1','1']]
    ['19' UNKNOWN19 ['1','1']]
    ['20' UNKNOWN20 ['1','1']]
    ['21' BYTE ['1','1']]
    ['22' WORD ['2','2']]
    ['23' DWORD ['4','3']]
    ['24' UNKNOWN24 ['1','1']]
    ['25' EBOOL ['1','1']]
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
