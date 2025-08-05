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
	"reflect"

	"github.com/rs/zerolog/log"
)

// BoxedDump dumps a 56+2 char wide hex string
func BoxedDump(data []byte, options ...func(*BoxOptions)) AsciiBox {
	var opts BoxOptions
	opts.BoxSet = DefaultBoxSet()
	for _, opt := range options {
		opt(&opts)
	}
	if opts.CharWidth <= 0 {
		opts.CharWidth = DefaultWidth + boxLineOverheat
	}
	// we substract the 2 lines at the side
	dumpWidth := opts.CharWidth - 1 - 1
	return AsciiBoxWriterDefault.BoxString(DumpFixedWidth(data, dumpWidth), WithAsciiBoxOptions(opts))
}

// BoxedDumpAnything dumps anything as hex into a beautiful box
func BoxedDumpAnything(anything any, options ...func(*BoxOptions)) AsciiBox {
	return AsciiBoxWriterDefault.BoxString(DumpAnything(anything, options...), options...)
}

// DumpAnything dumps anything as hex
func DumpAnything(anything any, options ...func(*BoxOptions)) string {
	var opts BoxOptions
	for _, opt := range options {
		opt(&opts)
	}
	if opts.CharWidth <= 0 {
		opts.CharWidth = DefaultWidth
	}
	convertedBytes, err := toBytes(anything)
	if err != nil {
		if DebugHex {
			log.Error().Err(err).Msg("Error converting to bytes")
		}
		return "<undumpable>"
	}
	return DumpFixedWidth(convertedBytes, opts.CharWidth)
}

func BoxAnything(anything any, options ...func(*BoxOptions)) AsciiBox {
	var opts BoxOptions
	for _, opt := range options {
		opt(&opts)
	}
	switch anything.(type) {
	case nil:
		return AsciiBox{asciiBoxWriter: AsciiBoxWriterDefault.(*asciiBoxWriter)}
	case AsciiBoxer:
		return anything.(AsciiBoxer).Box(options...)
	case bool:
		asInt := 0
		if anything.(bool) {
			asInt = 1
		}
		return AsciiBoxWriterDefault.BoxString(fmt.Sprintf("b%d %t", asInt, anything), options...)
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		hexDigits := reflect.TypeOf(anything).Bits() / 4
		return AsciiBoxWriterDefault.BoxString(fmt.Sprintf("%#0*x %d", hexDigits, anything, anything), options...)
	case []byte:
		//return AsciiBox{DumpFixedWidth(anything.([]byte), charWidth), AsciiBoxWriterDefault.(*asciiBoxWriter), AsciiBoxWriterDefault.(*asciiBoxWriter).compressBoxSet()}
		return AsciiBoxWriterDefault.BoxString(DumpFixedWidth(anything.([]byte), opts.CharWidth), options...)
	case string:
		return AsciiBoxWriterDefault.BoxString(anything.(string), options...)
	case fmt.Stringer:
		return AsciiBoxWriterDefault.BoxString(anything.(fmt.Stringer).String(), options...)
	default:
		valueOf := reflect.ValueOf(anything)
		switch valueOf.Kind() {
		case reflect.Slice, reflect.Array:
			boxes := make([]AsciiBox, valueOf.Len())
			for i := 0; i < valueOf.Len(); i++ {
				index := valueOf.Index(i)
				boxes[i] = BoxAnything(index.Interface(), WithAsciiBoxCharWidth(opts.CharWidth-2))
			}
			return AsciiBoxWriterDefault.BoxBox(AsciiBoxWriterDefault.AlignBoxes(boxes, opts.CharWidth), options...)
		case reflect.Ptr, reflect.Uintptr:
			return BoxAnything(valueOf.Elem().Interface(), options...)
		default:
			return AsciiBoxWriterDefault.BoxString(fmt.Sprintf("%v", anything), options...)
		}
	}
}
