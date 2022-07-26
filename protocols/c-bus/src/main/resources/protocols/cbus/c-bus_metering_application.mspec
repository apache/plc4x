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

[type MeteringData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsMeteringCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  MeteringCommandTypeContainer commandTypeContainer                                   ]
    [virtual MeteringCommandType          commandType          'commandTypeContainer.commandType']
    [simple byte argument                                                               ]
    [typeSwitch commandType, argument
        ['EVENT', '0x01'       *MeasureElectricity
        ]
        ['EVENT', '0x02'       *MeasureGas
        ]
        ['EVENT', '0x03'       *MeasureDrinkingWater
        ]
        ['EVENT', '0x04'       *MeasureOtherWater
        ]
        ['EVENT', '0x05'       *MeasureOil
        ]
        ['EVENT', '0x81'       *ElectricityConsumption
            [simple uint 32 kWhr      ] // kilo watt hours
        ]
        ['EVENT', '0x82'       *GasConsumption
            [simple uint 32 mJ        ] // mega joule
        ]
        ['EVENT', '0x83'       *DrinkingWaterConsumption
            [simple uint 32 kL        ] // kilo litre
        ]
        ['EVENT', '0x84'       *OtherWaterConsumption
            [simple uint 32 kL        ] // kilo litre
        ]
        ['EVENT', '0x85'       *OilConsumption
            [simple uint 32 L         ] // litre
        ]
    ]
]

[enum uint 8 MeteringCommandTypeContainer(MeteringCommandType commandType, uint 5 numBytes)
    ['0x08' MeteringCommandEvent_0Bytes                    ['EVENT',  '0']]
    ['0x09' MeteringCommandEvent_1Bytes                    ['EVENT',  '1']]
    ['0x0A' MeteringCommandEvent_2Bytes                    ['EVENT',  '2']]
    ['0x0B' MeteringCommandEvent_3Bytes                    ['EVENT',  '3']]
    ['0x0C' MeteringCommandEvent_4Bytes                    ['EVENT',  '4']]
    ['0x0D' MeteringCommandEvent_5Bytes                    ['EVENT',  '5']]
    ['0x0E' MeteringCommandEvent_6Bytes                    ['EVENT',  '6']]
    ['0x0F' MeteringCommandEvent_7Bytes                    ['EVENT',  '7']]
]

[enum uint 4 MeteringCommandType
    ['0x00' EVENT     ]
]
