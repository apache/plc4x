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

[type Plc4xConstants
    [const          uint 16     plc4xTcpDefaultPort 59837] // Hex of CAFE
]

[discriminatedType Plc4xMessage byteOrder='BIG_ENDIAN'
    [const         uint 8           version      0x01           ]
    [implicit      uint 16          packetLength 'lengthInBytes']
    [simple        uint 16          requestId                   ]
    [discriminator Plc4xRequestType requestType                 ]
    [typeSwitch requestType
        ['CONNECT_REQUEST' Plc4xConnectRequest
            [implicit uint 8                            connectionStringLen 'STR_LEN(connectionString)']
            [simple   vstring 'connectionStringLen * 8' connectionString                               ]
        ]
        ['CONNECT_RESPONSE' Plc4xConnectResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
        ]
        ['READ_REQUEST' Plc4xReadRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array    Plc4xTagRequest         tags       count 'numTags']
        ]
        ['READ_RESPONSE' Plc4xReadResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array    Plc4xTagValueResponse   tags       count 'numTags']
        ]
        ['WRITE_REQUEST' Plc4xWriteRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array    Plc4xTagValueRequest    tags       count 'numTags']
        ]
        ['WRITE_RESPONSE' Plc4xWriteResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array    Plc4xTagResponse        tags       count 'numTags']
        ]
        // TODO: Implement this later on.
        /*['SUBSCRIPTION_REQUEST' Plc4xSubscriptionRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array                            tags       count 'numTags']
        ]
        ['SUBSCRIPTION_RESPONSE' Plc4xSubscriptionResponse
            [simple   uint 16                 connectionId                  ]
            [simple   PlcResponseCode         responseCode                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array                            tags       count 'numTags']
        ]
        ['UNSUBSCRIPTION_REQUEST' Plc4xUnsubscriptionRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array                            tags       count 'numTags']
        ]
        ['UNSUBSCRIPTION_RESPONSE' Plc4xUnsubscriptionResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numTags    'COUNT(tags)'  ]
            [array                            tags       count 'numTags']
        ]*/
    ]
]

[type Plc4xTag
    [implicit uint 8                      nameLen       'STR_LEN(name)'      ]
    [simple   vstring 'nameLen * 8'       name                               ]
    [implicit uint 8                      tagQueryLen 'STR_LEN(tagQuery)']
    [simple   vstring 'tagQueryLen * 8' tagQuery                         ]
]

[type Plc4xTagRequest
    [simple Plc4xTag              tag       ]
]

[type Plc4xTagValueRequest
    [simple   Plc4xTag              tag                                       ]
    [simple   Plc4xValueType          valueType                                   ]
    [optional Plc4xValue('valueType') value     'valueType != Plc4xValueType.NULL']
]

[type Plc4xTagResponse
    [simple Plc4xTag              tag       ]
    [simple Plc4xResponseCode       responseCode]
]

[type Plc4xTagValueResponse
    [simple   Plc4xTag              tag                                          ]
    [simple   Plc4xResponseCode       responseCode                                   ]
    [simple   Plc4xValueType          valueType                                      ]
    [optional Plc4xValue('valueType') value        'valueType != Plc4xValueType.NULL']
]

[dataIo Plc4xValue(Plc4xValueType valueType)
    [typeSwitch valueType
        // Bit Strings
        ['BOOL'          BOOL
            [reserved uint 7                     '0x00'                          ]
            [simple   bit                        value                           ]
        ]
        ['BYTE'          BYTE
            [simple   uint 8                     value                           ]
        ]
        ['WORD'          WORD
            [simple   uint 16                    value                           ]
        ]
        ['DWORD'         DWORD
            [simple   uint 32                    value                           ]
        ]
        ['LWORD'         LWORD
            [simple   uint 64                    value                           ]
        ]

        // Unsigned Integers
        ['USINT'         USINT
            [simple   uint 8                     value                           ]
        ]
        ['UINT'          UINT
            [simple   uint 16                    value                           ]
        ]
        ['UDINT'         UDINT
            [simple   uint 32                    value                           ]
        ]
        ['ULINT'         ULINT
            [simple   uint 64                    value                           ]
        ]

        // Signed Integers
        ['SINT'          SINT
            [simple   int 8                      value                           ]
        ]
        ['INT'           INT
            [simple   int 16                     value                           ]
        ]
        ['DINT'          DINT
            [simple   int 32                     value                           ]
        ]
        ['LINT'          LINT
            [simple   int 64                     value                           ]
        ]

        // Floating Points
        ['REAL'          REAL
            [simple   float 32                   value                           ]
        ]
        ['LREAL'         LREAL
            [simple   float 64                   value                           ]
        ]

        // Chars and Strings
        ['CHAR'          STRING
            [simple   string 8                   value                           ]
        ]
        ['WCHAR'         STRING
            [simple   string 16                  value        encoding='"UTF-16"']
        ]
        ['STRING'        STRING
            [manual vstring value  'STATIC_CALL("parseString", readBuffer, _type.encoding)' 'STATIC_CALL("serializeString", writeBuffer, _value, _type.encoding)' '(STR_LEN(_value) + 1) * 8']
        ]
        ['WSTRING'       STRING
            [manual vstring value  'STATIC_CALL("parseString", readBuffer, _type.encoding)' 'STATIC_CALL("serializeString", writeBuffer, _value, _type.encoding)' '(STR_LEN(_value) + 1) * 16' encoding='"UTF-16"']
        ]

        // Times and Dates
        ['TIME'           TIME
            [simple uint 32 milliseconds]
        ]
        ['LTIME'          LTIME
            [simple uint 64 nanoseconds]
        ]
        ['DATE'           DATE
            [simple uint 32 secondsSinceEpoch]
        ]
        ['LDATE'          LDATE
            [simple uint 64 nanosecondsSinceEpoch]
        ]
        ['TIME_OF_DAY'    TIME_OF_DAY
            [simple uint 32 millisecondsSinceMidnight]
        ]
        ['LTIME_OF_DAY'   LTIME_OF_DAY
            [simple uint 64 nanosecondsSinceMidnight]
        ]
        ['DATE_AND_TIME'  DATE_AND_TIME
            [simple uint 32 secondsSinceEpoch]
        ]
        ['LDATE_AND_TIME' LDATE_AND_TIME
            [simple uint 64 nanosecondsSinceEpoch]
        ]

        // Derived Types
        ['Struct'        Struct       ]
        //['List'          List         ]
    ]
]

[enum uint 8 Plc4xRequestType
    ['0x01' CONNECT_REQUEST        ]
    ['0x02' CONNECT_RESPONSE       ]
    ['0x03' DISCONNECT_REQUEST     ]
    ['0x04' DISCONNECT_RESPONSE    ]
    ['0x05' READ_REQUEST           ]
    ['0x06' READ_RESPONSE          ]
    ['0x07' WRITE_REQUEST          ]
    ['0x08' WRITE_RESPONSE         ]
    ['0x09' SUBSCRIPTION_REQUEST   ]
    ['0x0A' SUBSCRIPTION_RESPONSE  ]
    ['0x0B' UNSUBSCRIPTION_REQUEST ]
    ['0x0C' UNSUBSCRIPTION_RESPONSE]
]

[enum uint 8 Plc4xResponseCode
    ['0x01' OK              ]
    ['0x02' NOT_FOUND       ]
    ['0x03' ACCESS_DENIED   ]
    ['0x04' INVALID_ADDRESS ]
    ['0x06' INVALID_DATATYPE]
    ['0x07' INVALID_DATA    ]
    ['0x08' INTERNAL_ERROR  ]
    ['0x09' REMOTE_BUSY     ]
    ['0x0A' REMOTE_ERROR    ]
    ['0x0B' UNSUPPORTED     ]
    ['0x0C' RESPONSE_PENDING]
]

[enum uint 8 Plc4xValueType
    ['0x00' NULL          ]

    // Bit Strings
    ['0x01' BOOL          ]
    ['0x02' BYTE          ]
    ['0x03' WORD          ]
    ['0x04' DWORD         ]
    ['0x05' LWORD         ]

    // Unsigned Integers
    ['0x11' USINT         ]
    ['0x12' UINT          ]
    ['0x13' UDINT         ]
    ['0x14' ULINT         ]

    // Signed Integers
    ['0x21' SINT          ]
    ['0x22' INT           ]
    ['0x23' DINT          ]
    ['0x24' LINT          ]

    // Floating Point Values
    ['0x31' REAL          ]
    ['0x32' LREAL         ]

    // Chars and Strings
    ['0x41' CHAR          ]
    ['0x42' WCHAR         ]
    ['0x43' STRING        ]
    ['0x44' WSTRING       ]

    // Times and Dates
    ['0x51' TIME          ]
    ['0x52' LTIME         ]
    ['0x53' DATE          ]
    ['0x54' LDATE         ]
    ['0x55' TIME_OF_DAY   ]
    ['0x56' LTIME_OF_DAY  ]
    ['0x57' DATE_AND_TIME ]
    ['0x58' LDATE_AND_TIME]

    // Complex types
    ['0x61' Struct        ]
    ['0x62' List          ]

    ['0x71' RAW_BYTE_ARRAY]
]

[enum uint 8 Plc4xSubscriptionType
   ['0x01' CYCLIC         ]
   ['0x02' CHANGE_OF_STATE]
   ['0x03' EVENT          ]
]
