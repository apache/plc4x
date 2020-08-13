//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

////////////////////////////////////////////////////////////////
// AMS/TCP Packet
////////////////////////////////////////////////////////////////

[type 'AmsTCPDiscoveryPacket'
    [const uint 32 'header' '0x03661471']
    [reserved   uint 32  '0x00000000']
    [simple     uint 16  'something']
    [reserved   uint 16  '0x0000']
    [simple     AmsNetId 'amsNetId']
    [reserved   uint 16  '0x1027']
    [simple     uint 16  'something2']
    [reserved   uint 32  '0x0000']
    [simple     uint 16  'something3']
    [implicit   uint 16  'hostnameLength' 'COUNT(hostname)']
    [array      int 8    'hostname' COUNT 'hostnameLength']
    [reserved   uint 16  '0x4000']
    [simple     uint 16  'something4']
    [simple     uint 16  'something5']
    [reserved   uint 32  '0x0000']
    [simple     uint 16  'something6']
    [reserved   uint 32  '0x0000']
    [reserved   uint 32  '0x00000000']
    [reserved   uint 32  '0x00000000']
    [reserved   uint 16  '0x0000']
    [implicit   uint 16  'ipLength' 'COUNT(ip)']
    [array      int 8    'ip' COUNT 'ipLength']
]

[type 'AmsNetId'
    [simple     uint        8   'octet1'            ]
    [simple     uint        8   'octet2'            ]
    [simple     uint        8   'octet3'            ]
    [simple     uint        8   'octet4'            ]
    [simple     uint        8   'octet5'            ]
    [simple     uint        8   'octet6'            ]
]
