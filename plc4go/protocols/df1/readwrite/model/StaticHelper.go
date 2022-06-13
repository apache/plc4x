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

package model

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/snksoft/crc"
)

var table *crc.Table

func init() {
	// CRC-16/DF-1
	table = crc.NewTable(&crc.Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000})
}

func CrcCheck(destinationAddress uint8, sourceAddress uint8, command *DF1Command) (uint16, error) {
	df1Crc := table.InitCrc()
	df1Crc = table.UpdateCrc(df1Crc, []byte{destinationAddress, sourceAddress})
	bufferByteBased := utils.NewWriteBufferByteBased()
	err := command.Serialize(bufferByteBased)
	if err != nil {
		return 0, err
	}
	bytes := bufferByteBased.GetBytes()
	df1Crc = table.UpdateCrc(df1Crc, bytes)
	df1Crc = table.UpdateCrc(df1Crc, []byte{0x03})
	return table.CRC16(df1Crc), nil
}

func DataTerminate(io utils.ReadBuffer) bool {
	rbbb := io.(utils.ReadBufferByteBased)
	// The byte sequence 0x10 followed by 0x03 indicates the end of the message,
	// so if we would read this, we abort the loop and stop reading data.
	return rbbb.PeekByte(0) == 0x10 && rbbb.PeekByte(1) == 0x03
}

func ReadData(io utils.ReadBuffer) uint8 {
	rbbb := io.(utils.ReadBufferByteBased)
	// If we read a 0x10, this has to be followed by another 0x10, which is how
	// this value is escaped in DF1, so if we encounter two 0x10, we simply ignore the first.
	if rbbb.PeekByte(0) == 0x10 && rbbb.PeekByte(1) == 0x10 {
		_, _ = io.ReadUint8("", 8)
	}
	data, _ := io.ReadUint8("", 8)
	return data
}

func WriteData(io utils.WriteBuffer, element uint8) {
	if element == 0x10 {
		// If a value is 0x10, this has to be duplicated in order to be escaped.
		_ = io.WriteUint8("", 8, element)
	}
	_ = io.WriteUint8("", 8, element)
}

func DataLength(data []byte) uint16 {
	length := uint16(0)
	for _, datum := range data {
		if datum == 0x10 {
			// If a value is 0x10, this has to be duplicated which increases the message size by one.
			length++
		}
		length++
	}
	return length
}
