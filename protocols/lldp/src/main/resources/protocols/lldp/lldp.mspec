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

// LLDP stands for Link Layer Discovery Protocol, it is a layer over ethernet frame which allows to communicate state of
// network enabled devices, but not only.
// Some parts of standard are available online:
// https://standards.ieee.org/standard/802_1AB-2016.html

[type 'LLDP'
    [array TLV 'tlv' terminated 'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.hasMore", io)']
]

[type 'TLV'
    [discriminator     uint 7  'tlvType'   ]
    [typeSwitch 'tlvType'
        ['0x00' EndTLV
            [reserved uint 9 '0x0'                 ]
        ]
        ['0x01' ChassisIdTLV
            [simple         uint 9        'len'                      ]
            [enum           ChassisIdType 'chassisIdType'            ]
            [simple         Text          'chassisId'     ['len - 1']]
        ]
        ['0x02' PortIdTLV
            [simple         uint 9        'len'                      ]
            [enum           PortIdType    'portIdType'               ]
            [simple         Text          'portId'        ['len - 1']]
        ]
        ['0x03' TimeToLiveTLV
            [simple         uint 9        'len'                      ]
            [simple         TTL           'ttl'               ['len']]
        ]
        ['0x04' PortDescriptionTLV
            [simple         uint 9        'len'                      ]
            [simple         Text          'description'       ['len']]
        ]
        ['0x05' SystemNameTLV
            [simple         uint 9        'len'                      ]
            [simple         Text          'name'              ['len']]
        ]
        ['0x06' SystemDescriptionTLV
            [simple         uint 9        'len'                      ]
            [simple         Text          'description'       ['len']]
        ]
        ['0x07' SystemCapabilitiesTLV
            [simple SystemCapabilitiesValueBlock 'systemCapabilities']
            [simple SystemCapabilitiesValueBlock 'enabledCapabilities']
        ]
        ['0x08' ManagementAddressTLV
            [simple         uint 9        'len'                      ]
            [simple         uint 8        'addressLength'            ]
            [simple         uint 8        'addressSubType'           ]
            [simple         IPv4Address   'ip'                       ]
            [simple         uint 8        'interfaceType'            ]
            [simple         uint 32       'interfaceNumber'          ]
            [simple         uint 8        'oid'                      ]
            [simple         Value         'objectIdentifier'  ['oid']]
        ]
        ['0x7F' CustomTLV
            [simple         uint 9      'len'      ]
            [simple CustomTLVBlock 'custom' ['len']]
        ]
        [UnknownTLV
            [simple uint 9 'len'                ]
            [array   int 8 'data' length 'len'  ]
        ]
    ]
]

[discriminatedType 'CustomTLVBlock' [uint 9 'len']
    [discriminator     uint 24  'oui'                     ]
    [typeSwitch 'oui'
        ['0x000ECF' ProfinetTLV
            [simple ProfinetDataBlock 'profinet' ['len']]
        ]
        [UnknownCustomSubTLV
            [simple  uint 8   'subType'                  ]
            [simple   Value   'data'    ['len - 4' ]]
        ]
    ]
]

[type 'ProfinetDataBlock' [uint 9 'len']
    [discriminator     uint 8   'subType'                ]
    [typeSwitch 'subType'
        ['0x02' ProfinetPort
            [enum            RTClass2PortStatus    'rtClass2'      ]
            [enum            RTClass3PortStatus    'rtClass3'      ]
            // there is some additional section, but can't determine available fields, it is 3 bits from rtClass3 status
            // plus additional 13 bits to form two bytes.
            [reserved uint 13 '0x00'                               ]
        ]
        ['0x05' ProfinetChassisMac
            [simple          MacAddress               'mac'        ]
        ]
        [UnknownProfinetDataBlock
            [simple          Value         'data'  ['len - 4' ]]
        ]
    ]
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

[type 'SystemCapabilitiesValueBlock'
    [simple bit 'other'             ]
    [simple bit 'repeater'          ]
    [simple bit 'bridge'            ]
    [simple bit 'wlanAccessPoint'   ]
    [simple bit 'router'            ]
    [simple bit 'telephone'         ]
    [simple bit 'docsisCableDevice' ]
    [simple bit 'stationOnly'       ]
    [simple bit 'cVLANComponent'    ]
    [simple bit 'sVLANComponent'    ]
    [simple bit 'twoPortMACRelay'   ]
]

[enum uint 8          'ChassisIdType'
    ['0x01' CHASSIS_COMPONENT        ]
    ['0x02' INTERFACE_ALIAS          ]
    ['0x03' PORT_COMPONENT           ]
    ['0x04' MAC_ADDRESS              ]
    ['0x05' NETWORK_ADDRESS          ]
    ['0x06' INTERFACE_NAME           ]
    ['0x07' LOCAL                    ]
]

[enum uint 8          'PortIdType'
    ['0x01' INTERFACE_ALIAS          ]
    ['0x01' PORT_COMPONENT           ]
    ['0x03' MAC_ADDRESS              ]
    ['0x04' NETWORK_ADDRESS          ]
    ['0x05' INTERFACE_NAME           ]
    ['0x06' AGENT_CIRCUIT_ID         ]
    ['0x07' LOCAL                    ]
]


[enum uint 16         'RTClass2PortStatus'
    ['0x0000' OFF]
    ['0x0001' ON ]
]

[enum uint 3           'RTClass3PortStatus'
    ['0x00' OFF      ]
    ['0x01' RESERVED ]
    ['0x02' UP       ]
    ['0x03' DOWN     ]
    ['0x04' RUN      ]
]

[type 'ValueBlock'
    [simple            uint 9   'len']
    [array uint 8 'data' length 'len']
]

[type 'Value' [uint 9  'len']
    [array int 8 'data' length 'len']
]

[type 'Text' [uint 9  'len']
    // fixme _type.encoding doesn't work for readString call, it ends up with &quot; sequence in generated Java
    [manual string 'UTF-8' 'text'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.readString", io, len)'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.writeString", io, _value, _type.encoding)'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.length", text)'
    ]
]

[type 'TTL' [uint 9  'len']
    // fixme Use a duration type instead!
    [manual time 'value'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.readTime", io, len)'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.writeTime", io, _value)'
        'STATIC_CALL("org.apache.plc4x.java.lldp.LLDPUtil.byteLength", value)'
    ]
]