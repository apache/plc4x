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

package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
)

// Serializable indicates that something is serializable
type Serializable interface {
	utils.Serializable
}

// ReadBufferByteBased is a byte based read buffer
type ReadBufferByteBased interface {
	utils.ReadBufferByteBased
}

// NewReadBufferByteBased creates a byte based read buffer
func NewReadBufferByteBased(data []byte) ReadBufferByteBased {
	return utils.NewReadBufferByteBased(data)
}

// NewLittleEndianReadBufferByteBased creates a little endian byte based read buffer
func NewLittleEndianReadBufferByteBased(data []byte) ReadBufferByteBased {
	return utils.NewLittleEndianReadBufferByteBased(data)
}

// WriteBufferByteBased is a byte based write buffer
type WriteBufferByteBased interface {
	utils.WriteBufferByteBased
}

// NewWriteBufferByteBased creates a byte based write buffer
func NewWriteBufferByteBased() WriteBufferByteBased {
	return utils.NewWriteBufferByteBased()
}

// NewLittleEndianWriteBufferByteBased creates a little endian byte write read buffer
func NewLittleEndianWriteBufferByteBased() WriteBufferByteBased {
	return utils.NewLittleEndianWriteBufferByteBased()
}

// NewCustomWriteBufferByteBased creates a byte base write buffer with a custom underlying bytes.Buffer
func NewCustomWriteBufferByteBased(buffer *bytes.Buffer, byteOrder binary.ByteOrder) WriteBufferByteBased {
	return utils.NewCustomWriteBufferByteBased(buffer, byteOrder)
}
