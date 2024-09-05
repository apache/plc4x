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

[type Ethernet_Frame byteOrder='BIG_ENDIAN'
    // When sending to the mac address prefix of 01:0e:cf are multicast packets
    [simple MacAddress            destination]
    [simple MacAddress            source     ]
    [simple Ethernet_FramePayload payload    ]
]

[discriminatedType Ethernet_FramePayload
    [discriminator uint 16 packetType]
    [typeSwitch packetType
        ['0x8100' Ethernet_FramePayload_VirtualLan
            [simple VirtualLanPriority    priority                                                   ]
            [simple bit                   ineligible                                                 ]
            [simple uint 12               id                                                         ]
            [simple Ethernet_FramePayload payload                                                    ]
        ]
        ['0x8892' Ethernet_FramePayload_PnDcp
            [simple PnDcp_Pdu             pdu                                                        ]
        ]
    ]
]

// There are some special MAC addresses reserved:
// 01-0E-CF-00-00-00:      As destination for DCP-Identify-ReqPDU (The FrameID is then required to be set to 0xFEFE)
// 01-0E-CF-00-00-01:      As destination for DCP-Helo-ReqPDU (The FrameID is then required to be set to 0xFEFC)
// 01-0E-CF-00-00-02:      Reserved
// 01-0E-CF-00-01-00:      Reserved for further multicast addresses within the Type 10 context
// 01-0E-CF-00-01-01:      As multicast destination for RT_CLASS_3
// 01-0E-CF-00-01-02:      As invalid frame multicast destination for RT_CLASS_3
// 01-0E-CF-00-01-03 - FF: Reserved for further multicast addresses within the Type 10 context
[type MacAddress
    [array byte address count '6']
]

[type IpAddress
    [array byte data count '4']
]

// https://de.wikipedia.org/wiki/IEEE_802.1p
[enum uint 3 VirtualLanPriority   (string 16 acronym)
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
[discriminatedType PnDcp_Pdu byteOrder='BIG_ENDIAN'
    [simple        uint 16           frameIdValue                      ]
    [typeSwitch frameIdValue
        ['0xFEFE' PnDcp_Pdu_IdentifyReq
            [const    uint 8      serviceId                    0x05                                ]
            // ServiceType Start
            [reserved uint 5      '0x00'                                                           ]
            [const    bit         notSupported                 false                               ]
            [reserved uint 1      '0x00'                                                           ]
            [const    bit         response                     false                               ]
            // ServiceType End
            // 4.3.1.3.3 (Page 95)
            [simple   uint 32     xid                                                              ]
            // 4.3.1.3.5 (Page 95ff)
            [simple   uint 16     responseDelay                                                    ]
            // 4.3.1.3.4 (Page 95)
            [implicit uint 16     dcpDataLength                'lengthInBytes - 12'                ]
            [array    PnDcp_Block blocks                        length              'dcpDataLength']
        ]
        ////////////////////////////////////////////////////////////////////////////
        // Multicast (Well theoretically)
        ////////////////////////////////////////////////////////////////////////////
        // The Identify request is valid in two options:
        // 1) One containing only an AllSelectorBlock
        // 2) One containing optionally either NameOfStationBlock or AliasNameBlock and another optional IdentifyReqBlock
        // (I assume, that if in case 2 both optionally aren't used, this might not be valid and option 1 should be sent instead)
        ['0xFEFF' PnDcp_Pdu_IdentifyRes
            [const    uint 8      serviceId                    0x05                                ]
            // ServiceType Start
            [reserved uint 5      '0x00'                                                           ]
            [simple   bit         notSupported                                                     ]
            [reserved uint 1      '0x00'                                                           ]
            [const    bit         response                     true                                ]
            // ServiceType End
            // 4.3.1.3.3 (Page 95)
            [simple   uint 32     xid                                                              ]
            // 4.3.1.3.5 (Page 95ff)
            [reserved uint 16     '0x0000'                                                         ]
            // 4.3.1.3.4 (Page 95)
            [implicit uint 16     dcpDataLength                'lengthInBytes - 12'                ]
            [array    PnDcp_Block blocks                        length              'dcpDataLength']
        ]
    ]
]

[discriminatedType PnDcp_Block byteOrder='BIG_ENDIAN'
    [discriminator PnDcp_BlockOptions option                   ]
    [discriminator uint 8       suboption                      ]
    [implicit      uint 16      blockLength 'lengthInBytes - 4']
    [typeSwitch option,suboption

        ////////////////////////////////////////////////////////////////////////////
        // IP_OPTION
        ////////////////////////////////////////////////////////////////////////////

        // 4.3.1.4.1 (Page 97)
        ['IP_OPTION','1' PnDcp_Block_IpMacAddress
            [reserved uint 16  '0x0000'                                               ]
            [simple MacAddress macAddress                                             ]
        ]
        ['IP_OPTION','2' PnDcp_Block_IpParameter
            // 4.3.1.4.12 (Page 105ff)
            [reserved uint 8 '0x00'                                                   ]
            [simple   bit    ipConflictDetected                                       ]
            [reserved uint 5 '0x00'                                                   ]
            [simple   bit    setViaDhcp                                               ]
            [simple   bit    setManually                                              ]
            [array    byte   ipAddress       count '4'                                ]
            [array    byte   subnetMask      count '4'                                ]
            [array    byte   standardGateway count '4'                                ]
        ]
        ['IP_OPTION','3' PnDcp_Block_FullIpSuite
            // TODO: Implement this ...
        ]

        ////////////////////////////////////////////////////////////////////////////
        // DEVICE_PROPERTIES_OPTION
        ////////////////////////////////////////////////////////////////////////////

        ['DEVICE_PROPERTIES_OPTION','1' PnDcp_Block_DevicePropertiesDeviceVendor(uint 16 blockLength)
            [reserved uint 16     '0x0000'                                            ]
            [array    byte        deviceVendorValue count 'blockLength-2'             ]
            [padding  uint 8      pad '0x00' 'STATIC_CALL("arrayLength", deviceVendorValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','2' PnDcp_Block_DevicePropertiesNameOfStation(uint 16 blockLength)
            [reserved uint 16     '0x0000'                                            ]
            [array    byte        nameOfStation count 'blockLength-2'                 ]
            [padding  uint 8      pad '0x00' 'STATIC_CALL("arrayLength", nameOfStation) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','3' PnDcp_Block_DevicePropertiesDeviceId
            [reserved uint 16 '0x0000'                                                ]
            [simple   uint 16 vendorId                                                ]
            [simple   uint 16 deviceId                                                ]
        ]
        ['DEVICE_PROPERTIES_OPTION','4' PnDcp_Block_DevicePropertiesDeviceRole
            [reserved uint 20 '0x000000'                                              ]
            [simple   bit     pnioSupervisor                                          ]
            [simple   bit     pnioMultidevive                                         ]
            [simple   bit     pnioController                                          ]
            [simple   bit     pnioDevice                                              ]
            [reserved uint 8  '0x00'                                                  ]
        ]
        // Contains a list of option combinations the device supports.
        ['DEVICE_PROPERTIES_OPTION','5' PnDcp_Block_DevicePropertiesDeviceOptions(uint 16 blockLength)
            [reserved uint 16               '0x0000'                                  ]
            [array    PnDcp_SupportedDeviceOption supportedOptions length 'blockLength - 2' ]
        ]
        ['DEVICE_PROPERTIES_OPTION','6' PnDcp_Block_DevicePropertiesAliasName(uint 16 blockLength)
            [reserved uint 16     '0x0000'                                            ]
            [array    byte        aliasNameValue count 'blockLength-2'                ]
            [padding  uint 8      pad '0x00' 'STATIC_CALL("arrayLength", aliasNameValue) % 2']
        ]
        ['DEVICE_PROPERTIES_OPTION','7' PnDcp_Block_DevicePropertiesDeviceInstance
            [reserved uint 16 '0x0000'                                                ]
            [simple   uint 8  deviceInstanceHigh                                      ]
            [simple   uint 8  deviceInstanceLow                                       ]
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

[type PnDcp_SupportedDeviceOption byteOrder='BIG_ENDIAN'
    [simple   PnDcp_BlockOptions option]
    [simple   uint 8       suboption   ]
]

// 4.3.1.4.1 (Page 97)
// All other values are "Reserved"
[enum uint 8 PnDcp_BlockOptions
    ['0x01' IP_OPTION               ]
    ['0x02' DEVICE_PROPERTIES_OPTION]
    ['0x03' DCP_OPTION              ]
    ['0x05' CONTROL_OPTION          ]
    ['0x06' DEVICE_INITIATIVE_OPTION]
    ['0xFF' ALL_SELECTOR_OPTION     ]
]

