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

//////////////////////////////////////////////////////////////////
///EthernetIP Header of size 24
/////////////////////////////////////////////////////////////////

[discriminatedType EipPacket (IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator uint 16 command]
    // TODO: was len before but that is a reserved keyword in golang and as long as we don't have a language agnostic neutralizer this clashes
    [implicit      uint 16 packetLength 'lengthInBytes - 24']
    [simple        uint 32 sessionHandle]
    [simple        uint 32 status]
    [array         byte    senderContext count '8']
    [simple        uint 32 options]
    [typeSwitch command
            ['0x0065' EipConnectionRequest
                [const  uint    16   protocolVersion   0x01]
                [const  uint    16   flags             0x00]
            ]
            ['0x0066' EipDisconnectRequest
            ]
            ['0x006F' CipRRData
                [simple     uint    32    interfaceHandle]
                [simple     uint    16    timeout]
                [simple     uint    16    itemCount]
                [array      TypeId('order')         typeId   count   'itemCount']
            ]
            ['0x0070' SendUnitData
                [const      uint    32     interfaceHandle  0x00000000]
                [simple     uint    16     timeout]
                [simple     uint    16     itemCount]
                [array      TypeId('order')         typeId   count   'itemCount']
            ]
        ]
]

[discriminatedType  TypeId(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator  uint    16  id]
    [typeSwitch id
        ['0x0000'   NullAddressItem
            [reserved       uint    16  '0x0000']
        ]
        ['0x00A1'   ConnectedAddressItem
            [reserved       uint    16  '0x0004']
            [simple         uint    32  connectionId]
        ]
        ['0x00B1'   ConnectedDataItem
            [implicit       uint    16  packetSize 'service.lengthInBytes + 2']
            [simple         uint    16  sequenceCount]
            [simple         CipService('true', 'packetSize', 'order')    service]
        ]
        ['0x00B2'   UnConnectedDataItem
            [implicit       uint    16  packetSize 'service.lengthInBytes']
            [simple         CipService('false', 'packetSize', 'order')    service]
        ]
    ]
]

[discriminatedType  CipService(bit connected, uint 16 serviceLen, IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator  bit     response]
    [discriminator  uint    7   service]
    [typeSwitch service,response,connected
        ['0x4C','false' CipReadRequest
            [implicit   int     8   requestPathSize 'COUNT(tag) / 2']
            [array      byte   tag   count  '(requestPathSize * 2)']
            [simple     uint    16  elementNb]
        ]
        ['0x4C','true' CipReadResponse
              [reserved   uint            8   '0x00']
              [simple     uint            8   status]
              [simple     uint            8   extStatus]
              [simple     CIPDataTypeCode     dataType]
              [array      byte   data  count  'serviceLen - 6']
        ]
        ['0x4D','false' CipWriteRequest
            [implicit   int     8       requestPathSize 'COUNT(tag) / 2']
            [array      byte            tag   length  'requestPathSize * 2']
            [simple     CIPDataTypeCode     dataType]
            [simple     uint    16          elementNb]
            [array      byte            data  length  'dataType.size * elementNb']
        ]
        ['0x4D','true' CipWriteResponse
            [reserved   uint        8   '0x00']
            [simple     uint        8   status]
            [simple     uint        8   extStatus]
        ]
        ['0x4E','false' CipConnectionManagerCloseRequest
              [simple      int     8           requestPathSize]
              [simple      PathSegment('order')         classSegment]
              [simple      PathSegment('order')         instanceSegment]
              [simple      uint    4           priority]
              [simple      uint    4           tickTime]
              [simple      uint    8           timeoutTicks]
              [simple      uint    16          connectionSerialNumber]
              [simple      uint    16          originatorVendorId]
              [simple      uint    32          originatorSerialNumber]
              [simple      uint    8           connectionPathSize]
              [reserved    byte                '0x00']
              [simple      PathSegment('order')         connectionPathPortSegment]
              [simple      PathSegment('order')         connectionPathClassSegment]
              [simple      PathSegment('order')         connectionPathInstanceSegment]
        ]
        ['0x4E','true' CipConnectionManagerCloseResponse
              [reserved uint    8   '0x00']
              [simple   uint    8   status]
              [simple   uint    8   additionalStatusWords]
              [simple      uint    16          connectionSerialNumber]
              [simple      uint    16          originatorVendorId]
              [simple      uint    32          originatorSerialNumber]
              [simple      uint     8          applicationReplySize]
              [reserved    uint     8          '0x00']
        ]
        ['0x0A','false' MultipleServiceRequest
               [const  int     8   requestPathSize   0x02]
               [const  uint    32  requestPath       0x01240220]   //Logical Segment: Class(0x20) 0x02, Instance(0x24) 01 (Message Router)
               [simple Services('serviceLen - 6 ', 'order')  data ]
        ]
        ['0x0A','true' MultipleServiceResponse
               [reserved   uint    8   '0x0']
               [simple     uint    8   status]
               [simple     uint    8   extStatus]
               [simple     uint    16  serviceNb]
               [array      uint    16  offsets       count  'serviceNb']
               [array      byte   servicesData count 'serviceLen - 6 - (2 * serviceNb)']
        ]
        ['0x52','false','false'   CipUnconnectedRequest
               [implicit      int     8         requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
               [simple      PathSegment('order')         classSegment]
               [simple      PathSegment('order')         instanceSegment]
               [reserved   uint    16  '0x9D05']   //Timeout 5s
               [implicit   uint    16  messageSize   'lengthInBytes - 10 - 4']   //subtract above and routing
               [simple     CipService('false','messageSize','order')  unconnectedService ]
               [const      uint    16  route 0x0001]
               [simple     int     8   backPlane]
               [simple     int     8   slot]
        ]
        ['0x52','false','true'   CipConnectedRequest
               [implicit   uint    8    requestPathSize 'COUNT(pathSegments) / 2']
               [array      byte         pathSegments    count 'requestPathSize * 2']
               [reserved   uint    16   '0x0001']
               [reserved   uint    32   '0x00000000']
        ]
        ['0x52','true','true'   CipConnectedResponse
               [reserved   uint    8    '0x00']
               [simple     uint    8    status]
               [simple     uint    8    additionalStatusWords]
               [simple     uint    32   value]
               [simple     uint    16   tagStatus]
        ]
        ['0x5B','false'     CipConnectionManagerRequest
               [implicit      int     8         requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
               [simple      PathSegment('order')         classSegment]
               [simple      PathSegment('order')         instanceSegment]
               [simple      uint    4           priority]
               [simple      uint    4           tickTime]
               [simple      uint    8           timeoutTicks]
               [simple      uint    32          otConnectionId]
               [simple      uint    32          toConnectionId]
               [simple      uint    16          connectionSerialNumber]
               [simple      uint    16          originatorVendorId]
               [simple      uint    32          originatorSerialNumber]
               [simple      uint    8           timeoutMultiplier]
               [reserved    uint    24          '0x000000']
               [simple      uint    32          otRpi]
               [simple      NetworkConnectionParameters('order') otConnectionParameters]
               [simple      uint    32          toRpi]
               [simple      NetworkConnectionParameters('order') toConnectionParameters]
               [simple      TransportType('order')       transportType]
               [simple      uint    8           connectionPathSize]
               [simple      PathSegment('order')         connectionPathPortSegment]
               [simple      PathSegment('order')         connectionPathClassSegment]
               [simple      PathSegment('order')         connectionPathInstanceSegment]

        ]
        ['0x5B','true'     CipConnectionManagerResponse
               [reserved    uint    24          '0x000000']
               [simple      uint    32          otConnectionId]
               [simple      uint    32          toConnectionId]
               [simple      uint    16          connectionSerialNumber]
               [simple      uint    16          originatorVendorId]
               [simple      uint    32          originatorSerialNumber]
               [simple      uint    32          otApi]
               [simple      uint    32          toApi]
               [implicit    uint    8           replySize   'serviceLen - 30']
               [reserved    uint    8           '0x00']
        ]
    ]
]

[discriminatedType PathSegment(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator  uint    3   pathSegment]
    [typeSwitch pathSegment
        ['0x00'      PortSegment
            [simple bit extendedLinkAddress]
            [simple uint 4  port]
            [simple uint    8   linkAddress]
        ]
        ['0x01'      LogicalSegment
            [simple LogicalSegmentType('order')  segmentType]
        ]
        ['0x04'      DataSegment
            [simple DataSegmentType('order') segmentType]
        ]
    ]
]

[discriminatedType LogicalSegmentType(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator  uint    3   logicalSegmentType]
    [typeSwitch logicalSegmentType
        ['0x00' ClassID
            [simple uint    2   format]
            [simple uint    8   segmentClass]
        ]
        ['0x01' InstanceID
            [simple uint    2   format]
            [simple uint    8   instance]
        ]
        ['0x02' MemberID
            [simple uint    2   format]
            [simple uint    8   instance]
        ]
    ]
]

[discriminatedType DataSegmentType(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [discriminator  uint    5   dataSegmentType]
    [typeSwitch dataSegmentType
        ['0x11'      AnsiExtendedSymbolSegment
            [implicit   uint    8   dataSize        'symbol.length']
            [simple     vstring     'dataSize * 8'      symbol]
            [optional   uint    8   pad         'symbol.length % 2 != 0']
        ]
    ]
]

[type   InstanceSegment(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [simple     uint    3   pathSegmentType]
    [simple     uint    3   logicalSegmentType]
    [simple     uint    2   logicalSegmentFormat]
    [simple     uint    8   instance]
]

[type   ClassSegment(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [simple     uint    3   pathSegmentType]
    [simple     uint    3   logicalSegmentType]
    [simple     uint    2   logicalSegmentFormat]
    [simple     uint    8   classSegment]
]


[type   NetworkConnectionParameters(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
   [simple      uint    16  connectionSize]
   [reserved    uint    8   '0x00']
   [simple      bit         owner]
   [simple      uint    2   connectionType]
   [reserved    bit         'false']
   [simple      uint    2   priority]
   [simple      bit         connectionSizeType]
   [reserved    bit         'false']
]

[type   TransportType(IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
   [simple      bit         direction]
   [simple      uint    3   trigger]
   [simple      uint    4   classTransport]
]

[type   Services  (uint   16   servicesLen, IntegerEncoding order) byteOrder='order == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
    [simple uint        16  serviceNb]
    [array  uint        16  offsets       count  'serviceNb']
    [array  CipService('false', 'servicesLen / serviceNb', 'order')   services    count  'serviceNb' ]
]

[enum uint   16   CIPDataTypeCode(uint 8  size)
    ['0X00C1'   BOOL            ['1']]
    ['0X00C2'   SINT            ['1']]
    ['0X00C3'   INT             ['2']]
    ['0X00C4'   DINT            ['4']]
    ['0X00C5'   LINT            ['8']]
    ['0X00CA'   REAL            ['4']]
    ['0X00D3'   DWORD           ['4']]
    ['0X02A0'   STRUCTURED      ['88']]
    ['0X02A0'   STRING          ['88']]
    ['0X02A0'   STRING36        ['40']]
    //TODO: -1 is not a valid value for uint
    //['-1'       UNKNOWN         ['-1']]
]

[enum uint 16 CIPStructTypeCode
    ['0x0FCE'   STRING]
]

[enum   uint    16  EiPCommand
    ['0x0065'   RegisterSession ]
    ['0x0066'   UnregisterSession ]
    ['0x006F'   SendRRData ]
]

[enum uint 4 IntegerEncoding
    ['0x0' BIG_ENDIAN]
    ['0x1' LITTLE_ENDIAN]
]

[enum   uint    32  CIPStatus
    ['0x00000000'   Success]
    ['0x00000001'   ConnectionFailure]
    ['0x00000002'   ResourceUnAvailable]
    ['0x00000003'   InvalidParameterValue]
    ['0x00000004'   PathSegmentError]
    ['0x00000005'   PathDestinationUnknown]
    ['0x00000006'   PartialTransfer]
    ['0x00000007'   ConnectionIDNotValid]
    ['0x00000008'   ServiceNotSupported]
    ['0x00000009'   InvalidAttributeValue]
    ['0x0000000A'   AttributeListError]
    ['0x0000000B'   AlreadyInRequestedState]
    ['0x0000000C'   ObjectStateConflict]
    ['0x0000000D'   ObjectAlreadyExists]
    ['0x0000000E'   AttributeNotSettable]
    ['0x0000000F'   PriviligeViolation]
    ['0x00000010'   DeviceStateConflict]
    ['0x00000011'   ReplyDataTooLarge]
    ['0x00000012'   FragmentationOfPrimitiveValue]
    ['0x00000013'   NotEnoughData]
    ['0x00000014'   AttributeNotSupported]
    ['0x00000015'   TooMuchData]
    ['0x00000016'   ObjectDoesNotExist]
    ['0x00000017'   ServiceFragmentation]
    ['0x00000018'   NoStoredAttributeData]
    ['0x00000019'   StoreOperationFailure]
    ['0x0000001A'   RequestPacketTooLarge]
    ['0x0000001B'   ResponsePacketTooLarge]
    ['0x0000001C'   MissingAttributeListEntryData]
    ['0x0000001D'   InvalidAttributeValueList]
    ['0x0000001E'   EmbeddedServiceError]
    ['0x0000001F'   VendorSpecificError]
]