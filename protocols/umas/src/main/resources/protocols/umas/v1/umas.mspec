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

////////////////////////////////////////////////////////////////
// UMAS Protocol (Unified Messaging Application Services)
//
// Schneider Electric's proprietary protocol for Modicon PLCs (M340, M580,
// Quantum, Premium). Tunneled inside Modbus/TCP using Function Code 90 (0x5A).
//
// Key detail: Modbus/TCP MBAP header is big-endian, UMAS payload is little-endian.
//
// References:
//   - Version 0 of the umas.mspec
//   - Kaspersky ICS CERT: "The secrets of Schneider Electric's UMAS protocol"
//   - Wireshark Lua dissectors (zaltzman, yanissec, biero-el-corridor)
//   - Recordings of the Schneider Electric OPC UA Server Expert communicating
//     with a Simulated PLC (Tag loading)
//   - Recordings of the Schneider Electric OPC UA Server Expert communicating
//     with a Simulated PLC (Using UAExpert to produce subscription traffic)
//   - Recordings of the Schneider Electric OPC UA Server Expert communicating
//     with a Simulated PLC (Using Prosys OPC UA Browser to produce read/write
//     traffic)
////////////////////////////////////////////////////////////////

[constants
    [const          uint 16     UmasTcpDefaultPort 502]
]

////////////////////////////////////////////////////////////////
// Modbus/TCP Transport Framing (BIG-ENDIAN)
////////////////////////////////////////////////////////////////

// Modbus/TCP Application Data Unit — wraps the UMAS payload.
// The MBAP header uses big-endian byte ordering per the Modbus specification.
[type ModbusTcpADU(uint 8 umasRequestFunctionKey) byteOrder='"BIG_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' floatEncoding='"IEEE754"' stringEncoding='"UTF8"'
    // Transaction identifier for request/response correlation
    [simple         uint 16     transactionIdentifier]

    // Protocol identifier — always 0x0000 for Modbus
    [const          uint 16     protocolIdentifier    0x0000]

    // Length of remaining bytes (unitIdentifier + pdu)
    [implicit       uint 16     length                'pdu.lengthInBytes + 1']

    // Unit/slave address
    [simple         uint 8      unitIdentifier]

    // UMAS payload wrapped in a Modbus PDU (FC 0x5A).
    // umasRequestFunctionKey is forwarded for response discrimination — pass 0
    // when the original request FC is unknown (the catch-all handles it).
    [simple         ModbusPDU('umasRequestFunctionKey', 'length - 1')   pdu]
]

////////////////////////////////////////////////////////////////
// Modbus PDU with UMAS discrimination
////////////////////////////////////////////////////////////////

// The Modbus PDU discriminates on error flag and function code.
// FC 0x5A (90) identifies a UMAS payload.
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

////////////////////////////////////////////////////////////////
// UMAS PDU Item (LITTLE-ENDIAN)
//
// The core UMAS protocol data unit. Byte order switches to little-endian
// here — all UMAS payload fields are little-endian.
//
// Requests are discriminated by umasFunctionKey alone.
// Responses use 0xFE (success) or 0xFD (error) as umasFunctionKey;
// the second discriminator umasRequestFunctionKey (the original request FC)
// selects the specific response type. A catch-all ['0xFE'] at the end
// handles responses when the original request FC is unknown (pass 0).
////////////////////////////////////////////////////////////////

[type UmasPDUItem(uint 8 umasRequestFunctionKey, uint 16 byteLength) byteOrder='"LITTLE_ENDIAN"'
    [simple     uint 8     pairingKey]
    [discriminator     uint 8     umasFunctionKey]
    [typeSwitch umasFunctionKey, umasRequestFunctionKey

        // ---- FC 0x01: Init Communication ----
        ['0x01'      UmasInitCommsRequest
            [simple     uint 8          subCode]
        ]
        ['0xFE', '0x01'     UmasInitCommsResponse
            [simple     uint 16         maxFrameSize]
            [simple     uint 16         firmwareVersion]
            [simple     uint 32         notSure]
            [simple     uint 32         internalCode]
            [simple     uint 8          hostnameLength]
            [simple     vstring         'hostnameLength*8' hostname]
        ]

        // ---- FC 0x02: PLC Identification ----
        ['0x02'      UmasPDUPlcIdentRequest
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

        // ---- FC 0x03: Project Info ----
        ['0x03'      UmasPDUProjectInfoRequest
            [simple     uint 8          subcode]
        ]
        ['0xFE', '0x03'     UmasPDUProjectInfoResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x04: PLC Status ----
        ['0x04'      UmasPDUPlcStatusRequest
        ]
        ['0xFE', '0x04'     UmasPDUPlcStatusResponse
            [simple     uint 8          notUsed1]
            [simple     uint 16         notUsed2]
            [simple     uint 8          numberOfBlocks]
            [array      uint 32         blocks count 'numberOfBlocks']
            // Additional status data present on M580 PLCs (absent on M340)
            [array      byte            additionalData count 'byteLength - 2 - 4 - (numberOfBlocks * 4)']
        ]

        // ---- FC 0x06: Read Card Info ----
        ['0x06'      UmasPDUReadCardInfoRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x06'     UmasPDUReadCardInfoResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x0A: Repeat / Echo ----
        ['0x0A'      UmasPDURepeatRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x0A'     UmasPDURepeatResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x10: Take PLC Reservation ----
        ['0x10'      UmasPDUTakePlcReservationRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x10'     UmasPDUTakePlcReservationResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x11: Release PLC Reservation ----
        ['0x11'      UmasPDUReleasePlcReservationRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x11'     UmasPDUReleasePlcReservationResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x12: Keep Alive ----
        ['0x12'      UmasPDUKeepAliveRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x12'     UmasPDUKeepAliveResponse
            [array      byte           block count 'byteLength - 2']
        ]

        // ---- FC 0x20: Read Memory Block ----
        ['0x20'      UmasPDUReadMemoryBlockRequest
            [simple     uint 8         range]
            [simple     uint 16        blockNumber]
            [simple     uint 16        offset]
            [simple     uint 16        unknownObject1]
            [simple     uint 16        numberOfBytes]
        ]
        ['0xFE', '0x20'     UmasPDUReadMemoryBlockResponse
            [simple     uint 8          range]
            [simple     uint 16         numberOfBytes]
            [array      byte            block count 'numberOfBytes']
        ]

        // ---- FC 0x22: Read Variables ----
        ['0x22'      UmasPDUReadVariableRequest
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableReadRequestReference variables count 'variableCount']
        ]
        ['0xFE', '0x22'     UmasPDUReadVariableResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x23: Write Variables ----
        ['0x23'      UmasPDUWriteVariableRequest
            [simple     uint 32        crc]
            [simple     uint 8        variableCount]
            [array      VariableWriteRequestReference variables count 'variableCount']
        ]
        ['0xFE', '0x23'     UmasPDUWriteVariableResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x24: Read Coils/Registers ----
        ['0x24'      UmasPDUReadCoilsRegistersRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x24'     UmasPDUReadCoilsRegistersResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x25: Write Coils/Registers ----
        ['0x25'      UmasPDUWriteCoilsRegistersRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0xFE', '0x25'     UmasPDUWriteCoilsRegistersResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x26: Data Dictionary ----
        ['0x26'     UmasPDUReadUnlocatedVariableNamesRequest
            [simple     uint 16         recordType]
            [simple     uint 8          index]
            [simple     uint 32         hardwareId]
            [simple     uint 16         blockNo]
            [simple     uint 16         offset]
            // Trailing padding present in DD02 requests but absent in DD03
            [optional   uint 16         blank 'recordType == 0xDD02']
        ]
        ['0xFE', '0x26'     UmasPDUReadUnlocatedVariableResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x30-0x35: Strategy Transfer ----
        ['0x30'      UmasPDUInitializeUploadRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0x31'      UmasPDUUploadBlockRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0x32'      UmasPDUEndStrategyUploadRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0x33'      UmasPDUInitializeDownloadRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0x34'      UmasPDUDownloadBlockRequest
            [array      byte           data count 'byteLength - 2']
        ]
        ['0x35'      UmasPDUEndStrategyDownloadRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x39: Read Ethernet Master Data ----
        ['0x39'      UmasPDUReadEthMasterDataRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x40: Start PLC ----
        ['0x40'      UmasPDUStartPlcRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x41: Stop PLC ----
        ['0x41'      UmasPDUStopPlcRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x50: Monitor PLC ----
        // Used by Schneider OPC UA Server Expert to read variable values from
        // M580 PLCs (instead of FC 0x22 ReadVariable which returns error 0xA1A1).
        // The request contains a list of sub-operations: register variables for
        // monitoring, read all registered values, or deregister variables.
        ['0x50'      UmasPDUMonitorPlcRequest
            [simple     uint 8          subCommand]
            [simple     uint 8          unknown]
            [simple     uint 8          numberOfSubOperations]
            [array      MonitorPlcSubOperation  subOperations count 'numberOfSubOperations']
        ]
        // Response payload is a flat byte array because the per-variable value
        // sizes depend on the data types of the registered variables, which are
        // only known at the driver layer (via the symbol table). The driver
        // parses the raw bytes using DataItem and the registered variable types.
        ['0xFE', '0x50'     UmasPDUMonitorPlcResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- FC 0x58: Check PLC ----
        ['0x58'      UmasPDUCheckPlcRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x70: Read IO Object ----
        ['0x70'      UmasPDUReadIoObjectRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x71: Write IO Object ----
        ['0x71'      UmasPDUWriteIoObjectRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- FC 0x73: Get Status Module ----
        ['0x73'      UmasPDUGetStatusModuleRequest
            [array      byte           data count 'byteLength - 2']
        ]

        // ---- Error Response (any FC) ----
        ['0xFD'     UmasPDUErrorResponse
            [array      byte            block count 'byteLength - 2']
        ]

        // ---- Generic Success Response ----
        // Catch-all for responses when umasRequestFunctionKey is unknown (pass 0).
        // Specific response types above take priority when the correct FC is supplied.
        ['0xFE'     UmasPDUSuccessResponse
            [array      byte            block count 'byteLength - 2']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// Helper Types
////////////////////////////////////////////////////////////////

// Sub-operation within an FC 0x50 MonitorPLC request.
// Discriminated by the first byte (operation type).
[discriminatedType MonitorPlcSubOperation byteOrder='"LITTLE_ENDIAN"'
    [discriminator uint 8  operationType]
    [typeSwitch operationType
        // Register or deregister a variable for monitoring.
        // After registration, the variable's current value is returned
        // in subsequent read (0x07) responses.
        ['0x05'  MonitorPlcRegisterVariable
            [simple     uint 8                   variableIndex]
            [simple     uint 16                  block]
            [simple     uint 16                  offset]
            [simple     MonitorPlcRegisterAction action]
        ]
        // Read current values for all registered variables.
        // Returns concatenated raw values in the response — the driver
        // must know each variable's data type to parse the byte stream.
        ['0x07'  MonitorPlcReadAll
        ]
        // Register a variable AND include its value in the response.
        // Combines registration with an immediate read for that variable.
        ['0x09'  MonitorPlcRegisterAndRead
            [simple     uint 8                  variableIndex]
            [simple     uint 16                 block]
            [simple     uint 16                 offset]
        ]
        // Clear/reset monitoring state. Observed in Modicon M340 captures
        // with no payload (single byte operation, like 0x07).
        ['0x0B'  MonitorPlcReset
        ]
    ]
]

// Memory block structure for specific block/offset combinations
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

// Parsed response for unlocated variable names (used by driver layer)
[type UmasPDUReadUnlocatedVariableNamesResponse
    [simple     uint 8          range]
    [simple     uint 16         nextAddress]
    [simple     uint 16         unknown1]
    [simple     uint 16         noOfRecords]
    [array      UmasUnlocatedVariableReference         records count 'noOfRecords']
]

// Parsed response for UDT definitions
[type UmasPDUReadUmasUDTDefinitionResponse
    [simple     uint 8          range]
    [simple     uint 32         unknown1]
    [simple     uint 16         noOfRecords]
    [array      UmasUDTDefinition         records count 'noOfRecords']
]

// Parsed response for datatype names
[type UmasPDUReadDatatypeNamesResponse
    [simple     uint 8         range]
    [simple     uint 16        nextAddress]
    [simple     uint 8         unknown1]
    [simple     uint 16        noOfRecords]
    [array      UmasDatatypeReference         records count 'noOfRecords']
]

// Variable read request reference — identifies a variable to read
[type VariableReadRequestReference
    [simple     uint 4           isArray]
    [simple     uint 4           dataSizeIndex]
    [simple     uint 16          block]
    [const      uint 8           unknown1 0x01]
    [simple     uint 16          baseOffset]
    [simple     uint 8           offset]
    [optional   uint 16          arrayLength 'isArray != 0']
]

// Variable write request reference — identifies a variable to write with data.
// The dataSizeIndex is a size index (not a byte count): 1→1B, 2→2B, 3→4B, 4→8B.
// The formula 2^(index-1) converts the index to the actual byte size.
[type VariableWriteRequestReference
    [simple     uint 4           isArray]
    [simple     uint 4           dataSizeIndex]
    [simple     uint 16          block]
    [simple     uint 16          baseOffset]
    [simple     uint 16          offset]
    [optional   uint 16          arrayLength 'isArray != 0']
    [virtual    uint 16          dataSize    'STATIC_CALL("writeSizeIndexToByteCount", dataSizeIndex)']
    [array      byte             recordData     length  'isArray == 1 ? dataSize * arrayLength : dataSize']
]

// Unlocated variable reference from the data dictionary
[type UmasUnlocatedVariableReference
    [simple     uint 16          dataType]
    [simple     uint 16          block]
    [simple     uint 32          offset]
    [simple     uint 8           flags]
    [simple     uint 8           unknown4]
    // Name is null-terminated (no length prefix)
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, -1)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, -1)' '(STR_LEN(value) + 1) * 8']
]

// User-Defined Type definition
[type UmasUDTDefinition
    [simple     uint 16          dataType]
    [simple     uint 16          offset]
    [simple     uint 16          unknown5]
    [simple     uint 16          unknown4]
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, -1)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, -1)' '(STR_LEN(value) + 1) * 8']
]

// Datatype reference from the datatype dictionary
[type UmasDatatypeReference
    [simple     uint 16          dataSize]
    [simple     uint 16          unknown1]
    [simple     uint 8           classIdentifier]
    [simple     uint 8           dataType]
    // Padding byte before the null-terminated name (always 0x00 in observed captures)
    [reserved   uint 8           '0x00']
    // Name is null-terminated (no length prefix) — uses -1 to signal "read until null"
    [manual vstring value  'STATIC_CALL("parseTerminatedString", readBuffer, -1)' 'STATIC_CALL("serializeTerminatedString", writeBuffer, value, -1)' '(STR_LEN(value) + 1) * 8']
]

// Memory bank identification in PLC identification responses
[type PlcMemoryBlockIdent
    [simple uint 8 blockType]
    [simple uint 8 folio]
    [simple uint 16 status]
    [simple uint 32 memoryLength]
]

////////////////////////////////////////////////////////////////
// Data Type Mapping (DataIO)
////////////////////////////////////////////////////////////////

// Maps UMAS data types to their wire format representations
[dataIo DataItem(UmasDataType dataType, uint 16 numberOfValues)
    [typeSwitch dataType,numberOfValues
        ['BOOL','1'  BOOL
            [reserved uint 7 '0x0000'                         ]
            [simple   bit     value                           ]
        ]
        ['EBOOL','1'  BOOL
            [reserved uint 7 '0x0000'                         ]
            [simple   bit     value                            ]
        ]
        ['BYTE','1'  BYTE
            [simple   uint 8  value]
        ]
        ['BYTE' List
            [array    uint 8  value count 'numberOfValues' ]
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
            [simple uint 32 milliseconds]
        ]
        ['TIME' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['DATE','1' DATE
            [simple uint 8 day encoding='"BCD"']
            [simple uint 8 month encoding='"BCD"']
            [simple uint 16 year encoding='"BCD"']
        ]
        ['TOD','1' TIME_OF_DAY
            [simple uint 8 centiseconds encoding='"BCD"']
            [simple uint 8 seconds encoding='"BCD"']
            [simple uint 8 minutes encoding='"BCD"']
            [simple uint 8 hours encoding='"BCD"']
        ]
        ['TOD' List
            [array uint 32 value count 'numberOfValues']
        ]
        ['DATE_AND_TIME','1' DATE_AND_TIME
            [reserved uint 8 '0x00']
            [simple uint 8 seconds encoding='"BCD"']
            [simple uint 8 minutes encoding='"BCD"']
            [simple uint 8 hour encoding='"BCD"']
            [simple uint 8 day encoding='"BCD"']
            [simple uint 8 month encoding='"BCD"']
            [simple uint 16 year encoding='"BCD"']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// Enumerations
////////////////////////////////////////////////////////////////

// Schneider PLC data types with size metadata for variable read/write operations
[enum uint 8 UmasDataType(uint 8 dataTypeSize, uint 8 requestSize, vstring data_type_conversion)
    ['1' BOOL       ['1','1','"BOOL"']]
    ['2' UNKNOWN2   ['1','1','"BOOL"']]
    ['3' UNKNOWN3   ['1','1','"BOOL"']]
    ['4' INT        ['2', '2','"INT"']]
    ['5' UINT       ['2','2','"UINT"']]
    ['6' DINT       ['4','3','"DINT"']]
    ['7' UDINT      ['4','3','"UDINT"']]
    ['8' REAL       ['4','3','"REAL"']]
    ['9' STRING     ['1','17','"STRING"']]
    ['10' TIME      ['4','3','"TIME"']]
    ['11' UNKNOWN11 ['1','1','"BYTE"']]
    ['12' UNKNOWN12 ['1','1','"BYTE"']]
    ['13' UNKNOWN13 ['1','1','"BYTE"']]
    ['14' DATE      ['4','3','"DATE"']]
    ['15' TOD       ['4','3','"TIME_OF_DAY"']]
    ['16' DATE_AND_TIME ['8','4','"DATE_AND_TIME"']]
    ['17' UNKNOWN17 ['1','1','"BYTE"']]
    ['18' UNKNOWN18 ['1','1','"BYTE"']]
    ['19' UNKNOWN19 ['1','1','"BYTE"']]
    ['20' UNKNOWN20 ['1','1','"BYTE"']]
    ['21' BYTE      ['1','1','"BYTE"']]
    ['22' WORD      ['2','2','"WORD"']]
    ['23' DWORD     ['4','3','"DWORD"']]
    ['24' UNKNOWN24 ['1','1','"BYTE"']]
    ['25' EBOOL     ['1','1','"BOOL"']]
]

// Action codes for MonitorPLC register sub-operation (type 0x05)
[enum uint 8 MonitorPlcRegisterAction
    ['0x01' DEREGISTER]
    ['0x02' REGISTER  ]
]

// Standard Modbus error codes used in UMAS error responses
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

// PLC memory block types reported in PlcMemoryBlockIdent (UmasPDUPlcIdentResponse)
[enum uint 8 PlcMemoryBlockType
    ['0x01' RAM_CPU ]
    ['0x04' SD_CARD ]
]

// Data dictionary record types used in FC 0x26 requests (recordType field)
[enum uint 16 UmasDataDictionaryRecordType
    ['0xDD02' VARIABLE_AND_UDT_DEFINITIONS ]
    ['0xDD03' DATATYPE_DEFINITIONS         ]
]

// Type class identifiers returned in datatype table responses (FC 0x26, recordType 0xDD03)
[enum uint 8 UmasTypeClass
    ['0x02' STRUCT ]
    ['0x04' ARRAY  ]
]

////////////////////////////////////////////////////////////////
// Data Dictionary Helper Types
//
// These types model the internal structure of data dictionary
// (FC 0x26) responses. They are parsed at the driver layer from
// the raw byte arrays returned by UmasPDUReadUnlocatedVariableResponse.
////////////////////////////////////////////////////////////////

// Array type definition — returned inside DD02 responses for array types.
// The classId distinguishes arrays (0x04) from structs (0x02).
[type UmasArrayTypeDefinition byteOrder='"LITTLE_ENDIAN"'
    // Always 0x04 for arrays
    [simple     uint 8           classId]

    // Element data type — primitive UmasDataType ID for simple arrays,
    // or a custom type table ID (>= 0x1A) for arrays of structs
    [simple     uint 16          elementTypeId]

    // Number of array dimensions (1D, 2D, or 3D observed in captures)
    [simple     uint 8           numberOfDimensions]

    // Bounds for each dimension (inclusive start and upper bound)
    [array      UmasArrayDimension dimensions count 'numberOfDimensions']
]

// A single array dimension with inclusive start and upper bounds.
// Element count for this dimension = upperBound - startIndex + 1
[type UmasArrayDimension byteOrder='"LITTLE_ENDIAN"'
    [simple     uint 32          startIndex]
    [simple     uint 32          upperBound]
]
