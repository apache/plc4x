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

[type 'EthernetFrame'
    // When sending to the mac address prefix of 01:0e:cf are multicast packets
    [simple MacAddress          'destination']
    [simple MacAddress          'source']
    [simple EthernetFamePayload 'payload']
]

[discriminatedType 'EthernetFamePayload'
    [discriminator uint 16 'packetType']
    [typeSwitch 'packetType'
        ['0x8100' VirtualLanEthernetFramePayload
            [simple VirtualLanPriority  'priority']
            [simple bit                 'ineligible']
            [simple uint 12             'id']
            [simple EthernetFamePayload 'payload']
        ]
        ['0x8892' ProfinetEthernetFramePayload
            [simple DCP_PDU         'pdu']
        ]
    ]
]

// Page 90
[discriminatedType 'DCP_PDU'
    [discriminator ProfinetFrameId  'frameId']
    [discriminator ServiceId        'serviceId']
    [simple        ServiceType      'serviceType']
    // 4.3.1.3.4 (Page 95)
    [simple        uint 32          'xid']
    // 4.3.1.3.5 (Page 95ff)
    [simple        uint 16          'responseDelayFactorOrPadding']
    // 4.3.1.3.4 (Page 95)
    [implicit      uint 16          'dcpDataLength' 'lengthInBytes - 12']
    [typeSwitch 'frameId','serviceId','serviceType.response'
        ////////////////////////////////////////////////////////////////////////////
        // Multicast (Well theoretically)
        ////////////////////////////////////////////////////////////////////////////
        // The Identify request is valid in two options:
        // 1) One containing only an AllSelectorBlock
        // 2) One containing optionally either NameOfStationBlock or AliasNameBlock and another optional IdentifyReqBlock
        // (I assume, that if in case 2 both optionally aren't used, this might not be valid and option 1 should be sent instead)
        ['DCP_Identify_ReqPDU','IDENTIFY','false' DCP_Identify_ReqPDU [uint 16 'dcpDataLength']
            [array DCP_Block 'blocks' length 'dcpDataLength']
        ]

        // Response to a Identify request
        ['DCP_Identify_ResPDU','IDENTIFY','true' DCP_Identify_ResPDU [uint 16 'dcpDataLength']
            [array DCP_Block 'blocks' length 'dcpDataLength']
        ]

        // Packet a Profinet station might emit once it is turned on
        ['DCP_Hello_ReqPDU','HELLO','false' DCP_Hello_ReqPDU
//            [simple NameOfStationBlockRes    'nameOfStationBlockRes']
//            [simple IPParameterBlockRes      'iPParameterBlockRes']
//            [simple DeviceIdBlockRes         'deviceIdBlockRes']
//            [simple DeviceVendorBlockRes     'deviceVendorBlockRes']
//            [simple DeviceOptionsBlockRes    'deviceOptionsBlockRes']
//            [simple DeviceRoleBlockRes       'deviceRoleBlockRes']
//            [simple DeviceInitiativeBlockRes 'deviceInitiativeBlockRes']
        ]

        ////////////////////////////////////////////////////////////////////////////
        // Unicast
        ////////////////////////////////////////////////////////////////////////////

        ['DCP_GetSet_PDU','GET','false' DCP_Get_ReqPDU
//            [simple GetReqBlock              'getReqBlock']
        ]
        ['DCP_GetSet_PDU','GET','true' DCP_Get_ResPDU
//            [simple GetResBlock              'getResBlock']
//            [simple GetNegResBlock           'getNegResBlock']
        ]

        ['DCP_GetSet_PDU','SET','false' DCP_Set_ReqPDU
//            [simple StartTransactionBlock    'startTransactionBlock']
//            [simple BlockQualifier           'blockQualifier']
//            [simple SetResetReqBlock         'setResetReqBlock']
//            [simple SetReqBlock              'setReqBlock']
//            [simple StopTransactionBlock     'stopTransactionBlock']
//            [simple BlockQualifier           'blockQualifier']
        ]
        ['DCP_GetSet_PDU','SET','true' DCP_Set_ResPDU
//            [simple SetResBlock              'setResBlock']
//            [simple SetNegResBlock           'setNegResBlock']
        ]
    ]
]

[discriminatedType 'DCP_Block'
    [discriminator BlockOptions 'option'                                                ]
    [discriminator uint 8       'suboption'                                             ]
    [implicit      uint 16      'blockLength' 'lengthInBytes - 4'                       ]
    [typeSwitch 'option','suboption'

        ////////////////////////////////////////////////////////////////////////////
        // IP_OPTION
        ////////////////////////////////////////////////////////////////////////////

        // 4.3.1.4.1 (Page 97)
        ['IP_OPTION','1' DCP_BlockIpMacAddress
            [reserved uint 16  '0x0000'                                                 ]
            [simple MacAddress 'macAddress'                                             ]
        ]
        ['IP_OPTION','2' DCP_BlockIpIpParameter
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
        ['IP_OPTION','3' DCP_BlockIpFullIpSuite
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DEVICE_PROPERTIES_OPTION
        ////////////////////////////////////////////////////////////////////////////

        ['DEVICE_PROPERTIES_OPTION','1' DCP_BlockDevicePropertiesDeviceVendor [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            // TODO: Figure out how to do this correctly.
            [array    byte        'deviceVendorValue' count 'blockLength-2'             ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", deviceVendorValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','2' DCP_BlockDevicePropertiesNameOfStation [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            // TODO: Figure out how to do this correctly.
            [array    byte        'nameOfStation' count 'blockLength-2'                 ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", nameOfStation) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','3' DCP_BlockDevicePropertiesDeviceId
            [reserved uint 16 '0x0000'                                                  ]
            [simple   uint 16 'vendorId'                                                ]
            [simple   uint 16 'deviceId'                                                ]
        ]
        ['DEVICE_PROPERTIES_OPTION','4' DCP_BlockDevicePropertiesDeviceRole
            [reserved uint 20 '0x000000'                                                ]
            [simple   bit     'pnioSupervisor'                                          ]
            [simple   bit     'pnioMultidevive'                                         ]
            [simple   bit     'pnioController'                                          ]
            [simple   bit     'pnioDevice'                                              ]
            [reserved uint 8  '0x00'                                                    ]
        ]
        // Contains a list of option combinations the device supports.
        ['DEVICE_PROPERTIES_OPTION','5' DCP_BlockDevicePropertiesDeviceOptions [uint 16 'blockLength']
            [reserved uint 16               '0x0000'                                    ]
            [array    SupportedDeviceOption 'supportedOptions' length 'blockLength - 2' ]
        ]
        ['DEVICE_PROPERTIES_OPTION','6' DCP_BlockDevicePropertiesAliasName [uint 16 'blockLength']
            [reserved uint 16     '0x0000'                                              ]
            [array    byte        'aliasNameValue' count 'blockLength-2'                ]
            [padding  uint 8      'pad' '0x00' 'STATIC_CALL("org.apache.plc4x.java.profinet.utils.StaticHelper.arrayLength", aliasNameValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','7' DCP_BlockDevicePropertiesDeviceInstance
            [reserved uint 16 '0x0000'                                                  ]
            [simple   uint 8  'deviceInstanceHigh'                                      ]
            [simple   uint 8  'deviceInstanceLow'                                       ]
        ]
        ['DEVICE_PROPERTIES_OPTION','8' DCP_BlockDevicePropertiesOemDeviceId
            // TODO: Implement this ...
        ]
        ['DEVICE_PROPERTIES_OPTION','9' DCP_BlockDevicePropertiesStandardGateway
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DCP_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98 & 100)

        // TODO: Check if these are really all DCP_OPTION
        ['DCP_OPTION','12' DCP_BlockDhcpOptionHostName
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','43' DCP_BlockDhcpOptionVendorSpecificInformation
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','54' DCP_BlockDhcpOptionServerIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','55' DCP_BlockDhcpOptionParameterRequestList
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','60' DCP_BlockDhcpOptionClassIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','61' DCP_BlockDhcpOptionDhcpClientIdentifier
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','81' DCP_BlockDhcpOptionFullyQualifiedDomainName
            // TODO: Implement this ...
        ]
        ['DCP_OPTION','97' DCP_BlockDhcpOptionUuidBasedClient
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // CONTROL_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98)

        ['CONTROL_OPTION','1' DCP_BlockControlOptionStart
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','2' DCP_BlockControlOptionStop
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','3' DCP_BlockControlOptionSignal
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','4' DCP_BlockControlOptionResponse
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','5' DCP_BlockControlOptionFactoryReset
            // TODO: Implement this ...
        ]
        ['CONTROL_OPTION','6' DCP_BlockControlOptionResetToFactory
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DEVICE_INITIATIVE_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 98)

        ['DEVICE_INITIATIVE_OPTION','1' DCP_BlockDeviceInitiativeOption
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // ALL_SELECTOR_OPTION
        ////////////////////////////////////////////////////////////////////////////
        // 4.3.1.4.1 (Page 99)

        ['ALL_SELECTOR_OPTION','0xFF' DCP_BlockALLSelector
            // This type of block is empty
        ]

        ////////////////////////////////////////////////////////////////////////////
        // Device manufacturer specific options 0x00-0xFF
        ////////////////////////////////////////////////////////////////////////////
   ]
]

[type 'SupportedDeviceOption'
    [simple   BlockOptions 'option'                                             ]
    [simple   uint 8       'suboption'                                          ]
]

// 4.3.1.3.2 (Page 94ff)
[type 'ServiceType'
    [reserved uint 5 '0x00']
    [simple   bit 'notSupported']
    [reserved uint 1 '0x00']
    [simple   bit 'response']
]

// Page 86ff: Coding of the field FrameID
[enum uint 16 'ProfinetFrameId'
    // Range 1
    ['0x0020' PTCP_RTSyncPDUWithFollowUp]
    // Range 2
    ['0x0080' PTCP_RTSyncPDU]
    // Range 3
    // 0x100-0x0FFF RT_CLASS_3
    // Range 6
    // 0x8000-BFFF RT_CLASS_1
    // Range 7
    // 0XC000-FBFF RT_CLASS_UDP
    // Range 8
    ['0xFC01' Alarm_High]
    ['0xFE01' Alarm_Low]
    ['0xFEFC' DCP_Hello_ReqPDU]
    ['0xFEFD' DCP_GetSet_PDU]
    ['0xFEFE' DCP_Identify_ReqPDU]
    ['0xFEFF' DCP_Identify_ResPDU]
    // Range 9
    ['0xFF00' PTCP_AnnouncePDU]
    ['0xFF20' PTCP_FollowUpPDU]
    ['0xFF40' PTCP_DelayReqPDU]
    ['0xFF41' PTCP_DelayResPDUWithFollowUp]
    ['0xFF42' PTCP_DelayFuResPDUWithFollowUp]
    ['0xFF43' PTCP_DelayResPDUWithoutFollowUp]
    // Range 12
    // 0xFF80 - 0xFF8F FragmentationFrameId
]

// Page 94
// All other values are "Reserved"
[enum uint 8 'ServiceId'
    ['0x03' GET]
    ['0x04' SET]
    ['0x05' IDENTIFY]
    ['0x06' HELLO]
    //[RESERVED]
]

// 4.3.1.4.1 (Page 97)
// All other values are "Reserved"
[enum uint 8 'BlockOptions'
    ['0x01' IP_OPTION]
    ['0x02' DEVICE_PROPERTIES_OPTION]
    ['0x03' DCP_OPTION]
    ['0x05' CONTROL_OPTION]
    ['0x06' DEVICE_INITIATIVE_OPTION]
    ['0xFF' ALL_SELECTOR_OPTION]
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