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
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type TPKTPacket byteOrder='BIG_ENDIAN'
    [const    uint 8                 protocolId 0x03]
    [reserved uint 8                 '0x00']
    [implicit uint 16                len       'payload.lengthInBytes + 4']
    [simple   COTPPacket('len - 4') payload]
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType COTPPacket (uint 16 cotpLen)
    [implicit      uint 8 headerLength 'lengthInBytes - (((payload != null) ? payload.lengthInBytes : 0) + 1)']
    [discriminator uint 8 tpduCode]
    [typeSwitch tpduCode
        ['0xF0' COTPPacketData
            [simple bit    eot]
            [simple uint 7 tpduRef]
        ]
        ['0xE0' COTPPacketConnectionRequest
            [simple uint 16           destinationReference]
            [simple uint 16           sourceReference]
            [simple COTPProtocolClass protocolClass]
        ]
        ['0xD0' COTPPacketConnectionResponse
            [simple uint 16           destinationReference]
            [simple uint 16           sourceReference]
            [simple COTPProtocolClass protocolClass]
        ]
        ['0x80' COTPPacketDisconnectRequest
            [simple uint 16           destinationReference]
            [simple uint 16           sourceReference]
            [simple COTPProtocolClass protocolClass]
        ]
        ['0xC0' COTPPacketDisconnectResponse
            [simple uint 16 destinationReference]
            [simple uint 16 sourceReference]
        ]
        ['0x70' COTPPacketTpduError
            [simple uint 16 destinationReference]
            [simple uint 8  rejectCause]
        ]
    ]
    [array    COTPParameter ('(headerLength + 1) - curPos') parameters length '(headerLength + 1) - curPos']
    [optional S7Message                                     payload    'curPos < cotpLen']
]

[discriminatedType COTPParameter (uint 8 rest)
    [discriminator uint 8 parameterType]
    [implicit      uint 8 parameterLength 'lengthInBytes - 2']
    [typeSwitch parameterType
        ['0xC0' COTPParameterTpduSize
            [simple COTPTpduSize tpduSize]
        ]
        ['0xC1' COTPParameterCallingTsap
            [simple uint 16 tsapId]
        ]
        ['0xC2' COTPParameterCalledTsap
            [simple uint 16 tsapId]
        ]
        ['0xC3' COTPParameterChecksum
            [simple uint 8 crc]
        ]
        ['0xE0' COTPParameterDisconnectAdditionalInformation
            [array byte data count 'rest']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// S7
////////////////////////////////////////////////////////////////

[discriminatedType S7Message
    [const         uint 8  protocolId      0x32]
    [discriminator uint 8  messageType]
    [reserved      uint 16 '0x0000']
    [simple        uint 16 tpduReference]
    [implicit      uint 16 parameterLength 'parameter != null ? parameter.lengthInBytes : 0']
    [implicit      uint 16 payloadLength   'payload != null ? payload.lengthInBytes : 0']
    [typeSwitch messageType
        ['0x01' S7MessageRequest
        ]
        ['0x02' S7MessageResponse
            [simple uint 8 errorClass]
            [simple uint 8 errorCode]
        ]
        ['0x03' S7MessageResponseData
            [simple uint 8 errorClass]
            [simple uint 8 errorCode]
        ]
        ['0x07' S7MessageUserData
        ]
    ]
    [optional S7Parameter ('messageType')              parameter 'parameterLength > 0']
    [optional S7Payload   ('messageType', 'parameter') payload   'payloadLength > 0'  ]
]

////////////////////////////////////////////////////////////////
// Parameters

[discriminatedType S7Parameter (uint 8 messageType)
    [discriminator uint 8 parameterType]
    [typeSwitch parameterType,messageType
        ['0xF0' S7ParameterSetupCommunication
            [reserved uint 8  '0x00']
            [simple   uint 16 maxAmqCaller]
            [simple   uint 16 maxAmqCallee]
            [simple   uint 16 pduLength]
        ]
        ['0x04','0x01' S7ParameterReadVarRequest
            [implicit uint 8                    numItems 'COUNT(items)']
            [array    S7VarRequestParameterItem items    count 'numItems']
        ]
        ['0x04','0x03' S7ParameterReadVarResponse
            [simple uint 8 numItems]
        ]
        ['0x05','0x01' S7ParameterWriteVarRequest
            [implicit uint 8                    numItems 'COUNT(items)']
            [array    S7VarRequestParameterItem items    count 'numItems']
        ]
        ['0x05','0x03' S7ParameterWriteVarResponse
            [simple uint 8 numItems]
        ]
        ['0x00','0x07' S7ParameterUserData
            [implicit uint 8                  numItems 'COUNT(items)']
            [array    S7ParameterUserDataItem items count 'numItems']
        ]
        ['0x01','0x07' S7ParameterModeTransition
            [reserved uint 16  '0x0010']
            [implicit uint 8  itemLength 'lengthInBytes - 2']
            [simple   uint 8  method]
            [simple   uint 4  cpuFunctionType]
            [simple   uint 4  cpuFunctionGroup]
            [simple   uint 8  currentMode]
            [simple   uint 8  sequenceNumber]
        ]
    ]
]

[discriminatedType S7VarRequestParameterItem
    [discriminator uint 8 itemType]
    [typeSwitch itemType
        ['0x12' S7VarRequestParameterItemAddress
            [implicit uint 8    itemLength 'address.lengthInBytes']
            [simple   S7Address address]
        ]
    ]
]

[discriminatedType S7Address
    [discriminator uint 8 addressType]
    [typeSwitch addressType
        ['0x10' S7AddressAny
            [enum     TransportSize transportSize code]
            [simple   uint 16       numberOfElements]
            [simple   uint 16       dbNumber]
            [simple   MemoryArea    area]
            [reserved uint 5        '0x00']
            [simple   uint 16       byteAddress]
            [simple   uint 3        bitAddress]
        ]
    ]
]

[discriminatedType S7ParameterUserDataItem
    [discriminator uint 8 itemType]
    [typeSwitch itemType
        ['0x12' S7ParameterUserDataItemCPUFunctions
            [implicit uint 8  itemLength 'lengthInBytes - 2']
            [simple   uint 8  method]
            [simple   uint 4  cpuFunctionType]
            [simple   uint 4  cpuFunctionGroup]
            [simple   uint 8  cpuSubfunction]
            [simple   uint 8  sequenceNumber]
            [optional uint 8  dataUnitReferenceNumber '(cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))']
            [optional uint 8  lastDataUnit '(cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))']
            [optional uint 16 errorCode '(cpuFunctionType == 8) || ((cpuFunctionType == 0) && (cpuFunctionGroup == 2))']
        ]
    ]
]

/*
 * SZL is used as a reference to the list of system states. 
 * Siemens literature and forums use SZL or SSL interchangeably.
 * SZL = System Zustand Liste
 * SSL = System Status List
 */
[type SzlId
    [simple SzlModuleTypeClass typeClass]
    [simple uint 4             sublistExtract]
    [simple SzlSublist         sublistList]
]

[type SzlDataTreeItem
    [simple uint 16 itemIndex]
    [array  byte    mlfb count '20']
    [simple uint 16 moduleTypeId]
    [simple uint 16 ausbg]
    [simple uint 16 ausbe]
]

////////////////////////////////////////////////////////////////
// Payloads

[discriminatedType S7Payload (uint 8 messageType, S7Parameter parameter)
    [typeSwitch parameter.parameterType, messageType
        ['0x04','0x03' S7PayloadReadVarResponse
            [array S7VarPayloadDataItem items count 'CAST(parameter, "S7ParameterReadVarResponse").numItems']
        ]
        ['0x05','0x01' S7PayloadWriteVarRequest
            [array S7VarPayloadDataItem items count 'COUNT(CAST(parameter, "S7ParameterWriteVarRequest").items)']
        ]
        ['0x05','0x03' S7PayloadWriteVarResponse
            [array S7VarPayloadStatusItem items count 'CAST(parameter, "S7ParameterWriteVarResponse").numItems']
        ]
        ['0x00','0x07' S7PayloadUserData
            [array S7PayloadUserDataItem('CAST(CAST(parameter, "S7ParameterUserData").items[0], "S7ParameterUserDataItemCPUFunctions").cpuFunctionGroup', 'CAST(CAST(parameter, "S7ParameterUserData").items[0], "S7ParameterUserDataItemCPUFunctions").cpuFunctionType', 'CAST(CAST(parameter, "S7ParameterUserData").items[0], "S7ParameterUserDataItemCPUFunctions").cpuSubfunction') items count 'COUNT(CAST(parameter, "S7ParameterUserData").items)']
        ]
    ]
]

// This is actually not quite correct as depending pon the transportSize the length is either defined in bits or bytes.
//@param hasNext In the serialization process, if you have multiple write
//               requests the last element does not require padding.
[type S7VarPayloadDataItem
    [simple   DataTransportErrorCode returnCode]
    [simple   DataTransportSize      transportSize]
    [implicit uint 16                dataLength 'COUNT(data) * ((transportSize == DataTransportSize.BIT) ? 1 : (transportSize.sizeInBits ? 8 : 1))']
    [array    byte                   data       count 'transportSize.sizeInBits ? CEIL(dataLength / 8.0) : dataLength']
    [padding  uint 8                 pad        '0x00' '(!_lastItem) ? (COUNT(data) % 2) : 0']
]

[type S7VarPayloadStatusItem
    [simple DataTransportErrorCode returnCode]
]


////////////////////////////////////////////////////////////////
// Event 7 Alarms Types
////////////////////////////////////////////////////////////////

//Under test
[discriminatedType  S7DataAlarmMessage(uint 4 cpuFunctionType)
    [const    uint 8 functionId       0x00]
    [const    uint 8 numberMessageObj 0x01]
    [typeSwitch cpuFunctionType
        ['0x04' S7MessageObjectRequest
            [const    uint 8       variableSpec  0x12]
            [const    uint 8       length        0x08]
            [simple   SyntaxIdType syntaxId]
            [reserved uint 8       '0x00']
            [simple   QueryType    queryType]
            [reserved uint 8       '0x34']
            [simple   AlarmType    alarmType]
        ]
        ['0x08' S7MessageObjectResponse
            [simple   DataTransportErrorCode returnCode]
            [simple   DataTransportSize      transportSize]
            [reserved uint 8                 '0x00']
        ]
    ]
]

//TODO: The calculation must be modified to include the type
//      . if it is type 0x07(REAL) or 0x09 (OCTET_STRING), the length is indicated
//      . another type uses scrolling
//      . verify calculation with the other types
[type AssociatedValueType
    [simple DataTransportErrorCode returnCode]
    [simple DataTransportSize      transportSize]
    //[manual uint 16                valueLength   'STATIC_CALL("RightShift3", readBuffer)' 'STATIC_CALL("LeftShift3", writeBuffer, valueLength)' '16']
    [manual uint 16                valueLength  'STATIC_CALL("RightShift3", readBuffer, transportSize)' 'STATIC_CALL("LeftShift3", writeBuffer, valueLength)' '2']
    [array  uint 8                 data          count    'STATIC_CALL("EventItemLength", readBuffer, valueLength)']
]

[type AssociatedQueryValueType
    [simple DataTransportErrorCode returnCode]
    [simple DataTransportSize      transportSize]
    [simple uint 16                valueLength]
    [array  uint 8                 data          count    'valueLength']
]

//TODO: Convert BCD to uint
[type DateAndTime
    [manual uint 8  year    'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, year)'    '8']
    [manual uint 8  month   'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, month)'   '8']
    [manual uint 8  day     'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, day)'     '8']
    [manual uint 8  hour    'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, hour)'    '8']
    [manual uint 8  minutes 'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, minutes)' '8']
    [manual uint 8  seconds 'STATIC_CALL("BcdToInt", readBuffer)'    'STATIC_CALL("ByteToBcd", writeBuffer, seconds)' '8']
    [manual uint 12 msec    'STATIC_CALL("S7msecToInt", readBuffer)' 'STATIC_CALL("IntToS7msec", writeBuffer, msec)'  '12']
    [simple uint 4  dow                                                                                                         ]
]

[type State
    [simple bit SIG_8]
    [simple bit SIG_7]
    [simple bit SIG_6]
    [simple bit SIG_5]
    [simple bit SIG_4]
    [simple bit SIG_3]
    [simple bit SIG_2]
    [simple bit SIG_1]
]

[type AlarmMessageObjectPushType
    [const  uint 8              variableSpec     0x12]
    [simple uint 8              lengthSpec]
    [simple SyntaxIdType        syntaxId]
    [simple uint 8              numberOfValues]
    [simple uint 32             eventId]
    [simple State               eventState]
    [simple State               localState]
    [simple State               ackStateGoing]
    [simple State               ackStateComing]
    [array  AssociatedValueType AssociatedValues count 'numberOfValues' ]
]

[type AlarmMessageAckObjectPushType
    [const  uint 8       variableSpec 0x12]
    [simple uint 8       lengthSpec]
    [simple SyntaxIdType syntaxId]
    [simple uint 8       numberOfValues]
    [simple uint 32      eventId]
    [simple State        ackStateGoing]
    [simple State        ackStateComing]
]

[type AlarmMessagePushType
    [simple DateAndTime                TimeStamp]
    [simple uint 8                     functionId]
    [simple uint 8                     numberOfObjects]
    [array  AlarmMessageObjectPushType messageObjects count 'numberOfObjects' ]
]

[type AlarmMessageAckPushType
    [simple DateAndTime                   TimeStamp]
    [simple uint 8                        functionId]
    [simple uint 8                        numberOfObjects]
    [array  AlarmMessageAckObjectPushType messageObjects count 'numberOfObjects' ]
]

//TODO: Apply for S7-300
[type AlarmMessageQueryType(uint 16 dataLength)
    [simple uint 8                      functionId]
    [simple uint 8                      numberOfObjects]
    [simple DataTransportErrorCode      returnCode]
    [simple DataTransportSize           transportSize]
    [const  uint 16                     DataLength     0xFFFF]
    [array  AlarmMessageObjectQueryType messageObjects count   'STATIC_CALL("countAMOQT", readBuffer, dataLength)' ]
]

//TODO: Apply for S7-400
[type Alarm8MessageQueryType
    [simple uint 8                      functionId]
    [simple uint 8                      numberOfObjects]
    [simple DataTransportErrorCode      returnCode]
    [simple DataTransportSize           transportSize]
    [simple  uint 16                    byteCount]
    [array  AlarmMessageObjectQueryType messageObjects count   'byteCount / 12' ]
]

//TODO: Check for Alarm_8
[type AlarmMessageObjectQueryType
    [simple   uint 8              lengthDataset]
    [reserved uint 16             '0x0000']
    [const    uint 8              variableSpec   0x12]
    [simple   State               eventState]
    [simple   State               ackStateGoing]
    [simple   State               ackStateComing]
    [simple   DateAndTime         timeComing]
    [simple   AssociatedValueType valueComing]
    [simple   DateAndTime         timeGoing]
    [simple   AssociatedValueType valueGoing]
]

[type AlarmMessageQueryType
    [simple uint 8                      functionId]
    [simple uint 8                      numberOfObjects]
    [simple DataTransportErrorCode      returnCode]
    [simple DataTransportSize           transportSize]
    [const  uint 16                     DataLength     0xFFFF]
    [array  AlarmMessageObjectQueryType messageObjects count    'numberOfObjects' ]
]

[type AlarmMessageObjectAckType
    [const  uint 8       variableSpec 0x12]
    [const  uint 8       length 0x08]
    [simple SyntaxIdType syntaxId]
    [simple uint 8       numberOfValues]
    [simple uint 32      eventId]
    [simple State        ackStateGoing]
    [simple State        ackStateComing]
]

[type AlarmMessageAckType
    [simple uint 8                    functionId]
    [simple uint 8                    numberOfObjects]
    [array  AlarmMessageObjectAckType messageObjects count 'numberOfObjects' ]
]

[type AlarmMessageAckResponseType
    [simple uint 8 functionId]
    [simple uint 8 numberOfObjects]
    [array  uint 8 messageObjects  count 'numberOfObjects' ]
]

////////////////////////////////////////////////////////////////
// Cycle service Payloads
////////////////////////////////////////////////////////////////
//Under test
[discriminatedType  CycServiceItemType
    [const    uint 8 functionId       0x12]
    [simple   uint 8 byteLength]
    [simple   uint 8 syntaxId]
    [typeSwitch syntaxId
        ['0x10' CycServiceItemAnyType
            //[simple  TransportSize   transportSize]
            [enum     TransportSize transportSize code]
            [simple uint 16 length]
            [simple uint 16 dbNumber]            
            [simple MemoryArea memoryArea]
            [simple uint 24 address]
        ]
        ['0xb0' CycServiceItemDbReadType
            [simple   uint 8 numberOfAreas]            
            [array SubItem items count 'numberOfAreas']
        ]
    ]
]

[type SubItem
    [simple uint 8 bytesToRead]
    [simple uint 16 dbNumber]
    [simple uint 16 startAddress]
]

////////////////////////////////////////////////////////////////
// DataItem by Function Group Type:
// 0x00 MODE_TRANSITION
// 0x04 CPU_FUNCTIONS
// 0x08 TYPE_RES
//
// DataItem by Function Type:
// 0x00 PUSH
// 0x04 REQUEST
// 0x08 RESPONSE
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

[discriminatedType S7PayloadUserDataItem(uint 4 cpuFunctionGroup, uint 4 cpuFunctionType, uint 8 cpuSubfunction)
    [simple     DataTransportErrorCode returnCode]
    [simple     DataTransportSize      transportSize]
    [simple         uint 16                dataLength]
    //[implicit   uint 16                dataLength    'lengthInBytes - 4']

    [typeSwitch cpuFunctionGroup, cpuFunctionType, cpuSubfunction, dataLength

        ['0x02', '0x00', '0x01' S7PayloadUserDataItemCyclicServicesPush
            [simple uint 16 itemsCount]
            [array AssociatedValueType items count 'itemsCount']
        ]

        ['0x02', '0x00', '0x05' S7PayloadUserDataItemCyclicServicesChangeDrivenPush
            [simple uint 16 itemsCount]
            [array AssociatedQueryValueType items count 'itemsCount']
        ]

        ['0x02', '0x04', '0x01' S7PayloadUserDataItemCyclicServicesSubscribeRequest
            [simple uint 16 itemsCount]
            [simple TimeBase timeBase]
            [simple uint 8 timeFactor]
            [array CycServiceItemType item count 'itemsCount']
        ]

        ['0x02', '0x04', '0x04' S7PayloadUserDataItemCyclicServicesUnsubscribeRequest
            [simple  uint 8  function]
            [simple  uint 8  jobId]
        ]

        ['0x02', '0x08', '0x01' S7PayloadUserDataItemCyclicServicesSubscribeResponse
            [simple uint 16 itemsCount]
            [array AssociatedValueType items count 'itemsCount']
        ]

        ['0x02', '0x08', '0x04' S7PayloadUserDataItemCyclicServicesUnsubscribeResponse
        ]

        ['0x02', '0x08', '0x05', '0x00' S7PayloadUserDataItemCyclicServicesErrorResponse
        ]

        ['0x02', '0x08', '0x05'  S7PayloadUserDataItemCyclicServicesChangeDrivenSubscribeResponse
            [simple uint 16 itemsCount]
            [array AssociatedQueryValueType  items count 'itemsCount']
        ]

        //USER and SYSTEM Messages
        ['0x04', '0x00', '0x03' S7PayloadDiagnosticMessage
            [simple uint 16     EventId]
            [simple uint 8      PriorityClass]
            [simple uint 8      ObNumber]
            [simple uint 16     DatId]
            [simple uint 16     Info1]
            [simple uint 32     Info2]
            [simple DateAndTime TimeStamp]
        ]

        //PUSH message reception S7300 & S7400 (ALARM_SQ, ALARM_S, ALARM_SC, ...)
        ['0x04', '0x00', '0x05' S7PayloadAlarm8
            [simple AlarmMessagePushType alarmMessage]
        ]
        ['0x04', '0x00', '0x06' S7PayloadNotify
            [simple AlarmMessagePushType alarmMessage]
        ]
        ['0x04', '0x00', '0x0c' S7PayloadAlarmAckInd
            [simple AlarmMessageAckPushType alarmMessage]
        ]
        ['0x04', '0x00', '0x11' S7PayloadAlarmSQ
            [simple AlarmMessagePushType alarmMessage]
        ]
        ['0x04', '0x00', '0x12' S7PayloadAlarmS
            [simple AlarmMessagePushType alarmMessage]
        ]
        ['0x04', '0x00', '0x13' S7PayloadAlarmSC
            [simple AlarmMessagePushType alarmMessage]
        ]
        ['0x04', '0x00', '0x16' S7PayloadNotify8
            [simple AlarmMessagePushType alarmMessage]
        ]

        //Request for specific functions of the SZL system
        ['0x04','0x04', '0x01', '0x00' S7PayloadUserDataItemCpuFunctionReadSzlNoDataRequest
        ]

        ['0x04', '0x04', '0x01' S7PayloadUserDataItemCpuFunctionReadSzlRequest
            [simple   SzlId                  szlId]
            [simple   uint 16                szlIndex]
        ]

        //['0x04', '0x08', '0x01' S7PayloadUserDataItemCpuFunctionReadSzlResponse
        //    [simple   SzlId           szlId]
        //    [simple   uint 16         szlIndex]
        //    [const    uint 16         szlItemLength 28]
        //    [implicit uint 16         szlItemCount  'COUNT(items)']
        //    [array    SzlDataTreeItem items         count 'szlItemCount']
        //]

        ['0x04', '0x08', '0x01' S7PayloadUserDataItemCpuFunctionReadSzlResponse(uint 16 dataLength)
            [array byte items count 'dataLength']
        ]

        //Subscription to PUSH messages
        ['0x04', '0x04', '0x02' S7PayloadUserDataItemCpuFunctionMsgSubscriptionRequest
            [simple   uint 8         Subscription]
            [reserved uint 8         '0x00']
            [simple   string         64             magicKey           ]
            [optional AlarmStateType Alarmtype    'Subscription >= 128']
            [optional uint 8         Reserve      'Subscription >= 128']
        ]

	['0x04', '0x08', '0x02', '0x00' S7PayloadUserDataItemCpuFunctionMsgSubscriptionResponse]

        ['0x04', '0x08', '0x02', '0x02' S7PayloadUserDataItemCpuFunctionMsgSubscriptionSysResponse
            [simple uint 8 result]
            [simple uint 8 reserved01]
        ]

        ['0x04', '0x08', '0x02', '0x05' S7PayloadUserDataItemCpuFunctionMsgSubscriptionAlarmResponse
            [simple uint 8    result]
            [simple uint 8    reserved01]
            [simple AlarmType alarmType]
            [simple uint 8    reserved02]
            [simple uint 8    reserved03]
        ]

        //ALARM_ACK Acknowledgment of alarms
        ['0x04', '0x04', '0x0b' S7PayloadUserDataItemCpuFunctionAlarmAckRequest
            [const    uint 8       functionId       0x09]
            [implicit uint 8                    numberOfObjects 'COUNT(messageObjects)']
            [array    AlarmMessageObjectAckType messageObjects  count 'numberOfObjects' ]
        ]

        ['0x04', '0x08', '0x0b', '0x00' S7PayloadUserDataItemCpuFunctionAlarmAckErrorResponse
        ]

        ['0x04', '0x08', '0x0b' S7PayloadUserDataItemCpuFunctionAlarmAckResponse
            [simple    uint 8 functionId]
            [implicit  uint 8 numberOfObjects 'COUNT(messageObjects)']
            [array     uint 8 messageObjects  count 'numberOfObjects' ]
        ]

        //ALARM_QUERY Request for alarms stored in the controller
        ['0x04', '0x04', '0x13' S7PayloadUserDataItemCpuFunctionAlarmQueryRequest
            [const    uint 8       functionId       0x00]
            [const    uint 8       numberMessageObj 0x01]
            [const    uint 8       variableSpec     0x12]
            [const    uint 8       length           0x08]
            [simple   SyntaxIdType syntaxId]
            [reserved uint 8       '0x00']
            [simple   QueryType    queryType]
            [reserved uint 8       '0x34']
            [simple   AlarmType    alarmType]
        ]

        ['0x04', '0x08', '0x13' S7PayloadUserDataItemCpuFunctionAlarmQueryResponse(uint 16 dataLength)
            [array byte items count 'dataLength']
        ]

    ]
]

[dataIo DataItem(vstring dataProtocolId, int 32 stringLength)
    [typeSwitch dataProtocolId
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['"IEC61131_BOOL"' BOOL
            [reserved uint 7 '0x00']
            [simple   bit    value]
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['"IEC61131_BYTE"' BYTE
            [simple uint 8 value]
        ]
        // 2 byte (16 bit)
        ['"IEC61131_WORD"' WORD
            [simple uint 16 value]
        ]
        // 4 byte (32 bit)
        ['"IEC61131_DWORD"' DWORD
            [simple uint 32 value]
        ]
        // 8 byte (64 bit)
        ['"IEC61131_LWORD"' LWORD
            [simple uint 64 value]
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['"IEC61131_SINT"' SINT
            [simple int 8 value]
        ]
        ['"IEC61131_USINT"' USINT
            [simple uint 8 value]
        ]
        // 16 bit:
        ['"IEC61131_INT"' INT
            [simple int 16 value]
        ]
        ['"IEC61131_UINT"' UINT
            [simple uint 16 value]
        ]
        // 32 bit:
        ['"IEC61131_DINT"' DINT
            [simple int 32 value]
        ]
        ['"IEC61131_UDINT"' UDINT
            [simple uint 32 value]
        ]
        // 64 bit:
        ['"IEC61131_LINT"' LINT
            [simple int 64 value]
        ]
        ['"IEC61131_ULINT"' ULINT
            [simple uint 64 value]
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['"IEC61131_REAL"' REAL
            [simple float 32  value]
        ]
        ['"IEC61131_LREAL"' LREAL
            [simple float 64 value]
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['"IEC61131_CHAR"' CHAR
            [simple string 8 value encoding='"UTF-8"']
        ]
        ['"IEC61131_WCHAR"' CHAR
            [simple string 16 value encoding='"UTF-16"']
        ]
        ['"IEC61131_STRING"' STRING
            // TODO: Fix this length
            [manual vstring value  'STATIC_CALL("parseS7String", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("serializeS7String", writeBuffer, _value, stringLength, _type.encoding)' 'STR_LEN(_value) + 2' encoding='"UTF-8"']
        ]
        ['"IEC61131_WSTRING"' STRING
            // TODO: Fix this length
            [manual vstring value 'STATIC_CALL("parseS7String", readBuffer, stringLength, _type.encoding)' 'STATIC_CALL("serializeS7String", writeBuffer, _value, stringLength, _type.encoding)' '(STR_LEN(_value) * 2) + 2' encoding='"UTF-16"']
        ]

        // -----------------------------------------
        // TIA Date-Formats
        // -----------------------------------------
        // - Duration: Interpreted as "milliseconds"
        ['"IEC61131_TIME"' TIME
            [simple uint 32 milliseconds]
        ]
        //['"S7_S5TIME"' TIME
        //    [reserved uint 2  '0x00']
        //    [uint     uint 2  'base']
        //    [simple   uint 12 value]
        //]
        // - Duration: Interpreted as "number of nanoseconds"
        ['"IEC61131_LTIME"' LTIME
            [simple uint 64 nanoseconds]
        ]
        // - Date: Interpreted as "number of days since 1990-01-01"
        ['"IEC61131_DATE"' DATE
            [simple uint 16 daysSinceSiemensEpoch]
            // Number of days between 1990-01-01 and 1970-01-01 according to https://www.timeanddate.com/
            //[virtual uint 16 daysSinceEpoch 'daysSinceSiemensEpoch + 7305']
        ]
        //['"IEC61131_LDATE"' LDATE
        //    [implicit uint 16 daysSinceSiemensEpoch 'daysSinceEpoch - 7305']
        //    [virtual uint 16 daysSinceEpoch 'daysSinceSiemensEpoch + 7305']
        //]
        // - Time: Interpreted as "milliseconds since midnight (0:00)"
        ['"IEC61131_TIME_OF_DAY"' TIME_OF_DAY
            [simple uint 32 millisecondsSinceMidnight]
        ]
        // - Time: Interpreted as "nanoseconds since midnight (0:00)"
        ['"IEC61131_LTIME_OF_DAY"' LTIME_OF_DAY
            [simple uint 64 nanosecondsSinceMidnight]
        ]
        // - Date & Time: interpreted as individual components.
        ['"IEC61131_DATE_AND_TIME"' DATE_AND_TIME
            [simple uint 16 year]
            [simple uint 8  month]
            [simple uint 8  day]
            [simple uint 8  dayOfWeek]
            [simple uint 8  hour]
            [simple uint 8  minutes]
            [simple uint 8  seconds]
            [simple uint 32 nanoseconds]
        ]
        // - Date & Time: Interpreted as "number of nanoseconds since 1990-01-01"
        //['"IEC61131_LDATE_AND_TIME"' LDATE_AND_TIME
        //    [implicit uint 16 nanosecondsSinceSiemensEpoch 'nanosecondsSinceEpoch ...']
        //    [virtual uint 16 nanosecondsSinceEpoch 'nanosecondsSinceSiemensEpoch ...']
        //]
    ]
]

[enum uint 8 COTPTpduSize(uint 16 sizeInBytes)
    ['0x07' SIZE_128 ['128']]
    ['0x08' SIZE_256 ['256']]
    ['0x09' SIZE_512 ['512']]
    ['0x0a' SIZE_1024 ['1024']]
    ['0x0b' SIZE_2048 ['2048']]
    ['0x0c' SIZE_4096 ['4096']]
    ['0x0d' SIZE_8192 ['8192']]
]

[enum uint 8 COTPProtocolClass
    ['0x00' CLASS_0]
    ['0x10' CLASS_1]
    ['0x20' CLASS_2]
    ['0x30' CLASS_3]
    ['0x40' CLASS_4]
]

[enum uint 8 DataTransportSize(bit sizeInBits)
    ['0x00' NULL            ['false']]
    ['0x03' BIT             ['true' ]]
    ['0x04' BYTE_WORD_DWORD ['true' ]]
    ['0x05' INTEGER         ['true' ]]
    ['0x06' DINTEGER        ['false']]
    ['0x07' REAL            ['false']]
    ['0x09' OCTET_STRING    ['false']]
]

[enum uint 8 DeviceGroup
    ['0x01' PG_OR_PC]
    ['0x02' OS      ]
    ['0x03' OTHERS  ]
]

[enum uint 8 TransportSize(uint 8 code, uint 8 shortName, uint 8 sizeInBytes, TransportSize baseType, DataTransportSize dataTransportSize, vstring dataProtocolId  , bit supported_S7_300, bit supported_S7_400, bit supported_S7_1200, bit supported_S7_1500, bit supported_LOGO)
    // Bit Strings
    ['0x01' BOOL          ['0x01'     , 'X'             , '1'               , 'null'                , 'BIT'                              , 'IEC61131_BOOL'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x02' BYTE          ['0x02'     , 'B'             , '1'               , 'null'                , 'BYTE_WORD_DWORD'                  , 'IEC61131_BYTE'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x03' WORD          ['0x04'     , 'W'             , '2'               , 'null'                , 'BYTE_WORD_DWORD'                  , 'IEC61131_WORD'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x04' DWORD         ['0x06'     , 'D'             , '4'               , 'WORD'                , 'BYTE_WORD_DWORD'                  , 'IEC61131_DWORD'        , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x05' LWORD         ['0x00'     , 'X'             , '8'               , 'null'                , 'null'                             , 'IEC61131_LWORD'        , 'false'             , 'false'             , 'false'              , 'true'               , 'false'             ]]

    // Integer values
    // INT and UINT moved out of order as the enum constant INT needs to be generated before it's used in java
    ['0x06' INT           ['0x05'     , 'W'             , '2'               , 'null'                , 'INTEGER'                          , 'IEC61131_INT'          , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x07' UINT          ['0x05'     , 'W'             , '2'               , 'INT'                 , 'INTEGER'                          , 'IEC61131_UINT'         , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]
    // ...
    ['0x08' SINT          ['0x02'     , 'B'             , '1'               , 'INT'                 , 'BYTE_WORD_DWORD'                  , 'IEC61131_SINT'         , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]
    ['0x09' USINT         ['0x02'     , 'B'             , '1'               , 'INT'                 , 'BYTE_WORD_DWORD'                  , 'IEC61131_USINT'        , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]
    ['0x0A' DINT          ['0x07'     , 'D'             , '4'               , 'INT'                 , 'INTEGER'                          , 'IEC61131_DINT'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x0B' UDINT         ['0x07'     , 'D'             , '4'               , 'INT'                 , 'INTEGER'                          , 'IEC61131_UDINT'        , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]
    ['0x0C' LINT          ['0x00'     , 'X'             , '8'               , 'INT'                 , 'null'                             , 'IEC61131_LINT'         , 'false'             , 'false'             , 'false'              , 'true'               , 'false'             ]]
    ['0x0D' ULINT         ['0x00'     , 'X'             , '16'              , 'INT'                 , 'null'                             , 'IEC61131_ULINT'        , 'false'             , 'false'             , 'false'              , 'true'               , 'false'             ]]

    // Floating point values
    ['0x0E' REAL          ['0x08'     , 'D'             , '4'               , 'null'                , 'REAL'                             , 'IEC61131_REAL'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x0F' LREAL         ['0x30'     , 'X'             , '8'               , 'REAL'                , 'null'                             , 'IEC61131_LREAL'        , 'false'             , 'false'             , 'true'               , 'true'               , 'false'             ]]

    // Characters and Strings
    ['0x10' CHAR          ['0x03'     , 'B'             , '1'               , 'null'                , 'BYTE_WORD_DWORD'                  , 'IEC61131_CHAR'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x11' WCHAR         ['0x13'     , 'X'             , '2'               , 'null'                , 'null'                             , 'IEC61131_WCHAR'        , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]
    ['0x12' STRING        ['0x03'     , 'X'             , '1'               , 'null'                , 'BYTE_WORD_DWORD'                  , 'IEC61131_STRING'       , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x13' WSTRING       ['0x00'     , 'X'             , '2'               , 'null'                , 'null'                             , 'IEC61131_WSTRING'      , 'false'             , 'false'             , 'true'               , 'true'               , 'true'              ]]

    // Dates and time values (Please note that we seem to have to rewrite queries for these types to reading bytes or we'll get "Data type not supported" errors)
    ['0x14' TIME          ['0x0B'     , 'X'             , '4'                 , 'null'                  , 'null'                         , 'IEC61131_TIME'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    //['0x15' S5TIME        ['0x0C'    , 'X'             , '4'                 , 'null'                  , 'null'                          , 'S7_S5TIME'             , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x16' LTIME         ['0x00'     , 'X'             , '8'                 , 'TIME'                  , 'null'                         , 'IEC61131_LTIME'        , 'false'             , 'false'             , 'false'              , 'true'               , 'false'             ]]
    ['0x17' DATE          ['0x09'     , 'X'             , '2'                 , 'null'                  , 'BYTE_WORD_DWORD'              , 'IEC61131_DATE'         , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x18' TIME_OF_DAY   ['0x06'     , 'X'             , '4'                 , 'null'                  , 'BYTE_WORD_DWORD'              , 'IEC61131_TIME_OF_DAY'  , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x19' TOD           ['0x06'     , 'X'             , '4'                 , 'null'                  , 'BYTE_WORD_DWORD'              , 'IEC61131_TIME_OF_DAY'  , 'true'              , 'true'              , 'true'               , 'true'               , 'true'              ]]
    ['0x1A' DATE_AND_TIME ['0x0F'     , 'X'             , '12'                , 'null'                  , 'null'                         , 'IEC61131_DATE_AND_TIME', 'true'              , 'true'              , 'false'              , 'true'               , 'false'             ]]
    ['0x1B' DT            ['0x0F'     , 'X'             , '12'                , 'null'                  , 'null'                         , 'IEC61131_DATE_AND_TIME', 'true'              , 'true'              , 'false'              , 'true'               , 'false'             ]]
]

[enum uint 8 MemoryArea(string 24 shortName)
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

[enum uint 8 DataTransportSize(bit sizeInBits)
    ['0x00' NULL                ['false']]
    ['0x03' BIT                 ['true']]
    ['0x04' BYTE_WORD_DWORD     ['true']]
    ['0x05' INTEGER             ['true']]
    ['0x06' DINTEGER            ['false']]
    ['0x07' REAL                ['false']]
    ['0x09' OCTET_STRING        ['false']]
]

[enum uint 8 DataTransportErrorCode
    ['0x00' RESERVED               ]
    ['0xFF' OK                     ]
    ['0x03' ACCESS_DENIED          ]
    ['0x05' INVALID_ADDRESS        ]
    ['0x06' DATA_TYPE_NOT_SUPPORTED]
    ['0x0A' NOT_FOUND              ]
]

[enum uint 4 SzlModuleTypeClass
    ['0x0' CPU]
    ['0x4' IM]
    ['0x8' FM]
    ['0xC' CP]
]

[enum uint 8 SzlSublist
    ['0x00' NONE]    
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
    ['0x71' H_CPU_GROUP_INFORMATION]
    ['0x74' STATUS_SINGLE_MODULE_LED]
    ['0x75' SWITCHED_DP_SLAVES_H_SYSTEM]
    ['0x90' DP_MASTER_SYSTEM_INFORMATION]
    ['0x91' MODULE_STATUS_INFORMATION]
    ['0x92' RACK_OR_STATION_STATUS_INFORMATION]
    ['0x94' RACK_OR_STATION_STATUS_INFORMATION_2]
    ['0x95' ADDITIONAL_DP_MASTER_SYSTEM_OR_PROFINET_IO_SYSTEM_INFORMATION]
    ['0x96' MODULE_STATUS_INFORMATION_PROFINET_IO_AND_PROFIBUS_DP]
    ['0x9C' TOOL_CHANGER_INFORMATION_PROFINET]
    ['0xA0' DIAGNOSTIC_BUFFER]
    ['0xB1' MODULE_DIAGNOSTIC_INFORMATION_DR0]
    ['0xB2' MODULE_DIAGNOSTIC_INFORMATION_DR1_GI]
    ['0xB3' MODULE_DIAGNOSTIC_INFORMATION_DR1_LA]
    ['0xB4' DIAGNOSTIC_DATA_DP_SLAVE]
]

[enum uint 8 CpuSubscribeEvents
    ['0x01' CPU]
    ['0x02' IM]
    ['0x04' FM]
    ['0x80' CP]
]

[enum uint 8 EventType
    ['0x01' MODE]
    ['0x02' SYS]
    ['0x04' USR]
    ['0x80' ALM]
    ['0x69' CYC] //Not from s7 standar, only for internal processing.
]

[enum uint 8 SyntaxIdType
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

[enum uint 8 AlarmType
    ['0x01' SCAN]
    ['0x02' ALARM_8]
    ['0x04' ALARM_S]
]

[enum uint 8 AlarmStateType
    ['0x00' SCAN_ABORT]
    ['0x01' SCAN_INITIATE]
    ['0x04' ALARM_ABORT]
    ['0x05' ALARM_INITIATE]
    ['0x08' ALARM_S_ABORT]
    ['0x09' ALARM_S_INITIATE]
]

[enum uint 8 QueryType
    ['0x01' BYALARMTYPE]
    ['0x02' ALARM_8]
    ['0x04' ALARM_S]
    ['0x09' ALARM_8P] //Under test with S7-400 PLC
]

[enum uint 8 ModeTransitionType
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

[enum uint 8 'TimeBase'
    ['0x00' B01SEC]
    ['0x01' B1SEC]
    ['0X02' B10SEC]
]

