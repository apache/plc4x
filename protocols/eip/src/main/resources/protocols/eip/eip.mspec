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
        ['0x01','true'  GetAttributeAllResponse
            [reserved   uint   8            '0x00'                                                                     ]
            [simple     uint   8            status                                                                     ]
            [simple     uint   8            extStatus                                                                  ]
            [optional   CIPAttributes('serviceLen - 4')          attributes '(serviceLen - 4) > 0'                     ]
        ]
        ['0x02','false' SetAttributeAllRequest
            // TODO: Implement
        ]
        ['0x02','true'  SetAttributeAllResponse
            // TODO: Implement
        ]
        ['0x03','false' GetAttributeListRequest
            // TODO: Implement
        ]
        ['0x03','true'  GetAttributeListResponse
            // TODO: Implement
        ]
        ['0x04','false' SetAttributeListRequest
            // TODO: Implement
        ]
        ['0x04','true'  SetAttributeListResponse
            // TODO: Implement
        ]
        ['0x0A','false' MultipleServiceRequest
            [const      uint    8           requestPathSize   0x02                                                     ]
            [const      uint    32          requestPath       0x01240220                                               ]   //Logical Segment: Class(0x20) 0x02, Instance(0x24) 01 (Message Router)
            [simple Services('serviceLen-6') data                                                                      ]
        ]
        ['0x0A','true'  MultipleServiceResponse
            [reserved   uint    8           '0x0'                                                                      ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           extStatus                                                                  ]
            [simple     uint    16          serviceNb                                                                  ]
            [array      uint    16          offsets       count  'serviceNb'                                           ]
            [array      byte   servicesData count 'serviceLen - 6 - (2 * serviceNb)'                                   ]
        ]
        ['0x0E','false' GetAttributeSingleRequest
            // TODO: Implement
        ]
        ['0x0E','true'  GetAttributeSingleResponse
            // TODO: Implement
        ]
        ['0x10','false' SetAttributeSingleRequest
            // TODO: Implement
        ]
        ['0x10','true'  SetAttributeSingleResponse
            // TODO: Implement
        ]
        ['0x4C','false' CipReadRequest
            [implicit   uint    8           requestPathSize 'COUNT(tag) / 2'                                           ]
            [array      byte                tag           count  '(requestPathSize * 2)'                               ]
            [simple     uint    16          elementNb                                                                  ]
        ]
        ['0x4C','true'  CipReadResponse
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
        ['0x4D','true'  CipWriteResponse
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
        ['0x4E','true'  CipConnectionManagerCloseResponse
            [reserved   uint    8           '0x00'                                                                     ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           additionalStatusWords                                                      ]
            [simple     uint    16          connectionSerialNumber                                                     ]
            [simple     uint    16          originatorVendorId                                                         ]
            [simple     uint    32          originatorSerialNumber                                                     ]
            [simple     uint    8           applicationReplySize                                                       ]
            [reserved   uint    8           '0x00'                                                                     ]
        ]
        ['0x52','false','false'  CipUnconnectedRequest
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
        ['0x52','true'  CipConnectedResponse
            [reserved   uint    8           '0x00'                                                                     ]
            [simple     uint    8           status                                                                     ]
            [simple     uint    8           additionalStatusWords                                                      ]
            [optional   CIPDataConnected    data    '(serviceLen - 4) > 0'                                             ]
        ]
        ['0x5B','false' CipConnectionManagerRequest
            [implicit   uint    8          requestPathSize '(classSegment.lengthInBytes + instanceSegment.lengthInBytes)/2']
            [simple     PathSegment        classSegment                                                                ]
            [simple     PathSegment        instanceSegment                                                             ]
            [simple     uint    4          priority                                                                    ]
            [simple     uint    4          tickTime                                                                    ]
            [simple     uint    8          timeoutTicks                                                                ]
            // ot = Originator (Client) Target (Server)
            [simple     uint    32         otConnectionId                                                              ]
            // to = Target (Server) Originator (Client)
            [simple     uint    32         toConnectionId                                                              ]
            [simple     uint    16         connectionSerialNumber                                                      ]
            [simple     uint    16         originatorVendorId                                                          ]
            [simple     uint    32         originatorSerialNumber                                                      ]
            [simple     uint    8          timeoutMultiplier                                                           ]
            [reserved   uint    24         '0x000000'                                                                  ]
            // ot = Originator (Client) Target (Server)
            [simple     uint    32         otRpi                                                                       ]
            [simple     NetworkConnectionParameters otConnectionParameters                                             ]
            // to = Target (Server) Originator (Client)
            [simple     uint    32         toRpi                                                                       ]
            [simple     NetworkConnectionParameters toConnectionParameters                                             ]
            [simple     TransportType      transportType                                                               ]
            [simple     uint    8          connectionPathSize                                                          ]
            [array      PathSegment        connectionPaths terminated  'STATIC_CALL("noMorePathSegments", readBuffer)' ]
        ]
        ['0x5B','true'  CipConnectionManagerResponse
            [reserved   uint    24         '0x000000'                                                                  ]
            // ot = Originator (Client) Target (Server)
            [simple     uint    32         otConnectionId                                                              ]
            // to = Target (Server) Originator (Client)
            [simple     uint    32         toConnectionId                                                              ]
            [simple     uint    16         connectionSerialNumber                                                      ]
            [simple     uint    16         originatorVendorId                                                          ]
            [simple     uint    32         originatorSerialNumber                                                      ]
            // ot = Originator (Client) Target (Server)
            [simple     uint    32         otApi                                                                       ]
            // to = Target (Server) Originator (Client)
            [simple     uint    32         toApi                                                                       ]
            [implicit   uint    8          replySize   'lengthInBytes - 30'                                            ]
            [reserved   uint    8          '0x00'                                                                      ]
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
    [implicit   uint    16              numberOfClasses 'COUNT(classId)'                             ]
    [array      uint    16              classId count   'numberOfClasses'                            ]
    [optional   uint    16              numberAvailable 'packetLength >= ((numberOfClasses * 2) + 4)']
    [optional   uint    16              numberActive    'packetLength >= ((numberOfClasses * 2) + 6)']
    [array      byte                    data    count   '(packetLength > ((numberOfClasses * 2) + 6)) ? packetLength - ((numberOfClasses * 2) + 6) : 0']
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
    ['0x00000000'   Success                      ]
    ['0x00000001'   ConnectionFailure            ]
    ['0x00000002'   ResourceUnAvailable          ]
    ['0x00000003'   InvalidParameterValue        ]
    ['0x00000004'   PathSegmentError             ]
    ['0x00000005'   PathDestinationUnknown       ]
    ['0x00000006'   PartialTransfer              ]
    ['0x00000007'   ConnectionIDNotValid         ]
    ['0x00000008'   ServiceNotSupported          ]
    ['0x00000009'   InvalidAttributeValue        ]
    ['0x0000000A'   AttributeListError           ]
    ['0x0000000B'   AlreadyInRequestedState      ]
    ['0x0000000C'   ObjectStateConflict          ]
    ['0x0000000D'   ObjectAlreadyExists          ]
    ['0x0000000E'   AttributeNotSettable         ]
    ['0x0000000F'   PrivilegeViolation           ]
    ['0x00000010'   DeviceStateConflict          ]
    ['0x00000011'   ReplyDataTooLarge            ]
    ['0x00000012'   FragmentationOfPrimitiveValue]
    ['0x00000013'   NotEnoughData                ]
    ['0x00000014'   AttributeNotSupported        ]
    ['0x00000015'   TooMuchData                  ]
    ['0x00000016'   ObjectDoesNotExist           ]
    ['0x00000017'   ServiceFragmentation         ]
    ['0x00000018'   NoStoredAttributeData        ]
    ['0x00000019'   StoreOperationFailure        ]
    ['0x0000001A'   RequestPacketTooLarge        ]
    ['0x0000001B'   ResponsePacketTooLarge       ]
    ['0x0000001C'   MissingAttributeListEntryData]
    ['0x0000001D'   InvalidAttributeValueList    ]
    ['0x0000001E'   EmbeddedServiceError         ]
    ['0x0000001F'   VendorSpecificError          ]
]

[enum   uint    16  CIPClassID
    // General Use Objects
    ['0x0001'   Identity                    ]
    ['0x0002'   MessageRouter               ]
    ['0x0004'   Assembly                    ]
    ['0x0005'   Connection                  ]
    ['0x0006'   ConnectionManager           ]
    ['0x0007'   Register                    ]
    ['0x000F'   Parameter                   ]
    ['0x0010'   ParameterGroup              ]
    ['0x002B'   AcknowledgeHandler          ]
    ['0x002E'   Selection                   ]
    ['0x0037'   File                        ]
    ['0x0045'   OriginatorConnectionList    ]
    ['0x00F3'   ConnectionConfiguration     ]
    ['0x00F4'   Port                        ]
    // Application Specific Objects
    ['0x0008'   DiscreteInputPoint          ]
    ['0x0009'   DiscreteOutputPoint         ]
    ['0x000A'   AnalogInputPoint            ]
    ['0x000B'   AnalogOutputPoint           ]
    ['0x000E'   PresenceSensing             ]
    ['0x0012'   Group                       ]
    ['0x001D'   DiscreteInputGroup          ]
    ['0x001E'   DiscreteOutputGroup         ]
    ['0x001F'   DiscreteGroup               ]
    ['0x0020'   AnalogInputGroup            ]
    ['0x0021'   AnalogOutputGroup           ]
    ['0x0022'   AnalogGroup                 ]
    ['0x0023'   PositionSensor              ]
    ['0x0024'   PositionControllerSupervisor]
    ['0x0025'   PositionController          ]
    ['0x0026'   BlockSequencer              ]
    ['0x0027'   CommandBlock                ]
    ['0x0028'   MotorData                   ]
    ['0x0029'   ControlSupervisor           ]
    ['0x002A'   AcDcDrive                   ]
    ['0x002C'   Overload                    ]
    ['0x0030'   SDeviceSupervisor           ]
    ['0x0031'   SAnalogSensor               ]
    ['0x0032'   SAnalogActuator             ]
    ['0x0033'   SSingleStageController      ]
    ['0x0034'   SGasCalibration             ]
    ['0x0038'   SPartialPressure            ]
    ['0x0040'   SSensorCalibration          ]
    ['0x0041'   EventLog                    ]
    ['0x0042'   MotionDeviceAxis            ]
    ['0x004A'   SafetyAnalogInputGroup      ]
    ['0x004E'   BaseEnergy                  ]
    ['0x004F'   ElectricalEnergy            ]
    ['0x0050'   NonElectricalEnergy         ]
    ['0x0053'   PowerManagementObject       ]
    ['0x005C'   PowerCurtailmentObject      ]
    // Network Specific Objects
    ['0x0003'   DeviceNet                   ]
    ['0x0044'   Modbus                      ]
    ['0x0046'   ModbusSerialLink            ]
    ['0x0047'   DeviceLevelRing             ]
    ['0x0048'   QOS                         ]
    ['0x004C'   Sercos3Link                 ]
    ['0x0051'   BaseSwitch                  ]
    ['0x0052'   Snmp                        ]
    ['0x0053'   PowerManagement             ]
    ['0x0054'   RstpBridge                  ]
    ['0x0055'   RstpPort                    ]
    ['0x0056'   ParallelRedundancyProtocol  ]
    ['0x0057'   PrpNodesTable               ]
    ['0x00F0'   ControlNet                  ]
    ['0x00F1'   ControlNetKeeper            ]
    ['0x00F2'   ControlNetScheduling        ]
    ['0x00F5'   TcpIpInterface              ]
    ['0x00F6'   EthernetLink                ]
    ['0x00F7'   CompoNetLink                ]
    ['0x00F8'   CompoNetRepeater            ]
]

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Types related to CIP segments (EIP-CIP Volume 1 Appendix C-1.4)
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*[discriminatedType CipSegment
    [discriminator SegmentType segmentType                                                                             ]
    [typeSwitch segmentType
        ['PortSegment'         *PortSegment
            [simple   bit     extendedLinkAddressSize                                                                  ]
            [simple   uint 4  portIdentifier                                                                           ]
            [optional uint 8  linkAddressSize         'extendedLinkAddressSize'                                        ]
            [optional uint 16 extendedPortIdentifier  'portIdentifier == 15'                                           ]
            [array    byte    linkAddress             count             'extendedLinkAddressSize ? 1 : linkAddressSize']
        ]
        ['LogicalSegment'      *LogicalSegment
            [simple        LogicalType   logicalType                                                                   ]
            [discriminator LogicalFormat logicalFormat                                                                 ]
            [typeSwitch logicalFormat
                ['LogicalAddress8Bit' *8Bit
                    [simple  uint 8 value                                                                              ]
                ]
                ['LogicalAddress8Bit' *16Bit
// TODO: If were using "padded", we need padding added, without this is "compact" format
//                    [const   uint 8  pad   '0x00'                                                                      ]
                    [simple  uint 16 value                                                                             ]
                ]
                ['LogicalAddress8Bit' *32Bit
// TODO: If were using "padded", we need padding added, without this is "compact" format
//                    [const   uint 8  pad   '0x00'                                                                      ]
                    [simple  uint 32 value                                                                             ]
                ]
            ]
        ]
        ['NetworkSegment'      *NetworkSegment
            [discriminator NetworkSegmentSubType networkSegmentSubType                                                 ]
            [typeSwitch networkSegmentSubType
                ['ScheduleSegment'        *ScheduleSegment
                    // TODO: Appendix C-2.1 of Volume 4
                ]
                ['FixedTagSegment'        *FixedTagSegment
                    // TODO: Appendix C-2.2 of Volume 4
                ]
                ['ProductionInhibitTime'  *ProductionInhibitTime
                    [simple   uint 8  minInterval                                                                      ]
                ]
                ['SafetySegment'          *SafetySegment
                    // TODO: Appendix C of Volume 5
                ]
                ['ExtendedNetworkSegment' *ExtendedNetworkSegment
                    [implicit uint 8  numWords                      'COUNT(data) / 2'                                  ]
                    [simple   uint 16 extendedNetworkSegmentSubtype                                                    ]
                    [array    byte    data     count                'numWords * 2'                                     ]
                ]
            ]
        ]
        ['SymbolicSegment'     *SymbolicSegment
            [peek uint 5                          symbolSizeInBytes                                                    ]
            [typeSwitch symbolSizeInBytes
                ['0' *ExtendedStringFormat
                    [const         uint 5         symbolLength      0x00                                               ]
                    [discriminator SymbolicFormat symbolicFormat                                                       ]
                    [typeSwitch symbolicFormat
                        ['SymbolicFormat16BitCharacters' *SymbolicFormat16BitCharacters
                        ]
                        ['SymbolicFormat24BitCharacters' *SymbolicFormat24BitCharacters
                        ]
                        ['SymbolicFormatNumeric'         *SymbolicFormatNumeric
                            [discriminator SymbolicFormatNumericSubType numericFormat                                  ]
                            [typeSwitch numericFormat
                                ['USINT' *USINT
                                    [simple uint 8  value                                                              ]
                                ]
                                ['UINT'  *UINT
                                    [simple uint 16 value                                                              ]
                                ]
                                ['UDINT' *UDINT
                                    [simple uint 32 value                                                              ]
                                ]
                            ]
                        ]
                    ]
                ]
                [    *StandardStringFormat
                    [implicit uint 5                     symbolLength 'STR_LEN(symbol)'                                ]
                    [simple   vstring 'symbolLength * 8' symbol                                                        ]
                ]
            ]
        ]
        ['DataSegment'         *DataSegment
            [discriminator DataSegmentSubType dataSegmentSubType                                                       ]
            [typeSwitch dataSegmentSubType
                ['Simple'             *Simple
                    [implicit uint 8  numWords 'COUNT(data)'                                                           ]
                    [array    uint 16 data     count          'numWords'                                               ]
                ]
                ['AnsiExtendedSymbol' *AnsiExtendedSymbol
                    [implicit uint 8                   dataLength 'STR_LEN(data)'                                      ]
                    [simple   vstring 'dataLength * 8' data                                                            ]
                ]
            ]
        ]
        ['DataTypeConstructed' *DataTypeConstructed
        ]
        ['DataTypeElementary'  *DataTypeElementary
        ]
    ]
]

[enum uint 3 SegmentType
    ['0x0' PortSegment                ]
    ['0x1' LogicalSegment             ]
    ['0x2' NetworkSegment             ]
    ['0x3' SymbolicSegment            ]
    ['0x4' DataSegment                ]
    ['0x5' DataTypeConstructed        ]
    ['0x6' DataTypeElementary         ]
    ['0x7' Reserved                   ]
]

[enum uint 5 NetworkSegmentSubType
    ['0x00' Reserved_0                ]
    ['0x01' ScheduleSegment           ]
    ['0x02' FixedTagSegment           ]
    ['0x03' ProductionInhibitTime     ]
    ['0x04' Reserved_4                ]
    ['0x05' Reserved_5                ]
    ['0x06' Reserved_6                ]
    ['0x07' Reserved_7                ]
    ['0x08' Reserved_8                ]
    ['0x09' Reserved_9                ]
    ['0x0A' Reserved_10               ]
    ['0x0B' Reserved_11               ]
    ['0x0C' Reserved_12               ]
    ['0x0D' Reserved_13               ]
    ['0x0E' Reserved_14               ]
    ['0x0F' Reserved_15               ]
    ['0x10' SafetySegment             ]
    ['0x11' Reserved_17               ]
    ['0x12' Reserved_18               ]
    ['0x13' Reserved_19               ]
    ['0x14' Reserved_20               ]
    ['0x15' Reserved_21               ]
    ['0x16' Reserved_22               ]
    ['0x17' Reserved_23               ]
    ['0x18' Reserved_24               ]
    ['0x19' Reserved_25               ]
    ['0x1A' Reserved_26               ]
    ['0x1B' Reserved_27               ]
    ['0x1C' Reserved_28               ]
    ['0x1D' Reserved_29               ]
    ['0x1E' Reserved_30               ]
    ['0x1F' ExtendedNetworkSegment    ]
]

[enum uint 3 LogicalType
    ['0x0' ClassId                    ]
    ['0x1' InstanceId                 ]
    ['0x2' MemberId                   ]
    ['0x3' ConnectionPoint            ]
    ['0x4' AttributeId                ]
    ['0x5' Special                    ]
    ['0x6' ServiceId                  ]
    ['0x7' Reserved                   ]
]

[enum uint 2 LogicalFormat
    ['0x0' LogicalAddress8Bit         ]
    ['0x1' LogicalAddress16Bit        ]
    ['0x2' LogicalAddress32Bit        ]
    ['0x3' Reserved                   ]
]

[enum uint 3 SymbolicFormat
    ['1' SymbolicFormat16BitCharacters]
    ['2' SymbolicFormat24BitCharacters]
    ['6' SymbolicFormatNumeric        ]
]

[enum uint 5 SymbolicFormatNumericSubType
    ['0x00' Reserved_0                ]
    ['0x01' Reserved_1                ]
    ['0x02' Reserved_2                ]
    ['0x03' Reserved_3                ]
    ['0x04' Reserved_4                ]
    ['0x05' Reserved_5                ]
    ['0x06' USINT                     ]
    ['0x07' UINT                      ]
    ['0x08' UDINT                     ]
    ['0x09' Reserved_9                ]
    ['0x0A' Reserved_10               ]
    ['0x0B' Reserved_11               ]
    ['0x0C' Reserved_12               ]
    ['0x0D' Reserved_13               ]
    ['0x0E' Reserved_14               ]
    ['0x0F' Reserved_15               ]
    ['0x10' Reserved_16               ]
    ['0x11' Reserved_17               ]
    ['0x12' Reserved_18               ]
    ['0x13' Reserved_19               ]
    ['0x14' Reserved_20               ]
    ['0x15' Reserved_21               ]
    ['0x16' Reserved_22               ]
    ['0x17' Reserved_23               ]
    ['0x18' Reserved_24               ]
    ['0x19' Reserved_25               ]
    ['0x1A' Reserved_26               ]
    ['0x1B' Reserved_27               ]
    ['0x1C' Reserved_28               ]
    ['0x1D' Reserved_29               ]
    ['0x1E' Reserved_30               ]
    ['0x1F' Reserved_31               ]
]

[enum uint 5 DataSegmentSubType
    ['0x00' Simple                    ]
    ['0x01' Reserved_1                ]
    ['0x02' Reserved_2                ]
    ['0x03' Reserved_3                ]
    ['0x04' Reserved_4                ]
    ['0x05' Reserved_5                ]
    ['0x06' Reserved_6                ]
    ['0x07' Reserved_7                ]
    ['0x08' Reserved_8                ]
    ['0x09' Reserved_9                ]
    ['0x0A' Reserved_10               ]
    ['0x0B' Reserved_11               ]
    ['0x0C' Reserved_12               ]
    ['0x0D' Reserved_13               ]
    ['0x0E' Reserved_14               ]
    ['0x0F' Reserved_15               ]
    ['0x10' Reserved_16               ]
    ['0x11' AnsiExtendedSymbol        ]
    ['0x12' Reserved_18               ]
    ['0x13' Reserved_19               ]
    ['0x14' Reserved_20               ]
    ['0x15' Reserved_21               ]
    ['0x16' Reserved_22               ]
    ['0x17' Reserved_23               ]
    ['0x18' Reserved_24               ]
    ['0x19' Reserved_25               ]
    ['0x1A' Reserved_26               ]
    ['0x1B' Reserved_27               ]
    ['0x1C' Reserved_28               ]
    ['0x1D' Reserved_29               ]
    ['0x1E' Reserved_30               ]
    ['0x1F' Reserved_31               ]
]
*/