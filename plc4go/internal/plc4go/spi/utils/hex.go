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
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// DefaultWidth defaults to a default screen dumps size
const DefaultWidth = 51

// boxLineOverheat Overheat per line when drawing boxes
const boxLineOverheat = 1 + 1

// Dump dumps a 46 char wide hex string
func BoxedDump(name string, data []byte) string {
	// we substract the 2 lines at the side
	dumpWidth := DefaultWidth - boxLineOverheat
	return string(BoxString(name, DumpFixedWidth(data, dumpWidth), DefaultWidth))
}

// Dump dumps a 46 char wide hex string
func Dump(data []byte) string {
	return DumpFixedWidth(data, DefaultWidth)
}

// BoxedDumpFixedWidth dumps a hex into a beautiful box
func BoxedDumpFixedWidth(name string, data []byte, charWidth int) string {
	// we substract the 2 lines at the side
	dumpWidth := charWidth - 1 - 1
	return string(BoxString(name, DumpFixedWidth(data, dumpWidth), charWidth))
}

// DumpAnything dumps anything as hex
func DumpAnything(anything interface{}) string {
	convertedBytes, err := toBytes(anything)
	if err != nil {
		panic(err)
	}
	return DumpFixedWidth(convertedBytes, 1)
}

// TODO: test with charWidht <12
// DumpFixedWidth dumps hex as hex string. Min width of string returned is 18 up to supplied charWidth
func DumpFixedWidth(data []byte, charWidth int) string {
	if charWidth <= 0 {
		panic("charWidth needs to be greater than 0")
	}
	hexString := ""
	// 3 digits index plus one blank
	const indexWidth = 3 + 1
	// 2 hex digits + 2 blanks
	const byteWidth = 2 + 2
	// strings get quoate by 2 chars
	const stringRenderOverheat = 2
	const minWidth = indexWidth + byteWidth + stringRenderOverheat + 1
	if charWidth < minWidth {
		charWidth = minWidth + 6
	}
	// Formulary to calculate max bytes per row...
	maxBytesPerRow := ((charWidth - indexWidth - stringRenderOverheat) / (byteWidth + 1)) - 1

	for byteIndex, rowIndex := 0, 0; byteIndex < len(data); byteIndex, rowIndex = byteIndex+maxBytesPerRow, rowIndex+1 {
		hexString += fmt.Sprintf("%03d 0x: ", byteIndex)
		for columnIndex := 0; columnIndex < maxBytesPerRow; columnIndex++ {
			absoluteIndex := byteIndex + columnIndex
			if absoluteIndex < len(data) {
				hexString += fmt.Sprintf("%02x  ", data[absoluteIndex])
			} else {
				// align with empty byte representation
				hexString += strings.Repeat(" ", byteWidth)
			}
		}
		endIndex := byteIndex + maxBytesPerRow
		if endIndex >= len(data) {
			endIndex = len(data)
		}
		stringRepresentation := maskString(data[byteIndex:endIndex])
		if len([]rune(stringRepresentation)) < maxBytesPerRow {
			stringRepresentation += strings.Repeat(" ", (maxBytesPerRow-len([]rune(stringRepresentation)))%maxBytesPerRow)
		}
		hexString += fmt.Sprintf("'%s'\n", stringRepresentation)
	}
	// remove last newline
	return hexString[:len(hexString)-1]
}

func maskString(data []byte) string {
	for i := range data {
		switch {
		case data[i] < 32:
			fallthrough
		case data[i] > 126:
			data[i] = '.'
		}
	}
	return string(data)
}

func toBytes(anything interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(anything)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding datatype")
	}
	return buffer.Bytes(), nil
}
