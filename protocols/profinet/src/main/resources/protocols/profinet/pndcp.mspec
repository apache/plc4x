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
    [virtual       PnDcp_FrameId     frameId       'STATIC_CALL("getFrameId", frameIdValue)']
    [typeSwitch frameId
        ['RT_CLASS_1' PnDcp_Pdu_RealTimeCyclic
            [manual   PnIo_CyclicServiceDataUnit
                                          dataUnit
                                                'STATIC_CALL("readDataUnit", readBuffer)'
                                                'STATIC_CALL("writeDataUnit", writeBuffer, dataUnit)'
                                                '(dataUnit.lengthInBytes)*8'      ]
            [simple   uint 16                    cycleCounter             ]
            // Data Status Start (4.7.2.1.3)
            [simple   bit                        ignore                   ]
            [reserved bit                        'false'                  ]
            [simple   bit                        stationProblemIndicatorOk]
            [simple   bit                        providerStateRun         ]
            [reserved bit                        'false'                  ]
            [simple   bit                        dataValid                ]
            [simple   bit                        redundancy               ]
            [simple   bit                        statePrimary             ]
            // Data Status End
            // "Transfer-Status" (Set to 0x00 for all RT-Classes except RT-CLASS-3,
            // which PLC4X will never be able to support
            [reserved uint 8                     '0x00'                   ] // transferStatus
        ]
        ['PTCP_DelayReqPDU' PcDcp_Pdu_DelayReq
            // Header Start
            [reserved uint 32 '0x00000000']
            [reserved uint 32 '0x00000000']
            [reserved uint 32 '0x00000000']
            [simple   uint 16 sequenceId  ]
            [reserved uint 16 '0x0000'    ]
            // Header End
            [simple   uint 32 delayInNs   ]
            // Delay Parameter Start
            // TODO: This seems to usually be an array of parameters terminated by an End-Parameter which is indicated by type and length being 0
            [const    uint 7     parameterType   6]
            [const    uint 9     parameterLength 6]
            [simple   MacAddress portMacAddress   ]
            [const    uint 7     endType         0]
            [const    uint 9     endLength       0]
            // Delay Parameter End
        ]
        ['Alarm_Low' PnDcp_Pdu_AlarmLow
                    [simple uint 16 alarmDstEndpoint]
                    [simple uint 16 alarmSrcEndpoint]
                    [simple uint 4  version]
                    [simple uint 4  errorType]
                    [simple uint 4  tAck]
                    [simple uint 4  windowSize]
                    [simple uint 16 senSeqNum]
                    [simple uint 16 ackSeqNum]
                    [implicit uint 16 varPartLen 'COUNT(varPart)']
                    [array    byte varPart                        length              'varPartLen']
                ]
        ['DCP_Identify_ReqPDU' PnDcp_Pdu_IdentifyReq
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
        ['DCP_Identify_ResPDU' PnDcp_Pdu_IdentifyRes
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

[discriminatedType PnDcp_Pdu_IdentifyRes_Payload byteOrder='BIG_ENDIAN'
    [discriminator PnDcp_ServiceId   serviceId                         ]
    [simple        PnDcp_ServiceType serviceType                       ]
    // 4.3.1.3.3 (Page 95)
    [simple        uint 32           xid                               ]
    // 4.3.1.3.5 (Page 95ff)
    [simple        uint 16           responseDelayFactorOrPadding      ]
    // 4.3.1.3.4 (Page 95)
    [implicit      uint 16           dcpDataLength 'lengthInBytes - 12']
    [typeSwitch serviceId,serviceType.response

        // Packet a Profinet station might emit once it is turned on
//        ['DCP_Hello_ReqPDU','HELLO','false' PnDcp_Pdu_HelloReq
//            [simple NameOfStationBlockRes    nameOfStationBlockRes   ]
//            [simple IPParameterBlockRes      iPParameterBlockRes     ]
//            [simple DeviceIdBlockRes         deviceIdBlockRes        ]
//            [simple DeviceVendorBlockRes     deviceVendorBlockRes    ]
//            [simple DeviceOptionsBlockRes    deviceOptionsBlockRes   ]
//            [simple DeviceRoleBlockRes       deviceRoleBlockRes      ]
//            [simple DeviceInitiativeBlockRes deviceInitiativeBlockRes]
//        ]

        ////////////////////////////////////////////////////////////////////////////
        // Unicast
        ////////////////////////////////////////////////////////////////////////////

//        ['DCP_GetSet_PDU','GET','false' PnDcp_Pdu_GetReq
//            [simple GetReqBlock              getReqBlock             ]
//        ]
//        ['DCP_GetSet_PDU','GET','true' PnDcp_Pdu_GetRes
//            [simple GetResBlock              getResBlock             ]
//            [simple GetNegResBlock           getNegResBlock          ]
//        ]

//        ['DCP_GetSet_PDU','SET','false' PnDcp_Pdu_SetReq
//            [simple StartTransactionBlock    startTransactionBlock   ]
//            [simple BlockQualifier           blockQualifier          ]
//            [simple SetResetReqBlock         setResetReqBlock        ]
//            [simple SetReqBlock              setReqBlock             ]
//            [simple StopTransactionBlock     stopTransactionBlock    ]
//            [simple BlockQualifier           blockQualifier          ]
//        ]
//        ['DCP_GetSet_PDU','SET','true' PnDcp_Pdu_SetRes
//            [simple SetResBlock              setResBlock             ]
//            [simple SetNegResBlock           setNegResBlock          ]
//        ]
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

// 4.3.1.3.2 (Page 94ff)
// 4.3.1.3.2 (Page 94ff)
// The spec lists meanings for request and response separately, but
// they are actually mergeable, which we did in this construct.
[type PnDcp_ServiceType byteOrder='BIG_ENDIAN'
    [reserved uint 5 '0x00'      ]
    [simple   bit    notSupported]
    [reserved uint 1 '0x00'      ]
    [simple   bit    response    ]
]

// Page 86ff: Coding of the field FrameID
[enum uint 16 PnDcp_FrameId
    ['0x0000' RESERVED                       ]
    // Range 1
    ['0x0020' PTCP_RTSyncPDUWithFollowUp     ]
    // Range 2
    ['0x0080' PTCP_RTSyncPDU                 ]
    // Range 3
    // 0x0100-0x0FFF RT_CLASS_3
    ['0x0100' RT_CLASS_3                     ]
    // Range 4
    // (Not used)
    // Range 5
    // (Not used)
    // Range 6
    // 0x8000-BFFF RT_CLASS_1
    ['0x8000' RT_CLASS_1                     ]
    // Range 7
    // 0XC000-FBFF RT_CLASS_UDP
    ['0xC000' RT_CLASS_UDP                   ]
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
    ['0xFF80' FragmentationFrameId           ]
]

// Page 94
// All other values are "Reserved"
[enum uint 8 PnDcp_ServiceId
    ['0x03' GET     ]
    ['0x04' SET     ]
    ['0x05' IDENTIFY]
    ['0x06' HELLO   ]
    //[RESERVED]
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
