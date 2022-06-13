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
	"github.com/rs/zerolog/log"
	"math"
	"regexp"
	"strings"
)

// AsciiBox is a string surrounded by an ascii border (and an optional name)
type AsciiBox struct {
	data           string
	asciiBoxWriter *asciiBoxWriter
}

// DebugAsciiBox set to true to get debug messages
var DebugAsciiBox bool

// AsciiBoxer is used to render something in a box
type AsciiBoxer interface {
	// Box where int param is the proposed width
	Box(string, int) AsciiBox
}

var AsciiBoxWriterDefault = NewAsciiBoxWriter()

var AsciiBoxWriterLight = NewAsciiBoxWriterWithCustomBorders(
	"╭",
	"╮",
	"┄",
	"┆",
	"╰",
	"╯",
)

type AsciiBoxWriter interface {
	BoxBox(name string, box AsciiBox, charWidth int) AsciiBox
	BoxString(name string, data string, charWidth int) AsciiBox
	AlignBoxes(asciiBoxes []AsciiBox, desiredWith int) AsciiBox
	BoxSideBySide(box1 AsciiBox, box2 AsciiBox) AsciiBox
	BoxBelowBox(box1 AsciiBox, box2 AsciiBox) AsciiBox
}

func NewAsciiBoxWriter() AsciiBoxWriter {
	return NewAsciiBoxWriterWithCustomBorders(
		"╔",
		"╗",
		"═",
		"║",
		"╚",
		"╝",
	)
}

func NewAsciiBoxWriterWithCustomBorders(upperLeftCorner string, upperRightCorner string, horizontalLine string, verticalLine string, lowerLeftCorner string, lowerRightCorner string) AsciiBoxWriter {
	return &asciiBoxWriter{
		upperLeftCorner:  upperLeftCorner,
		upperRightCorner: upperRightCorner,
		horizontalLine:   horizontalLine,
		verticalLine:     verticalLine,
		lowerLeftCorner:  lowerLeftCorner,
		lowerRightCorner: lowerRightCorner,
		newLine:          '\n',
		emptyPadding:     " ",
		// the name gets prefixed with a extra symbol for indent
		extraNameCharIndent: 1,
		borderWidth:         1,
		newLineCharWidth:    1,
		boxNameRegex:        regexp.MustCompile(`^` + upperLeftCorner + horizontalLine + `(?P<name>[\w /]+)` + horizontalLine + `*` + upperRightCorner),
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type asciiBoxWriter struct {
	upperLeftCorner  string
	upperRightCorner string
	horizontalLine   string
	verticalLine     string
	lowerLeftCorner  string
	lowerRightCorner string
	newLine          rune
	emptyPadding     string
	// the name gets prefixed with a extra symbol for indent
	extraNameCharIndent int
	borderWidth         int
	newLineCharWidth    int
	boxNameRegex        *regexp.Regexp
}

func (a *asciiBoxWriter) boxString(name string, data string, charWidth int) AsciiBox {
	rawBox := AsciiBox{data, a}
	longestLine := rawBox.Width()
	if charWidth < longestLine {
		if DebugAsciiBox {
			log.Debug().Msgf("Overflow by %d chars", longestLine-charWidth)
		}
		charWidth = longestLine + a.borderWidth + a.borderWidth
	}
	var boxedString strings.Builder
	boxedString.Grow((a.borderWidth + longestLine + a.borderWidth + a.newLineCharWidth) * rawBox.Height())
	namePadding := int(math.Max(float64(charWidth-countChars(name)-a.borderWidth-a.extraNameCharIndent-a.borderWidth), 0))
	boxedString.WriteString(a.upperLeftCorner + a.horizontalLine + name + strings.Repeat(a.horizontalLine, namePadding) + a.upperRightCorner)
	boxedString.WriteRune(a.newLine)
	// Name of the header stretches the box so we align to that
	charWidth = a.borderWidth + a.extraNameCharIndent + countChars(name) + namePadding + a.borderWidth
	for _, line := range rawBox.Lines() {
		linePadding := float64(charWidth - boxLineOverheat - countChars(line))
		if linePadding < 0 {
			linePadding = 0
		}
		frontPadding := math.Floor(linePadding / 2.0)
		backPadding := math.Ceil(linePadding / 2.0)
		boxedString.WriteString(a.verticalLine + strings.Repeat(a.emptyPadding, int(frontPadding)) + line + strings.Repeat(a.emptyPadding, int(backPadding)) + a.verticalLine)
		boxedString.WriteRune(a.newLine)
	}
	bottomPadding := namePadding + countChars(name) + a.extraNameCharIndent
	boxedString.WriteString(a.lowerLeftCorner + strings.Repeat(a.horizontalLine, bottomPadding) + a.lowerRightCorner)
	return AsciiBox{boxedString.String(), a}
}

func (a *asciiBoxWriter) getBoxName(box AsciiBox) string {
	subMatch := a.boxNameRegex.FindStringSubmatch(box.String())
	if subMatch == nil {
		return ""
	}
	if len(subMatch) != 2 {
		panic("should never occur as we only have one named group")
	}
	return subMatch[1]
}

func (a *asciiBoxWriter) changeBoxName(box AsciiBox, newName string) AsciiBox {
	if !a.hasBorders(box) {
		return a.boxString(newName, box.String(), 0)
	}
	minimumWidthWithNewName := countChars(a.upperLeftCorner + a.horizontalLine + newName + a.upperRightCorner)
	nameLengthDifference := minimumWidthWithNewName - (a.unwrap(box).Width() + a.borderWidth + a.borderWidth)
	return a.BoxString(newName, a.unwrap(box).String(), box.Width()+nameLengthDifference)
}

func (a *asciiBoxWriter) mergeHorizontal(boxes []AsciiBox) AsciiBox {
	switch len(boxes) {
	case 0:
		return AsciiBox{"", a}
	case 1:
		return boxes[0]
	case 2:
		return a.BoxSideBySide(boxes[0], boxes[1])
	default:
		return a.BoxSideBySide(boxes[0], a.mergeHorizontal(boxes[1:]))
	}
}

func (a *asciiBoxWriter) expandBox(box AsciiBox, desiredWidth int) AsciiBox {
	if box.Width() >= desiredWidth {
		return box
	}
	boxLines := box.Lines()
	numberOfLine := len(boxLines)
	boxWidth := box.Width()
	padding := strings.Repeat(" ", desiredWidth-boxWidth)
	var newBox strings.Builder
	newBox.Grow((boxWidth + a.newLineCharWidth) * numberOfLine)
	for i, line := range boxLines {
		newBox.WriteString(line)
		newBox.WriteString(padding)
		if i < numberOfLine-1 {
			newBox.WriteRune(a.newLine)
		}
	}
	return AsciiBox{newBox.String(), a}
}

func (a *asciiBoxWriter) unwrap(box AsciiBox) AsciiBox {
	if !a.hasBorders(box) {
		return box
	}
	originalLines := box.Lines()
	newLines := make([]string, len(originalLines)-2)
	for i, line := range originalLines {
		if i == 0 {
			// we ignore the first line
			continue
		}
		if i == len(originalLines)-1 {
			// we ignore the last line
			break
		}
		runes := []rune(line)
		// Strip the vertical Lines and trim the padding
		unwrappedLine := string(runes[1 : len(runes)-1])
		if !strings.ContainsAny(unwrappedLine, a.verticalLine+a.horizontalLine) {
			// only trim boxes witch don't contain other boxes
			unwrappedLine = strings.Trim(unwrappedLine, a.emptyPadding)
		}
		newLines[i-1] = unwrappedLine
	}
	return AsciiBox{strings.Join(newLines, string(a.newLine)), a}
}

func (a *asciiBoxWriter) hasBorders(box AsciiBox) bool {
	if len(box.String()) == 0 {
		return false
	}
	// Check if the first char is the upper left corner
	return []rune(box.String())[0] == []rune(a.upperLeftCorner)[0]
}

func countChars(s string) int {
	return len([]rune(s))
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

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
	return strings.Split(m.data, "\n")
}

func (m AsciiBox) GetBoxName() string {
	return m.asciiBoxWriter.getBoxName(m)
}

func (m AsciiBox) ChangeBoxName(newName string) AsciiBox {
	return m.asciiBoxWriter.changeBoxName(m, newName)
}

func (m AsciiBox) IsEmpty() bool {
	if m.asciiBoxWriter.hasBorders(m) {
		return m.asciiBoxWriter.unwrap(m).String() == ""
	}
	return m.String() == ""
}

// String returns the string of the box
func (m AsciiBox) String() string {
	return m.data
}

// BoxBox boxes a box
func (a *asciiBoxWriter) BoxBox(name string, box AsciiBox, charWidth int) AsciiBox {
	return a.BoxString(name, box.data, charWidth)
}

// BoxString boxes a newline separated string into a beautiful box
func (a *asciiBoxWriter) BoxString(name string, data string, charWidth int) AsciiBox {
	return a.boxString(name, data, charWidth)
}

// AlignBoxes aligns all boxes to a desiredWidth and orders them from left to right and top to bottom (size will be at min the size of the biggest box)
func (a *asciiBoxWriter) AlignBoxes(boxes []AsciiBox, desiredWidth int) AsciiBox {
	if len(boxes) == 0 {
		return AsciiBox{"", a}
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
	bigBox := AsciiBox{"", a}
	currentBoxRow := make([]AsciiBox, 0)
	currentRowLength := 0
	for _, box := range boxes {
		currentRowLength += box.Width()
		if currentRowLength > actualWidth {
			mergedBoxes := a.mergeHorizontal(currentBoxRow)
			if bigBox.IsEmpty() {
				bigBox = mergedBoxes
			} else {
				bigBox = a.BoxBelowBox(bigBox, mergedBoxes)
			}
			currentRowLength = box.Width()
			currentBoxRow = make([]AsciiBox, 0)
		}
		currentBoxRow = append(currentBoxRow, box)
	}
	if len(currentBoxRow) > 0 {
		// Special case where all boxes fit into one row
		mergedBoxes := a.mergeHorizontal(currentBoxRow)
		if bigBox.IsEmpty() {
			bigBox = mergedBoxes
		} else {
			bigBox = a.BoxBelowBox(bigBox, mergedBoxes)
		}
	}
	return bigBox
}

// BoxSideBySide renders two boxes side by side
func (a *asciiBoxWriter) BoxSideBySide(box1, box2 AsciiBox) AsciiBox {
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
	return AsciiBox{aggregateBox.String(), a}
}

// BoxBelowBox renders two boxes below
func (a *asciiBoxWriter) BoxBelowBox(box1, box2 AsciiBox) AsciiBox {
	box1Width := box1.Width()
	box2Width := box2.Width()
	if box1Width < box2Width {
		box1 = a.expandBox(box1, box2Width)
	} else if box2Width < box1Width {
		box2 = a.expandBox(box2, box1Width)
	}
	return AsciiBox{box1.String() + "\n" + box2.String(), a}
}
