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

[type AdsDiscoveryConstants
    [const          uint 16     adsDiscoveryUdpDefaultPort 58899]
]

////////////////////////////////////////////////////////////////
// AMS/TCP Packet
////////////////////////////////////////////////////////////////
//
// All discovery requests target the ams port SYSTEM_SERVICE
//
// A Discovery Request is usually targeted at the broadcast ip of the current network
// instead of a real IP address (last segment is "255").
// The PLCs then respond to that.
// The Discovery Request generally contains no blocks.
//
// A Discovery Response usually contains 3 or 4 blocks.
// In all cases does it contain:
// - AdsDiscoveryBlockHostname (Host name of the PLC)
// - AdsDiscoveryBlockOsData   (Some information on the type and OS)
// - AdsDiscoveryBlockVersion  (The TwinCat version)
// Some also send a AdsDiscoveryBlockFingerprint
//
// A request for adding or deleting a route usually contains:
// - AdsDiscoveryBlockRouteName
// - AdsDiscoveryBlockAmsNetId
// - AdsDiscoveryBlockUserName
// - AdsDiscoveryBlockPassword
// - AdsDiscoveryBlockHostName
//
// The responses to adding or deleting a route usually only contain one block,
// containing the AdsDiscoveryBlockStatus

[type AdsDiscovery byteOrder='LITTLE_ENDIAN'
    [const    uint 32           header                          0x71146603                 ]
    [simple   uint 32           requestId                                                  ]
    [simple   Operation         operation                                                  ]
    [simple   AmsNetId          amsNetId                                                   ]
    [simple   AdsPortNumbers    portNumber                                                 ]
    [implicit uint 32           numBlocks                       'COUNT(blocks)'            ]
    [array    AdsDiscoveryBlock blocks                          count           'numBlocks']
]

[discriminatedType AdsDiscoveryBlock
    [discriminator    AdsDiscoveryBlockType    blockType                                   ]
    [typeSwitch    blockType
        ['STATUS' *Status
            [const    uint 16   statusLength   0x0004                                      ]
            [simple   Status    status                                                     ]
        ]
        ['PASSWORD' *Password
            [simple   AmsString password                                                   ]
        ]
        ['VERSION' *Version
            [implicit uint 16   versionDataLen 'COUNT(versionData)'                        ]
            [array    byte      versionData    count                'versionDataLen'       ]
        ]
        ['OS_DATA' *OsData
            [implicit uint 16   osDataLen      'COUNT(osData)'                             ]
            [array    byte      osData         count                'osDataLen'            ]
        ]
        ['HOST_NAME' *HostName
            [simple   AmsString hostName                                                   ]
        ]
        ['AMS_NET_ID' *AmsNetId
            [const    uint 16   amsNetIdLength 0x0006                                      ]
            [simple   AmsString amsNetId                                                   ]
        ]
        ['ROUTE_NAME' *RouteName
            [simple   AmsString routeName                                                  ]
        ]
        ['USER_NAME' *UserName
            [simple   AmsString userName                                                   ]
        ]
        ['FINGERPRINT' *Fingerprint
            [implicit uint 16   dataLen        'COUNT(data)'                               ]
            [array    byte      data           count                'dataLen'              ]
        ]
    ]
]

[enum uint 16 AdsDiscoveryBlockType
    ['0x0001' STATUS]
    ['0x0002' PASSWORD]
    ['0x0003' VERSION]
    ['0x0004' OS_DATA]
    ['0x0005' HOST_NAME]
    ['0x0007' AMS_NET_ID]
    ['0x000C' ROUTE_NAME]
    ['0x000D' USER_NAME]
    ['0x0012' FINGERPRINT]
]

[enum uint 32 Operation
    ['0x00000000' GET_AMS_NET_ID_REQUEST ]
    ['0x80000000' GET_AMS_NET_ID_RESPONSE]
    ['0x00000001' DISCOVERY_REQUEST      ]
    ['0x80000001' DISCOVERY_RESPONSE     ]
    ['0x00000006' ADD_ROUTE_REQUEST      ]
    ['0x80000006' ADD_ROUTE_RESPONSE     ]
    ['0x00000007' DEL_ROUTE_REQUEST      ]
    ['0x80000007' DEL_ROUTE_RESPONSE     ]
]

[enum uint 32 Status
    ['0x00000000' SUCCESS]
    ['0x000407' FAILURE]
]

[type AmsString
    [implicit uint 16                 len    'STR_LEN(text) + 1']
    [simple   vstring '8 * (len - 1)' text   encoding='"UTF-8"' ]
    [reserved uint 8                  '0x00'                    ]
]

[type AmsNetId
    [simple uint 8 octet1]
    [simple uint 8 octet2]
    [simple uint 8 octet3]
    [simple uint 8 octet4]
    [simple uint 8 octet5]
    [simple uint 8 octet6]
]

// From: https://infosys.beckhoff.com/index.php?content=../content/1031/tcadscommon/html/tcadscommon_identadsdevice.htm&id=
[enum uint 16 AdsPortNumbers
    ['100'   LOGGER              ]
    ['110'   EVENT_LOGGER        ]
    ['300'   IO                  ]
    ['301'   ADDITIONAL_TASK_1   ]
    ['302'   ADDITIONAL_TASK_2   ]
    ['500'   NC                  ]
    ['801'   PLC_RUNTIME_SYSTEM_1]
    ['811'   PLC_RUNTIME_SYSTEM_2]
    ['821'   PLC_RUNTIME_SYSTEM_3]
    ['831'   PLC_RUNTIME_SYSTEM_4]
    ['900'   CAM_SWITCH          ]
    ['10000' SYSTEM_SERVICE      ]
    ['14000' SCOPE               ]
]
