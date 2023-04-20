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
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

// BoxedDump dumps a 56+2 char wide hex string
func BoxedDump(name string, data []byte) AsciiBox {
	return AsciiBoxWriterDefault.BoxString(name, DumpFixedWidth(data, DefaultWidth), DefaultWidth+boxLineOverheat)
}

// BoxedDumpFixedWidth dumps a hex into a beautiful box
func BoxedDumpFixedWidth(name string, data []byte, charWidth int) AsciiBox {
	// we substract the 2 lines at the side
	dumpWidth := charWidth - 1 - 1
	return AsciiBoxWriterDefault.BoxString(name, DumpFixedWidth(data, dumpWidth), charWidth)
}

// BoxedDumpAnything dumps anything as hex into a beautiful box
func BoxedDumpAnything(name string, anything any) AsciiBox {
	return AsciiBoxWriterDefault.BoxString(name, DumpAnything(anything), 0)
}

// BoxedDumpAnythingFixedWidth dumps anything as hex into a beautiful box with a given width
func BoxedDumpAnythingFixedWidth(name string, anything any, charWidth int) AsciiBox {
	return AsciiBoxWriterDefault.BoxString(name, DumpAnythingFixedWidth(anything, charWidth), 0)
}

// DumpAnything dumps anything as hex
func DumpAnything(anything any) string {
	return DumpAnythingFixedWidth(anything, DefaultWidth)
}

// DumpAnythingFixedWidth dumps anything as hex
func DumpAnythingFixedWidth(anything any, charWidth int) string {
	convertedBytes, err := toBytes(anything)
	if err != nil {
		if DebugHex {
			log.Error().Err(err).Msg("Error converting to bytes")
		}
		return "<undumpable>"
	}
	return DumpFixedWidth(convertedBytes, charWidth)
}

func BoxAnything(name string, anything any, charWidth int) AsciiBox {
	switch anything.(type) {
	case nil:
		return AsciiBox{asciiBoxWriter: AsciiBoxWriterDefault.(*asciiBoxWriter)}
	case AsciiBoxer:
		// A box usually has its own name
		return anything.(AsciiBoxer).Box(name, charWidth)
	case bool:
		asInt := 0
		if anything.(bool) {
			asInt = 1
		}
		return AsciiBoxWriterDefault.BoxString(name, fmt.Sprintf("b%d %t", asInt, anything), 0)
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		hexDigits := reflect.TypeOf(anything).Bits() / 4
		return AsciiBoxWriterDefault.BoxString(name, fmt.Sprintf("%#0*x %d", hexDigits, anything, anything), 0)
	case []byte:
		//return AsciiBox{DumpFixedWidth(anything.([]byte), charWidth), AsciiBoxWriterDefault.(*asciiBoxWriter), AsciiBoxWriterDefault.(*asciiBoxWriter).compressBoxSet()}
		return AsciiBoxWriterDefault.BoxString(name, DumpFixedWidth(anything.([]byte), charWidth), charWidth)
	case string:
		return AsciiBoxWriterDefault.BoxString(name, anything.(string), charWidth)
	case fmt.Stringer:
		return AsciiBoxWriterDefault.BoxString(name, anything.(fmt.Stringer).String(), 0)
	default:
		valueOf := reflect.ValueOf(anything)
		switch valueOf.Kind() {
		case reflect.Slice, reflect.Array:
			boxes := make([]AsciiBox, valueOf.Len())
			for i := 0; i < valueOf.Len(); i++ {
				index := valueOf.Index(i)
				boxes[i] = BoxAnything("", index.Interface(), charWidth-2)
			}
			return AsciiBoxWriterDefault.BoxBox(name, AsciiBoxWriterDefault.AlignBoxes(boxes, charWidth), 0)
		case reflect.Ptr, reflect.Uintptr:
			return BoxAnything(name, valueOf.Elem().Interface(), charWidth)
		default:
			return AsciiBoxWriterDefault.BoxString(name, fmt.Sprintf("%v", anything), charWidth)
		}
	}
}
