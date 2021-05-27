//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
	"unsafe"
)

// BoxedDump dumps a 56+2 char wide hex string
func BoxedDump(name string, data []byte) AsciiBox {
	return BoxString(name, DumpFixedWidth(data, DefaultWidth), DefaultWidth+boxLineOverheat)
}

// BoxedDumpFixedWidth dumps a hex into a beautiful box
func BoxedDumpFixedWidth(name string, data []byte, charWidth int) AsciiBox {
	// we substract the 2 lines at the side
	dumpWidth := charWidth - 1 - 1
	return BoxString(name, DumpFixedWidth(data, dumpWidth), charWidth)
}

// BoxedDumpAnything dumps anything as hex into a beautiful box
func BoxedDumpAnything(name string, anything interface{}) AsciiBox {
	return BoxString(name, DumpAnything(anything), 0)
}

// BoxedDumpAnythingFixedWidth dumps anything as hex into a beautiful box with a given width
func BoxedDumpAnythingFixedWidth(name string, anything interface{}, charWidth int) AsciiBox {
	return BoxString(name, DumpAnythingFixedWidth(anything, charWidth), 0)
}

// DumpAnything dumps anything as hex
func DumpAnything(anything interface{}) string {
	return DumpAnythingFixedWidth(anything, DefaultWidth)
}

// DumpAnythingFixedWidth dumps anything as hex
func DumpAnythingFixedWidth(anything interface{}, charWidth int) string {
	convertedBytes, err := toBytes(anything)
	if err != nil {
		if DebugHex {
			log.Error().Err(err).Msg("Error converting to bytes")
		}
		return "<undumpable>"
	}
	return DumpFixedWidth(convertedBytes, charWidth)
}

func BoxAnything(name string, anything interface{}, charWidth int) AsciiBox {
	switch anything.(type) {
	case nil:
		return ""
	case AsciiBoxer:
		if reflect.ValueOf(anything).IsNil() {
			return ""
		}
		// A box usually has its own name
		return anything.(AsciiBoxer).Box(name, charWidth)
	case bool:
		asInt := 0
		if anything.(bool) {
			asInt = 1
		}
		return BoxString(name, fmt.Sprintf("b%d %t", asInt, anything), 0)
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		hexDigits := reflect.TypeOf(anything).Bits() / 4
		return BoxString(name, fmt.Sprintf("%#0*x %d", hexDigits, anything, anything), 0)
	case []byte:
		return AsciiBox(DumpFixedWidth(anything.([]byte), charWidth))
	case string:
		return BoxString(name, anything.(string), charWidth)
	case fmt.Stringer:
		return BoxString(name, anything.(fmt.Stringer).String(), 0)
	default:
		valueOf := reflect.ValueOf(anything)
		if valueOf.IsNil() {
			return ""
		}
		switch valueOf.Kind() {
		case reflect.Bool:
			return BoxString(name, fmt.Sprintf("%t", anything), 0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			return BoxString(name, fmt.Sprintf("%#0*x %d", unsafe.Sizeof(valueOf.Elem())/2, valueOf.Elem(), anything), 0)
		case reflect.Slice, reflect.Array:
			boxes := make([]AsciiBox, valueOf.Len())
			for i := 0; i < valueOf.Len(); i++ {
				index := valueOf.Index(i)
				boxes[i] = BoxAnything("", index.Interface(), charWidth-2)
			}
			return BoxBox(name, AlignBoxes(boxes, charWidth), 0)
		case reflect.Ptr, reflect.Uintptr:
			return BoxAnything(name, valueOf.Elem().Interface(), charWidth)
		default:
			return BoxString(name, fmt.Sprintf("0x%x", anything.(interface{})), charWidth)
		}
	}
}
