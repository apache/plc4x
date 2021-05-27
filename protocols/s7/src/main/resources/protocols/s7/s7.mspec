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
        ['0xE0' COTPParameterDisconnectAdditionalInformation [uint 8 'rest']
            [array byte 'data' count 'rest']
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
        ['0x01','0x07' S7ParameterModeTransition
            [reserved uint 16  '0x0010']
            [implicit uint 8  'itemLength' 'lengthInBytes - 2']
            [simple   uint 8  'method']
            [simple   uint 4  'cpuFunctionType']
            [simple   uint 4  'cpuFunctionGroup']
            [simple   uint 8  'currentMode']
            [simple   uint 8  'sequenceNumber']             
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
            [enum     TransportSize 'transportSize' 'code']
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
    [array  byte    'mlfb' count '20']
    [simple uint 16 'moduleTypeId']
    [simple uint 16 'ausbg']
    [simple uint 16 'ausbe']
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatedType 'S7Payload' [uint 8 'messageType', S7Parameter 'parameter']
    [typeSwitch 'parameter.parameterType', 'messageType'
        ['0x04','0x03' S7PayloadReadVarResponse [S7Parameter 'parameter']
            [array S7VarPayloadDataItem 'items' count 'CAST(parameter, S7ParameterReadVarResponse).numItems' ['lastItem']]
        ]
        ['0x05','0x01' S7PayloadWriteVarRequest [S7Parameter 'parameter']
            [array S7VarPayloadDataItem 'items' count 'COUNT(CAST(parameter, S7ParameterWriteVarRequest).items)' ['lastItem']]
        ]
        ['0x05','0x03' S7PayloadWriteVarResponse [S7Parameter 'parameter']
            [array S7VarPayloadStatusItem 'items' count 'CAST(parameter, S7ParameterWriteVarResponse).numItems']
        ]
        ['0x00','0x07' S7PayloadUserData [S7Parameter 'parameter']
            [array S7PayloadUserDataItem 'items' count 'COUNT(CAST(parameter, S7ParameterUserData).items)' ['CAST(CAST(parameter, S7ParameterUserData).items[0], S7ParameterUserDataItemCPUFunctions).cpuFunctionType', 'CAST(CAST(parameter, S7ParameterUserData).items[0], S7ParameterUserDataItemCPUFunctions).cpuSubfunction']]
        ]
    ]
]

// This is actually not quite correct as depending pon the transportSize the length is either defined in bits or bytes.
[type 'S7VarPayloadDataItem' [bit 'lastItem']
    [enum     DataTransportErrorCode 'returnCode']
    [enum     DataTransportSize      'transportSize']
    [implicit uint 16                'dataLength' 'COUNT(data) * ((transportSize == DataTransportSize.BIT) ? 1 : (transportSize.sizeInBits ? 8 : 1))']
    [array    byte                   'data'       count 'transportSize.sizeInBits ? CEIL(dataLength / 8.0) : dataLength']
    [padding  uint 8                 'pad'        '0x00' 'lastItem ? 0 : COUNT(data) % 2']
]

[type 'S7VarPayloadStatusItem'
    [enum DataTransportErrorCode 'returnCode']
]


////////////////////////////////////////////////////////////////
// Event 7 Alarms Types
////////////////////////////////////////////////////////////////

//Under test
[discriminatedType  'S7DataAlarmMessage' [uint 4 'cpuFunctionType']
    [const    uint 8 'functionId' '0x00']
    [const    uint 8 'numberMessageObj' '0x01']
    [typeSwitch 'cpuFunctionType'
        ['0x04' S7MessageObjectRequest
            [const  uint 8  'variableSpec'  '0x12']
            [const  uint 8  'length'        '0x08']
            [enum   SyntaxIdType    'syntaxId']
            [reserved uint 8        '0x00']
            [enum   QueryType    'queryType']
            [reserved uint 8        '0x34']
            [enum   AlarmType    'alarmType']
        ]
        ['0x08' S7MessageObjectResponse
            [enum     DataTransportErrorCode 'returnCode']
            [enum     DataTransportSize      'transportSize']
            [reserved uint 8        '0x00']
        ]        
    ]
]

[type 'AssociatedValueType'
    [enum   DataTransportErrorCode 'returnCode']
    [enum   DataTransportSize      'transportSize']
    [manual uint 16 'valueLength' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.RightShift3", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.LeftShift3", writeBuffer, _value.valueLength)' '2']
    [array uint 8 'data' count 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.EventItemLength", readBuffer, valueLength)']
]

//TODO: Convert BCD to uint
[type 'DateAndTime'
    [manual uint 8 'year' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.year)' '1']
    [manual uint 8  'month' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.month)' '1']
    [manual uint 8  'day' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.day)' '1']
    [manual uint 8  'hour' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.hour)' '1']
    [manual uint 8  'minutes' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.minutes)' '1']
    [manual uint 8  'seconds' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.BcdToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.ByteToBcd", writeBuffer, _value.seconds)' '1']
    [manual uint 12 'msec' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.S7msecToInt", readBuffer)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.S7EventHelper.IntToS7msec", writeBuffer, _value.msec)' '2']
    [simple uint 4  'dow'] 
]

[type 'State'
    [simple bit 'SIG_8']
    [simple bit 'SIG_7']  
    [simple bit 'SIG_6']
    [simple bit 'SIG_5'] 
    [simple bit 'SIG_4']
    [simple bit 'SIG_3']  
    [simple bit 'SIG_2']
    [simple bit 'SIG_1']   
]

[type 'AlarmMessageObjectPushType'
    [const uint 8 'variableSpec' '0x12']
    [simple uint 8 'lengthspec']  
    [simple SyntaxIdType 'syntaxId'] 
    [simple uint 8 'numberOfValues']
    [simple uint 32 'eventId']
    [simple State 'eventState']
    [simple State 'localState']
    [simple State 'ackStateGoing']
    [simple State 'ackStateComing']
    [array AssociatedValueType 'AssociatedValues' count 'numberOfValues' ]
]

[type 'AlarmMessagePushType'
    [simple DateAndTime 'TimeStamp']
    [simple uint 8 'functionId']
    [simple uint 8 'numberOfObjects']
    [array AlarmMessageObjectPushType 'messageObjects' count 'numberOfObjects' ]  
]

[type 'AlarmMessageObjectQueryType'
    [simple uint 8 'lengthDataset'] 
    [reserved uint 16 '0x0000']
    [const uint 8 'variableSpec' '0x12']
    [simple State 'eventState']
    [simple State 'ackStateGoing']
    [simple State 'ackStateComing']
    [simple DateAndTime 'timeComing']
    [simple AssociatedValueType 'valueComing'] 
    [simple DateAndTime 'timeGoing']
    [simple AssociatedValueType 'valueGoing'] 
]

[type 'AlarmMessageQueryType'
    [simple uint 8 'functionId']
    [simple uint 8 'numberOfObjects']
    [enum   DataTransportErrorCode 'returnCode']
    [enum   DataTransportSize      'transportSize']
    [const uint 16 'DataLength' '0xFFFF']
    [array AlarmMessageObjectQueryType 'messageObjects' count 'numberOfObjects' ]  
]

[type 'AlarmMessageObjectAckType'
    [const uint 8 'variableSpec' '0x12']
    [const uint 8 'length' '0x08']  
    [simple SyntaxIdType 'syntaxId'] 
    [simple uint 8 'numberOfValues']
    [simple uint 32 'eventId']
    [simple State 'ackStateGoing']
    [simple State 'ackStateComing']
]

[type 'AlarmMessageAckType'
    [simple uint 8 'functionId']
    [simple uint 8 'numberOfObjects']
    [array AlarmMessageObjectAckType 'messageObjects' count 'numberOfObjects' ]  
]

[type 'AlarmMessageAckResponseType'
    [simple uint 8 'functionId']
    [simple uint 8 'numberOfObjects']
    [array uint 8 'messageObjects' count 'numberOfObjects' ]  
]

////////////////////////////////////////////////////////////////
// DataItem by Function Type:
// 0x00 PUSH
// 0x04 TYPE_REQ
// 0x08 TYPE_RES
//
// DataItem by Sub Function Type:
// 0x01 CPU_READSZL
// 0x02 CPU_MSGS 
// 0x03 CPU_DIAGMSG
// 0x05 ALARM8_IND
// 0x06 NOTIFY_IND
// 0x07 ALARM8LOCK
// 0x08 ALARM8UNLOCK
// 0x0b ALARMACK
// 0x0c ALARMACK_IND
// 0x0d ALARM8LOCK_IND
// 0x0e ALARM8UNLOCK_IND
// 0x11 ALARMSQ_IND
// 0x12 ALARMS_IND
// 0x13 ALARMQUERY
// 0x16 NOTIFY8_IND
////////////////////////////////////////////////////////////////

[discriminatedType 'S7PayloadUserDataItem' [uint 4 'cpuFunctionType', uint 8 'cpuSubfunction']
    [enum     DataTransportErrorCode 'returnCode']
    [enum     DataTransportSize      'transportSize']
    [implicit uint 16                'dataLength' 'lengthInBytes - 4']
    [typeSwitch 'cpuFunctionType', 'cpuSubfunction', 'dataLength'

        //USER and SYSTEM Messages
        ['0x00', '0x03' S7PayloadDiagnosticMessage
            [simple uint 16 'EventId']
            [simple uint 8  'PriorityClass']
            [simple uint 8  'ObNumber']
            [simple uint 16 'DatId']
            [simple uint 16 'Info1']
            [simple uint 32 'Info2']
            [simple DateAndTime 'TimeStamp']
        ]

        //PUSH message reception S7300 & S7400 (ALARM_SQ, ALARM_S, ALARM_SC, ...)
        ['0x00', '0x05' S7PayloadAlarm8
            [simple AlarmMessagePushType 'alarmMessage']
        ]
        ['0x00', '0x06' S7PayloadNotify
            [simple AlarmMessagePushType 'alarmMessage']
        ]
        ['0x00', '0x11' S7PayloadAlarmSQ
            [simple AlarmMessagePushType 'alarmMessage']
        ]
        ['0x00', '0x12' S7PayloadAlarmS
            [simple AlarmMessagePushType 'alarmMessage']
        ]
        ['0x00', '0x13' S7PayloadAlarmSC
            [simple AlarmMessagePushType 'alarmMessage']
        ]
        ['0x00', '0x16' S7PayloadNotify8
            [simple AlarmMessagePushType 'alarmMessage']
        ]

        //Request for specific functions of the SZL system
        ['0x04', '0x01' S7PayloadUserDataItemCpuFunctionReadSzlRequest
            [simple   SzlId                  'szlId']
            [simple   uint 16                'szlIndex']
        ]
        ['0x08', '0x01' S7PayloadUserDataItemCpuFunctionReadSzlResponse
            [simple   SzlId                  'szlId']
            [simple   uint 16                'szlIndex']
            [const    uint 16 'szlItemLength' '28']
            [implicit uint 16 'szlItemCount'  'COUNT(items)']
            [array SzlDataTreeItem 'items' count 'szlItemCount']
        ]

        //Subscription to PUSH messages
        ['0x04', '0x02' S7PayloadUserDataItemCpuFunctionMsgSubscription
            [simple   uint 8    'Subscription']
            [reserved uint 8    '0x00']
            [simple string '64' 'UTF-8' 'magicKey']
            [optional AlarmStateType 'Alarmtype' 'Subscription >= 128']
            [optional uint 8 'Reserve' 'Subscription >= 128']
        ]
	['0x08', '0x02', '0x00' S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse
        ]
	['0x08', '0x02', '0x02' S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse
            [simple uint 8 'result']
            [simple uint 8 'reserved01']
        ]
	['0x08', '0x02', '0x05' S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse
            [simple uint 8 'result']
            [simple uint 8 'reserved01']
            [simple AlarmType 'alarmType']
            [simple uint 8 'reserved02']
            [simple uint 8 'reserved03']
        ]

        //ALARM_ACK Acknowledgment of alarms
        ['0x04', '0x0b' S7PayloadUserDataItemCpuFunctionAlarmAck
            [simple uint 8 'functionId']
            [implicit uint 8 'numberOfObjects' 'COUNT(messageObjects)']
            [array AlarmMessageObjectAckType 'messageObjects' count 'numberOfObjects' ]  
        ]
        ['0x08', '0x0b' S7PayloadUserDataItemCpuFunctionAlarmAckResponse
            [simple uint 8 'functionId']
            [implicit  uint 8 'numberOfObjects' 'COUNT(messageObjects)']
            [array uint 8 'messageObjects' count 'numberOfObjects' ]  
        ]

        //ALARM_QUERY Request for alarms stored in the controller
        ['0x04', '0x13' S7PayloadUserDataItemCpuFunctionAlarmQuery
            [const    uint 8 'functionId' '0x00']
            [const    uint 8 'numberMessageObj' '0x01']
            [const  uint 8  'variableSpec'  '0x12']
            [const  uint 8  'length'        '0x08']
            [enum   SyntaxIdType    'syntaxId']
            [reserved uint 8        '0x00']
            [enum   QueryType    'queryType']
            [reserved uint 8        '0x34']
            [enum   AlarmType    'alarmType']
        ]
        ['0x08', '0x13' S7PayloadUserDataItemCpuFunctionAlarmQueryResponse
            [const    uint 8 'functionId' '0x00']
            [const    uint 8 'numberMessageObj' '0x01']
            [enum     DataTransportErrorCode 'pudicfReturnCode']
            [enum     DataTransportSize      'pudicftransportSize']
            [reserved uint 8        '0x00']
        ]
    ]
]

[dataIo 'DataItem' [string '-1' 'dataProtocolId', int 32 'stringLength']
    [typeSwitch 'dataProtocolId'
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['IEC61131_BOOL' BOOL
            [reserved uint 7 '0x00']
            [simple   bit    'value']
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['IEC61131_BYTE' List
            [array bit 'value' count '8']
        ]
        // 2 byte (16 bit)
        ['IEC61131_WORD' List
            [array bit 'value' count '16']
        ]
        // 4 byte (32 bit)
        ['IEC61131_DWORD' List
            [array bit 'value' count '32']
        ]
        // 8 byte (64 bit)
        ['IEC61131_LWORD' List
            [array bit 'value' count '64']
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['IEC61131_SINT' SINT
            [simple int 8 'value']
        ]
        ['IEC61131_USINT' USINT
            [simple uint 8 'value']
        ]
        // 16 bit:
        ['IEC61131_INT' INT
            [simple int 16 'value']
        ]
        ['IEC61131_UINT' UINT
            [simple uint 16 'value']
        ]
        // 32 bit:
        ['IEC61131_DINT' DINT
            [simple int 32 'value']
        ]
        ['IEC61131_UDINT' UDINT
            [simple uint 32 'value']
        ]
        // 64 bit:
        ['IEC61131_LINT' LINT
            [simple int 64 'value']
        ]
        ['IEC61131_ULINT' ULINT
            [simple uint 64 'value']
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['IEC61131_REAL' REAL
            [simple float 8.23  'value']
        ]
        ['IEC61131_LREAL' LREAL
            [simple float 11.52 'value']
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['IEC61131_CHAR' CHAR
            [manual string '-1' 'UTF-8' 'value'  'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7Char", readBuffer, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7Char", writeBuffer, _value, _type.encoding)' '1']
        ]
        ['IEC61131_WCHAR' CHAR
            [manual string '-1' 'UTF-16' 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7Char", readBuffer, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7Char", writeBuffer, _value, _type.encoding)' '2']
        ]
        ['IEC61131_STRING' STRING
            [manual string '-1' 'UTF-8' 'value'  'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7String", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7String", writeBuffer, _value, stringLength, _type.encoding)' '_value.string.length + 2']
        ]
        ['IEC61131_WSTRING' STRING
            [manual string '-1' 'UTF-16' 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.parseS7String", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.s7.utils.StaticHelper.serializeS7String", writeBuffer, _value, stringLength, _type.encoding)' '(_value.string.length * 2) + 2']
        ]

        // -----------------------------------------
        // TIA Date-Formats
        // -----------------------------------------
        // Interpreted as "milliseconds"
        ['IEC61131_TIME' TIME
            [simple uint 32 'value']
        ]
        //['S7_S5TIME' TIME
        //    [reserved uint 2  '0x00']
        //    [uint     uint 2  'base']
        //    [simple   uint 12 'value']
        //]
        // Interpreted as "number of nanoseconds"
        ['IEC61131_LTIME' LTIME
            [simple uint 64 'value']
        ]
        // Interpreted as "number of days since 1990-01-01"
        ['IEC61131_DATE' DATE
            [simple uint 16 'value']
        ]
        ['IEC61131_TIME_OF_DAY' TIME_OF_DAY
            [simple uint 32 'value']
        ]
        ['IEC61131_DATE_AND_TIME' DATE_AND_TIME
            [simple uint 16 'year']
            [simple uint 8  'month']
            [simple uint 8  'day']
            [simple uint 8  'dayOfWeek']
            [simple uint 8  'hour']
            [simple uint 8  'minutes']
            [simple uint 8  'seconds']
            [simple uint 32 'nanos']
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

[enum int 8 'TransportSize'  [uint 8 'code', uint 8 'shortName', uint 8 'sizeInBytes', TransportSize 'baseType', DataTransportSize 'dataTransportSize', string '-1' 'dataProtocolId', bit 'supported_S7_300', bit 'supported_S7_400', bit 'supported_S7_1200', bit 'supported_S7_1500', bit 'supported_LOGO']
    // Bit Strings
    ['0x01' BOOL             ['0x01'       , 'X'               , '1'                 , 'null'                  , 'BIT'              , 'IEC61131_BOOL'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x02' BYTE             ['0x02'       , 'B'               , '1'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_BYTE'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x03' WORD             ['0x04'       , 'W'               , '2'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_WORD'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x04' DWORD            ['0x06'       , 'D'               , '4'                 , 'WORD'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_DWORD'        , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x05' LWORD            ['0x00'       , 'X'               , '8'                 , 'null'                  , 'null'             , 'IEC61131_LWORD'        , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]

    // Integer values
    // INT and UINT moved out of order as the enum constant INT needs to be generated before it's used in java
    ['0x06' INT              ['0x05'       , 'W'               , '2'                 , 'null'                  , 'INTEGER'          , 'IEC61131_INT'          , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x07' UINT             ['0x05'       , 'W'               , '2'                 , 'INT'                   , 'INTEGER'          , 'IEC61131_UINT'         , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    // ...
    ['0x08' SINT             ['0x02'       , 'B'               , '1'                 , 'INT'                   , 'BYTE_WORD_DWORD'  , 'IEC61131_SINT'         , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x09' USINT            ['0x02'       , 'B'               , '1'                 , 'INT'                   , 'BYTE_WORD_DWORD'  , 'IEC61131_USINT'        , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x0A' DINT             ['0x07'       , 'D'               , '4'                 , 'INT'                   , 'INTEGER'          , 'IEC61131_DINT'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x0B' UDINT            ['0x07'       , 'D'               , '4'                 , 'INT'                   , 'INTEGER'          , 'IEC61131_UDINT'        , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x0C' LINT             ['0x00'       , 'X'               , '8'                 , 'INT'                   , 'null'             , 'IEC61131_LINT'         , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]
    ['0x0D' ULINT            ['0x00'       , 'X'               , '16'                , 'INT'                   , 'null'             , 'IEC61131_ULINT'        , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]

    // Floating point values
    ['0x0E' REAL             ['0x08'       , 'D'               , '4'                 , 'null'                  , 'REAL'             , 'IEC61131_REAL'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x0F' LREAL            ['0x30'       , 'X'               , '8'                 , 'REAL'                  , 'null'             , 'IEC61131_LREAL'        , 'false'               , 'false'               , 'true'                 , 'true'                 , 'false'             ]]

    // Characters and Strings
    ['0x10' CHAR             ['0x03'       , 'B'               , '1'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_CHAR'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x11' WCHAR            ['0x13'       , 'X'               , '2'                 , 'null'                  , 'null'             , 'IEC61131_WCHAR'        , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]
    ['0x12' STRING           ['0x03'       , 'X'               , '1'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_STRING'       , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x13' WSTRING          ['0x00'       , 'X'               , '2'                 , 'null'                  , 'null'             , 'IEC61131_WSTRING'      , 'false'               , 'false'               , 'true'                 , 'true'                 , 'true'              ]]

    // Dates and time values (Please note that we seem to have to rewrite queries for these types to reading bytes or we'll get "Data type not supported" errors)
    ['0x14' TIME             ['0x0B'       , 'X'               , '4'                 , 'null'                  , 'null'             , 'IEC61131_TIME'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    //['0x15' S5TIME           ['0x0C'      , 'X'               , '4'                 , 'null'                  , 'null'                               , 'S7_S5TIME'             , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x16' LTIME            ['0x00'       , 'X'               , '8'                 , 'TIME'                  , 'null'             , 'IEC61131_LTIME'        , 'false'               , 'false'               , 'false'                , 'true'                 , 'false'             ]]
    ['0x17' DATE             ['0x09'       , 'X'               , '2'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_DATE'         , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x18' TIME_OF_DAY      ['0x06'       , 'X'               , '4'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_TIME_OF_DAY'  , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x19' TOD              ['0x06'       , 'X'               , '4'                 , 'null'                  , 'BYTE_WORD_DWORD'  , 'IEC61131_TIME_OF_DAY'  , 'true'                , 'true'                , 'true'                 , 'true'                 , 'true'              ]]
    ['0x1A' DATE_AND_TIME    ['0x0F'       , 'X'               , '12'                , 'null'                  , 'null'             , 'IEC61131_DATE_AND_TIME', 'true'                , 'true'                , 'false'                , 'true'                 , 'false'             ]]
    ['0x1B' DT               ['0x0F'       , 'X'               , '12'                , 'null'                  , 'null'             , 'IEC61131_DATE_AND_TIME', 'true'                , 'true'                , 'false'                , 'true'                 , 'false'             ]]
]

[enum uint 8 'MemoryArea'             [string '24' 'shortName']
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

[enum uint 8 'DataTransportSize' [bit 'sizeInBits']
    ['0x00' NULL                ['false']]
    ['0x03' BIT                 ['true']]
    ['0x04' BYTE_WORD_DWORD     ['true']]
    ['0x05' INTEGER             ['true']]
    ['0x06' DINTEGER            ['false']]
    ['0x07' REAL                ['false']]
    ['0x09' OCTET_STRING        ['false']]
]

[enum uint 8 'DataTransportErrorCode'
    ['0x00' RESERVED               ]
    ['0xFF' OK                     ]
    ['0x03' ACCESS_DENIED          ]
    ['0x05' INVALID_ADDRESS        ]
    ['0x06' DATA_TYPE_NOT_SUPPORTED]
    ['0x0A' NOT_FOUND              ]
]

[enum uint 4 'SzlModuleTypeClass'
    ['0x0' CPU]
    ['0x4' IM]
    ['0x8' FM]
    ['0xC' CP]
]

[enum uint 8 'SzlSublist'
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

[enum uint 8 'CpuSubscribeEvents'
    ['0x01' CPU]
    ['0x02' IM]
    ['0x04' FM]
    ['0x80' CP]
]

[enum uint 8 'EventType'
    ['0x01' MODE]
    ['0x02' SYS]
    ['0x04' USR]
    ['0x80' ALM]
]

[enum uint 8 'SyntaxIdType'
    ['0x01' S7ANY]
    ['0x13' PBC_ID]
    ['0x15' ALARM_LOCKFREESET]
    ['0x16' ALARM_INDSET]
    ['0x19' ALARM_ACKSET]
    ['0x1A' ALARM_QUERYREQSET]
    ['0x1C' NOTIFY_INDSET]
    ['0x82' NCK]
    ['0x83' NCK_METRIC]
    ['0x84' NCK_INCH]
    ['0xA2' DRIVEESANY]
    ['0xB2' SYM1200]
    ['0xB0' DBREAD]
]

[enum uint 8 'AlarmType'
    ['0x01' SCAN]
    ['0x02' ALARM_8]
    ['0x04' ALARM_S]
]

[enum uint 8 'AlarmStateType'
    ['0x00' SCAN_ABORT]
    ['0x01' SCAN_INITIATE]
    ['0x04' ALARM_ABORT]
    ['0x05' ALARM_INITIATE]
    ['0x08' ALARM_S_ABORT]
    ['0x09' ALARM_S_INITIATE]
]

[enum uint 8 'QueryType'
    ['0x01' BYALARMTYPE]
    ['0x02' ALARM_8]
    ['0x04' ALARM_S]
]

[enum uint 8 'ModeTransitionType'
    ['0x00' STOP]
    ['0x01' WARM_RESTART]
    ['0x02' RUN]
    ['0x03' HOT_RESTART]
    ['0x04' HOLD]
    ['0x06' COLD_RESTART]
    ['0x09' RUN_R]
    ['0x11' LINK_UP]
    ['0x12' UPDATE]
]

