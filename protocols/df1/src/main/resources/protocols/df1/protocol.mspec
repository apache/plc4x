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


[discriminatedType 'DF1Message' [uint 8 'payloadSize']
    [const    uint 8     'messageStart' '0x10']
    [field    uint 8     'transmissionMode']
    [field    uint 8     'destinationAddress']
    [field    uint 8     'sourceAddress']
    [discriminator uint 8 'command']
    [field    uint 8     'status']
    [field    uint 16    'transactioncounter']
    [typeSwitch 'command'
        ['0x01' DF1ReadRequest
         [field uint 16   'address']
         [field uint 8    'size']
        ]
        ['0x41' DF1ReadResponse
         [arrayField uint 8 'data' length 'payloadSize']
        ]
    ]
    [const    uint 8     'messageEnd' '0x10']
    [const    uint 8     'messageTransmissionEnd' '0x03']
    [implicit uint 16    'crc' 'CRC(destinationAddress, sourceAddress, discriminatorValues, status, transactioncounter)']
]

