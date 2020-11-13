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

[dataIo 'DataItem' [uint 8 'dataType', uint 8 'numberOfValues']
    [typeSwitch 'dataType','numberOfValues'
        ['1','1' BOOL
            [reserved uint 7 '0x00']
            [simple   bit    'value']
        ]
        ['1' BOOL
            [array bit 'value' count 'numberOfValues']
        ]
        ['10','1' BYTE
            [simple uint 8 'value']
        ]
        ['10' BYTE
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['11','1' WORD
            [simple uint 16 'value']
        ]
        ['11' WORD
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['12','1' DWORD
            [simple uint 32 'value']
        ]
        ['12' DWORD
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['13','1' LWORD
            [simple uint 64 'value']
        ]
        ['13' LWORD
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['20','1' SINT
            [simple int 8 'value']
        ]
        ['20' SINT
            [array int 8 'value' count 'numberOfValues']
        ]
        ['21','1' INT
            [simple int 16 'value']
        ]
        ['21' INT
            [array int 16 'value' count 'numberOfValues']
        ]
        ['22','1' DINT
            [simple int 32 'value']
        ]
        ['22' DINT
            [array int 32 'value' count 'numberOfValues']
        ]
        ['23','1' LINT
            [simple int 64 'value']
        ]
        ['23' LINT
            [array int 64 'value' count 'numberOfValues']
        ]
        ['24','1' USINT
            [simple uint 8 'value']
        ]
        ['24' USINT
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['25','1' UINT
            [simple uint 16 'value']
        ]
        ['25' UINT
            [array uint 16 'value' count 'numberOfValues']
        ]
        ['26','1' UDINT
            [simple uint 32 'value']
        ]
        ['26' UDINT
            [array uint 32 'value' count 'numberOfValues']
        ]
        ['27','1' ULINT
            [simple uint 64 'value']
        ]
        ['27' ULINT
            [array uint 64 'value' count 'numberOfValues']
        ]
        ['30','1' REAL
            [simple float 8.23  'value']
        ]
        ['30' REAL
            [array float 8.23 'value' count 'numberOfValues']
        ]
        ['31','1' LREAL
            [simple float 11.52  'value']
        ]
        ['31' LREAL
            [array float 11.52 'value' count 'numberOfValues']
        ]
        ['80','1' CHAR
            [simple uint 8 'value']
        ]
        ['80' CHAR
            [array uint 8 'value' count 'numberOfValues']
        ]
        ['81','1' WCHAR
            [simple uint 16 'value']
        ]
        ['81' WCHAR
            [array uint 16 'value' count 'numberOfValues']
        ]
    ]
]

[enum uint 8 'SimulatedDataType' [uint 8 'dataTypeSize']
    ['00' NULL ['0']]
    ['01' BOOL ['1']]
    ['10' BYTE ['1']]
    ['11' WORD ['2']]
    ['12' DWORD ['4']]
    ['13' LWORD ['8']]
    ['20' SINT ['1']]
    ['21' INT ['2']]
    ['22' DINT ['4']]
    ['23' LINT ['8']]
    ['24' USINT ['1']]
    ['25' UINT ['2']]
    ['26' UDINT ['4']]
    ['27' ULINT ['8']]
    ['30' REAL ['4']]
    ['31' LREAL ['8']]
    ['40' TIME ['8']]
    ['41' LTIME ['8']]
    ['50' DATE ['8']]
    ['51' LDATE ['8']]
    ['60' TIME_OF_DAY ['8']]
    ['61' LTIME_OF_DAY ['8']]
    ['70' DATE_AND_TIME ['8']]
    ['71' LDATE_AND_TIME ['8']]
    ['80' CHAR ['1']]
    ['81' WCHAR ['2']]
    ['82' STRING ['1']]
    ['83' WSTRING ['2']]
]
