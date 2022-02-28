/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

[discriminatedType CBusCommand
    [const  uint 8     initiator '0x5C'] // 0x5C = "/"
    [simple CBusHeader header          ]
    [typeSwitch 'header.destinationAddressType'
        ['PointToPointToMultiPoint' CBusCommandPointToPointToMultiPoint
            [simple CBusPointToPointToMultipointCommand command]
        ]
        ['PointToMultiPoint'        CBusCommandPointToMultiPoint
            [simple CBusPointToMultiPointCommand        command]
        ]
        ['PointToPoint'             CBusCommandPointToPoint
            [simple CBusPointToPointCommand             command]
        ]
    ]
]

[type CBusHeader
    [simple   PriorityClass          priorityClass         ]
    [reserved bit                    'false'               ] // Reserved for internal C-Bus management purposes
    [reserved uint 2                 '0'                   ] // Reserved for internal C-Bus management purposes
    [simple   DestinationAddressType destinationAddressType]
]

[discriminatedType CBusPointToPointCommand
    [peek    uint 16     bridgeAddressCountPeek]
    [virtual UnitAddress unitAddress           ]
    [typeSwitch 'bridgeAddressCountPeek && 0x00FF'
        ['0x0000' CBusPointToPointCommandDirect
            [simple UnitAddress   unitAddress                                                  ]
        ]
        [         CBusPointToPointCommandIndirect
            [simple BridgeAddress firstBridgeAddress                                           ]
            [simple RouteType     routeType                                                    ]
            [array  BridgeAddress additionalBridgeAddresses count 'routeType.additionalBridges']
            [simple UnitAddress   unitAddress                                                  ]
        ]
    ]

]

[enum uint 2 PriorityClass
    ['0x00' Class4] // lowest
    ['0x01' Class3] // medium low
    ['0x02' Class2] // medium high
    ['0x03' Class1] // highest
]

[enum uint 3 DestinationAddressType
    ['0x03' PointToPointToMultiPoint] // P-P-M
    ['0x05' PointToMultiPoint       ] // P-M
    ['0x06' PointToPoint            ] // P-P
]

[enum uint 8 RouteType(uint 3 additionalBridges)
    ['0x00' NoBridgeAtAll         ['0']]
    ['0x09' NoAdditionalBridge    ['1']]
    ['0x12' OneAdditionalBridge   ['2']]
    ['0x1B' TwoAdditionalBridge   ['3']]
    ['0x24' ThreeAdditionalBridge ['4']]
    ['0x2D' FourAdditionalBridge  ['4']]
    ['0x36' FiveAdditionalBridge  ['4']]
