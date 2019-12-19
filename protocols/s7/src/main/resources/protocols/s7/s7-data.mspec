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

[discriminatedType 'DataItem' [uint 8 'dataProtocolId']
    [typeSwitch 'dataProtocolId'
        // -----------------------------------------
        // Bit
        // -----------------------------------------
        ['01' DataItemBOOL
            [reserved uint 7 '0x00']
            [simple   bit    'value']
        ]

        // -----------------------------------------
        // Bit-strings
        // -----------------------------------------
        // 1 byte
        ['11' DataItemBYTE
            [array bit 'value' count '8']
        ]
        // 2 byte (16 bit)
        ['12' DataItemWORD
            [array bit 'value' count '16']
        ]
        // 4 byte (32 bit)
        ['13' DataItemDWORD
            [array bit 'value' count '32']
        ]
        // 8 byte (64 bit)
        ['14' DataItemLWORD
            [array bit 'value' count '64']
        ]

        // -----------------------------------------
        // Integers
        // -----------------------------------------
        // 8 bit:
        ['21' DataItemSINT
            [simple int 8 'value']
        ]
        ['22' DataItemUSINT
            [simple uint 8 'value']
        ]
        // 16 bit:
        ['23' DataItemINT
            [simple int 16 'value']
        ]
        ['24' DataItemUINT
            [simple uint 16 'value']
        ]
        // 32 bit:
        ['25' DataItemDINT
            [simple int 32 'value']
        ]
        ['26' DataItemUDINT
            [simple uint 32 'value']
        ]
        // 64 bit:
        ['27' DataItemLINT
            [simple int 64 'value']
        ]
        ['28' DataItemULINT
            [simple uint 64 'value']
        ]

        // -----------------------------------------
        // Floating point values
        // -----------------------------------------
        ['31' DataItemREAL
            [simple float 8.23  'value']
        ]
        ['32' DataItemLREAL
            [simple float 11.52 'value']
        ]

        // -----------------------------------------
        // Characters & Strings
        // -----------------------------------------
        ['41' DataItemCHAR
        ]
        ['42' DataItemWCHAR
        ]
        ['43' DataItemSTRING
        ]
        ['44' DataItemWSTRING
        ]

        // -----------------------------------------
        // TIA Date-Formats
        // -----------------------------------------
        ['51' DataItemTime
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.parseTiaTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.serializeTiaTime", io, _value.value)' '4']
        ]
        // TODO: Check if this is really 8 bytes
        ['52' DataItemLTime
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.parseTiaLTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.serializeTiaLTime", io, _value.value)' '8']
        ]
        ['53' DataItemDate
            [manual date 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.parseTiaDate", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.serializeTiaDate", io, _value.value)' '2']
        ]
        ['54' DataItemTimeOfDay
            [manual time 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.parseTiaTimeOfDay", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.serializeTiaTimeOfDay", io, _value.value)' '4']
        ]
        ['55' DataItemDateAndTime
            [manual dateTime 'value' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.parseTiaDateTime", io)' 'STATIC_CALL("org.apache.plc4x.java.s7.readwrite.utils.StaticHelper.serializeTiaDateTime", io, _value.value)' '8']
        ]
    ]
]

