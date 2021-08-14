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

////////////////////////////////////////////////////////////////
// Simple Type
////////////////////////////////////////////////////////////////

//Showing below an example comment including which languages it is supported in and ascii doc reference.
// Bit Field Test Single
// Java +, C +, Go +
// tag::SimpleBitTypeTest[]
[type 'SimpleBitTypeTest'
    [simple bit 'bitField']
]
// end::SimpleBitTypeTest[]

[type 'FieldTypeTest'
    [simple         uint 8 'simpleField']
    //Abstract fields can only be used within discriminated base types.
    //[abstract       unit 8  'abstractField']
    [array          uint 8  'arrayField'        count      '5']
    //TODO: Checksums fields are not supported in C
    //[checksum       uint 8  'checksumField'     '100']
    [const          uint 8  'constField'        '5']
    // Discriminated Field can't be used in simple type
    //[discriminator  uint 8  'discriminatorField']
    [enum           EnumType  'enumField']
    [implicit       uint 8  'implicitField' 'simpleField']
    [optional       uint 8  'optionalField' 'simpleField == 5']
    [padding        uint 8  'paddingField'  '0x00'  'simpleField']
    [reserved       uint 8  '0x00']
    // TypeSwitch field can't be used in non discriminatedTypes
    //[typeSwitch 'simpleField' ]
]

// If a type has an unknown field, the entire serializer is reduced to firing an exception
[type 'FieldTypeTestWithUnknownField'
    [simple         uint 8 'simpleField']
    //Abstract fields can only be used within discriminated base types.
    //[abstract       unit 8  'abstractField']
    [array          uint 8  'arrayField'        count      '5']
    //TODO: Checksums fields are not supported in C
    //[checksum       uint 8  'checksumField'     '100']
    [const          uint 8  'constField'        '5']
    // Discriminated Field can't be used in simple type
    //[discriminator  uint 8  'discriminatorField']
    [enum           EnumType  'enumField']
    [implicit       uint 8  'implicitField' 'simpleField']
    [optional       uint 8  'optionalField' 'simpleField == 5']
    [padding        uint 8  'paddingField'  '0x00'  'simpleField']
    [reserved       uint 8  '0x00']
    [unknown        uint 16]
    // TypeSwitch field can't be used in non discriminatedTypes
    //[typeSwitch 'simpleField' ]
]

/*
 * TODO: doesn't compile in java
[type 'UFloatTypeTest'
    [simple ufloat 8.23 'ufloatField']
    [simple ufloat 11.52 'udoubleField']
]
*/

/*
 * TODO: doesn't compile in java
[type 'TimeTypeTest'
    [simple time 8 'timeField']
    [simple date 8 'dateField']
    [simple dateTime 8 'dateTimeField']
]
*/

[type 'SimpleTypeTest'
    [simple bit 'bitField']
    [simple byte 'byteField']
    [simple int 8 'intField']
    [simple uint 8 'uintField']
    [simple float 8.23 'floatField']
    [simple float 11.52 'doubleField']
    [simple string '8' 'UTF-8' 'stringField']
]

[type 'AbstractTypeTest'
    //Abstract fields can only be used within discriminated base types.
    [simple         uint 8 'simpleField']
    [abstract bit 'abstractBitField']
    [abstract int 8 'abstractIntField']
    [abstract uint 8 'abstractUintField']
    [abstract float 8.23 'abstractFloatField']
    [abstract float 11.52 'abstractDoubleField']
    [abstract string '8' 'UTF-8' 'abstractStringField']
    [typeSwitch 'simpleField'
        ['0' AbstractedType
            [simple bit 'abstractBitField']
            [simple int 8 'abstractIntField']
            [simple uint 8 'abstractUintField']
            [simple float 8.23 'abstractFloatField']
            [simple float 11.52 'abstractDoubleField']
            [simple string '8' 'UTF-8' 'abstractStringField']
        ]
    ]
]

[type 'AbstractTypeTest'
    //Abstract fields can only be used within discriminated base types.
    [simple   uint 8 'simpleField']
    [abstract bit 'abstractBitField']
    [abstract int 8 'abstractIntField']
    [abstract uint 8 'abstractUintField']
    [abstract float 8.23 'abstractFloatField']
    [abstract float 11.52 'abstractDoubleField']
    [abstract string '8' 'UTF-8' 'abstractStringField']
    [typeSwitch 'simpleField'
        ['0' AbstractedType
            //Abstract fields need to be overriden in child
            [simple bit 'abstractBitField']
            [simple int 8 'abstractIntField']
            [simple uint 8 'abstractUintField']
            [simple float 8.23 'abstractFloatField']
            [simple float 11.52 'abstractDoubleField']
            [simple string '8' 'UTF-8' 'abstractStringField']
        ]
    ]
]

[type 'ArrayTypeTest'
    [array bit 'bitField' count      '5']
    [array int 8 'intField' count      '5']
    [array uint 8 'uintField' count      '5']
    [array float 8.23 'floatField' count      '5']
    [array float 11.52 'doubleField' count      '5']
    [array string '8' 'UTF-8' 'stringField' count      '5']
]

//TODO: Checksums fields are not supported in C
//[type 'CheckSumTypeTest'
    //Bit field cannot be used for a checksum
    //[checksum bit 'bitField' true]
    //[checksum int 8 'intField' '100']
    //[checksum uint 8 'uintField' '100']
    //Float fields cannot be used as checksums
    //[checksum float 8.23 'floatField' '100.0f']
    //[checksum float 11.52 'doubleField' '100.0']
    //String field cannot be used as a checksum
    //[checksum string '11 * 8' 'UTF-8' 'stringField' '"HELLO TODDY"']
//]

[type 'ConstTypeTest'
    [const bit 'bitField' 'true']
    [const int 8 'intField' '100']
    [const uint 8 'uintField' '100']
    [const float 8.23 'floatField' '100.0']
    [const float 11.52 'doubleField' '100.0']
    [const string '8' 'UTF-8' 'stringField' '"HELLO TODDY"']
]

[type 'EnumTypeTest'
    [enum           EnumType  'enumField']
]

[type 'PascalStringTypeTest'
    [simple int 8 'stringLength']
    [simple string 'stringLength' 'UTF-8' 'stringField']
]

[type 'ImplicitPascalStringTypeTest'
    [implicit int 8 'stringLength' 'stringField.length']
    [simple string 'stringLength' 'UTF-8' 'stringField']
]

[type 'ImplicitTypeTest'
    //Implicit types have the requirement that the expression is of a similar type to the field
    //TODO: i.e Integers can't be cast to Booleans
    [simple   uint 8 'simpleField']

    [implicit bit 'bitField' 'simpleField > 0']
    [implicit int 8 'intField' 'simpleField']
    [implicit uint 8 'uintField' 'simpleField']
    [implicit float 8.23 'floatField' 'simpleField']
    [implicit float 11.52 'doubleField' 'simpleField']
    //TODO: String literals can't be used in the expression
    //[implicit string '8' 'UTF-8' 'stringField' 'simpleField > 0 ? "HELLO TODDY" : "BYE TODDY"']
]

[type 'OptionalTypeTest'
    [simple         uint 8 'simpleField']
    [optional       uint 8  'optionalField' 'simpleField == 5']
]

[type 'PaddingTypeTest'
    [simple         uint 8 'simpleField']
    [padding        uint 8  'paddingField'  '0x00'  'simpleField']
]

[type 'ReservedTypeTest'
    [reserved       uint 8  '0x00']
]

//TODO: Virtual fields fail for GO, haven't checked C assuming fails.
//[type 'VirtualFieldTest'
//    [simple  uint 8 'simpleField']
//    [virtual bit 'virtualBitField' 'simpleField == 0']
//    [virtual int 8 'virtualIntField' 'simpleField']
//    [virtual uint 8 'virtualUintField' 'simpleField']
//    [virtual float 8.23 'virtualFloatField' 'simpleField']
//    [virtual float 11.52 'virtualDoubleField' 'simpleField']
//    [virtual string '24' 'virtualStringField' 'simpleField']
//]

//TODO: Virtual fields fail for GO, haven't checked C assuming fails.
//[discriminatedType 'DiscriminatedVirtualTypeTest'
//    [simple  uint 8 'simpleField']
//    [virtual bit 'virtualBitField' 'simpleField == 0']
//    [virtual int 8 'virtualIntField' 'simpleField']
//    [virtual uint 8 'virtualUintField' 'simpleField']
//    [virtual float 8.23 'virtualFloatField' 'simpleField']
//    [virtual float 11.52 'virtualDoubleField' 'simpleField']
//    [virtual string '24' 'UTF-8' 'virtualStringField' 'simpleField']
//    [typeSwitch 'simpleField'
//        ['0' DiscriminatedVirtualType
//            [simple int 8 'intField']
//        ]
//    ]
//]

[type 'IntTypeTest'
    [simple int 3 'ThreeField']
    [simple int 8 'ByteField']
    [simple int 16 'WordField']
    [simple int 24 'WordPlusByteField']
    [simple int 32 'DoubleIntField']
    [simple int 64 'QuadIntField']
]

[type 'UIntTypeTest'
    [simple uint 3 'ThreeField']
    [simple uint 8 'ByteField']
    [simple uint 16 'WordField']
    [simple uint 24 'WordPlusByteField']
    [simple uint 32 'DoubleIntField']
    [simple uint 64 'QuadIntField']
]

//Specific test confirming a continous loop isn't formed when working out the length.
[type 'LentghLoopTest'
    [simple        uint 16 'commandType']
    [implicit      uint 16 'len' 'lengthInBytes - 8']
]

////////////////////////////////////////////////////////////////
// Discriminated Type Tests
////////////////////////////////////////////////////////////////

[discriminatedType 'EnumDiscriminatedType'
    [discriminator EnumType 'discr']
    [typeSwitch 'discr'
        ['BOOL' EnumDiscriminatedTypeA
            [simple        uint 8 'simpA']
        ]
        ['UINT' EnumDiscriminatedTypeB
            [simple        uint 8 'simpB']
        ]
        ['INT' EnumDiscriminatedTypeC
            [simple        uint 8 'simpC']
        ]
    ]
]

// Multiple Enumerated discriminators
[discriminatedType 'EnumDiscriminatedTypeMultiple'
    [discriminator EnumType 'discr1']
    [discriminator EnumTypeInt 'discr2']
    [typeSwitch 'discr1','discr2'
        ['BOOL','BOOLINT' EnumDiscriminatedTypeMultipleA
            [simple        uint 8 'simpA']
        ]
        ['UINT','UINTINT' EnumDiscriminatedTypeMultipleB
            [simple        uint 8 'simpB']
        ]
        ['INT','INTINT' EnumDiscriminatedTypeMultipleC
            [simple        uint 8 'simpC']
        ]
    ]
]

// Enumerated Parameter
[discriminatedType 'EnumDiscriminatedTypeParameter' [EnumType 'discr']
    [typeSwitch 'discr'
        ['BOOL' EnumDiscriminatedTypeAParameter
            [simple        uint 8 'simpA']
        ]
        ['UINT' EnumDiscriminatedTypeBParameter
            [simple        uint 8 'simpB']
        ]
        ['INT' EnumDiscriminatedTypeCParameter
            [simple        uint 8 'simpC']
        ]
    ]
]

// Multiple Enumerated Parameters
[discriminatedType 'EnumDiscriminatedTypeParameterMultiple' [EnumType 'discr1', EnumTypeInt 'discr2']
    [typeSwitch 'discr1','discr2'
        ['BOOL','BOOLINT' EnumDiscriminatedTypeAParameterMultiple
            [simple        uint 8 'simpA']
        ]
        ['UINT','UINTINT' EnumDiscriminatedTypeBParameterMultiple
            [simple        uint 8 'simpB']
        ]
        ['INT','INTINT' EnumDiscriminatedTypeCParameterMultiple
            [simple        uint 8 'simpC']
        ]
    ]
]

[discriminatedType 'SimpleDiscriminatedType'
    [discriminator uint 8 'discr']
    [typeSwitch 'discr'
        ['0x00' SimpleDiscriminatedTypeA
            [simple        uint 8 'simpA']
        ]
        ['0x01' SimpleDiscriminatedTypeB
            [simple        uint 8 'simpB']
        ]
        ['0x02' SimpleDiscriminatedTypeC
            [simple        uint 8 'simpC']
        ]
    ]
]


//Test to check if we can include concrete types as fields. Doesn't work in any language at the moment.
//[discriminatedType 'SimpleDiscriminatedType'
//    [discriminator uint 8 'discr']
//    [typeSwitch 'discr'
//        ['0x00' SimpleDiscriminatedTypeA
//            [simple        AnotherSimpleDiscriminatedTypeA 'simpA']
//        ]
//    ]
//]

//[discriminatedType 'AnotherSimpleDiscriminatedType'
//    [discriminator uint 8 'discr']
//    [typeSwitch 'discr'
//        ['0x00' AnotherSimpleDiscriminatedTypeA
//            [simple        uint 8 'simpA']
//        ]
//    ]
//]

////////////////////////////////////////////////////////////////
// Enumerated Type Tests
////////////////////////////////////////////////////////////////

[enum bit 'EnumTypeBit'
    ['true' TRUE]
    ['false' FALSE]
]

[enum int 8 'EnumTypeInt'
    ['0x01' BOOLINT]
    ['0x02' UINTINT]
    ['0x03' INTINT]
]

[enum uint 8 'EnumType'
    ['0x01' BOOL]
    ['0x02' UINT]
    ['0x03' INT]
]

//TODO:  C doesn't support non integer switch fields
//[enum float 8.23 'EnumTypeFloat'
//    ['100.0' LOW]
//    ['101.0' MID]
//    ['102.0' BIG]
//]

//TODO:  C doesn't support non integer switch fields
//[enum float 11.52 'EnumTypeDouble'
//    ['100.0' LOW]
//    ['101.0' MID]
//    ['102.0' BIG]
//]

//String based enum's needs some work in C, compiles but assigns 0 as values.
[enum string '-1' 'EnumTypeString'
    ['Toddy1' TODDY]
]

//TODO:  Fails to import the base Enum in C, need to find it in getComplexTypeReferences
//[enum EnumType 'EnumTypeEnum'
//    ['BOOL' BOOL]
//    ['UINT' UINT]
//    ['INT' INT]
//]

//TODO:  Float parameters aren't implemented for constants in enums in C
//[enum int 8 'EnumTypeAllTest'  [bit 'bitType', int 8 'intType', uint 8 'uintType', float 8.23 'floatType', float 11.52 'doubleType', string '-1' 'stringType', EnumType 'enumType']
//    ['0x01' BOOL             ['false'      , '1'               , '1'                 , '100.0'                  , '100.0'              , 'BOOL'         , 'BOOL']]
//    ['0x02' BYTE             ['true'       , '2'               , '2'                 , '101.1'                  , '101.1'              , 'BYTE'         , 'UINT']]
//]

//TODO:  Keyword named parameters aren't allowed
//[enum int 8 'EnumTypeIntTest'  [int 8 'int']
//    ['0x01' BOOL             ['1']]
//    ['0x02' BYTE             ['2']]
//]

//Showing allowed parameter types for enums
[enum int 8 'EnumTypeParameters'  [bit 'bitType', int 8 'intType', uint 8 'uintType', string '-1' 'stringType', EnumType 'enumType']
    ['0x01' BOOL             ['false'      , '1'               , '1'                 , 'BOOL'         , 'BOOL']]
    ['0x02' BYTE             ['true'       , '2'               , '2'                 , 'BYTE'         , 'UINT']]
]

////////////////////////////////////////////////////////////////
// Data IO Tests
////////////////////////////////////////////////////////////////

[dataIo 'DataIOType' [EnumType 'dataType']
    [typeSwitch 'dataType'
        ['BOOL' BOOL
            [simple bit 'value']
        ]
        ['UINT' USINT
            [simple uint 8 'value']
        ]
        ['INT' UINT
            [simple uint 16 'value']
        ]
    ]
]
