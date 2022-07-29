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

[type RequestContext
    // Useful for response parsing: Set this to true if you send a CAL before. This will change the way the response will be parsed
    [simple   bit       sendCalCommandBefore        ]
    // Useful for response parsing: Set this to true if you send a SAL status request request level before. This will change the way the response will be parsed
    [simple   bit       sendStatusRequestLevelBefore  ]
    // Useful for response parsing: Set this to true if you send a identify request before. This will change the way the response will be parsed
    [simple   bit       sendIdentifyRequestBefore   ]
]

[type CBusOptions
    // Defines that SAL messages can occur at any time
    [simple bit connect]
    // Disable echo of characters. When used with connect SAL have a long option. Select long from of most CAL replies
    [simple bit smart  ]
    // only works with smart. Select long form of CAL messages
    [simple bit idmon  ]
    // useful with smart. Select long form, extended format for all monitored and initiated status requests
    [simple bit exstat ]
    // monitors all traffic for status requests. Status requests will be returned as CAL. Replies are modified by exstat. Usually used in conjunction with connect.
    [simple bit monitor]
    // Same as connect. In addition it will return remote network SAL
    [simple bit monall ]
    // Serial interface will emit a power up notification
    [simple bit pun    ]
    // causes parameter change notifications to be emitted.
    [simple bit pcn    ]
    // enabled the checksum checks
    [simple bit srchk ]
]

[type CBusMessage(bit isResponse, RequestContext requestContext, CBusOptions cBusOptions)
    [validation 'requestContext != null' "requestContext required"  ]
    [validation 'cBusOptions != null'    "cBusOptions required"     ]
    [typeSwitch isResponse
       ['false' *ToServer
            [simple   Request('cBusOptions')         request        ]
       ]
       ['true' *ToClient
            [simple   ReplyOrConfirmation('cBusOptions', 'requestContext')  reply           ]
       ]
    ]
]

[type Request(CBusOptions cBusOptions)
    [peek     RequestType peekedByte                                        ]
    [optional RequestType startingCR       'peekedByte == RequestType.EMPTY']
    [optional RequestType resetMode        'peekedByte == RequestType.RESET']
    [peek     RequestType secondPeek                                        ]
    [virtual  RequestType actualPeek       '(startingCR==null&&resetMode==null)||(startingCR==null&&resetMode!=null&&secondPeek==RequestType.EMPTY)?peekedByte:secondPeek'  ]
    [typeSwitch actualPeek
        ['SMART_CONNECT_SHORTCUT' *SmartConnectShortcut
            [const    byte        pipe      0x7C                            ]
            [peek     RequestType pipePeek                                  ]
            [optional byte        secondPipe 'pipePeek == RequestType.SMART_CONNECT_SHORTCUT']
        ]
        ['RESET' *Reset
            [peek     RequestType tildePeek                                     ]
            [optional byte        secondTilde 'tildePeek == RequestType.RESET'  ]
            [peek     RequestType tildePeek2                                    ]
            [optional byte        thirdTilde 'tildePeek2 == RequestType.RESET'  ]
        ]
        ['DIRECT_COMMAND' *DirectCommandAccess
            [const    byte    at        0x40                                ]
            [manual   CALData
                              calData
                        'STATIC_CALL("readCALData", readBuffer)'
                        'STATIC_CALL("writeCALData", writeBuffer, calData)'
                        '(calData.lengthInBytes*2)*8'                       ]
        ]
        ['REQUEST_COMMAND' *Command
            [const    byte  initiator 0x5C                                  ] // 0x5C == "/"
            [manual   CBusCommand
                              cbusCommand
                        'STATIC_CALL("readCBusCommand", readBuffer, cBusOptions, cBusOptions.srchk)'
                        'STATIC_CALL("writeCBusCommand", writeBuffer, cbusCommand)'
                        '(cbusCommand.lengthInBytes*2)*8'                   ]
            [manual   Checksum
                              chksum
                        'STATIC_CALL("readAndValidateChecksum", readBuffer, cbusCommand, cBusOptions.srchk)'
                        'STATIC_CALL("calculateChecksum", writeBuffer, cbusCommand, cBusOptions.srchk)'
                        '(cBusOptions.srchk)?(16):(0)'                      ]
            [optional Alpha         alpha                                   ]
        ]
        ['NULL' *Null
            [const    uint 32             nullIndicator        0x6E756C6C   ] // "null"
        ]
        ['EMPTY' *Empty
        ]
        // TODO: we should check if we are in basic mode
        [* *Obsolete
            [manual   CALData
                              calData
                        'STATIC_CALL("readCALData", readBuffer)'
                        'STATIC_CALL("writeCALData", writeBuffer, calData)'
                        '(calData.lengthInBytes*2)*8'                       ]
            [optional Alpha   alpha                                         ]
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

[discriminatedType CBusCommand(CBusOptions cBusOptions)
    [simple  CBusHeader header           ]
    [virtual bit        isDeviceManagement 'header.dp']
    // TODO: header.destinationAddressType could be used directly but for this we need source type resolving to work (WIP)
    [virtual DestinationAddressType destinationAddressType 'header.destinationAddressType']
    [typeSwitch destinationAddressType, isDeviceManagement
        [*, 'true' *DeviceManagement
            [simple     Parameter paramNo                                 ]
            [const      byte      delimiter       0x0                     ]
            [simple     byte      parameterValue                          ]
        ]
        ['PointToPointToMultiPoint' *PointToPointToMultiPoint
            [simple CBusPointToPointToMultiPointCommand('cBusOptions') command]
        ]
        ['PointToMultiPoint'        *PointToMultiPoint
            [simple CBusPointToMultiPointCommand('cBusOptions')        command]
        ]
        ['PointToPoint'             *PointToPoint
            [simple CBusPointToPointCommand('cBusOptions')             command]
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
    [simple NetworkProtocolControlInformation networkPCI                                ]
    [array  BridgeAddress additionalBridgeAddresses count 'networkPCI.stackDepth-1'     ] // We substract 1 as when a route is used we always have one prefixed
]

[discriminatedType CBusPointToPointCommand(CBusOptions cBusOptions)
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
    [simple   CALData('null') calData                                                           ]
]

[discriminatedType CBusPointToMultiPointCommand(CBusOptions cBusOptions)
    [peek    byte     peekedApplication                                                                ]
    [typeSwitch peekedApplication
        ['0xFF'   *Status
            [reserved byte          '0xFF'                                                             ]
            [reserved byte          '0x00'                                                             ]
            [simple   StatusRequest statusRequest                                                      ]
        ]
        [*        *Normal
            [simple   ApplicationIdContainer                application                                ]
            [reserved byte                                  '0x00'                                     ]
            [simple   SALData('application.applicationId')  salData                                    ]
        ]
    ]
]

[discriminatedType CBusPointToPointToMultiPointCommand(CBusOptions cBusOptions)
    [simple BridgeAddress bridgeAddress                                                              ]
    [simple NetworkRoute  networkRoute                                                               ]
    [peek    byte       peekedApplication                                                            ]
    [typeSwitch peekedApplication
        ['0xFF'   *Status
            [reserved byte        '0xFF'                                                             ]
            [simple StatusRequest statusRequest                                                      ]
        ]
        [*        *Normal
            [simple   ApplicationIdContainer                application                              ]
            [simple   SALData('application.applicationId')  salData                                  ]
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
    ['0x14' MEDIA_TRANSPORT_CONTROL           ]
    ['0x15' ERROR_REPORTING                   ]
    ['0x16' HVAC_ACTUATOR                     ]
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
    ['0x73' RESERVED_73                           ['HVAC_ACTUATOR'                     , 'NA'                  ]] // HVAC_ACTUATOR
    ['0x74' RESERVED_74                           ['HVAC_ACTUATOR'                     , 'NA'                  ]] // HVAC_ACTUATOR
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
    ['0xC0' MEDIA_TRANSPORT_CONTROL_C0            ['MEDIA_TRANSPORT_CONTROL'           , 'NA'                  ]] // MEDIA_TRANSPORT_CONTROL
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
    ['0xCE' ERROR_REPORTING_CE                    ['ERROR_REPORTING'                   , 'NA'                  ]] // ERROR_REPORTING
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

[type CALData(RequestContext requestContext)
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsCALCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  CALCommandTypeContainer commandTypeContainer                                   ]
    [virtual CALCommandType          commandType          'commandTypeContainer.commandType']
    [virtual bit  sendIdentifyRequestBefore       'requestContext!=null?requestContext.sendIdentifyRequestBefore:false']
    [typeSwitch commandType, sendIdentifyRequestBefore
        ['RESET'            *Reset                                                              // Request
        ]
        ['RECALL'           *Recall                                                             // Request
            [simple Parameter paramNo                                                       ]
            [simple uint 8    count                                                         ]
        ]
        ['IDENTIFY'         *Identify                                                           // Request
            [simple Attribute attribute                                                     ]
        ]
        ['GET_STATUS'       *GetStatus // Request
            [simple Parameter paramNo                                                       ]
            [simple uint 8    count                                                         ]
        ]
        ['WRITE'            *Write(CALCommandTypeContainer commandTypeContainer)                // Request
            [simple Parameter paramNo                                                       ]
            [simple byte      code                                                          ]
            // TODO: we can decode this with the parametert above... e.g. INTERFACE_OPTIONS_1 is defined below
            [array  byte      data        count 'commandTypeContainer.numBytes - 2'         ]
        ]
        ['REPLY', 'true'    *IdentifyReply(CALCommandTypeContainer commandTypeContainer)        // Reply
            [simple Attribute   attribute                                                   ]
            [simple IdentifyReplyCommand('attribute', 'commandTypeContainer.numBytes - 1')
                                identifyReplyCommand                                        ]
        ]
        ['REPLY'            *Reply(CALCommandTypeContainer commandTypeContainer)                // Reply
            [simple Parameter paramNo                                                       ]
            [array  byte      data        count 'commandTypeContainer.numBytes-1'           ]
        ]
        ['ACKNOWLEDGE'      *Acknowledge // Reply
            [simple Parameter paramNo                                                       ]
            [simple uint 8    code                                                          ]
        ]
        ['STATUS'           *Status(CALCommandTypeContainer commandTypeContainer)               // Reply
            [simple ApplicationIdContainer application                                                 ]
            [simple uint 8                 blockStart                                                  ]
            [array  byte                   data        count 'commandTypeContainer.numBytes - 2'       ]
        ]
        ['STATUS_EXTENDED'  *StatusExtended(CALCommandTypeContainer commandTypeContainer)       // Reply
            [simple uint 8                 coding                                                      ]
            [virtual bit                   isBinaryBySerialInterface 'coding == 0x00'                  ]
            [virtual bit                   isBinaryByElsewhere       'coding == 0x40'                  ]
            [virtual bit                   isLevelBySerialInterface  'coding == 0x07'                  ]
            [virtual bit                   isLevelByElsewhere        'coding == 0x47'                  ]
            [virtual bit                   isReserved                '!isBinaryBySerialInterface && !isBinaryByElsewhere && !isLevelBySerialInterface && !isLevelByElsewhere']
            [simple ApplicationIdContainer application                                                 ]
            [simple uint 8                 blockStart                                                  ]
            [array  byte                   data        count 'commandTypeContainer.numBytes - 2'       ] // TODO: this should be -3 but somehow it is -2 with the examples
        ]
    ]
    // Note: we omit the request context as it is only useful for the first element
    [optional CALData('null') additionalData]
]

[enum uint 8 Parameter(vstring group, vstring parameterDescription, vstring form, bit isVolatile, ProtectionLevel protectionLevel)
    ['0x00' UNKNOWN_01                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x01' UNKNOWN_02                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x02' UNKNOWN_03                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x03' UNKNOWN_04                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x04' UNKNOWN_05                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x05' UNKNOWN_06                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x06' UNKNOWN_07                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x07' UNKNOWN_08                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x08' UNKNOWN_09                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x09' UNKNOWN_10                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0A' UNKNOWN_11                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0B' UNKNOWN_12                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0C' UNKNOWN_13                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0D' UNKNOWN_14                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0E' UNKNOWN_15                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x0F' UNKNOWN_16                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x10' UNKNOWN_17                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x11' UNKNOWN_18                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x12' UNKNOWN_19                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x13' UNKNOWN_20                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x14' UNKNOWN_21                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x15' UNKNOWN_22                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x16' UNKNOWN_23                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x17' UNKNOWN_24                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x18' UNKNOWN_25                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x19' UNKNOWN_26                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1A' UNKNOWN_27                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1B' UNKNOWN_28                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1C' UNKNOWN_29                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1D' UNKNOWN_30                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1E' UNKNOWN_31                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x1F' UNKNOWN_32                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x20' UNKNOWN_33                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x21' APPLICATION_ADDRESS_1                   ['"Mgmt"', '"Application Address 1"',                   '"Byte (Note 1)"',          'false', 'UNLOCK_REQUIRED']]
    ['0x22' APPLICATION_ADDRESS_2                   ['"Mgmt"', '"Application Address 2"',                   '"Byte (Note 1)"',          'false', 'UNLOCK_REQUIRED']]
    ['0x23' UNKOWN_35                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x24' UNKOWN_36                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x25' UNKOWN_37                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x26' UNKOWN_38                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x27' UNKOWN_39                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x28' UNKOWN_40                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x29' UNKOWN_41                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2A' UNKOWN_42                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2B' UNKOWN_43                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2C' UNKOWN_44                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2D' UNKOWN_45                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2E' UNKOWN_46                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x2F' UNKOWN_47                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x30' INTERFACE_OPTIONS_1                     ['"Unit"', '"Interface options 1"',                     '"8 Bits (Note 2)"',        'true',  'NO_WRITE_ACCESS']]
    ['0x31' UNKOWN_49                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x32' UNKOWN_50                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x33' UNKOWN_51                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x34' UNKOWN_52                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x35' UNKOWN_53                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x36' UNKOWN_54                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x37' UNKOWN_55                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x38' UNKOWN_56                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x39' UNKOWN_57                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x3A' UNKOWN_58                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x3B' UNKOWN_59                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x3C' UNKOWN_60                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x3D' BAUD_RATE_SELECTOR                      ['"Unit"', '"Baud rate selector"',                      '"Byte (Note 3)"',          'false', 'NO_WRITE_ACCESS']]
    ['0x3E' INTERFACE_OPTIONS_2                     ['"Unit"', '"Interface options 2"',                     '"Byte (Note 4)"',          'false', 'NONE'           ]]
    ['0x3F' UNKOWN_63                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x40' UNKOWN_64                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x41' INTERFACE_OPTIONS_1_POWER_UP_SETTINGS   ['"Unit"', '"Interface options 2 power up settings"',   '"8 Bits (Note 5)"',        'false', 'UNLOCK_REQUIRED']]
    ['0x42' INTERFACE_OPTIONS_3                     ['"Unit"', '"Interface options 3"',                     '"Byte (Note 6)"',          'false', 'UNLOCK_REQUIRED']]
    ['0x43' UNKOWN_67                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x44' UNKOWN_68                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x45' UNKOWN_69                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x46' UNKOWN_70                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x47' UNKOWN_71                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x48' UNKOWN_72                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x49' UNKOWN_73                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4A' UNKOWN_74                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4B' UNKOWN_75                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4C' UNKOWN_76                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4D' UNKOWN_77                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4E' UNKOWN_78                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x4F' UNKOWN_79                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x50' UNKOWN_80                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x51' UNKOWN_81                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x52' UNKOWN_82                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x53' UNKOWN_83                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x54' UNKOWN_84                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x55' UNKOWN_85                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x56' UNKOWN_86                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x57' UNKOWN_87                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x58' UNKOWN_88                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x59' UNKOWN_89                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5A' UNKOWN_90                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5B' UNKOWN_91                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5C' UNKOWN_92                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5D' UNKOWN_93                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5E' UNKOWN_94                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x5F' UNKOWN_95                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x60' UNKOWN_96                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x61' UNKOWN_97                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x62' UNKOWN_98                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x63' UNKOWN_99                               ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x64' UNKOWN_100                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x65' UNKOWN_101                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x66' UNKOWN_102                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x67' UNKOWN_103                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x68' UNKOWN_104                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x69' UNKOWN_105                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6A' UNKOWN_106                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6B' UNKOWN_107                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6C' UNKOWN_108                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6D' UNKOWN_109                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6E' UNKOWN_110                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x6F' UNKOWN_111                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x70' UNKOWN_112                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x71' UNKOWN_113                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x72' UNKOWN_114                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x73' UNKOWN_115                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x74' UNKOWN_116                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x75' UNKOWN_117                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x76' UNKOWN_118                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x77' UNKOWN_119                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x78' UNKOWN_120                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x79' UNKOWN_121                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7A' UNKOWN_122                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7B' UNKOWN_123                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7C' UNKOWN_124                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7D' UNKOWN_125                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7E' UNKOWN_126                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x7F' UNKOWN_127                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x80' UNKOWN_128                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x81' UNKOWN_129                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x82' UNKOWN_130                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x83' UNKOWN_131                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x84' UNKOWN_132                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x85' UNKOWN_133                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x86' UNKOWN_134                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x87' UNKOWN_135                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x88' UNKOWN_136                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x89' UNKOWN_137                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8A' UNKOWN_138                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8B' UNKOWN_139                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8C' UNKOWN_140                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8D' UNKOWN_141                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8E' UNKOWN_142                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x8F' UNKOWN_143                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x90' UNKOWN_144                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x91' UNKOWN_145                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x92' UNKOWN_146                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x93' UNKOWN_147                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x94' UNKOWN_148                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x95' UNKOWN_149                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x96' UNKOWN_150                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x97' UNKOWN_151                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x98' UNKOWN_152                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x99' UNKOWN_153                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9A' UNKOWN_154                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9B' UNKOWN_155                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9C' UNKOWN_156                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9D' UNKOWN_157                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9E' UNKOWN_158                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0x9F' UNKOWN_159                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA0' UNKOWN_160                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA1' UNKOWN_161                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA2' UNKOWN_162                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA3' UNKOWN_163                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA4' UNKOWN_164                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA5' UNKOWN_165                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA6' UNKOWN_166                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA7' UNKOWN_167                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA8' UNKOWN_168                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xA9' UNKOWN_169                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAA' UNKOWN_170                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAB' UNKOWN_171                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAC' UNKOWN_172                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAD' UNKOWN_173                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAE' UNKOWN_174                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xAF' UNKOWN_175                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB0' UNKOWN_176                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB1' UNKOWN_177                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB2' UNKOWN_178                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB3' UNKOWN_179                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB4' UNKOWN_180                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB5' UNKOWN_181                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB6' UNKOWN_182                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB7' UNKOWN_183                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB8' UNKOWN_184                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xB9' UNKOWN_185                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBA' UNKOWN_186                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBB' UNKOWN_187                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBC' UNKOWN_188                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBD' UNKOWN_189                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBE' UNKOWN_190                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xBF' UNKOWN_191                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC0' UNKOWN_192                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC1' UNKOWN_193                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC2' UNKOWN_194                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC3' UNKOWN_195                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC4' UNKOWN_196                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC5' UNKOWN_197                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC6' UNKOWN_198                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC7' UNKOWN_199                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC8' UNKOWN_200                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xC9' UNKOWN_201                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCA' UNKOWN_202                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCB' UNKOWN_203                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCC' UNKOWN_204                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCD' UNKOWN_205                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCE' UNKOWN_206                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xCF' UNKOWN_207                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD0' UNKOWN_208                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD1' UNKOWN_209                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD2' UNKOWN_210                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD3' UNKOWN_211                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD4' UNKOWN_212                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD5' UNKOWN_213                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD6' UNKOWN_214                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD7' UNKOWN_215                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD8' UNKOWN_216                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xD9' UNKOWN_217                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDA' UNKOWN_218                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDB' UNKOWN_219                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDC' UNKOWN_220                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDD' UNKOWN_221                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDE' UNKOWN_222                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xDF' UNKOWN_223                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE0' UNKOWN_224                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE1' UNKOWN_225                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE2' UNKOWN_226                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE3' UNKOWN_227                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE4' UNKOWN_228                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE5' UNKOWN_229                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE6' UNKOWN_230                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE7' UNKOWN_231                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE8' UNKOWN_232                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xE9' UNKOWN_233                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xEA' UNKOWN_234                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
    ['0xEB' CUSTOM_MANUFACTURER_1                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xEC' CUSTOM_MANUFACTURER_2                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xED' CUSTOM_MANUFACTURER_3                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xEE' CUSTOM_MANUFACTURER_4                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xEF' CUSTOM_MANUFACTURER_5                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xF0' CUSTOM_MANUFACTURER_6                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xF1' CUSTOM_MANUFACTURER_7                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'UNLOCK_REQUIRED']]
    ['0xF2' CUSTOM_MANUFACTURER_8                   ['"Mgmt"', '"Custom Manufacturer (8 bytes)"',           '"ASCII Chars (Note 7)"',   'false', 'READ_ONLY'      ]]
    ['0xF3' SERIAL_NUMBER_1                         ['"Mgmt"', '"Serial Number"',                           '"Bytes (Note 8)"',         'false', 'READ_ONLY'      ]]
    ['0xF4' SERIAL_NUMBER_2                         ['"Mgmt"', '"Serial Number"',                           '"Bytes (Note 8)"',         'false', 'READ_ONLY'      ]]
    ['0xF5' SERIAL_NUMBER_3                         ['"Mgmt"', '"Serial Number"',                           '"Bytes (Note 8)"',         'false', 'READ_ONLY'      ]]
    ['0xF6' SERIAL_NUMBER_4                         ['"Mgmt"', '"Serial Number"',                           '"Bytes (Note 8)"',         'false', 'READ_ONLY'      ]]
    ['0xF7' CUSTOM_TYPE_1                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xF8' CUSTOM_TYPE_2                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xF9' CUSTOM_TYPE_3                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFA' CUSTOM_TYPE_4                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFB' CUSTOM_TYPE_5                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFC' CUSTOM_TYPE_6                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFD' CUSTOM_TYPE_7                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFE' CUSTOM_TYPE_8                           ['"Mgmt"', '"Custom Type (8 bytes)"',                   '"ASCII Chars (Note 9)"',   'false', 'READ_ONLY'      ]]
    ['0xFF' UNKOWN_255                              ['""',     '""',                                        '""',                       'false', 'NONE'           ]]
]

[enum uint 4 ProtectionLevel(vstring description)
    ['0' UNLOCK_REQUIRED    ['"Unlock required from C-BUS port"']]
    ['1' NO_WRITE_ACCESS    ['"No write access via C-BUS port"' ]]
    ['2' NONE               ['"None"'                           ]]
    ['3' READ_ONLY          ['"Read only"'                      ]]
]

[type ApplicationAddress1 // Note 1
    [simple  byte address                       ]
    // if wildcard is set address 2 should set to wildcard as well
    [virtual bit  isWildcard 'address == 0xFF'  ]
]

[type ApplicationAddress2 // Note 1
    [simple  byte address                       ]
    [virtual bit  isWildcard 'address == 0xFF'  ]
]

[type InterfaceOptions1 // Note 2
    [reserved bit  'false'                       ]
    [simple   bit  idmon                         ]
    [simple   bit  monitor                       ]
    [simple   bit  smart                         ]
    [simple   bit  srchk                         ]
    [simple   bit  xonXoff                       ]
    [reserved bit  'false'                       ]
    [simple   bit  connect                       ]
]

// Undefined values default to 0xFF
[enum uint 8 BaudRateSelector
    ['0x01' SELECTED_4800_BAUD]
    ['0x02' SELECTED_2400_BAUD]
    ['0x03' SELECTED_1200_BAUD]
    ['0x04' SELECTED_600_BAUD ]
    ['0x05' SELECTED_300_BAUD ]
    ['0xFF' SELECTED_9600_BAUD]
]

[type InterfaceOptions2 // Note 4
    [reserved bit  'false'                       ]
    [simple   bit  burden                        ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [simple   bit  clockGen                      ]
]

[type InterfaceOptions1PowerUpSettings // Note 5
    [simple InterfaceOptions1 interfaceOptions1  ]
]

[type InterfaceOptions3 // Note 6
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [reserved bit  'false'                       ]
    [simple   bit  exstat                        ]
    [simple   bit  pun                           ]
    [simple   bit  localSal                      ]
    [simple   bit  pcn                           ]
]

[type CustomManufacturer // Note 7
    // TODO: 8 is a placeholder at the moment
    [simple vstring '8' customString                 ]
]

[type SerialNumber // Note 8
    [simple byte octet1]
    [simple byte octet2]
    [simple byte octet3]
    [simple byte octet4]
]

[type CustomTypes // Note 9
    // TODO: 8 is a placeholder at the moment
    [simple vstring '8' customString                 ]
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

[type IdentifyReplyCommand(Attribute attribute, uint 5 numBytes)
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
            [simple   ApplicationIdContainer  lowApplication         ]
            [simple   ApplicationIdContainer  highApplication        ]
            [simple   byte                    area                   ]
            [simple   uint 16                 crc                    ]
            [simple   uint 32                 serialNumber           ]
            [simple   byte                    networkVoltage         ]
            [virtual  float 32                networkVoltageInVolts 'networkVoltage/6.375']
            [simple   bit                     unitInLearnMode        ]
            [simple   bit                     networkVoltageLow      ]
            [simple   bit                     networkVoltageMarginal ]
            [reserved uint 1                  '0'                    ]
            [reserved uint 1                  '0'                    ]
            [reserved uint 1                  '0'                    ]
            [simple   bit                     enableChecksumAlarm    ]
            [simple   bit                     outputUnit             ]
            [simple   bit                     installationMMIError   ]
            [simple   bit                     EEWriteError           ]
            [simple   bit                     EEChecksumError        ]
            [simple   bit                     EEDataError            ]
            [simple   bit                     microReset             ]
            [simple   bit                     commsTxError           ]
            [simple   bit                     internalStackOverflow  ]
            [simple   bit                     microPowerReset        ]
        ]
        ['NetworkTerminalLevels'        IdentifyReplyCommandNetworkTerminalLevels
            [array  byte        minimumLevels        count 'numBytes'       ] // TODO: check datatype
        ]
        ['TerminalLevel'                IdentifyReplyCommandTerminalLevels
            [array  byte        terminalLevels        count 'numBytes'       ] // TODO: check datatype
        ]
        ['NetworkVoltage'               IdentifyReplyCommandNetworkVoltage
           [simple string 2     volts                   ]
           [const  byte         dot     0x2C            ]
           [simple string 2     voltsDecimalPlace       ]
           [const  byte         v       0x56            ]
        ]
        ['GAVValuesCurrent'             IdentifyReplyCommandGAVValuesCurrent
            [array  byte        values  count   'numBytes'    ] // TODO: check datatype
        ]
        ['GAVValuesStored'              IdentifyReplyCommandGAVValuesStored
            [array  byte        values  count   'numBytes'    ] // TODO: check datatype
        ]
        ['GAVPhysicalAddresses'         IdentifyReplyCommandGAVPhysicalAddresses
            [array  byte        values  count   'numBytes'    ] // TODO: check datatype
        ]
        ['LogicalAssignment'            IdentifyReplyCommandLogicalAssignment
            [array  LogicAssignment   logicAssigment        count 'numBytes'       ]
        ]
        ['Delays'                       IdentifyReplyCommandDelays
            [array  byte        terminalLevels        count 'numBytes-1'       ]
            [simple byte        reStrikeDelay                   ]
        ]
        ['MinimumLevels'                IdentifyReplyCommandMinimumLevels
            [array  byte        minimumLevels       count 'numBytes'       ]
        ]
        ['MaximumLevels'                IdentifyReplyCommandMaximumLevels
            [array  byte        maximumLevels       count 'numBytes'       ]
        ]
        ['CurrentSenseLevels'           IdentifyReplyCommandCurrentSenseLevels
            [array  byte        currentSenseLevels  count 'numBytes'       ]
        ]
        ['OutputUnitSummary'            IdentifyReplyCommandOutputUnitSummary
            // TODO: we can use the bytes from above, but how is that dynamic? repeat the complete block here?
            [simple   IdentifyReplyCommandUnitSummary
                             unitFlags                              ]
            [simple   byte   gavStoreEnabledByte1                   ]
            [simple   byte   gavStoreEnabledByte2                   ]
            [simple   uint 8 timeFromLastRecoverOfMainsInSeconds    ]
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

[type IdentifyReplyCommandUnitSummary
    [simple bit assertingNetworkBurden  ]
    [simple bit restrikeTimingActive    ]
    [simple bit remoteOFFInputAsserted  ]
    [simple bit remoteONInputAsserted   ]
    [simple bit localToggleEnabled      ]
    [simple bit localToggleActiveState  ]
    [simple bit clockGenerationEnabled  ]
    [simple bit unitGeneratingClock     ]
]

[type LogicAssignment
    [simple   bit greaterOfOrLogic  ]
    [simple   bit reStrikeDelay     ]
    [reserved bit 'false'           ]
    [reserved bit 'false'           ]
    [simple   bit assignedToGav16   ]
    [simple   bit assignedToGav15   ]
    [simple   bit assignedToGav14   ]
    [simple   bit assignedToGav13   ]
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
    ['0x32' CALCommandAcknowledge            ['ACKNOWLEDGE',      '0']]
    ['0x80' CALCommandReply_0Bytes           ['REPLY',            '0']]
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
    ['0xA0' CALCommandWrite_0Bytes           ['WRITE',            '0']]
    ['0xA1' CALCommandWrite_1Bytes           ['WRITE',            '1']]
    ['0xA2' CALCommandWrite_2Bytes           ['WRITE',            '2']]
    ['0xA3' CALCommandWrite_3Bytes           ['WRITE',            '3']]
    ['0xA4' CALCommandWrite_4Bytes           ['WRITE',            '4']]
    ['0xA5' CALCommandWrite_5Bytes           ['WRITE',            '5']]
    ['0xA6' CALCommandWrite_6Bytes           ['WRITE',            '6']]
    ['0xA7' CALCommandWrite_7Bytes           ['WRITE',            '7']]
    ['0xA8' CALCommandWrite_8Bytes           ['WRITE',            '8']]
    ['0xA9' CALCommandWrite_9Bytes           ['WRITE',            '9']]
    ['0xAA' CALCommandWrite_10Bytes          ['WRITE',           '10']]
    ['0xAB' CALCommandWrite_11Bytes          ['WRITE',           '11']]
    ['0xAC' CALCommandWrite_12Bytes          ['WRITE',           '12']]
    ['0xAD' CALCommandWrite_13Bytes          ['WRITE',           '13']]
    ['0xAE' CALCommandWrite_14Bytes          ['WRITE',           '14']]
    ['0xAF' CALCommandWrite_15Bytes          ['WRITE',           '15']]
    ['0xC0' CALCommandStatus_0Bytes          ['STATUS',           '0']]
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
    ['0xE0' CALCommandStatusExtended_0Bytes  ['STATUS_EXTENDED',  '0']]
    ['0xE1' CALCommandStatusExtended_1Bytes  ['STATUS_EXTENDED',  '1']]
    ['0xE2' CALCommandStatusExtended_2Bytes  ['STATUS_EXTENDED',  '2']]
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

[enum uint 8 CALCommandType
    // Request
    ['0x00' RESET          ]
    ['0x01' RECALL         ]
    ['0x02' IDENTIFY       ]
    ['0x03' GET_STATUS     ]
    ['0x04' WRITE          ]
    // Response
    ['0x0F' REPLY          ]
    ['0x10' ACKNOWLEDGE    ]
    ['0x11' STATUS         ]
    ['0x12' STATUS_EXTENDED]
]

[type StatusRequest
    [peek    byte     statusType           ]
    [typeSwitch statusType
        ['0x7A' *BinaryState
            [reserved   byte                    '0x7A'                                              ]
            [simple     ApplicationIdContainer  application                                         ]
            [reserved   byte                    '0x00'                                              ]
        ]
        ['0xFA' *BinaryStateDeprecated
            [reserved   byte                    '0xFA'                                              ]
            [simple     ApplicationIdContainer  application                                         ]
            [reserved   byte                    '0x00'                                              ]
        ]
        ['0x73' *Level
            [reserved   byte                    '0x73'                                              ]
            [reserved   byte                    '0x07'                                              ]
            [simple     ApplicationIdContainer  application                                         ]
            [simple     byte                    startingGroupAddressLabel                           ]
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

// TODO: this is currently lightning only so we need more typeSwitched based on the applicationid
[type SALData(ApplicationId applicationId)
    [typeSwitch applicationId
        ['RESERVED'                             *Reserved
            [validation '1==2' "RESERVED Not yet implemented"] // TODO: implement me
        ]
        ['FREE_USAGE'                           *FreeUsage
            [validation '1==2' "FREE_USAGE Not yet implemented"] // TODO: implement me
        ]
        ['TEMPERATURE_BROADCAST'                *TemperatureBroadcast
            [simple TemperatureBroadcastData temperatureBroadcastData]
        ]
        ['ROOM_CONTROL_SYSTEM'                  *RoomControlSystem
            [validation '1==2' "ROOM_CONTROL_SYSTEM Not yet implemented"] // TODO: implement me
        ]
        ['LIGHTING'                             *Lighting
            [simple LightingData lightingData]
        ]
        ['VENTILATION'                          *Ventilation
            // Note: the documentation states that the data for ventilation uses LightingData
            [simple LightingData ventilationData]
        ]
        ['IRRIGATION_CONTROL'                   *IrrigationControl
             // Note: the documentation states that the data for irrigation control uses LightingData
            [simple LightingData irrigationControlData]
        ]
        ['POOLS_SPAS_PONDS_FOUNTAINS_CONTROL'   *PoolsSpasPondsFountainsControl
             // Note: the documentation states that the data for pools spas ponds fountains uses LightingData
            [simple LightingData poolsSpaPondsFountainsData]
        ]
        ['HEATING'                              *Heating
            // Note: the documentation states that the data for ventilation uses LightingData
            [simple LightingData heatingData]
        ]
        ['AIR_CONDITIONING'                     *AirConditioning
            [simple AirConditioningData airConditioningData]
        ]
        ['TRIGGER_CONTROL'                      *TriggerControl
            [simple TriggerControlData triggerControlData]
        ]
        ['ENABLE_CONTROL'                       *EnableControl
            [simple EnableControlData enableControlData]
        ]
        ['AUDIO_AND_VIDEO'                      *AudioAndVideo
             // Note: the documentation states that the data for ventilation uses LightingData
            [simple LightingData audioVideoData]
        ]
        ['SECURITY'                             *Security
            [simple SecurityData securityData]
        ]
        ['METERING'                             *Metering
            [simple MeteringData meteringData]
        ]
        ['ACCESS_CONTROL'                       *AccessControl
            [simple AccessControlData accessControlData]
        ]
        ['CLOCK_AND_TIMEKEEPING'                *ClockAndTimekeeping
            [simple ClockAndTimekeepingData clockAndTimekeepingData]
        ]
        ['TELEPHONY_STATUS_AND_CONTROL'         *TelephonyStatusAndControl
            [simple TelephonyData telephonyData]
        ]
        ['MEASUREMENT'                          *Measurement
            [simple MeasurementData measurementData]
        ]
        ['TESTING'                              *Testing
            [validation '1==2' "TESTING Not yet implemented"] // TODO: implement me
        ]
        ['MEDIA_TRANSPORT_CONTROL'              *MediaTransport
            [simple MediaTransportControlData   mediaTransportControlData]
        ]
        ['ERROR_REPORTING'                      *ErrorReporting
            [simple ErrorReportingData   errorReportingData]
        ]
        ['HVAC_ACTUATOR'                        *HvacActuator
             // Note: the documentation states that the data for hvac actuator uses LightingData
            [simple LightingData ventilationData]
        ]
    ]
    [optional SALData('applicationId') salData                                  ]
]

[type ReplyOrConfirmation(CBusOptions cBusOptions, RequestContext requestContext)
    [peek    byte peekedByte                                                ]
    [virtual bit  isAlpha '(peekedByte >= 0x67) && (peekedByte <= 0x7A)'    ]
    [typeSwitch isAlpha
        ['true' *Confirmation
            [simple   Confirmation                      confirmation        ]
            [optional ReplyOrConfirmation('cBusOptions', 'requestContext') embeddedReply]
        ]
        ['false' *Reply
            [simple   Reply('cBusOptions', 'requestContext')    reply               ]
            [simple   ResponseTermination               termination         ]
        ]
    ]
]

[type Reply(CBusOptions cBusOptions, RequestContext requestContext)
    [peek    byte peekedByte                                                                ]
    [typeSwitch peekedByte
        ['0x2B' PowerUpReply // is a +
            [simple PowerUp isA]
        ]
        ['0x3D' ParameterChangeReply // is a =
            [simple ParameterChange isA                 ]
        ]
        ['0x21' ServerErrorReply // is a !
            [const  byte    errorMarker     0x21        ]
        ]
        [*      *EncodedReply
            [manual   EncodedReply
                              encodedReply
                                    'STATIC_CALL("readEncodedReply", readBuffer, cBusOptions, requestContext, cBusOptions.srchk)'
                                    'STATIC_CALL("writeEncodedReply", writeBuffer, encodedReply)'
                                    '(encodedReply.lengthInBytes*2)*8'                                   ]
            [manual   Checksum
                              chksum
                        'STATIC_CALL("readAndValidateChecksum", readBuffer, encodedReply, cBusOptions.srchk)'
                        'STATIC_CALL("calculateChecksum", writeBuffer, encodedReply, cBusOptions.srchk)'
                        '(cBusOptions.srchk)?(16):(0)'        ]
        ]
    ]
]

[type EncodedReply(CBusOptions cBusOptions, RequestContext requestContext)
    [peek    byte peekedByte                                                        ]
    // TODO: if we reliable can detect this with the mask we don't need the request context anymore
    [virtual bit  isMonitoredSAL            '(peekedByte & 0x3F) == 0x05 || peekedByte == 0x00 || (peekedByte & 0xF8) == 0x00'] // First check if it is in long mode, second for short mode, third for bridged short mode
    [virtual bit  isCalCommand              '(peekedByte & 0x3F) == 0x06 || requestContext.sendCalCommandBefore'    ] // The 0x3F and 0x06 doesn't seem to work always
    [virtual bit  isStandardFormatStatus    '(peekedByte & 0xC0) == 0xC0 && !cBusOptions.exstat'                    ]
    [virtual bit  isExtendedFormatStatus    '(peekedByte & 0xE0) == 0xE0 && (cBusOptions.exstat || requestContext.sendStatusRequestLevelBefore)']
    [typeSwitch isMonitoredSAL, isCalCommand, isStandardFormatStatus, isExtendedFormatStatus
        ['true', 'false', 'false'   MonitoredSALReply
            [simple   MonitoredSAL('cBusOptions')                   monitoredSAL    ]
        ]
        [*, *, 'true', 'false'      *StandardFormatStatusReply
            [simple   StandardFormatStatusReply                     reply           ]
        ]
        [*, *, *, 'true'            *ExtendedFormatStatusReply
            [simple   ExtendedFormatStatusReply                     reply           ]
        ]
        [*, 'true', *, *            *CALReply
            [simple   CALReply('cBusOptions', 'requestContext')     calReply        ]
        ]
    ]
]

[type CALReply(CBusOptions cBusOptions, RequestContext requestContext)
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
    [simple   CALData('requestContext')   calData                                                ]
]

[type MonitoredSAL(CBusOptions cBusOptions)
    [peek    byte     salType             ]
    [typeSwitch salType
        ['0x05' *LongFormSmartMode
            [reserved byte '0x05']
            [peek    uint 24     terminatingByte                                ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual  bit isUnitAddress '(terminatingByte & 0xff) == 0x00'      ]
            [optional UnitAddress            unitAddress     'isUnitAddress'    ]
            [optional BridgeAddress          bridgeAddress   '!isUnitAddress'   ]
            [simple   ApplicationIdContainer application                        ]
            [optional byte                   reservedByte    'isUnitAddress'    ]
            [validation 'isUnitAddress && reservedByte == 0x00 || !isUnitAddress' "invalid unit address"]
            [optional ReplyNetwork           replyNetwork       '!isUnitAddress']
            [optional SALData('application.applicationId')   salData            ]
        ]
        [*      *ShortFormBasicMode
            [peek     byte                   counts                             ]
            [optional uint 8                 bridgeCount     'counts != 0x00'   ]
            [optional uint 8                 networkNumber   'counts != 0x00'   ]
            [optional byte                   noCounts        'counts == 0x00'   ] // TODO: add validation that this is 0x00 when no bridge and network number are set
            [simple   ApplicationIdContainer application                        ]
            [optional SALData('application.applicationId')  salData             ]
        ]
    ]
]

[type Confirmation
    [simple   Alpha           alpha                                                     ]
    // TODO: seem like sometimes there are two alphas in a confirmation... check that
    [optional Alpha           secondAlpha                                               ]
    [simple  ConfirmationType confirmationType                                          ]
    [virtual bit              isSuccess 'confirmationType == ConfirmationType.CONFIRMATION_SUCCESSFUL'   ]
]

[enum byte ConfirmationType
    ['0x2E'    CONFIRMATION_SUCCESSFUL                  ] // "."
    ['0x23'    NOT_TRANSMITTED_TO_MANY_RE_TRANSMISSIONS ] // "#"
    ['0x24'    NOT_TRANSMITTED_CORRUPTION               ] // "$"
    ['0x25'    NOT_TRANSMITTED_SYNC_LOSS                ] // "%"
    ['0x27'    NOT_TRANSMITTED_TOO_LONG                 ] // "'"
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
    [simple NetworkRoute  networkRoute                              ]
    [simple UnitAddress   unitAddress                               ]
]

[type Checksum
    [simple byte value]
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
    [virtual    uint 5  numberOfStatusBytes '(coding == StatusCoding.BINARY_BY_THIS_SERIAL_INTERFACE || coding == StatusCoding.BINARY_BY_ELSEWHERE)?(statusHeader.numberOfCharacterPairs - 3):(0)']
    [virtual    uint 5  numberOfLevelInformation '(coding == StatusCoding.LEVEL_BY_THIS_SERIAL_INTERFACE || coding == StatusCoding.LEVEL_BY_ELSEWHERE)?((statusHeader.numberOfCharacterPairs - 3) / 2):(0)']
    [array      StatusByte
                        statusBytes
                            count
                            'numberOfStatusBytes'                   ]
    [array      LevelInformation
                        levelInformation
                            count
                            'numberOfLevelInformation'              ]
]

[type ExtendedStatusHeader
    [reserved   uint 3                 '0x7'                        ]
    [simple     uint 5  numberOfCharacterPairs                      ]
]

[enum byte StatusCoding
    ['0x00' BINARY_BY_THIS_SERIAL_INTERFACE     ]
    ['0x40' BINARY_BY_ELSEWHERE                 ]
    ['0x07' LEVEL_BY_THIS_SERIAL_INTERFACE      ]
    ['0x47' LEVEL_BY_ELSEWHERE                  ]
]

[type StatusByte
    [simple GAVState    gav3                                        ]
    [simple GAVState    gav2                                        ]
    [simple GAVState    gav1                                        ]
    [simple GAVState    gav0                                        ]
]

[type LevelInformation
    [peek    uint 16    raw                                         ]
    [virtual uint 4     nibble1 '(raw & 0xF000) >> 12'              ]
    [virtual uint 4     nibble2 '(raw & 0x0F00) >> 8'               ]
    [virtual uint 4     nibble3 '(raw & 0x00F0) >> 4'               ]
    [virtual uint 4     nibble4 '(raw & 0x000F) >> 0'               ]
    [virtual bit        isAbsent 'nibble1 == 0x0 && nibble2 == 0x0 && nibble3 == 0x0 && nibble4 == 0x0']
    [virtual bit        isCorruptedByNoise '!isAbsent && (((nibble1 < 0x5) || (nibble1 == 0x8) || (nibble1 == 0xC)) || ((nibble2 < 0x5) || (nibble2 == 0x8) || (nibble2 == 0xC)) || ((nibble3 < 0x5) || (nibble3 == 0x8) || (nibble3 == 0xC)) || ((nibble4 < 0x5) || (nibble4 == 0x8) || (nibble4 == 0xC)))']
    [virtual bit        isCorruptedByNoiseOrLevelsDiffer '!isAbsent && (((nibble1 == 0x7) || (nibble1 == 0xB) || (nibble1 > 0xC)) || ((nibble2 == 0x7) || (nibble2 == 0xB) || (nibble2 > 0xC)) || ((nibble3 == 0x7) || (nibble3 == 0xB) || (nibble3 > 0xC)) || ((nibble4 == 0x7) || (nibble4 == 0xB) || (nibble4 > 0xC)))']
    [virtual bit        isCorrupted 'isCorruptedByNoise || isCorruptedByNoiseOrLevelsDiffer']
    [typeSwitch isAbsent, isCorrupted
        ['true'     *Absent
            [reserved uint 16 '0x0000'                                      ]
        ]
        [*, 'true'  *Corrupted
            [simple  uint 4    corruptedNibble1]
            [simple  uint 4    corruptedNibble2]
            [simple  uint 4    corruptedNibble3]
            [simple  uint 4    corruptedNibble4]
        ]
        [*          *Normal
            [simple  LevelInformationNibblePair  pair1                      ]
            [simple  LevelInformationNibblePair  pair2                      ]
            [virtual uint 8  actualLevel 'pair2.nibbleValue << 4 | pair1.nibbleValue']
        ]
    ]
]

[enum uint 8 LevelInformationNibblePair(uint 4 nibbleValue)
    ['0x55' Value_F ['0xF']]
    ['0x56' Value_E ['0xE']]
    ['0x59' Value_D ['0xD']]
    ['0x5A' Value_C ['0xC']]
    ['0x65' Value_B ['0xB']]
    ['0x66' Value_A ['0xA']]
    ['0x69' Value_9 ['0x9']]
    ['0x6A' Value_8 ['0x8']]
    ['0x95' Value_7 ['0x7']]
    ['0x96' Value_6 ['0x6']]
    ['0x99' Value_5 ['0x5']]
    ['0x9A' Value_4 ['0x4']]
    ['0xA5' Value_3 ['0x3']]
    ['0xA6' Value_2 ['0x2']]
    ['0xA9' Value_1 ['0x1']]
    ['0xAA' Value_0 ['0x0']]
]

[enum uint 2 GAVState
    ['0' DOES_NOT_EXIST                                             ]
    ['1' ON                                                         ]
    ['2' OFF                                                        ]
    ['3' ERROR                                                      ]
]

[type NetworkProtocolControlInformation
    [reserved   uint 2  '0x0'           ]
    [simple     uint 3  stackCounter    ] // Number of bridges required to transmit information from source to destination
    [simple     uint 3  stackDepth      ] // Number of bridges required to complete the transmission from source to destination
]

[type RequestTermination
    [const      byte    cr  0x0D                                    ] // 0xD == "<cr>"
]

[type ResponseTermination
    [const      byte    cr  0x0D                                    ] // 0xD == "<cr>"
    [const      byte    lf  0x0A                                    ] // 0xA == "<lf>"
]
