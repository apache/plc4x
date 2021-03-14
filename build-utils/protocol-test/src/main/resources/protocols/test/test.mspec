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

////////////////////////////////////////////////////////////////
// Simple Type
////////////////////////////////////////////////////////////////

[type 'FieldTypeTest'
    [simple         uint 8 'simpleField']
    [abstract       unit 8  'abstractField']
    [array          uint 8  'arrayField'        count      '5']
    //Checksums fields are not supported in C
    //[checksum       uint 8  'checksumField'     '100']
    [const          uint 8  'constField'        '5']

    //Discriminated Field can't be used in simple type
    //[discriminator  uint 8  'discriminatorField']

    [enum           EnumType  'enumField']
    [implicit       uint 8  'implicitField' 'simpleField']
    [optional       uint 8  'optionalField' 'simpleField == 5']
    [padding        uint 8  'paddingField'  '0x00'  'simpleField']
    [reserved       uint 8  '0x00']

    //TypeSwitch field can't be used in non disriminatedTypes
    //[typeSwitch 'simpleField' ]
]

//The following data types are not yet implemented but have a reference
//Not Implemented Yet In
//-Java
//[simple ufloat 8.23 'ufloatField']

//Not Implemented Yet In
//-Java
//[simple ufloat 11.52 'udoubleField']

//Not Implemented Yet In
//-Java
//[simple time 8 'timeField']

//Not Implemented Yet In
//-Java
//[simple date 8 'dateField']

//Not Implemented Yet In
//-Java
//[simple dateTime 8 'dateTimeField']

[type 'SimpleTypeTest'
    [simple bit 'bitField']
    [simple int 8 'intField']
    [simple uint 8 'uintField']
    [simple float 8.23 'floatField']
    [simple float 11.52 'doubleField']
    [simple string '8' 'UTF-8' 'stringField']
]

[type 'AbstractTypeTest'
    [abstract bit 'bitField']
    [abstract int 8 'intField']
    [abstract uint 8 'uintField']
    [abstract float 8.23 'floatField']
    [abstract float 11.52 'doubleField']
    [abstract string '8' 'UTF-8' 'stringField']
]

[type 'ArrayTypeTest'
    [array bit 'bitField' count      '5']
    [array int 8 'intField' count      '5']
    [array uint 8 'uintField' count      '5']
    [array float 8.23 'floatField' count      '5']
    [array float 11.52 'doubleField' count      '5']
    [array string '8' 'UTF-8' 'stringField' count      '5']
]

//Checksums fields are not supported in C
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
    [const bit 'bitField' true]
    [const int 8 'intField' '100']
    [const uint 8 'uintField' '100']
    [const float 8.23 'floatField' '100.0f']
    [const float 11.52 'doubleField' '100.0']
    [const string '8' 'UTF-8' 'stringField' '"HELLO TODDY"']
]

[type 'EnumTypeTest'
    [enum           EnumType  'enumField']
]

[type 'ImplicitTypeTest'
    //Implicit types have the requirement that the expression is of a similar type to the field
    //i.e Integers can't be cast to Booleans
    [simple   uint 8 'simpleField']

    [implicit bit 'bitField' 'simpleField > 0']
    [implicit int 8 'intField' 'simpleField']
    [implicit uint 8 'uintField' 'simpleField']
    [implicit float 8.23 'floatField' 'simpleField']
    [implicit float 11.52 'doubleField' 'simpleField']
    //String literals can't be used in the expression
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

////////////////////////////////////////////////////////////////
// Discriminated Type
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



////////////////////////////////////////////////////////////////
// Arguments Type
////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////
// Discriminated Type with multiple conditions
////////////////////////////////////////////////////////////////


[enum uint 8 'EnumType'
    ['0x01' BOOL]
    ['0x02' UINT]
    ['0x03' INT]
]


[dataIo 'DataIOTypeEmpty'

]

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