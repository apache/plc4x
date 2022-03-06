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

// https://modbus.org/docs/Modbus_Application_Protocol_V1_1b.pdf

// Remark: The different fields are encoded in Big-endian.

[type ModbusConstants
    [const          uint 16     modbusTcpDefaultPort 502]
]

[type ModbusTcpADU(bit response) byteOrder='BIG_ENDIAN'
    // It is used for transaction pairing, the MODBUS server copies in the response the transaction
    // identifier of the request.
    [simple         uint 16     transactionIdentifier]

    // It is used for intra-system multiplexing. The MODBUS protocol is identified by the value 0.
    [const          uint 16     protocolIdentifier    0x0000]

    // The length field is a byte count of the following fields, including the Unit Identifier and
    // data fields.
    [implicit       uint 16     length                'pdu.lengthInBytes + 1']

    // This field is used for intra-system routing purpose. It is typically used to communicate to
    // a MODBUS+ or a MODBUS serial line slave through a gateway between an Ethernet TCP-IP network
    // and a MODBUS serial line. This field is set by the MODBUS Client in the request and must be
    // returned with the same value in the response by the server.
    [simple         uint 8      unitIdentifier]

    // The actual modbus payload
    [simple         ModbusPDU('response')   pdu]
]

[type ModbusRtuADU(bit response) byteOrder='LITTLE_ENDIAN'
    // The start is indicated by more than 3,5 chars of value 0x00 ...
    // The protocol will take care of consuming all except the last 4 empty chars
    // The Length is determined by starting at the last 4 empty characters
    // to the start of the next 4 empty characters
    [const          uint 8      space1 0x00           ] // Character '\0'
    [const          uint 8      space2 0x00           ] // Character '\0'
    [const          uint 8      space3 0x00           ] // Character '\0'
    [const          uint 8      space4 0x00           ] // Character '\0'
    [simple         uint 8      address]

    // The actual modbus payload
    [simple         ModbusPDU('response')   pdu]
    //[checksum       uint 16     crc            ]

    // The at least 4 silence chars is a separator ... so we don't actually consume them.
    // They will be consumed by the next packet
    //[const          uint 8      0x00           ] // Character '\0'
    //[const          uint 8      0x00           ] // Character '\0'
    //[const          uint 8      0x00           ] // Character '\0'
    //[const          uint 8      0x00           ] // Character '\0'
]

[type ModbusAsciiADU(bit response) byteOrder='BIG_ENDIAN'
    [const          uint 8      start 0x3A           ] // Character ':'

    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
    // Starting the ASCII encoded part where every byte is encoded as the ASCII representation of the raw byte.
    [simple         uint 8      address        ]

    // The actual modbus payload
    [simple         ModbusPDU('response')   pdu]

    //[checksum       uint 8      lrc            ]

    // End the ASCII encoded part ...
    ////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

    [const          uint 8      cr 0x0D           ] // Character '\r'
    [const          uint 8      lf 0x0A           ] // Character '\n'
]

[discriminatedType ModbusPDU(bit response)
    [discriminator bit         errorFlag]
    [discriminator uint 7      functionFlag]
    [typeSwitch errorFlag,functionFlag,response
        ['true'                     ModbusPDUError
            [simple ModbusErrorCode  exceptionCode]
        ]

        // Bit Access
        ['false','0x02','false'     ModbusPDUReadDiscreteInputsRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]
        ['false','0x02','true'      ModbusPDUReadDiscreteInputsResponse
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        ['false','0x01','false'     ModbusPDUReadCoilsRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]
        ['false','0x01','true'      ModbusPDUReadCoilsResponse
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        ['false','0x05','false'     ModbusPDUWriteSingleCoilRequest
            [simple     uint 16     address]
            [simple     uint 16     value]
        ]
        ['false','0x05','true'      ModbusPDUWriteSingleCoilResponse
            [simple     uint 16     address]
            [simple     uint 16     value]
        ]

        ['false','0x0F','false'     ModbusPDUWriteMultipleCoilsRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]
        ['false','0x0F','true'      ModbusPDUWriteMultipleCoilsResponse
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]

        // Uint 16 Access (short)
        ['false','0x04','false'     ModbusPDUReadInputRegistersRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]
        ['false','0x04','true'      ModbusPDUReadInputRegistersResponse
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        ['false','0x03','false'     ModbusPDUReadHoldingRegistersRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]
        ['false','0x03','true'      ModbusPDUReadHoldingRegistersResponse
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        ['false','0x06','false'     ModbusPDUWriteSingleRegisterRequest
            [simple     uint 16     address]
            [simple     uint 16     value]
        ]
        ['false','0x06','true'      ModbusPDUWriteSingleRegisterResponse
            [simple     uint 16     address]
            [simple     uint 16     value]
        ]

        ['false','0x10','false'     ModbusPDUWriteMultipleHoldingRegistersRequest
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]
        ['false','0x10','true'      ModbusPDUWriteMultipleHoldingRegistersResponse
            [simple     uint 16     startingAddress]
            [simple     uint 16     quantity]
        ]

        ['false','0x17','false'     ModbusPDUReadWriteMultipleHoldingRegistersRequest
            [simple     uint 16     readStartingAddress]
            [simple     uint 16     readQuantity]
            [simple     uint 16     writeStartingAddress]
            [simple     uint 16     writeQuantity]
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]
        ['false','0x17','true'      ModbusPDUReadWriteMultipleHoldingRegistersResponse
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        ['false','0x16','false'     ModbusPDUMaskWriteHoldingRegisterRequest
            [simple     uint 16     referenceAddress]
            [simple     uint 16     andMask]
            [simple     uint 16     orMask]
        ]
        ['false','0x16','true'      ModbusPDUMaskWriteHoldingRegisterResponse
            [simple     uint 16     referenceAddress]
            [simple     uint 16     andMask]
            [simple     uint 16     orMask]
        ]

        ['false','0x18','false'     ModbusPDUReadFifoQueueRequest
            [simple     uint 16     fifoPointerAddress]
        ]
        ['false','0x18','true'      ModbusPDUReadFifoQueueResponse
            [implicit   uint 16     byteCount     '(COUNT(fifoValue) * 2) + 2']
            [implicit   uint 16     fifoCount     '(COUNT(fifoValue) * 2) / 2']
            [array      uint 16     fifoValue     count   'fifoCount']
        ]

        // File Record Access
        ['false','0x14','false'     ModbusPDUReadFileRecordRequest
            [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPDUReadFileRecordRequestItem      items length 'byteCount']
        ]
        ['false','0x14','true'      ModbusPDUReadFileRecordResponse
            [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPDUReadFileRecordResponseItem     items length 'byteCount']
        ]

        ['false','0x15','false'     ModbusPDUWriteFileRecordRequest
            [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPDUWriteFileRecordRequestItem     items length 'byteCount']
        ]
        ['false','0x15','true'      ModbusPDUWriteFileRecordResponse
            [implicit   uint 8      byteCount                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPDUWriteFileRecordResponseItem    items length 'byteCount']
        ]

        // Diagnostics (Serial Line Only)
        ['false','0x07','false'     ModbusPDUReadExceptionStatusRequest
        ]
        ['false','0x07','true'      ModbusPDUReadExceptionStatusResponse
            [simple     uint 8      value]
        ]

        ['false','0x08','false'     ModbusPDUDiagnosticRequest
            [simple     uint 16     subFunction]
            [simple     uint 16     data]
        ]
        ['false','0x08','true'      ModbusPDUDiagnosticResponse
            [simple     uint 16     subFunction]
            [simple     uint 16     data]
        ]

        ['false','0x0B','false'     ModbusPDUGetComEventCounterRequest
        ]
        ['false','0x0B','true'      ModbusPDUGetComEventCounterResponse
            [simple     uint 16     status]
            [simple     uint 16     eventCount]
        ]

        ['false','0x0C','false'     ModbusPDUGetComEventLogRequest
        ]
        ['false','0x0C','true'      ModbusPDUGetComEventLogResponse
            [implicit   uint 8      byteCount    'COUNT(events) + 6']
            [simple     uint 16     status]
            [simple     uint 16     eventCount]
            [simple     uint 16     messageCount]
            [array      byte        events       count   'byteCount - 6']
        ]

        ['false','0x11','false'     ModbusPDUReportServerIdRequest
        ]
        ['false','0x11','true'      ModbusPDUReportServerIdResponse
            // TODO: This is not specified very well in the spec ... investigate.
            [implicit   uint 8      byteCount     'COUNT(value)']
            [array      byte        value         count   'byteCount']
        ]

        // Remark: Even if the Modbus spec states that supporting this type of request is mandatory
        // I have not come across a single device that really supported it. Some devices just reacted
        // with an error.
        ['false','0x2B','false'     ModbusPDUReadDeviceIdentificationRequest
            [const  uint 8                       meiType  0x0E]
            [simple ModbusDeviceInformationLevel level        ]
            [simple uint 8                       objectId     ]
        ]
        ['false','0x2B','true'      ModbusPDUReadDeviceIdentificationResponse
            [const    uint 8                                 meiType          0x0E                              ]
            [simple   ModbusDeviceInformationLevel           level                                              ]
            [simple   bit                                    individualAccess                                   ]
            [simple   ModbusDeviceInformationConformityLevel conformityLevel                                    ]
            [simple   ModbusDeviceInformationMoreFollows     moreFollows                                        ]
            [simple   uint 8                                 nextObjectId                                       ]
            [implicit uint 8                                 numberOfObjects  'COUNT(objects)'                  ]
            [array    ModbusDeviceInformationObject          objects          count            'numberOfObjects']
        ]
    ]
]

[type ModbusPDUReadFileRecordRequestItem
    [simple     uint 8     referenceType]
    [simple     uint 16    fileNumber   ]
    [simple     uint 16    recordNumber ]
    [simple     uint 16    recordLength ]
]

[type ModbusPDUReadFileRecordResponseItem
    [implicit   uint 8     dataLength     'COUNT(data) + 1'       ]
    [simple     uint 8     referenceType                          ]
    [array      byte       data           length  'dataLength - 1']
]

[type ModbusPDUWriteFileRecordRequestItem
    [simple     uint 8     referenceType]
    [simple     uint 16    fileNumber]
    [simple     uint 16    recordNumber]
    [implicit   uint 16    recordLength   'COUNT(recordData) / 2'   ]
    [array      byte       recordData     length  'recordLength * 2']
]

[type ModbusPDUWriteFileRecordResponseItem
    [simple     uint 8     referenceType]
    [simple     uint 16    fileNumber]
    [simple     uint 16    recordNumber]
    [implicit   uint 16    recordLength   'COUNT(recordData) / 2']
    [array      byte       recordData     length  'recordLength']
]

[type ModbusDeviceInformationObject
    [simple   uint 8 objectId                                  ]
    [implicit uint 8 objectLength  'COUNT(data)'               ]
    [array    byte   data          count         'objectLength']
]

[dataIo DataItem(ModbusDataType dataType, uint 16 numberOfValues)
    [typeSwitch dataType,numberOfValues
        ['BOOL','1' BOOL
            [reserved uint 15 '0x0000']
            [simple   bit     value   ]
        ]
        ['BOOL' List
            [array bit value count 'numberOfValues']
        ]
        ['BYTE','1' BitString
            [reserved uint 8 '0x00']
            [simple   uint 8 value ]
        ]
        ['BYTE' List
            [array bit value count 'numberOfValues * 8']
        ]
        ['WORD','1' BitString
            [simple uint 16 value]
        ]
        ['WORD' List
            [array bit value count 'numberOfValues * 16']
        ]
        ['DWORD','1' BitString
            [simple uint 32 value]
        ]
        ['DWORD' List
            [array bit value count 'numberOfValues * 32']
        ]
        ['LWORD','1' BitString
            [simple uint 64 value]
        ]
        ['LWORD' List
            [array bit value count 'numberOfValues * 64']
        ]
        ['SINT','1' SINT
            [reserved uint 8 '0x00']
            [simple   int 8  value ]
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
            [simple uint 8 value]
        ]
        ['CHAR' List
            [array uint 8 value count 'numberOfValues']
        ]
        ['WCHAR','1' WCHAR
            [simple uint 16 value]
        ]
        ['WCHAR' List
            [array uint 16 value count 'numberOfValues']
        ]
    ]
]

[enum uint 8 ModbusDataType(uint 8 dataTypeSize)
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

[enum uint 8 ModbusDeviceInformationLevel
    ['0x01' BASIC     ]
    ['0x02' REGULAR   ]
    ['0x03' EXTENDED  ]
    ['0x04' INDIVIDUAL]
]

[enum uint 7 ModbusDeviceInformationConformityLevel
    ['0x01' BASIC_STREAM_ONLY   ]
    ['0x02' REGULAR_STREAM_ONLY ]
    ['0x03' EXTENDED_STREAM_ONLY]
]

[enum uint 8 ModbusDeviceInformationMoreFollows
    ['0x00' NO_MORE_OBJECTS_AVAILABLE]
    ['0xFF' MORE_OBJECTS_AVAILABLE   ]
]
