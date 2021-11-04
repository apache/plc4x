/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

[enum 'GenericCANDataType' (uint 8 'numBits')
    [BYTE        [ '8'] ]
    [BOOLEAN     [ '1'] ]
    [UNSIGNED8   [ '8'] ]
    [UNSIGNED16  ['16'] ]
    [UNSIGNED24  ['24'] ]
    [UNSIGNED32  ['32'] ]
    [UNSIGNED40  ['40'] ]
    [UNSIGNED48  ['48'] ]
    [UNSIGNED56  ['56'] ]
    [UNSIGNED64  ['64'] ]
    [INTEGER8    [ '8'] ]
    [INTEGER16   ['16'] ]
    [INTEGER24   ['24'] ]
    [INTEGER32   ['32'] ]
    [INTEGER40   ['40'] ]
    [INTEGER48   ['48'] ]
    [INTEGER56   ['56'] ]
    [INTEGER64   ['64'] ]
    [REAL32      ['32'] ]
    [REAL64      ['64'] ]
]

[dataIo DataItem byteOrder='"LITTLE_ENDIAN"' (GenericCANDataType 'dataType')
    [typeSwitch 'dataType'
        ['BYTE' BYTE
            [simple byte 'value']
        ]
        ['BOOLEAN' BOOL
            [simple bit 'value']
        ]
        ['UNSIGNED8' USINT
            [simple uint 8 'value']
        ]
        ['UNSIGNED16' UINT
            [simple uint 16 'value']
        ]
        ['UNSIGNED24' UDINT
            [simple uint 24 'value']
        ]
        ['UNSIGNED32' UDINT
            [simple uint 32 'value']
        ]
        ['UNSIGNED40' ULINT
            [simple uint 40 'value']
        ]
        ['UNSIGNED48' ULINT
            [simple uint 48 'value']
        ]
        ['UNSIGNED56' ULINT
            [simple uint 56 'value']
        ]
        ['UNSIGNED64' ULINT
            [simple uint 64 'value']
        ]
        ['INTEGER8' SINT
            [simple int 8 'value']
        ]
        ['INTEGER16' INT
            [simple int 16 'value']
        ]
        ['INTEGER24' DINT
            [simple int 24 'value']
        ]
        ['INTEGER32' DINT
            [simple int 32 'value']
        ]
        ['INTEGER40' LINT
            [simple int 40 'value']
        ]
        ['INTEGER48' LINT
            [simple int 48 'value']
        ]
        ['INTEGER56' LINT
            [simple int 56 'value']
        ]
        ['INTEGER64' LINT
            [simple int 64 'value']
        ]
        ['REAL32' REAL
            [simple float 32 'value']
        ]
        ['REAL64' LREAL
            [simple float 64 'value']
        ]
    ]
]
