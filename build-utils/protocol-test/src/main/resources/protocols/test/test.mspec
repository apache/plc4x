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

[type 'SimpleType'
    [array         uint 8 'arrCount'      count      '5']
    [array         uint 8 'arrLen'        length     '6']
    [array         uint 8 'arrTerm'       terminated                'terminationExpression']
    [checksum      uint 8 'crc'           'checksumExpression']
    [const         uint 8 'con'           '0x03']
    [implicit      uint 8 'impl'          'serializationExpression']
    [manualArray   uint 8 'manArrayCount' count                     'loopExpression'            'serializationExpression' 'deserializationExpression' 'lengthExpression']
    [manualArray   uint 8 'manArrayLen'   length                    'loopExpression'            'serializationExpression' 'deserializationExpression' 'lengthExpression']
    [manualArray   uint 8 'manArrayTerm'  terminated                'loopExpression'            'serializationExpression' 'deserializationExpression' 'lengthExpression']
    [manual        uint 8 'man'           'serializationExpression' 'deserializationExpression' 'lengthExpression']
    [optional      uint 8 'opt'           'optionalExpression']
    [padding       uint 8 'pad'           '0'                       'paddingExpression']
    [reserved      uint 8 '0x00']
    [simple        uint 8 'simp']
    [virtual       uint 8 'virt'          'valueExpression']
]

////////////////////////////////////////////////////////////////
// Discriminated Type
////////////////////////////////////////////////////////////////

[discriminatedType 'DiscriminatedType'
    [discriminator uint 8 'discr']
    [typeSwitch 'discr'
        ['0x01' DiscriminatedTypeA
            [simple        uint 8 'simpA']
        ]
        ['0x02' DiscriminatedTypeA
            [simple        uint 8 'simpB']
        ]
        ['0x03' DiscriminatedTypeA
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


