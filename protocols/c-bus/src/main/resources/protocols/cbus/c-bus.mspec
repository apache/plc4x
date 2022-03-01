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
    [simple bit smart  ]
    [simple bit idmon  ]
    [simple bit exstat ]
    [simple bit monitor]
    [simple bit monall ]
    [simple bit pun    ]
    [simple bit pcn    ]
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
    ['0x09' NoAdditionalBridge    ['0']]
    ['0x12' OneAdditionalBridge   ['1']]
    ['0x1B' TwoAdditionalBridge   ['2']]
    ['0x24' ThreeAdditionalBridge ['3']]
    ['0x2D' FourAdditionalBridge  ['4']]
    ['0x36' FiveAdditionalBridge  ['5']]
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
    [peek     byte          peekAlpha                                                           ]
    [optional Alpha         alpha    '(peekAlpha >= 0x67) && (peekAlpha <= 0x7A)'               ] // Read if the peeked byte is between 'g' and 'z'
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
            [peek     byte          peekAlpha                                                  ]
            [optional Alpha         alpha         '(peekAlpha >= 0x67) && (peekAlpha <= 0x7A)' ] // Read if the peeked byte is between 'g' and 'z'
            [const    byte          cr            0xD                                          ] // 0xD == "<cr>"
        ]
        [         CBusPointToMultiPointCommandNormal
            [simple   Application   application                                                ]
            [reserved byte          '0x00'                                                     ]
            [simple   SALData       salData                                                    ]
            [optional Checksum      crc         'srchk'                                        ] // crc      is optional but mspec crc      isn't
            [peek     byte          peekAlpha                                                  ]
            [optional Alpha         alpha       '(peekAlpha >= 0x67) && (peekAlpha <= 0x7A)'   ] // Read if the peeked byte is between 'g' and 'z'
            [const    byte          cr          0xD                                            ] // 0xD == "<cr>"
        ]
    ]
]

[discriminatedType CBusPointToPointToMultipointCommand(bit srchk)
    [simple BridgeAddress bridgeAddress                                                         ]
    [simple NetworkRoute  networkRoute                                                          ]
    [peek    byte       peekedApplication                                                       ]
    [typeSwitch peekedApplication
        ['0xFF'   CBusCommandPointToPointToMultiPointStatus
            [reserved byte        '0xFF'                                                    ]
            [simple StatusRequest statusRequest                                             ]
            [optional Checksum    crc           'srchk'                                     ] // crc      is optional but mspec crc      isn't
            [peek     byte        peekAlpha                                                 ]
            [optional Alpha       alpha         '(peekAlpha >= 0x67) && (peekAlpha <= 0x7A)'] // Read if the peeked byte is between 'g' and 'z'
            [const    byte        cr            0xD                                         ] // 0xD == "<cr>"
        ]
        [         CBusCommandPointToPointToMultiPointNormal
            [simple   Application application                                               ]
            [simple   SALData     salData                                                   ]
            [optional Checksum    crc         'srchk'                                       ] // crc      is optional but mspec crc      isn't
            [peek     byte        peekAlpha                                                 ]
            [optional Alpha       alpha       '(peekAlpha >= 0x67) && (peekAlpha <= 0x7A)'  ] // Read if the peeked byte is between 'g' and 'z'
            [const    byte        cr          0xD                                           ] // 0xD == "<cr>"
        ]
    ]
]

[type Application
    [simple byte id]
]

[type CALData
    [simple  CALCommandTypeContainer commandTypeContainer                                   ]
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
            [simple Application application                                                 ]
            [simple uint 8      blockStart                                                  ]
            [array  byte        data        count 'commandTypeContainer.numBytes'           ]
        ]
        ['STATUS_EXTENDED' CALDataReplyStatusExtended(CALCommandTypeContainer commandTypeContainer)
            [simple uint 8      encoding                                                    ]
            [simple Application application                                                 ]
            [simple uint 8      blockStart                                                  ]
            [array  byte        data        count 'commandTypeContainer.numBytes'           ]
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

[type Reply
    [peek   byte magicByte]
    [typeSwitch magicByte
        ['0x0' CALReplyReply
            [simple CALReply isA]
        ]
        ['0x0' MonitoredSALReply
            [simple MonitoredSAL isA]
        ]
        ['0x0' ConfirmationReply
            [simple Confirmation isA]
        ]
        ['0x0' PowerUpReply
            [simple PowerUp isA]
        ]
        ['0x0' ParameterChangeReply
            [simple ParameterChange isA]
        ]
        ['0x0' ExclamationMarkReply
            [simple ExclamationMark isA]
        ]
    ]
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
    [const    byte      cr      0x0D                                                             ] // 0xD == "<cr>"
    [const    byte      lf      0x0A                                                             ] // 0xA == "<lf>"
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
// TODO: implement garbage reading
//    [array    byte        garbage   terminated  '0x2B'                              ] // "+"
    [const    byte        plus 0x02B                                                  ] // 0xD == "<cr>"
    [const    byte        cr   0x0D                                                   ] // 0xD == "<cr>"
    [const    byte        lf   0x0A                                                   ] // 0xA == "<lf>"
]

[type ParameterChange
    [const    byte        specialChar1      0x3D                                    ] // "="
    [const    byte        specialChar2      0x3D                                    ] // "="
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