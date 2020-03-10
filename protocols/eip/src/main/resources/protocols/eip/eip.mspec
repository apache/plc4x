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

//////////////////////////////////////////////////////////////////
///EthernetIP Header of size 24
/////////////////////////////////////////////////////////////////
[discriminatedType 'EipPacket'
    [discriminator uint 16 'command']
    [implicit      uint 16 'len' 'lengthInBytes - 24']
    [simple        uint 32 'sessionHandle']
    [simple        uint 32 'status']
    [array         uint 8  'senderContext' count '8']
    [simple        uint 32 'options']
    [typeSwitch 'command'
            ['0x0065' EipConnectionRequest
                [const  uint    16   'protocolVersion'   '0x01']
                [const  uint    16   'flags'             '0x00']
            ]
            ['0x0066' EipDisconnectRequest
            ]
        ]
]

[enum int   16   'CIPDataTypeCode' [uint 8  'size']
    ['0X00C1'   BOOL            ['1']]
    ['0X00CA'   REAL            ['4']]
    ['0X00C4'   DINT            ['4']]
    ['0X00C3'   INT             ['2']]
    ['0X00C2'   SINT            ['1']]
    ['0X02A0'   STRUCTURED      ['88']]
    ['0X02A0'   STRING          ['88']]
    ['0X02A0'   STRING36        ['40']]
    ['-1'       UNKNOWN         ['-1']]
]