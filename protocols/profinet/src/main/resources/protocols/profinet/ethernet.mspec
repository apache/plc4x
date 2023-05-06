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

[type PascalString
    [implicit int 8 sLength          'stringValue.length == 0 ? -1 : stringValue.length']
    [simple vstring 'sLength == -1 ? 0 : sLength * 8' stringValue]
    [virtual  int 8 stringLength     'stringValue.length == -1 ? 0 : stringValue.length']
]

[type PascalString16BitLength
    [implicit int 16 sLength          'stringValue.length == 0 ? -1 : stringValue.length']
    [simple vstring 'sLength == -1 ? 0 : sLength * 8' stringValue]
    [virtual  int 16 stringLength     'stringValue.length == -1 ? 0 : stringValue.length']
]

[type Uuid
    [array byte data count '16']
]
