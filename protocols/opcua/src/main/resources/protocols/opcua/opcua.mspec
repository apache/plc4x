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

[enum uint 8 'OpcuaDataType'
    ['0' NULL ]
    ['1' BOOL ]
    ['2' BYTE ]
    ['3' WORD ]
    ['4' DWORD ]
    ['5' LWORD ]
    ['6' SINT ]
    ['7' INT ]
    ['8' DINT ]
    ['9' LINT ]
    ['10' USINT ]
    ['11' UINT ]
    ['12' UDINT ]
    ['13' ULINT ]
    ['14' REAL ]
    ['15' LREAL ]
    ['16' TIME ]
    ['17' LTIME ]
    ['18' DATE ]
    ['19' LDATE ]
    ['20' TIME_OF_DAY ]
    ['21' LTIME_OF_DAY ]
    ['22' DATE_AND_TIME ]
    ['23' LDATE_AND_TIME ]
    ['24' CHAR ]
    ['25' WCHAR ]
    ['26' STRING ]
    ['27' WSTRING ]
]


[enum string '1' 'OpcuaIdentifierType'
    ['s' STRING_IDENTIFIER]
    ['i' NUMBER_IDENTIFIER]
    ['g' GUID_IDENTIFIER]
    ['b' BINARY_IDENTIFIER]
]
