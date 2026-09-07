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

import (
	"context"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// ParseZeroTerminatedString reads `stringValueLength` bytes from the buffer as
// an ASCII string and then consumes a single trailing 0x00 terminator byte.
// It returns an error if the terminator byte is not 0x00.
func ParseZeroTerminatedString(ctx context.Context, io utils.ReadBuffer, stringValueLength uint16) (string, error) {
	// Read the raw bytes that make up the string content.
	rawBytes, err := io.ReadByteArray("stringValueLength", int(stringValueLength))
	if err != nil {
		return "", errors.Wrap(err, "Error reading string bytes")
	}

	// Consume the terminator.
	terminatorByte, err := io.ReadByte("terminator")
	if err != nil {
		return "", errors.Wrap(err, "Error reading terminator byte")
	}
	if terminatorByte != 0x00 {
		return "", errors.Errorf("Expected 0x00, but found 0x%02x", terminatorByte)
	}

	return string(rawBytes), nil
}

// SerializeZeroTerminatedString writes the given string as ASCII followed by a
// single 0x00 terminator byte.
func SerializeZeroTerminatedString(ctx context.Context, io utils.WriteBuffer, data string) error {
	if err := io.WriteString("stringValue", uint32(len(data)*8), data, utils.WithEncoding("ASCII")); err != nil {
		return errors.Wrap(err, "Error serializing string value")
	}
	if err := io.WriteByte("terminator", 0x00); err != nil {
		return errors.Wrap(err, "Error serializing terminator byte")
	}
	return nil
}

// LengthZeroTerminatedString returns the length in bits of a zero-terminated
// string (string bytes + 1 terminator byte, times 8).
func LengthZeroTerminatedString(ctx context.Context, data string) uint16 {
	return uint16((len(data) + 1) * 8)
}
