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

////////////////////////////////////////////////////////////////
// AMS/TCP Packet
////////////////////////////////////////////////////////////////

[discriminatedType AdsDiscovery byteOrder='BIG_ENDIAN'
    [const uint 32 header '0x03661471L']
    [reserved   uint 32  '0x00000000L']
    [simple     Operation operation]
    [reserved   uint 16  '0x0000']
    [simple     Direction direction]
    [typeSwitch 'operation', 'direction'
        ['DISCOVERY', 'REQUEST' DiscoveryRequest
            [simple AmsNetId amsNetId]
            [reserved uint 16 '0x1027']
            [reserved uint 32 '0x00000000L']
        ]
        ['DISCOVERY', 'RESPONSE' DiscoveryResponse
            [simple AmsNetId amsNetId]
            [reserved uint 16 '0x1027']
            [reserved uint 16 '0x0400']
            [reserved uint 24 '0x000005L']
            [simple AmsMagicString name]
        ]
        ['ROUTE', 'REQUEST' RouteRequest
            [simple     AmsNetId sender]
            [reserved   uint 16  '0x1027']
            [reserved   uint 16  '0x0500']
            [reserved   uint 24  '0x000C']
            [simple AmsMagicString routeName ]
            [reserved   uint 16 '0x0700']
            [implicit   uint 8 amsSize 'target.lengthInBytes']
            [const uint 8 targetPrefix '0x00']
            [simple AmsNetId target]
            [const uint 8 usernamePrefix '0x0D']
            [simple AmsMagicString username]
            [const uint 8 passwordPrefix '0x02']
            [simple AmsMagicString password]
            [const uint 8 routePrefix '0x05']
            [simple AmsMagicString address]

        ]
        ['ROUTE', 'RESPONSE' RouteResponse
            [simple AmsNetId amsNetId]
            [reserved uint 16 '0x1027']
            [reserved uint 16 '0x0100']
            [reserved uint 32 '0x00000100']
            [simple   RouteStatus status]
            [reserved uint 24 '0x000000']
        ]
    ]
]

[enum uint 8 Operation
    ['0x01' DISCOVERY]
    ['0x06' ROUTE    ]
]

[enum uint 8 Direction
    ['0x00' REQUEST ]
    ['0x80' RESPONSE]
]

[enum uint 24 RouteStatus
    ['0x040000' SUCCESS]
    ['0x000407' FAILURE]
]

[type AmsMagicString
    [implicit uint 16 len 'COUNT(text) + 1']
    [reserved uint 8 '0x00']
    [array int 8 text count 'len - 1']
    [reserved uint 8 '0x00']
]

[type AmsNetId
    [simple     uint        8   octet1            ]
    [simple     uint        8   octet2            ]
    [simple     uint        8   octet3            ]
    [simple     uint        8   octet4            ]
    [simple     uint        8   octet5            ]
    [simple     uint        8   octet6            ]
]
