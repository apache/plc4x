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

// Remark: The different fields are encoded in Big-endian.

[dataIo 'DataItem' [string 'dataType', uint 16 'numberOfValues']
    [typeSwitch 'dataType','numberOfValues'
        ['IEC61131_BOOL','1' BOOL
            [simple   bit    'value']
        ]
        ['IEC61131_BOOL' List
            [array bit 'value' count 'numberOfValues']
        ]
        ['IEC61131_BYTE','1' BYTE
            [simple uint 8 'value']
        ]
        ['IEC61131_BYTE' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['IEC61131_WORD','1' WORD
            [simple uint 16 'value']
        ]
        ['IEC61131_WORD' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['IEC61131_DWORD','1' DWORD
            [simple uint 32 'value']
        ]
        ['IEC61131_DWORD' List
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['IEC61131_LWORD','1' LWORD
            [simple uint 64 'value']
        ]
        ['IEC61131_LWORD' List
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['IEC61131_SINT','1' SINT
            [simple int 8 'value']
        ]
        ['IEC61131_SINT' List
            [array int 8 'value' count 'numberOfValues']
        ]
        ['IEC61131_INT','1' INT
            [simple int 16 'value']
        ]
        ['IEC61131_INT' List
            [array int 16 'value' count 'numberOfValues']
        ]
        ['IEC61131_DINT','1' DINT
            [simple int 32 'value']
        ]
        ['IEC61131_DINT' List
            [array int 32 'value' count 'numberOfValues']
        ]
        ['IEC61131_LINT','1' LINT
            [simple int 64 'value']
        ]
        ['IEC61131_LINT' List
            [array int 64 'value' count 'numberOfValues']
        ]
        ['IEC61131_USINT','1' USINT
            [simple uint 8 'value']
        ]
        ['IEC61131_USINT' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['IEC61131_UINT','1' UINT
            [simple uint 16 'value']
        ]
        ['IEC61131_UINT' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['IEC61131_UDINT','1' UDINT
            [simple uint 32 'value']
        ]
        ['IEC61131_UDINT' List
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['IEC61131_ULINT','1' ULINT
            [simple uint 64 'value']
        ]
        ['IEC61131_ULINT' List
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['IEC61131_REAL','1' REAL
            [simple float 8.23  'value']
        ]
        ['IEC61131_REAL' List
            [array float 8.23 'value' count 'numberOfValues']
        ]
        ['IEC61131_LREAL','1' LREAL
            [simple float 11.52  'value']
        ]
        ['IEC61131_LREAL' List
            [array float 11.52 'value' count 'numberOfValues']
        ]
        ['IEC61131_CHAR','1' CHAR
            [simple uint 8 'value']
        ]
        ['IEC61131_CHAR' List
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['IEC61131_WCHAR','1' WCHAR
            [simple uint 16 'value']
        ]
        ['IEC61131_WCHAR' List
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['IEC61131_STRING' STRING
            [manual string 'UTF-8' 'value' 'STATIC_CALL("org.apache.plc4x.java.simulated.utils.StaticHelper.parsePascalString", io, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.simulated.utils.StaticHelper.serializePascalString", io, _value, _type.encoding)' '_value.length + 2']
        ]
        ['IEC61131_WSTRING' STRING
            [manual string 'UTF-16' 'value''STATIC_CALL("org.apache.plc4x.java.simulated.utils.StaticHelper.parsePascalString", io, _type.encoding)' 'STATIC_CALL("org.apache.plc4x.java.simulated.utils.StaticHelper.serializePascalString", io, _value, _type.encoding)' '(_value.length * 2) + 2']
        ]
    ]
]

[enum string 'SimulatedDataTypeSizes' [uint 8 'dataTypeSize']
    ['IEC61131_BOOL' BOOL ['1']]
    ['IEC61131_BYTE' BYTE ['1']]
    ['IEC61131_WORD' WORD ['2']]
    ['IEC61131_DWORD' DWORD ['4']]
    ['IEC61131_LWORD' LWORD ['8']]
    ['IEC61131_SINT' SINT ['1']]
    ['IEC61131_INT' INT ['2']]
    ['IEC61131_DINT' DINT ['4']]
    ['IEC61131_LINT' LINT ['8']]
    ['IEC61131_USINT' USINT ['1']]
    ['IEC61131_UINT' UINT ['2']]
    ['IEC61131_UDINT' UDINT ['4']]
    ['IEC61131_ULINT' ULINT ['8']]
    ['IEC61131_REAL' REAL ['4']]
    ['IEC61131_LREAL' LREAL ['8']]
    ['IEC61131_TIME' TIME ['8']]
    ['IEC61131_LTIME' LTIME ['8']]
    ['IEC61131_DATE' DATE ['8']]
    ['IEC61131_LDATE' LDATE ['8']]
    ['IEC61131_TIME_OF_DAY' TIME_OF_DAY ['8']]
    ['IEC61131_LTIME_OF_DAY' LTIME_OF_DAY ['8']]
    ['IEC61131_DATE_AND_TIME' DATE_AND_TIME ['8']]
    ['IEC61131_LDATE_AND_TIME' LDATE_AND_TIME ['8']]
    ['IEC61131_CHAR' CHAR ['1']]
    ['IEC61131_WCHAR' WCHAR ['2']]
    ['IEC61131_STRING' STRING ['256']]
    ['IEC61131_WSTRING' WSTRING ['127']]
]
