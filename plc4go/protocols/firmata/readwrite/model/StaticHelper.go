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
	"context"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

func IsSysexEnd(ctx context.Context, io utils.ReadBuffer) func([]byte) bool {
	return func(bytes []byte) bool {
		return io.(utils.ReadBufferByteBased).PeekByte(0) == 0xF7
	}
}

func ParseSysexString(ctx context.Context, io utils.ReadBuffer) func(context.Context) (byte, error) {
	return func(context.Context) (byte, error) {
		aByte, err := io.ReadByte("")
		if err != nil {
			return 0, err
		}
		// Skip the empty byte.
		_, _ = io.ReadByte("")
		return aByte, err
	}
}

func SerializeSysexString(ctx context.Context, io utils.WriteBuffer, data byte) error {
	if err := io.WriteByte("", data); err != nil {
		return err
	}
	if err := io.WriteByte("", 0x00); err != nil {
		return err
	}
	return nil
}

func LengthSysexString(ctx context.Context, data []byte) uint16 {
	return uint16(len(data) * 2)
}
