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

////////////////////////////////////////////////////////////////
// AMS/TCP Paket
////////////////////////////////////////////////////////////////

[type 'AMSPacket'
    [reserved   uint       16       '0x0000'                                                   ]
    [implicit   uint       32       'packetLength' 'header.lengthInBytes + data.lengthInBytes' ]
    [simple     AMSHeader  'header'                                                            ]
    [simple     ADSData    'data'   ['header.commandId', 'header.state.response']              ]
]

[type 'AMSHeader'
    [simple     AMSNetId        'targetAmsNetId'                            ]
    [simple     uint        16  'targetAmsPort'                             ]
    [simple     AMSNetId        'sourceAmsNetId'                            ]
    [simple     uint        16  'sourceAmsPort'                             ]
    [enum       CommandId       'commandId'                                 ]
    [simple     State           'state'                                     ]
    [simple     uint        32  'dataLength'                                ]
    [simple     uint        32  'errorCode'                                 ]
    // free usable field of 4 bytes
    [simple      uint        32  'invokeId'                                 ]
]

[enum uint 16 'CommandId'
    ['0x00' INVALID]
    ['0x01' ADS_READ_DEVICE_INFO]
    ['0x02' ADS_READ]
    ['0x03' ADS_WRITE]
    ['0x04' ADS_READ_STATE]
    ['0x05' ADS_WRITE_CONTROL]
    ['0x06' ADS_ADD_DEVICE_NOTIFICATION]
    ['0x07' ADS_DELETE_DEVICE_NOTIFICATION]
    ['0x08' ADS_DEVICE_NOTIFICATION]
    ['0x09' ADS_READ_WRITE]
]

[type 'State'
    [simple     bit 'broadcast'             ]
    [reserved   int 7 '0x0'                 ]
    [simple     bit 'initCommand'           ]
    [simple     bit 'updCommand'            ]
    [simple     bit 'timestampAdded'        ]
    [simple     bit 'highPriorityCommand'   ]
    [simple     bit 'systemCommand'         ]
    [simple     bit 'adsCommand'            ]
    [simple     bit 'noReturn'              ]
    [simple     bit 'response'              ]
]

[type 'AMSNetId'
    [simple     uint        8   'octet1'            ]
    [simple     uint        8   'octet2'            ]
    [simple     uint        8   'octet3'            ]
    [simple     uint        8   'octet4'            ]
    [simple     uint        8   'octet5'            ]
    [simple     uint        8   'octet6'            ]
]

[discriminatedType 'ADSData' [CommandId 'commandId', bit 'response']
    [typeSwitch 'commandId', 'response'
        ['0x00', 'true' AdsInvalidResponse]
        ['0x00', 'false' AdsInvalidRequest]
        ['0x01', 'true' AdsReadDeviceInfoResponse
            // 4 bytes	ADS error number.
            [simple uint 32 'result']
            // Version	1 byte	Major version number
            [simple uint 8  'majorVersion']
            // Version	1 byte	Minor version number
            [simple uint 8  'minorVersion']
            // Build	2 bytes	Build number
            [simple uint 16  'version']
            // Name	16 bytes	Name of ADS device
            [array int 8  'device' count '16']
        ]
        ['0x01', 'false' AdsReadDeviceInfoRequest]
        ['0x02', 'true' Adstodo4]
        ['0x02', 'false' Adstodo5]
        ['0x03', 'true' Adstodo6]
        ['0x03', 'false' Adstodo7]
        ['0x04', 'true' Adstodo8]
        ['0x04', 'false' Adstodo9]
        ['0x05', 'true' Adstodo06]
        ['0x05', 'false' Adstodo60]
        ['0x06', 'true' Adstodo58]
        ['0x06', 'false' Adstodo45]
        ['0x07', 'true' Adstodo34]
        ['0x07', 'false' Adstodo23]
        ['0x08', 'true' Adstodo12]
        ['0x08', 'false' Adstodo23]
        ['0x09', 'true' Adstodo34]
        ['0x09', 'false' Adstodo45]
    ]
]