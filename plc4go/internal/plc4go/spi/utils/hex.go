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
	"github.com/rs/zerolog/log"
	"math"
	"strings"
)

// DefaultWidth defaults to a default screen dumps size
const DefaultWidth = 46 // 10 bytes per line on a []byte < 999

// boxLineOverheat Overheat per line when drawing boxes
const boxLineOverheat = 1 + 1

// blankWidth blank size of blank
const blankWidth = 1

// byteWidth required size of runes required to print one bytes 2 hex digits + 1 blanks
const byteWidth = 2 + 1

// pipeWidth size of the pipe char
const pipeWidth = 1

// DebugHex set to true to get debug messages
var DebugHex bool

// BoxedDump dumps a 56+2 char wide hex string
func BoxedDump(name string, data []byte) string {
	return string(BoxString(name, DumpFixedWidth(data, DefaultWidth), DefaultWidth+boxLineOverheat))
}

// Dump dumps a 56 char wide hex string
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
		if DebugHex {
			log.Error().Err(err).Msg("Error converting to bytes")
		}
		return "<undumpable>"
	}
	return Dump(convertedBytes)
}

// DumpFixedWidth dumps hex as hex string. Min width of string returned is 18 up to supplied charWidth
func DumpFixedWidth(data []byte, desiredCharWidth int) string {
	if data == nil || len(data) < 1 {
		return ""
	}
	hexString := ""
	maxBytesPerRow, indexWidth := calculateBytesPerRowAndIndexWidth(len(data), desiredCharWidth)

	for byteIndex, rowIndex := 0, 0; byteIndex < len(data); byteIndex, rowIndex = byteIndex+maxBytesPerRow, rowIndex+1 {
		indexString := fmt.Sprintf("%0*d|", indexWidth, byteIndex)
		hexString += indexString
		for columnIndex := 0; columnIndex < maxBytesPerRow; columnIndex++ {
			absoluteIndex := byteIndex + columnIndex
			if absoluteIndex < len(data) {
				hexString += fmt.Sprintf("%02x ", data[absoluteIndex])
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

func calculateBytesPerRowAndIndexWidth(numberOfBytes, desiredStringWidth int) (int, int) {
	if DebugHex {
		log.Debug().Msgf("Calculating max row and index for %d number of bytes and a desired string width of %d", numberOfBytes, desiredStringWidth)
	}
	indexDigits := int(math.Log10(float64(numberOfBytes))) + 1
	requiredIndexWidth := indexDigits + pipeWidth
	if DebugHex {
		log.Debug().Msgf("index width %d for indexDigits %d for bytes %d", requiredIndexWidth, indexDigits, numberOfBytes)
	}
	// strings get quoted by 2 chars
	const quoteRune = 1
	const potentialStringRenderRune = 1
	// 0 00  '.'
	availableSpace := requiredIndexWidth + byteWidth + quoteRune + potentialStringRenderRune + quoteRune
	if DebugHex {
		log.Debug().Msgf("calculated %d minimal width for number of bytes %d", availableSpace, numberOfBytes)
	}
	if desiredStringWidth >= availableSpace {
		availableSpace = desiredStringWidth
	} else {
		if DebugHex {
			log.Debug().Msgf("Overflow by %d runes", desiredStringWidth-availableSpace)
		}
	}
	if DebugHex {
		log.Debug().Msgf("Actual space %d", availableSpace)
	}

	z := float64(availableSpace)
	y := float64(requiredIndexWidth)
	a := float64(byteWidth)
	b := float64(quoteRune)
	// c = needed space for bytes x * byteWidth
	// x = maxBytesPerRow
	// x = (z - (y + b + x * 1 + b)) / a == x = (-2 * b - y + z)/(a + 1) and a + 1!=0 and a!=0
	x := ((-2 * b) - y + z) / (a + 1)
	if DebugHex {
		log.Debug().Msgf("Calculated number of bytes per row %f in int %d", x, int(x))
	}
	return int(x), indexDigits
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
