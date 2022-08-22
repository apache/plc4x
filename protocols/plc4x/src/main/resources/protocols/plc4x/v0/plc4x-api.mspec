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

[enum uint 8 Plc4xRequestType
    ['0x01' CONNECT_REQUEST        ]
    ['0x02' CONNECT_RESPONSE       ]
    ['0x03' DISCONNECT_REQUEST     ]
    ['0x04' DISCONNECT_RESPONSE    ]
    ['0x05' READ_REQUEST           ]
    ['0x06' READ_RESPONSE          ]
    ['0x07' WRITE_REQUEST          ]
    ['0x08' WRITE_RESPONSE         ]
    ['0x09' SUBSCRIPTION_REQUEST   ]
    ['0x0A' SUBSCRIPTION_RESPONSE  ]
    ['0x0B' UNSUBSCRIPTION_REQUEST ]
    ['0x0C' UNSUBSCRIPTION_RESPONSE]
]

[enum uint 8 Plc4xReturnCode
    ['0x01' OK              ]
    ['0x02' NOT_FOUND       ]
    ['0x03' ACCESS_DENIED   ]
    ['0x04' INVALID_ADDRESS ]
    ['0x05' INVALID_DATATYPE]
    ['0x06' INVALID_DATA    ]
    ['0x07' INTERNAL_ERROR  ]
    ['0x08' REMOTE_BUSY     ]
    ['0x09' REMOTE_ERROR    ]
    ['0x0A' UNSUPPORTED     ]
    ['0x0B' RESPONSE_PENDING]
]

[enum uint 8 Plc4xValueType
    ['0x00' NULL         ]

    // Bit Strings
    ['0x01' BOOL         ]
    ['0x02' BYTE         ]
    ['0x03' WORD         ]
    ['0x04' DWORD        ]
    ['0x05' LWORD        ]

    // Unsigned Integers
    ['0x11' USINT        ]
    ['0x12' UINT         ]
    ['0x13' UDINT        ]
    ['0x14' ULINT        ]

    // Signed Integers
    ['0x21' SINT         ]
    ['0x22' INT          ]
    ['0x23' DINT         ]
    ['0x24' LINT         ]

    // Floating Point Values
    ['0x31' REAL         ]
    ['0x32' LREAL        ]

    // Chars and Strings
    ['0x41' CHAR         ]
    ['0x42' WCHAR        ]
    ['0x43' STRING       ]
    ['0x44' WSTRING      ]

    // Times and Dates
    ['0x51' TIME         ]
    ['0x52' TIME_OF_DAY  ]
    ['0x53' DATE         ]
    ['0x54' DATE_AND_TIME]

    // Complex types
    ['0x61' Struct       ]
    ['0x62' List         ]
]

[enum uint 8 Plc4xResponseCode
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

[enum uint 8 Plc4xSubscriptionType
   ['0x01' CYCLIC         ]
   ['0x02' CHANGE_OF_STATE]
   ['0x03' EVENT          ]
]