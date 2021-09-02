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

 /*
 Overview of the Protocols involved in this driver:

 Ethernet
 Udp                                                        (Based on Ethernet)
 DceRpc                                                     (Based on Udp and Ethernet)

 ARP		Address Resolution Protocol                     (Based on Ethernet)
 LLDP		Link Layer Discovery Protocol                   (Based on Ethernet)

 PnDcp		PROFINET Discovery and Configuration Protocol   (Based on Ethernet)
 PnIoCm 	PROFINET IO Context Manager                     (Based on DCP/RPC, UDP and Ethernet)
 PnIo		PROFINET IO                                     (Based on Ethernet)
 PnIoPs 	PROFIsafe protocol                              (Based on Ethernet)
 PnIoAl 	PROFINET Alarm Events                           (Based on Ethernet)

 // Not really handled in this driver (Just listed for the sake of completeness)
 PnPtcp     PROFINET Precision Transparent Clock Protocol
*/

[type 'Ethernet_Frame'
    // When sending to the mac address prefix of 01:0e:cf are multicast packets
    [simple MacAddress            'destination']
    [simple MacAddress            'source'     ]
    [simple Ethernet_FramePayload 'payload'    ]
]

[discriminatedType 'Ethernet_FramePayload'
    [discriminator uint 16 'packetType']
    [typeSwitch 'packetType'
        ['0x0800' Ethernet_FramePayload_IPv4
            [const    uint 4              'version'                         '0x4'                      ]
            [const    uint 4              'headerLength'                    '0x5'                      ]
            [const    uint 6              'differentiatedServicesCodepoint' '0x00'                     ]
            [const    uint 2              'explicitCongestionNotification'  '0x0'                      ]
            [implicit uint 16             'totalLength'                     '20 + packet.lengthInBytes']
            [simple   uint 15             'identification'                                             ]
            [const    uint 3              'flags'                           '0x00'                     ]
            [const    uint 13             'fragmentOffset'                  '0x00'                     ]
            // Time to live: 64
            [const    uint 8              'timeToLive'                      '0x40'                     ]
            // Protocol: UDP
            [const    uint 8              'protocol'                        '0x11'                     ]
            // TODO: Implement
            //[checksum uint 16           'headerChecksum'                                             ]
            [simple   IpAddress           'sourceAddress'                                              ]
            [simple   IpAddress           'destinationAddress'                                         ]
            [simple   Udp_Packet          'packet'                                                     ]
        ]
        ['0x8100' Ethernet_FramePayload_VirtualLan
            [simple VirtualLanPriority    'priority'                                                   ]
            [simple bit                   'ineligible'                                                 ]
            [simple uint 12               'id'                                                         ]
            [simple Ethernet_FramePayload 'payload'                                                    ]
        ]
        ['0x8892' Ethernet_FramePayload_PnDcp
            [simple PnDcp_Pdu             'pdu'                                                        ]
        ]
    ]
]

[type 'Udp_Packet'
    [simple   uint 16       'sourcePort'                                        ]
    [simple   uint 16       'destinationPort'                                   ]
    [implicit uint 16       'packetLength'    'lengthInBytes'                   ]
    // TODO: Implement
    //[checksum uint 16       'headerChecksum'                                    ]
    [simple   DceRpc_Packet 'payload'                                           ]
]

// 4.10.3.2
[discriminatedType 'DceRpc_Packet'
    [const         uint 8            'version'           '0x04'                 ]
    [discriminator DceRpc_PacketType 'packetType'                               ]
    [typeSwitch 'packetType'
        ['REQUEST'  DceRpc_Packet_Req
            [reserved uint 1         '0x0'                                      ]
            [const    uint 1         'broadcast'                        '0'     ]
            [const    uint 1         'idempotent'                       '1'     ]
            [const    uint 1         'maybe'                            '0'     ]
            [const    uint 1         'noFragmentAcknowledgeRequested'   '0'     ]
            [const    uint 1         'fragment'                         '0'     ]
            [const    uint 1         'lastFragment'                     '0'     ]
            [reserved uint 1         '0x0'                                      ]
        ]

        ['RESPONSE' DceRpc_Packet_Res
            [reserved uint 1         '0x0'                                      ]
            [const    uint 1         'broadcast'                        '0'     ]
            [const    uint 1         'idempotent'                       '0'     ]
            [const    uint 1         'maybe'                            '0'     ]
            [const    uint 1         'noFragmentAcknowledgeRequested'   '1'     ]
            [const    uint 1         'fragment'                         '0'     ]
            [const    uint 1         'lastFragment'                     '1'     ]
            [reserved uint 1         '0x0'                                      ]
        ]
    ]
    [reserved      uint 6            '0x00'                                     ]
    [const         uint 1            'cancelWasPending'                 '0'     ]
    [reserved      uint 1            '0x0'                                      ]
    // Byte Order:
    // 0x0 = Big-Endian
    // 0x1 = Little-Endian
    [const         uint 4            'integerEncoding'   '0x0'                  ]
    // Character Type: Ascii
    [const         uint 4            'characterEncoding' '0x0'                  ]
    // Floating Point Type: IEEE
    [const         uint 8            'floatingPointRepresentation' '0x00'       ]
    [const         uint 8            'serialHigh'        '0x00'                 ]
    // FIXME: Strangely the spec seems to also be referencing a 'serialLow' but I can't see that in communictaion
    // 4.10.3.2.8 Coding of the field RPCObjectUUID DEA00000-6C97-11D1-8271-{instanceOrNodeNumber}{deviceId}{vendorId}
    // Apache Vendor Id: 0x060B
    // PLC4X Profinet Driver Device ID (can be chosen freely): 0xCAFE
    // NOTE: We can get the Device-Id and Vendor-Id from the PN-DCP search result of the browser.
    [const         uint 32           'uuid1'             '0xDEA00000'           ]
    [const         uint 16           'uuid2'             '0x6C97'               ]
    [const         uint 16           'uuid3'             '0x11D1'               ]
    [const         uint 16           'uuid4'             '0x8271'               ]
    [simple        uint 16           'instanceOrNodeNumber'                     ]
    [simple        uint 16           'deviceId'                                 ]
    [simple        uint 16           'vendorId'                                 ]
    // 4.10.3.2.9
    // Device Interface:            DEA00001-6C97-11D1-8271-00A02442DF7D
    // Controller Interface:        DEA00002-6C97-11D1-8271-00A02442DF7D
    // Supervisor Interface:        DEA00003-6C97-11D1-8271-00A02442DF7D
    // Parameter Server Interface:  DEA00004-6C97-11D1-8271-00A02442DF7D
    [const         uint 32           'interface1'        '0xDEA00001'           ]
    [const         uint 16           'interface2'        '0x6C97'               ]
    [const         uint 16           'interface3'        '0x11D1'               ]
    [const         uint 16           'interface4'        '0x8271'               ]
    [const         uint 16           'interface5'        '0x00A0'               ]
    [const         uint 32           'interface6'        '0x2442DF7D'           ]
    // 4.10.3.2.10
    // The Controller and the Device generate the uuid for each AR (Application Relationship) and use them as long as the AR exists
    [simple        uint 32           'activity'                                 ]
    [const         uint 16           'activity2'         '0x0000'               ]
    [const         uint 16           'activity3'         '0x1010'               ]
    [const         uint 16           'activity4'         '0x77BE'               ]
    [const         uint 32           'activity5'         '0x3D3C6D60'           ]
    [const         uint 16           'activity6'         '0xA3A9'               ]
    [simple        uint 32           'serverBootTime'                           ]
    [const         uint 32           'interfaceVer'      '0x00000001'           ]
    [simple        uint 32           'sequenceNumber'                           ]
    [simple        DceRpc_Operation  'operation'                                ]
    [const         uint 16           'interfaceHint'     '0xFFFF'               ]
    [const         uint 16           'activityHint'      '0xFFFF'               ]
    [implicit      uint 16           'fragmentLength'    'payload.lengthInBytes']
    [const         uint 16           'fragmentNum'       '0x0000'               ]
    [const         uint 8            'authProto'         '0x00'                 ]
    [const         uint 8            'serialLow'         '0x00'                 ]
    [simple        PnIoCm_Packet     'payload'           ['packetType']         ]
]

[type 'Uuid'
    [array byte 'data' count '16']
]

// There are some special MAC addresses reserved:
// 01-0E-CF-00-00-00:      As destination for DCP-Identify-ReqPDU (The FrameID is then required to be set to 0xFEFE)
// 01-0E-CF-00-00-01:      As destination for DCP-Helo-ReqPDU (The FrameID is then required to be set to 0xFEFC)
// 01-0E-CF-00-00-02:      Reserved
// 01-0E-CF-00-01-00:      Reserved for further multicast addresses within the Type 10 context
// 01-0E-CF-00-01-01:      As multicast destination for RT_CLASS_3
// 01-0E-CF-00-01-02:      As invalid frame multicast destination for RT_CLASS_3
// 01-0E-CF-00-01-03 - FF: Reserved for further multicast addresses within the Type 10 context

[type 'MacAddress'
    [array uint 8 'address' count '6']
]

[type 'IpAddress'
    [array uint 8 'data' count '4']
]

// 4.10.3.2.2
[enum uint 8 'DceRpc_PacketType'
    ['0x00' REQUEST              ]
    ['0x01' PING                 ]
    ['0x02' RESPONSE             ]
    ['0x03' FAULT                ]
    ['0x04' WORKING              ]
    // Response to PING
    ['0x05' NO_CALL              ]
    ['0x06' REJECT               ]
    ['0x07' ACKNOWLEDGE          ]
    ['0x08' CONNECTIONLESS_CANCEL]
    ['0x09' FRAGMENT_ACKNOWLEDGE ]
    ['0x0A' CANCEL_ACKNOWLEDGE   ]
]

[enum uint 16 'DceRpc_Operation'
    ['0x0000' CONNECT]
    ['0x0003' WRITE  ]
]

// https://de.wikipedia.org/wiki/IEEE_802.1p
[enum uint 3 'VirtualLanPriority'   [string '2' 'acronym']
    ['0x0' BEST_EFFORT              ['BE'                ]]
    ['0x1' BACKGROUND               ['BK'                ]]
    ['0x2' EXCELLENT_EFFORT         ['EE'                ]]
    ['0x3' CRITICAL_APPLICATIONS    ['CA'                ]]
    ['0x4' VIDEO                    ['VI'                ]]
    ['0x5' VOICE                    ['VO'                ]]
    ['0x6' INTERNETWORK_CONTROL     ['IC'                ]]
    ['0x7' NETWORK_CONTROL          ['NC'                ]]
]

/////////////////////////////////////////////////////////////////////////////////////////
//
//   PROFINET DCP
//
// Discovery and basic configuration
//
/////////////////////////////////////////////////////////////////////////////////////////

// Page 90
[discriminatedType 'PnDcp_Pdu'
    [discriminator PnDcp_FrameId     'frameId'                           ]
    [discriminator PnDcp_ServiceId   'serviceId'                         ]
    [simple        PnDcp_ServiceType 'serviceType'                       ]
    // 4.3.1.3.4 (Page 95)
    [simple        uint 32           'xid'                               ]
    // 4.3.1.3.5 (Page 95ff)
    [simple        uint 16           'responseDelayFactorOrPadding'      ]
    // 4.3.1.3.4 (Page 95)
    [implicit      uint 16           'dcpDataLength' 'lengthInBytes - 12']
    [typeSwitch 'frameId','serviceId','serviceType.response'
        ////////////////////////////////////////////////////////////////////////////
        // Multicast (Well theoretically)
        ////////////////////////////////////////////////////////////////////////////
        // The Identify request is valid in two options:
        // 1) One containing only an AllSelectorBlock
        // 2) One containing optionally either NameOfStationBlock or AliasNameBlock and another optional IdentifyReqBlock
        // (I assume, that if in case 2 both optionally aren't used, this might not be valid and option 1 should be sent instead)
        ['DCP_Identify_ReqPDU','IDENTIFY','false' PnDcp_Pdu_IdentifyReq [uint 16 'dcpDataLength']
            [array PnDcp_Block 'blocks' length 'dcpDataLength'           ]
        ]

        // Response to a Identify request
        ['DCP_Identify_ResPDU','IDENTIFY','true' PnDcp_Pdu_IdentifyRes [uint 16 'dcpDataLength']
            [array PnDcp_Block 'blocks' length 'dcpDataLength'           ]
        ]

        // Packet a Profinet station might emit once it is turned on
        ['DCP_Hello_ReqPDU','HELLO','false' PnDcp_Pdu_HelloReq
//            [simple NameOfStationBlockRes    'nameOfStationBlockRes'   ]
//            [simple IPParameterBlockRes      'iPParameterBlockRes'     ]
//            [simple DeviceIdBlockRes         'deviceIdBlockRes'        ]
//            [simple DeviceVendorBlockRes     'deviceVendorBlockRes'    ]
//            [simple DeviceOptionsBlockRes    'deviceOptionsBlockRes'   ]
//            [simple DeviceRoleBlockRes       'deviceRoleBlockRes'      ]
//            [simple DeviceInitiativeBlockRes 'deviceInitiativeBlockRes']
        ]

        ////////////////////////////////////////////////////////////////////////////
        // Unicast
        ////////////////////////////////////////////////////////////////////////////

        ['DCP_GetSet_PDU','GET','false' PnDcp_Pdu_GetReq
//            [simple GetReqBlock              'getReqBlock'             ]
        ]
        ['DCP_GetSet_PDU','GET','true' PnDcp_Pdu_GetRes
//            [simple GetResBlock              'getResBlock'             ]
//            [simple GetNegResBlock           'getNegResBlock'          ]
        ]

        ['DCP_GetSet_PDU','SET','false' PnDcp_Pdu_SetReq
//            [simple StartTransactionBlock    'startTransactionBlock'   ]
//            [simple BlockQualifier           'blockQualifier'          ]
//            [simple SetResetReqBlock         'setResetReqBlock'        ]
//            [simple SetReqBlock              'setReqBlock'             ]
//            [simple StopTransactionBlock     'stopTransactionBlock'    ]
//            [simple BlockQualifier           'blockQualifier'          ]
        ]
        ['DCP_GetSet_PDU','SET','true' PnDcp_Pdu_SetRes
//            [simple SetResBlock              'setResBlock'             ]
//            [simple SetNegResBlock           'setNegResBlock'          ]
        ]
    ]
]

[discriminatedType 'PnDcp_Block'
    [discriminator PnDcp_BlockOptions 'option'                   ]
    [discriminator uint 8       'suboption'                      ]
    [implicit      uint 16      'blockLength' 'lengthInBytes - 4']
    [typeSwitch 'option','suboption'

        ////////////////////////////////////////////////////////////////////////////
        // IP_OPTION
        ////////////////////////////////////////////////////////////////////////////

        // 4.3.1.4.1 (Page 97)
        ['IP_OPTION','1' PnDcp_Block_IpMacAddress
            [reserved uint 16  '0x0000'                                                 ]
            [simple MacAddress 'macAddress'                                             ]
        ]
        ['IP_OPTION','2' PnDcp_Block_IpParameter
            // 4.3.1.4.12 (Page 105ff)
            [reserved uint 8 '0x00'                                                     ]
            [simple   bit    'ipConflictDetected'                                       ]
            [reserved uint 5 '0x00'                                                     ]
            [simple   bit    'setViaDhcp'                                               ]
            [simple   bit    'setManually'                                              ]
            [array    uint 8 'ipAddress'       count '4'                                ]
            [array    uint 8 'subnetMask'      count '4'                                ]
            [array    uint 8 'standardGateway' count '4'                                ]
        ]
        ['IP_OPTION','3' PnDcp_Block_FullIpSuite
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DEVICE_PROPERTIES_OPTION
        ////////////////////////////////////////////////////////////////////////////

        ['DEVICE_PROPERTIES_OPTION','1' PnDcp_Block_DevicePropertiesDeviceVendor [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            // TODO: Figure out how to do this correctly.
            [array    byte        'deviceVendorValue' count 'blockLength-2'             ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", deviceVendorValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','2' PnDcp_Block_DevicePropertiesNameOfStation [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            // TODO: Figure out how to do this correctly.
            [array    byte        'nameOfStation' count 'blockLength-2'                 ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", nameOfStation) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','3' PnDcp_Block_DevicePropertiesDeviceId
            [reserved uint 16 '0x0000'                                                  ]
            [simple   uint 16 'vendorId'                                                ]
            [simple   uint 16 'deviceId'                                                ]
        ]
        ['DEVICE_PROPERTIES_OPTION','4' PnDcp_Block_DevicePropertiesDeviceRole
            [reserved uint 20 '0x000000'                                                ]
            [simple   bit     'pnioSupervisor'                                          ]
            [simple   bit     'pnioMultidevive'                                         ]
            [simple   bit     'pnioController'                                          ]
            [simple   bit     'pnioDevice'                                              ]
            [reserved uint 8  '0x00'                                                    ]
        ]
        // Contains a list of option combinations the device supports.
        ['DEVICE_PROPERTIES_OPTION','5' PnDcp_Block_DevicePropertiesDeviceOptions [uint 16 'blockLength']
            [reserved uint 16               '0x0000'                                    ]
            [array    PnDcp_SupportedDeviceOption 'supportedOptions' length 'blockLength - 2' ]
        ]
        ['DEVICE_PROPERTIES_OPTION','6' PnDcp_Block_DevicePropertiesAliasName [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            [array    byte        'aliasNameValue' count 'blockLength-2'                ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", aliasNameValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','7' PnDcp_Block_DevicePropertiesDeviceInstance
            [reserved uint 16 '0x0000'                                                  ]
            [simple   uint 8  'deviceInstanceHigh'                                      ]
            [simple   uint 8  'deviceInstanceLow'                                       ]
        ]
        ['DEVICE_PROPERTIES_OPTION','8' PnDcp_Block_DevicePropertiesOemDeviceId
            // TODO: Implement this ...
        ]
        ['DEVICE_PROPERTIES_OPTION','9' PnDcp_Block_DevicePropertiesStandardGateway
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DHCP_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98 & 100)

        // TODO: Check if these are really all DCP_OPTION
        ['DCP_OPTION','12' PnDcp_Block_DhcpOptionHostName
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','43' PnDcp_Block_DhcpOptionVendorSpecificInformation
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','54' PnDcp_Block_DhcpOptionServerIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','55' PnDcp_Block_DhcpOptionParameterRequestList
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','60' PnDcp_Block_DhcpOptionClassIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','61' PnDcp_Block_DhcpOptionDhcpClientIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','81' PnDcp_Block_DhcpOptionFullyQualifiedDomainName
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','97' PnDcp_Block_DhcpOptionUuidBasedClient
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // CONTROL_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98)

        ['CONTROL_OPTION','1' PnDcp_Block_ControlOptionStart
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','2' PnDcp_Block_ControlOptionStop
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','3' PnDcp_Block_ControlOptionSignal
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','4' PnDcp_Block_ControlOptionResponse
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','5' PnDcp_Block_ControlOptionFactoryReset
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','6' PnDcp_Block_ControlOptionResetToFactory
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DEVICE_INITIATIVE_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98)

        ['DEVICE_INITIATIVE_OPTION','1' PnDcp_Block_DeviceInitiativeOption
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // ALL_SELECTOR_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 99)

        ['ALL_SELECTOR_OPTION','0xFF' PnDcp_Block_ALLSelector
            // This type of block is empty
        ]

        ////////////////////////////////////////////////////////////////////////////
        // Device manufacturer specific options 0x00-0xFF
        ////////////////////////////////////////////////////////////////////////////
   ]
]

[type 'PnDcp_SupportedDeviceOption'
    [simple   PnDcp_BlockOptions 'option']
    [simple   uint 8       'suboption'   ]
]

// 4.3.1.3.2 (Page 94ff)
[type 'PnDcp_ServiceType'
    [reserved uint 5 '0x00'     ]
    [simple   bit 'notSupported']
    [reserved uint 1 '0x00'     ]
    [simple   bit 'response'    ]
]

// Page 86ff: Coding of the field FrameID
[enum uint 16 'PnDcp_FrameId'
    // Range 1
    ['0x0020' PTCP_RTSyncPDUWithFollowUp     ]
    // Range 2
    ['0x0080' PTCP_RTSyncPDU                 ]
    // Range 3
    // 0x100-0x0FFF RT_CLASS_3
    // Range 6
    // 0x8000-BFFF RT_CLASS_1
    // Range 7
    // 0XC000-FBFF RT_CLASS_UDP
    // Range 8
    ['0xFC01' Alarm_High                     ]
    ['0xFE01' Alarm_Low                      ]
    ['0xFEFC' DCP_Hello_ReqPDU               ]
    ['0xFEFD' DCP_GetSet_PDU                 ]
    ['0xFEFE' DCP_Identify_ReqPDU            ]
    ['0xFEFF' DCP_Identify_ResPDU            ]
    // Range 9
    ['0xFF00' PTCP_AnnouncePDU               ]
    ['0xFF20' PTCP_FollowUpPDU               ]
    ['0xFF40' PTCP_DelayReqPDU               ]
    ['0xFF41' PTCP_DelayResPDUWithFollowUp   ]
    ['0xFF42' PTCP_DelayFuResPDUWithFollowUp ]
    ['0xFF43' PTCP_DelayResPDUWithoutFollowUp]
    // Range 12
    // 0xFF80 - 0xFF8F FragmentationFrameId
]

// Page 94
// All other values are "Reserved"
[enum uint 8 'PnDcp_ServiceId'
    ['0x03' GET     ]
    ['0x04' SET     ]
    ['0x05' IDENTIFY]
    ['0x06' HELLO   ]
    //[RESERVED]
]

// 4.3.1.4.1 (Page 97)
// All other values are "Reserved"
[enum uint 8 'PnDcp_BlockOptions'
    ['0x01' IP_OPTION               ]
    ['0x02' DEVICE_PROPERTIES_OPTION]
    ['0x03' DCP_OPTION              ]
    ['0x05' CONTROL_OPTION          ]
    ['0x06' DEVICE_INITIATIVE_OPTION]
    ['0xFF' ALL_SELECTOR_OPTION     ]
]

/////////////////////////////////////////////////////////////////////////////////////////
//
//   PROFINET IO
//
// CM: Context Manager
//
/////////////////////////////////////////////////////////////////////////////////////////

// TODO: Check if it's really Little Endian
// 5.1.2
// 5.5.2.2
[discriminatedType 'PnIoCm_Packet' [DceRpc_PacketType 'packetType']
    [typeSwitch 'packetType'
        ['REQUEST' PnIoCm_Packet_Req
            [simple uint 32            'argsMaximum'            ]
        ]
        ['RESPONSE' PnIoCm_Packet_Res
            [simple uint 8  'errorCode2'             ]
            [simple uint 8  'errorCode1'             ]
            [simple uint 8  'errorDecode'            ]
            [simple uint 8  'errorCode'              ]
        ]
    ]
    [simple uint 32      'argsLength'                           ]
    [simple uint 32      'arrayMaximumCount'                    ]
    [simple uint 32      'arrayOffset'                          ]
    [simple uint 32      'arrayActualCount'                     ]
    [array  PnIoCm_Block 'blocks'            length 'argsLength']
]

// Big Endian
[type 'PnIoCm_Block'
    [discriminator PnIoCm_BlockType 'blockType'                           ]
    [implicit      uint 16          'blockLength'      'lengthInBytes - 4']
    [simple        uint 8           'blockVersionHigh'                    ]
    [simple        uint 8           'blockVersionLow'                     ]
    [typeSwitch 'blockType'
        ['AR_BLOCK_REQ' PnIoCm_Block_ArReq
            [simple   PnIoCm_ArType          'arType'                                                 ]
            [simple   Uuid                   'arUuid'                                                 ]
            [simple   uint 16                'sessionKey'                                             ]
            [simple   MacAddress             'cmInitiatorMacAddr'                                     ]
            [simple   Uuid                   'cmInitiatorObjectUuid'                                  ]
            // Begin ARProperties
            [simple   bit                    'pullModuleAlarmAllowed'                                 ]
            [simple   bit                    'nonLegacyStartupMode'                                   ]
            [simple   bit                    'combinedObjectContainerUsed'                            ]
            [reserved uint 17                '0x00000'                                                ]
            [simple   bit                    'acknowledgeCompanionAr'                                 ]
            [simple   PnIoCm_CompanionArType 'companionArType'                                        ]
            [simple   bit                    'deviceAccess'                                           ]
            [reserved uint 3                 '0x0'                                                    ]
            [simple   bit                    'cmInitiator'                                            ]
            [simple   bit                    'supervisorTakeoverAllowed'                              ]
            [simple   PnIoCm_State           'state'                                                  ]
            // End ARProperties
            [simple   uint 16                'cmInitiatorActivityTimeoutFactor'                       ]
            [simple   uint 16                'cmInitiatorUdpRtPort'                                   ]
            //[implicit uint 16                'stationNameLength'    'STR_LEN(cmInitiatorStationName)' ]
            //[simple   string                 'stationNameLength * 8' 'cmInitiatorStationName'         ]
        ]
        ['AR_BLOCK_RES' PnIoCm_Block_ArRes
            [simple   PnIoCm_ArType          'arType'                                                 ]
            [simple   Uuid                   'arUuid'                                                 ]
            [simple   uint 16                'sessionKey'                                             ]
            [simple   MacAddress             'cmResponderMacAddr'                                     ]
            [simple   Uuid                   'cmResponderObjectUuid'                                  ]
        ]
        ['IO_CR_BLOCK_REQ' PnIoCm_Block_IoCrReq
            [simple PnIoCm_IoCrType          'ioCrType'                                               ]
            [simple uint 16                  'ioCrReference'                                          ]
            [simple uint 16                  'lt'                                                     ]
            // Begin IOCRProperties
            [simple   bit                    'fullSubFrameStructure'                                  ]
            [simple   bit                    'distributedSubFrameWatchDog'                            ]
            [simple   bit                    'fastForwardingMacAdr'                                   ]
            [reserved uint 16                '0x0000'                                                 ]
            [simple   bit                    'mediaRedundancy'                                        ]
            [reserved uint 7                 '0x00'                                                   ]
            [simple   PnIoCm_RtClass         'rtClass'                                                ]
            // End IOCRProperties
            [simple   uint 16                'dataLength'                                             ]
            [simple   uint 16                'frameId'                                                ]
            [simple   uint 16                'sendClockFactor'                                        ]
            [simple   uint 16                'reductionRatio'                                         ]
            [simple   uint 16                'phase'                                                  ]
            [simple   uint 16                'sequence'                                               ]
            [simple   uint 32                'frameSendOffset'                                        ]
            [simple   uint 16                'watchDogFactor'                                         ]
            [simple   uint 16                'dataHoldFactor'                                         ]
            [simple   uint 16                'ioCrTagHeader'                                          ]
            [simple   MacAddress             'ioCrMulticastMacAdr'                                    ]
            [implicit uint 16                'numberOfApis'        'COUNT(apis)'                      ]
            [array    PnIoCm_IoCrBlockReqApi 'apis'                count         'numberOfApis'       ]
        ]
        ['IO_CR_BLOCK_RES' PnIoCm_Block_IoCrRes
            [simple PnIoCm_IoCrType          'ioCrType'                                               ]
            [simple uint 16                  'ioCrReference'                                          ]
            [simple   uint 16                'frameId'                                                ]
        ]
        ['ALARM_CR_BLOCK_REQ' PnIoCm_Block_AlarmCrReq
            [simple   PnIoCm_AlarmCrType     'alarmType'                                              ]
            [simple   uint 16                'lt'                                                     ]
            // Begin AlarmCrProperties
            [reserved uint 30                '0x00000000'                                             ]
            [simple   bit                    'transport'                                              ]
            [simple   bit                    'priority'                                               ]
            // End AlarmCrProperties
            [simple   uint 16                'rtaTimeoutFactor'                                       ]
            [simple   uint 16                'rtaRetries'                                             ]
            [simple   uint 16                'localAlarmReference'                                    ]
            [simple   uint 16                'maxAlarmDataLength'                                     ]
            [simple   uint 16                'alarmCtrTagHeaderHigh'                                  ]
            [simple   uint 16                'alarmCtrTagHeaderLow'                                   ]
        ]
        ['ALARM_CR_BLOCK_RES' PnIoCm_Block_AlarmCrRes
            [simple   PnIoCm_AlarmCrType     'alarmType'                                              ]
            [simple   uint 16                'localAlarmReference'                                    ]
            [simple   uint 16                'maxAlarmDataLength'                                     ]
        ]
        ['EXPECTED_SUBMODULE_BLOCK_REQ' PnIoCm_Block_ExpectedSubmoduleReq
            [implicit uint 16                'numberOfApis'         'COUNT(apis)'                     ]
            [array    PnIoCm_ExpectedSubmoduleBlockReqApi 'apis'   count         'numberOfApis'       ]
        ]
        ['MODULE_DIFF_BLOCK' PnIoCm_Block_ModuleDiff
            [implicit uint 16                'numberOfApis'         'COUNT(apis)'                     ]
            [array    PnIoCm_ModuleDiffBlockApi 'apis'              count         'numberOfApis'      ]
        ]
        ['AR_SERVER_BLOCK' PnIoCm_Block_ArServer
            //[implicit uint 16                'stationNameLength'     'STR_LEN(cmInitiatorStationName)']
            //[simple   string                 'stationNameLength * 8' 'cmInitiatorStationName'         ]
            //[padding  byte '0x00'                                                                     ]
        ]
    ]
]

[type 'PnIoCm_IoCrBlockReqApi'
    [const    uint 32            'api'              '0x00000000'             ]
    [implicit uint 16            'numIoDataObjects' 'COUNT(ioDataObjects)'   ]
    [array    PnIoCm_IoDataObject 'ioDataObjects'    count 'numIoDataObjects']
    [implicit uint 16            'numIoCss'         'COUNT(ioCss)'           ]
    [array    PnIoCm_IoCs         'ioCss'            count 'numIoCss'        ]
]

[type 'PnIoCm_IoDataObject'
    [simple   uint 16 'slotNumber'             ]
    [simple   uint 16 'subSlotNumber'          ]
    [simple   uint 16 'ioDataObjectFrameOffset']
]

[type 'PnIoCm_IoCs'
    [simple   uint 16 'slotNumber'   ]
    [simple   uint 16 'subSlotNumber']
    [simple   uint 16 'ioFrameOffset']
]

[type 'PnIoCm_ExpectedSubmoduleBlockReqApi'
    [const    uint 32          'api'               '0x00000000'                       ]
    [simple   uint 16          'slotNumber'                                           ]
    [simple   uint 32          'moduleIdentNumber'                                    ]
    [simple   uint 16          'moduleProperties'                                     ]
    [implicit uint 16          'numSubmodules'     'COUNT(submodules)'                ]
    [array    PnIoCm_Submodule 'submodules'        count               'numSubmodules']
]

[type 'PnIoCm_ModuleDiffBlockApi'
    [const    uint 32                          'api'        '0x00000000'                    ]
    [implicit uint 16                          'numModules' 'COUNT(modules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Module 'modules'    count               'numModules']
]

[type 'PnIoCm_ModuleDiffBlockApi_Module'
    [simple   uint 16                             'slotNumber'                                           ]
    [simple   uint 32                             'moduleIdentNumber'                                    ]
    [simple   PnIoCm_ModuleState                  'moduleState'                                          ]
    [implicit uint 16                             'numSubmodules'     'COUNT(submodules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Submodule 'submodules'        count               'numSubmodules']
]

[type 'PnIoCm_ModuleDiffBlockApi_Submodule'
    [simple uint 16          'subslotNumber'       ]
    [simple uint 32          'submoduleIdentNumber']
    [simple bit              'codingUsesIdentInfo' ]
    [simple PnIoCm_IdentInfo 'identInfo'           ]
    [simple PnIoCm_ArInfo    'arInfo'              ]
    [simple bit              'diagInfoAvailable'   ]
    [simple bit              'maintenanceDemanded' ]
    [simple bit              'maintenanceRequired' ]
    [simple bit              'qualifiedInfo'       ]
    [simple PnIoCm_AddInfo   'addInfo'             ]
]

[type 'PnIoCm_Submodule'
    [simple   uint 16                'slotNumber'                  ]
    [simple   uint 32                'submoduleIdentNumber'        ]
    // Begin SubmoduleProperties
    [reserved uint 10 '0x000']
    [simple   bit                    'discardIoxs'                 ]
    [simple   bit                    'reduceOutputModuleDataLength']
    [simple   bit                    'reduceInputModuleDataLength' ]
    [simple   bit                    'sharedInput'                 ]
    [simple   PnIoCm_SubmoduleType   'submoduleType'               ]
    // End SubmoduleProperties
    [simple   PnIoCm_DescriptionType 'descriptionType'             ]
    [simple   uint 16                'submoduleDataLength'         ]
    [simple   uint 8                 'lengthIoCs'                  ]
    [simple   uint 8                 'lengthIoPs'                  ]
]

[enum uint 16 'PnIoCm_BlockType'
    ['0x0101' AR_BLOCK_REQ                ]
    ['0x8101' AR_BLOCK_RES                ]
    ['0x0102' IO_CR_BLOCK_REQ             ]
    ['0x8102' IO_CR_BLOCK_RES             ]
    ['0x0103' ALARM_CR_BLOCK_REQ          ]
    ['0x8103' ALARM_CR_BLOCK_RES          ]
    ['0x0104' EXPECTED_SUBMODULE_BLOCK_REQ]
    ['0x8104' MODULE_DIFF_BLOCK           ]
    ['0x8106' AR_SERVER_BLOCK             ]
]

[enum uint 16 'PnIoCm_ArType'
    ['0x0001' IO_CONTROLLER]
]

[enum uint 2 'PnIoCm_CompanionArType'
    ['0x0' SINGLE_AR]
]

[enum uint 3 'PnIoCm_State'
    ['0x1' ACTIVE]
]

[enum uint 16 'PnIoCm_IoCrType'
    ['0x0001' INPUT_CR]
]

[enum uint 4 'PnIoCm_RtClass'
    ['0x0010' RT_CLASS_2]
]

[enum uint 16 'PnIoCm_AlarmCrType'
    ['0x0001' ALARM_CR]
]

[enum uint 16 'PnIoCm_ModuleState'
    ['0x0002' PROPER_MODULE]
]

[enum uint 2 'PnIoCm_SubmoduleType'
    ['0x0' NO_INPUT_NO_OUTPUT_DATA]
]

[enum uint 16 'PnIoCm_DescriptionType'
    ['0x0001' INPUT]
]

[enum uint 4 'PnIoCm_IdentInfo'
    ['0x0' OK]
]

[enum uint 4 'PnIoCm_ArInfo'
    ['0x0' OWN]
]

[enum uint 3 'PnIoCm_AddInfo'
    ['0x0' NONE]
]