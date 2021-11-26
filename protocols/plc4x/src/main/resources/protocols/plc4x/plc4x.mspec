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

[discriminatedType Request byteOrder='BIG_ENDIAN'
    // TODO: Replace this with an discriminator field asap
    [simple RequestType type]
    [typeSwitch type
        ['RequestType.ReadRequest' 'ReadRequest'
            [simple
        ]
        ['RequestType.ReadResponse' 'ReadResponse'

        ]
        ['RequestType.WriteRequest' 'WriteRequest'

        ]
        ['RequestType.WriteResponse' 'WriteResponse'

        ]
    ]
]

[enum uint 8 RequestType
    ['0x01' ReadRequest]
    ['0x02' ReadResponse]
    ['0x03' WriteRequest]
    ['0x04' WriteResponse]
]