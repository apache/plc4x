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

[type TelephonyData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsTelephonyCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  TelephonyCommandTypeContainer      commandTypeContainer                                   ]
    [virtual TelephonyCommandType               commandType          'commandTypeContainer.commandType']
    [simple  byte   argument]
    [typeSwitch commandType, argument
        ['EVENT', '0x01'  *LineOnHook
        ]
        ['EVENT', '0x02'  *LineOffHook(TelephonyCommandTypeContainer commandTypeContainer)
            [simple   LineOffHookReason     reason]
            [simple   vstring '(commandTypeContainer.numBytes-2)*8' number ]
        ]
        ['EVENT', '0x03'  *DialOutFailure
            [simple   DialOutFailureReason  reason]
        ]
        ['EVENT', '0x04'  *DialInFailure
            [simple   DialInFailureReason   reason]
        ]
        ['EVENT', '0x05'  *Ringing(TelephonyCommandTypeContainer commandTypeContainer)
            [reserved byte    '0x01'              ]
            [simple   vstring '(commandTypeContainer.numBytes-2)*8' number ]
        ]
        ['EVENT', '0x06'  *RecallLastNumber(TelephonyCommandTypeContainer commandTypeContainer)
            [simple   byte    recallLastNumberType  ]
            [virtual  bit     isNumberOfLastOutgoingCall    'recallLastNumberType == 0x01'  ]
            [virtual  bit     isNumberOfLastIncomingCall    'recallLastNumberType == 0x02'  ]
            [simple   vstring '(commandTypeContainer.numBytes-2)*8' number ]
        ]
        ['EVENT', '0x07'  *InternetConnectionRequestMade
        ]
        ['EVENT', '0x80'  *IsolateSecondaryOutlet
            [simple   byte    isolateStatus]
            [virtual  bit     isBehaveNormal 'isolateStatus == 0x00']
            [virtual  bit     isToBeIsolated 'isolateStatus == 0x01']
        ]
        ['EVENT', '0x81'  *RecallLastNumberRequest
            [simple   byte    recallLastNumberType  ]
            [virtual  bit     isNumberOfLastOutgoingCall    'recallLastNumberType == 0x01'  ]
            [virtual  bit     isNumberOfLastIncomingCall    'recallLastNumberType == 0x02'  ]
        ]
        ['EVENT', '0x82'  *RejectIncomingCall
        ]
        ['EVENT', '0x83'  *Divert(TelephonyCommandTypeContainer commandTypeContainer)
            [simple   vstring '(commandTypeContainer.numBytes-1)*8' number ]
        ]
        ['EVENT', '0x84'  *ClearDiversion
        ]
    ]
]

[enum uint 8 TelephonyCommandTypeContainer(TelephonyCommandType commandType, uint 5 numBytes)
    ['0x09' TelephonyCommandLineOnHook            ['EVENT',  '1']]
    ['0xA0' TelephonyCommandLineOffHook_0Bytes    ['EVENT',  '0']]
    ['0xA1' TelephonyCommandLineOffHook_1Bytes    ['EVENT',  '1']]
    ['0xA2' TelephonyCommandLineOffHook_2Bytes    ['EVENT',  '2']]
    ['0xA3' TelephonyCommandLineOffHook_3Bytes    ['EVENT',  '3']]
    ['0xA4' TelephonyCommandLineOffHook_4Bytes    ['EVENT',  '4']]
    ['0xA5' TelephonyCommandLineOffHook_5Bytes    ['EVENT',  '5']]
    ['0xA6' TelephonyCommandLineOffHook_6Bytes    ['EVENT',  '6']]
    ['0xA7' TelephonyCommandLineOffHook_7Bytes    ['EVENT',  '7']]
    ['0xA8' TelephonyCommandLineOffHook_8Bytes    ['EVENT',  '8']]
    ['0xA9' TelephonyCommandLineOffHook_9Bytes    ['EVENT',  '9']]
    ['0xAA' TelephonyCommandLineOffHook_10Bytes   ['EVENT', '10']]
    ['0xAB' TelephonyCommandLineOffHook_11Bytes   ['EVENT', '11']]
    ['0xAC' TelephonyCommandLineOffHook_12Bytes   ['EVENT', '12']]
    ['0xAD' TelephonyCommandLineOffHook_13Bytes   ['EVENT', '13']]
    ['0xAE' TelephonyCommandLineOffHook_14Bytes   ['EVENT', '14']]
    ['0xAF' TelephonyCommandLineOffHook_15Bytes   ['EVENT', '15']]
    ['0xB0' TelephonyCommandLineOffHook_16Bytes   ['EVENT', '16']]
    ['0xB1' TelephonyCommandLineOffHook_17Bytes   ['EVENT', '17']]
    ['0xB2' TelephonyCommandLineOffHook_18Bytes   ['EVENT', '18']]
    ['0xB3' TelephonyCommandLineOffHook_19Bytes   ['EVENT', '19']]
    ['0xB4' TelephonyCommandLineOffHook_20Bytes   ['EVENT', '20']]
    ['0xB5' TelephonyCommandLineOffHook_21Bytes   ['EVENT', '21']]
    ['0xB6' TelephonyCommandLineOffHook_22Bytes   ['EVENT', '22']]
    ['0xB7' TelephonyCommandLineOffHook_23Bytes   ['EVENT', '23']]
    ['0xB8' TelephonyCommandLineOffHook_24Bytes   ['EVENT', '24']]
    ['0xB9' TelephonyCommandLineOffHook_25Bytes   ['EVENT', '25']]
    ['0xBA' TelephonyCommandLineOffHook_26Bytes   ['EVENT', '26']]
    ['0xBB' TelephonyCommandLineOffHook_27Bytes   ['EVENT', '27']]
    ['0xBC' TelephonyCommandLineOffHook_28Bytes   ['EVENT', '28']]
    ['0xBD' TelephonyCommandLineOffHook_29Bytes   ['EVENT', '29']]
    ['0xBE' TelephonyCommandLineOffHook_30Bytes   ['EVENT', '30']]
    ['0xBF' TelephonyCommandLineOffHook_31Bytes   ['EVENT', '31']]
]

[enum uint 4 TelephonyCommandType
    ['0x00' EVENT ]
]

[enum uint 8 LineOffHookReason
    ['0x01' INCOMING_VOICE_CALL ]
    ['0x02' INCOMING_DATA_CALL  ]
    ['0x03' INCOMING_CALL       ]
    ['0x10' OUTGOING_VOICE_CALL ]
    ['0x20' OUTGOING_DATA_CALL  ]
    ['0x30' OUTGOING_CALL       ]
    ['0x40' CBTI_IS_SETTING     ]
    ['0x50' CBTI_IS_CLEARING    ]
]

[enum uint 8 DialOutFailureReason
    ['0x01' NO_DIAL_TONE                            ]
    ['0x02' NO_ANSWER                               ]
    ['0x03' NO_VALID_ACKNOWLEDGEMENT_OF_PROMPTS     ]
    ['0x04' NUMBER_WAS_UNOBTAINABLE_DOES_NOT_EXIST  ]
    ['0x05' NUMBER_WAS_BUSY                         ]
    ['0x06' INTERNAL_FAILURE                        ]
]

[enum uint 8 DialInFailureReason
    ['0x01' PHONE_STOPPED_RINGING                   ]
]
