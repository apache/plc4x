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

[type CBusConstants
    [const    uint 16     cbusTcpDefaultPort 10001]
]

// TODO: check if that can be used in combination with srchk
[type CBusOptions
    [simple bit connect]
    [simple bit smart  ]
    [simple bit idmon  ]
    [simple bit exstat ]
    [simple bit monitor]
    [simple bit monall ]
    [simple bit pun    ]
    [simple bit pcn    ]
]

[type CBusMessage(bit response, bit srchk, uint 16 messageLength)
    [typeSwitch response
       ['false' *ToServer
            [simple   Request('srchk', 'messageLength')         request         ]
       ]
       ['true' *ToClient
            [simple   Reply('messageLength')                    reply           ]
       ]
    ]
]

[type Request(bit srchk, uint 16 messageLength)
    [peek    RequestType peekedByte                                         ]
    [virtual uint 16 payloadLength 'messageLength-2'                        ] // We substract the command itself and the termination
    [typeSwitch peekedByte
        ['SMART_CONNECT_SHORTCUT' *SmartConnectShortcut
            [const    byte                pipe      0x7C                    ]
        ]
        ['RESET' *Reset
            [const    byte                tilde     0x7E                    ]
        ]
        ['DIRECT_COMMAND' *DirectCommandAccess(uint 16 payloadLength)
            [const    byte                at        0x40                    ]
            // Usually you would read the command now here but we need to decode ascii first
            //[simple   CBusCommandPointToPoint('srchk')     cbusCommand    ]
            [manual   CBusCommand
                                          cbusCommand
                        'STATIC_CALL("readCBusCommand", readBuffer, payloadLength, srchk)'
                        'STATIC_CALL("writeCBusCommand", writeBuffer, cbusCommand)'
                        '_value.lengthInBytes*2'                                     ]
        ]
        ['REQUEST_COMMAND' *Command(uint 16 payloadLength)
            [const    byte                initiator 0x5C                    ] // 0x5C == "/"
            // Usually you would read the command now here but we need to decode ascii first
            //[simple   CBusCommand('srchk')     cbusCommand                ]
            [manual   CBusCommand
                                          cbusCommand
                        'STATIC_CALL("readCBusCommand", readBuffer, payloadLength, srchk)'
                        'STATIC_CALL("writeCBusCommand", writeBuffer, cbusCommand)'
                        '_value.lengthInBytes*2'                                     ]
        ]
        ['NULL' *Null
            [const    uint 32             nullIndicator        0x6E756C6C   ] // "null"
        ]
        ['EMPTY' *Empty
        ]
    ]
    [simple   RequestTermination  termination                               ]
]

[enum uint 8 RequestType(uint 8 controlChar)
    ['0x00' UNKNOWN                 ['0x00']]
    ['0x7C' SMART_CONNECT_SHORTCUT  ['0x7C']] // control char = '|'
    ['0x7E' RESET                   ['0x7E']] // control char = '~'
    ['0x40' DIRECT_COMMAND          ['0x40']] // control char = '@'
    ['0x5C' REQUEST_COMMAND         ['0x5C']] // control char = '/'
    ['0x6E' NULL                    ['0x00']] // null doesn't have a "control char" so we just consume the rest
    ['0x0D' EMPTY                   ['0x00']] // empty doesn't have a "control char" so we just consume the rest
]

[discriminatedType CBusCommand(bit srchk)
    [simple  CBusHeader header           ]
    [virtual bit        isDeviceManagement 'header.dp']
    // TODO: header.destinationAddressType could be used directly but for this we need source type resolving to work (WIP)
    [virtual DestinationAddressType destinationAddressType 'header.destinationAddressType']
    [typeSwitch destinationAddressType, isDeviceManagement
        [*, 'true' *DeviceManagement
            [simple     uint 8  parameterNumber                         ]
            [const      byte    delimiter       0x0                     ]
            [simple     byte    parameterValue                          ]
        ]
        ['PointToPointToMultiPoint' *PointToPointToMultiPoint
            [simple CBusPointToPointToMultipointCommand('srchk') command]
        ]
        ['PointToMultiPoint'        *PointToMultiPoint
            [simple CBusPointToMultiPointCommand('srchk')        command]
        ]
        ['PointToPoint'             *PointToPoint
            [simple CBusPointToPointCommand('srchk')             command]
        ]
    ]
]

[type CBusHeader
    [simple   PriorityClass          priorityClass         ]
    [simple   bit                    dp                    ] // Reserved for internal C-Bus management purposes (Referred to as special packet attribute)
    [simple   uint 2                 rc                    ] // Reserved for internal C-Bus management purposes (Referred to as special packet attribute)
    [simple   DestinationAddressType destinationAddressType]
]

[enum uint 2 PriorityClass
    ['0x00' Class4] // lowest
    ['0x01' Class3] // medium low
    ['0x02' Class2] // medium high
    ['0x03' Class1] // highest
]

[enum uint 3 DestinationAddressType
    ['0x03' PointToPointToMultiPoint] // P-P-M
    ['0x05' PointToMultiPoint       ] // P-M
    ['0x06' PointToPoint            ] // P-P
]

[type UnitAddress
    [simple byte address]
]

[type BridgeAddress
    [simple byte address]
]

[type SerialInterfaceAddress
    [simple byte address]
]

[type Alpha
    [simple byte character]
    [validation '(character >= 0x67) && (character <= 0x7A)' "character not in alpha space" shouldFail=false] // Read if the peeked byte is between 'g' and 'z'
]

[type NetworkRoute
    [reserved uint 2      '0x00'                                                       ]
    [simple RouteType     reverseRouteType                                             ]
    [simple RouteType     routeType                                                    ]
    [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
]

// The last 3 bits are the total number of bridges ... subtracting 1 results in the number of additional bridges.
// It also seems as if these are generally 2 empty bits and then twice the number of bridges as 3-bit numbers.
// The block of 3 in bit's 3..5 seem to be the reverse route.
//
// Observations on failing packets:
// - In the first case the first two bits are not empty, but 01 ... the first block of bridges is then increased by one.
// - In another packet the first two bits are 0, but the first group ist set to 111 and the number of bridges is set to 000.
// - In another packet the first two bits are 0, the first group is set to 001 and the second to 101
[enum uint 3 RouteType(uint 3 additionalBridges)
    ['0x0' NoBridgeAtAll         ['0']]
    ['0x1' NoAdditionalBridge    ['0']]
    ['0x2' OneAdditionalBridge   ['1']]
    ['0x3' TwoAdditionalBridge   ['2']]
    ['0x4' ThreeAdditionalBridge ['3']]
    ['0x5' FourAdditionalBridge  ['4']]
    ['0x6' FiveAdditionalBridge  ['5']]
    ['0x7' SixAdditionalBridge   ['6']]
]

[discriminatedType CBusPointToPointCommand(bit srchk)
    [peek    uint 16     bridgeAddressCountPeek ]
    [virtual bit         isDirect  '(bridgeAddressCountPeek & 0x00FF) == 0x0000']
    [typeSwitch isDirect
        ['true'  CBusPointToPointCommandDirect
            [simple   UnitAddress   unitAddress                                                 ]
            [reserved uint 8        '0x00'                                                      ]
        ]
        ['false' CBusPointToPointCommandIndirect
            [simple BridgeAddress bridgeAddress                                                 ]
            [simple NetworkRoute  networkRoute                                                  ]
            [simple UnitAddress   unitAddress                                                   ]
        ]
    ]
    [simple   CALData calData                                                                   ]
    [optional Checksum      crc      'srchk'                                                    ] // checksum is optional but mspec checksum isn't
    [optional Alpha         alpha                                                               ]
]

[discriminatedType CBusPointToMultiPointCommand(bit srchk)
    [peek    byte     peekedApplication                                                                ]
    [typeSwitch peekedApplication
        ['0xFF'   CBusPointToMultiPointCommandStatus
            [reserved byte          '0xFF'                                                             ]
            [reserved byte          '0x00'                                                             ]
            [simple   StatusRequest statusRequest                                                      ]
            [optional Checksum      crc           'srchk'                                              ] // checksum is optional but mspec checksum isn't
            [optional Alpha         alpha                                                              ]
        ]
        [         CBusPointToMultiPointCommandNormal
            [simple   ApplicationIdContainer   application                                             ]
            [reserved byte                     '0x00'                                                  ]
            [simple   SALData                  salData                                                 ]
            [optional Checksum                 crc         'srchk'                                     ] // crc      is optional but mspec crc      isn't
            [optional Alpha         alpha                                                               ]
        ]
    ]
]

[discriminatedType CBusPointToPointToMultipointCommand(bit srchk)
    [simple BridgeAddress bridgeAddress                                                              ]
    [simple NetworkRoute  networkRoute                                                               ]
    [peek    byte       peekedApplication                                                            ]
    [typeSwitch peekedApplication
        ['0xFF'   CBusCommandPointToPointToMultiPointStatus
            [reserved byte        '0xFF'                                                             ]
            [simple StatusRequest statusRequest                                                      ]
            [optional Checksum    crc           'srchk'                                              ] // crc      is optional but mspec crc      isn't
            [optional Alpha         alpha                                                            ]
        ]
        [         CBusCommandPointToPointToMultiPointNormal
            [simple   ApplicationIdContainer application                                             ]
            [simple   SALData                salData                                                 ]
            [optional Checksum               crc         'srchk'                                     ] // crc      is optional but mspec crc      isn't
            [optional Alpha         alpha                                                            ]
        ]
    ]
]

/*
    Application ID Ranges 
    From https://updates.clipsal.com/ClipsalSoftwareDownload/mainsite/cis/__data/page/4129/ApplicationIDNumbers.pdf
*/
[enum uint 8 ApplicationId
    ['0x00' RESERVED                          ]
    ['0x01' FREE_USAGE                        ]
    ['0x02' TEMPERATURE_BROADCAST             ]
    ['0x03' ROOM_CONTROL_SYSTEM               ]
    ['0x04' LIGHTING                          ]
    ['0x05' VENTILATION                       ]
    ['0x06' IRRIGATION_CONTROL                ]
    ['0x07' POOLS_SPAS_PONDS_FOUNTAINS_CONTROL]
    ['0x08' HEATING                           ]
    ['0x09' AIR_CONDITIONING                  ]
    ['0x0A' TRIGGER_CONTROL                   ]
    ['0x0B' ENABLE_CONTROL                    ]
    ['0x0C' AUDIO_AND_VIDEO                   ]
    ['0x0D' SECURITY                          ]
    ['0x0E' METERING                          ]
    ['0x0F' ACCESS_CONTROL                    ]
    ['0x10' CLOCK_AND_TIMEKEEPING             ]
    ['0x11' TELEPHONY_STATUS_AND_CONTROL      ]
    ['0x12' MEASUREMENT                       ]
    ['0x13' TESTING                           ]
]

[enum uint 4 LightingCompatible
    ['0x0' NO                  ]
    ['0x1' YES                 ]
    ['0x2' YES_BUT_RESTRICTIONS]
    ['0x3' NA                  ]
]

[enum uint 8 ApplicationIdContainer(ApplicationId applicationId,                         LightingCompatible lightingCompatible)
    ['0x00' RESERVED_00                           ['RESERVED'                          , 'NA'                  ]]
    ['0x01' FREE_USAGE_01                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x02' FREE_USAGE_02                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x03' FREE_USAGE_03                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x04' FREE_USAGE_04                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x05' FREE_USAGE_05                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x06' FREE_USAGE_06                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x07' FREE_USAGE_07                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x08' FREE_USAGE_08                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x09' FREE_USAGE_09                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0A' FREE_USAGE_0A                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0B' FREE_USAGE_0B                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0C' FREE_USAGE_0C                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0D' FREE_USAGE_0D                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0E' FREE_USAGE_0E                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x0F' FREE_USAGE_0F                         ['FREE_USAGE'                        , 'NA'                  ]]
    ['0x10' RESERVED_10                           ['RESERVED'                          , 'NA'                  ]] // Actually the spec says: 0x1F-0x18but it wouldn't make sense that way.
    ['0x11' RESERVED_11                           ['RESERVED'                          , 'NA'                  ]]
    ['0x12' RESERVED_12                           ['RESERVED'                          , 'NA'                  ]]
    ['0x13' RESERVED_13                           ['RESERVED'                          , 'NA'                  ]]
    ['0x14' RESERVED_14                           ['RESERVED'                          , 'NA'                  ]]
    ['0x15' RESERVED_15                           ['RESERVED'                          , 'NA'                  ]]
    ['0x16' RESERVED_16                           ['RESERVED'                          , 'NA'                  ]]
    ['0x17' RESERVED_17                           ['RESERVED'                          , 'NA'                  ]]
    ['0x18' RESERVED_18                           ['RESERVED'                          , 'NA'                  ]]
    ['0x19' TEMPERATURE_BROADCAST_19              ['TEMPERATURE_BROADCAST'             , 'NO'                  ]]
    ['0x1A' RESERVED_1A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x1B' RESERVED_1B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x1C' RESERVED_1C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x1D' RESERVED_1D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x1E' RESERVED_1E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x1F' RESERVED_1F                           ['RESERVED'                          , 'NA'                  ]]
    ['0x20' RESERVED_20                           ['RESERVED'                          , 'NA'                  ]]
    ['0x21' RESERVED_21                           ['RESERVED'                          , 'NA'                  ]]
    ['0x22' RESERVED_22                           ['RESERVED'                          , 'NA'                  ]]
    ['0x23' RESERVED_23                           ['RESERVED'                          , 'NA'                  ]]
    ['0x24' RESERVED_24                           ['RESERVED'                          , 'NA'                  ]]
    ['0x25' RESERVED_25                           ['RESERVED'                          , 'NA'                  ]]
    ['0x26' ROOM_CONTROL_SYSTEM_26                ['ROOM_CONTROL_SYSTEM'               , 'YES'                 ]]
    ['0x27' RESERVED_27                           ['RESERVED'                          , 'NA'                  ]]
    ['0x28' RESERVED_28                           ['RESERVED'                          , 'NA'                  ]]
    ['0x29' RESERVED_29                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2A' RESERVED_2A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2B' RESERVED_2B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2C' RESERVED_2C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2D' RESERVED_2D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2E' RESERVED_2E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x2F' RESERVED_2F                           ['RESERVED'                          , 'NA'                  ]]
    ['0x30' LIGHTING_30                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x31' LIGHTING_31                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x32' LIGHTING_32                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x33' LIGHTING_33                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x34' LIGHTING_34                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x35' LIGHTING_35                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x36' LIGHTING_36                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x37' LIGHTING_37                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x38' LIGHTING_38                           ['LIGHTING'                          , 'YES'                 ]] // Default
    ['0x39' LIGHTING_39                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3A' LIGHTING_3A                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3B' LIGHTING_3B                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3C' LIGHTING_3C                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3D' LIGHTING_3D                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3E' LIGHTING_3E                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x3F' LIGHTING_3F                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x40' LIGHTING_40                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x41' LIGHTING_41                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x42' LIGHTING_42                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x43' LIGHTING_43                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x44' LIGHTING_44                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x45' LIGHTING_45                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x46' LIGHTING_46                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x47' LIGHTING_47                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x48' LIGHTING_48                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x49' LIGHTING_49                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4A' LIGHTING_4A                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4B' LIGHTING_4B                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4C' LIGHTING_4C                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4D' LIGHTING_4D                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4E' LIGHTING_4E                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x4F' LIGHTING_4F                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x50' LIGHTING_50                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x51' LIGHTING_51                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x52' LIGHTING_52                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x53' LIGHTING_53                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x54' LIGHTING_54                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x55' LIGHTING_55                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x56' LIGHTING_56                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x57' LIGHTING_57                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x58' LIGHTING_58                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x59' LIGHTING_59                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5A' LIGHTING_5A                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5B' LIGHTING_5B                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5C' LIGHTING_5C                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5D' LIGHTING_5D                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5E' LIGHTING_5E                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x5F' LIGHTING_5F                           ['LIGHTING'                          , 'YES'                 ]]
    ['0x60' RESERVED_60                           ['RESERVED'                          , 'NA'                  ]]
    ['0x61' RESERVED_61                           ['RESERVED'                          , 'NA'                  ]]
    ['0x62' RESERVED_62                           ['RESERVED'                          , 'NA'                  ]]
    ['0x63' RESERVED_63                           ['RESERVED'                          , 'NA'                  ]]
    ['0x64' RESERVED_64                           ['RESERVED'                          , 'NA'                  ]]
    ['0x65' RESERVED_65                           ['RESERVED'                          , 'NA'                  ]]
    ['0x66' RESERVED_66                           ['RESERVED'                          , 'NA'                  ]]
    ['0x67' RESERVED_67                           ['RESERVED'                          , 'NA'                  ]]
    ['0x68' RESERVED_68                           ['RESERVED'                          , 'NA'                  ]]
    ['0x69' RESERVED_69                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6A' RESERVED_6A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6B' RESERVED_6B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6C' RESERVED_6C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6D' RESERVED_6D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6E' RESERVED_6E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x6F' RESERVED_6F                           ['RESERVED'                          , 'NA'                  ]]
    ['0x70' VENTILATION_70                        ['VENTILATION'                       , 'YES'                 ]]
    ['0x71' IRRIGATION_CONTROL_71                 ['IRRIGATION_CONTROL'                , 'YES'                 ]]
    ['0x72' POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72 ['POOLS_SPAS_PONDS_FOUNTAINS_CONTROL', 'YES'                 ]]
    ['0x73' RESERVED_73                           ['RESERVED'                          , 'NA'                  ]] // HVAC_ACTUATOR
    ['0x74' RESERVED_74                           ['RESERVED'                          , 'NA'                  ]] // HVAC_ACTUATOR
    ['0x75' RESERVED_75                           ['RESERVED'                          , 'NA'                  ]]
    ['0x76' RESERVED_76                           ['RESERVED'                          , 'NA'                  ]]
    ['0x77' RESERVED_77                           ['RESERVED'                          , 'NA'                  ]]
    ['0x78' RESERVED_78                           ['RESERVED'                          , 'NA'                  ]]
    ['0x79' RESERVED_79                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7A' RESERVED_7A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7B' RESERVED_7B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7C' RESERVED_7C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7D' RESERVED_7D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7E' RESERVED_7E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x7F' RESERVED_7F                           ['RESERVED'                          , 'NA'                  ]]
    ['0x80' RESERVED_80                           ['RESERVED'                          , 'NA'                  ]]
    ['0x81' RESERVED_81                           ['RESERVED'                          , 'NA'                  ]]
    ['0x82' RESERVED_82                           ['RESERVED'                          , 'NA'                  ]]
    ['0x83' RESERVED_83                           ['RESERVED'                          , 'NA'                  ]]
    ['0x84' RESERVED_84                           ['RESERVED'                          , 'NA'                  ]]
    ['0x85' RESERVED_85                           ['RESERVED'                          , 'NA'                  ]]
    ['0x86' RESERVED_86                           ['RESERVED'                          , 'NA'                  ]]
    ['0x87' RESERVED_87                           ['RESERVED'                          , 'NA'                  ]]
    ['0x88' HEATING_88                            ['HEATING'                           , 'YES'                 ]]
    ['0x89' RESERVED_89                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8A' RESERVED_8A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8B' RESERVED_8B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8C' RESERVED_8C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8D' RESERVED_8D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8E' RESERVED_8E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x8F' RESERVED_8F                           ['RESERVED'                          , 'NA'                  ]]
    ['0x90' RESERVED_90                           ['RESERVED'                          , 'NA'                  ]]
    ['0x91' RESERVED_91                           ['RESERVED'                          , 'NA'                  ]]
    ['0x92' RESERVED_92                           ['RESERVED'                          , 'NA'                  ]]
    ['0x93' RESERVED_93                           ['RESERVED'                          , 'NA'                  ]]
    ['0x94' RESERVED_94                           ['RESERVED'                          , 'NA'                  ]]
    ['0x95' RESERVED_95                           ['RESERVED'                          , 'NA'                  ]]
    ['0x96' RESERVED_96                           ['RESERVED'                          , 'NA'                  ]]
    ['0x97' RESERVED_97                           ['RESERVED'                          , 'NA'                  ]]
    ['0x98' RESERVED_98                           ['RESERVED'                          , 'NA'                  ]]
    ['0x99' RESERVED_99                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9A' RESERVED_9A                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9B' RESERVED_9B                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9C' RESERVED_9C                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9D' RESERVED_9D                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9E' RESERVED_9E                           ['RESERVED'                          , 'NA'                  ]]
    ['0x9F' RESERVED_9F                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA0' RESERVED_A0                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA1' RESERVED_A1                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA2' RESERVED_A2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA3' RESERVED_A3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA4' RESERVED_A4                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA5' RESERVED_A5                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA6' RESERVED_A6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA7' RESERVED_A7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA8' RESERVED_A8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xA9' RESERVED_A9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xAA' RESERVED_AA                           ['RESERVED'                          , 'NA'                  ]]
    ['0xAB' RESERVED_AB                           ['RESERVED'                          , 'NA'                  ]]
    ['0xAC' AIR_CONDITIONING_AC                   ['AIR_CONDITIONING'                  , 'NO'                  ]]
    ['0xAD' RESERVED_AD                           ['RESERVED'                          , 'NA'                  ]] // INFO_MESSAGES
    ['0xAE' RESERVED_AE                           ['RESERVED'                          , 'NA'                  ]]
    ['0xAF' RESERVED_AF                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB0' RESERVED_B0                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB1' RESERVED_B1                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB2' RESERVED_B2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB3' RESERVED_B3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB4' RESERVED_B4                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB5' RESERVED_B5                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB6' RESERVED_B6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB7' RESERVED_B7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB8' RESERVED_B8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xB9' RESERVED_B9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBA' RESERVED_BA                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBB' RESERVED_BB                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBC' RESERVED_BC                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBD' RESERVED_BD                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBE' RESERVED_BE                           ['RESERVED'                          , 'NA'                  ]]
    ['0xBF' RESERVED_BF                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC0' MEDIA_TRANSPORT_C0                    ['RESERVED'                          , 'NA'                  ]] // MEDIA_TRANSPORT
    ['0xC1' RESERVED_C1                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC2' RESERVED_C2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC3' RESERVED_C3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC4' RESERVED_C4                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC5' RESERVED_C5                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC6' RESERVED_C6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC7' RESERVED_C7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC8' RESERVED_C8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xC9' RESERVED_C9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xCA' TRIGGER_CONTROL_CA                    ['TRIGGER_CONTROL'                   , 'YES_BUT_RESTRICTIONS']]
    ['0xCB' ENABLE_CONTROL_CB                     ['ENABLE_CONTROL'                    , 'YES_BUT_RESTRICTIONS']]
    ['0xCC' I_HAVE_NO_IDEA_CC                     ['RESERVED'                          , 'NA'                  ]] // This is the only value actually not defined in the spec.
    ['0xCD' AUDIO_AND_VIDEO_CD                    ['AUDIO_AND_VIDEO'                   , 'YES_BUT_RESTRICTIONS']]
    ['0xCE' ERROR_REPORTING_CE                    ['RESERVED'                          , 'NA'                  ]] // ERROR_REPORTING
    ['0xCF' RESERVED_CF                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD0' SECURITY_D0                           ['SECURITY'                          , 'NO'                  ]]
    ['0xD1' METERING_D1                           ['METERING'                          , 'NO'                  ]]
    ['0xD2' RESERVED_D2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD3' RESERVED_D3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD4' RESERVED_D4                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD5' ACCESS_CONTROL_D5                     ['ACCESS_CONTROL'                    , 'NO'                  ]]
    ['0xD6' RESERVED_D6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD7' RESERVED_D7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD8' RESERVED_D8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xD9' RESERVED_D9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDA' RESERVED_DA                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDB' RESERVED_DB                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDC' RESERVED_DC                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDD' RESERVED_DD                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDE' RESERVED_DE                           ['RESERVED'                          , 'NA'                  ]]
    ['0xDF' CLOCK_AND_TIMEKEEPING_DF              ['CLOCK_AND_TIMEKEEPING'             , 'NO'                  ]]
    ['0xE0' TELEPHONY_STATUS_AND_CONTROL_E0       ['TELEPHONY_STATUS_AND_CONTROL'      , 'NO'                  ]]
    ['0xE1' RESERVED_E1                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE2' RESERVED_E2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE3' RESERVED_E3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE4' MEASUREMENT_E4                        ['MEASUREMENT'                       , 'NO'                  ]]
    ['0xE5' RESERVED_E5                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE6' RESERVED_E6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE7' RESERVED_E7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE8' RESERVED_E8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xE9' RESERVED_E9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xEA' RESERVED_EA                           ['RESERVED'                          , 'NA'                  ]]
    ['0xEB' RESERVED_EB                           ['RESERVED'                          , 'NA'                  ]]
    ['0xEC' RESERVED_EC                           ['RESERVED'                          , 'NA'                  ]]
    ['0xED' RESERVED_ED                           ['RESERVED'                          , 'NA'                  ]]
    ['0xEE' RESERVED_EE                           ['RESERVED'                          , 'NA'                  ]]
    ['0xEF' RESERVED_EF                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF0' RESERVED_F0                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF1' RESERVED_F1                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF2' RESERVED_F2                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF3' RESERVED_F3                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF4' RESERVED_F4                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF5' RESERVED_F5                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF6' RESERVED_F6                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF7' RESERVED_F7                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF8' RESERVED_F8                           ['RESERVED'                          , 'NA'                  ]]
    ['0xF9' RESERVED_F9                           ['RESERVED'                          , 'NA'                  ]]
    ['0xFA' TESTING_FA                            ['TESTING'                           , 'NA'                  ]]
    ['0xFB' RESERVED_FB                           ['RESERVED'                          , 'NO'                  ]]
    ['0xFC' RESERVED_FC                           ['RESERVED'                          , 'NO'                  ]]
    ['0xFD' RESERVED_FD                           ['RESERVED'                          , 'NO'                  ]]
    ['0xFE' RESERVED_FE                           ['RESERVED'                          , 'NO'                  ]]
    ['0xFF' RESERVED_FF                           ['RESERVED'                          , 'NO'                  ]] // NETWORK_CONTROL
]

[type CALData
    [simple  CALCommandTypeContainer commandTypeContainer                                   ]
    //TODO: golang doesn't like checking for 0
    //[validation 'commandTypeContainer!=null' "no command type could be found"               ]
    [virtual CALCommandType          commandType          'commandTypeContainer.commandType']
    [typeSwitch commandType
        ['RESET' CALDataRequestReset
        ]
        ['RECALL' CALDataRequestRecall
            [simple uint 8 paramNo                                                          ]
            [simple uint 8 count                                                            ]
        ]
        ['IDENTIFY' CALDataRequestIdentify
            [simple Attribute attribute                                                     ]
        ]
        ['GET_STATUS' CALDataRequestGetStatus
            [simple uint 8 paramNo                                                          ]
            [simple uint 8 count                                                            ]
        ]
        ['REPLY' CALDataReplyReply(CALCommandTypeContainer commandTypeContainer)
            [simple uint 8 paramNumber                                                      ]
            [array  byte   data        count 'commandTypeContainer.numBytes'                ]
        ]
        ['ACKNOWLEDGE' CALDataReplyAcknowledge
            [simple uint 8 paramNo                                                          ]
            [simple uint 8 code                                                             ]
        ]
        ['STATUS' CALDataReplyStatus(CALCommandTypeContainer commandTypeContainer)
            [simple ApplicationIdContainer application                                                 ]
            [simple uint 8                 blockStart                                                  ]
            [array  byte                   data        count 'commandTypeContainer.numBytes'           ]
        ]
        ['STATUS_EXTENDED' CALDataReplyStatusExtended(CALCommandTypeContainer commandTypeContainer)
            [simple uint 8                 encoding                                                    ]
            [simple ApplicationIdContainer application                                                 ]
            [simple uint 8                 blockStart                                                  ]
            [array  byte                   data        count 'commandTypeContainer.numBytes'           ]
        ]
    ]
]

[enum uint 8 Attribute(uint 8 bytesReturned)
    ['0x00' Manufacturer              [ '8']]
    ['0x01' Type                      [ '8']]
    ['0x02' FirmwareVersion           [ '8']]
    ['0x03' Summary                   [ '9']]
    ['0x04' ExtendedDiagnosticSummary ['13']]
    ['0x05' NetworkTerminalLevels     ['13']]
    ['0x06' TerminalLevel             ['13']]
    ['0x07' NetworkVoltage            [ '5']]
    ['0x08' GAVValuesCurrent          ['16']]
    ['0x09' GAVValuesStored           ['16']]
    ['0x0A' GAVPhysicalAddresses      ['16']]
    ['0x0B' LogicalAssignment         ['13']]
    ['0x0C' Delays                    ['14']]
    ['0x0D' MinimumLevels             ['13']]
    ['0x0E' MaximumLevels             ['13']]
    ['0x0F' CurrentSenseLevels        [ '8']]
    ['0x10' OutputUnitSummary         [ '4']]
    ['0x11' DSIStatus                 ['10']]
]

[type IdentifyReplyCommand(Attribute attribute)
    [typeSwitch attribute
        ['Manufacturer'                 IdentifyReplyCommandManufacturer
            [simple string 64  manufacturerName ]
        ]
        ['Type'                         IdentifyReplyCommandType
            [simple string 64  unitType         ]
        ]
        ['FirmwareVersion'              IdentifyReplyCommandFirmwareVersion
            [simple string 64  firmwareVersion  ]
        ]
        ['Summary'                      IdentifyReplyCommandFirmwareSummary
            [simple string 48  firmwareVersion  ]
            [simple byte       unitServiceType  ]
            [simple string 32  version          ]
        ]
        ['ExtendedDiagnosticSummary'    IdentifyReplyCommandExtendedDiagnosticSummary
            [simple ApplicationIdContainer  lowApplication         ]
            [simple ApplicationIdContainer  highApplication        ]
            [simple byte                    area                   ]
            [simple uint 16                 crc                    ]
            [simple uint 32                 serialNumber           ]
            [simple byte                    networkVoltage         ]
            [simple bit                     outputUnit             ]
            [simple bit                     enableChecksumAlarm    ]
            [reserved uint 1                '0'                    ]
            [reserved uint 1                '0'                    ]
            [reserved uint 1                '0'                    ]
            [simple bit                     networkVoltageMarginal ]
            [simple bit                     networkVoltageLow      ]
            [simple bit                     unitInLearnMode        ]
            [simple bit                     microPowerReset        ]
            [simple bit                     internalStackOverflow  ]
            [simple bit                     commsTxError           ]
            [simple bit                     microReset             ]
            [simple bit                     EEDataError            ]
            [simple bit                     EEChecksumError        ]
            [simple bit                     EEWriteError           ]
            [simple bit                     installationMMIError   ]
        ]
        ['NetworkTerminalLevels'        IdentifyReplyCommandNetworkTerminalLevels
            //TODO: read dynamic
        ]
        ['TerminalLevel'                IdentifyReplyCommandTerminalLevels
            //TODO: read dynamic
        ]
        ['NetworkVoltage'               IdentifyReplyCommandNetworkVoltage
           [simple string 2     volts                   ]
           [const  byte         dot     0x2C            ]
           [simple string 2     voltsDecimalPlace       ]
           [const  byte         v       0x56            ]
        ]
        ['GAVValuesCurrent'             IdentifyReplyCommandGAVValuesCurrent
           [array  byte         values  count   '16'    ] // TODO: check datatype
        ]
        ['GAVValuesStored'              IdentifyReplyCommandGAVValuesStored
           [array  byte         values  count   '16'    ] // TODO: check datatype
        ]
        ['GAVPhysicalAddresses'         IdentifyReplyCommandGAVPhysicalAddresses
           [array  byte         values  count   '16'    ] // TODO: check datatype
        ]
        ['LogicalAssignment'            IdentifyReplyCommandLogicalAssignment
            //TODO: read dynamic
        ]
        ['Delays'                       IdentifyReplyCommandDelays
            //TODO: read dynamic
        ]
        ['MinimumLevels'                IdentifyReplyCommandMinimumLevels
            //TODO: read dynamic
        ]
        ['MaximumLevels'                IdentifyReplyCommandMaximumLevels
            //TODO: read dynamic
        ]
        ['CurrentSenseLevels'           IdentifyReplyCommandCurrentSenseLevels
            //TODO: read dynamic
        ]
        ['OutputUnitSummary'            IdentifyReplyCommandOutputUnitSummary
            //TODO: read dynamic
        ]
        ['DSIStatus'                    IdentifyReplyCommandDSIStatus
            [simple ChannelStatus   channelStatus1          ]
            [simple ChannelStatus   channelStatus2          ]
            [simple ChannelStatus   channelStatus3          ]
            [simple ChannelStatus   channelStatus4          ]
            [simple ChannelStatus   channelStatus5          ]
            [simple ChannelStatus   channelStatus6          ]
            [simple ChannelStatus   channelStatus7          ]
            [simple ChannelStatus   channelStatus8          ]
            [simple UnitStatus      unitStatus              ]
            [simple byte            dimmingUCRevisionNumber ]
        ]
    ]
]

[enum uint 8 ChannelStatus
    ['0'    OK                      ]
    ['2'    LAMP_FAULT              ]
    ['3'    CURRENT_LIMIT_OR_SHORT  ]
]

[enum uint 8 UnitStatus
    ['0'    OK                      ]
    ['1'    NACK                    ]
    ['2'    NO_RESPONSE             ]
]

// 1------: Long Form Command (Length is in the 5 least significant bits)
// 0------: Short Form Command (Length is in the 3 least significant bits)
// The invalid packets are receiving a value of 13 / 0x0D -> Short form command: length = 5 (no idea what the bit number 4 means, which is set)
[enum uint 8 CALCommandTypeContainer(CALCommandType commandType, uint 5 numBytes)
    ['0x08' CALCommandReset                  ['RESET',            '0']]
    ['0x1A' CALCommandRecall                 ['RECALL',           '0']]
    ['0x21' CALCommandIdentify               ['IDENTIFY',         '0']]
    ['0x2A' CALCommandGetStatus              ['GET_STATUS',       '0']]
    ['0x81' CALCommandReply_1Bytes           ['REPLY',            '1']]
    ['0x82' CALCommandReply_2Bytes           ['REPLY',            '2']]
    ['0x83' CALCommandReply_3Bytes           ['REPLY',            '3']]
    ['0x84' CALCommandReply_4Bytes           ['REPLY',            '4']]
    ['0x85' CALCommandReply_5Bytes           ['REPLY',            '5']]
    ['0x86' CALCommandReply_6Bytes           ['REPLY',            '6']]
    ['0x87' CALCommandReply_7Bytes           ['REPLY',            '7']]
    ['0x88' CALCommandReply_8Bytes           ['REPLY',            '8']]
    ['0x89' CALCommandReply_9Bytes           ['REPLY',            '9']]
    ['0x8A' CALCommandReply_10Bytes          ['REPLY',           '10']]
    ['0x8B' CALCommandReply_11Bytes          ['REPLY',           '11']]
    ['0x8C' CALCommandReply_12Bytes          ['REPLY',           '12']]
    ['0x8D' CALCommandReply_13Bytes          ['REPLY',           '13']]
    ['0x8E' CALCommandReply_14Bytes          ['REPLY',           '14']]
    ['0x8F' CALCommandReply_15Bytes          ['REPLY',           '15']]
    ['0x90' CALCommandReply_16Bytes          ['REPLY',           '16']]
    ['0x91' CALCommandReply_17Bytes          ['REPLY',           '17']]
    ['0x92' CALCommandReply_18Bytes          ['REPLY',           '18']]
    ['0x93' CALCommandReply_19Bytes          ['REPLY',           '19']]
    ['0x94' CALCommandReply_20Bytes          ['REPLY',           '20']]
    ['0x95' CALCommandReply_21Bytes          ['REPLY',           '21']]
    ['0x96' CALCommandReply_22Bytes          ['REPLY',           '22']]
    ['0x97' CALCommandReply_23Bytes          ['REPLY',           '23']]
    ['0x98' CALCommandReply_24Bytes          ['REPLY',           '24']]
    ['0x99' CALCommandReply_25Bytes          ['REPLY',           '25']]
    ['0x9A' CALCommandReply_26Bytes          ['REPLY',           '26']]
    ['0x9B' CALCommandReply_27Bytes          ['REPLY',           '27']]
    ['0x9C' CALCommandReply_28Bytes          ['REPLY',           '28']]
    ['0x9D' CALCommandReply_29Bytes          ['REPLY',           '29']]
    ['0x9E' CALCommandReply_30Bytes          ['REPLY',           '30']]
    ['0x9F' CALCommandReply_31Bytes          ['REPLY',           '31']]
    ['0x32' CALCommandAcknowledge            ['ACKNOWLEDGE',      '0']]
    ['0xC1' CALCommandStatus_1Bytes          ['STATUS',           '1']]
    ['0xC2' CALCommandStatus_2Bytes          ['STATUS',           '2']]
    ['0xC3' CALCommandStatus_3Bytes          ['STATUS',           '3']]
    ['0xC4' CALCommandStatus_4Bytes          ['STATUS',           '4']]
    ['0xC5' CALCommandStatus_5Bytes          ['STATUS',           '5']]
    ['0xC6' CALCommandStatus_6Bytes          ['STATUS',           '6']]
    ['0xC7' CALCommandStatus_7Bytes          ['STATUS',           '7']]
    ['0xC8' CALCommandStatus_8Bytes          ['STATUS',           '8']]
    ['0xC9' CALCommandStatus_9Bytes          ['STATUS',           '9']]
    ['0xCA' CALCommandStatus_10Bytes         ['STATUS',          '10']]
    ['0xCB' CALCommandStatus_11Bytes         ['STATUS',          '11']]
    ['0xCC' CALCommandStatus_12Bytes         ['STATUS',          '12']]
    ['0xCD' CALCommandStatus_13Bytes         ['STATUS',          '13']]
    ['0xCE' CALCommandStatus_14Bytes         ['STATUS',          '14']]
    ['0xCF' CALCommandStatus_15Bytes         ['STATUS',          '15']]
    ['0xD0' CALCommandStatus_16Bytes         ['STATUS',          '16']]
    ['0xD1' CALCommandStatus_17Bytes         ['STATUS',          '17']]
    ['0xD2' CALCommandStatus_18Bytes         ['STATUS',          '18']]
    ['0xD3' CALCommandStatus_19Bytes         ['STATUS',          '19']]
    ['0xD4' CALCommandStatus_20Bytes         ['STATUS',          '20']]
    ['0xD5' CALCommandStatus_21Bytes         ['STATUS',          '21']]
    ['0xD6' CALCommandStatus_22Bytes         ['STATUS',          '22']]
    ['0xD7' CALCommandStatus_23Bytes         ['STATUS',          '23']]
    ['0xD8' CALCommandStatus_24Bytes         ['STATUS',          '24']]
    ['0xD9' CALCommandStatus_25Bytes         ['STATUS',          '25']]
    ['0xDA' CALCommandStatus_26Bytes         ['STATUS',          '26']]
    ['0xDB' CALCommandStatus_27Bytes         ['STATUS',          '27']]
    ['0xDC' CALCommandStatus_28Bytes         ['STATUS',          '28']]
    ['0xDD' CALCommandStatus_29Bytes         ['STATUS',          '29']]
    ['0xDE' CALCommandStatus_30Bytes         ['STATUS',          '30']]
    ['0xDF' CALCommandStatus_31Bytes         ['STATUS',          '31']]
    ['0xE1' CALCommandStatusExtended_1Bytes  ['STATUS_EXTENDED',  '1']]
    ['0xE1' CALCommandStatusExtended_2Bytes  ['STATUS_EXTENDED',  '2']]
    ['0xE3' CALCommandStatusExtended_3Bytes  ['STATUS_EXTENDED',  '3']]
    ['0xE4' CALCommandStatusExtended_4Bytes  ['STATUS_EXTENDED',  '4']]
    ['0xE5' CALCommandStatusExtended_5Bytes  ['STATUS_EXTENDED',  '5']]
    ['0xE6' CALCommandStatusExtended_6Bytes  ['STATUS_EXTENDED',  '6']]
    ['0xE7' CALCommandStatusExtended_7Bytes  ['STATUS_EXTENDED',  '7']]
    ['0xE8' CALCommandStatusExtended_8Bytes  ['STATUS_EXTENDED',  '8']]
    ['0xE9' CALCommandStatusExtended_9Bytes  ['STATUS_EXTENDED',  '9']]
    ['0xEA' CALCommandStatusExtended_10Bytes ['STATUS_EXTENDED', '10']]
    ['0xEB' CALCommandStatusExtended_11Bytes ['STATUS_EXTENDED', '11']]
    ['0xEC' CALCommandStatusExtended_12Bytes ['STATUS_EXTENDED', '12']]
    ['0xED' CALCommandStatusExtended_13Bytes ['STATUS_EXTENDED', '13']]
    ['0xEE' CALCommandStatusExtended_14Bytes ['STATUS_EXTENDED', '14']]
    ['0xEF' CALCommandStatusExtended_15Bytes ['STATUS_EXTENDED', '15']]
    ['0xF0' CALCommandStatusExtended_16Bytes ['STATUS_EXTENDED', '16']]
    ['0xF1' CALCommandStatusExtended_17Bytes ['STATUS_EXTENDED', '17']]
    ['0xF2' CALCommandStatusExtended_18Bytes ['STATUS_EXTENDED', '18']]
    ['0xF3' CALCommandStatusExtended_19Bytes ['STATUS_EXTENDED', '19']]
    ['0xF4' CALCommandStatusExtended_20Bytes ['STATUS_EXTENDED', '20']]
    ['0xF5' CALCommandStatusExtended_21Bytes ['STATUS_EXTENDED', '21']]
    ['0xF6' CALCommandStatusExtended_22Bytes ['STATUS_EXTENDED', '22']]
    ['0xF7' CALCommandStatusExtended_23Bytes ['STATUS_EXTENDED', '23']]
    ['0xF8' CALCommandStatusExtended_24Bytes ['STATUS_EXTENDED', '24']]
    ['0xF9' CALCommandStatusExtended_25Bytes ['STATUS_EXTENDED', '25']]
    ['0xFA' CALCommandStatusExtended_26Bytes ['STATUS_EXTENDED', '26']]
    ['0xFB' CALCommandStatusExtended_27Bytes ['STATUS_EXTENDED', '27']]
    ['0xFC' CALCommandStatusExtended_28Bytes ['STATUS_EXTENDED', '28']]
    ['0xFD' CALCommandStatusExtended_29Bytes ['STATUS_EXTENDED', '29']]
    ['0xFE' CALCommandStatusExtended_30Bytes ['STATUS_EXTENDED', '30']]
    ['0xFF' CALCommandStatusExtended_31Bytes ['STATUS_EXTENDED', '31']]
]

[enum uint 4 CALCommandType
    // Request
    ['0x0' RESET          ] //00001000
    ['0x0' RECALL         ] //00011010
    ['0x1' IDENTIFY       ] //00100001
    ['0x2' GET_STATUS     ] //01000001
    // Response
    ['0x3' REPLY          ] //100xxxxx
    ['0x4' ACKNOWLEDGE    ] //00110010
    ['0x5' STATUS         ] //110xxxxx
    ['0x5' STATUS_EXTENDED] //111xxxxx
]

[type StatusRequest
    [peek    byte     statusType           ]
    [typeSwitch statusType
        ['0x7A' StatusRequestBinaryState
            [reserved   byte      '0x7A'                                              ]
            [simple     byte      application                                         ]
            [reserved   byte      '0x00'                                              ]
        ]
        ['0x73' StatusRequestLevel
            [reserved   byte      '0x73'                                              ]
            [reserved   byte      '0x07'                                              ]
            [simple     byte      application                                         ]
            [simple     byte      startingGroupAddressLabel                           ]
            [validation           'startingGroupAddressLabel == 0x00
                                || startingGroupAddressLabel == 0x20
                                || startingGroupAddressLabel == 0x40
                                || startingGroupAddressLabel == 0x60
                                || startingGroupAddressLabel == 0x80
                                || startingGroupAddressLabel == 0xA0
                                || startingGroupAddressLabel == 0xC0
                                || startingGroupAddressLabel == 0xE0'  "invalid label"]
        ]
    ]
]

[type SALData
    [simple  SALCommandTypeContainer commandTypeContainer                                   ]
    [virtual SALCommandType          commandType          'commandTypeContainer.commandType']
    [typeSwitch commandType
        ['OFF'            SALDataOff
            [simple byte group                                                              ]
        ]
        ['ON'             SALDataOn
            [simple byte group                                                              ]
        ]
        ['RAMP_TO_LEVEL'  SALDataRampToLevel
            [simple byte group                                                              ]
            [simple byte level                                                              ]
        ]
        ['TERMINATE_RAMP' SALDataTerminateRamp
            [simple byte group                                                              ]
        ]
    ]
    // TODO: According to spec this could be recursive
    //[optional SALData salData 'what decides if this is present?']
]

[enum uint 8 SALCommandTypeContainer(SALCommandType commandType)
    ['0x01' SALCommandOff                       ['OFF'           ]]
    ['0x79' SALCommandOn                        ['ON'            ]]
    ['0x02' SALCommandRampToLevel_Instantaneous ['RAMP_TO_LEVEL' ]]
    ['0x0A' SALCommandRampToLevel_4Second       ['RAMP_TO_LEVEL' ]]
    ['0x12' SALCommandRampToLevel_8Second       ['RAMP_TO_LEVEL' ]]
    ['0x1A' SALCommandRampToLevel_12Second      ['RAMP_TO_LEVEL' ]]
    ['0x22' SALCommandRampToLevel_20Second      ['RAMP_TO_LEVEL' ]]
    ['0x2A' SALCommandRampToLevel_30Second      ['RAMP_TO_LEVEL' ]]
    ['0x32' SALCommandRampToLevel_40Second      ['RAMP_TO_LEVEL' ]]
    ['0x3A' SALCommandRampToLevel_60Second      ['RAMP_TO_LEVEL' ]]
    ['0x42' SALCommandRampToLevel_90Second      ['RAMP_TO_LEVEL' ]]
    ['0x4A' SALCommandRampToLevel_120Second     ['RAMP_TO_LEVEL' ]]
    ['0x52' SALCommandRampToLevel_180Second     ['RAMP_TO_LEVEL' ]]
    ['0x5A' SALCommandRampToLevel_300Second     ['RAMP_TO_LEVEL' ]]
    ['0x62' SALCommandRampToLevel_420Second     ['RAMP_TO_LEVEL' ]]
    ['0x6A' SALCommandRampToLevel_600Second     ['RAMP_TO_LEVEL' ]]
    ['0x72' SALCommandRampToLevel_900Second     ['RAMP_TO_LEVEL' ]]
    ['0x7A' SALCommandRampToLevel_1020Second    ['RAMP_TO_LEVEL' ]]
    ['0x09' SALCommandTerminateRamp             ['TERMINATE_RAMP']]
]

[enum uint 4 SALCommandType
    ['0x00' OFF           ]
    ['0x01' ON            ]
    ['0x02' RAMP_TO_LEVEL ]
    ['0x03' TERMINATE_RAMP]
]

[type CommandHeader
    [simple byte value]
]

[type Reply(uint 16 messageLength)
    [peek    byte peekedByte                                              ]
    [virtual bit  isAlpha '(peekedByte >= 0x67) && (peekedByte <= 0x7A)'  ]
    [typeSwitch peekedByte, isAlpha
        [*, 'true' ConfirmationReply
            [simple Confirmation isA]
        ]
        ['0x2B' PowerUpReply // is a +
            [simple PowerUp isA]
        ]
        ['0x3D' ParameterChangeReply // is a =
            [simple ParameterChange isA]
        ]
        ['0x21' ServerErrorReply // is a !
            [const  byte    errorMarker     0x21        ]
        ]
        ['0x0' MonitoredSALReply
            [simple MonitoredSAL isA]
        ]
        ['0x0' StandardFormatStatusReplyReply
            [simple StandardFormatStatusReply reply]
        ]
        ['0x0' ExtendedFormatStatusReplyReply
            [simple ExtendedFormatStatusReply reply]
        ]
        [* CALReplyReply
            [virtual uint 16 payloadLength 'messageLength-2'                        ] // We substract the termination \r\n
            [manual   CALReply
                              calReply
                                    'STATIC_CALL("readCALReply", readBuffer, payloadLength)'
                                    'STATIC_CALL("writeCALReply", writeBuffer, calReply)'
                                    '_value.lengthInBytes*2'                                     ]
        ]
    ]
    [simple   ResponseTermination termination                       ]
]

[type CALReply
    [peek    byte     calType                                                                    ]
    [typeSwitch calType
        ['0x86' CALReplyLong
            [reserved   byte                   '0x86'                                            ]
            [peek       uint 24                terminatingByte                                   ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual    bit                    isUnitAddress   '(terminatingByte & 0xff) == 0x00']
            [optional   UnitAddress            unitAddress     'isUnitAddress'                   ]
            [optional   BridgeAddress          bridgeAddress   '!isUnitAddress'                  ]
            [simple     SerialInterfaceAddress serialInterfaceAddress                            ]
            [optional   byte                   reservedByte    'isUnitAddress'                   ]
            [validation                        'isUnitAddress && reservedByte == 0x00 || !isUnitAddress' "wrong reservedByte"]
            [optional   ReplyNetwork           replyNetwork   '!isUnitAddress'                   ]
        ]
        [       CALReplyShort
        ]
    ]
    [simple   CALData   calData                                                                  ]
    //[checksum byte crc   '0x00'                                                                ] // TODO: Fix this
]

[type BridgeCount
    [simple uint 8 count]
]

[type NetworkNumber
    [simple uint 8 number]
]

[type MonitoredSAL
    [peek    byte     salType             ]
    [typeSwitch salType
        ['0x05' MonitoredSALLongFormSmartMode
            [reserved byte '0x05']
            [peek    uint 24     terminatingByte                        ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual bit isUnitAddress '(terminatingByte & 0xff) == 0x00' ]
            [optional   UnitAddress
                         unitAddress     'isUnitAddress'                ]
            [optional   BridgeAddress
                         bridgeAddress   '!isUnitAddress'               ]
            [simple     SerialInterfaceAddress
                         serialInterfaceAddress                         ]
            [optional   byte    reservedByte    'isUnitAddress'         ]
            [validation 'isUnitAddress && reservedByte == 0x00 || !isUnitAddress' "invalid unit address"]
            [optional   ReplyNetwork     replyNetwork       '!isUnitAddress'        ]
        ]
        [    MonitoredSALShortFormBasicMode
            [peek    byte                  counts                                  ]
            [optional BridgeCount          bridgeCount     'counts != 0x00'        ]
            [optional NetworkNumber        networkNumber   'counts != 0x00'        ]
            [optional byte                 noCounts        'counts == 0x00'        ] // TODO: add validation that this is 0x00 when no bridge and network number are set
            [simple ApplicationIdContainer application                             ]
        ]
    ]
    [optional SALData salData                                               ]
    //[checksum byte crc   '0x00'                                                                ] // TODO: Fix this
]

[type Confirmation
    [simple Alpha alpha]
    [discriminator   byte confirmationType]
    [typeSwitch confirmationType
        ['0x2E'    ConfirmationSuccessful              ] // "."
        ['0x23'    NotTransmittedToManyReTransmissions ] // "#"
        ['0x24'    NotTransmittedCorruption            ] // "$"
        ['0x25'    NotTransmittedSyncLoss              ] // "%"
        ['0x27'    NotTransmittedTooLong               ] // "'"
        [*         *Unknown
        ]
    ]
]

[type PowerUp
    [const    byte        powerUpIndicator       0x2B                  ] // "+"
    // TODO: do we really need a static helper to peek for terminated?=
    //[array    uint 8        garbage   terminated  '0x0D'                 ] // read all following +
    [simple   RequestTermination  reqTermination                       ] // TODO: maybe should be externalized
]

[type ParameterChange
    [const    byte        specialChar1      0x3D                    ] // "="
    [const    byte        specialChar2      0x3D                    ] // "="
]

[type ReplyNetwork
     [simple RouteType     routeType                                                    ]
     [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
     [simple UnitAddress   unitAddress                                                  ]
]

[type Checksum
    [simple byte crc]
]

[type StandardFormatStatusReply
    [simple     StatusHeader
                        statusHeader                                ]
    [simple     ApplicationIdContainer
                        application                                 ]
    [simple     uint 8  blockStart                                  ]
    [array      StatusByte
                        statusBytes
                        count
                        'statusHeader.numberOfCharacterPairs - 2'   ]
    [simple     Checksum
                        crc                                         ]
]

[type StatusHeader
    [reserved   uint 2                 '0x3'                        ]
    [simple     uint 6  numberOfCharacterPairs                      ]
]

[type ExtendedFormatStatusReply
    [simple     ExtendedStatusHeader
                        statusHeader                                ]
    [simple     StatusCoding
                        coding                                      ]
    [simple     ApplicationIdContainer
                        application                                 ]
    [simple     uint 8  blockStart                                  ]
    [array      StatusByte
                        statusBytes
                        count
                        'statusHeader.numberOfCharacterPairs - 3'   ]
    [simple     Checksum
                        crc                                         ]
]

[type ExtendedStatusHeader
    [reserved   uint 3                 '0x7'                        ]
    [simple     uint 5  numberOfCharacterPairs                      ]
]

[type StatusByte
    [simple GAVState    gav3                                        ]
    [simple GAVState    gav2                                        ]
    [simple GAVState    gav1                                        ]
    [simple GAVState    gav0                                        ]
]

[enum uint 2 GAVState
    ['0' DOES_NOT_EXIST                                             ]
    ['1' ON                                                         ]
    ['2' OFF                                                        ]
    ['3' ERROR                                                      ]
]

[enum byte StatusCoding
    ['0x00' BINARY_BY_THIS_SERIAL_INTERFACE     ]
    ['0x40' BINARY_BY_ELSEWHERE                 ]
    ['0x07' LEVEL_BY_THIS_SERIAL_INTERFACE      ]
    ['0x47' LEVEL_BY_ELSEWHERE                  ]
]

[type NetworkProtocolControlInformation
    [reserved   uint 2  '0x0'           ]
    [simple     uint 3  stackCounter    ]
    [simple     uint 3  stackDepth      ]
]

[type RequestTermination
    [const      byte    cr  0x0D                                    ] // 0xD == "<cr>"
]

[type ResponseTermination
    [const      byte    cr  0x0D                                    ] // 0xD == "<cr>"
    [const      byte    lf  0x0A                                    ] // 0xA == "<lf>"
]