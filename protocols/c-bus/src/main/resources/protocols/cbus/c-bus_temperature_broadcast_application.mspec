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

[type TemperatureBroadcastData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsTemperatureBroadcastCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  TemperatureBroadcastCommandTypeContainer   commandTypeContainer                                   ]
    [virtual TemperatureBroadcastCommandType            commandType          'commandTypeContainer.commandType']
    [simple  byte                                       temperatureGroup                                       ]
    [simple  byte                                       temperatureByte                                        ]
    [virtual float 32                                   temperatureInCelsius 'temperatureByte/4'               ]
]

[enum uint 8 TemperatureBroadcastCommandTypeContainer(TemperatureBroadcastCommandType commandType, uint 5 numBytes)
    ['0x02' TemperatureBroadcastCommandSetBroadcastEvent0_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x0A' TemperatureBroadcastCommandSetBroadcastEvent1_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x12' TemperatureBroadcastCommandSetBroadcastEvent2_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x1A' TemperatureBroadcastCommandSetBroadcastEvent3_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x22' TemperatureBroadcastCommandSetBroadcastEvent4_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x2A' TemperatureBroadcastCommandSetBroadcastEvent5_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x32' TemperatureBroadcastCommandSetBroadcastEvent6_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x3A' TemperatureBroadcastCommandSetBroadcastEvent7_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x42' TemperatureBroadcastCommandSetBroadcastEvent8_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x4A' TemperatureBroadcastCommandSetBroadcastEvent9_2Bytes       ['BROADCAST_EVENT',   '2']]
    ['0x52' TemperatureBroadcastCommandSetBroadcastEvent10_2Bytes      ['BROADCAST_EVENT',   '2']]
    ['0x5A' TemperatureBroadcastCommandSetBroadcastEvent11_2Bytes      ['BROADCAST_EVENT',   '2']]
    ['0x62' TemperatureBroadcastCommandSetBroadcastEvent12_2Bytes      ['BROADCAST_EVENT',   '2']]
    ['0x6A' TemperatureBroadcastCommandSetBroadcastEvent13_2Bytes      ['BROADCAST_EVENT',   '2']]
    ['0x72' TemperatureBroadcastCommandSetBroadcastEvent14_2Bytes      ['BROADCAST_EVENT',   '2']]
    ['0x7A' TemperatureBroadcastCommandSetBroadcastEvent15_2Bytes      ['BROADCAST_EVENT',   '2']]
]

[enum uint 4 TemperatureBroadcastCommandType
    ['0x00' BROADCAST_EVENT   ]
]
