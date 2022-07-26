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

[type MeasurementData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsMeasurementCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  MeasurementCommandTypeContainer      commandTypeContainer                                   ]
    [virtual MeasurementCommandType               commandType          'commandTypeContainer.commandType']
    [typeSwitch commandType
        ['MEASUREMENT_EVENT'              *ChannelMeasurementData
            [simple   uint 8            deviceId   ]
            [simple   uint 8            channel    ]
            [simple   MeasurementUnits  units      ]
            [simple    int 8            multiplier ]
            [simple   uint 8            msb        ]
            [simple   uint 8            lsb        ]
            [virtual  uint 16           rawValue    'msb<<8|lsb'            ]
            [virtual  float 64          value       'rawValue*multiplier*10']
        ]
    ]
]

[enum uint 8 MeasurementCommandTypeContainer(MeasurementCommandType commandType, uint 5 numBytes)
    ['0x0E' MeasurementCommandChannelMeasurementData    ['MEASUREMENT_EVENT',  '6']]
]

[enum uint 4 MeasurementCommandType
    ['0x00' MEASUREMENT_EVENT              ]
]

[enum uint 8 MeasurementUnits
    ['0x00' CELSIUS                 ]
    ['0x01' AMPS                    ]
    ['0x02' ANGLE_DEGREES           ]
    ['0x03' COULOMB                 ]
    ['0x04' BOOLEANLOGIC            ]
    ['0x05' FARADS                  ]
    ['0x06' HENRYS                  ]
    ['0x07' HERTZ                   ]
    ['0x08' JOULES                  ]
    ['0x09' KATAL                   ]
    ['0x0A' KG_PER_M3               ]
    ['0x0B' KILOGRAMS               ]
    ['0x0C' LITRES                  ]
    ['0x0D' LITRES_PER_HOUR         ]
    ['0x0E' LITRES_PER_MINUTE       ]
    ['0x0F' LITRES_PER_SECOND       ]
    ['0x10' LUX                     ]
    ['0x11' METRES                  ]
    ['0x12' METRES_PER_MINUTE       ]
    ['0x13' METRES_PER_SECOND       ]
    ['0x14' METRES_PER_S_SQUARED    ]
    ['0x15' MOLE                    ]
    ['0x16' NEWTON_METRE            ]
    ['0x17' NEWTONS                 ]
    ['0x18' OHMS                    ]
    ['0x19' PASCAL                  ]
    ['0x1A' PERCENT                 ]
    ['0x1B' DECIBELS                ]
    ['0x1C' PPM                     ]
    ['0x1D' RPM                     ]
    ['0x1E' SECOND                  ]
    ['0x1F' MINUTES                 ]
    ['0x20' HOURS                   ]
    ['0x21' SIEVERTS                ]
    ['0x22' STERADIAN               ]
    ['0x23' TESLA                   ]
    ['0x24' VOLTS                   ]
    ['0x25' WATT_HOURS              ]
    ['0x26' WATTS                   ]
    ['0x27' WEBERS                  ]
    ['0xFE' NO_UNITS                ]
    ['0xFF' CUSTOM                  ]
]
