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

package model

import "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"

func FirmataUtilsIsSysexEnd(io utils.ReadBuffer) bool {
	return io.(utils.ReadBufferByteBased).PeekByte(0) == 0xF7
}

func FirmataUtilsParseSysexString(io utils.ReadBuffer) int8 {
	aByte, err := io.ReadInt8("", 8)
	if err != nil {
		return 0
	}
	// Skip the empty byte.
	_, _ = io.ReadByte("")
	return aByte
}

func FirmataUtilsSerializeSysexString(io utils.WriteBuffer, data int8) {
	_ = io.WriteByte("", byte(data))
	_ = io.WriteByte("", 0x00)
}

func FirmataUtilsLengthSysexString(data []int8) uint16 {
	return uint16(len(data) * 2)
}
