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
// Simple Type
////////////////////////////////////////////////////////////////

//Showing below an example comment including which languages it is supported in and ascii doc reference.
// Bit Field Test Single
// Java +, C +, Go +
// tag::SimpleBitTypeTest[]
[type SimpleBitTypeTest
    [simple bit bitField]
]
// end::SimpleBitTypeTest[]

[type FieldTypeTest
    [simple         uint 8 simpleField                         ]
    //Abstract fields can only be used within discriminated base types.
    //[abstract       unit 8  abstractField                      ]
    [array          uint 8  arrayFieldCount      count      '5'   ]
    [array          uint 8  arrayFieldLength     length     '5'   ]
    [array          uint 8  arrayFieldTerminated terminated 'true']
    [assert         uint 8  assertField          '42'             ]
    [checksum       uint 8  checksumField        'STATIC_CALL("crcUint8" , -23)']
    [const          uint 8  constField           5                ]
    // Discriminated Field can't be used in simple type
    //[discriminator  uint 8  discriminatorField]
    [enum           EnumTypeParameters  enumField intType      ]
    [implicit       uint 8  implicitField 'simpleField'        ]
    // TODO: Currently manual array seems to only support the flavor "terminated"
    //[manualArray    uint 8  manualArrayField count      '1'    ] // TODO: Implement ...
    //[manualArray    uint 8  manualArrayField length     '1'    ] // TODO: Implement ...
    //[manualArray    uint 8  manualArrayField terminated '1'    ] // TODO: Implement ...
    [manual         uint 8  manualField 'STATIC_CALL("readAManualField", readBuffer, simpleField)'
                                        'STATIC_CALL("writeAManualField", writeBuffer, simpleField)'
                                        'simpleField*8'                   ]
    [optional       uint 8  optionalField 'simpleField == 5'   ]
    [padding        uint 8  paddingField  '0x00'  'simpleField']
    [reserved       uint 8  '0x00'                             ]
    // TypeSwitch field can't be used in non discriminatedTypes
    //[typeSwitch simpleField                                    ]
    [unknown        uint 8                                     ]
    [virtual        uint 8  virtualField  'constField + 10'    ]
    [validation     'false' "Validation Message"               ]
    [peek           uint 8  peekField                          ]
]

// If a type has an unknown field, the entire serializer is reduced to firing an exception
[type FieldTypeTestWithUnknownField
    [simple         uint 8 simpleField]
    //Abstract fields can only be used within discriminated base types.
    //[abstract       unit 8  abstractField]
    [array          uint 8  arrayField        count      '5']
    [checksum       uint 8  checksumField     'STATIC_CALL("crcUint8" , -23)']
    [const          uint 8  constField        5]
    // Discriminated Field can't be used in simple type
    //[discriminator  uint 8  discriminatorField]
    [enum           EnumTypeParameters  enumField intType]
    [implicit       uint 8  implicitField 'simpleField']
    [optional       uint 8  optionalField 'simpleField == 5']
    [padding        uint 8  paddingField  '0x00'  'simpleField']
    [reserved       uint 8  '0x00']
    [unknown        uint 16]
    // TypeSwitch field can't be used in non discriminatedTypes
    //[typeSwitch simpleField ]
]

/*
 * TODO: doesn't compile in golang and c
[type TimeTypeTest
    [simple time 8 timeField]
    [simple date 8 dateField]
    [simple dateTime 8 dateTimeField]
]
*/

[type SimpleTypeTest
    [simple bit bitField]
    [simple byte byteField]
    [simple int 8 intField]
    [simple uint 8 uintField]
    [simple float 32 floatField]
    [simple float 64 doubleField]
    [simple string 8 stringField]
]

[type AbstractTypeTest
    //Abstract fields can only be used within discriminated base types.
    [simple   uint 8 simpleField]
    [abstract bit abstractBitField]
    [abstract int 8 abstractIntField]
    [abstract uint 8 abstractUintField]
    [abstract float 32 abstractFloatField]
    [abstract float 64 abstractDoubleField]
    [abstract string 8 abstractStringField]
    [typeSwitch simpleField
        ['0' *SubType
            //Abstract fields need to be overridden in child
            [simple bit abstractBitField]
            [simple int 8 abstractIntField]
            [simple uint 8 abstractUintField]
            [simple float 32 abstractFloatField]
            [simple float 64 abstractDoubleField]
            [simple string 8 abstractStringField]
        ]
    ]
    [abstract bit oneMoreBit] // TODO: apparently this breaks java
]

[type ArrayTypeTest
    [array bit bitField count              '5']
    [array byte byteField count            '5']
    [array int 8 intField count            '5']
    [array uint 8 uintField count          '5']
    [array float 32 floatField count       '5']
    [array float 64 doubleField count      '5']
    [array string 8 stringField count      '5']
]

/* TODO: we need to check the expressions. The loop expression for count is a number and not true. True fits to terminated. So it seems the java impl is of somehow
[type ManualArrayTypeTest
    [manualArray bit        bitField    count      'true'  'STATIC_CALL("parseBit",    readBuffer)' 'STATIC_CALL("serializeBit",    writeBuffer, _value)' '1'   ]
    [manualArray byte       byteField   count      'true'  'STATIC_CALL("parseByte",   readBuffer)' 'STATIC_CALL("serializeByte",   writeBuffer, _value)' '5'   ]
    [manualArray int     8  intField    count      'true'  'STATIC_CALL("parseInt8",   readBuffer)' 'STATIC_CALL("serializeInt8",   writeBuffer, _value)' '5'   ]
    [manualArray uint    8  uintField   count      'true'  'STATIC_CALL("parseUint8",  readBuffer)' 'STATIC_CALL("serializeUint8",  writeBuffer, _value)' '5'   ]
    [manualArray float  32  floatField  count      'true'  'STATIC_CALL("parseFloat",  readBuffer)' 'STATIC_CALL("serializeFloat",  writeBuffer, _value)' '20'  ]
    [manualArray float  64  doubleField count      'true'  'STATIC_CALL("parseDouble", readBuffer)' 'STATIC_CALL("serializeDouble", writeBuffer, _value)' '40'  ]
    [manualArray string  8  stringField count      'true'  'STATIC_CALL("parseString", readBuffer)' 'STATIC_CALL("serializeString", writeBuffer, _value)' '5'   ]
]*/

[type CheckSumTypeTest
    //Bit field cannot be used for a checksum: [checksum bit bitField 'true']
    [checksum int 8  intField  'STATIC_CALL("crcInt8" , -23)']
    [checksum uint 8 uintField 'STATIC_CALL("crcUint8", 42)']
    // Float fields cannot be used as checksums: [checksum float 32 floatField '100.0']
    // Float fields cannot be used as checksums: [checksum float 64 doubleField '100.0']
    // String field cannot be used as a checksum: [checksum vstring '11 * 8' stringField '"HELLO TODDY"']
]

[type ConstTypeTest
    [const bit bitField true]
    [const int 8 intField 100]
    [const uint 8 uintField 100]
    [const float 32 floatField 100.0]
    [const float 64 doubleField 100.0]
    [const string 8 stringField "HELLO TODDY"]
]

[type EnumTypeTest
    [enum           EnumTypeParameters  enumField intType]
]

[type PascalStringTypeTest
    [simple int 8 stringLength]
    [simple vstring 'stringLength' stringField]
]

[type ImplicitPascalStringTypeTest
    [implicit int 8 stringLength 'stringField.length']
    [simple vstring 'stringLength' stringField]
]

[type ImplicitTypeTest
    //Implicit types have the requirement that the expression is of a similar type to the field
    [simple   uint 8 simpleField]

    [implicit bit bitField 'simpleField > 0']
    [implicit int 8 intField 'simpleField']
    [implicit uint 8 uintField 'simpleField']
    [implicit float 32 floatField 'simpleField']
    [implicit float 64 doubleField 'simpleField']
    [implicit string 8 stringField 'simpleField > 0 ? "HELLO TODDY" : "BYE TODDY"']
]

[type OptionalTypeTest
    [simple         uint 8 simpleField]
    [optional       uint 8  optionalField 'simpleField == 5']
]

[type PaddingTypeTest
    [simple         uint 8 simpleField]
    [padding        uint 8  paddingField  '0x00'  'simpleField']
]

[type ReservedTypeTest
    [reserved       uint 8  '0x00']
]

[type ReservedTypeTestParent
    [reserved       uint 8  '0x00']
    [simple  uint 8 simpleField]
    [typeSwitch simpleField
        ['0' ReservedTypeTestChild
            [reserved       uint 8  '0x00']
        ]
    ]
]

// Check field with arguments
[type Struct(bit signed, bit unsigned)
  [typeSwitch signed
    ['true' SignedType
      [simple int 8 data]
    ]
    ['false' UnsignedType
      [simple uint 8 data]
    ]
  ]
]

[type StructContainer
    [simple bit signed]
    [simple Struct('signed', '!signed') struct]
]

// TODO: So far only trouble in GO, C seems OK.
[type VirtualFieldTest
    [simple  uint 8 simpleField]
    [virtual bit  virtualBitField 'simpleField == 0']
    [virtual int 8  virtualIntField 'simpleField']
    [virtual uint 8  virtualUintField 'simpleField']
    [virtual float 32  virtualFloatField 'simpleField']
    [virtual float 64  virtualDoubleField 'simpleField']
    [virtual string 24  virtualStringField 'simpleField']
]

//TODO: Virtual fields fail for GO, haven't checked C assuming fails.
//[discriminatedType DiscriminatedVirtualTypeTest
//    [simple  uint 8 simpleField]
//    [virtual bit  virtualBitField 'simpleField == 0']
//    [virtual int 8  virtualIntField 'simpleField']
//    [virtual uint 8  virtualUintField 'simpleField']
//    [virtual float 32  virtualFloatField 'simpleField']
//    [virtual float 64  virtualDoubleField 'simpleField']
//    [virtual string 24  virtualStringField 'simpleField']
//    [typeSwitch simpleField
//        ['0' DiscriminatedVirtualType
//            [simple int 8 intField]
//        ]
//    ]
//]

[type IntTypeTest
    [simple int 3 ThreeField]
    [simple int 8 ByteField]
    [simple int 16 WordField]
    [simple int 24 WordPlusByteField]
    [simple int 32 DoubleIntField]
    [simple int 64 QuadIntField]
]

[type UIntTypeTest
    [simple uint 3 ThreeField]
    [simple uint 8 ByteField]
    [simple uint 16 WordField]
    [simple uint 24 WordPlusByteField]
    [simple uint 32 DoubleIntField]
    [simple uint 64 QuadIntField]
]

//Specific test confirming a continuous loop isn't formed when working out the length.
[type LengthLoopTest
    [simple        uint 16 commandType]
    [implicit      uint 16 len 'lengthInBytes - 8']
]

// Complex of complex usage
[type ComplexOfComplexUsageA
    [simple         ComplexOfComplexUsageB      b]
]

[type ComplexOfComplexUsageB
    [optional       ComplexOfComplexUsageC      c]
    [array          ComplexOfComplexUsageD      ds count      '5']
]

[type ComplexOfComplexUsageC
    [simple uint 8  irrelevant]
]

[type ComplexOfComplexUsageD
    [simple uint 8  irrelevant]
]

////////////////////////////////////////////////////////////////
// Discriminated Type Tests
////////////////////////////////////////////////////////////////

[discriminatedType EnumDiscriminatedType
    [discriminator EnumType discr]
    [typeSwitch discr
        ['BOOL' *A
            [simple        uint 8 simpA]
        ]
        ['UINT' *B
            [simple        uint 8 simpB]
        ]
        ['INT' *C
            [simple        uint 8 simpC]
        ]
    ]
]

// Multiple Enumerated discriminators
[discriminatedType EnumDiscriminatedTypeMultiple
    [discriminator EnumType discr1]
    [discriminator EnumTypeInt discr2]
    [typeSwitch discr1,discr2
        ['BOOL','BOOLINT' *A
            [simple        uint 8 simpA]
        ]
        ['UINT','UINTINT' *B
            [simple        uint 8 simpB]
        ]
        ['INT','INTINT'   *C
            [simple        uint 8 simpC]
        ]
    ]
]

// Enumerated Parameter
[discriminatedType EnumDiscriminatedTypeParameter(EnumType discr)
    [typeSwitch discr
        ['BOOL' *A
            [simple        uint 8 simpA]
        ]
        ['UINT' *B
            [simple        uint 8 simpB]
        ]
        ['INT' *C
            [simple        uint 8 simpC]
        ]
    ]
]

// Multiple Enumerated Parameters
[discriminatedType EnumDiscriminatedTypeParameterMultiple(EnumType discr1, EnumTypeInt discr2)
    [typeSwitch discr1,discr2
        ['BOOL','BOOLINT' *A
            [simple        uint 8 simpA]
        ]
        ['UINT','UINTINT' *B
            [simple        uint 8 simpB]
        ]
        ['INT','INTINT'   *C
            [simple        uint 8 simpC]
        ]
    ]
]

[discriminatedType SimpleDiscriminatedType
    [discriminator uint 8 discr]
    [typeSwitch discr
        ['0x00' *A
            [simple        uint 8 simpA]
        ]
        ['0x01' *B
            [simple        uint 8 simpB]
        ]
        ['0x02' *C
            [simple        uint 8 simpC]
        ]
    ]
]


//Test to check if we can include concrete types as fields. Doesn't work in any language at the moment.
// TODO: Currently C doesn't support directly referencing child types
/*[discriminatedType SimpleDiscriminatedType
    [discriminator uint 8 discr]
    [typeSwitch discr
        ['0x00' *A
            [simple        AnotherSimpleDiscriminatedTypeA simpA]
        ]
    ]
]

[discriminatedType AnotherSimpleDiscriminatedType
    [discriminator uint 8 discr]
    [typeSwitch discr
        ['0x00' *A
            [simple        uint 8 simpA]
        ]
    ]
]*/

////////////////////////////////////////////////////////////////
// Enumerated Type Tests
////////////////////////////////////////////////////////////////

[enum bit EnumTypeBit
    ['true' TRUE]
    ['false' FALSE]
]

[enum int 8 EnumTypeInt
    ['0x01' BOOLINT]
    ['0x02' UINTINT]
    ['0x03' INTINT]
]

[enum uint 8 EnumType
    ['0x01' BOOL]
    ['0x02' UINT]
    ['0x03' INT]
    ['0x11' BYTE_ARRAY]
]

//TODO:  C doesn't support non integer switch fields (However this feature might not make sense)
//[enum float 32 EnumTypeFloat
//    ['100.0' LOW]
//    ['101.0' MID]
//    ['102.0' BIG]
//]

//TODO:  C doesn't support non integer switch fields (However this feature might not make sense)
//[enum float 64 EnumTypeDouble
//    ['100.0' LOW]
//    ['101.0' MID]
//    ['102.0' BIG]
//]

//TODO:  C doesn't support non integer switch fields
//[enum string 48 EnumTypeString
//    ['Toddy1' TODDY]
//]

//TODO:  Fails to import the base Enum in C, need to find it in getComplexTypeReferences (However this feature might not make sense)
//[enum EnumType EnumTypeEnum
//    ['BOOL' BOOL]
//    ['UINT' UINT]
//    ['INT' INT]
//]

//TODO:  Float parameters aren't implemented for constants in enums in C
//[enum int 8 EnumTypeAllTest(bit bitType, int 8 intType, uint 8 uintType, float 32 floatType, float 64 doubleType, string '-1' stringType, EnumType enumType)
//    ['0x01' BOOL             ['false'      , '1'               , '1'                 , '100.0'                  , '100.0'              , 'BOOL'         , 'BOOL']]
//    ['0x02' BYTE             ['true'       , '2'               , '2'                 , '101.1'                  , '101.1'              , 'BYTE'         , 'UINT']]
//]

//TODO:  Keyword named parameters aren't allowed (This won't be fixed)
//[enum int 8 EnumTypeIntTest(int 8 int)
//    ['0x01' BOOL             ['1']]
//    ['0x02' BYTE             ['2']]
//]

//Showing allowed parameter types for enums
[enum int 8 EnumTypeParameters(bit bitType, int 8 intType, uint 8 uintType, string 32 stringType, EnumType enumType)
    ['0x01' BOOL             ['false'      , '1'               , '1'                 , 'BOOL'         , 'BOOL']]
    ['0x02' BYTE             ['true'       , '2'               , '2'                 , 'BYTE'         , 'UINT']]
]

////////////////////////////////////////////////////////////////
// Data IO Tests
////////////////////////////////////////////////////////////////

[dataIo DataIOType(EnumType dataType)
    [typeSwitch dataType
        ['BOOL' BOOL
            [simple bit value]
        ]
        ['UINT' USINT
            [simple uint 8 value]
        ]
        ['INT' UINT
            [simple uint 16 value]
        ]
        // TODO: Causes problems in C<<<>>a <
        /*['BYTE_ARRAY' List
            [array byte value length '5']
        ]*/
    ]
]

////////////////////////////////////////////////////////////////
// TypeSwitch in TypeSwitch
////////////////////////////////////////////////////////////////

[discriminatedType TTGranddad
    [discriminator uint 8 dadNumber]
    [simple uint 8 warStories]
    [typeSwitch dadNumber
        ['0x00' TTDad
            [discriminator uint 8 sonNumber]
            [simple uint 8 beerBottles]
            [typeSwitch sonNumber
                ['0x01' TTSon
                    [simple uint 8 gameConsoles]
                    [discriminator uint 8 babyNumber]
                    [typeSwitch babyNumber
                        ['0x02' TTBaby
                            [simple uint 8 lalalala]
                        ]
                    ]
                ]
            ]
        ]
    ]
]

/* Needs to be ported to C and GO
[discriminatedType TypeSwitchInTypeSwitchParentType
    [discriminator uint 8 typeNumber]
    [simple uint 8 parentFieldHurz]
    [typeSwitch typeNumber
        ['0x00' *Child0
            [discriminator uint 8 childNumber]
            [simple uint 8 childFieldwolf]
            [typeSwitch childNumber
                ['0x01' *Infant0
                    [discriminator uint 8 infantNumber]
                    [simple uint 8 infantSomeField00]
                    [typeSwitch infantNumber
                        ['0x98' *InfantsChild0
                            [simple uint 8 infantsChild0SomeField00]
                        ]
                        ['0x99' *InfantsChild1
                            [simple uint 8 infantsChild1SomeField00]
                        ]
                    ]
                ]
                ['0x02' *Infant1
                    [simple uint 8 infantSomeField01]
                ]
                ['0x03' *Infant2
                    [simple uint 8 infantSomeField02]
                ]
            ]
        ]
        ['0x10' *Child1
            [discriminator uint 8 childNumber]
            [simple uint 8 childFieldLamm]
            [typeSwitch childNumber
                ['0x11' *Infant3
                    [simple uint 8 infantSomeField1ß]
                ]
                ['0x12' *Infant4
                    [simple uint 8 infantSomeField11]
                ]
                ['0x13' *Infant5
                    [simple uint 8 infantSomeField12]
                ]
            ]
        ]
        ['0x20' *Child2
            [discriminator uint 8 childNumber]
            [simple uint 8 childFieldGrueneWiese]
            [typeSwitch childNumber
                ['0x21' *Infant6
                    [simple uint 8 infantSomeField21]
                ]
                ['0x22' *Infant7
                    [simple uint 8 infantSomeField22]
                ]
                ['0x23' *Infant8
                    [simple uint 8 infantSomeField23]
                ]
            ]
        ]
    ]
]*/

[discriminatedType TypeSwitchWithArg(bit arg1, uint 8 arg2)
    [abstract bit isItTrue]
    [simple uint 8 value]
    [typeSwitch arg1
        ['true' TrustfulTypeSwitch(uint 8 value)
            [virtual bit isItTrue 'true']
        ]
        ['false' UnTrustfulTypeSwitch(uint 8 value)
            [virtual bit isItTrue 'false']
        ]
    ]
]

////////////////////////////////////////////////////////////////
// Missing Tests
////////////////////////////////////////////////////////////////

// TODO: Test case where a dataIo type is used as type of a simple field (Linking normal model and dataIo)