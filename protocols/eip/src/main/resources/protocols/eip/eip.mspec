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

[discriminatedType EipPacket (bit response)
    [discriminator uint 16 command                                                                                     ]
    [implicit      uint 16 packetLength 'lengthInBytes - 24'                                                           ]
    [simple        uint 32 sessionHandle                                                                               ]
    [simple        uint 32 status                                                                                      ]
    [array         byte    senderContext count '8'                                                                     ]
    [simple        uint 32 options                                                                                     ]
    [typeSwitch command,response,packetLength
        ['0x0001','false'   NullCommandRequest
        ]
        ['0x0001','true'    NullCommandResponse
        ]
        ['0x0004','false' ListServicesRequest
        ]
        ['0x0004','true','0' NullListServicesResponse
        ]
        ['0x0004','true' ListServicesResponse
            [implicit uint    16   typeIdCount  'COUNT(typeIds)'                                                       ]
            [array    TypeId       typeIds   count   'typeIdCount'                                                     ]
        ]
        ['0x0065','false' EipConnectionRequest
            [const    uint    16   protocolVersion   0x01                                                              ]
            [const    uint    16   flags             0x00                                                              ]
        ]
        ['0x0065','true','0' NullEipConnectionResponse
        ]
        ['0x0065','true' EipConnectionResponse
            [const    uint    16   protocolVersion   0x01                                                              ]
            [const    uint    16   flags             0x00                                                              ]
        ]
        ['0x0066' EipDisconnectRequest
        ]
        ['0x006F' CipRRData
            [simple   uint    32   interfaceHandle                                                                     ]
            [simple   uint    16   timeout                                                                             ]
            [implicit uint    16   typeIdCount  'COUNT(typeIds)'                                                       ]
            [array    TypeId       typeIds   count   'typeIdCount'                                                     ]
        ]
        ['0x0070' SendUnitData
            [const    uint    32   interfaceHandle  0x00000000                                                         ]
            [simple   uint    16   timeout                                                                             ]
            [implicit uint    16   typeIdCount  'COUNT(typeIds)'                                                       ]
            [array    TypeId       typeIds   count   'typeIdCount'                                                     ]
        ]
    ]
]

[discriminatedType  TypeId
    [discriminator  uint    16  id                                                                                     ]
    [typeSwitch id
        ['0x0000'   NullAddressItem
            [reserved       uint    16  '0x0000'                                                                       ]
        ]
        ['0x0100'   ServicesResponse
            [implicit       uint    16     serviceLen 'lengthInBytes - 4'                                              ]
            [simple         uint    16     encapsulationProtocol                                                       ]
            [reserved       uint    2      '0x00'                                                                      ]
            [simple         bit            supportsCIPEncapsulation                                                    ]
            [reserved       uint    12     '0x00'                                                                      ]
            [simple         bit            supportsUDP                                                                 ]
            [array          byte           data        count  'serviceLen - 4'                                         ]
        ]
        ['0x00A1'   ConnectedAddressItem
            [reserved       uint    16  '0x0004'                                                                       ]
            [simple         uint    32  connectionId                                                                   ]
        ]
        ['0x00B1'   ConnectedDataItem
            [implicit       uint    16  packetSize 'service.lengthInBytes + 2'                                         ]
            [simple         uint    16  sequenceCount                                                                  ]
            [simple         CipService('true', 'packetSize - 2')    service                                            ]
        ]
        ['0x00B2'   UnConnectedDataItem
            [implicit       uint    16  packetSize 'service.lengthInBytes'                                             ]
            [simple         CipService('false', 'packetSize')    service                                               ]
        ]
    ]
]

[discriminatedType  CipService(bit connected, uint 16 serviceLen)
    [discriminator  bit     response                                                                                   ]
    [discriminator  uint    7   service                                                                                ]
    [typeSwitch service,response,connected
        ['0x01','false' GetAttributeAllRequest
            [implicit   uint   8            requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
            [simple     PathSegment         classSegment                                                               ]
            [simple     PathSegment         instanceSegment                                                            ]
        ]
        ['0x01','true' GetAttributeAllResponse
            [reserved   uint   8            '0x00'                                                                     ]
            [simple     uint   8            status                                                                     ]
            [simple     uint   8            extStatus                                                                  ]
            [optional   CIPAttributes('serviceLen - 4')          attributes '(serviceLen - 4) > 0'                     ]
        ]
        ['0x4C','false' CipReadRequest
            [implicit   uint    8           requestPathSize 'COUNT(tag) / 2'                                           ]
            [array      byte                tag           count  '(requestPathSize * 2)'                               ]
            [simple     uint    16          elementNb                                                                  ]
        ]
        ['0x4C','true' CipReadResponse
            [reserved   uint    8           '0x00'                                                                     ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           extStatus                                                                  ]
            [optional   CIPData('serviceLen - 4')   data    '(serviceLen - 4) > 0'                                      ]
        ]
        ['0x4D','false' CipWriteRequest
            [implicit   uint    8           requestPathSize 'COUNT(tag) / 2'                                           ]
            [array      byte                tag   length  'requestPathSize * 2'                                        ]
            [simple     CIPDataTypeCode     dataType                                                                   ]
            [simple     uint    16          elementNb                                                                  ]
            [array      byte                data  length  'dataType.size * elementNb'                                  ]
        ]
        ['0x4D','true' CipWriteResponse
            [reserved   uint        8       '0x00'                                                                     ]
            [simple     uint        8       status                                                                     ]
            [simple     uint        8       extStatus                                                                  ]
        ]
        ['0x4E','false' CipConnectionManagerCloseRequest
            [simple     uint     8          requestPathSize                                                            ]
            [simple     PathSegment         classSegment                                                               ]
            [simple     PathSegment         instanceSegment                                                            ]
            [simple     uint    4           priority                                                                   ]
            [simple     uint    4           tickTime                                                                   ]
            [simple     uint    8           timeoutTicks                                                               ]
            [simple     uint    16          connectionSerialNumber                                                     ]
            [simple     uint    16          originatorVendorId                                                         ]
            [simple     uint    32          originatorSerialNumber                                                     ]
            [simple     uint    8           connectionPathSize                                                         ]
            [reserved   byte                '0x00'                                                                     ]
            [array      PathSegment         connectionPaths     terminated   'STATIC_CALL("noMorePathSegments", readBuffer)']
        ]
        ['0x4E','true' CipConnectionManagerCloseResponse
            [reserved   uint    8           '0x00'                                                                     ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           additionalStatusWords                                                      ]
            [simple     uint    16          connectionSerialNumber                                                     ]
            [simple     uint    16          originatorVendorId                                                         ]
            [simple     uint    32          originatorSerialNumber                                                     ]
            [simple     uint    8           applicationReplySize                                                       ]
            [reserved   uint    8           '0x00'                                                                     ]
        ]
        ['0x0A','false' MultipleServiceRequest
            [const      uint    8           requestPathSize   0x02                                                     ]
            [const      uint    32          requestPath       0x01240220                                               ]   //Logical Segment: Class(0x20) 0x02, Instance(0x24) 01 (Message Router)
            [simple Services('serviceLen-6') data                                                                      ]
        ]
        ['0x0A','true' MultipleServiceResponse
            [reserved   uint    8           '0x0'                                                                      ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           extStatus                                                                  ]
            [simple     uint    16          serviceNb                                                                  ]
            [array      uint    16          offsets       count  'serviceNb'                                           ]
            [array      byte   servicesData count 'serviceLen - 6 - (2 * serviceNb)'                                   ]
        ]
        ['0x52','false','false'   CipUnconnectedRequest
            [implicit   uint    8           requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
            [simple     PathSegment         classSegment                                                               ]
            [simple     PathSegment         instanceSegment                                                            ]
            [reserved   uint    16          '0x9D05'                                                                   ]   //Timeout 5s
            [implicit   uint    16          messageSize   'lengthInBytes - 10 - 4'                                     ]   //subtract above and routing
            [simple     CipService('false','messageSize')  unconnectedService                                          ]
            [const      uint    16          route 0x0001                                                               ]
            [simple     int     8           backPlane                                                                  ]
            [simple     int     8           slot                                                                       ]
        ]
        ['0x52','false','true'   CipConnectedRequest
            [implicit   uint    8           requestPathSize 'COUNT(pathSegments) / 2'                                  ]
            [array      byte                pathSegments    count 'requestPathSize * 2'                                ]
            [reserved   uint    16          '0x0001'                                                                   ]
            [reserved   uint    32          '0x00000000'                                                               ]
        ]
        ['0x52','true'   CipConnectedResponse
            [reserved   uint    8           '0x00'                                                                     ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           additionalStatusWords                                                      ]
            [optional   CIPDataConnected    data    '(serviceLen - 4) > 0'                                             ]
        ]
        ['0x5B','false'     CipConnectionManagerRequest
            [implicit    uint    8          requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
            [simple      PathSegment        classSegment                                                               ]
            [simple      PathSegment        instanceSegment                                                            ]
            [simple      uint    4          priority                                                                   ]
            [simple      uint    4          tickTime                                                                   ]
            [simple      uint    8          timeoutTicks                                                               ]
            [simple      uint    32         otConnectionId                                                             ]
            [simple      uint    32         toConnectionId                                                             ]
            [simple      uint    16         connectionSerialNumber                                                     ]
            [simple      uint    16         originatorVendorId                                                         ]
            [simple      uint    32         originatorSerialNumber                                                     ]
            [simple      uint    8          timeoutMultiplier                                                          ]
            [reserved    uint    24         '0x000000'                                                                 ]
            [simple      uint    32         otRpi                                                                      ]
            [simple      NetworkConnectionParameters otConnectionParameters                                            ]
            [simple      uint    32         toRpi                                                                      ]
            [simple      NetworkConnectionParameters toConnectionParameters                                            ]
            [simple      TransportType      transportType                                                              ]
            [simple      uint    8          connectionPathSize                                                         ]
            [array       PathSegment        connectionPaths terminated  'STATIC_CALL("noMorePathSegments", readBuffer)']
        ]
        ['0x5B','true'     CipConnectionManagerResponse
            [reserved    uint    24         '0x000000'                                                                 ]
            [simple      uint    32         otConnectionId                                                             ]
            [simple      uint    32         toConnectionId                                                             ]
            [simple      uint    16         connectionSerialNumber                                                     ]
            [simple      uint    16         originatorVendorId                                                         ]
            [simple      uint    32         originatorSerialNumber                                                     ]
            [simple      uint    32         otApi                                                                      ]
            [simple      uint    32         toApi                                                                      ]
            [implicit    uint    8          replySize   'lengthInBytes - 30'                                           ]
            [reserved    uint    8          '0x00'                                                                     ]
        ]
    ]
]

[discriminatedType PathSegment
    [discriminator  uint    3   pathSegment         ]
    [typeSwitch pathSegment
        ['0x00'      PortSegment
            [simple PortSegmentType     segmentType]
        ]
        ['0x01'      LogicalSegment
            [simple LogicalSegmentType  segmentType]
        ]
        ['0x04'      DataSegment
            [simple DataSegmentType     segmentType]
        ]
    ]
]

[discriminatedType  PortSegmentType
    [discriminator     bit         extendedLinkAddress]
    [typeSwitch extendedLinkAddress
        ['false'    PortSegmentNormal
            [simple     uint    4   port]
            [simple     uint    8   linkAddress]
        ]
        ['true'     PortSegmentExtended
            [simple     uint    4   port]
            [simple   uint    8   linkAddressSize]
            [virtual  uint 8    paddingByte 'linkAddressSize % 2']
            [simple   vstring  '(linkAddressSize * 8) + (paddingByte * 8)'    address]
        ]
    ]
]

[discriminatedType LogicalSegmentType
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

[discriminatedType DataSegmentType
    [discriminator  uint    5   dataSegmentType]
    [typeSwitch dataSegmentType
        ['0x11'      AnsiExtendedSymbolSegment
            [implicit   uint    8                   dataSize    'STR_LEN(symbol)'         ]
            [simple     vstring 'dataSize * 8'      symbol                                ]
            [optional   uint    8                   pad         'STR_LEN(symbol) % 2 != 0']
        ]
    ]
]

[type   CIPAttributes(uint 16 packetLength)
    [implicit   uint    16              numberOfClasses 'COUNT(classId)']
    [array      uint    16              classId count   'numberOfClasses']
    [simple     uint    16              numberAvailable]
    [simple     uint    16              numberActive]
    [array      byte                    data count 'packetLength - 2 - (COUNT(classId) * 2) - 4']
]

[type   CIPData(uint 16 packetLength)
    [simple     CIPDataTypeCode     dataType]
    [array      byte                data  count  'packetLength - 2']
]

[type   CIPDataConnected
    [simple     uint    32   value]
    [simple     uint    16   tagStatus]
]

[type   InstanceSegment
    [simple     uint    3    pathSegmentType]
    [simple     uint    3    logicalSegmentType]
    [simple     uint    2    logicalSegmentFormat]
    [simple     uint    8    instance]
]

[type   ClassSegment
    [simple     uint    3   pathSegmentType]
    [simple     uint    3   logicalSegmentType]
    [simple     uint    2   logicalSegmentFormat]
    [simple     uint    8   classSegment]
]


[type   NetworkConnectionParameters
   [simple      uint    16  connectionSize]
   [reserved    uint    8   '0x00']
   [simple      bit         owner]
   [simple      uint    2   connectionType]
   [reserved    bit         'false']
   [simple      uint    2   priority]
   [simple      bit         connectionSizeType]
   [reserved    bit         'false']
]

[type   TransportType
   [simple      bit         direction]
   [simple      uint    3   trigger]
   [simple      uint    4   classTransport]
]

[type   Services  (uint   16   servicesLen)
    [implicit   uint        16  serviceNb 'COUNT(offsets)']
    [array      uint        16  offsets       count  'serviceNb']
    [array      CipService('false', 'servicesLen / serviceNb')   services    count  'serviceNb' ]
]

[enum uint   16   CIPDataTypeCode(uint 8  size)
    ['0X00C1'   BOOL            ['1']]
    ['0X00C2'   SINT            ['1']]
    ['0X00C3'   INT             ['2']]
    ['0X00C4'   DINT            ['4']]
    ['0X00C5'   LINT            ['8']]
    ['0X00C6'   USINT           ['1']]
    ['0X00C7'   UINT            ['2']]
    ['0X00C8'   UDINT           ['4']]
    ['0X00C9'   ULINT           ['8']]
    ['0X00CA'   REAL            ['4']]
    ['0X00CB'   LREAL           ['8']]
    ['0X00CC'   STIME           ['4']] // Synchronous time information
    ['0X00CD'   DATE            ['2']]
    ['0X00CE'   TIME_OF_DAY     ['4']]
    ['0X00CF'   DATE_AND_TIME   ['6']]
    ['0X00D0'   STRING          ['0']] // Character string, 1-byte per character, 2-byte length
    ['0X00D1'   BYTE            ['1']]
    ['0X00D2'   WORD            ['2']]
    ['0X00D3'   DWORD           ['4']]
    ['0X00D3'   LWORD           ['8']]
    ['0X00D5'   STRING2         ['0']] // Character string, 2-bytes per character
    ['0X00D6'   FTIME           ['4']] // Duration - high resolution
    ['0X00D7'   LTIME           ['8']] // Duration - long
    ['0X00D8'   ITIME           ['2']] // Duration - short
    ['0X00D9'   STRINGN         ['0']] // Character string, n-bytes per character
    ['0X00DA'   SHORT_STRING    ['0']] // Character string, 1-byte per character, 1-byte length
    ['0X00DB'   TIME            ['4']] // Duration - milliseconds
    ['0X00DC'   EPATH           ['0']] // CIP path segments
    ['0X00DD'   ENGUNIT         ['0']] // Engineering units
    ['0X00DD'   STRINGI         ['0']] // International character string
    // ARRAY
    ['0X02A0'   STRUCTURED      ['88']]
    //['0X02A0'   STRING          ['88']]
    //['0X02A0'   STRING36        ['40']]
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

[enum   uint    16  CIPClassID
    ['0x0001'   Identity]
    ['0x0002'   MessageRouter]
    ['0x0006'   ConnectionManager]
]