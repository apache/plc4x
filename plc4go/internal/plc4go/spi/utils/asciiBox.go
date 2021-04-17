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
	"math"
	"reflect"
	"strings"
)

// AsciiBox is a string surrounded by a ascii border (and a optional name)
type AsciiBox string

// DebugAsciiBox set to true to get debug messages
var DebugAsciiBox bool

// Width returns the width of the box without the newlines
func (m AsciiBox) Width() int {
	maxWidth := 0
	for _, line := range strings.Split(string(m), "\n") {
		currentLength := countChars(line)
		if maxWidth < currentLength {
			maxWidth = currentLength
		}
	}
	return maxWidth
}

// AsciiBoxer is used to render something in a box
type AsciiBoxer interface {
	// Box where int param is the proposed width
	Box(string, int) AsciiBox
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
		return BoxString(name, fmt.Sprintf("%t", anything), 0)
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64:
		// TODO: include hex later with this line
		//return BoxString(name, fmt.Sprintf("%#0*x %d", unsafe.Sizeof(anything)/2, anything, anything), 0)
		return BoxString(name, fmt.Sprintf("%d", anything), 0)
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
			// TODO: include hex here somehow. Seems that %x does print strange hex values here
			return BoxString(name, fmt.Sprintf("%d", anything), 0)
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

// BoxBox boxes a box
func BoxBox(name string, box AsciiBox, charWidth int) AsciiBox {
	return BoxString(name, string(box), charWidth)
}

// BoxString boxes a newline separated string into a beautiful box
func BoxString(name string, data string, charWidth int) AsciiBox {
	longestLine := AsciiBox(data).Width()
	if charWidth < longestLine {
		if DebugAsciiBox {
			log.Debug().Msgf("Overflow by %d chars", longestLine-charWidth)
		}
		charWidth = longestLine + 2
	}
	boxedString := ""
	namePadding := int(math.Max(float64(charWidth-countChars(name)-2-1), 0))
	bottomPadding := namePadding + countChars(name) + 1
	boxedString += "╔═" + name + strings.Repeat("═", namePadding) + "╗\n"
	// Name of the header stretches the box so we align to that
	charWidth = 2 + countChars(name) + namePadding + 1
	for _, line := range strings.Split(data, "\n") {
		linePadding := float64(charWidth - boxLineOverheat - countChars(line))
		if linePadding < 0 {
			linePadding = 0
		}
		frontPadding := math.Floor(linePadding / 2.0)
		backPadding := math.Ceil(linePadding / 2.0)
		boxedString += "║" + strings.Repeat(" ", int(frontPadding)) + line + strings.Repeat(" ", int(backPadding)) + "║\n"
	}
	boxedString += "╚" + strings.Repeat("═", bottomPadding) + "╝"
	return AsciiBox(boxedString)
}

func AlignBoxes(boxes []AsciiBox, desiredWidth int) AsciiBox {
	if len(boxes) == 0 {
		return boxes[0]
	}
	actualWidth := desiredWidth
	for _, box := range boxes {
		boxWidth := box.Width()
		if boxWidth > desiredWidth {
			if DebugAsciiBox {
				log.Debug().Msgf("Overflow by %d chars", boxWidth-desiredWidth)
			}
			actualWidth = boxWidth
		}
	}
	if DebugAsciiBox {
		log.Debug().Msgf("Working with %d chars", actualWidth)
	}
	bigBox := AsciiBox("")
	currentBoxRow := make([]AsciiBox, 0)
	currentRowLength := 0
	for _, box := range boxes {
		currentRowLength += box.Width()
		if currentRowLength > actualWidth {
			mergedBoxes := mergeHorizontal(currentBoxRow)
			if bigBox == "" {
				bigBox = mergedBoxes
			} else {
				bigBox = BoxBelowBox(bigBox, mergedBoxes)
			}
			currentRowLength = box.Width()
			currentBoxRow = make([]AsciiBox, 0)
		}
		currentBoxRow = append(currentBoxRow, box)
	}
	if len(currentBoxRow) > 0 {
		// Special case where all boxes fit into one row
		mergedBoxes := mergeHorizontal(currentBoxRow)
		if bigBox == "" {
			bigBox = mergedBoxes
		} else {
			bigBox = BoxBelowBox(bigBox, mergedBoxes)
		}
	}
	return bigBox
}

func mergeHorizontal(boxes []AsciiBox) AsciiBox {
	switch len(boxes) {
	case 1:
		return boxes[0]
	case 2:
		return BoxSideBySide(boxes[0], boxes[1])
	default:
		return BoxSideBySide(boxes[0], mergeHorizontal(boxes[1:]))
	}
}

func BoxSideBySide(box1, box2 AsciiBox) AsciiBox {
	aggregateBox := ""
	box1Width := box1.Width()
	box1Lines := strings.Split(string(box1), "\n")
	box2Width := box2.Width()
	box2Lines := strings.Split(string(box2), "\n")
	maxRows := int(math.Max(float64(len(box1Lines)), float64(len(box2Lines))))
	for row := 0; row < maxRows; row++ {
		ranOutOfLines := false
		if row >= len(box1Lines) {
			ranOutOfLines = true
			aggregateBox += strings.Repeat(" ", box1Width)
		} else {
			split1Row := box1Lines[row]
			padding := box1Width - countChars(split1Row)
			aggregateBox += split1Row + strings.Repeat(" ", padding)
		}
		if row >= len(box2Lines) {
			if ranOutOfLines {
				break
			}
			aggregateBox += strings.Repeat(" ", box2Width)
		} else {
			split2Row := box2Lines[row]
			padding := box2Width - countChars(split2Row)
			aggregateBox += split2Row + strings.Repeat(" ", padding)
		}
		aggregateBox += "\n"
	}
	return AsciiBox(aggregateBox[:len(aggregateBox)-1])
}

func BoxBelowBox(box1, box2 AsciiBox) AsciiBox {
	box1Width := box1.Width()
	box2Width := box2.Width()
	if box1Width < box2Width {
		box1 = expandBox(box1, box2Width)
	} else if box2Width < box1Width {
		box2 = expandBox(box2, box1Width)
	}
	return AsciiBox(string(box1) + "\n" + string(box2))
}

func expandBox(box AsciiBox, desiredWidth int) AsciiBox {
	if box.Width() >= desiredWidth {
		return box
	}
	// TODO: should we expand the borders?
	boxLines := strings.Split(string(box), "\n")
	newBox := ""
	for _, line := range boxLines {
		padding := desiredWidth - countChars(line)
		newBox += line + strings.Repeat(" ", padding) + "\n"
	}
	return AsciiBox(newBox[:len(newBox)-1])
}

func countChars(s string) int {
	return len([]rune(s))
}
