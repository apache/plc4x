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
	"github.com/rs/zerolog/log"
	"math"
	"strings"
)

// AsciiBox is a string surrounded by a ascii border (and a optional name)
type AsciiBox string

// DebugAsciiBox set to true to get debug messages
var DebugAsciiBox bool

// AsciiBoxer is used to render something in a box
type AsciiBoxer interface {
	// Box where int param is the proposed width
	Box(string, int) AsciiBox
}

// Width returns the width of the box without the newlines
func (m AsciiBox) Width() int {
	maxWidth := 0
	for _, line := range m.Lines() {
		currentLength := countChars(line)
		if maxWidth < currentLength {
			maxWidth = currentLength
		}
	}
	return maxWidth
}

// Height returns the height of the box without
func (m AsciiBox) Height() int {
	return len(m.Lines())
}

// Lines returns the lines of the box
func (m AsciiBox) Lines() []string {
	return strings.Split(string(m), "\n")
}

// String returns the string of the box
func (m AsciiBox) String() string {
	return string(m)
}

// BoxBox boxes a box
func BoxBox(name string, box AsciiBox, charWidth int) AsciiBox {
	return BoxString(name, string(box), charWidth)
}

// BoxString boxes a newline separated string into a beautiful box
func BoxString(name string, data string, charWidth int) AsciiBox {
	// the name gets prefixed with a extra symbol for indent
	const extraNameCharWidth = 1
	const borderWidth = 1
	const newLineCharWidth = 1
	rawBox := AsciiBox(data)
	longestLine := rawBox.Width()
	if charWidth < longestLine {
		if DebugAsciiBox {
			log.Debug().Msgf("Overflow by %d chars", longestLine-charWidth)
		}
		charWidth = longestLine + borderWidth + borderWidth
	}
	var boxedString strings.Builder
	boxedString.Grow((borderWidth + longestLine + borderWidth + newLineCharWidth) * rawBox.Height())
	namePadding := int(math.Max(float64(charWidth-countChars(name)-borderWidth-extraNameCharWidth-borderWidth), 0))
	boxedString.WriteString("╔═" + name + strings.Repeat("═", namePadding) + "╗\n")
	// Name of the header stretches the box so we align to that
	charWidth = borderWidth + extraNameCharWidth + countChars(name) + namePadding + borderWidth
	for _, line := range rawBox.Lines() {
		linePadding := float64(charWidth - boxLineOverheat - countChars(line))
		if linePadding < 0 {
			linePadding = 0
		}
		frontPadding := math.Floor(linePadding / 2.0)
		backPadding := math.Ceil(linePadding / 2.0)
		boxedString.WriteString("║" + strings.Repeat(" ", int(frontPadding)) + line + strings.Repeat(" ", int(backPadding)) + "║\n")
	}
	bottomPadding := namePadding + countChars(name) + extraNameCharWidth
	boxedString.WriteString("╚" + strings.Repeat("═", bottomPadding) + "╝")
	return AsciiBox(boxedString.String())
}

// AlignBoxes aligns all boxes to a desiredWidth and orders them from left to right and top to bottom (size will be at min the size of the biggest box)
func AlignBoxes(boxes []AsciiBox, desiredWidth int) AsciiBox {
	if len(boxes) == 0 {
		return ""
	}
	actualWidth := desiredWidth
	for _, box := range boxes {
		boxWidth := box.Width()
		if boxWidth > actualWidth {
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
	case 0:
		return ""
	case 1:
		return boxes[0]
	case 2:
		return BoxSideBySide(boxes[0], boxes[1])
	default:
		return BoxSideBySide(boxes[0], mergeHorizontal(boxes[1:]))
	}
}

// BoxSideBySide renders two boxes side by side
func BoxSideBySide(box1, box2 AsciiBox) AsciiBox {
	const newLineCharWidth = 1
	var aggregateBox strings.Builder
	box1Width := box1.Width()
	box1Lines := box1.Lines()
	box2Width := box2.Width()
	box2Lines := box2.Lines()
	maxRows := int(math.Max(float64(len(box1Lines)), float64(len(box2Lines))))
	aggregateBox.Grow((box1Width + box2Width + newLineCharWidth) * maxRows)
	for row := 0; row < maxRows; row++ {
		ranOutOfLines := false
		if row >= len(box1Lines) {
			ranOutOfLines = true
			aggregateBox.WriteString(strings.Repeat(" ", box1Width))
		} else {
			split1Row := box1Lines[row]
			padding := box1Width - countChars(split1Row)
			aggregateBox.WriteString(split1Row + strings.Repeat(" ", padding))
		}
		if row >= len(box2Lines) {
			if ranOutOfLines {
				break
			}
			aggregateBox.WriteString(strings.Repeat(" ", box2Width))
		} else {
			split2Row := box2Lines[row]
			padding := box2Width - countChars(split2Row)
			aggregateBox.WriteString(split2Row + strings.Repeat(" ", padding))
		}
		if row < maxRows-1 {
			// Only write newline if we are not the last line
			aggregateBox.WriteRune('\n')
		}
	}
	return AsciiBox(aggregateBox.String())
}

// BoxBelowBox renders two boxes below
func BoxBelowBox(box1, box2 AsciiBox) AsciiBox {
	box1Width := box1.Width()
	box2Width := box2.Width()
	if box1Width < box2Width {
		box1 = expandBox(box1, box2Width)
	} else if box2Width < box1Width {
		box2 = expandBox(box2, box1Width)
	}
	return AsciiBox(box1.String() + "\n" + box2.String())
}

func expandBox(box AsciiBox, desiredWidth int) AsciiBox {
	const newLineCharWidth = 1
	if box.Width() >= desiredWidth {
		return box
	}
	// TODO: should we expand the borders?
	boxLines := box.Lines()
	numberOfLine := len(boxLines)
	boxWidth := box.Width()
	padding := strings.Repeat(" ", desiredWidth-boxWidth)
	var newBox strings.Builder
	newBox.Grow((boxWidth + newLineCharWidth) * numberOfLine)
	for i, line := range boxLines {
		newBox.WriteString(line)
		newBox.WriteString(padding)
		if i < numberOfLine-1 {
			newBox.WriteRune('\n')
		}
	}
	return AsciiBox(newBox.String())
}

func countChars(s string) int {
	return len([]rune(s))
}
