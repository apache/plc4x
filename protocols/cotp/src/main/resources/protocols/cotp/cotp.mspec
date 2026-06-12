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

////////////////////////////////////////////////////////////////
// IsoOnTcp/TPKT
////////////////////////////////////////////////////////////////

[type TPKTPacket byteOrder='"BIG_ENDIAN"' unsignedIntegerEncoding='"unsigned-binary"' signedIntegerEncoding='"twos-complement"' floatEncoding='"IEEE754"' stringEncoding='"UTF8"'
    [const    uint 8                 protocolId 0x03                                                                   ]
    [reserved uint 8                 '0x00'                                                                            ]
    [implicit uint 16                len        'payload.lengthInBytes + 4'                                            ]
    [simple   COTPPacket('len - 4')  payload                                                                           ]
]

////////////////////////////////////////////////////////////////
// COTP
////////////////////////////////////////////////////////////////

[discriminatedType COTPPacket (uint 16 cotpLen)
    [implicit      uint 8 headerLength 'lengthInBytes - (((payload != null) ? COUNT(payload) : 0) + 1)'                ]
    [discriminator uint 8 tpduCode                                                                                     ]
    [typeSwitch tpduCode
        ['0xF0' COTPPacketData
            [simple bit               eot                                                                              ]
            [simple uint 7            tpduRef                                                                          ]
        ]
        ['0xE0' COTPPacketConnectionRequest
            [simple uint 16           destinationReference                                                             ]
            [simple uint 16           sourceReference                                                                  ]
            [simple COTPProtocolClass protocolClass                                                                    ]
        ]
        ['0xD0' COTPPacketConnectionResponse
            [simple uint 16           destinationReference                                                             ]
            [simple uint 16           sourceReference                                                                  ]
            [simple COTPProtocolClass protocolClass                                                                    ]
        ]
        ['0x80' COTPPacketDisconnectRequest
            [simple uint 16           destinationReference                                                             ]
            [simple uint 16           sourceReference                                                                  ]
            [simple COTPProtocolClass protocolClass                                                                    ]
        ]
        ['0xC0' COTPPacketDisconnectResponse
            [simple uint 16           destinationReference                                                             ]
            [simple uint 16           sourceReference                                                                  ]
        ]
        ['0x70' COTPPacketTpduError
            [simple uint 16           destinationReference                                                             ]
            [simple uint 8            rejectCause                                                                      ]
        ]
    ]
    [array    COTPParameter ('(headerLength + 1) - (curPos / 8)') parameters length '(headerLength + 1) - (curPos / 8)']
    [array    byte                                                payload    count 'cotpLen - (curPos / 8)'            ]
]

[discriminatedType COTPParameter (uint 8 restSizeInBytes)
    [discriminator uint 8        parameterType                                                                         ]
    [implicit      uint 8        parameterLength    'lengthInBytes - 2'                                                ]
    [typeSwitch parameterType
        ['0xC0' COTPParameterTpduSize
            [simple COTPTpduSize tpduSize                                                                              ]
        ]
        ['0xC1' COTPParameterCallingTsap (uint 8 parameterLength)
            [array  byte         tsapId             count               'parameterLength'                              ]
        ]
        ['0xC2' COTPParameterCalledTsap (uint 8 parameterLength)
            [array  byte         tsapId             count               'parameterLength'                              ]
        ]
        ['0xC3' COTPParameterChecksum
            [simple uint 8       crc                                                                                   ]
        ]
        // In general the way the parameter is constructed, it would gobble up all bytes of following parameters
        // However, this is only used on disconnects and always seems to be coming on it's own. So it shouldn't
        // matter.
        ['0xE0' COTPParameterDisconnectAdditionalInformation
            [array byte          data               count                   'restSizeInBytes'                          ]
        ]
    ]
]


[enum uint 8 COTPTpduSize(uint 16 sizeInBytes)
    ['0x07' SIZE_128     ['128'              ]                                                                         ]
    ['0x08' SIZE_256     ['256'              ]                                                                         ]
    ['0x09' SIZE_512     ['512'              ]                                                                         ]
    ['0x0a' SIZE_1024    ['1024'             ]                                                                         ]
    ['0x0b' SIZE_2048    ['2048'             ]                                                                         ]
    ['0x0c' SIZE_4096    ['4096'             ]                                                                         ]
    ['0x0d' SIZE_8192    ['8192'             ]                                                                         ]
]

[enum uint 8 COTPProtocolClass
    ['0x00' CLASS_0                                                                                                    ]
    ['0x10' CLASS_1                                                                                                    ]
    ['0x20' CLASS_2                                                                                                    ]
    ['0x30' CLASS_3                                                                                                    ]
    ['0x40' CLASS_4                                                                                                    ]
]

