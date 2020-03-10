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

[type 'EipPkt'
    [simple uint    16      'command']
    [simple uint    16      'len']
    [simple uint    32      'sessionHandle']
    [simple uint    32      'status']
    [simple uint    64      'senderContext']
    [simple uint    32      'options']
    [optional EipData         'commandData']
]

[discriminatedType 'EipData'
    [simple uint    32      'interface']
    [simple uint    16      'timeout']
    [simple uint    16      'itemsCount']
    [array  EipItem         'items' length  'itemsCount']
    [switchType 'interface'
        ['0x00000000'   CipInterface
        ]
    ]
]

[discriminatedType  'CipInterface'
    [simple uint    8   'service']
    [simple uint    8   'pathSize']
    [simple RequestPath 'path']
    [typeSwitch 'service'
        ['0x52' CipUnconnectedRead
            [simple ReadData    'data']
        ]
    ]
]
[discriminatedType  'ReadData'
    [constant   uint    8       'priority'  '0x05']
    [simple     uint    16      'timeout']
    [simple     uint    16      'requestSize']
    [simple     Request         'request']
]

[discriminatedType  'Request'
    [simple uint    8   'service']
    [simple uint    8   'pathSize']
    [array  uint    'requestPath'   length  'pathSize']
    [simple uint    16   'commandData']
    [const  uint    8   'reserved'  '0x00']
    [simple uint    16  'route']
]


[discriminatedType  'RequestPath'
    [uint   8   'size']
    [array  PathSegment length  'size']
]

[discriminatedType  'PathSegment'
    [simple uint    8   'segment']
]

[discriminatedType  'EipItem'
    [simple uint    16  'type']
    [simple uint    16  'length']
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
