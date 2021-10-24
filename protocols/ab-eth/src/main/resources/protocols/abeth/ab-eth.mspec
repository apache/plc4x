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

[discriminatedType 'CIPEncapsulationPacket' byteOrder='"BIG_ENDIAN"'
    [discriminator uint 16 'commandType']
    [implicit      uint 16 'len' 'lengthInBytes - 28']
    [simple        uint 32 'sessionHandle']
    [simple        uint 32 'status']
    [array         uint 8  'senderContext' count '8']
    [simple        uint 32 'options']
    [reserved      uint 32 '0x00000000']
    [typeSwitch 'commandType'
        ['0x0101' CIPEncapsulationConnectionRequest
        ]
        ['0x0201' CIPEncapsulationConnectionResponse
        ]
        ['0x0107' CIPEncapsulationReadRequest
            [simple   DF1RequestMessage  'request']
        ]
        ['0x0207' CIPEncapsulationReadResponse (uint 16 'len')
            [simple   DF1ResponseMessage('len') 'response']
        ]
    ]
]

[discriminatedType 'DF1RequestMessage'
    [simple        uint 8  'destinationAddress']
    [simple        uint 8  'sourceAddress']
    [reserved      uint 16 '0x0000']
    [discriminator uint 8  'commandCode']
    [simple        uint 8  'status']
    [simple        uint 16 'transactionCounter']
    [typeSwitch 'commandCode'
        ['0x0F' DF1CommandRequestMessage
            [simple DF1RequestCommand 'command']
        ]
    ]
]

[discriminatedType 'DF1ResponseMessage' (uint 16 'payloadLength')
    [reserved      uint 8  '0x00']
    [simple        uint 8  'destinationAddress']
    [simple        uint 8  'sourceAddress']
    [reserved      uint 8  '0x00']
    [discriminator uint 8  'commandCode']
    [simple        uint 8  'status']
    [simple        uint 16 'transactionCounter']
    [typeSwitch 'commandCode'
        ['0x4F' DF1CommandResponseMessageProtectedTypedLogicalRead (uint 16 'payloadLength', uint 8 'status')
            [array    uint 8 'data' length 'payloadLength - 8']
        ]
    ]
]

[discriminatedType 'DF1RequestCommand'
    [discriminator    uint 8 'functionCode']
    [typeSwitch 'functionCode'
        ['0xA2' DF1RequestProtectedTypedLogicalRead
            [simple uint 8 'byteSize']
            [simple uint 8 'fileNumber']
            [simple uint 8 'fileType']
            [simple uint 8 'elementNumber']
            [simple uint 8 'subElementNumber']
        ]
    ]
]


