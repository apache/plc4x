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
[discriminatedType 'MBusFrame'
    [discriminator uint 8 'header'               ]
    [typeSwitch 'header'
        ['0xE5' AcknowledgeFrame
        ]
        ['0x10' ShortFrame
        ]
        ['0x68' DataFrame
            [simple  uint 8 'len'                   ]
            [padding uint 8 'lenCheck' 'len' '1'    ]
            [const   uint 8 'headerCheck' '0x68'    ]
            [simple  uint 8 'controlField'          ]
            [simple  uint 8 'address'               ]
            [simple  uint 8 'controlInformation'    ]
            [simple  VariableDataStructure 'data' ['controlField', 'controlInformation', 'len'] ]
            [checksum int 8  'crc' 'STATIC_CALL("org.apache.plc4x.java.mbus.helper.MBusHelper.crc", checksumRawData)']
            [const   uint 8 'endFlag' '0x16'        ]
        ]
    ]
]

[type 'VariableDataStructure' [uint 8 'controlField', uint 8 'ci', uint 8 'len']
    [typeSwitch 'ci'
        ['0x72' LongHeaderDataFrame
            [simple SecondaryAddress 'address']
            [simple ShortHeader 'header']
        ]
    ]
    [array DataRecord 'records' terminated 'curPos > len - 7']
]

[type 'SecondaryAddress'
    [simple uint 32 'identifier']
    [simple uint 16 'manufacturer']
    [simple uint 8 'version']
    [simple uint 8 'deviceType']
]

[type 'ShortHeader'
    [simple uint 8 'accessNumber']
    [simple uint 8 'status']
    [simple uint 8 'encryptedBlocks']
    [simple uint 8 'encryptionMode' ]
]

[type 'DataRecord'
    [implicit bit 'extension' 'dife != null']
    [simple uint 1 'storage']
    [enum Function 'function']
    [enum Coding 'coding']
    [optional DataInformationFieldExtension 'dife' 'extension']
    [simple ValueInformationBlock 'vif']
    [array uint 8 'data' length 'coding.size / 8']
]

[type 'DataInformationField'
    [implicit bit 'extension' 'dife != null']
    [simple uint 1 'storage']
    [enum Function 'function']
    [enum Coding 'coding']
    [optional DataInformationFieldExtension 'dife' 'extension']
]

[type 'DataInformationFieldExtension'
    [implicit bit 'extension' 'dife != null']
    [simple uint 1 'unit']
    [simple uint 2 'tariff']
    [simple uint 4 'storage']
    [optional DataInformationFieldExtension 'dife' 'extension']
]

[type 'ValueInformationBlock'
    [simple ValueInformationField 'vif']
]

[type 'ValueInformationField'
    [implicit bit    'extension' 'vife != null']
    [simple uint 7 'unitAndMultiplier']
    [optional ValueInformationField 'vife' 'extension']
]

[enum uint 4 'Coding' [uint 8 'size']
    ['0x00' NO_DATA            ['0'] ]
    ['0x01' INT8               ['8'] ]
    ['0x02' INT16             ['16'] ]
    ['0x03' INT24             ['24'] ]
    ['0x04' INT32             ['32'] ]
    ['0x05' REAL              ['32'] ]
    ['0x06' INT48             ['48'] ]
    ['0x07' INT64             ['64'] ]
    ['0x08' READOUT            ['0'] ]
    ['0x09' BCD_2_DIGIT        ['8'] ]
    ['0x0A' BCD_4_DIGIT       ['16'] ]
    ['0x0B' BCD_6_DIGIT       ['24'] ]
    ['0x0C' BCD_8_DIGIT       ['32'] ]
    ['0x0D' VARIABLE_LEN      ['-1'] ]
    ['0x0E' BCD_12_DIGIT      ['48'] ]
    ['0x0F' SPECIAL_FUNCTION  ['-2'] ]
]

[enum uint 2 'Function'
    ['0x00' INSTANT_VALUE]
    ['0x01' MAXIMUM_VALUE]
    ['0x02' MINIMUM_VALUE]
    ['0x03' ERROR_VALUE]
]