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

//////////////////////////////////////////////////////////////////
///EthernetIP Header of size 24
/////////////////////////////////////////////////////////////////

[discriminatedType EipPacket byteOrder='BIG_ENDIAN'
    [discriminator uint 16 command]
    // TODO: was len before but that is a reserved keyword in golang and as long as we don't have a language agnostic neutralizer this clashes
    [implicit      uint 16 packetLength 'lengthInBytes - 24']
    [simple        uint 32 sessionHandle]
    [simple        uint 32 status]
    [array         uint 8  senderContext count '8']
    [simple        uint 32 options]
    [typeSwitch command
            ['0x0065' EipConnectionRequest
                [const  uint    16   protocolVersion   0x01]
                [const  uint    16   flags             0x00]
            ]
            ['0x0066' EipDisconnectRequest
            ]
            ['0x006F' CipRRData(uint 16 packetLength)
                [reserved  uint    32    '0x00000000']
                [reserved  uint    16    '0x0000']
                [simple    CipExchange('packetLength - 6')   exchange]
            ]
        ]
]

[type  CipExchange (uint 16 exchangeLen)  //We pass then length down to evey sub-type to be able to provide the remaining data size
    [const          uint 16                         itemCount           0x02                  ]  //2 items
    [const          uint 32                         nullPtr             0x0                   ]  //NullPointerAddress
    [const          uint 16                         unconnectedData     0x00B2                ]  //Connection Manager
    [implicit       uint 16                         size                'lengthInBytes - 8 - 2' ]  //remove fields above and routing
    [simple         CipService('exchangeLen - 10')  service                                     ]
]

[discriminatedType  CipService(uint 16 serviceLen)
    [discriminator  uint    8   service]
    [typeSwitch service
        ['0x4C' CipReadRequest
            [simple     int     8   requestPathSize]
            [array      byte   tag   length  '(requestPathSize * 2)']
            [simple     uint    16  elementNb]
        ]
        ['0xCC' CipReadResponse
              [reserved   uint            8   '0x00']
              [simple     uint            8   status]
              [simple     uint            8   extStatus]
              [simple     CIPDataTypeCode     dataType]
              [array      byte   data  count  'serviceLen - 6']
        ]
        ['0x4D' CipWriteRequest
            [simple     int     8           requestPathSize]
            [array      byte           tag   length  'requestPathSize * 2']
            [simple     CIPDataTypeCode     dataType]
            [simple     uint    16          elementNb]
            [array      byte            data  length  'dataType.size * elementNb']
        ]
        ['0xCD' CipWriteResponse
            [reserved   uint        8   '0x00']
            [simple     uint        8   status]
            [simple     uint        8   extStatus]
        ]
        ['0x0A' MultipleServiceRequest
               [const  int     8   requestPathSize   0x02]
               [const  uint    32  requestPath       0x01240220]   //Logical Segment: Class(0x20) 0x02, Instance(0x24) 01 (Message Router)
               [simple Services('serviceLen - 6 ')  data ]
        ]
        ['0x8A' MultipleServiceResponse
               [reserved   uint    8   '0x0']
               [simple     uint    8   status]
               [simple     uint    8   extStatus]
               [simple     uint    16  serviceNb]
               [array      uint    16  offsets       count  'serviceNb']
               [array      byte   servicesData count 'serviceLen - 6 - (2 * serviceNb)']
        ]
        ['0x52'   CipUnconnectedRequest
               [reserved   uint    8   '0x02']
               [reserved   uint    8   '0x20']   // setRequestPathLogicalClassSegment
               [reserved   uint    8   '0x06']   // set request class path
               [reserved   uint    8   '0x24']   // setRequestPathLogicalInstanceSegment
               [reserved   uint    8   '0x01']   // setRequestPathInstance
               [reserved   uint    16  '0x9D05']   //Timeout 5s
               [implicit   uint    16  messageSize   'lengthInBytes - 10 - 4']   //subtract above and routing
               [simple     CipService('messageSize')  unconnectedService ]
               [const      uint    16  route 0x0001]
               [simple     int     8   backPlane]
               [simple     int     8   slot]
        ]
    ]
]

[type   Services  (uint   16   servicesLen)
    [simple uint        16  serviceNb]
    [array  uint        16  offsets       count  'serviceNb']
    [array  CipService('servicesLen / serviceNb')   services    count  'serviceNb' ]
]

[enum uint   16   CIPDataTypeCode(uint 8  size)
    ['0X00C1'   BOOL            ['1']]
    ['0X00C2'   SINT            ['1']]
    ['0X00C3'   INT             ['2']]
    ['0X00C4'   DINT            ['4']]
    ['0X00C5'   LINT            ['8']]
    ['0X00CA'   REAL            ['4']]
    ['0X00D3'   DWORD           ['4']]
    ['0X02A0'   Struct          ['88']]
    ['0X02A0'   STRING          ['88']]
    //['0X02A0'   STRING36        ['40']]
    //TODO: -1 is not a valid value for uint
    //['-1'       UNKNOWN         ['-1']]
]

[enum uint 16 CIPStructTypeCode
    ['0x0FCE'   STRING]
]

[enum   uint    16  EiPCommand
    ['0x0065'   RegisterSession ]
    ['0x0066'   UnregisterSession ]
    ['0x006F'   SendRRData ]
]