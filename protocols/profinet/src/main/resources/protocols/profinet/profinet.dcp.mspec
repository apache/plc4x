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
    [simple MacAddress 'destination' ]
    [simple MacAddress 'source'      ]
    [simple uint 16    'ethernetType']
    [simple ProfinetFrame 'payload'  ]
]

[type 'ProfinetFrame'
    [enum FrameType 'frameType'                ]
    [simple ProfinetData 'frame'  ['frameType']]
]

[discriminatedType 'ProfinetData' [FrameType 'frameType']
    [typeSwitch 'frameType'
        ['FrameType.IDENTIFY_MULTICAST_REQUEST' DcpIdentRequestPDU
            [enum DCPServiceID 'serviceId'                 ]
            [enum DCPServiceType 'serviceType'             ]
            [simple uint 32 'xid'                          ]
            [simple uint 16 'responseDelay'                ]
            [simple uint 16 'dcpDataLength'                ]
            [array DCPBlock 'blocks' length 'dcpDataLength']
        ]
        ['FrameType.IDENTIFY_RESPONSE' DcpIdentResponsePDU
            [enum DCPServiceID 'serviceId'                 ]
            [enum DCPServiceType 'serviceType'             ]
            [simple uint 32 'xid'                          ]
            [simple uint 16 'responseDelay'                ]
            [simple uint 16 'dcpDataLength'                ]
            [array DCPBlock 'blocks' length 'dcpDataLength']
        ]
    ]
]

[discriminatedType 'DCPBlock'
    [enum DCPBlockOption 'blockType']
    [typeSwitch 'blockType'
        ['DCPBlockOption.IP' IP
            [enum   IpSubOption 'subOption']
            [simple uint 16 'length'   ]
            [simple uint 16 'info'     ]
            [simple IPv4Address 'ipAddress'      ]
            [simple IPv4Address 'subnetMask'     ]
            [simple IPv4Address 'standardGateway']
        ]
        ['DCPBlockOption.DEVICE_PROPERTIES' DeviceProperties
            [enum   DevicePropertiesSubOption 'subOption' ]
            [simple DCPDeviceProperties       'properties' ['subOption']]
        ]
        ['DCPBlockOption.ALL_SELECTOR' AllSelector
            [const uint 8  'subOption'  '0xFF' ]
            [const uint 16 'length'     '0x00' ]
        ]
    ]
]

[discriminatedType 'DCPDeviceProperties' [DevicePropertiesSubOption 'subOptionType']
    [simple uint 16 'length' ]
    [simple uint 16 'info'   ]
    [typeSwitch 'subOptionType'
        ['DevicePropertiesSubOption.STATION_TYPE' StationType [uint 16 'length']
            [simple Text 'vendorNameForDevice' ['length - 2']]
        ]
        ['DevicePropertiesSubOption.STATION_NAME' StationName [uint 16 'length']
            [simple Text 'name' ['length - 2']]
        ]
        ['DevicePropertiesSubOption.DEVICE_ID' DeviceId
            [simple uint 16 'vendorId'         ]
            [simple uint 16 'deviceId'         ]
        ]
        ['DevicePropertiesSubOption.DEVICE_ROLE' DeviceRole
            [simple uint 8   'role'       ]
            [reserved uint 8 '0x00'       ]
        ]
        ['DevicePropertiesSubOption.DEVICE_OPTIONS' DeviceOptions
            [enum DCPBlockOption            'blockType'         ]
            [enum DevicePropertiesSubOption 'subOption'         ]
        ]
        ['DevicePropertiesSubOption.DEVICE_INSTANCE' DeviceInstance
            [simple uint 8 'instanceLow'         ]
            [simple uint 8 'instanceHigh'        ]
        ]
    ]
    [padding uint 8                 'pad' '0x00' '(length % 2 == 0 ? 0 : 1)']
]


[enum uint 8 'DCPBlockOption'
    ['0x01' IP                          ]
    ['0x02' DEVICE_PROPERTIES           ]
    ['0xFF' ALL_SELECTOR                ]
]

[enum uint 8 'DevicePropertiesSubOption'
    ['0x01' STATION_TYPE                ]
    ['0x02' STATION_NAME                ]
    ['0x03' DEVICE_ID                   ]
    ['0x04' DEVICE_ROLE                 ]
    ['0x05' DEVICE_OPTIONS              ]
    ['0x07' DEVICE_INSTANCE             ]
]

[enum uint 8 'IpSubOption'
    ['0x02' IP_PARAMETER                ]
]

[enum uint 8 'DCPServiceID'
    ['0x05' IDENTIFY]
    ['0xFEFE' IDENTIFY_RESPONSE         ]
]

[enum uint 8 'DCPServiceType'
    ['0x00' REQUEST                     ]
    ['0x01' RESPONSE_SUCCESS            ]
]

[enum uint 16 'FrameType'
    ['0xFEFE' IDENTIFY_MULTICAST_REQUEST]
    ['0xFEFF' IDENTIFY_RESPONSE         ]
]

[type 'MacAddress'
    [simple     uint        8   'octet1'            ]
    [simple     uint        8   'octet2'            ]
    [simple     uint        8   'octet3'            ]
    [simple     uint        8   'octet4'            ]
    [simple     uint        8   'octet5'            ]
    [simple     uint        8   'octet6'            ]
]

[type 'IPv4Address'
    [simple     uint        8   'octet1'            ]
    [simple     uint        8   'octet2'            ]
    [simple     uint        8   'octet3'            ]
    [simple     uint        8   'octet4'            ]
]

[type 'Text' [uint 16  'len']
    // fixme _type.encoding doesn't work for readString call, it ends up with &quot; sequence in generated Java
    [manual string 'UTF-8' 'text'
        'STATIC_CALL("org.apache.plc4x.java.profinet.dcp.DCPUtil.readString", io, len)'
        'STATIC_CALL("org.apache.plc4x.java.profinet.dcp.DCPUtil.writeString", io, _value, _type.encoding)'
        'STATIC_CALL("org.apache.plc4x.java.profinet.dcp.DCPUtil.length", text)'
    ]
]