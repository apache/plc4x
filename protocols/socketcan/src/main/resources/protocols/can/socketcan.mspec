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

[type 'SocketCANFrame'
    [simple int 32 'rawId']
    [virtual int 32 'identifier'
        'STATIC_CALL("org.apache.plc4x.java.transport.socketcan.helper.HeaderParser.readIdentifier", rawId)'
    ]
    [virtual bit 'extended'
        'STATIC_CALL("org.apache.plc4x.java.transport.socketcan.helper.HeaderParser.isExtended", rawId)'
    ]
    [virtual bit 'remote'
        'STATIC_CALL("org.apache.plc4x.java.transport.socketcan.helper.HeaderParser.isRemote", rawId)'
    ]
    [virtual bit 'error'
        'STATIC_CALL("org.apache.plc4x.java.transport.socketcan.helper.HeaderParser.isError", rawId)'
    ]
    [implicit uint 8 'size' 'COUNT(data)']
    [reserved uint 8 '0x0'] //flags
    [reserved uint 8 '0x0'] // padding 1
    [reserved uint 8 '0x0'] // padding 2
    [array byte 'data' count 'size']
    [padding uint 8 'alignment' '0x00' '8 - (COUNT(data))']
]