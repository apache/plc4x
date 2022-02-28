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

[discriminatedType CBusCommand
    [const  byte     initiator '0x5C'   ] // 0x5C == "/"
    [simple CBusHeader header           ]
    [typeSwitch 'header.destinationAddressType'
        ['PointToPointToMultiPoint' CBusCommandPointToPointToMultiPoint
            [simple CBusPointToPointToMultipointCommand command]
        ]
        ['PointToMultiPoint'        CBusCommandPointToMultiPoint
            [simple CBusPointToMultiPointCommand        command]
        ]
        ['PointToPoint'             CBusCommandPointToPoint
            [simple CBusPointToPointCommand             command]
        ]
    ]
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

[discriminatedType CBusPointToPointCommand
    [peek    uint 16     bridgeAddressCountPeek ]
    [typeSwitch 'bridgeAddressCountPeek && 0x00FF'
        ['0x0000' CBusPointToPointCommandDirect
            [simple UnitAddress   unitAddress                                                   ]
        ]
        ['*'      CBusPointToPointCommandIndirect
            [simple BridgeAddress bridgeAddress                                                 ]
            [simple NetworkRoute  networkRoute                                                  ]
            [simple UnitAddress   unitAddress                                                   ]
        ]
    ]
    [simple CALData calData                                                                     ]
    [optional Checksum      checksum                                                            ] // TODO: checksum is optional but mspec checksum isn't
    [optional Alpha          alpha                                                               ]
    [const    byte        cr '0xD'                                                            ] // 0xD == "<cr>"
]

[discriminatedType CBusPointToMultiPointCommand
    [peek    byte     application             ]
    [typeSwitch 'application'
        ['0xFF'   CBusPointToMultiPointCommandStatus
            [reserved byte                 '0xFF'                                              ]
            [reserved byte                 '0x00'                                              ]
            [simple StatusRequest   statusRequest                                              ]
            [optional Checksum      checksum                                                   ] // TODO: checksum is optional but mspec checksum isn't
            [optional Alpha          alpha                                                      ]
            [const    byte        cr '0xD'                                                   ] // 0xD == "<cr>"
        ]
        ['*'      CBusPointToMultiPointCommandNormal
            [simple Application          application                                                  ]
            [reserved byte                 '0x00'                                              ]
            [simple SALData         salData                                                    ]
            [optional Checksum      checksum                                                   ] // TODO: checksum is optional but mspec checksum isn't
            [optional Alpha          alpha                                                      ]
            [const    byte        cr '0xD'                                                   ] // 0xD == "<cr>"
        ]
    ]
]

[discriminatedType CBusCommandPointToPointToMultiPoint
    [simple BridgeAddress bridgeAddress                                                 ]
    [simple NetworkRoute  networkRoute                                                  ]
    [peek    byte     application             ]
    [typeSwitch 'application'
            ['0xFF'   CBusCommandPointToPointToMultiPointStatus
                [reserved byte                 '0xFF'                                              ]
                [reserved byte                 '0x00'                                              ]
                [simple StatusRequest   statusRequest                                              ]
                [optional Checksum      checksum                                                   ] // TODO: checksum is optional but mspec checksum isn't
                [optional Alpha          alpha                                                      ]
                [const    byte        cr '0xD'                                                   ] // 0xD == "<cr>"
            ]
            ['*'      CBusCommandPointToPointToMultiPointNormal
                [simple Application          application                                                  ]
                [reserved byte                 '0x00'                                              ]
                [simple SALData         salData                                                    ]
                [optional Checksum      checksum                                                   ] // TODO: checksum is optional but mspec checksum isn't
                [optional Alpha          alpha                                                      ]
                [const    byte        cr '0xD'                                                   ] // 0xD == "<cr>"
            ]
        ]
]

[type Application
    [simple byte id]
]

[type CALData
    // TODO: implement me
]

[type StatusRequest
    [peek    byte     type             ]
    [typeSwitch 'type'
        ['0x7A' StatusRequestBinaryState
            [reserved byte                 '0x7A'                                              ]
            [simple byte          application                                                  ]
            [reserved byte                 '0x00'                                              ]
        ]
        ['0x73' StatusRequestLevel
            [reserved byte                 '0x73'                                              ]
            [reserved byte                 '0x07'                                              ]
            [simple byte          application                                                  ]
            [simple byte          startingGroupAddressLabel                                    ]
            [validation           'startingGroupAddressLabel == 0x00
                                || startingGroupAddressLabel == 0x20
                                || startingGroupAddressLabel == 0x40
                                || startingGroupAddressLabel == 0x60
                                || startingGroupAddressLabel == 0x80
                                || startingGroupAddressLabel == 0xA0
                                || startingGroupAddressLabel == 0xC0
                                || startingGroupAddressLabel == 0xE0'                          ]
        ]
    ]
]

[type SALData
    // TODO: implement me
]

[type Reply
    [peek   byte    magicByte]
    [typeSwitch 'magicByte'
        ['??' CALReply]
        ['??' MonitoredSAL]
        ['??' Confirmation]
        ['??' PowerUp]
        ['??' ParameterChange]
        ['??' ExclamationMark]
    ]
]

[type CALReply
    [peek    byte     type             ]
    [typeSwitch 'type'
        ['0x86' CALReplyLong
            [reserved byte '0x86']
            [peek    uint 24     terminatingByte                        ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual bit isUnitAddress 'terminatingByte & 0xff == 0x00' ]
            [optional   UnitAddress
                         unitAddress     'isUnitAddress'                ]
            [optional   BridgeAddress
                         bridgeAddress   '!isUnitAddress'               ]
            [simple     SerialInterfaceAddress
                         serialInterfaceAddress                         ]
            [optional   byte    reservedByte    'isUnitAddress'         ]
            [validation 'isUnitAddress && reservedByte == 0x00 || !isUnitAddress']
            [optional   ReplyNetwork            '!isUnitAddress'        ]
        ]
        ['*'    CALReplyShort
        ]
    ]
    [simple CALData calData]
    [checksum checksum]
    [const    byte        cr '0x0D'                                                   ] // 0xD == "<cr>"
    [const    byte        lf '0x0A'                                                   ] // 0xA == "<lf>"
]

[type BridgeCount
    [simple uint 8 count]
]

[type NetworkNumber
    [simple uint 8 number]
]

[type MonitoredSAL
    [peek    byte     type             ]
    [typeSwitch 'type'
        ['0x05' MonitoredSALLongFormSmartMode
            [reserved byte '0x05']
            [peek    uint 24     terminatingByte                        ]
            // TODO: this should be subSub type but mspec doesn't support that yet directly
            [virtual bit isUnitAddress 'terminatingByte & 0xff == 0x00' ]
            [optional   UnitAddress
                         unitAddress     'isUnitAddress'                ]
            [optional   BridgeAddress
                         bridgeAddress   '!isUnitAddress'               ]
            [simple     SerialInterfaceAddress
                         serialInterfaceAddress                         ]
            [optional   byte    reservedByte    'isUnitAddress'         ]
            [validation 'isUnitAddress && reservedByte == 0x00 || !isUnitAddress']
            [optional   ReplyNetwork            '!isUnitAddress'        ]
        ]
        ['*' MonitoredSALShortFormBasicMode
            [peek    byte     counts                                        ]
            [optional BridgeCount   bridgeCount     'counts != 0x00'        ]
            [optional NetworkNumber networkNumber   'counts != 0x00'        ]
            [optional byte    noCounts              'counts == 0x00'        ] // TODO: add validation that this is 0x00 when no bridge and network number are set
            [simple Application application                                 ]
        ]
    ]
    [optional SALData salData                                               ]
    [checksum checksum]
    [const    byte        cr '0x0D'                                                   ] // 0xD == "<cr>"
    [const    byte        lf '0x0A'                                                   ] // 0xA == "<lf>"
]

[type Confirmation
    [simple Alpha alpha]
    [dicriminator   byte type]
    [typeSwitch 'type'
        ['.'    ConfirmationSuccessful]
        ['#'    NotTransmittedToManyReTransmissions]
        ['$'    NotTransmittedCorruption]
        ['%'    NotTransmittedSyncLoss]
        ['''    NotTransmittedTooLong]
    ]
]

[type PowerUp
    // TODO: skip potential garbage as first reserved might be wrong
    [const    byte        something1     '+']
    [const    byte        something2    '+']
    [const    byte        cr '0x0D'                                                   ] // 0xD == "<cr>"
    [const    byte        lf '0x0A'                                                   ] // 0xA == "<lf>"
]

[type ParameterChange
    [const    byte        something1    '=']
    [const    byte        something2    '=']
    [const    byte        cr '0x0D'                                                   ] // 0xD == "<cr>"
    [const    byte        lf '0x0A'                                                   ] // 0xA == "<lf>"
]

[type ExclamationMark
    // TODO: implement me
]

[type ReplyNetwork
     [simple RouteType     routeType                                                    ]
     [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
     [simple UnitAddress    unitAddress                                                 ]
]
