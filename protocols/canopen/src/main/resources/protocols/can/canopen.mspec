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

[enum uint 4 CANOpenService(uint 8 'min', uint 8 'max', bit 'pdo')
    ['0b0000' NMT             ['0',          '0'    ,      'false']]
    ['0b0001' SYNC            ['0x80',       '0x80' ,      'false']]
    ['0b0001' EMCY            ['0x81',       '0xFF' ,      'false']]
    ['0b0010' TIME            ['0x100',      '0x100',      'false']]
    ['0b0011' TRANSMIT_PDO_1  ['0x180',      '0x1FF',      'true' ]]
    ['0b0100' RECEIVE_PDO_1   ['0x200',      '0x27F',      'true' ]]
    ['0b0101' TRANSMIT_PDO_2  ['0x280',      '0x2FF',      'true' ]]
    ['0b0110' RECEIVE_PDO_2   ['0x300',      '0x37F',      'true' ]]
    ['0b0111' TRANSMIT_PDO_3  ['0x380',      '0x3FF',      'true' ]]
    ['0b1000' RECEIVE_PDO_3   ['0x400',      '0x47F',      'true' ]]
    ['0b1001' TRANSMIT_PDO_4  ['0x480',      '0x4FF',      'true' ]]
    ['0b1010' RECEIVE_PDO_4   ['0x500',      '0x57F',      'true' ]]
    ['0b1011' TRANSMIT_SDO    ['0x580',      '0x5FF',      'false']]
    ['0b1100' RECEIVE_SDO     ['0x600',      '0x67F',      'false']]
    ['0b1110' HEARTBEAT       ['0x700',      '0x77F',      'false']]
]

[enum uint 8 NMTStateRequest
    ['0x01' START]
    ['0x02' STOP]
    ['0x80' PRE_OPERATIONAL]
    ['0x81' RESET_NODE]
    ['0x82' RESET_COMMUNICATION]
]

[enum uint 8 NMTState
    ['0x00' BOOTED_UP]
    ['0x04' STOPPED]
    ['0x05' OPERATIONAL]
    ['0x7f' PRE_OPERATIONAL]
]

[discriminatedType CANOpenPayload(CANOpenService 'service')
    [typeSwitch 'service'
        ['NMT' CANOpenNetworkPayload
            [simple NMTStateRequest 'request']
            [reserved uint 1 '0x00']
            [simple uint 7 'node']
        ]
        ['TIME' CANOpenTimeSynchronization
            [simple CANOpenTime 'timeOfDay']
        ]
        ['RECEIVE_PDO_1' CANOpenPDOPayload
            [simple CANOpenPDO/*('1', 'true')*/ 'pdo']
        ]
        ['TRANSMIT_PDO_1' CANOpenPDOPayload
            [simple CANOpenPDO/*('1', 'false')*/ 'pdo']
        ]
        ['RECEIVE_PDO_2' CANOpenPDOPayload
            [simple CANOpenPDO/*('2', 'true')*/ 'pdo']
        ]
        ['TRANSMIT_PDO_2' CANOpenPDOPayload
            [simple CANOpenPDO/*('1', 'false')*/ 'pdo']
        ]
        ['RECEIVE_PDO_3' CANOpenPDOPayload
            [simple CANOpenPDO/*('3', 'true')*/ 'pdo']
        ]
        ['TRANSMIT_PDO_3' CANOpenPDOPayload
            [simple CANOpenPDO/*('1', 'false')*/ 'pdo']
        ]
        ['RECEIVE_PDO_4' CANOpenPDOPayload
            [simple CANOpenPDO/*('4', 'true')*/ 'pdo']
        ]
        ['TRANSMIT_PDO_4' CANOpenPDOPayload
            [simple CANOpenPDO/*('1', 'false')*/ 'pdo']
        ]
        ['RECEIVE_SDO' CANOpenSDORequest
            [simple SDORequestCommand 'command']
            [simple SDORequest('command') 'request']
        ]
        ['TRANSMIT_SDO' CANOpenSDOResponse
            [simple SDOResponseCommand 'command']
            [simple SDOResponse('command') 'response']
        ]
        ['HEARTBEAT' CANOpenHeartbeatPayload
            [simple NMTState 'state']
        ]
    ]
]

[type SDORequest(SDORequestCommand 'command')
    [typeSwitch 'command'
        ['SEGMENT_DOWNLOAD' SDOSegmentDownloadRequest
            [simple bit 'toggle']
            [implicit uint 3 'size' '7 - COUNT(data)']
            [simple bit 'last']
            [array byte 'data' count '7 - size']
            [padding uint 8 'alignment' '0x00' '7 - COUNT(data)']
        ]
        ['INITIATE_DOWNLOAD' SDOInitiateDownloadRequest
            [reserved uint 1 '0x00']
            [implicit uint 2 'size' 'STATIC_CALL("count", expedited, indicated, payload)']
            [simple bit 'expedited']
            [simple bit 'indicated']
            [simple IndexAddress 'address']
            [simple SDOInitiateUploadResponsePayload('expedited', 'indicated', 'size') 'payload']
        ]
        ['INITIATE_UPLOAD' SDOInitiateUploadRequest
            [reserved uint 5 '0x00']
            [simple IndexAddress 'address']
            [reserved int 32 '0x00'] // padding
        ]
        ['SEGMENT_UPLOAD' SDOSegmentUploadRequest
            [simple bit 'toggle']
            [reserved uint 4 '0x00']
            [reserved int 56 '0x00'] // padding
        ]
        ['ABORT' SDOAbortRequest
            [simple SDOAbort 'abort']
        ]
        ['BLOCK' SDOBlockRequest
            [simple SDOBlockData 'block']
        ]
    ]
]

[type SDOBlockData
    [simple uint 5 'flags']
    [array byte 'data' count '7']
]

[type SDOResponse(SDOResponseCommand 'command')
    [typeSwitch 'command'
        ['SEGMENT_UPLOAD' SDOSegmentUploadResponse
            [simple bit 'toggle']
            [implicit uint 3 'size' '7 - COUNT(data)']
            [simple bit 'last']
            [array byte 'data' count '7 - size']
            [padding uint 8 'alignment' '0x00' '7 - COUNT(data)']
        ]
        ['SEGMENT_DOWNLOAD' SDOSegmentDownloadResponse
            [simple bit 'toggle']
            [reserved uint 4 '0x00'] // fill first byte
            [reserved int 56 '0x00'] // padding
        ]
        ['INITIATE_UPLOAD' SDOInitiateUploadResponse
            [reserved uint 1 '0x00']
            [implicit uint 2 'size' 'STATIC_CALL("count", expedited, indicated, payload)']
            [simple bit 'expedited']
            [simple bit 'indicated']
            [simple IndexAddress 'address']
            [simple SDOInitiateUploadResponsePayload('expedited', 'indicated', 'size') 'payload']
        ]
        ['INITIATE_DOWNLOAD' SDOInitiateDownloadResponse
            [reserved uint 5 '0x00']
            [simple IndexAddress 'address']
            [reserved int 32 '0x00'] // padding
        ]
        ['ABORT' SDOAbortResponse
            [simple SDOAbort 'abort']
        ]
        ['BLOCK' SDOBlockResponse
            [simple SDOBlockData 'block']
        ]
    ]
]

[type SDOInitiateUploadResponsePayload(bit 'expedited', bit 'indicated', uint 2 'size')
    [typeSwitch 'expedited', 'indicated'
        ['true', 'true' SDOInitiateExpeditedUploadResponse
            [array byte 'data' count '4 - size']
            [padding uint 8 'alignment' '0x00' '4 - COUNT(data)']
        ]
        ['false', 'true' SDOInitiateSegmentedUploadResponse
            [simple uint 32 'bytes']
        ]
        ['false', 'false' SDOInitiateSegmentedReservedResponse
            [reserved int 32 '0x00']
        ]
    ]
]

[type SDOAbort
    [reserved uint 5 '0x00']
    [simple IndexAddress 'address']
    [simple uint 32 'code']
]

[type SDOSegment
    [reserved uint 1 '0x00']
    [implicit uint 2 'size' 'expedited && indicated ? 4 - COUNT(data) : 0']
    [simple bit 'expedited']
    [simple bit 'indicated']
    [simple IndexAddress 'address']
    [array byte 'data' count '(expedited && indicated) ? 4 - size : 0']
    [padding uint 8 'alignment' '0x00' '4 - (COUNT(data))']
]

[type IndexAddress
    [simple uint 16 'index']
    [simple uint 8 'subindex']
]

[enum uint 3 SDORequestCommand
    ['0x00' SEGMENT_DOWNLOAD  ]
    ['0x01' INITIATE_DOWNLOAD ]
    ['0x02' INITIATE_UPLOAD   ]
    ['0x03' SEGMENT_UPLOAD    ]
    ['0x04' ABORT             ]
    ['0x05' BLOCK             ]
]

[enum uint 3 SDOResponseCommand
    ['0x00' SEGMENT_UPLOAD    ]
    ['0x01' SEGMENT_DOWNLOAD  ]
    ['0x02' INITIATE_UPLOAD   ]
    ['0x03' INITIATE_DOWNLOAD ]
    ['0x04' ABORT             ]
    ['0x06' BLOCK             ]
]

[type CANOpenPDO /*(uint 2 'index', bit 'receive')*/
    [array byte 'data' count '8']
]

[type CANOpenTime
    // CiA 301 - section 7.1.6.5 and 7.1.6.6
    [simple uint 28 'millis']
    [reserved int 4 '0x00']
    [simple uint 16 'days']
]

[enum 'CANOpenDataType' (uint 8 'numBits')
    [BOOLEAN     [ '1'] ]
    [UNSIGNED8   [ '8'] ]
    [UNSIGNED16  ['16'] ]
    [UNSIGNED24  ['24'] ]
    [UNSIGNED32  ['32'] ]
    [UNSIGNED40  ['40'] ]
    [UNSIGNED48  ['48'] ]
    [UNSIGNED56  ['56'] ]
    [UNSIGNED64  ['64'] ]
    [INTEGER8    [ '8'] ]
    [INTEGER16   ['16'] ]
    [INTEGER24   ['24'] ]
    [INTEGER32   ['32'] ]
    [INTEGER40   ['40'] ]
    [INTEGER48   ['48'] ]
    [INTEGER56   ['56'] ]
    [INTEGER64   ['64'] ]
    [REAL32      ['32'] ]
    [REAL64      ['64'] ]

    // compound/complex types
    [RECORD           [ '8'] ]
    [OCTET_STRING     [ '8'] ]
    [VISIBLE_STRING   [ '8'] ]
    [UNICODE_STRING   ['16'] ]
    [TIME_OF_DAY      ['48'] ]
    [TIME_DIFFERENCE  ['48'] ]
]

[dataIo DataItem(CANOpenDataType 'dataType', int 32 'size')
    [typeSwitch 'dataType'
        ['BOOLEAN' BOOL
            [simple bit 'value']
        ]
        ['UNSIGNED8' USINT
            [simple uint 8 'value']
        ]
        ['UNSIGNED16' UINT
            [simple uint 16 'value']
        ]
        ['UNSIGNED24' UDINT
            [simple uint 24 'value']
        ]
        ['UNSIGNED32' UDINT
            [simple uint 32 'value']
        ]
        ['UNSIGNED40' ULINT
            [simple uint 40 'value']
        ]
        ['UNSIGNED48' ULINT
            [simple uint 48 'value']
        ]
        ['UNSIGNED56' ULINT
            [simple uint 56 'value']
        ]
        ['UNSIGNED64' ULINT
            [simple uint 64 'value']
        ]
        ['INTEGER8' SINT
            [simple int 8 'value']
        ]
        ['INTEGER16' INT
            [simple int 16 'value']
        ]
        ['INTEGER24' DINT
            [simple int 24 'value']
        ]
        ['INTEGER32' DINT
            [simple int 32 'value']
        ]
        ['INTEGER40' LINT
            [simple int 40 'value']
        ]
        ['INTEGER48' LINT
            [simple int 48 'value']
        ]
        ['INTEGER56' LINT
            [simple int 56 'value']
        ]
        ['INTEGER64' LINT
            [simple int 64 'value']
        ]
        ['REAL32' REAL
            [simple float 32 'value']
        ]
        ['REAL64' LREAL
            [simple float 64 'value']
        ]
        ['RECORD' List(int 32 'size')
            [array byte 'value' length 'size']
        ]
        ['OCTET_STRING' STRING
           [simple vstring 'size' 'test']
        ]
        ['VISIBLE_STRING' STRING
            [simple vstring 'size' 'value']
        ]
        //CANOpenDataType.TIME_OF_DAY' CANOpenTime
        //CANOpenDataType.TIME_DIFFERENCE' CANOpenTime
        ['UNICODE_STRING' STRING
            [simple vstring 'size/8' 'value']
        ]
    ]
]

// utility type quickly write data for mapped/manufacturer PDOs
[type CANOpenMPDO
    [simple uint 8 'node']
    [simple IndexAddress 'address']
    [array byte 'data' count '4']
]

// A compact, byte aligned structure for test and embedding purposes
[type CANOpenFrame byteOrder='"LITTLE_ENDIAN"'
    [simple uint 8 'nodeId']
    [simple CANOpenService 'service']
    [const uint 4 'alignment' '0x00']
    [simple CANOpenPayload('service') 'payload']
    [padding uint 8 'alignment' '0x00' '8 - (payload.lengthInBytes)']
]
