/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

[discriminatedType CBusCommand(bit srchk)
    [const  byte       initiator 0x5C   ] // 0x5C == "/"
    [simple CBusHeader header           ]
    [typeSwitch header.destinationAddressType
        ['PointToPointToMultiPoint' CBusCommandPointToPointToMultiPoint
            [simple CBusPointToPointToMultipointCommand('srchk') command]
        ]
        ['PointToMultiPoint'        CBusCommandPointToMultiPoint
            [simple CBusPointToMultiPointCommand('srchk')        command]
        ]
        ['PointToPoint'             CBusCommandPointToPoint
            [simple CBusPointToPointCommand('srchk')             command]
        ]
    ]
]

// TODO: check if that can be used in combination with srchk
[type CBusOptions
    [simple bit connect]
    [simple bit smart]
    [simple bit idmon]
    [simple bit exstat]
    [simple bit monitor]
    [simple bit monall]
    [simple bit pun]
    [simple bit pcn]
]

[type CBusHeader
    [simple   PriorityClass          priorityClass         ]
    [reserved bit                    'false'               ] // Reserved for internal C-Bus management purposes
    [reserved uint 2                 '0'                   ] // Reserved for internal C-Bus management purposes
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
]

[type NetworkRoute
    [simple RouteType     routeType                                                    ]
    [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
]

[enum byte RouteType(uint 3 additionalBridges)
    ['0x00' NoBridgeAtAll         ['0']]
    ['0x09' NoAdditionalBridge    ['1']]
    ['0x12' OneAdditionalBridge   ['2']]
    ['0x1B' TwoAdditionalBridge   ['3']]
    ['0x24' ThreeAdditionalBridge ['4']]
    ['0x2D' FourAdditionalBridge  ['4']]
    ['0x36' FiveAdditionalBridge  ['4']]
]

[discriminatedType CBusPointToPointCommand(bit srchk)
    [peek    uint 16     bridgeAddressCountPeek ]
    [virtual bit         isDirect  '(bridgeAddressCountPeek & 0x00FF) == 0x0000']
    [typeSwitch isDirect
        ['true'  CBusPointToPointCommandDirect
            [simple UnitAddress   unitAddress                                                   ]
        ]
        ['false' CBusPointToPointCommandIndirect
            [simple BridgeAddress bridgeAddress                                                 ]
            [simple NetworkRoute  networkRoute                                                  ]
            [simple UnitAddress   unitAddress                                                   ]
        ]
    ]
    [simple   CALData calData                                                                   ]
    [optional Checksum      crc      'srchk'                                                    ] // checksum is optional but mspec checksum isn't
    [optional Alpha         alpha    'curPos <= 123123'                                  ] // TODO:Read if there's still two bytes left (was lengthInBytes but that didn't compile here)
    [const    byte          cr       0xD                                                        ] // 0xD == "<cr>"
]

[discriminatedType CBusPointToMultiPointCommand(bit srchk)
    [peek    byte     peekedApplication             ]
    [typeSwitch peekedApplication
        ['0xFF'   CBusPointToMultiPointCommandStatus
            [reserved byte          '0xFF'                                                     ]
            [reserved byte          '0x00'                                                     ]
            [simple   StatusRequest statusRequest                                              ]
            [optional Checksum      crc           'srchk'                                      ] // checksum is optional but mspec checksum isn't
            [optional Alpha         alpha         'curPos <= 123123'                    ] // TODO:Read if there's still two bytes left (was lengthInBytes but that didn't compile here)
            [const    byte          cr            0xD                                          ] // 0xD == "<cr>"
        ]
        [         CBusPointToMultiPointCommandNormal
            [simple   Application   application                                                ]
            [reserved byte          '0x00'                                                     ]
            [simple   SALData       salData                                                    ]
            [optional Checksum      crc         'srchk'                                        ] // crc      is optional but mspec crc      isn't
            [optional Alpha         alpha       'curPos <= 123123'                      ] // TODO:Read if there's still two bytes left (was lengthInBytes but that didn't compile here)
            [const    byte          cr          0xD                                            ] // 0xD == "<cr>"
        ]
    ]
]

[discriminatedType CBusPointToPointToMultipointCommand(bit srchk)
    [simple BridgeAddress bridgeAddress                                                        ]
    [simple NetworkRoute  networkRoute                                                         ]
    [peek    byte       peekedApplication                                                      ]
    [typeSwitch peekedApplication
            ['0xFF'   CBusCommandPointToPointToMultiPointStatus
                [reserved byte        '0xFF'                                                   ]
                [reserved byte        '0x00'                                                   ]
                [simple StatusRequest statusRequest                                            ]
                [optional Checksum    crc           'srchk'                                    ] // crc      is optional but mspec crc      isn't
                [optional Alpha       alpha         'curPos <= 123123'                  ] // TODO:Read if there's still two bytes left (was lengthInBytes but that didn't compile here)
                [const    byte        cr            0xD                                        ] // 0xD == "<cr>"
            ]
            [         CBusCommandPointToPointToMultiPointNormal
                [simple Application   application                                              ]
                [reserved byte        '0x00'                                                   ]
                [simple SALData       salData                                                  ]
                [optional Checksum    crc      'srchk'                                         ] // crc      is optional but mspec crc      isn't
                [optional Alpha       alpha    'curPos <= 123123'                       ] // TODO:Read if there's still two bytes left (was lengthInBytes but that didn't compile here)
                [const    byte        cr       0xD                                             ] // 0xD == "<cr>"
            ]
        ]
]

[type Application
    [simple byte id]
]

[type CALData
    [simple CommandHeader commandHeader]
    [typeSwitch commandHeader.value
        ['0x08' CALDataRequestReset
        ]
        ['0x1A' CALDataRequestRecall
            [simple uint 8 paramNo]
            [simple uint 8 count]
        ]
        ['0x21' CALDataRequestIdentify
            [simple byte attribute]
        ]
        ['0x09' CALDataRequestGetStatus
            [simple uint 8 paramNo]
            [simple uint 8 count]
        ]
        ['0x80' CALDataReplyReply
            // TODO: how to parse this?
        ]
        ['0x32' CALDataReplyAcknowledge
            [simple uint 8 paramNo]
            [simple uint 8 code]
        ]
        ['0xC0' CALDataReplyStatus
            // TODO: how to parse this?
        ]
        ['0xE0' CALDataReplyStatusExtended
            // TODO: how to parse this?
        ]
    ]
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
                                || startingGroupAddressLabel == 0xE0'                 "invalid label"]
        ]
    ]
]

[type SALData
    [simple CommandHeader commandHeader]
    [typeSwitch commandHeader.value
        ['0x01' SALDataOff
            [simple byte group]
        ]
        ['0x79' SALDataOn
            [simple byte group]
        ]
        ['0x0' SALDataRampToLevel
            [simple byte groupLevel]
        ]
        ['0x09' SALDataTerminateRamp
            [simple byte group]
        ]
    ]
]

[type CommandHeader
    [simple byte value]
]

[type Reply
    [peek   byte magicByte]
    [typeSwitch magicByte
        ['0x0' CALReplyReply
        ]
        ['0x0' MonitoredSALReply
        ]
        ['0x0' ConfirmationReply
        ]
        ['0x0' PowerUpReply
        ]
        ['0x0' ParameterChangeReply
        ]
        ['0x0' ExclamationMarkReply
        ]
    ]
]

[type CALReply
    [peek    byte     calType             ]
    [typeSwitch calType
        ['0x86' CALReplyLong
            [reserved   byte                   '0x86'                                          ]
            [peek       uint 24                terminatingByte                                 ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual    bit                    isUnitAddress   '(terminatingByte & 0xff) == 0x00']
            [optional   UnitAddress            unitAddress     'isUnitAddress'                 ]
            [optional   BridgeAddress          bridgeAddress   '!isUnitAddress'                ]
            [simple     SerialInterfaceAddress serialInterfaceAddress                          ]
            [optional   byte                   reservedByte    'isUnitAddress'                 ]
            [validation                        'isUnitAddress && reservedByte == 0x00 || !isUnitAddress' "wrong reservedByte"]
            [optional   ReplyNetwork           replyNetwork   '!isUnitAddress'                 ]
        ]
        [       CALReplyShort
        ]
    ]
    [simple   CALData   calData                                                                ]
    //[checksum byte crc   '0x00'                                                                ] // TODO: Fix this
    [const    byte      cr      0x0D                                                           ] // 0xD == "<cr>"
    [const    byte      lf      0x0A                                                           ] // 0xA == "<lf>"
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
            [peek    byte     counts                                        ]
            [optional BridgeCount   bridgeCount     'counts != 0x00'        ]
            [optional NetworkNumber networkNumber   'counts != 0x00'        ]
            [optional byte    noCounts              'counts == 0x00'        ] // TODO: add validation that this is 0x00 when no bridge and network number are set
            [simple Application application                                 ]
        ]
    ]
    [optional SALData salData                                               ]
    //[checksum byte crc   '0x00'                                                                ] // TODO: Fix this
    [const    byte        cr 0x0D                                                     ] // 0xD == "<cr>"
    [const    byte        lf 0x0A                                                     ] // 0xA == "<lf>"
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
    ]
]

[type PowerUp
    // TODO: skip potential garbage as first reserved might be wrong
    [const    byte        something1     0x2B] // "+"
    [const    byte        something2     0x2B] // "+"
    [const    byte        cr 0x0D                                                   ] // 0xD == "<cr>"
    [const    byte        lf 0x0A                                                   ] // 0xA == "<lf>"
]

[type ParameterChange
    [const    byte        something1    0x3D] // "="
    [const    byte        something2    0x3D] // "="
    [const    byte        cr 0x0D                                                   ] // 0xD == "<cr>"
    [const    byte        lf 0x0A                                                   ] // 0xA == "<lf>"
]

[type ExclamationMark
    // TODO: implement me
]

[type ReplyNetwork
     [simple RouteType     routeType                                                    ]
     [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
     [simple UnitAddress   unitAddress                                                  ]
]

[type Checksum
    [simple byte crc]
]