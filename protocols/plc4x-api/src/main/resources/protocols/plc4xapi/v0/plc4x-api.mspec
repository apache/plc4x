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

[enum uint 8 PlcResponseCode
    ['0x01' OK              ]
    ['0x02' NOT_FOUND       ]
    ['0x03' ACCESS_DENIED   ]
    ['0x04' INVALID_ADDRESS ]
    ['0x06' INVALID_DATATYPE]
    ['0x07' INVALID_DATA    ]
    ['0x08' INTERNAL_ERROR  ]
    ['0x09' REMOTE_BUSY     ]
    ['0x0A' REMOTE_ERROR    ]
    ['0x0B' UNSUPPORTED     ]
    ['0x0C' RESPONSE_PENDING]
]

[enum uint 8 PlcValueType
    ['0x00' NULL          ]

    // Bit Strings
    ['0x01' BOOL          ]
    ['0x02' BYTE          ]
    ['0x03' WORD          ]
    ['0x04' DWORD         ]
    ['0x05' LWORD         ]

    // Unsigned Integers
    ['0x11' USINT         ]
    ['0x12' UINT          ]
    ['0x13' UDINT         ]
    ['0x14' ULINT         ]

    // Signed Integers
    ['0x21' SINT          ]
    ['0x22' INT           ]
    ['0x23' DINT          ]
    ['0x24' LINT          ]

    // Floating Point Values
    ['0x31' REAL          ]
    ['0x32' LREAL         ]

    // Chars and Strings
    ['0x41' CHAR          ]
    ['0x42' WCHAR         ]
    ['0x43' STRING        ]
    ['0x44' WSTRING       ]

    // Times and Dates
    ['0x51' TIME          ]
    ['0x52' LTIME         ]
    ['0x53' DATE          ]
    ['0x54' LDATE         ]
    ['0x55' TIME_OF_DAY   ]
    ['0x56' LTIME_OF_DAY  ]
    ['0x57' DATE_AND_TIME ]
    ['0x58' LDATE_AND_TIME]

    // Complex types
    ['0x61' Struct        ]
    ['0x62' List          ]

    ['0x71' RAW_BYTE_ARRAY]
]

[enum uint 8 PlcSubscriptionType
   ['0x01' CYCLIC         ]
   ['0x02' CHANGE_OF_STATE]
   ['0x03' EVENT          ]
]