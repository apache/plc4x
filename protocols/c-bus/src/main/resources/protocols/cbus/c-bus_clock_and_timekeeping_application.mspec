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

[type ClockAndTimekeepingData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsClockAndTimekeepingCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  ClockAndTimekeepingCommandTypeContainer    commandTypeContainer                                   ]
    [virtual ClockAndTimekeepingCommandType             commandType          'commandTypeContainer.commandType']
    [simple  byte   argument]
    [typeSwitch commandType, argument
        ['UPDATE_NETWORK_VARIABLE', '0x01'  *UpdateTime
            [simple   uint 8 hours          ]
            [simple   uint 8 minute         ]
            [simple   uint 8 second         ]
            [simple   byte   daylightSaving ]
            [virtual  bit    isNoDaylightSavings 'daylightSaving == 0x00']
            [virtual  bit    isAdvancedBy1Hour   'daylightSaving == 0x01']
            [virtual  bit    isReserved          'daylightSaving > 0x01 && daylightSaving <= 0xFE']
            [virtual  bit    isUnknown           'daylightSaving > 0xFE']
        ]
        ['UPDATE_NETWORK_VARIABLE', '0x02'  *UpdateDate
            [simple   byte   yearHigh       ]
            [simple   byte   yearLow        ]
            [simple   uint 8 month          ]
            [simple   uint 8 day            ]
            [simple   uint 8 dayOfWeek      ]
        ]
        ['REQUEST_REFRESH', '0x03'          *RequestRefresh
        ]
    ]
]

[enum uint 8 ClockAndTimekeepingCommandTypeContainer(ClockAndTimekeepingCommandType commandType, uint 5 numBytes)
    ['0x08' MediaTransportControlCommandUpdateNetworkVariable_0Bytes    ['UPDATE_NETWORK_VARIABLE', '0']]
    ['0x09' MediaTransportControlCommandUpdateNetworkVariable_1Bytes    ['UPDATE_NETWORK_VARIABLE', '1']]
    ['0x0A' MediaTransportControlCommandUpdateNetworkVariable_2Bytes    ['UPDATE_NETWORK_VARIABLE', '2']]
    ['0x0B' MediaTransportControlCommandUpdateNetworkVariable_3Bytes    ['UPDATE_NETWORK_VARIABLE', '3']]
    ['0x0C' MediaTransportControlCommandUpdateNetworkVariable_4Bytes    ['UPDATE_NETWORK_VARIABLE', '4']]
    ['0x0D' MediaTransportControlCommandUpdateNetworkVariable_5Bytes    ['UPDATE_NETWORK_VARIABLE', '5']]
    ['0x0E' MediaTransportControlCommandUpdateNetworkVariable_6Bytes    ['UPDATE_NETWORK_VARIABLE', '6']]
    ['0x0F' MediaTransportControlCommandUpdateNetworkVariable_7Bytes    ['UPDATE_NETWORK_VARIABLE', '7']]
    ['0x11' MediaTransportControlCommandRequestRefresh                  ['REQUEST_REFRESH',         '1']]
]

[enum uint 4 ClockAndTimekeepingCommandType
    ['0x00' UPDATE_NETWORK_VARIABLE ]
    ['0x01' REQUEST_REFRESH         ]
]
