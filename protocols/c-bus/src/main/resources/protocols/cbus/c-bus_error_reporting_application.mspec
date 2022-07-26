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

[type ErrorReportingData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsErrorReportingCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  ErrorReportingCommandTypeContainer      commandTypeContainer                                   ]
    [virtual ErrorReportingCommandType               commandType          'commandTypeContainer.commandType']
    [typeSwitch commandType
        [*              *Generic
            [simple   ErrorReportingSystemCategory  systemCategory    ]
            [simple   bit                           mostRecent        ]
            [simple   bit                           acknowledge       ]
            [simple   bit                           mostSevere        ]
            [validation 'mostRecent || mostSevere' "Invalid Error condition"]
            [virtual  bit                           isMostSevereError 'mostSevere']
            [virtual  bit                           isMostRecentError 'mostRecent']
            [virtual  bit                           isMostRecentAndMostSevere 'isMostRecentError && isMostSevereError']
            [simple   ErrorReportingSeverity        severity          ]
            [simple   uint 8                        deviceId          ]
            // TODO: maybe split them up according to appendix A
            [simple   uint 8                        errorData1        ]
            [simple   uint 8                        errorData2        ]
        ]
    ]
]

[enum uint 8 ErrorReportingCommandTypeContainer(ErrorReportingCommandType commandType, uint 5 numBytes)
    ['0x05' ErrorReportingCommandDeprecated         ['DEPRECATED',        '5']]
    ['0x15' ErrorReportingCommandErrorReport        ['ERROR_REPORT',      '5']]
    ['0x25' ErrorReportingCommandAcknowledge        ['ACKNOWLEDGE',       '5']]
    ['0x35' ErrorReportingCommandClearMostSevere    ['CLEAR_MOST_SEVERE', '5']]
]

[enum uint 4 ErrorReportingCommandType
    ['0x00' DEPRECATED          ]
    ['0x01' ERROR_REPORT        ]
    ['0x02' ACKNOWLEDGE         ]
    ['0x03' CLEAR_MOST_SEVERE   ]
]

[enum uint 3 ErrorReportingSeverity
    ['0x0' ALL_OK           ]
    ['0x1' OK               ]
    ['0x2' MINOR_FAILURE    ]
    ['0x3' GENERAL_FAILURE  ]
    ['0x4' EXTREME_FAILURE  ]
    ['0x5' RESERVED_1       ]
    ['0x6' RESERVED_2       ]
    ['0x7' RESERVED_3       ]
]

[type ErrorReportingSystemCategory
    [simple ErrorReportingSystemCategoryClass                       systemCategoryClass     ]
    [simple ErrorReportingSystemCategoryType('systemCategoryClass') systemCategoryType      ]
    [simple ErrorReportingSystemCategoryVariant                     systemCategoryVariant   ]
]

[enum uint 4 ErrorReportingSystemCategoryClass
    ['0x0'  RESERVED_0                  ]
    ['0x1'  RESERVED_1                  ]
    ['0x2'  RESERVED_2                  ]
    ['0x3'  RESERVED_3                  ]
    ['0x4'  RESERVED_4                  ]
    ['0x5'  INPUT_UNITS                 ]
    ['0x6'  RESERVED_6                  ]
    ['0x7'  RESERVED_7                  ]
    ['0x8'  RESERVED_8                  ]
    ['0x9'  SUPPORT_UNITS               ]
    ['0xA'  RESERVED_10                 ]
    ['0xB'  BUILDING_MANAGEMENT_SYSTEMS ]
    ['0xC'  RESERVED_12                 ]
    ['0xD'  OUTPUT_UNITS                ]
    ['0xE'  RESERVED_14                 ]
    ['0xF'  CLIMATE_CONTROLLERS         ]
]

[type ErrorReportingSystemCategoryType(ErrorReportingSystemCategoryClass errorReportingSystemCategoryClass)
    [typeSwitch errorReportingSystemCategoryClass
        ['INPUT_UNITS'                  *InputUnits
            [simple ErrorReportingSystemCategoryTypeForInputUnits                   categoryForType ]
        ]
        ['SUPPORT_UNITS'                *SupportUnits
            [simple ErrorReportingSystemCategoryTypeForSupportUnits                 categoryForType ]
        ]
        ['BUILDING_MANAGEMENT_SYSTEMS'  *BuildingManagementSystems
            [simple ErrorReportingSystemCategoryTypeForBuildingManagementSystems    categoryForType ]
        ]
        ['OUTPUT_UNITS'                 *OutputUnits
            [simple ErrorReportingSystemCategoryTypeForOutputUnits                  categoryForType ]
        ]
        ['CLIMATE_CONTROLLERS'          *ClimateControllers
            [simple ErrorReportingSystemCategoryTypeForClimateControllers           categoryForType ]
        ]
        [*                              *Reserved
            [simple uint 4  reservedValue]
        ]
    ]
]

[enum uint 4 ErrorReportingSystemCategoryTypeForInputUnits
    ['0x0'  KEY_UNITS                   ]
    ['0x1'  TELECOMMAND_AND_REMOTE_ENTRY]
    ['0x2'  RESERVED_2                  ]
    ['0x3'  RESERVED_3                  ]
    ['0x4'  RESERVED_4                  ]
    ['0x5'  RESERVED_5                  ]
    ['0x6'  RESERVED_6                  ]
    ['0x7'  RESERVED_7                  ]
    ['0x8'  RESERVED_8                  ]
    ['0x9'  RESERVED_9                  ]
    ['0xA'  RESERVED_10                 ]
    ['0xB'  RESERVED_11                 ]
    ['0xC'  RESERVED_12                 ]
    ['0xD'  RESERVED_13                 ]
    ['0xE'  RESERVED_14                 ]
    ['0xF'  RESERVED_15                 ]
]

[enum uint 4 ErrorReportingSystemCategoryTypeForSupportUnits
    ['0x0'  POWER_SUPPLIES              ]
    ['0x1'  RESERVED_1                  ]
    ['0x2'  RESERVED_2                  ]
    ['0x3'  RESERVED_3                  ]
    ['0x4'  RESERVED_4                  ]
    ['0x5'  RESERVED_5                  ]
    ['0x6'  RESERVED_6                  ]
    ['0x7'  RESERVED_7                  ]
    ['0x8'  RESERVED_8                  ]
    ['0x9'  RESERVED_9                  ]
    ['0xA'  RESERVED_10                 ]
    ['0xB'  RESERVED_11                 ]
    ['0xC'  RESERVED_12                 ]
    ['0xD'  RESERVED_13                 ]
    ['0xE'  RESERVED_14                 ]
    ['0xF'  RESERVED_15                 ]
]

[enum uint 4 ErrorReportingSystemCategoryTypeForBuildingManagementSystems
    ['0x0'  BMS_DIAGNOSTIC_REPORTING    ]
    ['0x1'  RESERVED_1                  ]
    ['0x2'  RESERVED_2                  ]
    ['0x3'  RESERVED_3                  ]
    ['0x4'  RESERVED_4                  ]
    ['0x5'  RESERVED_5                  ]
    ['0x6'  RESERVED_6                  ]
    ['0x7'  RESERVED_7                  ]
    ['0x8'  RESERVED_8                  ]
    ['0x9'  RESERVED_9                  ]
    ['0xA'  RESERVED_10                 ]
    ['0xB'  RESERVED_11                 ]
    ['0xC'  RESERVED_12                 ]
    ['0xD'  RESERVED_13                 ]
    ['0xE'  RESERVED_14                 ]
    ['0xF'  RESERVED_15                 ]
]

[enum uint 4 ErrorReportingSystemCategoryTypeForOutputUnits
    ['0x0'  LE_MONOBLOCK_DIMMERS                        ]
    ['0x1'  TE_MONOBLOCK_DIMMERS                        ]
    ['0x2'  RESERVED_2                                  ]
    ['0x3'  RESERVED_3                                  ]
    ['0x4'  RELAYS_AND_OTHER_ON_OFF_SWITCHING_DEVICES   ]
    ['0x5'  RESERVED_5                                  ]
    ['0x6'  PWM_DIMMERS_INCLUDES_LED_CONTROL            ]
    ['0x7'  SINEWAVE_MONOBLOCK_DIMMERS                  ]
    ['0x8'  RESERVED_8                                  ]
    ['0x9'  RESERVED_9                                  ]
    ['0xA'  DALI_DSI_AND_OTHER_BALLAST_CONTROL_GATEWAYS ]
    ['0xB'  MODULAR_DIMMERS                             ]
    ['0xC'  RESERVED_12                                 ]
    ['0xD'  UNIVERSAL_MONOBLOCK_DIMMERS                 ]
    ['0xE'  DEVICE_CONTROLLERS_IR_RS_232_etc            ]
    ['0xF'  RESERVED_15                                 ]
]

[enum uint 4 ErrorReportingSystemCategoryTypeForClimateControllers
    ['0x0'  AIR_CONDITIONING_SYSTEM     ]
    ['0x1'  RESERVED_1                  ]
    ['0x2'  RESERVED_2                  ]
    ['0x3'  RESERVED_3                  ]
    ['0x4'  RESERVED_4                  ]
    ['0x5'  RESERVED_5                  ]
    ['0x6'  RESERVED_6                  ]
    ['0x7'  RESERVED_7                  ]
    ['0x8'  RESERVED_8                  ]
    ['0x9'  RESERVED_9                  ]
    ['0xA'  RESERVED_10                 ]
    ['0xB'  RESERVED_11                 ]
    ['0xC'  GLOBAL_WARMING_MODULATOR    ]
    ['0xD'  RESERVED_13                 ]
    ['0xE'  RESERVED_14                 ]
    ['0xF'  RESERVED_15                 ]
]

[enum uint 2 ErrorReportingSystemCategoryVariant
    ['0x0' RESERVED_0   ]
    ['0x1' RESERVED_1   ]
    ['0x2' RESERVED_2   ]
    ['0x3' RESERVED_3   ]
]
