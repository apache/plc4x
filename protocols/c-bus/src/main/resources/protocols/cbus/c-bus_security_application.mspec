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

[type SecurityData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsSecurityCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  SecurityCommandTypeContainer commandTypeContainer                                   ]
    [virtual SecurityCommandType          commandType          'commandTypeContainer.commandType']
    [simple byte argument                                                               ]
    [typeSwitch commandType, argument
        ['ON', '0x80'       *SystemArmedDisarmed
            [simple SecurityArmCode armCodeType ]
        ]
        ['OFF', '0x80'      *SystemDisarmed
        ]
        ['EVENT', '0x81'    *ExitDelayStarted
        ]
        ['EVENT', '0x82'    *EntryDelayStarted
        ]
        ['ON', '0x83'       *AlarmOn
        ]
        ['OFF', '0x83'      *AlarmOff
        ]
        ['ON', '0x84'       *TamperOn
        ]
        ['OFF', '0x84'      *TamperOff
        ]
        ['ON', '0x85'       *PanicActivated
        ]
        ['OFF', '0x85'      *PanicCleared
        ]
        ['EVENT', '0x86'    *ZoneUnsealed
            [simple uint 8 zoneNumber]
        ]
        ['EVENT', '0x87'    *ZoneSealed
            [simple uint 8 zoneNumber]
        ]
        ['EVENT', '0x88'    *ZoneOpen
            [simple uint 8 zoneNumber]
        ]
        ['EVENT', '0x89'    *ZoneShort
            [simple uint 8 zoneNumber]
        ]
        ['EVENT', '0x89'    *ZoneIsolated
            [simple uint 8 zoneNumber]
        ]
        ['ON', '0x8B'       *LowBatteryDetected
        ]
        ['OFF', '0x8B'      *LowBatteryCorrected
        ]
        ['EVENT', '0x8C'    *LowBatteryCharging
            [simple  byte startStop                     ]
            [virtual bit  chargeStopped 'startStop==0x00'    ]
            [virtual bit  chargeStarted 'startStop>0xFE'     ]
        ]
        ['EVENT', '0x8D'    *ZoneName
            [simple uint 8      zoneNumber  ]
            [simple string 88   zoneName    ]
        ]
        ['EVENT', '0x8E'    *StatusReport1
            [simple  SecurityArmCode    armCodeType                                                 ]
            [simple  TamperStatus       tamperStatus                                                ]
            [simple  PanicStatus        panicStatus                                                 ]
            [array   ZoneStatus         zoneStatus        count '32'                                ]
        ]
        ['EVENT', '0x8F'    *StatusReport2
            [array   ZoneStatus         zoneStatus        count '48'                                ]
        ]
        ['EVENT', '0x90'    *PasswordEntryStatus
            [simple  byte code  ]
            [virtual bit  isPasswordEntrySucceeded      'code == 0x01']
            [virtual bit  isPasswordEntryFailed         'code == 0x02']
            [virtual bit  isPasswordEntryDisabled       'code == 0x03']
            [virtual bit  isPasswordEntryEnabledAgain   'code == 0x04']
            [virtual bit  isReserved                    'code >= 0x05']
        ]
        ['ON', '0x91'       *MainsFailure
        ]
        ['OFF', '0x91'      *MainsRestoredOrApplied
        ]
        ['EVENT', '0x92'    *ArmReadyNotReady
            [simple uint 8      zoneNumber  ]
        ]
        ['EVENT', '0x93'    *CurrentAlarmType
        ]
        ['ON', '0x94'       *LineCutAlarmRaised
        ]
        ['OFF', '0x94'      *LineCutAlarmCleared
        ]
        ['ON', '0x95'       *ArmFailedRaised
        ]
        ['OFF', '0x95'      *ArmFailedCleared
        ]
        ['ON', '0x96'       *FireAlarmRaised
        ]
        ['OFF', '0x96'      *FireAlarmCleared
        ]
        ['ON', '0x97'       *GasAlarmRaised
        ]
        ['OFF', '0x97'      *GasAlarmCleared
        ]
        ['ON', '0x98'       *OtherAlarmRaised
        ]
        ['OFF', '0x98'      *OtherAlarmCleared
        ]
        ['EVENT', '0xA0'    *Status1Request
        ]
        ['EVENT', '0xA1'    *Status2Request
        ]
        ['EVENT', '0xA2'    *ArmSystem
            [simple  byte armMode                                   ]
            [virtual bit  isReserved            'armMode == 0x00 || (armMode >= 0x05 && armMode <= 0xFE)'      ]
            [virtual bit  isArmToAwayMode       'armMode == 0x01'   ]
            [virtual bit  isArmToNightMode      'armMode == 0x02'   ]
            [virtual bit  isArmToDayMode        'armMode == 0x03'   ]
            [virtual bit  isArmToVacationMode   'armMode == 0x04'   ]
            [virtual bit  isArmToHighestLevelOfProtection   'armMode > 0xFE'   ]
        ]
        ['ON', '0xA3'       *RaiseTamper
        ]
        ['OFF', '0xA3'      *DropTamper
        ]
        ['ON', '0xA4'       *RaiseAlarm
        ]
        ['EVENT', '0xA5'    *EmulatedKeypad
            [simple  byte key                                        ]
            [virtual bit  isAscii       'key >= 0x00 && key <= 0x7F' ]
            [virtual bit  isCustom      'key >= 0x80'                ]
            [virtual bit  isEnter       'key == 0x0D'                ]
            [virtual bit  isShift       'key == 0x80'                ]
            [virtual bit  isPanic       'key == 0x81'                ]
            [virtual bit  isFire        'key == 0x82'                ]
            [virtual bit  isARM         'key == 0x83'                ]
            [virtual bit  isAway        'key == 0x84'                ]
            [virtual bit  isNight       'key == 0x85'                ]
            [virtual bit  isDay         'key == 0x86'                ]
            [virtual bit  isVacation    'key == 0x87'                ]
        ]
        ['ON', '0xA6'       *DisplayMessage(SecurityCommandTypeContainer commandTypeContainer)
            [simple vstring '(commandTypeContainer.numBytes-1)*8' message                           ]
        ]
        ['EVENT', '0xA7'    *RequestZoneName
            [simple uint 8      zoneNumber  ]
        ]
        ['OFF'                                  *Off(SecurityCommandTypeContainer commandTypeContainer)
            [array  byte data        count 'commandTypeContainer.numBytes-1'                    ]
        ]
        ['ON'                                   *On(SecurityCommandTypeContainer commandTypeContainer)
            [array  byte data        count 'commandTypeContainer.numBytes-1'                    ]
        ]
        ['EVENT'                                *Event(SecurityCommandTypeContainer commandTypeContainer)
            [array  byte data        count 'commandTypeContainer.numBytes-1'                    ]
        ]
    ]
]

[enum uint 8 SecurityCommandTypeContainer(SecurityCommandType commandType, uint 5 numBytes)
    ['0x00' SecurityCommandOff_0Bytes                    ['OFF',    '0']]
    ['0x01' SecurityCommandOff_1Bytes                    ['OFF',    '1']]
    ['0x02' SecurityCommandOff_2Bytes                    ['OFF',    '2']]
    ['0x03' SecurityCommandOff_3Bytes                    ['OFF',    '3']]
    ['0x04' SecurityCommandOff_4Bytes                    ['OFF',    '4']]
    ['0x05' SecurityCommandOff_5Bytes                    ['OFF',    '5']]
    ['0x06' SecurityCommandOff_6Bytes                    ['OFF',    '6']]
    ['0x07' SecurityCommandOff_7Bytes                    ['OFF',    '7']]
    ['0x08' SecurityCommandEvent_0Bytes                  ['EVENT',  '0']]
    ['0x09' SecurityCommandEvent_1Bytes                  ['EVENT',  '1']]
    ['0x0A' SecurityCommandEvent_2Bytes                  ['EVENT',  '2']]
    ['0x0B' SecurityCommandEvent_3Bytes                  ['EVENT',  '3']]
    ['0x0C' SecurityCommandEvent_4Bytes                  ['EVENT',  '4']]
    ['0x0D' SecurityCommandEvent_5Bytes                  ['EVENT',  '5']]
    ['0x0E' SecurityCommandEvent_6Bytes                  ['EVENT',  '6']]
    ['0x0F' SecurityCommandEvent_7Bytes                  ['EVENT',  '7']]
    ['0x78' SecurityCommandOn_0Bytes                     ['ON',     '0']]
    ['0x79' SecurityCommandOn_1Bytes                     ['ON',     '1']]
    ['0x7A' SecurityCommandOn_2Bytes                     ['ON',     '2']]
    ['0x7B' SecurityCommandOn_3Bytes                     ['ON',     '3']]
    ['0x7C' SecurityCommandOn_4Bytes                     ['ON',     '4']]
    ['0x7D' SecurityCommandOn_5Bytes                     ['ON',     '5']]
    ['0x7E' SecurityCommandOn_6Bytes                     ['ON',     '6']]
    ['0x7F' SecurityCommandOn_7Bytes                     ['ON',     '7']]
    ['0x80' SecurityCommandLongOff_0Bytes                ['OFF',    '8']]
    ['0x81' SecurityCommandLongOff_1Bytes                ['OFF',    '1']]
    ['0x82' SecurityCommandLongOff_2Bytes                ['OFF',    '2']]
    ['0x83' SecurityCommandLongOff_3Bytes                ['OFF',    '3']]
    ['0x84' SecurityCommandLongOff_4Bytes                ['OFF',    '4']]
    ['0x85' SecurityCommandLongOff_5Bytes                ['OFF',    '5']]
    ['0x86' SecurityCommandLongOff_6Bytes                ['OFF',    '6']]
    ['0x87' SecurityCommandLongOff_7Bytes                ['OFF',    '7']]
    ['0x88' SecurityCommandLongOff_8Bytes                ['OFF',    '8']]
    ['0x89' SecurityCommandLongOff_9Bytes                ['OFF',    '9']]
    ['0x8A' SecurityCommandLongOff_10Bytes               ['OFF',   '10']]
    ['0x8B' SecurityCommandLongOff_11Bytes               ['OFF',   '11']]
    ['0x8C' SecurityCommandLongOff_12Bytes               ['OFF',   '12']]
    ['0x8D' SecurityCommandLongOff_13Bytes               ['OFF',   '13']]
    ['0x8E' SecurityCommandLongOff_14Bytes               ['OFF',   '14']]
    ['0x8F' SecurityCommandLongOff_15Bytes               ['OFF',   '15']]
    ['0x90' SecurityCommandLongOff_16Bytes               ['OFF',   '16']]
    ['0x91' SecurityCommandLongOff_17Bytes               ['OFF',   '17']]
    ['0x92' SecurityCommandLongOff_18Bytes               ['OFF',   '18']]
    ['0x93' SecurityCommandLongOff_19Bytes               ['OFF',   '19']]
    ['0x94' SecurityCommandLongOff_20Bytes               ['OFF',   '20']]
    ['0x95' SecurityCommandLongOff_21Bytes               ['OFF',   '21']]
    ['0x96' SecurityCommandLongOff_22Bytes               ['OFF',   '22']]
    ['0x97' SecurityCommandLongOff_23Bytes               ['OFF',   '23']]
    ['0x98' SecurityCommandLongOff_24Bytes               ['OFF',   '24']]
    ['0x99' SecurityCommandLongOff_25Bytes               ['OFF',   '25']]
    ['0x9A' SecurityCommandLongOff_26Bytes               ['OFF',   '26']]
    ['0x9B' SecurityCommandLongOff_27Bytes               ['OFF',   '27']]
    ['0x9C' SecurityCommandLongOff_28Bytes               ['OFF',   '28']]
    ['0x9D' SecurityCommandLongOff_29Bytes               ['OFF',   '29']]
    ['0x9E' SecurityCommandLongOff_30Bytes               ['OFF',   '30']]
    ['0x9F' SecurityCommandLongOff_31Bytes               ['OFF',   '31']]
    ['0xA0' SecurityCommandLongEvent_0Bytes              ['EVENT',  '0']]
    ['0xA1' SecurityCommandLongEvent_1Bytes              ['EVENT',  '1']]
    ['0xA2' SecurityCommandLongEvent_2Bytes              ['EVENT',  '2']]
    ['0xA3' SecurityCommandLongEvent_3Bytes              ['EVENT',  '3']]
    ['0xA4' SecurityCommandLongEvent_4Bytes              ['EVENT',  '4']]
    ['0xA5' SecurityCommandLongEvent_5Bytes              ['EVENT',  '5']]
    ['0xA6' SecurityCommandLongEvent_6Bytes              ['EVENT',  '6']]
    ['0xA7' SecurityCommandLongEvent_7Bytes              ['EVENT',  '7']]
    ['0xA8' SecurityCommandLongEvent_8Bytes              ['EVENT',  '8']]
    ['0xA9' SecurityCommandLongEvent_9Bytes              ['EVENT',  '9']]
    ['0xAA' SecurityCommandLongEvent_10Bytes             ['EVENT', '10']]
    ['0xAB' SecurityCommandLongEvent_11Bytes             ['EVENT', '11']]
    ['0xAC' SecurityCommandLongEvent_12Bytes             ['EVENT', '12']]
    ['0xAD' SecurityCommandLongEvent_13Bytes             ['EVENT', '13']]
    ['0xAE' SecurityCommandLongEvent_14Bytes             ['EVENT', '14']]
    ['0xAF' SecurityCommandLongEvent_15Bytes             ['EVENT', '15']]
    ['0xB0' SecurityCommandLongEvent_16Bytes             ['EVENT', '16']]
    ['0xB1' SecurityCommandLongEvent_17Bytes             ['EVENT', '17']]
    ['0xB2' SecurityCommandLongEvent_18Bytes             ['EVENT', '18']]
    ['0xB3' SecurityCommandLongEvent_19Bytes             ['EVENT', '19']]
    ['0xB4' SecurityCommandLongEvent_20Bytes             ['EVENT', '20']]
    ['0xB5' SecurityCommandLongEvent_21Bytes             ['EVENT', '21']]
    ['0xB6' SecurityCommandLongEvent_22Bytes             ['EVENT', '22']]
    ['0xB7' SecurityCommandLongEvent_23Bytes             ['EVENT', '23']]
    ['0xB8' SecurityCommandLongEvent_24Bytes             ['EVENT', '24']]
    ['0xB9' SecurityCommandLongEvent_25Bytes             ['EVENT', '25']]
    ['0xBA' SecurityCommandLongEvent_26Bytes             ['EVENT', '26']]
    ['0xBB' SecurityCommandLongEvent_27Bytes             ['EVENT', '27']]
    ['0xBC' SecurityCommandLongEvent_28Bytes             ['EVENT', '28']]
    ['0xBD' SecurityCommandLongEvent_29Bytes             ['EVENT', '29']]
    ['0xBE' SecurityCommandLongEvent_30Bytes             ['EVENT', '30']]
    ['0xBF' SecurityCommandLongEvent_31Bytes             ['EVENT', '31']]
    ['0xE0' SecurityCommandLongOn_0Bytes                 ['ON',     '0']]
    ['0xE1' SecurityCommandLongOn_1Bytes                 ['ON',     '1']]
    ['0xE2' SecurityCommandLongOn_2Bytes                 ['ON',     '2']]
    ['0xE3' SecurityCommandLongOn_3Bytes                 ['ON',     '3']]
    ['0xE4' SecurityCommandLongOn_4Bytes                 ['ON',     '4']]
    ['0xE5' SecurityCommandLongOn_5Bytes                 ['ON',     '5']]
    ['0xE6' SecurityCommandLongOn_6Bytes                 ['ON',     '6']]
    ['0xE7' SecurityCommandLongOn_7Bytes                 ['ON',     '7']]
    ['0xE8' SecurityCommandLongOn_8Bytes                 ['ON',     '8']]
    ['0xE9' SecurityCommandLongOn_9Bytes                 ['ON',     '9']]
    ['0xEA' SecurityCommandLongOn_10Bytes                ['ON',    '10']]
    ['0xEB' SecurityCommandLongOn_11Bytes                ['ON',    '11']]
    ['0xEC' SecurityCommandLongOn_12Bytes                ['ON',    '12']]
    ['0xED' SecurityCommandLongOn_13Bytes                ['ON',    '13']]
    ['0xEE' SecurityCommandLongOn_14Bytes                ['ON',    '14']]
    ['0xEF' SecurityCommandLongOn_15Bytes                ['ON',    '15']]
    ['0xF0' SecurityCommandLongOn_16Bytes                ['ON',    '16']]
    ['0xF1' SecurityCommandLongOn_17Bytes                ['ON',    '17']]
    ['0xF2' SecurityCommandLongOn_18Bytes                ['ON',    '18']]
    ['0xF3' SecurityCommandLongOn_19Bytes                ['ON',    '19']]
    ['0xF4' SecurityCommandLongOn_20Bytes                ['ON',    '20']]
    ['0xF5' SecurityCommandLongOn_21Bytes                ['ON',    '21']]
    ['0xF6' SecurityCommandLongOn_22Bytes                ['ON',    '22']]
    ['0xF7' SecurityCommandLongOn_23Bytes                ['ON',    '23']]
    ['0xF8' SecurityCommandLongOn_24Bytes                ['ON',    '24']]
    ['0xF9' SecurityCommandLongOn_25Bytes                ['ON',    '25']]
    ['0xFA' SecurityCommandLongOn_26Bytes                ['ON',    '26']]
    ['0xFB' SecurityCommandLongOn_27Bytes                ['ON',    '27']]
    ['0xFC' SecurityCommandLongOn_28Bytes                ['ON',    '28']]
    ['0xFD' SecurityCommandLongOn_29Bytes                ['ON',    '29']]
    ['0xFE' SecurityCommandLongOn_30Bytes                ['ON',    '30']]
    ['0xFF' SecurityCommandLongOn_31Bytes                ['ON',    '31']]
]

[enum uint 4 SecurityCommandType
    ['0x00' OFF     ]
    ['0x01' ON      ]
    ['0x02' EVENT   ]
]

[type SecurityArmCode
    [simple  uint 8 code                                 ]
    [virtual bit    isDisarmed          'code == 0x00'   ]
    [virtual bit    isFullyArmed        'code == 0x01'   ]
    [virtual bit    isPartiallyArmed    'code == 0x02'   ]
    [virtual bit    isArmSubtype        'code >= 0x03 && code <= 0x7F'   ]
    [virtual bit    isReserved          'code > 0x7F'    ]
]

[type TamperStatus
    [simple  uint 8 status                                          ]
    [virtual bit    isNoTamper 'status == 0x00'                     ]
    [virtual bit    isReserved 'status >= 0x01 && status <= 0xFE']
    [virtual bit    isTamperActive 'status > 0xFE'                  ]
]

[type PanicStatus
    [simple  uint 8 status                                          ]
    [virtual bit    isNoPanic  'status == 0x00'                     ]
    [virtual bit    isReserved 'status >= 0x01 && status <= 0xFE'    ]
    [virtual bit    isPanicCurrentlyActive 'status > 0xFE'          ]
]

[type ZoneStatus
    [simple ZoneStatusTemp value]
]

// TODO: we can't use ZoneStatus directly as nobody used enums in list till now so we just wrap it
[enum uint 2 ZoneStatusTemp
    ['0x0' ZONE_SEALED      ]
    ['0x1' ZONE_UNSEALED    ]
    ['0x2' ZONE_OPEN        ]
    ['0x3' ZONE_SHORT       ]
]
