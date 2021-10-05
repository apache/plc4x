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

// Remark: The different fields are encoded in Big-endian.

[type   'Dummy'
    [simple uint        16  'dummy']
]

[dataIo 'DataItem' [vstring 'dataType', uint 16 'numberOfValues']
    [typeSwitch 'dataType','numberOfValues'
        ['BOOL','1' BOOL
            [simple   bit    'value']
        ]
        ['BOOL' List
            [array bit 'value' count 'numberOfValues']
        ]
        ['BYTE','1' BYTE
            [simple uint 8 'value']
        ]
        ['BYTE' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['WORD','1' WORD
            [simple uint 16 'value']
        ]
        ['WORD' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['DWORD','1' DWORD
            [simple uint 32 'value']
        ]
        ['DWORD' List
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['LWORD','1' LWORD
            [simple uint 64 'value']
        ]
        ['LWORD' List
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['SINT','1' SINT
            [simple int 8 'value']
        ]
        ['SINT' List
            [array int 8 'value' count 'numberOfValues']
        ]
        ['INT','1' INT
            [simple int 16 'value']
        ]
        ['INT' List
            [array int 16 'value' count 'numberOfValues']
        ]
        ['DINT','1' DINT
            [simple int 32 'value']
        ]
        ['DINT' List
            [array int 32 'value' count 'numberOfValues']
        ]
        ['LINT','1' LINT
            [simple int 64 'value']
        ]
        ['LINT' List
            [array int 64 'value' count 'numberOfValues']
        ]
        ['USINT','1' USINT
            [simple uint 8 'value']
        ]
        ['USINT' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['UINT','1' UINT
            [simple uint 16 'value']
        ]
        ['UINT' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['UDINT','1' UDINT
            [simple uint 32 'value']
        ]
        ['UDINT' List
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['ULINT','1' ULINT
            [simple uint 64 'value']
        ]
        ['ULINT' List
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['REAL','1' REAL
            [simple float 8.23  'value']
        ]
        ['REAL' List
            [array float 8.23 'value' count 'numberOfValues']
        ]
        ['LREAL','1' LREAL
            [simple float 11.52  'value']
        ]
        ['LREAL' List
            [array float 11.52 'value' count 'numberOfValues']
        ]
        ['CHAR','1' CHAR
            [simple uint 8 'value']
        ]
        ['CHAR' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['WCHAR','1' WCHAR
            [simple uint 16 'value']
        ]
        ['WCHAR' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['STRING' STRING
            [simple string 255 'value']
        ]
        ['WSTRING' STRING
            [simple string 255 'value']
        ]
    ]
]

[enum uint 8 'SimulatedDataTypeSizes' [uint 8 'dataTypeSize']
    ['1' BOOL ['1']]
    ['2' BYTE ['1']]
    ['3' WORD ['2']]
    ['4' DWORD ['4']]
    ['5' LWORD ['8']]
    ['6' SINT ['1']]
    ['7' INT ['2']]
    ['8' DINT ['4']]
    ['9' LINT ['8']]
    ['10' USINT ['1']]
    ['11' UINT ['2']]
    ['12' UDINT ['4']]
    ['13' ULINT ['8']]
    ['14' REAL ['4']]
    ['15' LREAL ['8']]
    ['16' TIME ['8']]
    ['17' LTIME ['8']]
    ['18' DATE ['8']]
    ['19' LDATE ['8']]
    ['20' TIME_OF_DAY ['8']]
    ['21' LTIME_OF_DAY ['8']]
    ['22' DATE_AND_TIME ['8']]
    ['23' LDATE_AND_TIME ['8']]
    ['24' CHAR ['1']]
    ['25' WCHAR ['2']]
    ['26' STRING ['255']]
    ['27' WSTRING ['127']]
]
