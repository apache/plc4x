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

package profinet

import (
	"encoding/binary"
	"math"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/profinet/gsdml"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// itemLocation is where a tag's data item sits in the frame.
type itemLocation struct {
	item   gsdml.DataItem
	offset int // into the C_SDU
	size   int
}

// locate finds a tag's data item in the plan.
func (p *plan) locate(tag Tag) (itemLocation, error) {
	s, ok := p.Find(tag.Slot, tag.Subslot)
	if !ok {
		return itemLocation{}, errors.Errorf("the device has no submodule in slot %d subslot %d", tag.Slot, tag.Subslot)
	}
	sets := s.Resolved.Submodule.IOData.Inputs
	base := s.InputData
	if tag.Output {
		sets = s.Resolved.Submodule.IOData.Outputs
		base = s.OutputData
	}
	offset, index := 0, 0
	for _, set := range sets {
		for _, item := range set.DataItems {
			size, err := item.Size()
			if err != nil {
				return itemLocation{}, err
			}
			if index == tag.Index {
				if base < 0 {
					return itemLocation{}, errors.Errorf("slot %d subslot %d has no %s data", tag.Slot, tag.Subslot, strings.ToLower(tag.direction()))
				}
				return itemLocation{item: item, offset: base + offset, size: size}, nil
			}
			offset += size
			index++
		}
	}
	return itemLocation{}, errors.Errorf("slot %d subslot %d has %d %s data items, not item %d", tag.Slot, tag.Subslot, index, strings.ToLower(tag.direction()), tag.Index)
}

// widths are the byte sizes of the value types a data item can be read as.
var widths = map[apiValues.PlcValueType]int{
	apiValues.BYTE: 1, apiValues.USINT: 1, apiValues.SINT: 1,
	apiValues.WORD: 2, apiValues.UINT: 2, apiValues.INT: 2,
	apiValues.DWORD: 4, apiValues.UDINT: 4, apiValues.DINT: 4, apiValues.REAL: 4,
	apiValues.LWORD: 8, apiValues.ULINT: 8, apiValues.LINT: 8, apiValues.LREAL: 8,
}

// decodeItem reads a tag's value from its data item's bytes. Profinet data is big-endian. BOOLs
// are the item's bits, bit 0 being the least significant bit of the item read as an integer.
func decodeItem(tag Tag, data []byte) (apiValues.PlcValue, error) {
	switch tag.ValueType {
	case apiValues.BOOL:
		if tag.Count > len(data)*8 {
			return nil, errors.Errorf("the item has %d bits, the tag asks for %d", len(data)*8, tag.Count)
		}
		bits := make([]apiValues.PlcValue, tag.Count)
		for i := range bits {
			bits[i] = spiValues.NewPlcBOOL(data[len(data)-1-i/8]&(1<<(i%8)) != 0)
		}
		if tag.Count == 1 {
			return bits[0], nil
		}
		return spiValues.NewPlcList(bits), nil
	case apiValues.STRING:
		return spiValues.NewPlcSTRING(strings.TrimRight(string(data), "\x00")), nil
	}
	width, ok := widths[tag.ValueType]
	if !ok {
		return nil, errors.Errorf("a data item can't be read as %s", tag.ValueType)
	}
	if width*tag.Count != len(data) {
		return nil, errors.Errorf("the item is %d bytes, %d x %s is %d", len(data), tag.Count, tag.ValueType, width*tag.Count)
	}
	values := make([]apiValues.PlcValue, tag.Count)
	for i := range values {
		values[i] = scalar(tag.ValueType, data[i*width:(i+1)*width])
	}
	if tag.Count == 1 {
		return values[0], nil
	}
	return spiValues.NewPlcList(values), nil
}

func scalar(valueType apiValues.PlcValueType, b []byte) apiValues.PlcValue {
	var u uint64
	for _, x := range b {
		u = u<<8 | uint64(x)
	}
	switch valueType {
	case apiValues.BYTE:
		return spiValues.NewPlcBYTE(uint8(u))
	case apiValues.USINT:
		return spiValues.NewPlcUSINT(uint8(u))
	case apiValues.SINT:
		return spiValues.NewPlcSINT(int8(u))
	case apiValues.WORD:
		return spiValues.NewPlcWORD(uint16(u))
	case apiValues.UINT:
		return spiValues.NewPlcUINT(uint16(u))
	case apiValues.INT:
		return spiValues.NewPlcINT(int16(u))
	case apiValues.DWORD:
		return spiValues.NewPlcDWORD(uint32(u))
	case apiValues.UDINT:
		return spiValues.NewPlcUDINT(uint32(u))
	case apiValues.DINT:
		return spiValues.NewPlcDINT(int32(u))
	case apiValues.REAL:
		return spiValues.NewPlcREAL(math.Float32frombits(uint32(u)))
	case apiValues.LWORD:
		return spiValues.NewPlcLWORD(u)
	case apiValues.ULINT:
		return spiValues.NewPlcULINT(u)
	case apiValues.LINT:
		return spiValues.NewPlcLINT(int64(u))
	}
	return spiValues.NewPlcLREAL(math.Float64frombits(u))
}

// encodeItem turns a value into a data item's bytes, with a mask of the bits it sets: writing some
// of an item's bits leaves the others as they are.
func encodeItem(tag Tag, size int, value apiValues.PlcValue) (data, mask []byte, err error) {
	elements := []apiValues.PlcValue{value}
	if value.IsList() {
		elements = value.GetList()
	}
	if len(elements) != tag.Count {
		return nil, nil, errors.Errorf("the tag has %d elements, the value %d", tag.Count, len(elements))
	}
	data = make([]byte, size)
	mask = make([]byte, size)
	switch tag.ValueType {
	case apiValues.BOOL:
		if tag.Count > size*8 {
			return nil, nil, errors.Errorf("the item has %d bits, the tag %d", size*8, tag.Count)
		}
		for i, element := range elements {
			bit := byte(1 << (i % 8))
			at := size - 1 - i/8
			mask[at] |= bit
			if element.GetBool() {
				data[at] |= bit
			}
		}
		return data, mask, nil
	case apiValues.STRING:
		text := value.GetString()
		if len(text) > size {
			return nil, nil, errors.Errorf("%d characters don't fit the %d byte item", len(text), size)
		}
		copy(data, text)
		for i := range mask {
			mask[i] = 0xFF
		}
		return data, mask, nil
	}
	width, ok := widths[tag.ValueType]
	if !ok || width*tag.Count != size {
		return nil, nil, errors.Errorf("the item is %d bytes, %d x %s doesn't fit it", size, tag.Count, tag.ValueType)
	}
	for i, element := range elements {
		var u uint64
		switch tag.ValueType {
		case apiValues.REAL:
			u = uint64(math.Float32bits(element.GetFloat32()))
		case apiValues.LREAL:
			u = math.Float64bits(element.GetFloat64())
		case apiValues.SINT, apiValues.INT, apiValues.DINT, apiValues.LINT:
			u = uint64(element.GetInt64())
		default:
			u = element.GetUint64()
		}
		chunk := data[i*width : (i+1)*width]
		switch width {
		case 1:
			chunk[0] = byte(u)
		case 2:
			binary.BigEndian.PutUint16(chunk, uint16(u))
		case 4:
			binary.BigEndian.PutUint32(chunk, uint32(u))
		case 8:
			binary.BigEndian.PutUint64(chunk, u)
		}
	}
	for i := range mask {
		mask[i] = 0xFF
	}
	return data, mask, nil
}
