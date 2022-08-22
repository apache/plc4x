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
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array    Plc4xFieldRequest       fields       count 'numFields']
        ]
        ['READ_RESPONSE' Plc4xReadResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array    Plc4xFieldValueResponse fields       count 'numFields']
        ]
        ['WRITE_REQUEST' Plc4xWriteRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array    Plc4xFieldValueRequest  fields       count 'numFields']
        ]
        ['WRITE_RESPONSE' Plc4xWriteResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array    Plc4xFieldResponse      fields       count 'numFields']
        ]
        // TODO: Implement this later on.
        /*['SUBSCRIPTION_REQUEST' Plc4xSubscriptionRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array                            fields       count 'numFields']
        ]
        ['SUBSCRIPTION_RESPONSE' Plc4xSubscriptionResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array                            fields       count 'numFields']
        ]
        ['UNSUBSCRIPTION_REQUEST' Plc4xUnsubscriptionRequest
            [simple   uint 16                 connectionId                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array                            fields       count 'numFields']
        ]
        ['UNSUBSCRIPTION_RESPONSE' Plc4xUnsubscriptionResponse
            [simple   uint 16                 connectionId                  ]
            [simple   Plc4xResponseCode       responseCode                  ]
            [implicit uint 8                  numFields    'COUNT(fields)'  ]
            [array                            fields       count 'numFields']
        ]*/
    ]
]

[type Plc4xField
    [implicit uint 8                      nameLen       'STR_LEN(name)'      ]
    [simple   vstring 'nameLen * 8'       name                               ]
    [implicit uint 8                      fieldQueryLen 'STR_LEN(fieldQuery)']
    [simple   vstring 'fieldQueryLen * 8' fieldQuery                         ]
]

[type Plc4xFieldRequest
    [simple Plc4xField              field       ]
]

[type Plc4xFieldValueRequest
    [simple   Plc4xField              field                                       ]
    [simple   Plc4xValueType          valueType                                   ]
    [optional Plc4xValue('valueType') value     'valueType != Plc4xValueType.NULL']
]

[type Plc4xFieldResponse
    [simple Plc4xField              field       ]
    [simple Plc4xResponseCode       responseCode]
]

[type Plc4xFieldValueResponse
    [simple   Plc4xField              field                                          ]
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
        ['BYTE'          List
            [array    bit                        value     count '8'             ]
        ]
        ['WORD'          List
            [array    bit                        value     count '16'            ]
        ]
        ['DWORD'         List
            [array    bit                        value     count '32'            ]
        ]
        ['LWORD'         List
            [array    bit                        value     count '64'            ]
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
//            [implicit uint 8                     stringLength 'STR_LEN(value)'   ]
//            [simple   vstring 'stringLength'     value                           ]
        ]
        ['WSTRING'       STRING
//            [implicit uint 8                     stringLength 'STR_LEN(value)'   ]
//            [simple   vstring 'stringLength * 2' value        encoding='"UTF-16"']
        ]

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