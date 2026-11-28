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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// fixedSizes are the byte sizes of the GSDML I/O data types with a fixed size.
var fixedSizes = map[string]int{
	"Integer8": 1, "Unsigned8": 1,
	"Integer16": 2, "Unsigned16": 2,
	"Integer32": 4, "Unsigned32": 4,
	"Integer64": 8, "Unsigned64": 8,
	"Float32": 4, "Float64": 8,
	"Date":                                   7,
	"TimeOfDay with date indication":         6,
	"TimeOfDay without date indication":      4,
	"TimeDifference with date indication":    6,
	"TimeDifference without date indication": 4,
	"NetworkTime":                            8,
	"NetworkTimeDifference":                  8,
	"Unsigned8+Unsigned8":                    2,
	"Float32+Unsigned8":                      5,
	"Float32+Status8":                        5,
	"F_MessageTrailer4Byte":                  4,
	"F_MessageTrailer5Byte":                  5,
}

// Size is how many bytes a data item takes in the cyclic data.
func (d DataItem) Size() (int, error) {
	switch d.DataType {
	case "OctetString", "VisibleString":
		if d.Length <= 0 {
			return 0, errors.Errorf("%s without a Length", d.DataType)
		}
		return d.Length, nil
	}
	if size, ok := fixedSizes[d.DataType]; ok {
		return size, nil
	}
	return 0, errors.Errorf("unknown GSDML data type %q", d.DataType)
}

// ValueType is the value type and element count a data item maps to, ported from plc4j's
// ProfinetDataTypeMapper: integers used as bits become that many BOOLs, OctetStrings bytes (or one
// BYTE/WORD/DWORD/LWORD when used as bits), the numeric types their IEC types. Where plc4j gives up
// (dates, times, F-trailers, the compound types) this returns the item's raw bytes, so a submodule
// with such an item can still be connected and its other items used.
func (d DataItem) ValueType() (apiValues.PlcValueType, int, error) {
	size, err := d.Size()
	if err != nil {
		return apiValues.NULL, 0, err
	}
	if d.UseAsBits {
		switch d.DataType {
		case "Integer8", "Unsigned8", "Integer16", "Unsigned16", "Integer32", "Unsigned32", "Integer64", "Unsigned64":
			return apiValues.BOOL, size * 8, nil
		case "OctetString":
			switch size {
			case 1:
				return apiValues.BYTE, 1, nil
			case 2:
				return apiValues.WORD, 1, nil
			case 4:
				return apiValues.DWORD, 1, nil
			case 8:
				return apiValues.LWORD, 1, nil
			}
			return apiValues.BYTE, size, nil
		}
	}
	switch d.DataType {
	case "Integer8":
		return apiValues.SINT, 1, nil
	case "Integer16":
		return apiValues.INT, 1, nil
	case "Integer32":
		return apiValues.DINT, 1, nil
	case "Integer64":
		return apiValues.LINT, 1, nil
	case "Unsigned8":
		return apiValues.USINT, 1, nil
	case "Unsigned16":
		return apiValues.UINT, 1, nil
	case "Unsigned32":
		return apiValues.UDINT, 1, nil
	case "Unsigned64":
		return apiValues.ULINT, 1, nil
	case "Float32":
		return apiValues.REAL, 1, nil
	case "Float64":
		return apiValues.LREAL, 1, nil
	case "VisibleString":
		return apiValues.STRING, 1, nil
	}
	return apiValues.BYTE, size, nil
}

// Layout is a submodule's cyclic data sizes, in bytes.
type Layout struct {
	InputLength, OutputLength int
	IOPSLength, IOCSLength    int
}

// Layout adds up the I/O data. IOPS and IOCS default to 1 byte, as the GSDML schema says.
func (d IOData) Layout() (Layout, error) {
	layout := Layout{IOPSLength: 1, IOCSLength: 1}
	if d.IOPSLength != nil {
		layout.IOPSLength = *d.IOPSLength
	}
	if d.IOCSLength != nil {
		layout.IOCSLength = *d.IOCSLength
	}
	for _, set := range d.Inputs {
		for _, item := range set.DataItems {
			size, err := item.Size()
			if err != nil {
				return Layout{}, err
			}
			layout.InputLength += size
		}
	}
	for _, set := range d.Outputs {
		for _, item := range set.DataItems {
			size, err := item.Size()
			if err != nil {
				return Layout{}, err
			}
			layout.OutputLength += size
		}
	}
	return layout, nil
}
