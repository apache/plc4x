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

use std::io::{Error, ErrorKind, Read, Write};
use crate::{Message, NoOption, plc4x_enum, ReadBuffer, WriteBuffer};

// [enum uint 8 ModbusErrorCode
//     ['1'    ILLEGAL_FUNCTION]
//     ['2'    ILLEGAL_DATA_ADDRESS]
//     ['3'    ILLEGAL_DATA_VALUE]
//     ['4'    SLAVE_DEVICE_FAILURE]
//     ['5'    ACKNOWLEDGE]
//     ['6'    SLAVE_DEVICE_BUSY]
//     ['7'    NEGATIVE_ACKNOWLEDGE]
//     ['8'    MEMORY_PARITY_ERROR]
//     ['10'   GATEWAY_PATH_UNAVAILABLE]
//     ['11'   GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
// ]
plc4x_enum!
[enum u8 : ModbusErrorCode
    [1 =>   ILLEGAL_FUNCTION],
    [2 =>   ILLEGAL_DATA_ADDRESS],
    [3 =>   ILLEGAL_DATA_VALUE],
    [4 =>   SLAVE_DEVICE_FAILURE],
    [5 =>   ACKNOWLEDGE],
    [6 =>   SLAVE_DEVICE_BUSY],
    [7 =>   NEGATIVE_ACKNOWLEDGE],
    [8 =>   MEMORY_PARITY_ERROR],
    [10 =>  GATEWAY_PATH_UNAVAILABLE],
    [11 =>  GATEWAY_TARGET_DEVICE_FAILED_TO_RESPOND]
];