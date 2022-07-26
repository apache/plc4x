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

[type TriggerControlData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsTriggerControlCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  TriggerControlCommandTypeContainer commandTypeContainer                                   ]
    [virtual TriggerControlCommandType          commandType          'commandTypeContainer.commandType']
    [simple  byte triggerGroup                                                                         ]
    [virtual bit  isUnused 'triggerGroup > 0xFE'                                                       ]
    [typeSwitch commandType
        ['TRIGGER_EVENT'       *TriggerEvent
            [simple byte actionSelector]
        ]
        ['TRIGGER_MIN'          *TriggerMin
        ]
        ['TRIGGER_MAX'          *TriggerMin
        ]
        ['INDICATOR_KILL'       *IndicatorKill
        ]
        ['LABEL'                *Label(TriggerControlCommandTypeContainer commandTypeContainer)
           [simple   TriggerControlLabelOptions triggerControlOptions                                   ]
           [simple   byte                       actionSelector                                          ]
           [optional Language                   language      'triggerControlOptions.labelType != TriggerControlLabelType.LOAD_DYNAMIC_ICON']
           [array    byte                       data          count '(commandTypeContainer.numBytes-((triggerControlOptions.labelType != TriggerControlLabelType.LOAD_DYNAMIC_ICON)?(4):(3)))'           ]
        ]
    ]
]

[enum uint 8 TriggerControlCommandTypeContainer(TriggerControlCommandType commandType, uint 5 numBytes)
    ['0x01' TriggerControlCommandTriggerMin_1Bytes          ['TRIGGER_MIN',     '1']]
    ['0x09' TriggerControlCommandIndicatorKill_1Bytes       ['INDICATOR_KILL',  '1']]
    ['0x79' TriggerControlCommandTriggerMax_1Bytes          ['TRIGGER_MAX',     '1']]
    ['0x02' TriggerControlCommandTriggerEvent0_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x0A' TriggerControlCommandTriggerEvent1_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x12' TriggerControlCommandTriggerEvent2_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x1A' TriggerControlCommandTriggerEvent3_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x22' TriggerControlCommandTriggerEvent4_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x2A' TriggerControlCommandTriggerEvent5_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x32' TriggerControlCommandTriggerEvent6_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x3A' TriggerControlCommandTriggerEvent7_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x42' TriggerControlCommandTriggerEvent8_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x4A' TriggerControlCommandTriggerEvent9_2Bytes       ['TRIGGER_EVENT',   '2']]
    ['0x52' TriggerControlCommandTriggerEvent10_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0x5A' TriggerControlCommandTriggerEvent11_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0x62' TriggerControlCommandTriggerEvent12_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0x6A' TriggerControlCommandTriggerEvent13_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0x72' TriggerControlCommandTriggerEvent14_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0x7A' TriggerControlCommandTriggerEvent15_2Bytes      ['TRIGGER_EVENT',   '2']]
    ['0xA0' TriggerControlCommandLabel_0Bytes               ['LABEL',           '0']]
    ['0xA1' TriggerControlCommandLabel_1Bytes               ['LABEL',           '1']]
    ['0xA2' TriggerControlCommandLabel_2Bytes               ['LABEL',           '2']]
    ['0xA3' TriggerControlCommandLabel_3Bytes               ['LABEL',           '3']]
    ['0xA4' TriggerControlCommandLabel_4Bytes               ['LABEL',           '4']]
    ['0xA5' TriggerControlCommandLabel_5Bytes               ['LABEL',           '5']]
    ['0xA6' TriggerControlCommandLabel_6Bytes               ['LABEL',           '6']]
    ['0xA7' TriggerControlCommandLabel_7Bytes               ['LABEL',           '7']]
    ['0xA8' TriggerControlCommandLabel_8Bytes               ['LABEL',           '8']]
    ['0xA9' TriggerControlCommandLabel_9Bytes               ['LABEL',           '9']]
    ['0xAA' TriggerControlCommandLabel_10Bytes              ['LABEL',          '10']]
    ['0xAB' TriggerControlCommandLabel_11Bytes              ['LABEL',          '11']]
    ['0xAC' TriggerControlCommandLabel_12Bytes              ['LABEL',          '12']]
    ['0xAD' TriggerControlCommandLabel_13Bytes              ['LABEL',          '13']]
    ['0xAE' TriggerControlCommandLabel_14Bytes              ['LABEL',          '14']]
    ['0xAF' TriggerControlCommandLabel_15Bytes              ['LABEL',          '15']]
    ['0xB0' TriggerControlCommandLabel_16Bytes              ['LABEL',          '16']]
    ['0xB1' TriggerControlCommandLabel_17Bytes              ['LABEL',          '17']]
    ['0xB2' TriggerControlCommandLabel_18Bytes              ['LABEL',          '18']]
    ['0xB3' TriggerControlCommandLabel_19Bytes              ['LABEL',          '19']]
    ['0xB4' TriggerControlCommandLabel_20Bytes              ['LABEL',          '20']]
    ['0xB5' TriggerControlCommandLabel_21Bytes              ['LABEL',          '21']]
    ['0xB6' TriggerControlCommandLabel_22Bytes              ['LABEL',          '22']]
    ['0xB7' TriggerControlCommandLabel_23Bytes              ['LABEL',          '23']]
    ['0xB8' TriggerControlCommandLabel_24Bytes              ['LABEL',          '24']]
    ['0xB9' TriggerControlCommandLabel_25Bytes              ['LABEL',          '25']]
    ['0xBA' TriggerControlCommandLabel_26Bytes              ['LABEL',          '26']]
    ['0xBB' TriggerControlCommandLabel_27Bytes              ['LABEL',          '27']]
    ['0xBC' TriggerControlCommandLabel_28Bytes              ['LABEL',          '28']]
    ['0xBD' TriggerControlCommandLabel_29Bytes              ['LABEL',          '29']]
    ['0xBE' TriggerControlCommandLabel_30Bytes              ['LABEL',          '30']]
    ['0xBF' TriggerControlCommandLabel_31Bytes              ['LABEL',          '31']]
]

[enum uint 4 TriggerControlCommandType
    ['0x00' TRIGGER_EVENT   ]
    ['0x01' TRIGGER_MIN     ]
    ['0x02' TRIGGER_MAX     ]
    ['0x03' INDICATOR_KILL  ]
    ['0x04' LABEL           ]
]

// TODO: maybe can be merged with lightning labels
[type TriggerControlLabelOptions
    [simple   bit                           reservedBit7] // only for dynamic icon loading can switch to 1 (note this could use mspec reserved field but sadly this discards data)
    [simple   TriggerControlLabelFlavour    labelFlavour]
    [reserved bit                           'false'     ]
    [simple   bit                           reservedBit3] // For Control Trigger, this bit must be 0 (note this could use mspec reserved field but sadly this discards data)
    [simple   TriggerControlLabelType       labelType   ]
    [simple   bit                           reservedBit0] // For Control Trigger, this bit must be 1 (note this could use mspec reserved field but sadly this discards data)
]

// TODO: maybe can be merged with lightning labels
[enum uint 2 TriggerControlLabelFlavour
    ['0' FLAVOUR_0              ]
    ['1' FLAVOUR_1              ]
    ['2' FLAVOUR_2              ]
    ['3' FLAVOUR_3              ]
]

[enum uint 2 TriggerControlLabelType
    ['0' TEXT_LABEL             ]
    ['1' PREDEFINED_ICON        ]
    ['2' LOAD_DYNAMIC_ICON      ]
    ['3' SET_PREFERRED_LANGUAGE ]
]
