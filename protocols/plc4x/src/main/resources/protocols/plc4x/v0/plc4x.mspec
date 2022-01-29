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

[discriminatedType Plc4xRequest byteOrder='BIG_ENDIAN'
    [const         uint 8           version 0x01                       ]
    [implicit      uint 16          packetLength 'lengthInBytes'       ]
    [discriminator Plc4xRequestType requestType                        ]
    [typeSwitch requestType
        ['READ_REQUEST' Plc4xReadRequest
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array  Plc4xFieldRequest       fields    count 'numFields']
        ]
        ['READ_RESPONSE' Plc4xReadResponse
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array  Plc4xFieldValueResponse fields    count 'numFields']
        ]
        ['WRITE_REQUEST' Plc4xWriteRequest
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array  Plc4xFieldValueRequest  fields    count 'numFields']
        ]
        ['WRITE_RESPONSE' Plc4xWriteResponse
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array  Plc4xFieldResponse      fields    count 'numFields']
        ]
        // TODO: Implement this later on.
        /*['SUBSCRIPTION_REQUEST' Plc4xSubscriptionRequest
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array                          fields    count 'numFields']
        ]
        ['SUBSCRIPTION_RESPONSE' Plc4xSubscriptionResponse
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array                          fields    count 'numFields']
        ]
        ['UNSUBSCRIPTION_REQUEST' Plc4xUnsubscriptionRequest
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array                          fields    count 'numFields']
        ]
        ['UNSUBSCRIPTION_RESPONSE' Plc4xUnsubscriptionResponse
            [simple uint 16                 requestId                  ]
            [simple uint 8                  numFields                  ]
            [array                          fields    count 'numFields']
        ]*/
    ]
]

[type Plc4xField
    [simple uint 8                  nameLen      ]
    [simple vstring 'nameLen'       name         ]
    [simple uint 8                  fieldQueryLen]
    [simple vstring 'fieldQueryLen' fieldQuery   ]
]

[type Plc4xFieldRequest
    [simple Plc4xField              field       ]
]

[type Plc4xFieldValueRequest
    [simple Plc4xField              field       ]
    [simple Plc4xValueType          valueType   ]
    [simple Plc4xValue('valueType') value       ]
]

[type Plc4xFieldResponse
    [simple Plc4xField              field       ]
    [simple Plc4xResponseCode       responseCode]
]

[type Plc4xFieldValueResponse
    [simple Plc4xField              field       ]
    [simple Plc4xResponseCode       responseCode]
    [simple Plc4xValueType          valueType   ]
    [simple Plc4xValue('valueType') value       ]
]

[dataIo Plc4xValue(Plc4xValueType valueType)
    [typeSwitch valueType
        // Bit Strings
        ['BOOL'          BOOL         ]
        ['BYTE'          BYTE         ]
        ['WORD'          WORD         ]
        ['DWORD'         DWORD        ]

        // Unsigned Integers
        ['USINT'         USINT        ]
        ['UINT'          UINT         ]
        ['UDINT'         UDINT        ]
        ['ULINT'         ULINT        ]

        // Signed Integers
        ['SINT'          SINT         ]
        ['INT'           INT          ]
        ['DINT'          DINT         ]
        ['LINT'          LINT         ]

        // Floating Points
        ['REAL'          REAL         ]
        ['LREAL'         LREAL        ]

        // Chars and Strings
        ['CHAR'          CHAR         ]
        ['WCHAR'         WCHAR        ]
        ['STRING'        STRING       ]
        ['WSTRING'       WSTRING      ]

        // Times and Dates
        ['TIME'          TIME         ]
        ['TIME_OF_DAY'   TIME_OF_DAY  ]
        ['DATE'          DATE         ]
        ['DATE_AND_TIME' DATE_AND_TIME]

        // Derived Types
        ['Struct'        Struct       ]
        ['List'          List         ]
    ]
]

[enum uint 8 Plc4xRequestType
    ['0x01' READ_REQUEST           ]
    ['0x02' READ_RESPONSE          ]
    ['0x03' WRITE_REQUEST          ]
    ['0x04' WRITE_RESPONSE         ]
    ['0x05' SUBSCRIPTION_REQUEST   ]
    ['0x06' SUBSCRIPTION_RESPONSE  ]
    ['0x07' UNSUBSCRIPTION_REQUEST ]
    ['0x08' UNSUBSCRIPTION_RESPONSE]
]

[enum uint 8 Plc4xValueType
    // Bit Strings
    ['0x01' BOOL         ]
    ['0x02' BYTE         ]
    ['0x03' WORD         ]
    ['0x04' DWORD        ]

    // Unsigned Integers
    ['0x11' USINT        ]
    ['0x12' UINT         ]
    ['0x13' UDINT        ]
    ['0x14' ULINT        ]

    // Signed Integers
    ['0x21' SINT         ]
    ['0x22' INT          ]
    ['0x23' DINT         ]
    ['0x24' LINT         ]

    // Floating Point Values
    ['0x31' REAL         ]
    ['0x32' LREAL        ]

    // Chars and Strings
    ['0x41' CHAR         ]
    ['0x42' WCHAR        ]
    ['0x43' STRING       ]
    ['0x44' WSTRING      ]

    // Times and Dates
    ['0x51' TIME         ]
    ['0x52' TIME_OF_DAY  ]
    ['0x53' DATE         ]
    ['0x54' DATE_AND_TIME]

    // Complex types
    ['0x61' Struct       ]
    ['0x62' List         ]
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