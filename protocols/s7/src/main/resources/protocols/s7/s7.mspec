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
    [implicit      uint 8 'headerLength' 'lengthInBytes - (((payload != null) ? payload.lengthInBytes : 0) + 1)']
    [discriminator uint 8 'tpduCode']
    [typeSwitch 'tpduCode'
        ['0xF0' COTPPacketData
            [simple bit    'eot']
            [simple uint 7 'tpduRef']
        ]
        ['0xE0' COTPPacketConnectionRequest
            [simple uint 16           'destinationReference']
            [simple uint 16           'sourceReference']
            [enum   COTPProtocolClass 'protocolClass']
        ]
        ['0xD0' COTPPacketConnectionResponse
            [simple uint 16           'destinationReference']
            [simple uint 16           'sourceReference']
            [enum   COTPProtocolClass 'protocolClass']
        ]
        ['0x80' COTPPacketDisconnectRequest
            [simple uint 16           'destinationReference']
            [simple uint 16           'sourceReference']
            [enum   COTPProtocolClass 'protocolClass']
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
            [enum COTPTpduSize 'tpduSize']
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
    [implicit      uint 16 'parameterLength' 'parameter != null ? parameter.lengthInBytes : 0']
    [implicit      uint 16 'payloadLength'   'payload != null ? payload.lengthInBytes : 0']
    [typeSwitch 'messageType'
        ['0x01' S7MessageRequest
        ]
        ['0x02' S7MessageResponse
            [simple uint 8 'errorClass']
            [simple uint 8 'errorCode']
        ]
        ['0x03' S7MessageResponseData
            [simple uint 8 'errorClass']
            [simple uint 8 'errorCode']
        ]
        ['0x07' S7MessageUserData
        ]
    ]
    [optional S7Parameter 'parameter' 'parameterLength > 0' ['messageType']]
    [optional S7Payload   'payload'   'payloadLength > 0'   ['messageType', 'parameter']]
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
            [implicit uint 8                  'numItems' 'COUNT(items)']
            [array    S7ParameterUserDataItem 'items' count 'numItems']
        ]
    ]
]

[discriminatedType 'S7VarRequestParameterItem'
    [discriminator uint 8 'itemType']
    [typeSwitch 'itemType'
        ['0x12' S7VarRequestParameterItemAddress
            [implicit uint 8    'itemLength' 'address.lengthInBytes']
            [simple   S7Address 'address']
        ]
    ]
]

[discriminatedType 'S7Address'
    [discriminator uint 8 'addressType']
    [typeSwitch 'addressType'
        ['0x10' S7AddressAny
            [enum     TransportSize 'transportSize']
            [simple   uint 16       'numberOfElements']
            [simple   uint 16       'dbNumber']
            [enum     MemoryArea    'area']
            [reserved uint 5        '0x00']
            [simple   uint 16       'byteAddress']
            [simple   uint 3        'bitAddress']
        ]
    ]
]

[discriminatedType 'S7ParameterUserDataItem'
    [discriminator uint 8 'itemType']
    [typeSwitch 'itemType'
        ['0x12' S7ParameterUserDataItemCPUFunctions
            [implicit uint 8  'itemLength' 'lengthInBytes - 2']
            [simple   uint 8  'method']
            [simple   uint 4  'cpuFunctionType']
            [simple   uint 4  'cpuFunctionGroup']
            [simple   uint 8  'cpuSubfunction']
            [simple   uint 8  'sequenceNumber']
            [optional uint 8  'dataUnitReferenceNumber' 'cpuFunctionType == 8']
            [optional uint 8  'lastDataUnit' 'cpuFunctionType == 8']
            [optional uint 16 'errorCode' 'cpuFunctionType == 8']
        ]
    ]
]

[type 'SzlId'
    [enum   SzlModuleTypeClass 'typeClass']
    [simple uint 4             'sublistExtract']
    [enum   SzlSublist         'sublistList']
]

[type 'SzlDataTreeItem'
    [simple uint 16 'itemIndex']
    [array  int 8   'mlfb' count '20']
    [simple uint 16 'moduleTypeId']
    [simple uint 16 'ausbg']
    [simple uint 16 'ausbe']
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatedType 'S7Payload' [uint 8 'messageType', S7Parameter 'parameter']
    [typeSwitch 'parameter.parameterType', 'messageType'
        ['0x04','0x03' S7PayloadReadVarResponse
            [array S7VarPayloadDataItem 'items' count 'CAST(parameter, S7ParameterReadVarResponse).numItems' ['lastItem']]
        ]
        ['0x05','0x01' S7PayloadWriteVarRequest
            [array S7VarPayloadDataItem 'items' count 'COUNT(CAST(parameter, S7ParameterWriteVarRequest).items)' ['lastItem']]
        ]
        ['0x05','0x03' S7PayloadWriteVarResponse
            [array S7VarPayloadStatusItem 'items' count 'CAST(parameter, S7ParameterWriteVarResponse).numItems']
        ]
        ['0x00','0x07' S7PayloadUserData
            [array S7PayloadUserDataItem 'items' count 'COUNT(CAST(parameter, S7ParameterUserData).items)' ['CAST(CAST(parameter, S7ParameterUserData).items[0], S7ParameterUserDataItemCPUFunctions).cpuFunctionType']]
        ]
    ]
]

// This is actually not quite correct as depending pon the transportSize the length is either defined in bits or bytes.
[type 'S7VarPayloadDataItem' [bit 'lastItem']
    [enum     DataTransportErrorCode 'returnCode']
    [enum     DataTransportSize      'transportSize']
    [implicit uint 16                'dataLength' 'COUNT(data) * ((transportSize == DataTransportSize.BIT) ? 1 : (transportSize.sizeInBits ? 8 : 1))']
    [array    int  8                 'data'       count 'transportSize.sizeInBits ? CEIL(dataLength / 8.0) : dataLength']
    [padding  uint 8                 'pad'        '0x00' '!lastItem && ((COUNT(data) % 2) == 1)']
]

[type 'S7VarPayloadStatusItem'
    [enum DataTransportErrorCode 'returnCode']
]

[discriminatedType 'S7PayloadUserDataItem' [uint 4 'cpuFunctionType']
    [enum     DataTransportErrorCode 'returnCode']
    [enum     DataTransportSize      'transportSize']
    [implicit uint 16                'dataLength' 'lengthInBytes - 4']
    [simple   SzlId                  'szlId']
    [simple   uint 16                'szlIndex']
    [typeSwitch 'cpuFunctionType'
        ['0x04' S7PayloadUserDataItemCpuFunctionReadSzlRequest
        ]
        ['0x08' S7PayloadUserDataItemCpuFunctionReadSzlResponse
            [const    uint 16 'szlItemLength' '28']
            [implicit uint 16 'szlItemCount'  'COUNT(items)']
            [array SzlDataTreeItem 'items' count 'szlItemCount']
        ]
    ]
]

[dataIo 'DataItem' [uint 8 'dataProtocolId']
    [typeSwitch 'dataProtocolId'
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['01' Boolean
            [reserved uint 7 '0x00']
            [simple   bit    'value']
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['11' List
            [array bit 'value' count '8']
        ]
        // 2 byte (16 bit)
        ['12' List
            [array bit 'value' count '16']
        ]
        // 4 byte (32 bit)
        ['13' List
            [array bit 'value' count '32']
        ]
        // 8 byte (64 bit)
        ['14' List
            [array bit 'value' count '64']
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['21' Integer
            [simple int 8 'value']
        ]
        ['22' Integer
            [simple uint 8 'value']
        ]
        // 16 bit:
        ['23' Integer
            [simple int 16 'value']
        ]
        ['24' Integer
            [simple uint 16 'value']
        ]
        // 32 bit:
        ['25' Integer
            [simple int 32 'value']
        ]
        ['26' Long
            [simple uint 32 'value']
        ]
        // 64 bit:
        ['27' Long
            [simple int 64 'value']
        ]
        ['28' BigInteger
            [simple uint 64 'value']
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['31' Float
            [simple float 8.23  'value']
        ]
        ['32' Double
            [simple float 11.52 'value']
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['41' String
        ]
        ['42' String
        ]
        ['43' String
            [manual string 'UTF-8' 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7String", io, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7String", io, _value, _type.encoding)' '_value.length + 2']
        ]
        ['44' String
            [manual string 'UTF-16' 'value''STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7String", io, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7String", io, _value, _type.encoding)' '(_value.length * 2) + 2']
        ]

        // -----------------------------------------
        // TIA Date-Formats
        // -----------------------------------------
        ['51' Time
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseTiaTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeTiaTime", io, _value)' '4']
        ]
        // TODO: Check if this is really 8 bytes
        ['52' Time
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseTiaLTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeTiaLTime", io, _value)' '8']
        ]
        ['53' Date
            [manual date 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseTiaDate", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeTiaDate", io, _value)' '2']
        ]
        ['54' Time
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseTiaTimeOfDay", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeTiaTimeOfDay", io, _value)' '4']
        ]
        ['55' DateTime
            [manual dateTime 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseTiaDateTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeTiaDateTime", io, _value)' '8']
        ]
    ]
]

[enum int 8 'COTPTpduSize' [uint 16 'sizeInBytes']
    ['0x07' SIZE_128 ['128']]
    ['0x08' SIZE_256 ['256']]
    ['0x09' SIZE_512 ['512']]
    ['0x0a' SIZE_1024 ['1024']]
    ['0x0b' SIZE_2048 ['2048']]
    ['0x0c' SIZE_4096 ['4096']]
    ['0x0d' SIZE_8192 ['8192']]
]

[enum int 8 'COTPProtocolClass'
    ['0x00' CLASS_0]
    ['0x10' CLASS_1]
    ['0x20' CLASS_2]
    ['0x30' CLASS_3]
    ['0x40' CLASS_4]
]

[enum int 8 'DataTransportSize' [bit 'sizeInBits']
    ['0x00' NULL            ['false']]
    ['0x03' BIT             ['true']]
    ['0x04' BYTE_WORD_DWORD ['true']]
    ['0x05' INTEGER         ['true']]
    ['0x06' DINTEGER        ['false']]
    ['0x07' REAL            ['false']]
    ['0x09' OCTET_STRING    ['false']]
]

[enum int 8 'DeviceGroup'
    ['0x01' PG_OR_PC]
    ['0x02' OS      ]
    ['0x03' OTHERS  ]
]

[enum int 8 'TransportSize'  [uint 8 'sizeCode', uint 8 'sizeInBytes', TransportSize 'baseType', DataTransportSize 'dataTransportSize', uint 8 'dataProtocolId', bit 'supported_S7_300', bit 'supported_S7_400', bit 'supported_S7_1200', bit 'supported_S7_1500', bit 'supported_LOGO']
    // Bit Strings
    ['0x01' BOOL             ['X'              , '1'                 , 'null'                  , 'DataTransportSize.BIT'              , '01'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x02' BYTE             ['B'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '11'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x04' WORD             ['W'              , '2'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '12'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x06' DWORD            ['D'              , '4'                 , 'TransportSize.WORD'    , 'DataTransportSize.BYTE_WORD_DWORD'  , '13'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x00' LWORD            ['X'              , '8'                 , 'null'                  , 'null'                               , '14'                   , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]

    // Integer values
    // INT and UINT moved out of order as the enum constant INT needs to be generated before it's used in java
    ['0x05' INT              ['W'              , '2'                 , 'null'                  , 'DataTransportSize.INTEGER'          , '23'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x05' UINT             ['W'              , '2'                 , 'TransportSize.INT'     , 'DataTransportSize.INTEGER'          , '24'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    // ...
    ['0x02' SINT             ['B'              , '1'                 , 'TransportSize.INT'     , 'DataTransportSize.BYTE_WORD_DWORD'  , '21'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x02' USINT            ['B'              , '1'                 , 'TransportSize.INT'     , 'DataTransportSize.BYTE_WORD_DWORD'  , '22'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x07' DINT             ['D'              , '4'                 , 'TransportSize.INT'     , 'DataTransportSize.INTEGER'          , '25'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x07' UDINT            ['D'              , '4'                 , 'TransportSize.INT'     , 'DataTransportSize.INTEGER'          , '26'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x00' LINT             ['X'              , '8'                 , 'TransportSize.INT'     , 'null'                               , '27'                   , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]
    ['0x00' ULINT            ['X'              , '16'                , 'TransportSize.INT'     , 'null'                               , '28'                   , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]

    // Floating point values
    ['0x08' REAL             ['D'              , '4'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '31'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x30' LREAL            ['X'              , '8'                 , 'TransportSize.REAL'    , 'null'                               , '32'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'false'             ]]

    // Characters and Strings
    ['0x03' CHAR             ['B'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '41'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x13' WCHAR            ['X'              , '2'                 , 'null'                  , 'null'                               , '42'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x03' STRING           ['X'              , '1'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '43'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x00' WSTRING          ['X'              , '1'                 , 'null'                  , 'null'                               , '44'                   , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]

    // Dates and time values
    ['0x0B' TIME             ['X'              , '4'                 , 'null'                  , 'null'                               , '51'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x00' LTIME            ['X'              , '8'                 , 'TransportSize.TIME'    , 'null'                               , '52'                   , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]
    ['0x02' DATE             ['X'              , '2'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '53'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x02' TIME_OF_DAY      ['X'              , '4'                 , 'null'                  , 'DataTransportSize.BYTE_WORD_DWORD'  , '54'                   , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x02' DATE_AND_TIME    ['X'              , '8'                 , 'null'                  , 'null'                               , '55'                   , 'true'                , 'true'                , 'false'                , 'true'                 , 'false'             ]]
]

[enum int 8 'MemoryArea'             [string 24 'utf8' 'shortName']
    ['0x1C' COUNTERS                 ['C']]
    ['0x1D' TIMERS                   ['T']]
    ['0x80' DIRECT_PERIPHERAL_ACCESS ['D']]
    ['0x81' INPUTS                   ['I']]
    ['0x82' OUTPUTS                  ['Q']]
    ['0x83' FLAGS_MARKERS            ['M']]
    ['0x84' DATA_BLOCKS              ['DB']]
    ['0x85' INSTANCE_DATA_BLOCKS     ['DBI']]
    ['0x86' LOCAL_DATA               ['LD']]
]

[enum int 8 'DataTransportSize' [bit 'sizeInBits']
    ['0x00' NULL                ['false']]
    ['0x03' BIT                 ['true']]
    ['0x04' BYTE_WORD_DWORD     ['true']]
    ['0x05' INTEGER             ['true']]
    ['0x06' DINTEGER            ['false']]
    ['0x07' REAL                ['false']]
    ['0x09' OCTET_STRING        ['false']]
]

[enum int 8 'DataTransportErrorCode'
    ['0x00' RESERVED               ]
    ['0xFF' OK                     ]
    ['0x03' ACCESS_DENIED          ]
    ['0x05' INVALID_ADDRESS        ]
    ['0x06' DATA_TYPE_NOT_SUPPORTED]
    ['0x0A' NOT_FOUND              ]
]

[enum int 4 'SzlModuleTypeClass'
    ['0x0' CPU]
    ['0x4' IM]
    ['0x8' FM]
    ['0xC' CP]
]

[enum int 8 'SzlSublist'
    ['0x11' MODULE_IDENTIFICATION]
    ['0x12' CPU_FEATURES]
    ['0x13' USER_MEMORY_AREA]
    ['0x14' SYSTEM_AREAS]
    ['0x15' BLOCK_TYPES]
    ['0x19' STATUS_MODULE_LEDS]
    ['0x1C' COMPONENT_IDENTIFICATION]
    ['0x22' INTERRUPT_STATUS]
    ['0x25' ASSIGNMENT_BETWEEN_PROCESS_IMAGE_PARTITIONS_AND_OBS]
    ['0x32' COMMUNICATION_STATUS_DATA]
    ['0x74' STATUS_SINGLE_MODULE_LED]
    ['0x90' DP_MASTER_SYSTEM_INFORMATION]
    ['0x91' MODULE_STATUS_INFORMATION]
    ['0x92' RACK_OR_STATION_STATUS_INFORMATION]
    ['0x94' RACK_OR_STATION_STATUS_INFORMATION_2]
    ['0x95' ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION]
    ['0x96' MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP]
    ['0xA0' DIAGNOSTIC_BUFFER]
    ['0xB1' MODULE_DIAGNOSTIC_DATA]
]



