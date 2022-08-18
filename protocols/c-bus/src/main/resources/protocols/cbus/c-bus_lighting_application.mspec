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

[type LightingData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsLightingCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  LightingCommandTypeContainer commandTypeContainer                                   ]
    [virtual LightingCommandType          commandType          'commandTypeContainer.commandType']
    [typeSwitch commandType
        ['OFF'            *Off
            [simple byte group                                                              ]
        ]
        ['ON'             *On
            [simple byte group                                                              ]
        ]
        ['RAMP_TO_LEVEL'  *RampToLevel
            [simple byte group                                                              ]
            [simple byte level                                                              ]
        ]
        ['TERMINATE_RAMP' *TerminateRamp
            [simple byte group                                                              ]
        ]
        ['LABEL'          *Label(LightingCommandTypeContainer commandTypeContainer)
            [simple   byte                  group                                                   ]
            [simple   LightingLabelOptions  labelOptions                                            ]
            [optional Language              language      'labelOptions.labelType != LightingLabelType.LOAD_DYNAMIC_ICON']
            [array    byte                  data        count '(commandTypeContainer.numBytes-((labelOptions.labelType != LightingLabelType.LOAD_DYNAMIC_ICON)?(3):(2)))'           ]
        ]
    ]
]

[type LightingLabelOptions
    [reserved bit                   'false'     ] // only for dynamic icon loading can switch to 1
    [simple   LightingLabelFlavour  labelFlavour]
    [reserved bit                   'false'     ]
    [reserved bit                   'false'     ] // For Lighting, this bit must be 0
    [simple   LightingLabelType     labelType   ]
    [reserved bit                   'false'     ] // For Lighting, this bit must be 0
]

[enum uint 2 LightingLabelFlavour
    ['0' FLAVOUR_0              ]
    ['1' FLAVOUR_1              ]
    ['2' FLAVOUR_2              ]
    ['3' FLAVOUR_3              ]
]

[enum uint 2 LightingLabelType
    ['0' TEXT_LABEL             ]
    ['1' PREDEFINED_ICON        ]
    ['2' LOAD_DYNAMIC_ICON      ]
    ['3' SET_PREFERRED_LANGUAGE ]
]

[enum uint 8 Language
    ['0x00' NO_LANGUAGE                 ]
    ['0x01' ENGLISH                     ]
    ['0x02' ENGLISH_AUSTRALIA           ]
    ['0x03' ENGLISH_BELIZE              ]
    ['0x04' ENGLISH_CANADA              ]
    ['0x05' ENGLISH_CARRIBEAN           ]
    ['0x06' ENGLISH_IRELAND             ]
    ['0x07' ENGLISH_JAMAICA             ]
    ['0x08' ENGLISH_NEW_ZEALAND         ]
    ['0x09' ENGLISH_PHILIPPINES         ]
    ['0x0A' ENGLISH_SOUTH_AFRICA        ]
    ['0x0B' ENGLISH_TRINIDAD            ]
    ['0x0C' ENGLISH_UK                  ]
    ['0x0D' ENGLISH_USA                 ]
    ['0x0E' ENGLISH_ZIMBABWE            ]
    ['0x40' AFRIKAANS                   ]
    ['0x41' BASQUE                      ]
    ['0x42' CATALAN                     ]
    ['0x43' DANISH                      ]
    ['0x44' DUTCH_BELGIUM               ]
    ['0x45' DUTCH_NETHERLANDS           ]
    ['0x46' FAEROESE                    ]
    ['0x47' FINNISH                     ]
    ['0x48' FRENCH_BELGIUM              ]
    ['0x49' FRENCH_CANADA               ]
    ['0x4A' FRENCH                      ]
    ['0x4B' FRENCH_LUXEMBOURG           ]
    ['0x4C' FRENCH_MONACO               ]
    ['0x4D' FRENCH_SWITZERLAND          ]
    ['0x4E' GALICIAN                    ]
    ['0x4F' GERMAN_AUSTRIA              ]
    ['0x50' GERMAN                      ]
    ['0x51' GERMAN_LIECHTENSTEIN        ]
    ['0x52' GERMAN_LUXEMBOURG           ]
    ['0x53' GERMAN_SWITZERLAND          ]
    ['0x54' ICELANDIC                   ]
    ['0x55' INDONESIAN                  ]
    ['0x56' ITALIAN                     ]
    ['0x57' ITALIAN_SWITZERLAND         ]
    ['0x58' MALAY_BRUNEI                ]
    ['0x59' MALAY                       ]
    ['0x5A' NORWEGIAN                   ]
    ['0x5B' NORWEGIAN_NYNORSK           ]
    ['0x5C' PORTUGUESE_BRAZIL           ]
    ['0x5D' PORTUGUESE                  ]
    ['0x5E' SPANISH_ARGENTINE           ]
    ['0x5F' SPANISH_BOLIVIA             ]
    ['0x60' SPANISH_CHILE               ]
    ['0x61' SPANISH_COLOMBIA            ]
    ['0x62' SPANISH_COSTA_RICA          ]
    ['0x63' SPANISH_DOMINICAN_REPUBLIC  ]
    ['0x64' SPANISH_ECUADOR             ]
    ['0x65' SPANISH_EL_SALVADOR         ]
    ['0x66' SPANISH_GUATEMALA           ]
    ['0x67' SPANISH_HONDURAS            ]
    ['0x68' SPANISH                     ]
    ['0x69' SPANISH_MEXICO              ]
    ['0x6A' SPANISH_NICARAGUA           ]
    ['0x6B' SPANISH_PANAMA              ]
    ['0x6C' SPANISH_PARAGUAY            ]
    ['0x6D' SPANISH_PERU                ]
    ['0x6E' SPANISH_PERTO_RICO          ]
    ['0x6F' SPANISH_TRADITIONAL         ]
    ['0x70' SPANISH_URUGUAY             ]
    ['0x71' SPANISH_VENEZUELA           ]
    ['0x72' SWAHILI                     ]
    ['0x73' SWEDISH                     ]
    ['0x74' SWEDISH_FINLAND             ]
    ['0xCA' CHINESE_CP936               ]
]

[enum uint 8 LightingCommandTypeContainer(LightingCommandType commandType, uint 5 numBytes)
    ['0x01' LightingCommandOff                       ['OFF',             '1' ]]
    ['0x79' LightingCommandOn                        ['ON',              '1' ]]
    ['0x02' LightingCommandRampToLevel_Instantaneous ['RAMP_TO_LEVEL',   '1' ]]
    ['0x0A' LightingCommandRampToLevel_4Second       ['RAMP_TO_LEVEL',   '2' ]]
    ['0x12' LightingCommandRampToLevel_8Second       ['RAMP_TO_LEVEL',   '2' ]]
    ['0x1A' LightingCommandRampToLevel_12Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x22' LightingCommandRampToLevel_20Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x2A' LightingCommandRampToLevel_30Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x32' LightingCommandRampToLevel_40Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x3A' LightingCommandRampToLevel_60Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x42' LightingCommandRampToLevel_90Second      ['RAMP_TO_LEVEL',   '2' ]]
    ['0x4A' LightingCommandRampToLevel_120Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x52' LightingCommandRampToLevel_180Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x5A' LightingCommandRampToLevel_300Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x62' LightingCommandRampToLevel_420Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x6A' LightingCommandRampToLevel_600Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x72' LightingCommandRampToLevel_900Second     ['RAMP_TO_LEVEL',   '2' ]]
    ['0x7A' LightingCommandRampToLevel_1020Second    ['RAMP_TO_LEVEL',   '2' ]]
    ['0x09' LightingCommandTerminateRamp             ['TERMINATE_RAMP',  '1' ]]
    ['0xA0' LightingCommandLabel_0Bytes              ['LABEL',           '0' ]]
    ['0xA1' LightingCommandLabel_1Bytes              ['LABEL',           '1' ]]
    ['0xA2' LightingCommandLabel_2Bytes              ['LABEL',           '2' ]]
    ['0xA3' LightingCommandLabel_3Bytes              ['LABEL',           '3' ]]
    ['0xA4' LightingCommandLabel_4Bytes              ['LABEL',           '4' ]]
    ['0xA5' LightingCommandLabel_5Bytes              ['LABEL',           '5' ]]
    ['0xA6' LightingCommandLabel_6Bytes              ['LABEL',           '6' ]]
    ['0xA7' LightingCommandLabel_7Bytes              ['LABEL',           '7' ]]
    ['0xA8' LightingCommandLabel_8Bytes              ['LABEL',           '8' ]]
    ['0xA9' LightingCommandLabel_9Bytes              ['LABEL',           '9' ]]
    ['0xAA' LightingCommandLabel_10Bytes             ['LABEL',          '10' ]]
    ['0xAB' LightingCommandLabel_11Bytes             ['LABEL',          '11' ]]
    ['0xAC' LightingCommandLabel_12Bytes             ['LABEL',          '12' ]]
    ['0xAD' LightingCommandLabel_13Bytes             ['LABEL',          '13' ]]
    ['0xAE' LightingCommandLabel_14Bytes             ['LABEL',          '14' ]]
    ['0xAF' LightingCommandLabel_15Bytes             ['LABEL',          '15' ]]
    ['0xB0' LightingCommandLabel_16Bytes             ['LABEL',          '16' ]]
    ['0xB1' LightingCommandLabel_17Bytes             ['LABEL',          '17' ]]
    ['0xB2' LightingCommandLabel_18Bytes             ['LABEL',          '18' ]]
    ['0xB3' LightingCommandLabel_19Bytes             ['LABEL',          '19' ]]
    ['0xB4' LightingCommandLabel_20Bytes             ['LABEL',          '20' ]]
    ['0xB5' LightingCommandLabel_21Bytes             ['LABEL',          '21' ]]
    ['0xB6' LightingCommandLabel_22Bytes             ['LABEL',          '22' ]]
    ['0xB7' LightingCommandLabel_23Bytes             ['LABEL',          '23' ]]
    ['0xB8' LightingCommandLabel_24Bytes             ['LABEL',          '24' ]]
    ['0xB9' LightingCommandLabel_25Bytes             ['LABEL',          '25' ]]
    ['0xBA' LightingCommandLabel_26Bytes             ['LABEL',          '26' ]]
    ['0xBB' LightingCommandLabel_27Bytes             ['LABEL',          '27' ]]
    ['0xBC' LightingCommandLabel_28Bytes             ['LABEL',          '28' ]]
    ['0xBD' LightingCommandLabel_29Bytes             ['LABEL',          '29' ]]
    ['0xBE' LightingCommandLabel_30Bytes             ['LABEL',          '30' ]]
    ['0xBF' LightingCommandLabel_32Bytes             ['LABEL',          '31' ]]
]

[enum uint 4 LightingCommandType(uint 8 numberOfArguments)
    ['0x00' OFF            ['1']]
    ['0x01' ON             ['1']]
    ['0x02' RAMP_TO_LEVEL  ['2']]
    ['0x03' TERMINATE_RAMP ['1']]
    ['0x04' LABEL          ['4']]
]
