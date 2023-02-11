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
            [simple     uint 8                        chassisIdSubType              ]
            [simple     vstring     '(tlvIdLength - 1) * 8' chassisId               ]
        ]
        ['PORT_ID'   TlvPortId(uint 9 tlvIdLength)
            [simple     uint 8          portIdSubType                               ]
            [simple     vstring     '(tlvIdLength - 1) * 8' portId                  ]
        ]
        ['TIME_TO_LIVE'   TlvTimeToLive
            [simple     uint 16         tlvTimeToLiveUnit                           ]
        ]
        ['PORT_DESCRIPTION'   TlvPortDescription(uint 9 tlvIdLength)
            [simple     vstring     '(tlvIdLength) * 8' chassisId                   ]
        ]
        ['SYSTEM_NAME'   TlvSystemName(uint 9 tlvIdLength)
            [simple     vstring     '(tlvIdLength) * 8' chassisId                   ]
        ]
        ['SYSTEM_DESCRIPTION'   TlvSystemDescription(uint 9 tlvIdLength)
            [simple     vstring     '(tlvIdLength) * 8' chassisId                   ]
        ]
        ['SYSTEM_CAPABILITIES'   TlvSystemCapabilities
            [reserved   uint 8                          '0x00'                      ]
            [simple     bit             stationOnlyCapable                          ]
            [simple     bit             docsisCableDeviceCapable                    ]
            [simple     bit             telephoneCapable                            ]
            [simple     bit             routerCapable                               ]
            [simple     bit             wlanAccessPointCapable                      ]
            [simple     bit             bridgeCapable                               ]
            [simple     bit             repeaterCapable                             ]
            [simple     bit             otherCapable                                ]
            [reserved   uint 8                          '0x00'                      ]
            [simple     bit             stationOnlyEnabled                          ]
            [simple     bit             docsisCableDeviceEnabled                    ]
            [simple     bit             telephoneEnabled                            ]
            [simple     bit             routerEnabled                               ]
            [simple     bit             wlanAccessPointEnabled                      ]
            [simple     bit             bridgeEnabled                               ]
            [simple     bit             repeaterEnabled                             ]
            [simple     bit             otherEnabled                                ]
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
            [simple     TlvOrgSpecificIeee8023Unit      specificUnit               ]
        ]
    ]
]

[discriminatedType TlvOrgSpecificIeee8023Unit
    [discriminator  TlvIEEESubType  subType]
    [typeSwitch subType
        ['MAC_PHY_CONFIG_STATUS'  TlvIeee8023MacPhyConfigStatus
            [simple     uint 8                          negotiationSupport         ]
            [simple     uint 16                         negotiationCapability      ]
            [simple     uint 16                         operationalMauType         ]
        ]
        ['MAX_FRAME_SIZE'  TlvIeee8023MaxFrameSize
            [simple     uint 16                         maxSize                    ]
        ]
    ]
]

[discriminatedType TlvOrgSpecificProfibusUnit
    [discriminator  TlvProfibusSubType  subType]
    [typeSwitch subType
        ['MEASURED_DELAY'  TlvProfibusSubTypeMeasuredDelay
            [simple     uint 32                         localPortRxDelay]
            [simple     uint 32                         remotePortRxDelay]
            [simple     uint 32                         localPortTxDelay]
            [simple     uint 32                         remotePortTxDelay]
        ]
        ['PORT_STATUS'  TlvProfibusSubTypePortStatus
            [simple     uint 16                         rtClass2PortStatus]
            [reserved   uint 2                          '0x00'           ]
            [simple     bit                             preample         ]
            [simple     bit                             fragmentation    ]
            [reserved   uint 9                          '0x00'           ]
            [simple     uint 3                          rtClass3PortStatus]
        ]
        ['MRP_PORT_STATUS'  TlvProfibusSubTypeMrpPortStatus
            [simple     Uuid                      macAddress]
            [simple     uint 16                   Status]
        ]
        ['CHASSIS_MAC'  TlvProfibusSubTypeChassisMac
            [simple     MacAddress                      macAddress]
        ]
    ]
]

[enum  uint 8 TlvIEEESubType
    ['0x01' MAC_PHY_CONFIG_STATUS]
    ['0x04' MAX_FRAME_SIZE]
]

[enum  uint 8 TlvProfibusSubType
    ['0x01' MEASURED_DELAY]
    ['0x02' PORT_STATUS]
    ['0x04' MRP_PORT_STATUS]
    ['0x05' CHASSIS_MAC]
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
