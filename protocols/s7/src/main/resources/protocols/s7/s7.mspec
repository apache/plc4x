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

////////////////////////////////////////////////////////////////
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type 'TPKTPacket'
    [const    uint 8     'protocolId' '0x03']
    [reserved uint 8     '0x00']
    [implicit uint 16    'len'        'payload.lengthInBytes + 4']
    [simple   COTPPacket 'payload' ['len - 4']]
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType 'COTPPacket' [uint 16 'cotpLen']
    [implicit      uint 8 'headerLength' 'lengthInBytes - (payload.lengthInBytes + 1)']
    [discriminator uint 8 'tpduCode']
    [typeSwitch 'tpduCode'
        ['0xF0' COTPPacketData
            [simple bit    'eot']
            [simple uint 7 'tpduRef']
        ]
        ['0xE0' COTPPacketConnectionRequest
            [simple uint 16 'destinationReference']
            [simple uint 16 'sourceReference']
            [simple uint 8  'protocolClass']
        ]
        ['0xD0' COTPPacketConnectionResponse
            [simple uint 16 'destinationReference']
            [simple uint 16 'sourceReference']
            [simple uint 8  'protocolClass']
        ]
        ['0x80' COTPPacketDisconnectRequest
            [simple uint 16 'destinationReference']
            [simple uint 16 'sourceReference']
            [simple uint 8  'protocolClass']
        ]
        ['0xC0' COTPPacketDisconnectResponse
            [simple uint 16 'destinationReference']
            [simple uint 16 'sourceReference']
        ]
        ['0x70' COTPPacketTpduError
            [simple uint 16 'destinationReference']
            [simple uint 8  'rejectCause']
        ]
    ]
    [array    COTPParameter 'parameters' length '(headerLength + 1) - curPos' ['(headerLength + 1) - curPos']]
    [optional S7Message     'payload'    'curPos < cotpLen']
]

[discriminatedType 'COTPParameter' [uint 8 'rest']
    [discriminator uint 8 'parameterType']
    [implicit      uint 8 'parameterLength' 'lengthInBytes - 2']
    [typeSwitch 'parameterType'
        ['0xC0' COTPParameterTpduSize
            [simple uint 8 'tpduSize']
        ]
        ['0xC1' COTPParameterCallingTsap
            [simple uint 16 'tsapId']
        ]
        ['0xC2' COTPParameterCalledTsap
            [simple uint 16 'tsapId']
        ]
        ['0xC3' COTPParameterChecksum
            [simple uint 8 'crc']
        ]
        ['0xE0' COTPParameterDisconnectAdditionalInformation
            [array  uint 8 'data' count 'rest']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// S7
////////////////////////////////////////////////////////////////

[discriminatedType 'S7Message'
    [const         uint 8  'protocolId'      '0x32']
    [discriminator uint 8  'messageType']
    [reserved      uint 16 '0x0000']
    [simple        uint 16 'tpduReference']
    [implicit      uint 16 'parameterLength' 'parameter.lengthInBytes']
    [implicit      uint 16 'payloadLength'   'payload.lengthInBytes']
    [typeSwitch 'messageType'
        ['0x01' S7MessageRequest
        ]
        ['0x03' S7MessageResponse
            [simple uint 8 'errorClass']
            [simple uint 8 'errorCode']
        ]
        ['0x07' S7MessageUserData
        ]
    ]
    [simple S7Parameter 'parameter' ['messageType']]
    [simple S7Payload   'payload'   ['messageType', 'parameter']]
]

////////////////////////////////////////////////////////////////
// Parameters

[discriminatedType 'S7Parameter' [uint 8 'messageType']
    [discriminator uint 8 'parameterType']
    [typeSwitch 'parameterType','messageType'
        ['0xF0' S7ParameterSetupCommunication
            [reserved uint 8  '0x00']
            [simple   uint 16 'maxAmqCaller']
            [simple   uint 16 'maxAmqCallee']
            [simple   uint 16 'pduLength']
        ]
        ['0x04','0x01' S7ParameterReadVarRequest
            [implicit uint 8                    'numItems' 'COUNT(items)']
            [array    S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x04','0x03' S7ParameterReadVarResponse
            [simple uint 8 'numItems']
        ]
        ['0x05','0x01' S7ParameterWriteVarRequest
            [implicit uint 8                    'numItems' 'COUNT(items)']
            [array    S7VarRequestParameterItem 'items'    count 'numItems']
        ]
        ['0x05','0x03' S7ParameterWriteVarResponse
            [simple uint 8 'numItems']
        ]
        ['0x00','0x07' S7ParameterUserData
            [implicit uint 8       'numItems' 'COUNT(items)']
            [array    UserDataItem 'items' count 'numItems']
        ]
    ]
]

[discriminatedType 'S7VarRequestParameterItem'
    [discriminator uint 8 'parameterItemType']
    [typeSwitch 'parameterItemType'
        ['0x12' S7VarRequestParameterItemAddress
            [implicit uint 8    'addressLength' 'address.lengthInBytes']
            [simple   S7Address 'address']
        ]
    ]
]

[discriminatedType 'S7Address'
    [discriminator uint 8 'addressType']
    [typeSwitch 'addressType'
        ['0x10' S7AddressAny
            [simple   uint 8  'transportSize']
            [simple   uint 16 'numberOfElements']
            [simple   uint 16 'dbNumber']
            [simple   uint 8  'area']
            [reserved uint 5  '0x00']
            [simple   uint 16 'byteAddress']
            [simple   uint 3  'bitAddress']
        ]
    ]
]

// TODO: CPUFunctions still need some love ...
[discriminatedType 'UserDataItem'
    [discriminator uint 8 'itemType']
    [typeSwitch 'itemType'
        ['0x12' UserDataItemCPUFunctions
            [implicit uint 8  'parameterLength' 'lengthInBytes']
            [simple   uint 16 'cpuFunctionType']
            [simple   uint 8  'subFunctionGroup']
            [simple   uint 8  'sequenceNumber']
            [optional uint 8  'dataUnitReferenceNumber' 'parameterLength == 8']
            [optional uint 8  'lastDataUnit' 'parameterLength == 8']
            [optional uint 8  'errorCode' 'parameterLength == 8']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatedType 'S7Payload' [uint 8 'messageType', S7Parameter 'parameter']
    [typeSwitch 'parameter.discriminatorValues[0]', 'messageType'
        ['0xF0' S7PayloadSetupCommunication]
        ['0x04','0x01' S7PayloadReadVarRequest]
        ['0x04','0x03' S7PayloadReadVarResponse
            [array S7VarPayloadDataItem 'items' count 'CAST(parameter, S7ParameterReadVarResponse).numItems']
        ]
        ['0x05','0x01' S7PayloadWriteVarRequest
            [array S7VarPayloadDataItem 'items' count 'COUNT(CAST(parameter, S7ParameterWriteVarRequest).items)']
        ]
        ['0x05','0x03' S7PayloadWriteVarResponse
            [array S7VarPayloadStatusItem 'items' count 'CAST(parameter, S7ParameterWriteVarResponse).numItems']
        ]
        ['0x00','0x07' S7PayloadUserData
        ]
    ]
]

// This is actually not quite correct as depending pon the transportSize the length is either defined in bits or bytes.
[type 'S7VarPayloadDataItem'
    [simple  uint 8  'returnCode']
    [simple  uint 8  'transportSize']
    [simple  uint 16 'dataLength']
    [array   uint 8  'data' count 'dataLength / 8']
    [padding uint 8  'pad' '0x00' '(dataLength / 8) % 2 == 1']
]

[type 'S7VarPayloadStatusItem'
    [simple uint 8 'returnCode']
]

[enum uint 8 'DataTransportSize' [bit 'sizeInBits']
    ['0x00' NULL            ['false']]
    ['0x03' BIT             ['true']]
    ['0x04' BYTE_WORD_DWORD ['true']]
    ['0x05' INTEGER         ['true']]
    ['0x06' DINTEGER        ['false']]
    ['0x07' REAL            ['false']]
    ['0x09' OCTET_STRING    ['false']]
]

[enum uint 8 'TransportSize' [uint 8 'sizeCode', uint 8 'sizeInBytes', TransportSize 'baseType', DataTransportSize 'dataTransportSize']
    ['0x01' BOOL             ['X'              , '1'                 , 'null'                  , 'DataTransportSize.BIT']]
    ['0x02' BYTE             ['B'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x04' WORD             ['W'              , '2'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x06' DWORD            ['D'              , '4'                 , 'WORD'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x00' LWORD            ['X'              , '8'                 , 'null'                  , 'null']]
    ['0x05' INT              ['W'              , '2'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x05' UINT             ['W'              , '2'                 , 'INT'                   , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x02' SINT             ['B'              , '1'                 , 'INT'                   , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x02' USINT            ['B'              , '1'                 , 'INT'                   , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x07' DINT             ['D'              , '4'                 , 'INT'                   , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x07' UDINT            ['D'              , '4'                 , 'INT'                   , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x00' LINT             ['X'              , '8'                 , 'INT'                   , 'null']]
    ['0x00' ULINT            ['X'              , '16'                , 'INT'                   , 'null']]
    ['0x08' REAL             ['D'              , '4'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x00' LREAL            ['X'              , '8'                 , 'REAL'                  , 'null']]
    ['0x0B' TIME             ['X'              , '4'                 , 'null'                  , 'null']
    ['0x00' LTIME            ['X'              , '8'                 , 'TIME'                  , 'null']]
    ['0x02' DATE             ['X'              , '2'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x02' TIME_OF_DAY      ['X'              , '4'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x02' DATE_AND_TIME    ['X'              , '8'                 , 'null'                  , 'null']]
    ['0x03' CHAR             ['B'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x13' WCHAR            ['X'              , '2'                 , 'null'                  , 'null']]
    ['0x03' STRING           ['X'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD']]
    ['0x00' WSTRING          ['X'              , '1'                 , 'null'                  , 'null']]
]

