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

[type DceRpc_Packet byteOrder='BIG_ENDIAN'
// RPC Header {
    // RPCVersion 4.10.3.2.1
    [const         uint 8                version                        0x04                 ]
    // RPCPacketType 4.10.3.2.2 (8 bit)
    [simple        DceRpc_PacketType     packetType                                            ]
    // PRCFlags 4.10.3.2.3
    [reserved      bit                                                  'false'                ]
    [const         bit                   broadcast                      false                ]
    [simple        bit                   idempotent                                            ]
    [const         bit                   maybe                          false                ]
    [simple        bit                   noFragmentAcknowledgeRequested                        ]
    [const         bit                   fragment                       false                ]
    [simple        bit                   lastFragment                                          ]
    [reserved      bit                                                  'false'                ]
    // PRCFlags2 4.10.3.2.4
    [reserved      uint 6                                               '0x00'                 ]
    [const         bit                   cancelWasPending               false                ]
    [reserved      bit                                                  'false'                ]
    // RPCDRep 4.10.3.2.5 (4 bit & 4 bit)
    [simple        IntegerEncoding       integerEncoding                                       ]
    [simple        CharacterEncoding     characterEncoding                                     ]
    // RPCDRep2 4.10.3.2.5 (8 bit)
    [simple        FloatingPointEncoding floatingPointEncoding                                 ]
    // RPCDRep3 (8 bit shall be 0)
    [reserved      uint 8                                               '0x00'                 ]
    // RPCSerialHigh 4.10.3.2.6
    [const         uint 8                serialHigh                     0x00                 ]
    [batchSet byteOrder='integerEncoding == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN'
        // RPCObjectUUID 4.10.3.2.8
        // RPCObjectUUID 4.10.3.2.8
        [simple DceRpc_ObjectUuid        objectUuid                                            ]
        // RPCInterfaceUUID 4.10.3.2.9
        [simple DceRpc_InterfaceUuid     interfaceUuid                                         ]
        // RPCActivityUUID 4.10.3.2.10
        [simple DceRpc_ActivityUuid      activityUuid                                          ]
        // RPCServerBootTime 4.10.3.2.11
        [simple uint 32                  serverBootTime                                        ]
        // RPCInterfaceVersion 4.10.3.2.12
        [const  uint 32                  interfaceVer                   0x00000001           ]
        // RPCSequenceNmb 4.10.3.2.13
        [simple uint 32                  sequenceNumber                                        ]
        // RPCOperationNmb 4.10.3.2.14
        [simple DceRpc_Operation         operation                                             ]
        // RPCInterfaceHint 4.10.3.2.15
        [const        uint 16            interfaceHint                  0xFFFF               ]
        // RPCActivityHint 4.10.3.2.16
        [const        uint 16            activityHint                   0xFFFF               ]
        // RPCLengthOfBody 4.10.3.2.17
        [implicit     uint 16            lengthOfBody                   'payload.lengthInBytes']
        // RPCFragmentNmb 4.10.3.2.18 (Setting this to 0 as we will probably never have anything but 0 here
        [const        uint 16            fragmentNum                    0x0000               ]
        // RPCAuthenticationProtocol 4.10.3.2.19
        [const        uint 8             authProto                      0x00                 ]
    ]
    // RPCSerialLow 4.10.3.2.7
    [const            uint 8             serialLow                      0x00                 ]
// RPC Header }
// RPC Payload {
    [simple PnIoCm_Packet('packetType') payload byteOrder='integerEncoding == IntegerEncoding.BIG_ENDIAN ? BIG_ENDIAN : LITTLE_ENDIAN' ]
// RPC Payload }
]

// RPCObjectUUID 4.10.3.2.8
[type DceRpc_ObjectUuid
    [const  uint 32 data1      0xDEA00000                       ]
    [const  uint 16 data2      0x6C97                           ]
    [const  uint 16 data3      0x11D1                           ]
    // This part is described as a byte array, so the byte order is always big-endian
    [const  uint 16 data4      0x8271     byteOrder='BIG_ENDIAN']
    [simple uint 4  interfaceNumber       byteOrder='BIG_ENDIAN']
    [simple uint 12 nodeNumber            byteOrder='BIG_ENDIAN']
    [simple uint 16 deviceId              byteOrder='BIG_ENDIAN']
    [simple uint 16 vendorId              byteOrder='BIG_ENDIAN']
]

// RPCInterfaceUUID 4.10.3.2.9
// NOTE: If we would have been only using Big Endian encoding, we would have
//       implemented this via an enum. However as the first 8 bytes are
//       dynamically endianed and the last 8 bytes are set to Big Endian, we
//       had to do this trick.
[discriminatedType DceRpc_InterfaceUuid
    [discriminator  uint 32 interfaceType                               ]
    [const          uint 16 data1      0x6C97                           ]
    [const          uint 16 data2      0x11D1                           ]
    // This part is described as a byte array, so the byte order is always big-endian
    [const          uint 16 data3      0x8271     byteOrder='BIG_ENDIAN']
    [const          uint 16 data4      0x00A0     byteOrder='BIG_ENDIAN']
    [const          uint 16 data5      0x2442     byteOrder='BIG_ENDIAN']
    [const          uint 16 data6      0xDF7D     byteOrder='BIG_ENDIAN']
    [typeSwitch interfaceType
        ['0xDEA00001' DceRpc_InterfaceUuid_DeviceInterface
        ]
        ['0xDEA00002' DceRpc_InterfaceUuid_ControllerInterface
        ]
        ['0xDEA00003' DceRpc_InterfaceUuid_SupervisorInterface
        ]
        ['0xDEA00004' DceRpc_InterfaceUuid_ParameterInterface
        ]
    ]
]

// RPCActivityUUID 4.10.3.2.10
// NOTE: This value is generally randomly generated by the initiator
//       and used throughout the entire communication. Unfortunately,
//       the first parts are effected by endianess, and the last 8
//       bytes are fixed big-endian. Therefore the complicated notation.
[type DceRpc_ActivityUuid
    [simple  uint 32 data1          ]
    [simple  uint 16 data2          ]
    [simple  uint 16 data3          ]
    // This part is described as a byte array, so the byte order is always big-endian
    [array   byte    data4 count '8']
]

// 4.10.3.2.2
[enum uint 8 DceRpc_PacketType
    ['0x00' REQUEST              ]
    ['0x01' PING                 ]
    ['0x02' RESPONSE             ]
    ['0x03' FAULT                ]
    ['0x04' WORKING              ]
    // Response to PING
    ['0x05' NO_CALL              ]
    ['0x06' REJECT               ]
    ['0x07' ACKNOWLEDGE          ]
    ['0x08' CONNECTIONLESS_CANCEL]
    ['0x09' FRAGMENT_ACKNOWLEDGE ]
    ['0x0A' CANCEL_ACKNOWLEDGE   ]
]

// 4.10.3.2.14
[enum uint 16 DceRpc_Operation
    ['0x0000' CONNECT      ]
    ['0x0001' RELEASE      ]
    ['0x0002' READ         ]
    ['0x0003' WRITE        ]
    ['0x0004' CONTROL      ]
    ['0x0005' READ_IMPLICIT]
]
