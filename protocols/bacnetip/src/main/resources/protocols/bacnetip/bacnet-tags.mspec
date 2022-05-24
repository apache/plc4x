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

[discriminatedType BACnetTagHeader
    [simple        uint 4   tagNumber                                                                                   ]
    [simple        TagClass tagClass                                                                                    ]
    [simple        uint 3   lengthValueType                                                                             ]
    [optional      uint 8   extTagNumber    'tagNumber == 15'                                                           ]
    [virtual       uint 8   actualTagNumber 'tagNumber < 15 ? tagNumber : extTagNumber'                                 ]
    [virtual       bit      isBoolean       'tagNumber == 1 && tagClass == TagClass.APPLICATION_TAGS'                   ]
    [virtual       bit      isConstructed   'tagClass == TagClass.CONTEXT_SPECIFIC_TAGS && lengthValueType == 6'        ]
    [virtual       bit      isPrimitiveAndNotBoolean '!isConstructed && !isBoolean'                                     ]
    [optional      uint 8   extLength       'isPrimitiveAndNotBoolean && lengthValueType == 5'                          ]
    [optional      uint 16  extExtLength    'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 254'      ]
    [optional      uint 32  extExtExtLength 'isPrimitiveAndNotBoolean && lengthValueType == 5 && extLength == 255'      ]
    [virtual       uint 32  actualLength    'lengthValueType == 5 && extLength == 255 ? extExtExtLength : (lengthValueType == 5 && extLength == 254 ? extExtLength : (lengthValueType == 5 ? extLength : lengthValueType))']
]

[discriminatedType BACnetApplicationTag
    [simple        BACnetTagHeader
                            header
    ]
    [validation    'header.tagClass == TagClass.APPLICATION_TAGS'    "should be a application tag"                      ]
    [virtual       uint 8   actualTagNumber 'header.actualTagNumber'                                                    ]
    [virtual       uint 32  actualLength    'header.actualLength'                                                       ]
    [typeSwitch actualTagNumber
        ['0x0' BACnetApplicationTagNull
        ]
        ['0x1' BACnetApplicationTagBoolean(BACnetTagHeader header)
            [simple BACnetTagPayloadBoolean('header.actualLength')
                                payload                                                                                 ]
            [virtual    bit     actualValue 'payload.value'                                                             ]
        ]
        ['0x2' BACnetApplicationTagUnsignedInteger(BACnetTagHeader header)
            [simple BACnetTagPayloadUnsignedInteger('header.actualLength')
                                payload                                                                                 ]
            [virtual    uint 64 actualValue   'payload.actualValue'                                                     ]
        ]
        ['0x3' BACnetApplicationTagSignedInteger(BACnetTagHeader header)
            [simple BACnetTagPayloadSignedInteger('header.actualLength')
                                payload                                                                                 ]
            [virtual    uint 64    actualValue   'payload.actualValue'                                                  ]
        ]
        ['0x4' BACnetApplicationTagReal
            [simple BACnetTagPayloadReal
                                payload                                                                                 ]

            [virtual    float 32     actualValue 'payload.value'                                                        ]
        ]
        ['0x5' BACnetApplicationTagDouble
            [simple BACnetTagPayloadDouble
                                payload                                                                                 ]
            [virtual    float 64     actualValue 'payload.value'                                                        ]
        ]
        ['0x6' BACnetApplicationTagOctetString(BACnetTagHeader header)
            [simple BACnetTagPayloadOctetString('header.actualLength')
                                payload                                                                                 ]
        ]
        ['0x7' BACnetApplicationTagCharacterString(BACnetTagHeader header)
            [simple BACnetTagPayloadCharacterString('header.actualLength')
                                payload                                                                                 ]
            [virtual vstring     value             'payload.value'                                                      ]
        ]
        ['0x8' BACnetApplicationTagBitString(BACnetTagHeader header)
            [simple BACnetTagPayloadBitString('header.actualLength')
                                payload                                                                                 ]
        ]
        ['0x9' BACnetApplicationTagEnumerated(BACnetTagHeader header)
            [simple BACnetTagPayloadEnumerated('header.actualLength')
                                payload                                                                                 ]
            [virtual  uint 32   actualValue 'payload.actualValue'                                                       ]
        ]
        ['0xA' BACnetApplicationTagDate
            [simple BACnetTagPayloadDate
                                payload                                                                                 ]
        ]
        ['0xB' BACnetApplicationTagTime
            [simple BACnetTagPayloadTime
                                payload                                                                                 ]
        ]
        ['0xC' BACnetApplicationTagObjectIdentifier
            [simple BACnetTagPayloadObjectIdentifier
                                payload                                                                                 ]
            [virtual    BACnetObjectType
                                objectType
                                               'payload.objectType'                                                     ]
            [virtual  uint 22   instanceNumber
                                               'payload.instanceNumber'                                                 ]
        ]
    ]
]

[discriminatedType BACnetContextTag(uint 8 tagNumberArgument, BACnetDataType dataType)
    [simple        BACnetTagHeader
                            header                                                                                      ]
    [validation    'header.actualTagNumber == tagNumberArgument' "tagnumber doesn't match" shouldFail=false             ]
    [validation    'header.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS' "should be a context tag"                        ]
    [virtual       uint 4   tagNumber     'header.tagNumber'                                                            ]
    [virtual       uint 32  actualLength  'header.actualLength'                                                         ]
    [validation    'header.lengthValueType != 6 && header.lengthValueType != 7'
                   "length 6 and 7 reserved for opening and closing tag" shouldFail=false                               ]
    [typeSwitch dataType
        ['NULL' BACnetContextTagNull(BACnetTagHeader header)
            [validation 'header.actualLength == 0' "length field should be 0"                                           ]
        ]
        ['BOOLEAN' BACnetContextTagBoolean(BACnetTagHeader header)
            [validation 'header.actualLength == 1' "length field should be 1"                                           ]
            [simple  uint 8 value                                                                                       ]
            [simple BACnetTagPayloadBoolean('value')
                            payload                                                                                     ]
            [virtual bit    actualValue 'payload.value'                                                                 ]
        ]
        ['UNSIGNED_INTEGER' BACnetContextTagUnsignedInteger(BACnetTagHeader header)
            [simple BACnetTagPayloadUnsignedInteger('header.actualLength')
                                payload                                                                                 ]
            [virtual    uint 64 actualValue 'payload.actualValue'                                                       ]
        ]
        ['SIGNED_INTEGER' BACnetContextTagSignedInteger(BACnetTagHeader header)
            [simple BACnetTagPayloadSignedInteger('header.actualLength')
                                payload                                                                                 ]
            [virtual    uint 64     actualValue 'payload.actualValue'                                                   ]
        ]
        ['REAL' BACnetContextTagReal
            [simple BACnetTagPayloadReal
                                    payload                                                                             ]
            [virtual    float 32     actualValue 'payload.value'                                                        ]
        ]
        ['DOUBLE' BACnetContextTagDouble
            [simple BACnetTagPayloadDouble
                                payload                                                                                 ]

            [virtual    float 64     actualValue 'payload.value'                                                        ]
        ]
        ['OCTET_STRING' BACnetContextTagOctetString(BACnetTagHeader header)
            [simple BACnetTagPayloadOctetString('header.actualLength')
                                payload                                                                                 ]
        ]
        ['CHARACTER_STRING' BACnetContextTagCharacterString(BACnetTagHeader header)
            [simple BACnetTagPayloadCharacterString('header.actualLength')
                                payload                                                                                 ]
            [virtual vstring     value             'payload.value'                                                      ]
        ]
        ['BIT_STRING' BACnetContextTagBitString(BACnetTagHeader header)
            [simple BACnetTagPayloadBitString('header.actualLength')
                                payload                                                                                 ]
        ]
        ['ENUMERATED' BACnetContextTagEnumerated(BACnetTagHeader header)
            [simple BACnetTagPayloadEnumerated('header.actualLength')
                                payload                                                                                 ]
            [virtual  uint 32   actualValue 'payload.actualValue'                                                       ]
        ]
        ['DATE' BACnetContextTagDate
            [simple BACnetTagPayloadDate
                                payload                                                                                 ]
        ]
        ['TIME' BACnetContextTagTime
            [simple     BACnetTagPayloadTime
                                payload                                                                                 ]
        ]
        ['BACNET_OBJECT_IDENTIFIER' BACnetContextTagObjectIdentifier
            [simple  BACnetTagPayloadObjectIdentifier
                                payload                                                                                 ]
            [virtual BACnetObjectType
                                objectType 'payload.objectType'                                                         ]
            [virtual uint 22    instanceNumber
                                               'payload.instanceNumber'                                                 ]
        ]
        ['UNKNOWN' BACnetContextTagUnknown(uint 32 actualLength)
            [array byte unknownData length 'actualLength'                                                               ]
        ]
    ]
]

[type BACnetOpeningTag(uint 8 tagNumberArgument)
    [simple        BACnetTagHeader header                                                                               ]
    [validation    'header.actualTagNumber == tagNumberArgument' "tagnumber doesn't match" shouldFail=false             ]
    [validation    'header.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS' "should be a context tag"                        ]
    [validation    'header.lengthValueType == 6' "opening tag should have a value of 6"                                 ]
]

[type BACnetClosingTag(uint 8 tagNumberArgument)
    [simple        BACnetTagHeader header                                                                               ]
    [validation    'header.actualTagNumber == tagNumberArgument' "tagnumber doesn't match" shouldFail=false             ]
    [validation    'header.tagClass == TagClass.CONTEXT_SPECIFIC_TAGS' "should be a context tag"                        ]
    [validation    'header.lengthValueType == 7' "closing tag should have a value of 7"                                 ]
]

[type BACnetTagPayloadBoolean(uint 32 actualLength)
    [virtual bit value   'actualLength == 1'    ]
    [virtual bit isTrue  'value'                ]
    [virtual bit isFalse '!value'               ]
]

[type BACnetTagPayloadUnsignedInteger(uint 32 actualLength)
    [virtual    bit         isUint8         'actualLength == 1'  ]
    [optional   uint  8     valueUint8      'isUint8'            ]
    [virtual    bit         isUint16        'actualLength == 2'  ]
    [optional   uint 16     valueUint16     'isUint16'           ]
    [virtual    bit         isUint24        'actualLength == 3'  ]
    [optional   uint 24     valueUint24     'isUint24'           ]
    [virtual    bit         isUint32        'actualLength == 4'  ]
    [optional   uint 32     valueUint32     'isUint32'           ]
    [virtual    bit         isUint40        'actualLength == 5'  ]
    [optional   uint 40     valueUint40     'isUint40'           ]
    [virtual    bit         isUint48        'actualLength == 6'  ]
    [optional   uint 48     valueUint48     'isUint48'           ]
    [virtual    bit         isUint56        'actualLength == 7'  ]
    [optional   uint 56     valueUint56     'isUint56'           ]
    [virtual    bit         isUint64        'actualLength == 8'  ]
    [optional   uint 64     valueUint64     'isUint64'           ]
    [validation 'isUint8 || isUint16 || isUint24 || isUint32 || isUint40 || isUint48 || isUint56 || isUint64' "unmapped integer length"]
    [virtual    uint 64     actualValue     'isUint8?valueUint8:(isUint16?valueUint16:(isUint24?valueUint24:(isUint32?valueUint32:(isUint40?valueUint40:(isUint48?valueUint48:(isUint56?valueUint56:valueUint64))))))']
]

[type BACnetTagPayloadSignedInteger(uint 32 actualLength)
    [virtual    bit         isInt8          'actualLength == 1'  ]
    [optional   int 8       valueInt8       'isInt8'             ]
    [virtual    bit         isInt16         'actualLength == 2'  ]
    [optional   int 16      valueInt16      'isInt16'            ]
    [virtual    bit         isInt24         'actualLength == 3'  ]
    [optional   int 24      valueInt24      'isInt24'            ]
    [virtual    bit         isInt32         'actualLength == 4'  ]
    [optional   int 32      valueInt32      'isInt32'            ]
    [virtual    bit         isInt40         'actualLength == 5'  ]
    [optional   int 40      valueInt40      'isInt40'            ]
    [virtual    bit         isInt48         'actualLength == 6'  ]
    [optional   int 48      valueInt48      'isInt48'            ]
    [virtual    bit         isInt56         'actualLength == 7'  ]
    [optional   int 56      valueInt56      'isInt56'            ]
    [virtual    bit         isInt64         'actualLength == 8'  ]
    [optional   int 64      valueInt64      'isInt64'            ]
    [validation 'isInt8 || isInt16 || isInt24 || isInt32 || isInt40 || isInt48 || isInt56 || isInt64' "unmapped integer length"]
    [virtual    uint 64     actualValue     'isInt8?valueInt8:(isInt16?valueInt16:(isInt24?valueInt24:(isInt32?valueInt32:(isInt40?valueInt40:(isInt48?valueInt48:(isInt56?valueInt56:valueInt64))))))']
]

[type BACnetTagPayloadReal
    [simple float 32 value]
]

[type BACnetTagPayloadDouble
    [simple float 64 value]
]

[type BACnetTagPayloadOctetString(uint 32 actualLength)
    [array   byte    octets  length 'actualLength'              ]
]

[type BACnetTagPayloadCharacterString(uint 32 actualLength)
    [simple     BACnetCharacterEncoding      encoding]
    // TODO: The reader expects int but uint32 gets mapped to long so even uint32 would easily overflow...
    [virtual    uint     16                  actualLengthInBit 'actualLength * 8 - 8']
    // TODO: call to string on encoding or add type conversion so we can use the enum above
    [simple     vstring 'actualLengthInBit'  value encoding='"UTF-8"']
]

[type BACnetTagPayloadBitString(uint 32 actualLength)
    [simple     uint 8      unusedBits                                           ]
    [array      bit         data count '((actualLength - 1) * 8) - unusedBits'   ]
    [array      bit         unused count 'unusedBits'                            ]
]

[type BACnetTagPayloadEnumerated(uint 32 actualLength)
    [array   byte       data length 'actualLength']
    [virtual uint 32    actualValue 'STATIC_CALL("parseVarUint", data)'  ]
]

[type BACnetTagPayloadDate
    [virtual uint  8 wildcard '0xFF'                                    ]
    [simple  uint  8 yearMinus1900                                      ]
    [virtual bit    yearIsWildcard 'yearMinus1900 == wildcard'          ]
    [virtual uint 16 year 'yearMinus1900 + 1900'                        ]
    [simple  uint  8 month                                              ]
    [virtual bit    monthIsWildcard 'month == wildcard'                 ]
    [virtual bit    oddMonthWildcard 'month == 13'                      ]
    [virtual bit    evenMonthWildcard 'month == 14'                     ]
    [simple  uint  8 dayOfMonth                                         ]
    [virtual bit    dayOfMonthIsWildcard 'dayOfMonth == wildcard'       ]
    [virtual bit    lastDayOfMonthWildcard 'dayOfMonth == 32'           ]
    [virtual bit    oddDayOfMonthWildcard 'dayOfMonth == 33'            ]
    [virtual bit    evenDayOfMonthWildcard 'dayOfMonth == 34'           ]
    [simple  uint  8 dayOfWeek                                          ]
    [virtual bit    dayOfWeekIsWildcard 'dayOfWeek == wildcard'         ]
]

[type BACnetTagPayloadTime
    [virtual uint  8 wildcard '0xFF'                                    ]
    [simple  uint  8 hour                                               ]
    [virtual bit    hourIsWildcard 'hour == wildcard'                   ]
    [simple  uint  8 minute                                             ]
    [virtual bit    minuteIsWildcard 'minute == wildcard'               ]
    [simple  uint  8 second                                             ]
    [virtual bit    secondIsWildcard 'second == wildcard'               ]
    [simple  uint  8 fractional                                         ]
    [virtual bit    fractionalIsWildcard 'fractional == wildcard'       ]
]

[type BACnetTagPayloadObjectIdentifier
    [manual     BACnetObjectType    objectType         'STATIC_CALL("readObjectType", readBuffer)' 'STATIC_CALL("writeObjectType", writeBuffer, objectType)' '10']
    [manual     uint 10             proprietaryValue   'STATIC_CALL("readProprietaryObjectType", readBuffer, objectType)' 'STATIC_CALL("writeProprietaryObjectType", writeBuffer, objectType, proprietaryValue)' '0']
    [virtual    bit                 isProprietary      'objectType == BACnetObjectType.VENDOR_PROPRIETARY_VALUE']
    [simple     uint 22             instanceNumber                      ]
]

// plc4x helper enum
[enum uint 1 TagClass
    ['0x0' APPLICATION_TAGS                     ]
    ['0x1' CONTEXT_SPECIFIC_TAGS                ]
]

// plc4x helper enum
[enum uint 8 BACnetDataType
    ['0' NULL                                   ]
    ['1' BOOLEAN                                ]
    ['2' UNSIGNED_INTEGER                       ]
    ['3' SIGNED_INTEGER                         ]
    ['4' REAL                                   ]
    ['5' DOUBLE                                 ]
    ['6' OCTET_STRING                           ]
    ['7' CHARACTER_STRING                       ]
    ['8' BIT_STRING                             ]
    ['9' ENUMERATED                             ]
    ['10' DATE                                  ]
    ['11' TIME                                  ]
    ['12' BACNET_OBJECT_IDENTIFIER              ]
    ['33' UNKNOWN                               ]
]

// plc4x helper enum
[enum byte BACnetCharacterEncoding
    ['0x0' ISO_10646                            ] // UTF-8
    ['0x1' IBM_Microsoft_DBCS                   ]
    ['0x2' JIS_X_0208                           ]
    ['0x3' ISO_10646_4                          ] // (UCS-4)
    ['0x4' ISO_10646_2                          ] // (UCS-2)
    ['0x5' ISO_8859_1                           ]
]
