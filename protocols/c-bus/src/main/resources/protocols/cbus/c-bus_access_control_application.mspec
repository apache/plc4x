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

[type AccessControlData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsAccessControlCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  AccessControlCommandTypeContainer          commandTypeContainer                                   ]
    [virtual AccessControlCommandType                   commandType          'commandTypeContainer.commandType']
    [simple  byte                                       networkId                                              ]
    [simple  byte                                       accessPointId                                          ]
    [typeSwitch commandType
        ['VALID_ACCESS'             *ValidAccessRequest(AccessControlCommandTypeContainer commandTypeContainer)
            [simple   AccessControlDirection    accessControlDirection]
            [array    byte                      data          count 'commandTypeContainer.numBytes-3'          ]
        ]
        ['INVALID_ACCESS'           *InvalidAccessRequest(AccessControlCommandTypeContainer commandTypeContainer)
            [simple   AccessControlDirection    accessControlDirection]
            [array    byte                      data          count 'commandTypeContainer.numBytes-3'          ]
        ]
        ['ACCESS_POINT_LEFT_OPEN'   *AccessPointLeftOpen
        ]
        ['ACCESS_POINT_FORCED_OPEN' *AccessPointForcedOpen
        ]
        ['ACCESS_POINT_CLOSED'      *AccessPointClosed
        ]
        ['REQUEST_TO_EXIT'          *RequestToExit
        ]
        ['CLOSE_ACCESS_POINT'       *CloseAccessPoint
        ]
        ['LOCK_ACCESS_POINT'        *LockAccessPoint
        ]
    ]
]

[enum uint 8 AccessControlCommandTypeContainer(AccessControlCategory category,AccessControlCommandType commandType, uint 5 numBytes)
    ['0x02' AccessControlCommandCloseAccessPoint                ['SYSTEM_REQUEST',  'CLOSE_ACCESS_POINT',       '2']]
    ['0x0A' AccessControlCommandLockAccessPoint                 ['SYSTEM_REQUEST',  'LOCK_ACCESS_POINT',        '2']]
    ['0x12' AccessControlCommandAccessPointLeftOpen             ['SYSTEM_ACTIVITY', 'ACCESS_POINT_LEFT_OPEN',   '2']]
    ['0x1A' AccessControlCommandAccessPointForcedOpen           ['SYSTEM_ACTIVITY', 'ACCESS_POINT_FORCED_OPEN', '2']]
    ['0x22' AccessControlCommandAccessPointClosed               ['SYSTEM_ACTIVITY', 'ACCESS_POINT_CLOSED',      '2']]
    ['0x32' AccessControlCommandRequestToExit                   ['SYSTEM_ACTIVITY', 'REQUEST_TO_EXIT',          '2']]
    ['0xA0' AccessControlCommandValidAccessRequest_0Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '0']]
    ['0xA1' AccessControlCommandValidAccessRequest_1Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '1']]
    ['0xA2' AccessControlCommandValidAccessRequest_2Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '2']]
    ['0xA3' AccessControlCommandValidAccessRequest_3Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '3']]
    ['0xA4' AccessControlCommandValidAccessRequest_4Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '4']]
    ['0xA5' AccessControlCommandValidAccessRequest_5Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '5']]
    ['0xA6' AccessControlCommandValidAccessRequest_6Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '6']]
    ['0xA7' AccessControlCommandValidAccessRequest_7Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '7']]
    ['0xA8' AccessControlCommandValidAccessRequest_8Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '8']]
    ['0xA9' AccessControlCommandValidAccessRequest_9Bytes       ['SYSTEM_ACTIVITY', 'VALID_ACCESS',             '9']]
    ['0xAA' AccessControlCommandValidAccessRequest_10Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '10']]
    ['0xAB' AccessControlCommandValidAccessRequest_11Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '11']]
    ['0xAC' AccessControlCommandValidAccessRequest_12Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '12']]
    ['0xAD' AccessControlCommandValidAccessRequest_13Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '13']]
    ['0xAE' AccessControlCommandValidAccessRequest_14Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '14']]
    ['0xAF' AccessControlCommandValidAccessRequest_15Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '15']]
    ['0xB0' AccessControlCommandValidAccessRequest_16Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '16']]
    ['0xB1' AccessControlCommandValidAccessRequest_17Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '17']]
    ['0xB2' AccessControlCommandValidAccessRequest_18Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '18']]
    ['0xB3' AccessControlCommandValidAccessRequest_19Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '19']]
    ['0xB4' AccessControlCommandValidAccessRequest_20Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '20']]
    ['0xB5' AccessControlCommandValidAccessRequest_21Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '21']]
    ['0xB6' AccessControlCommandValidAccessRequest_22Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '22']]
    ['0xB7' AccessControlCommandValidAccessRequest_23Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '23']]
    ['0xB8' AccessControlCommandValidAccessRequest_24Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '24']]
    ['0xB9' AccessControlCommandValidAccessRequest_25Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '25']]
    ['0xBA' AccessControlCommandValidAccessRequest_26Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '26']]
    ['0xBB' AccessControlCommandValidAccessRequest_27Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '27']]
    ['0xBC' AccessControlCommandValidAccessRequest_28Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '28']]
    ['0xBD' AccessControlCommandValidAccessRequest_29Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '29']]
    ['0xBE' AccessControlCommandValidAccessRequest_30Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '30']]
    ['0xBF' AccessControlCommandValidAccessRequest_31Bytes      ['SYSTEM_ACTIVITY', 'VALID_ACCESS',            '31']]
    ['0xC0' AccessControlCommandInvalidAccessRequest_0Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '0']]
    ['0xC1' AccessControlCommandInvalidAccessRequest_1Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '1']]
    ['0xC2' AccessControlCommandInvalidAccessRequest_2Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '2']]
    ['0xC3' AccessControlCommandInvalidAccessRequest_3Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '3']]
    ['0xC4' AccessControlCommandInvalidAccessRequest_4Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '4']]
    ['0xC5' AccessControlCommandInvalidAccessRequest_5Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '5']]
    ['0xC6' AccessControlCommandInvalidAccessRequest_6Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '6']]
    ['0xC7' AccessControlCommandInvalidAccessRequest_7Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '7']]
    ['0xC8' AccessControlCommandInvalidAccessRequest_8Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '8']]
    ['0xC9' AccessControlCommandInvalidAccessRequest_9Bytes     ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',           '9']]
    ['0xCA' AccessControlCommandInvalidAccessRequest_10Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '10']]
    ['0xCB' AccessControlCommandInvalidAccessRequest_11Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '11']]
    ['0xCC' AccessControlCommandInvalidAccessRequest_12Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '12']]
    ['0xCD' AccessControlCommandInvalidAccessRequest_13Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '13']]
    ['0xCE' AccessControlCommandInvalidAccessRequest_14Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '14']]
    ['0xCF' AccessControlCommandInvalidAccessRequest_15Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '15']]
    ['0xD0' AccessControlCommandInvalidAccessRequest_16Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '16']]
    ['0xD1' AccessControlCommandInvalidAccessRequest_17Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '17']]
    ['0xD2' AccessControlCommandInvalidAccessRequest_18Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '18']]
    ['0xD3' AccessControlCommandInvalidAccessRequest_19Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '19']]
    ['0xD4' AccessControlCommandInvalidAccessRequest_20Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '20']]
    ['0xD5' AccessControlCommandInvalidAccessRequest_21Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '21']]
    ['0xD6' AccessControlCommandInvalidAccessRequest_22Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '22']]
    ['0xD7' AccessControlCommandInvalidAccessRequest_23Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '23']]
    ['0xD8' AccessControlCommandInvalidAccessRequest_24Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '24']]
    ['0xD9' AccessControlCommandInvalidAccessRequest_25Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '25']]
    ['0xDA' AccessControlCommandInvalidAccessRequest_26Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '26']]
    ['0xDB' AccessControlCommandInvalidAccessRequest_27Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '27']]
    ['0xDC' AccessControlCommandInvalidAccessRequest_28Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '28']]
    ['0xDD' AccessControlCommandInvalidAccessRequest_29Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '29']]
    ['0xDE' AccessControlCommandInvalidAccessRequest_30Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '30']]
    ['0xDF' AccessControlCommandInvalidAccessRequest_31Bytes    ['SYSTEM_ACTIVITY', 'INVALID_ACCESS',          '31']]
]

[enum uint 4 AccessControlCommandType(uint 8 numberOfArguments)
    ['0x00' CLOSE_ACCESS_POINT          ['0']]
    ['0x01' LOCK_ACCESS_POINT           ['0']]
    ['0x02' ACCESS_POINT_LEFT_OPEN      ['0']]
    ['0x03' ACCESS_POINT_FORCED_OPEN    ['0']]
    ['0x04' ACCESS_POINT_CLOSED         ['0']]
    ['0x05' REQUEST_TO_EXIT             ['0']]
    ['0x06' VALID_ACCESS                ['2']]
    ['0x07' INVALID_ACCESS              ['2']]
]

[enum uint 4 AccessControlCategory
    ['0x00' SYSTEM_ACTIVITY   ]
    ['0x01' SYSTEM_REQUEST    ]
]

[enum uint 8 AccessControlDirection
    ['0x00' NOT_USED    ]
    ['0x01' IN          ]
    ['0x02' OUT         ]
]
