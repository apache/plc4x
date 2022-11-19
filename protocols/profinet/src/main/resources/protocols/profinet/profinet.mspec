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

[type Ethernet_Frame byteOrder='BIG_ENDIAN'
    // When sending to the mac address prefix of 01:0e:cf are multicast packets
    [simple MacAddress            destination]
    [simple MacAddress            source     ]
    [simple Ethernet_FramePayload payload    ]
]

[discriminatedType Ethernet_FramePayload
    [discriminator uint 16 packetType]
    [typeSwitch packetType
        ['0x0800' Ethernet_FramePayload_IPv4
            [const    uint 4              version                         0x4                        ]
            // 5 = 5 x 32bit = 5 x 4byte = 20byte
            [const    uint 4              headerLength                    0x5                        ]
            [const    uint 6              differentiatedServicesCodepoint 0x00                       ]
            [const    uint 2              explicitCongestionNotification  0x0                        ]
            // Length of the header + payload
            [implicit uint 16             totalLength                     '28 + payload.lengthInBytes']
            [simple   uint 16             identification                                             ]
            [reserved bit                                                 'false'                    ]
            [simple   bit                 dontFragment                                               ]
            [simple   bit                 moreFragments                                              ]
            [const    uint 13             fragmentOffset                  0x00                       ]
            [simple   uint 8              timeToLive                                                 ]
            // Protocol: UDP
            [const    uint 8              protocol                        0x11                       ]
            // It seems that modern NICs mess this up as they take care of the validation in dedicated hardware.
            // This results in the wrong values being read. Using a 'checksum' field would fail most incoming packets.
            [implicit uint 16             headerChecksum                 'STATIC_CALL("calculateIPv4Checksum", totalLength, identification, timeToLive, sourceAddress, destinationAddress)']
            [simple   IpAddress           sourceAddress                                              ]
            [simple   IpAddress           destinationAddress                                         ]
            // Begin of the UDP packet part
            [simple   uint 16             sourcePort                                                 ]
            [simple   uint 16             destinationPort                                            ]
            [implicit uint 16             packetLength    '8 + payload.lengthInBytes'                ]
            [implicit uint 16             bodyChecksum                    'STATIC_CALL("calculateUdpChecksum", sourceAddress, destinationAddress, sourcePort, destinationPort, packetLength, payload)']
            [simple   DceRpc_Packet       payload                                                    ]
        ]
        ['0x8100' Ethernet_FramePayload_VirtualLan
            [simple VirtualLanPriority    priority                                                   ]
            [simple bit                   ineligible                                                 ]
            [simple uint 12               id                                                         ]
            [simple Ethernet_FramePayload payload                                                    ]
        ]
        ['0x8892' Ethernet_FramePayload_PnDcp
            [simple PnDcp_Pdu             pdu                                                        ]
        ]
        ['0x88cc' Ethernet_FramePayload_LLDP
            [simple Lldp_Pdu      pdu                                               ]
        ]
    ]
]

[type Lldp_Pdu
    [manualArray LldpUnit lldpParameters terminated 'STATIC_CALL("isSysexEnd", readBuffer)' 'STATIC_CALL("parseSysexString", readBuffer)' 'STATIC_CALL("serializeSysexString", writeBuffer, _value)' 'STATIC_CALL("lengthSysexString", lldpParameters)']
]

[discriminatedType LldpUnit
    [discriminator     TlvType                  tlvId                                ]
    [simple            uint 9                   tlvIdLength                          ]
    [typeSwitch tlvId
        ['END_OF_LLDP'  EndOfLldp
        ]
        ['CHASSIS_ID'   TlvChassisId(uint 9 tlvIdLength)
            [simple     uint 8                        chassisIdSubType               ]
            [simple     vstring     '(tlvIdLength - 1) * 8' chassisId                      ]
        ]
        ['PORT_ID'   TlvPortId(uint 9 tlvIdLength)
            [simple     uint 8          portIdSubType                              ]
            [simple     vstring     '(tlvIdLength - 1) * 8' portId                           ]
        ]
        ['TIME_TO_LIVE'   TlvTimeToLive
            [simple     uint 16         tlvTimeToLiveUnit                          ]
        ]
        ['MANAGEMENT_ADDRESS' TlvManagementAddress
            [implicit   uint 8          addressStringLength    '5' ]
            [simple     ManagementAddressSubType  addressSubType                   ]
            [simple     IpAddress       ipAddress                                  ]
            [simple     uint 8          interfaceSubType                           ]
            [simple     uint 32         interfaceNumber                            ]
            [simple     uint 8          oidStringLength                            ]
        ]
        ['ORGANIZATION_SPECIFIC' TlvOrganizationSpecific
            [simple     TlvOrganizationSpecificUnit     organizationSpecificUnit   ]
        ]
    ]
]

[discriminatedType TlvOrganizationSpecificUnit
    [discriminator      uint 24         uniqueCode]
    [typeSwitch uniqueCode
        ['0x000ECF' TlvOrgSpecificProfibus
            [simple     TlvOrgSpecificProfibusUnit      specificUnit               ]
        ]
        ['0x00120F' TlvOrgSpecificIeee8023
            [simple     uint 8                          subType                    ]
            [simple     uint 8                          negotiationSupport         ]
            [simple     uint 16                         negotiationCapability      ]
            [simple     uint 16                         operationalMauType         ]
        ]
    ]
]

[discriminatedType TlvOrgSpecificProfibusUnit
    [discriminator  TlvProfibusSubType  subType]
    [typeSwitch subType
        ['PORT_STATUS'  TlvProfibusSubTypePortStatus
            [simple     uint 16                         rtClass2PortStatus]
            [reserved   uint 2                          '0x00'           ]
            [simple     bit                             preample         ]
            [simple     bit                             fragmentation    ]
            [reserved   uint 9                          '0x00'           ]
            [simple     uint 3                          rtClass3PortStatus]
        ]
        ['CHASSIS_MAC'  TlvProfibusSubTypeChassisMac
            [simple     MacAddress                      macAddress]
        ]
    ]
]

[enum  uint 8 TlvProfibusSubType
    ['0x02' PORT_STATUS]
    ['0x05' CHASSIS_MAC]
]

// 4.10.3.2
// A lot of the fields are set to constant values, which would
// usually be dynamic. However are we trying to limit the number of
// arguments needed to construct the messages and Profinet only seems
// be using a very limited subset of all possible DCE/RPC packets.
[type DceRpc_Packet byteOrder='BIG_ENDIAN'
// RPC Header {
    // RPCVersion 4.10.3.2.1
    [const         uint 8                version                        0x04                 ]
    // RPCPacketType 4.10.3.2.2 (8 bit)
    [simple        DceRpc_PacketType     packetType                                            ]
    // PRCFlags 4.10.3.2.3
    [reserved      bit                                                  'false'                ]
    [const         bit                   broadcast                      false                ]
    [simple        bit                   idempotent                                            ]
    [const         bit                   maybe                          false                ]
    [simple        bit                   noFragmentAcknowledgeRequested                        ]
    [const         bit                   fragment                       false                ]
    [simple        bit                   lastFragment                                          ]
    [reserved      bit                                                  'false'                ]
    // PRCFlags2 4.10.3.2.4
    [reserved      uint 6                                               '0x00'                 ]
    [const         bit                   cancelWasPending               false                ]
    [reserved      bit                                                  'false'                ]
    // RPCDRep 4.10.3.2.5 (4 bit & 4 bit)
    [simple        IntegerEncoding       integerEncoding                                       ]
    [simple        CharacterEncoding     characterEncoding                                     ]
    // RPCDRep2 4.10.3.2.5 (8 bit)
    [simple        FloatingPointEncoding floatingPointEncoding                                 ]
    // RPCDRep3 (8 bit shall be 0)
    [reserved      uint 8                                               '0x00'                 ]
    // RPCSerialHigh 4.10.3.2.6
    [const         uint 8                serialHigh                     0x00                 ]
    [batchSet byteOrder='integerEncoding == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
        // RPCObjectUUID 4.10.3.2.8
        // RPCObjectUUID 4.10.3.2.8
        [simple DceRpc_ObjectUuid        objectUuid                                            ]
        // RPCInterfaceUUID 4.10.3.2.9
        [simple DceRpc_InterfaceUuid     interfaceUuid                                         ]
        // RPCActivityUUID 4.10.3.2.10
        [simple DceRpc_ActivityUuid      activityUuid                                          ]
        // RPCServerBootTime 4.10.3.2.11
        [simple uint 32                  serverBootTime                                        ]
        // RPCInterfaceVersion 4.10.3.2.12
        [const  uint 32                  interfaceVer                   0x00000001           ]
        // RPCSequenceNmb 4.10.3.2.13
        [simple uint 32                  sequenceNumber                                        ]
        // RPCOperationNmb 4.10.3.2.14
        [simple DceRpc_Operation         operation                                             ]
        // RPCInterfaceHint 4.10.3.2.15
        [const        uint 16            interfaceHint                  0xFFFF               ]
        // RPCActivityHint 4.10.3.2.16
        [const        uint 16            activityHint                   0xFFFF               ]
        // RPCLengthOfBody 4.10.3.2.17
        [implicit     uint 16            lengthOfBody                   'payload.lengthInBytes']
        // RPCFragmentNmb 4.10.3.2.18 (Setting this to 0 as we will probably never have anything but 0 here
        [const        uint 16            fragmentNum                    0x0000               ]
        // RPCAuthenticationProtocol 4.10.3.2.19
        [const        uint 8             authProto                      0x00                 ]
    ]
    // RPCSerialLow 4.10.3.2.7
    [const            uint 8             serialLow                      0x00                 ]
// RPC Header }
// RPC Payload {
    [simple PnIoCm_Packet('packetType') payload byteOrder='integerEncoding == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN' ]
// RPC Payload }
]

// RPCObjectUUID 4.10.3.2.8
[type DceRpc_ObjectUuid
    [const  uint 32 data1      0xDEA00000                       ]
    [const  uint 16 data2      0x6C97                           ]
    [const  uint 16 data3      0x11D1                           ]
    // This part is described as a byte array, so the byte order is always big-endian
    [const  uint 16 data4      0x8271     byteOrder='BIG_ENDIAN']
    [simple uint 4  interfaceNumber       byteOrder='BIG_ENDIAN']
    [simple uint 12 nodeNumber            byteOrder='BIG_ENDIAN']
    [simple uint 16 deviceId              byteOrder='BIG_ENDIAN']
    [simple uint 16 vendorId              byteOrder='BIG_ENDIAN']
]

// RPCInterfaceUUID 4.10.3.2.9
// NOTE: If we would have been only using Big Endian encoding, we would have
//       implemented this via an enum. However as the first 8 bytes are
//       dynamically endianed and the last 8 bytes are set to Big Endian, we
//       had to do this trick.
[discriminatedType DceRpc_InterfaceUuid
    [discriminator  uint 32 interfaceType                               ]
    [const          uint 16 data1      0x6C97                           ]
    [const          uint 16 data2      0x11D1                           ]
    // This part is described as a byte array, so the byte order is always big-endian
    [const          uint 16 data3      0x8271     byteOrder='BIG_ENDIAN']
    [const          uint 16 data4      0x00A0     byteOrder='BIG_ENDIAN']
    [const          uint 16 data5      0x2442     byteOrder='BIG_ENDIAN']
    [const          uint 16 data6      0xDF7D     byteOrder='BIG_ENDIAN']
    [typeSwitch interfaceType
        ['0xDEA00001' DceRpc_InterfaceUuid_DeviceInterface
        ]
        ['0xDEA00002' DceRpc_InterfaceUuid_ControllerInterface
        ]
        ['0xDEA00003' DceRpc_InterfaceUuid_SupervisorInterface
        ]
        ['0xDEA00004' DceRpc_InterfaceUuid_ParameterInterface
        ]
    ]
]

// RPCActivityUUID 4.10.3.2.10
// NOTE: This value is generally randomly generated by the initiator
//       and used throughout the entire communication. Unfortunately,
//       the first parts are effected by endianess, and the last 8
//       bytes are fixed big-endian. Therefore the complicated notation.
[type DceRpc_ActivityUuid
    [simple  uint 32 data1          ]
    [simple  uint 16 data2          ]
    [simple  uint 16 data3          ]
    // This part is described as a byte array, so the byte order is always big-endian
    [array   byte    data4 count '8']
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

// 4.10.3.2.2
[enum uint 8 DceRpc_PacketType
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

//LLDP Specific
[enum uint 7 TlvType
    ['0x00' END_OF_LLDP          ]
    ['0x01' CHASSIS_ID           ]
    ['0x02' PORT_ID              ]
    ['0x03' TIME_TO_LIVE         ]
    ['0x04' PORT_DESCRIPTION     ]
    ['0x05' SYSTEM_NAME          ]
    ['0x06' SYSTEM_DESCRIPTION   ]
    ['0x07' SYSTEM_CAPABILITIES  ]
    ['0x08' MANAGEMENT_ADDRESS    ]
    ['0x7F' ORGANIZATION_SPECIFIC]
]

[enum uint 8 ManagementAddressSubType
    ['0x00' UNKNOWN              ]
    ['0x01' IPV4                 ]
]

// 4.10.3.2.14
[enum uint 16 DceRpc_Operation
    ['0x0000' CONNECT      ]
    ['0x0001' RELEASE      ]
    ['0x0002' READ         ]
    ['0x0003' WRITE        ]
    ['0x0004' CONTROL      ]
    ['0x0005' READ_IMPLICIT]
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
[discriminatedType PnDcp_Pdu
    [simple        uint 16           frameIdValue                      ]
    [virtual       PnDcp_FrameId     frameId       'STATIC_CALL("getFrameId", frameIdValue)']
    [typeSwitch frameId
        ['RT_CLASS_1' PnDcp_Pdu_RealTimeCyclic
            // TODO: This type needs to be implemented based of the configuration and gsd file ...
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

[discriminatedType PnDcp_Pdu_IdentifyRes_Payload
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

[type PnIo_CyclicServiceDataUnit(int 16 dataUnitLength)
    [array    byte   data       count 'dataUnitLength'                 ]
]

[discriminatedType PnDcp_Block
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

[type PnDcp_SupportedDeviceOption
    [simple   PnDcp_BlockOptions option]
    [simple   uint 8       suboption   ]
]

// 4.3.1.3.2 (Page 94ff)
// 4.3.1.3.2 (Page 94ff)
// The spec lists meanings for request and response separately, but
// they are actually mergeable, which we did in this construct.
[type PnDcp_ServiceType
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
[discriminatedType PnIoCm_Packet(DceRpc_PacketType packetType)
    [typeSwitch packetType
        ['REQUEST' PnIoCm_Packet_Req
            [simple uint 32      argsMaximum                          ]
            [implicit uint 32    argsLength       'lengthInBytes - 20']
            [simple uint 32      arrayMaximumCount                    ]
            [simple uint 32      arrayOffset                          ]
            [implicit uint 32    arrayActualCount  'lengthInBytes - 20']
            [array  PnIoCm_Block blocks            length 'argsLength']
        ]
        ['RESPONSE' PnIoCm_Packet_Res
            [simple uint 8       errorCode2                           ]
            [simple uint 8       errorCode1                           ]
            [simple uint 8       errorDecode                          ]
            [simple uint 8       errorCode                            ]
            [implicit uint 32    argsLength       'lengthInBytes - 1 - 1 - 1 - 1 - 4 - 4 - 4 - 4']
            [simple uint 32      arrayMaximumCount                    ]
            [simple uint 32      arrayOffset                          ]
            [implicit uint 32    arrayActualCount  'lengthInBytes - 1 - 1 - 1 - 1 - 4 - 4 - 4 - 4'    ]
            [array  PnIoCm_Block blocks            length 'argsLength']
        ]
        ['REJECT'   PnIoCm_Packet_Rej
            [simple uint 32      status                               ]
        ]
    ]
]

[type UserData(uint 32 recordDataLength)
    [array              byte      data count         'recordDataLength'       ]
]

// Big Endian
[discriminatedType PnIoCm_Block byteOrder='BIG_ENDIAN'
    [discriminator PnIoCm_BlockType blockType                           ]
    [typeSwitch blockType
        ['IOD_WRITE_REQUEST_HEADER' IODWriteRequestHeader
            [implicit      uint 16          blockLength      'index < 0x8000 ? lengthInBytes - 4 - recordDataLength : lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   uint 16                         sequenceNumber                                         ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 32                         api                                                    ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         index                                                  ]
            [simple   uint 32                         recordDataLength                                       ]
            [padding  uint 8      pad '0x00'          'index < 0x8000 ? 64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4 : 64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4']
            [optional UserData('recordDataLength')     userData          'index < 0x8000'                    ]
        ]
        ['IOD_WRITE_RESPONSE_HEADER' IODWriteResponseHeader
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   uint 16                         sequenceNumber                                         ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 32                         api                                                    ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         index                                                  ]
            [simple   uint 32                         recordDataLength                                       ]
            [padding  uint 8      pad '0x00'          '64 - 6 - 2 - 16 - 4 - 2 - 2 - 2 - 2 - 4']
        ]
        ['PD_INTERFACE_ADJUST' PDInterfaceAdjust
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 16                         padField                  0x0000                       ]
            [const    uint 16                         multipleInterfaceModeReserved2                  0x0000 ]
            [const    uint 15                         multipleInterfaceModeReserved1                  0x0000 ]
            [simple   MultipleInterfaceModeNameOfDevice multipleInterfaceModeNameOfDevice                    ]
        ]
        ['PD_PORT_DATA_CHECK' PDPortDataCheck
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 16                         padField                  0x0000                       ]
            [simple   uint 16                         slotNumber                                             ]
            [simple   uint 16                         subSlotNumber                                          ]
            [simple   PnIoCm_Block                    checkPeers                                             ]
        ]
        ['CHECK_PEERS'  CheckPeers
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [const    uint 8                          noOfPeers                 0x01                         ]
            [simple   PascalString                    peerPortId                                             ]
            [simple   PascalString                    peerChassisId                                          ]
        ]
        ['AR_BLOCK_REQ' PnIoCm_Block_ArReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_ArType                   arType                                                 ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [simple   MacAddress                      cmInitiatorMacAddr                                     ]
            [simple   DceRpc_ObjectUuid               cmInitiatorObjectUuid                                  ]
            // Begin ARProperties
            [simple   bit                             pullModuleAlarmAllowed                                 ]
            [simple   bit                             nonLegacyStartupMode                                   ]
            [simple   bit                             combinedObjectContainerUsed                            ]
            [reserved uint 17                         '0x00000'                                              ]
            [simple   bit                             acknowledgeCompanionAr                                 ]
            [simple   PnIoCm_CompanionArType          companionArType                                        ]
            [simple   bit                             deviceAccess                                           ]
            [reserved uint 3                          '0x0'                                                  ]
            [simple   bit                             cmInitiator                                            ]
            [simple   bit                             supervisorTakeoverAllowed                              ]
            [simple   PnIoCm_State                    state                                                  ]
            // End ARProperties
            [simple   uint 16                         cmInitiatorActivityTimeoutFactor                       ]
            [simple   uint 16                         cmInitiatorUdpRtPort                                   ]
            [implicit uint 16                         stationNameLength     'STR_LEN(cmInitiatorStationName)']
            [simple   vstring 'stationNameLength * 8' cmInitiatorStationName                                 ]
        ]
        ['AR_BLOCK_RES' PnIoCm_Block_ArRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_ArType          arType                                                          ]
            [simple   Uuid                   arUuid                                                          ]
            [simple   uint 16                sessionKey                                                      ]
            [simple   MacAddress             cmResponderMacAddr                                              ]
            [simple   uint 16                responderUDPRTPort                                              ]
        ]
        ['IOD_CONTROL_REQ' PnIoCm_Control_Request
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [reserved uint 16                         '0x0000'                                               ]
        ]
        ['IOX_BLOCK_REQ'    PnIoCM_Block_Request
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [simple   uint 16                         controlBlockProperties                                 ]
        ]
        ['IOX_BLOCK_RES'    PnIoCM_Block_Response
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [simple   uint 16                         controlBlockProperties                                 ]
        ]
        ['IOD_CONTROL_RES' PnIoCm_Control_Response
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   Uuid                            arUuid                                                 ]
            [simple   uint 16                         sessionKey                                             ]
            [reserved uint 16                         '0x0000'                                               ]
            [simple   uint 16                         controlCommand                                         ]
            [reserved uint 16                         '0x0000'                                               ]
        ]
        ['IO_CR_BLOCK_REQ' PnIoCm_Block_IoCrReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple PnIoCm_IoCrType          ioCrType                                               ]
            [simple uint 16                  ioCrReference                                          ]
            [simple uint 16                  lt                                                     ]
            // Begin IOCRProperties
            [simple   bit                    fullSubFrameStructure                                  ]
            [simple   bit                    distributedSubFrameWatchDog                            ]
            [simple   bit                    fastForwardingMacAdr                                   ]
            [reserved uint 17                '0x0000'                                               ]
            [simple   bit                    mediaRedundancy                                        ]
            [reserved uint 7                 '0x00'                                                 ]
            [simple   PnIoCm_RtClass         rtClass                                                ]
            // End IOCRProperties
            [simple   uint 16                dataLength                                             ]
            [simple   uint 16                frameId                                                ]
            [simple   uint 16                sendClockFactor                                        ]
            [simple   uint 16                reductionRatio                                         ]
            [simple   uint 16                phase                                                  ]
            [simple   uint 16                sequence                                               ]
            [simple   uint 32                frameSendOffset                                        ]
            [simple   uint 16                watchDogFactor                                         ]
            [simple   uint 16                dataHoldFactor                                         ]
            [simple   uint 16                ioCrTagHeader                                          ]
            [simple   MacAddress             ioCrMulticastMacAdr                                    ]
            [implicit uint 16                numberOfApis        'COUNT(apis)'                      ]
            [array    PnIoCm_IoCrBlockReqApi apis                count         'numberOfApis'       ]
        ]
        ['IO_CR_BLOCK_RES' PnIoCm_Block_IoCrRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple PnIoCm_IoCrType          ioCrType                                               ]
            [simple uint 16                  ioCrReference                                          ]
            [simple   uint 16                frameId                                                ]
        ]
        ['ALARM_CR_BLOCK_REQ' PnIoCm_Block_AlarmCrReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_AlarmCrType     alarmType                                              ]
            [simple   uint 16                lt                                                     ]
            // Begin AlarmCrProperties
            [reserved uint 30                '0x00000000'                                           ]
            [simple   bit                    transport                                              ]
            [simple   bit                    priority                                               ]
            // End AlarmCrProperties
            [simple   uint 16                rtaTimeoutFactor                                       ]
            [simple   uint 16                rtaRetries                                             ]
            [simple   uint 16                localAlarmReference                                    ]
            [simple   uint 16                maxAlarmDataLength                                     ]
            [simple   uint 16                alarmCtrTagHeaderHigh                                  ]
            [simple   uint 16                alarmCtrTagHeaderLow                                   ]
        ]
        ['ALARM_CR_BLOCK_RES' PnIoCm_Block_AlarmCrRes
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PnIoCm_AlarmCrType     alarmType                                              ]
            [simple   uint 16                localAlarmReference                                    ]
            [simple   uint 16                maxAlarmDataLength                                     ]
        ]
        ['EXPECTED_SUBMODULE_BLOCK_REQ' PnIoCm_Block_ExpectedSubmoduleReq
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [implicit uint 16                numberOfApis         'COUNT(apis)'                     ]
            [array    PnIoCm_ExpectedSubmoduleBlockReqApi apis   count         'numberOfApis'       ]
        ]
        ['MODULE_DIFF_BLOCK' PnIoCm_Block_ModuleDiff
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [implicit uint 16                numberOfApis         'COUNT(apis)'                     ]
            [array    PnIoCm_ModuleDiffBlockApi apis              count         'numberOfApis'      ]
        ]
        ['AR_SERVER_BLOCK' PnIoCm_Block_ArServer
            [implicit      uint 16          blockLength      'lengthInBytes - 4']
            [simple        uint 8           blockVersionHigh                    ]
            [simple        uint 8           blockVersionLow                     ]
            [simple   PascalString                    stationName                                   ]
            [padding  uint 8      pad '0x00'          '20 - 6 - (stationName.stringLength)'              ]
        ]
    ]
]

[type PascalString
    [implicit int 8 sLength          'stringValue.length == 0 ? -1 : stringValue.length']
    [simple vstring 'sLength == -1 ? 0 : sLength * 8' stringValue]
    [virtual  int 8 stringLength     'stringValue.length == -1 ? 0 : stringValue.length']
]

[type PnIoCm_IoCrBlockReqApi
    [const    uint 32             api              0x00000000             ]
    [implicit uint 16             numIoDataObjects 'COUNT(ioDataObjects)'   ]
    [array    PnIoCm_IoDataObject ioDataObjects    count 'numIoDataObjects' ]
    [implicit uint 16             numIoCss         'COUNT(ioCss)'           ]
    [array    PnIoCm_IoCs         ioCss            count 'numIoCss'         ]
]

[type PnIoCm_IoDataObject
    [simple   uint 16 slotNumber             ]
    [simple   uint 16 subSlotNumber          ]
    [simple   uint 16 ioDataObjectFrameOffset]
]

[type PnIoCm_IoCs
    [simple   uint 16 slotNumber   ]
    [simple   uint 16 subSlotNumber]
    [simple   uint 16 ioFrameOffset]
]

[type PnIoCm_DataUnitIoCs
    [simple   bit               dataState]
    [simple   uint 2            instance ]
    [reserved uint 4            '0x00'   ]
    [simple   bit               extension]
]

[type PnIoCm_DataUnitDataObject(uint 16 dataObjectLength)
    [array    byte              dataState   count  'dataObjectLength']
    [simple   PnIoCm_DataUnitIoCs iops   ]
]

[type PnIoCm_ExpectedSubmoduleBlockReqApi
    [const    uint 32          api               0x00000000                       ]
    [simple   uint 16          slotNumber                                           ]
    [simple   uint 32          moduleIdentNumber                                    ]
    [simple   uint 16          moduleProperties                                     ]
    [implicit uint 16          numSubmodules     'COUNT(submodules)'                ]
    [array    PnIoCm_Submodule submodules        count               'numSubmodules']
]

[type PnIoCm_ModuleDiffBlockApi
    [const    uint 32                          api        0x00000000                    ]
    [implicit uint 16                          numModules 'COUNT(modules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Module modules    count               'numModules']
]

[type PnIoCm_ModuleDiffBlockApi_Module
    [simple   uint 16                             slotNumber                                           ]
    [simple   uint 32                             moduleIdentNumber                                    ]
    [simple   PnIoCm_ModuleState                  moduleState                                          ]
    [implicit uint 16                             numSubmodules     'COUNT(submodules)'                ]
    [array    PnIoCm_ModuleDiffBlockApi_Submodule submodules        count               'numSubmodules']
]

[type PnIoCm_ModuleDiffBlockApi_Submodule
    [simple uint 16          subslotNumber       ]
    [simple uint 32          submoduleIdentNumber]
    [simple bit              codingUsesIdentInfo ]
    [simple PnIoCm_IdentInfo identInfo           ]
    [simple PnIoCm_ArInfo    arInfo              ]
    [simple bit              diagInfoAvailable   ]
    [simple bit              maintenanceDemanded ]
    [simple bit              maintenanceRequired ]
    [simple bit              qualifiedInfo       ]
    [simple PnIoCm_AddInfo   addInfo             ]
]

[discriminatedType PnIoCm_Submodule
    [simple        uint 16                slotNumber                    ]
    [simple        uint 32                submoduleIdentNumber          ]
    // Begin SubmoduleProperties
    [reserved      uint 10                '0x000'                       ]
    [simple        bit                    discardIoxs                   ]
    [simple        bit                    reduceOutputModuleDataLength  ]
    [simple        bit                    reduceInputModuleDataLength   ]
    [simple        bit                    sharedInput                   ]
    [discriminator PnIoCm_SubmoduleType   submoduleType                 ]
    // End SubmoduleProperties
    [typeSwitch submoduleType
        ['NO_INPUT_NO_OUTPUT_DATA' PnIoCm_Submodule_NoInputNoOutputData
            [const    uint 16             dataDescription       0x0001]
            [const    uint 16             submoduleDataLength   0x0000]
            [const    uint 8              lengthIoCs            0x01  ]
            [const    uint 8              lengthIoPs            0x01  ]
        ]
        ['INPUT_DATA' PnIoCm_Submodule_InputData
            [const    uint 16             inputDataDescription  0x0001]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
        ]
        ['OUTPUT_DATA' PnIoCm_Submodule_OutputData
            [const    uint 16             inputDataDescription  0x0002]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
        ]
        ['INPUT_AND_OUTPUT_DATA' PnIoCm_Submodule_InputAndOutputData
            [const    uint 16             inputDataDescription  0x0001]
            [simple   uint 16             inputSubmoduleDataLength      ]
            [simple   uint 8              inputLengthIoCs               ]
            [simple   uint 8              inputLengthIoPs               ]
            [const    uint 16             outputDataDescription 0x0002]
            [simple   uint 16             outputSubmoduleDataLength     ]
            [simple   uint 8              outputLengthIoCs              ]
            [simple   uint 8              outputLengthIoPs              ]
        ]
    ]
]

[type Uuid
    [array byte data count '16']
]

[enum uint 16 PnIoCm_BlockType
    ['0x0008' IOD_WRITE_REQUEST_HEADER    ]
    ['0x0101' AR_BLOCK_REQ                ]
    ['0x0102' IO_CR_BLOCK_REQ             ]
    ['0x0103' ALARM_CR_BLOCK_REQ          ]
    ['0x0104' EXPECTED_SUBMODULE_BLOCK_REQ]
    ['0x0110' IOD_CONTROL_REQ             ]
    ['0x0112' IOX_BLOCK_REQ               ]
    ['0x0200' PD_PORT_DATA_CHECK          ]
    ['0x020a' CHECK_PEERS                 ]
    ['0x0250' PD_INTERFACE_ADJUST         ]
    ['0x8008' IOD_WRITE_RESPONSE_HEADER    ]
    ['0x8101' AR_BLOCK_RES                ]
    ['0x8102' IO_CR_BLOCK_RES             ]
    ['0x8103' ALARM_CR_BLOCK_RES          ]
    ['0x8104' MODULE_DIFF_BLOCK           ]
    ['0x8106' AR_SERVER_BLOCK             ]
    ['0x8110' IOD_CONTROL_RES             ]
    ['0x8112' IOX_BLOCK_RES               ]
]

[enum uint 16 ProfinetDeviceState
    ['0x00'     IDLE]
    ['0x01'     STARTUP]
    ['0x02'     PRMEND]
    ['0x03'     APPLRDY]
    ['0x04'     ABORT]
]

[enum uint 16 PnIoCm_ArType
    ['0x0001' IO_CONTROLLER]
]

[enum uint 2 PnIoCm_CompanionArType
    ['0x0' SINGLE_AR]
]

[enum uint 3 PnIoCm_State
    ['0x1' ACTIVE]
]

[enum uint 16 PnIoCm_IoCrType
    ['0x0001' INPUT_CR]
    ['0x0002' OUTPUT_CR]
]

[enum uint 4 PnIoCm_RtClass
    ['0x2' RT_CLASS_2]
]

[enum uint 16 PnIoCm_AlarmCrType
    ['0x0001' ALARM_CR]
]

[enum uint 16 PnIoCm_ModuleState
    ['0x0002' PROPER_MODULE]
]

[enum uint 2 PnIoCm_SubmoduleType
    ['0x0' NO_INPUT_NO_OUTPUT_DATA]
    ['0x1' INPUT_DATA]
    ['0x2' OUTPUT_DATA]
    ['0x3' INPUT_AND_OUTPUT_DATA]
]

[enum bit MultipleInterfaceModeNameOfDevice
    ['false' PORT_PROVIDED_BY_LLDP]
    ['true'  NAME_PROVIDED_BY_LLDP]
]

[enum uint 16 PnIoCm_DescriptionType
    ['0x0001' INPUT]
]

[enum uint 4 PnIoCm_IdentInfo
    ['0x0' OK]
]

[enum uint 4 PnIoCm_ArInfo
    ['0x0' OWN]
]

[enum uint 3 PnIoCm_AddInfo
    ['0x0' NONE]
]

[enum uint 4 IntegerEncoding
    ['0x0' BIG_ENDIAN]
    ['0x1' LITTLE_ENDIAN]
]

[enum uint 4 CharacterEncoding
    ['0x0' ASCII]
    ['0x1' EBCDIC]
]

[enum uint 8 FloatingPointEncoding
    ['0x00' IEEE]
    ['0x01' VAX ]
    ['0x02' CRAY]
    ['0x03' IBM ]
]