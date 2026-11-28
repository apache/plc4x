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

package gsdml

import (
	"encoding/binary"
	"math"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// DefaultData builds the record a controller writes to parameterize a submodule when nobody
// changed a parameter: the Const bytes, then each Ref's DefaultValue on top. Profinet is big-endian
// on the wire, and so are multi-byte parameters.
func (r RecordDataItem) DefaultData() ([]byte, error) {
	if r.Length <= 0 {
		return nil, errors.Errorf("record %d has no Length", r.Index)
	}
	data := make([]byte, r.Length)
	for _, c := range r.Consts {
		bytes, err := parseConstData(c.Data)
		if err != nil {
			return nil, errors.Wrapf(err, "record %d", r.Index)
		}
		if c.ByteOffset < 0 || c.ByteOffset+len(bytes) > len(data) {
			return nil, errors.Errorf("record %d: Const at %d with %d bytes overruns its %d bytes", r.Index, c.ByteOffset, len(bytes), len(data))
		}
		copy(data[c.ByteOffset:], bytes)
	}
	for _, ref := range r.Refs {
		if err := ref.apply(data); err != nil {
			return nil, errors.Wrapf(err, "record %d", r.Index)
		}
	}
	return data, nil
}

// parseConstData reads "0x00,0x01,..." (decimal works too).
func parseConstData(text string) ([]byte, error) {
	var bytes []byte
	for _, field := range strings.Split(text, ",") {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		value, err := ParseNumber(field)
		if err != nil || value > 0xFF {
			return nil, errors.Errorf("bad Const byte %q", field)
		}
		bytes = append(bytes, byte(value))
	}
	return bytes, nil
}

func (ref RecordRef) apply(data []byte) error {
	need := func(n int) error {
		if ref.ByteOffset < 0 || ref.ByteOffset+n > len(data) {
			return errors.Errorf("%s parameter at byte %d overruns the %d byte record", ref.DataType, ref.ByteOffset, len(data))
		}
		return nil
	}
	defaultValue := strings.TrimSpace(ref.DefaultValue)
	switch ref.DataType {
	case "Bit", "BitArea":
		width := 1
		if ref.DataType == "BitArea" {
			width = ref.BitLength
		}
		if width < 1 || width > 16 || ref.BitOffset < 0 || ref.BitOffset+width > 16 {
			return errors.Errorf("bad bit field at byte %d bit %d width %d", ref.ByteOffset, ref.BitOffset, width)
		}
		value, err := parseInteger(defaultValue)
		if err != nil {
			return err
		}
		if value < 0 || value >= 1<<width {
			return errors.Errorf("default %d doesn't fit %d bits", value, width)
		}
		// A field within the byte sits in that byte; a wider one spans it and the next, read as one
		// big-endian 16 bit value whose bit 0 is bit 0 of the second byte.
		if ref.BitOffset+width <= 8 {
			if err := need(1); err != nil {
				return err
			}
			mask := byte((1<<width)-1) << ref.BitOffset
			data[ref.ByteOffset] = data[ref.ByteOffset]&^mask | byte(value)<<ref.BitOffset&mask
			return nil
		}
		if err := need(2); err != nil {
			return err
		}
		word := binary.BigEndian.Uint16(data[ref.ByteOffset:])
		mask := uint16((1<<width)-1) << ref.BitOffset
		binary.BigEndian.PutUint16(data[ref.ByteOffset:], word&^mask|uint16(value)<<ref.BitOffset&mask)
		return nil
	case "Unsigned8", "Integer8", "Unsigned16", "Integer16", "Unsigned32", "Integer32", "Unsigned64", "Integer64":
		size := fixedSizes[ref.DataType]
		if err := need(size); err != nil {
			return err
		}
		value, err := parseInteger(defaultValue)
		if err != nil {
			return err
		}
		bits := uint(size * 8)
		signed := strings.HasPrefix(ref.DataType, "Integer")
		if signed && bits < 64 && (value < -(1<<(bits-1)) || value >= 1<<(bits-1)) ||
			!signed && bits < 64 && (value < 0 || value >= 1<<bits) {
			return errors.Errorf("default %d doesn't fit %s", value, ref.DataType)
		}
		putBigEndian(data[ref.ByteOffset:ref.ByteOffset+size], uint64(value))
		return nil
	case "Float32", "Float64":
		size := fixedSizes[ref.DataType]
		if err := need(size); err != nil {
			return err
		}
		value := 0.0
		if defaultValue != "" {
			var err error
			if value, err = strconv.ParseFloat(defaultValue, 64); err != nil {
				return errors.Errorf("bad default %q", defaultValue)
			}
		}
		if size == 4 {
			binary.BigEndian.PutUint32(data[ref.ByteOffset:], math.Float32bits(float32(value)))
		} else {
			binary.BigEndian.PutUint64(data[ref.ByteOffset:], math.Float64bits(value))
		}
		return nil
	case "OctetString", "VisibleString":
		if ref.Length <= 0 {
			return errors.Errorf("%s parameter without a Length", ref.DataType)
		}
		if err := need(ref.Length); err != nil {
			return err
		}
		var bytes []byte
		if ref.DataType == "OctetString" {
			var err error
			if bytes, err = parseConstData(defaultValue); err != nil {
				return err
			}
		} else {
			bytes = []byte(defaultValue)
		}
		if len(bytes) > ref.Length {
			return errors.Errorf("default is %d bytes, the parameter %d", len(bytes), ref.Length)
		}
		copy(data[ref.ByteOffset:ref.ByteOffset+ref.Length], bytes)
		return nil
	}
	return errors.Errorf("unsupported parameter data type %q", ref.DataType)
}

func parseInteger(text string) (int64, error) {
	if text == "" {
		return 0, nil
	}
	negative := strings.HasPrefix(text, "-")
	value, err := ParseNumber(strings.TrimPrefix(text, "-"))
	if err != nil || value > math.MaxInt64 {
		return 0, errors.Errorf("bad default value %q", text)
	}
	if negative {
		return -int64(value), nil
	}
	return int64(value), nil
}

func putBigEndian(target []byte, value uint64) {
	for i := len(target) - 1; i >= 0; i-- {
		target[i] = byte(value)
		value >>= 8
	}
}
