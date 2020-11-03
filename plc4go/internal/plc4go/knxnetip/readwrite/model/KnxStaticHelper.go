//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package model

import (
    "math"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

func KnxHelperBytesToF16(io *utils.ReadBuffer) (float32, error) {
    negative, err := io.ReadBit()
    if err != nil {
        return 0.0, err
    }
    exponent, err := io.ReadUint64(4)
    if err != nil {
        return 0.0, err
    }
    mantissa, err := io.ReadUint64(11)
    if err != nil {
        return 0.0, err
    }
    mantissaPart := 0.01 * float32(mantissa)
    powPart := math.Pow(float64(2), float64(exponent))
    if negative {
        return -1 * mantissaPart * float32(powPart), nil
    } else {
        return mantissaPart * float32(powPart), nil
    }
}
