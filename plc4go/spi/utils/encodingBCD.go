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
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// encodingNameBCD is the encoding name the mspec compiler emits for
// `encoding='"BCD"'` fields. It matches
// org.apache.plc4x.java.spi.buffers.bytebased.encoding.EncodingBCD#NAME.
const encodingNameBCD = "BCD"

// maxBCDDigits is the number of BCD digits that still fit into the uint64 the
// bit buffers exchange (16 nibbles == 64 bits).
const maxBCDDigits = 16

// bcdPow10 holds 10^n for n in [0, maxBCDDigits]. 10^16 still fits into a uint64.
var bcdPow10 = [maxBCDDigits + 1]uint64{
	1,
	10,
	100,
	1_000,
	10_000,
	100_000,
	1_000_000,
	10_000_000,
	100_000_000,
	1_000_000_000,
	10_000_000_000,
	100_000_000_000,
	1_000_000_000_000,
	10_000_000_000_000,
	100_000_000_000_000,
	1_000_000_000_000_000,
	10_000_000_000_000_000,
}

// readerArgsSelectBCD reports whether the supplied reader args select the BCD
// encoding.
//
// Semantics are "the first arg that carries an encoding decides", which is
// exactly what the canonical BufferCommons.ExtractEncoding does; both the direct
// and the nested (upcast) case therefore RETURN the comparison instead of
// continuing the scan, so a leading non-BCD encoding shadows a later BCD one in
// both implementations. TestArgsSelectBCD_AgreesWithExtractEncoding pins that
// equivalence, including for multi-encoding arg lists.
//
// This deliberately does NOT call BufferCommons.ExtractEncoding, which takes
// []WithReaderWriterArgs and would force the buffer method to upcast its own arg
// slice. Note that the caller one layer up already pays for an upcast
// (FieldReaderSimple.ReadSimpleField / FieldWriterSimple.WriteSimpleField call
// UpcastReaderArgs/UpcastWriterArgs to extract the byte order), so this saves
// the SECOND such allocation per field, not the first: the point is that the
// buffer methods themselves stay allocation free on a path every numeric field
// in every driver walks, whether or not it carries an encoding.
func readerArgsSelectBCD(args []WithReaderArgs) bool {
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case withEncoding:
			return typedArg.encoding == encodingNameBCD
		case readerWriterArg:
			// the nested form produced by UpcastReaderArgs; matches the
			// readerWriterArg case of BufferCommons.ExtractEncoding, reader
			// side first, and likewise returns rather than continuing the scan
			if nested, ok := typedArg.WithReaderArgs.(withEncoding); ok {
				return nested.encoding == encodingNameBCD
			}
			if nested, ok := typedArg.WithWriterArgs.(withEncoding); ok {
				return nested.encoding == encodingNameBCD
			}
		}
	}
	return false
}

// writerArgsSelectBCD is the WithWriterArgs twin of readerArgsSelectBCD and
// follows the same "first encoding arg wins" rule.
func writerArgsSelectBCD(args []WithWriterArgs) bool {
	for _, arg := range args {
		switch typedArg := arg.(type) {
		case withEncoding:
			return typedArg.encoding == encodingNameBCD
		case readerWriterArg:
			// the nested form produced by UpcastWriterArgs; see
			// readerArgsSelectBCD
			if nested, ok := typedArg.WithReaderArgs.(withEncoding); ok {
				return nested.encoding == encodingNameBCD
			}
			if nested, ok := typedArg.WithWriterArgs.(withEncoding); ok {
				return nested.encoding == encodingNameBCD
			}
		}
	}
	return false
}

// decodeBCD interprets raw as a stream of bitLength/4 BCD nibbles, most
// significant digit first, and returns the decimal value they denote.
//
// raw is the value as returned by ReadBitBuffer.ReadBits, i.e. right aligned in
// a uint64 with the first wire bit in the most significant position of the
// bitLength-wide window. So the nibble stream is simply walked from the top of
// that window downwards.
//
// The two languages agree here now. plc4j used to disagree for every field whose
// bit length is not a multiple of 8: its buffers right align a partial read
// inside the byte array while EncodingBCD walked nibbles from the high nibble of
// byte 0, so a 12 bit 0x123 decoded to 12 and a 4 bit nibble decoded to 0, with
// the write path broken symmetrically. That was fixed in EncodingBCD by indexing
// nibbles from the padding offset instead, which is bit identical for the byte
// multiple widths that were always correct. Only the S7 "msec" uint 12 and "dow"
// uint 4 fields are non byte multiple BCD fields in this repository (s7.mspec
// lines 307-316 for the DateAndTime type, again at 851/853 for the DATE_AND_TIME
// data item), which is why the bug stayed invisible for so long.
func decodeBCD(raw uint64, bitLength uint8) (uint64, error) {
	if bitLength%4 != 0 {
		return 0, errors.Errorf("bit length %d must be a multiple of 4 for BCD encoding", bitLength)
	}
	digits := int(bitLength / 4)
	if digits > maxBCDDigits {
		return 0, errors.Errorf("bit length %d exceeds the %d bits supported by BCD encoding", bitLength, maxBCDDigits*4)
	}
	var value uint64
	for i := range digits {
		nibble := (raw >> uint(4*(digits-1-i))) & 0x0F
		if nibble > 9 {
			return 0, errors.Errorf("invalid BCD digit: %d", nibble)
		}
		value = value*10 + nibble
	}
	return value, nil
}

// decodeBCDBounded decodes like decodeBCD and additionally rejects a value that
// does not fit the target type, mirroring EncodingBCD.decodeByte /
// EncodingBCD.decodeShort ("Decoded value too large for ...").
func decodeBCDBounded(raw uint64, bitLength uint8, maxValue uint64) (uint64, error) {
	value, err := decodeBCD(raw, bitLength)
	if err != nil {
		return 0, err
	}
	if value > maxValue {
		return 0, errors.Errorf("decoded BCD value %d is too large for the target type (max value: %d)", value, maxValue)
	}
	return value, nil
}

// encodeBCD packs value into bitLength/4 BCD nibbles, most significant digit
// first, and returns them right aligned in a uint64 ready to be handed to
// WriteBitBuffer.WriteBits. It is the exact inverse of decodeBCD.
//
// As with decodeBCD this deliberately diverges from plc4j for bit lengths that
// are not a multiple of 8; see the note on decodeBCD for the mechanism. Concretely
// encodeBCD(123, 12) == 0x123 and encodeBCD(4, 4) == 0x4, where plc4j puts 0x230
// and 0x0 on the wire respectively.
func encodeBCD(value uint64, bitLength uint8) (uint64, error) {
	if bitLength%4 != 0 {
		return 0, errors.Errorf("bit length %d must be a multiple of 4 for BCD encoding", bitLength)
	}
	digits := int(bitLength / 4)
	if digits > maxBCDDigits {
		return 0, errors.Errorf("bit length %d exceeds the %d bits supported by BCD encoding", bitLength, maxBCDDigits*4)
	}
	if value >= bcdPow10[digits] {
		return 0, errors.Errorf("value %d cannot be encoded in %d BCD digits (max value: %d)", value, digits, bcdPow10[digits]-1)
	}
	var raw uint64
	for i := digits - 1; i >= 0; i-- {
		raw |= (value % 10) << uint(4*(digits-1-i))
		value /= 10
	}
	return raw, nil
}

// encodeBCDSigned is the signed entry point. plc4j's EncodingBCD.encodeInt
// rejects negative values outright ("BCD encoding only supports non-negative
// integers"), so this does the same.
func encodeBCDSigned(value int64, bitLength uint8) (uint64, error) {
	if value < 0 {
		return 0, errors.Errorf("BCD encoding only supports non-negative integers, got %d", value)
	}
	return encodeBCD(uint64(value), bitLength)
}
