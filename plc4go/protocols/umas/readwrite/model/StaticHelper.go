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
	"unicode/utf8"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// This file provides the STATIC_CALL implementations referenced by the UMAS
// mspec. It is hand written (the code generator only emits the call sites) and
// is deliberately not marked as generated so regeneration leaves it alone.

// umasStringLength is the set of integer types the code generator may pass as a
// string length. The generator picks the Go type from the mspec (a literal
// becomes an untyped constant, a parser argument keeps its declared width), so
// the helpers accept any of them rather than forcing a cast into the generated
// call sites.
type umasStringLength interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// terminatorByte is the NUL byte that ends a UMAS string.
const terminatorByte = 0x00

// ParseTerminatedString reads a NUL terminated UMAS string.
//
// A negative stringLength selects variable length mode: bytes are consumed one
// at a time until the NUL terminator is found (and the terminator is consumed).
// Running out of buffer before the terminator is an error.
//
// A non negative stringLength selects fixed width mode: exactly stringLength
// bytes are consumed and everything up to the first NUL byte is returned. If no
// NUL byte is present the whole field is the string.
func ParseTerminatedString[L umasStringLength](ctx context.Context, io utils.ReadBuffer, stringLength L) (string, error) {
	if stringLength < 0 {
		var stringBytes []byte
		for {
			curByte, err := io.ReadByte("value")
			if err != nil {
				return "", errors.Wrap(err, "Error reading string byte (no terminator found)")
			}
			if curByte == terminatorByte {
				break
			}
			stringBytes = append(stringBytes, curByte)
		}
		return string(stringBytes), nil
	}
	return ParseTerminatedStringBytes(ctx, io, stringLength)
}

// ParseTerminatedStringBytes reads a fixed width field of numberOfValues bytes
// and returns the content up to the first NUL byte. It backs the STRING variant
// of the DataItem dataIo, where numberOfValues is the number of characters the
// tag declares.
func ParseTerminatedStringBytes[L umasStringLength](_ context.Context, io utils.ReadBuffer, numberOfValues L) (string, error) {
	numberOfBytes := int(numberOfValues)
	if numberOfBytes < 0 {
		return "", errors.Errorf("Invalid negative string length %d", numberOfBytes)
	}
	if numberOfBytes == 0 {
		return "", nil
	}
	rawBytes, err := io.ReadByteArray("value", numberOfBytes)
	if err != nil {
		return "", errors.Wrap(err, "Error reading string bytes")
	}
	for i, curByte := range rawBytes {
		if curByte == terminatorByte {
			return string(rawBytes[:i]), nil
		}
	}
	return string(rawBytes), nil
}

// SerializeTerminatedString writes a NUL terminated UMAS string. The value is
// either a plain string (regular types) or a PlcValue (the DataItem dataIo).
//
// A negative stringLength selects variable length mode: the UTF-8 bytes of the
// value followed by a single NUL terminator.
//
// A non negative stringLength selects fixed width mode: exactly stringLength
// bytes are written, the value NUL padded on the right. A value that does not
// fit is truncated on a rune boundary, so the field never carries half a
// multi byte character.
func SerializeTerminatedString[T any, L umasStringLength](_ context.Context, io utils.WriteBuffer, value T, stringLength L) error {
	stringValue, err := umasStringValue(value)
	if err != nil {
		return err
	}

	if stringLength < 0 {
		if err := io.WriteString("value", uint32(len(stringValue)*8), stringValue, utils.WithEncoding("UTF8")); err != nil {
			return errors.Wrap(err, "Error serializing string value")
		}
		if err := io.WriteByte("terminator", terminatorByte); err != nil {
			return errors.Wrap(err, "Error serializing terminator byte")
		}
		return nil
	}

	numberOfBytes := int(stringLength)
	// WriteString pads the remaining bits of the field with NUL bytes, so
	// truncating the value is all that is left to do here.
	if err := io.WriteString("value", uint32(numberOfBytes*8), truncateToBytes(stringValue, numberOfBytes), utils.WithEncoding("UTF8")); err != nil {
		return errors.Wrap(err, "Error serializing string value")
	}
	return nil
}

// WriteSizeIndexToByteCount converts a UMAS size index (the requestSize of
// UmasDataType) into the number of bytes one element of that type occupies:
// 1 -> 1, 2 -> 2, 3 -> 4, 4 -> 8. Index 17 marks STRING, whose elements are one
// byte each. Any other index is passed through as a raw byte count.
func WriteSizeIndexToByteCount[L umasStringLength](_ context.Context, sizeIndex L) uint16 {
	index := int64(sizeIndex)
	switch {
	case index == 17:
		// STRING: one byte per character, the element count comes from the
		// arrayLength field of the write reference.
		return 1
	case index >= 1 && index <= 4:
		return uint16(1) << uint16(index-1)
	default:
		if index < 0 {
			return 0
		}
		return uint16(index)
	}
}

// umasStringValue extracts the string to serialize from whatever the generated
// call site passes: the generated types hand over a plain string, the DataItem
// dataIo hands over a PlcValue.
func umasStringValue(value any) (string, error) {
	switch typedValue := value.(type) {
	case nil:
		return "", nil
	case string:
		return typedValue, nil
	case apiValues.PlcValue:
		return typedValue.GetString(), nil
	default:
		return "", errors.Errorf("Unsupported value type %T for a terminated string", value)
	}
}

// truncateToBytes shortens value to at most maxBytes bytes without splitting a
// multi byte rune.
func truncateToBytes(value string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(value) <= maxBytes {
		return value
	}
	cutoff := maxBytes
	for cutoff > 0 && !utf8.RuneStart(value[cutoff]) {
		cutoff--
	}
	return value[:cutoff]
}
