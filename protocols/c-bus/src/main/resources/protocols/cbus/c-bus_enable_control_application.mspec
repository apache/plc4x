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

[type EnableControlData
    //TODO: golang doesn't like checking for null so we use that static call to check that the enum is known
    [validation 'STATIC_CALL("knowsEnableControlCommandTypeContainer", readBuffer)' "no command type could be found" shouldFail=false]
    [simple  EnableControlCommandTypeContainer commandTypeContainer                                   ]
    [virtual EnableControlCommandType          commandType          'commandTypeContainer.commandType']
    [simple  byte                              enableNetworkVariable                                  ]
    [simple  byte                              value                                                  ]
]

[enum uint 8 EnableControlCommandTypeContainer(EnableControlCommandType commandType, uint 5 numBytes)
    ['0x02' EnableControlCommandSetNetworkVariable0_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x0A' EnableControlCommandSetNetworkVariable1_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x12' EnableControlCommandSetNetworkVariable2_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x1A' EnableControlCommandSetNetworkVariable3_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x22' EnableControlCommandSetNetworkVariable4_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x2A' EnableControlCommandSetNetworkVariable5_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x32' EnableControlCommandSetNetworkVariable6_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x3A' EnableControlCommandSetNetworkVariable7_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x42' EnableControlCommandSetNetworkVariable8_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x4A' EnableControlCommandSetNetworkVariable9_2Bytes       ['SET_NETWORK_VARIABLE',   '2']]
    ['0x52' EnableControlCommandSetNetworkVariable10_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
    ['0x5A' EnableControlCommandSetNetworkVariable11_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
    ['0x62' EnableControlCommandSetNetworkVariable12_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
    ['0x6A' EnableControlCommandSetNetworkVariable13_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
    ['0x72' EnableControlCommandSetNetworkVariable14_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
    ['0x7A' EnableControlCommandSetNetworkVariable15_2Bytes      ['SET_NETWORK_VARIABLE',   '2']]
]

[enum uint 4 EnableControlCommandType
    ['0x00' SET_NETWORK_VARIABLE   ]
]
