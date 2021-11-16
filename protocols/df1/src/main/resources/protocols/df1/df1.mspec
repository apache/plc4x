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

[discriminatedType DF1Symbol byteOrder='"BIG_ENDIAN"'
    [const            uint  8      messageStart '0x10']
    [discriminator    uint  8      symbolType]
    [typeSwitch 'symbolType'
        ['0x02' DF1SymbolMessageFrame
            [simple   uint  8      destinationAddress]
            [simple   uint  8      sourceAddress]
            [simple   DF1Command   command]
            [const    uint  8      messageEnd '0x10']
            [const    uint  8      endTransaction '0x03']
            [checksum uint 16      crc 'STATIC_CALL("crcCheck", destinationAddress, sourceAddress, command)']
        ]
        ['0x06' DF1SymbolMessageFrameACK
        ]
        ['0x15' DF1SymbolMessageFrameNAK
        ]
    ]
]

[discriminatedType DF1Command
    [discriminator  uint  8     commandCode         ]
    [simple         uint  8     status              ]
    [simple         uint 16     transactionCounter  ]
    [typeSwitch 'commandCode'
        ['0x01' DF1UnprotectedReadRequest
            [simple uint 16    address  ]
            [simple uint  8    size     ]
        ]
        ['0x41' DF1UnprotectedReadResponse
            [manualArray byte data terminated 'STATIC_CALL("dataTerminate", readBuffer)' 'STATIC_CALL("readData", readBuffer)' 'STATIC_CALL("writeData", writeBuffer, _value)' 'STATIC_CALL("dataLength", data)']
        ]
    ]
]

