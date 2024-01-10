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

// Dump dumps a 56 char wide hex string
func Dump(data []byte, highlights ...int) string {
	return DumpFixedWidth(data, DefaultWidth, highlights...)
}

// DumpFixedWidth dumps hex as hex string. Min width of string returned is 18 up to supplied charWidth
func DumpFixedWidth(data []byte, desiredCharWidth int, highlights ...int) string {
	if data == nil || len(data) < 1 {
		return ""
	}
	// We copy the array to avoid mutations
	data = append(data[:0:0], data...)
	hexString := ""
	maxBytesPerRow, indexWidth := calculateBytesPerRowAndIndexWidth(len(data), desiredCharWidth)
	highlightsSet := map[int]struct{}{}
	for _, highlight := range highlights {
		highlightsSet[highlight] = struct{}{}
	}
	for byteIndex, rowIndex := 0, 0; byteIndex < len(data); byteIndex, rowIndex = byteIndex+maxBytesPerRow, rowIndex+1 {
		indexString := fmt.Sprintf("%0*d|", indexWidth, byteIndex)
		hexString += indexString
		for columnIndex := 0; columnIndex < maxBytesPerRow; columnIndex++ {
			absoluteIndex := byteIndex + columnIndex
			if absoluteIndex < len(data) {
				if _, ok := highlightsSet[absoluteIndex]; ok {
					hexString += "\033[0;31m"
				}
				hexString += fmt.Sprintf("%02x ", data[absoluteIndex])
				if _, ok := highlightsSet[absoluteIndex]; ok {
					hexString += "\033[0m"
				}
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

// DiffHex produces a hex diff AsciiBox of two byte arrays
func DiffHex(expectedBytes, actualBytes []byte) AsciiBox {
	numBytes := int(math.Min(float64(len(expectedBytes)), float64(len(actualBytes))))
	brokenAt := -1
	var diffIndexes []int
	for i := 0; i < numBytes; i++ {
		if expectedBytes[i] != actualBytes[i] {
			if brokenAt < 0 {
				brokenAt = i
			}
			diffIndexes = append(diffIndexes, i)
		}
	}
	expectedHex := DumpFixedWidth(expectedBytes, 46, diffIndexes...)
	actialHex := DumpFixedWidth(actualBytes, 46, diffIndexes...)
	return AsciiBoxWriterDefault.BoxSideBySide(AsciiBoxWriterDefault.BoxString("expected", expectedHex, 0), AsciiBoxWriterDefault.BoxString("actual", actialHex, 0))

}

func calculateBytesPerRowAndIndexWidth(numberOfBytes, desiredStringWidth int) (int, int) {
	if DebugHex {
		log.Debug().
			Int("numberOfBytes", numberOfBytes).
			Int("desiredStringWidth", desiredStringWidth).
			Msg("Calculating max row and index for numberOfBytes number of bytes and a desired string width of desiredStringWidth")
	}
	indexDigits := int(math.Log10(float64(numberOfBytes))) + 1
	requiredIndexWidth := indexDigits + pipeWidth
	if DebugHex {
		log.Debug().
			Int("requiredIndexWidth", requiredIndexWidth).
			Int("indexDigits", indexDigits).
			Int("numberOfBytes", numberOfBytes).
			Msg("index width requiredIndexWidth for indexDigits for numberOfBytes")
	}
	// strings get quoted by 2 chars
	const quoteRune = 1
	const potentialStringRenderRune = 1
	// 0 00  '.'
	availableSpace := requiredIndexWidth + byteWidth + quoteRune + potentialStringRenderRune + quoteRune
	if DebugHex {
		log.Debug().
			Int("availableSpace", availableSpace).
			Int("numberOfBytes", numberOfBytes).
			Msg("calculated availableSpace minimal width for number of bytes numberOfBytes")
	}
	if desiredStringWidth >= availableSpace {
		availableSpace = desiredStringWidth
	} else {
		if DebugHex {
			log.Debug().
				Int("n", desiredStringWidth-availableSpace).
				Msg("Overflow by n runes")
		}
	}
	if DebugHex {
		log.Debug().
			Int("availableSpace", availableSpace).
			Msg("Actual space")
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
		log.Debug().
			Float64("x", x).
			Int("xInt", int(x)).
			Msg("Calculated number of bytes per row x in int xInt")
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

func toBytes(anything any) ([]byte, error) {
	var buffer bytes.Buffer
	err := gob.NewEncoder(&buffer).Encode(anything)
	if err != nil {
		return nil, errors.Wrap(err, "error encoding datatype")
	}
	return buffer.Bytes(), nil
}
