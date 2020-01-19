//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

// Remark: The different fields are encoded in Big-endian.

[type 'ModbusConstants'
    [const          uint 16     'modbusTcpDefaultPort' '502']
]

[type 'ModbusTcpADU' [bit 'response']
    // It is used for transaction pairing, the MODBUS server copies in the response the transaction
    // identifier of the request.
    [simple         uint 16     'transactionIdentifier']

    // It is used for intra-system multiplexing. The MODBUS protocol is identified by the value 0.
    [const          uint 16     'protocolIdentifier'    '0x0000']

    // The length field is a byte count of the following fields, including the Unit Identifier and
    // data fields.
    [simple         uint 16     'length']

    // This field is used for intra-system routing purpose. It is typically used to communicate to
    // a MODBUS+ or a MODBUS serial line slave through a gateway between an Ethernet TCP-IP network
    // and a MODBUS serial line. This field is set by the MODBUS Client in the request and must be
    // returned with the same value in the response by the server.
    [simple         uint 8      'unitItendifier']

    // The actual modbus payload
    [simple         ModbusPDU   'pdu' ['response']]
]

[type 'ModbusSerialADU' [bit 'response']
    [simple         uint 16     'transactionId']
    [reserved       uint 16     '0x0000']
    [simple         uint 16     'length']
    [simple         uint 8      'address']

    // The actual modbus payload
    [simple         ModbusPDU   'pdu' ['response']]
]

[discriminatedType 'ModbusPDU' [bit 'response']
    [implicit       uint 7      'function'  'DISCRIMINATOR_VALUES[1]']
    [implicit       bit         'error'     'DISCRIMINATOR_VALUES[0]']
    [typeSwitch 'error','function','response'
        ['true'                     ModbusPduError
            [simple     uint 8      'exceptionCode']
        ]

        // Bit Access
        ['false','0x02','false'     ModbusPduReadDiscreteInputsRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]
        ['false','0x02','true'      ModbusPduReadDiscreteInputsResponse
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x01','false'     ModbusPduReadCoilsRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]
        ['false','0x01','true'      ModbusPduReadCoilsResponse
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x05','false'     ModbusPduWriteSingleCoilRequest
            [simple     uint 16     'address']
            [simple     uint 16     'value']
        ]
        ['false','0x05','true'      ModbusPduWriteSingleCoilResponse
            [simple     uint 16     'address']
            [simple     uint 16     'value']
        ]

        ['false','0x0F','false'     ModbusPduWriteMultipleCoilsRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]
        ['false','0x0F','true'      ModbusPduWriteMultipleCoilsResponse
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]

        // Uint 16 Access (short)
        ['false','0x04','false'     ModbusPduReadInputRegistersRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]
        ['false','0x04','true'      ModbusPduReadInputRegistersResponse
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x03','false'     ModbusPduReadHoldingRegistersRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]
        ['false','0x03','true'      ModbusPduReadHoldingRegistersResponse
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x06','false'     ModbusPduWriteSingleRegisterRequest
            [simple     uint 16     'address']
            [simple     uint 16     'value']
        ]
        ['false','0x06','true'      ModbusPduWriteSingleRegisterResponse
            [simple     uint 16     'address']
            [simple     uint 16     'value']
        ]

        ['false','0x10','false'     ModbusPduWriteMultipleRegistersRequest
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]
        ['false','0x10','true'      ModbusPduWriteMultipleRegistersResponse
            [simple     uint 16     'startingAddress']
            [simple     uint 16     'quantity']
        ]

        ['false','0x17','false'     ModbusPduReadWriteMultipleRegistersRequest
            [simple     uint 16     'readStartingAddress']
            [simple     uint 16     'readQuantity']
            [simple     uint 16     'writeStartingAddress']
            [simple     uint 16     'writeQuantity']
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]
        ['false','0x17','true'      ModbusPduReadWriteMultipleRegistersResponse
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x16','false'     ModbusPduMaskWriteRegisterRequest
            [simple     uint 16     'referenceAddress']
            [simple     uint 16     'andMask']
            [simple     uint 16     'orMask']
        ]
        ['false','0x16','true'      ModbusPduMaskWriteRegisterResponse
            [simple     uint 16     'referenceAddress']
            [simple     uint 16     'andMask']
            [simple     uint 16     'orMask']
        ]

        ['false','0x18','false'     ModbusPduReadFifoQueueRequest
            [simple     uint 16     'fifoPointerAddress']
        ]
        ['false','0x18','true'      ModbusPduReadFifoQueueResponse
            [implicit   uint 16     'byteCount'     '(COUNT(fifoValue) * 2) + 2']
            [implicit   uint 16     'fifoCount'     '(COUNT(fifoValue) * 2) / 2']
            [array      uint 16     'fifoValue'     count   'fifoCount']
        ]

        // File Record Access
        ['false','0x14','false'     ModbusPduReadFileRecordRequest
            [implicit   uint 8      'byteCount'                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPduReadFileRecordRequestItem      'items' length 'byteCount']
        ]
        ['false','0x14','true'      ModbusPduReadFileRecordResponse
            [implicit   uint 8      'byteCount'                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPduReadFileRecordResponseItem     'items' length 'byteCount']
        ]

        ['false','0x15','false'     ModbusPduWriteFileRecordRequest
            [implicit   uint 8      'byteCount'                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPduWriteFileRecordRequestItem     'items' length 'byteCount']
        ]
        ['false','0x15','true'      ModbusPduWriteFileRecordResponse
            [implicit   uint 8      'byteCount'                 'ARRAY_SIZE_IN_BYTES(items)']
            [array      ModbusPduWriteFileRecordResponseItem    'items' length 'byteCount']
        ]

        // Diagnostics (Serial Line Only)
        ['false','0x07','false'     ModbusPduReadExceptionStatusRequest
        ]
        ['false','0x07','true'      ModbusPduReadExceptionStatusResponse
            [simple     uint 8      'value']
        ]

        ['false','0x08','false'     ModbusPduDiagnosticRequest
            // TODO: Implement the sub-request discriminated type [simple uint 8  'subfunction']
        ]
        ['false','0x08','true'      ModbusPduDiagnosticResponse
            // TODO: Implement the sub-request discriminated type [simple uint 8  'subfunction']
        ]

        ['false','0x0B','false'     ModbusPduGetComEventCounterRequest
        ]
        ['false','0x0B','true'      ModbusPduGetComEventCounterResponse
            [simple     uint 16     'status']
            [simple     uint 16     'eventCount']
        ]

        ['false','0x0C','false'     ModbusPduGetComEventLogRequest
        ]
        ['false','0x0C','true'      ModbusPduGetComEventLogResponse
            [implicit   uint 8      'byteCount'    'COUNT(events) + 6']
            [simple     uint 16     'status']
            [simple     uint 16     'eventCount']
            [simple     uint 16     'messageCount']
            [array      int 8       'events'       count   'byteCount - 6']
        ]

        ['false','0x11','false'     ModbusPduReportServerIdRequest
        ]
        ['false','0x11','true'      ModbusPduReportServerIdResponse
            // TODO: This is not specified very well in the spec ... investigate.
            [implicit   uint 8      'byteCount'     'COUNT(value)']
            [array      int 8       'value'         count   'byteCount']
        ]

        ['false','0x2B','false'     ModbusPduReadDeviceIdentificationRequest
        ]
        ['false','0x2B','true'      ModbusPduReadDeviceIdentificationResponse
        ]
    ]
]

[type 'ModbusPduReadFileRecordRequestItem'
    [simple     uint 8     'referenceType']
    [simple     uint 16    'fileNumber']
    [simple     uint 16    'recordNumber']
    [simple     uint 16    'recordLength']
]

[type 'ModbusPduReadFileRecordResponseItem'
    [implicit   uint 8     'dataLength'     '(COUNT(data) * 2) + 1']
    [simple     uint 8     'referenceType']
    [array      uint 16    'data'           length  'dataLength - 1']
]

[type 'ModbusPduWriteFileRecordRequestItem'
    [simple     uint 8     'referenceType']
    [simple     uint 16    'fileNumber']
    [simple     uint 16    'recordNumber']
    [implicit   uint 16    'recordLength'   '(COUNT(recordData) * 2) / 2']
    [array      uint 16    'recordData'     length  'recordLength * 2']
]

[type 'ModbusPduWriteFileRecordResponseItem'
    [simple     uint 8     'referenceType']
    [simple     uint 16    'fileNumber']
    [simple     uint 16    'recordNumber']
    [implicit   uint 16    'recordLength'   '(COUNT(recordData) * 2) / 2']
    [array      uint 16    'recordData'     length  'recordLength * 2']
]
