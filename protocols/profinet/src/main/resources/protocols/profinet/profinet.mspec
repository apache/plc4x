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
            [simple ProfinetFrameId 'frameId']
            [simple DCP_PDU         'pdu']
        ]
    ]
]

// Page 90
[discriminatedType 'DCP_PDU'
    [discriminator ServiceId   'serviceId']
    [discriminator ServiceType 'serviceType']
    // 4.3.1.3.4 (Page 95)
    [simple        uint 32     'xid']
    // 4.3.1.3.5 (Page 95ff)
    [simple        uint 16     'responseDelayFactorOrPadding']
    // 4.3.1.3.4 (Page 95)
    [implicit      uint 16     'dcpDataLength' 'lengthInBytes - 10']
    [typeSwitch 'serviceId','serviceType.response'
        // Multicast
        ['ServiceId.IDENTIFY','false' DCP_Identify_ReqPDU
            // For a DCP-IdentifyFilter-ReqPDU this can contain an optional NameOfStationBlock or AliasNameBlock and an optional IdentifyReqBlock (in total max 2)
            // For a DCP-IdentifyAll_ReqPDU this must contain an AllSelectorBlock (in total 1)
            [array DCP_Block 'blocks' length 'dcpDataLength']
        ]
        ['ServiceId.HELLO','false' DCP_Hello_ReqPDU
            [simple NameOfStationBlockRes    'nameOfStationBlockRes']
            [simple IPParameterBlockRes      'iPParameterBlockRes']
            [simple DeviceIdBlockRes         'deviceIdBlockRes']
            [simple DeviceVendorBlockRes     'deviceVendorBlockRes']
            [simple DeviceOptionsBlockRes    'deviceOptionsBlockRes']
            [simple DeviceRoleBlockRes       'deviceRoleBlockRes']
            [simple DeviceInitiativeBlockRes 'deviceInitiativeBlockRes']
        ]
        // Unicast
        /*['','' DCP_Get_ReqPDU]
            [simple GetReqBlock              'getReqBlock']
        ['','' DCP_Set_ReqPDU]
            [simple StartTransactionBlock    'startTransactionBlock']
                [simple BlockQualifier           'blockQualifier']
            [simple SetResetReqBlock         'setResetReqBlock']
            [simple SetReqBlock              'setReqBlock']
            [simple StopTransactionBlock     'stopTransactionBlock']
                [simple BlockQualifier           'blockQualifier']
        ['','' DCP_Get_ResPDU]
            [simple GetResBlock              'getResBlock']
            [simple GetNegResBlock           'getNegResBlock']
        ['','' DCP_Set_ResPDU]
            [simple SetResBlock              'setResBlock']
            [simple SetNegResBlock           'setNegResBlock']
        ['','' DCP_Identify_ResPDU
            [simple IdentifyResBlock         'identifyRes']
            [simple NameOfStationBlockRes    'nameOfStationBlockRes']
            [simple IPParameterBlockRes      'iPParameterBlockRes']
            [simple DeviceIdBlockRes         'deviceIdBlockRes']
            [simple DeviceVendorBlockRes     'deviceVendorBlockRes']
            [simple DeviceOptionsBlockRes    'deviceOptionsBlockRes']
            [simple DeviceRoleBlockRes       'deviceRoleBlockRes']
            [simple DeviceInitiativeBlockRes 'deviceInitiativeBlockRes']
            [simple DeviceInstanceBlockRes   'deviceInstanceBlockRes']
            [simple OemDeviceIdBlockRes      'oemDeviceIdBlockRes']
        ]*/
    ]
]

[discriminatedType 'DCP_Block'
    [discriminator BlockOptions 'option']
    [discriminator uint 8       'suboption']
    [implicit      uint 16      'blockLength' 'lengthInBytes']
    [typeSwitch 'option','suboption'

        // 4.3.1.4.1 (Page 97)
        ['BlockOptions.IP_OPTION','1' DCP_BlockIpMacAddress
        ]
        ['BlockOptions.IP_OPTION','2' DCP_BlockIpIpParameter
            // 4.3.1.4.12 (Page 105ff)
            [reserved uint 8 '0x00']
            [simple   bit    'ipConflictDetected']
            [reserved uint 5 '0x00']
            [simple   bit    'setViaDhcp']
            [simple   bit    'setManually']
            [array    uint 8 'ipAddress'       count '4']
            [array    uint 8 'subnetMask'      count '4']
            [array    uint 8 'standardGateway' count '4']
        ]
        ['BlockOptions.IP_OPTION','3' DCP_BlockIpFullIpSuite
        ]

        ['BlockOptions.ALL_SELECTOR_OPTION','0xFF' DCP_BlockALLSelector
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','1' DCP_BlockDevicePropertiesDeviceVendor
            [reserved uint 16 '0x0000']
            // TODO: Put a correct number here
            [simple   string 'length - 42' 'deviceVendorValue']
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','2' DCP_BlockDevicePropertiesNameOfStation
            [reserved uint 16 '0x0000']
            // TODO: Put a correct number here
            [simple   string 'length - 42' 'nameOfStation']
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','3' DCP_BlockDevicePropertiesDeviceId
            [reserved uint 16 '0x0000']
            [simple   uint 16 'vendorId']
            [simple   uint 16 'deviceId']
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','4' DCP_BlockDevicePropertiesDeviceRole
            [reserved uint 16 '0x0000']
            [simple   uint 8  'deviceRoleDetails']
            [reserved uint 8  '0x00']
        ]
        // TODO: Investigate why this has an option and suboption inside again ...
        // TODO: Seems to be an array of tuples of option+suboptions
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','5' DCP_BlockDevicePropertiesDeviceOptions
            [reserved uint 16      '0x0000']
            [simple   BlockOptions 'option']
            [simple   uint 8       'suboption']
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','6' DCP_BlockDevicePropertiesAliasName
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','7' DCP_BlockDevicePropertiesStandardGateway
            [reserved uint 16 '0x0000']
            [simple   uint 8  'deviceInstanceHigh']
            [simple   uint 8  'deviceInstanceLow']

        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','8' DCP_BlockDevicePropertiesOemDeviceId
        ]
        ['BlockOptions.DEVICE_PROPERTIES_OPTION','9' DCP_BlockDevicePropertiesOemDeviceId
        ]
    ]
]

// 4.3.1.3.2 (Page 94ff)
[type 'ServiceType'
    [simple   bit 'response']
    [reserved bit '0x00']
    [simple   bit 'notSupported']
    [reserved uint 5 '0x00']
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
    ['0x06' DEVICE_INITIALIVE_OPTION]
    ['0xFF' ALL_SELECTOR_OPTION]
    //[RESERVED]
]


// 4.3.1.4.1 (Page 98 & 100)
// All other values are "Reserved"
[enum uint 8 'BlockOptionsDhcpSuboptions' [bit 'read', bit 'write', bit 'opt']
    ['12' HOST_NAME                   ['true', 'true', 'true']]
    ['43' VENDOR_SPECIFIC_INFORMATION ['true', 'true', 'true']]
    ['54' SERVER_IDENTIFIER           ['true', 'true', 'true']]
    ['55' PARAMETER_REQUEST_LIST      ['true', 'true', 'true']]
    ['60' CLASS_IDENTIFIER            ['true', 'true', 'true']]
    ['61' DHCP_CLIENT_IDENTIFIER      ['true', 'true', 'true']]
    ['81' FULLY_QUALIFIED_DOMAIN_NAME ['true', 'true', 'true']]
    ['97' UUID_BASED_CLIENT           ['true', 'true', 'true']]
    //[RESERVED]
]

// 4.3.1.4.1 (Page 98)
// All other values are "Reserved"
[enum uint 8 'BlockOptionsControlSuboptions' [bit 'read', bit 'write', bit 'opt']
    ['0x01' START            ['false', 'true',  'false']]
    ['0x02' STOP             ['false', 'true',  'false']]
    ['0x03' SIGNAL           ['false', 'true',  'false']]
    ['0x04' RESPONSE         ['false', 'false', 'false']]
    ['0x05' FACTORY_RESET    ['false', 'true',  'true' ]]
    ['0x06' RESET_TO_FACTORY ['false', 'true',  'false']]
    //[RESERVED]
]

// 4.3.1.4.1 (Page 98)
// All other values are "Reserved"
[enum uint 8 'BlockOptionsDeviceInitiativeSuboptions' [bit 'read', bit 'write', bit 'opt']
    ['0x01' DEVICE_INITIATIVE ['true', 'false',  'false']]
    //[RESERVED]
]

// 4.3.1.4.1 (Page 99)
// All other values are "Reserved"
[enum uint 8 'BlockOptionsDeviceInitiativeSuboptions' [bit 'read', bit 'write', bit 'opt']
    ['0xFF' ALL_SELECTOR ['false', 'false',  'false']]
    //[RESERVED]
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