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

[enum ChunkHeaderMarker
    ['0x05' ChunkHeader  ]
]

[enum TSDataType
    ['0x05' INT32  ]
]

[enum CompressionType
    ['0x00' UNCOMPRESSED]
]

[enum TSEncoding
    ['0x00' PLAIN]
]

[type IoTDBString
    [implicit uint 8 length 'COUNT(content)']
    [array byte content count   'length']
]

[type ChunkGroupHeader
    [const uint 8 marker 0x00]
    [simple IoTDBString deviceId]
]

[type ChunkHeader(ChunkHeaderMarker marker)
    [simple     IoTDBString     measurementId]
    [simple     uint 8    dataSize]
    [simple     TSDataType    dataType]
    [simple     CompressionType    compression]
    [simple     TSEncoding    encoding]
]