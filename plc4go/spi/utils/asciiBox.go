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
	"math"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

var AsciiBoxLogger zerolog.Logger

// AsciiBox is a string surrounded by an ascii border (and an optional name)
type AsciiBox struct {
	data             string
	asciiBoxWriter   *asciiBoxWriter
	compressedBoxSet string
}

func WithAsciiBoxName(name string) func(*BoxOptions) {
	return func(opts *BoxOptions) {
		opts.Name = name
	}
}

func WithAsciiBoxHeader(header string) func(*BoxOptions) {
	return func(box *BoxOptions) {
		box.Header = header
	}
}

func WithAsciiBoxFooter(footer string) func(*BoxOptions) {
	return func(box *BoxOptions) {
		box.Footer = footer
	}
}

func WithAsciiBoxCharWidth(charWidth int) func(*BoxOptions) {
	return func(opts *BoxOptions) {
		opts.CharWidth = charWidth
	}
}

func WithAsciiBoxOptions(boxOptions BoxOptions) func(*BoxOptions) {
	return func(opts *BoxOptions) {
		*opts = boxOptions
	}
}

func WithAsciiBoxBoxSet(boxSet BoxSet) func(*BoxOptions) {
	return func(opts *BoxOptions) {
		opts.BoxSet = boxSet
	}
}

type BoxOptions struct {
	// The name of the box
	Name string
	// The additional header of the box appearing on the right upper side
	Header string
	// The additional footer of the box appearing on the right lower side
	Footer string
	// The desired CharWidth
	CharWidth int
	// The BoxSet used to print this box
	BoxSet BoxSet
}

type BoxSet struct {
	UpperLeftCorner  string
	UpperRightCorner string
	HorizontalLine   string
	VerticalLine     string
	LowerLeftCorner  string
	LowerRightCorner string
}

func DefaultBoxSet() BoxSet {
	return BoxSet{
		"╔",
		"╗",
		"═",
		"║",
		"╚",
		"╝",
	}
}

func DefaultLightBoxSet() BoxSet {
	return BoxSet{
		"╭",
		"╮",
		"┄",
		"┆",
		"╰",
		"╯",
	}
}

// DebugAsciiBox set to true to get debug messages
var DebugAsciiBox bool

// ANSI_PATTERN source: https://github.com/chalk/ansi-regex/blob/main/index.js#L3
var ANSI_PATTERN = regexp.MustCompile("[\u001b\u009b][\\[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]")

// AsciiBoxer is used to render something in a box
type AsciiBoxer interface {
	// Box with options
	Box(...func(*BoxOptions)) AsciiBox
}

var AsciiBoxWriterDefault = NewAsciiBoxWriter()

var AsciiBoxWriterLight = NewAsciiBoxWriter(WithAsciiBoxWriterDefaultBoxSet(DefaultLightBoxSet()))

type AsciiBoxWriter interface {
	BoxBox(box AsciiBox, options ...func(*BoxOptions)) AsciiBox
	BoxString(data string, options ...func(*BoxOptions)) AsciiBox
	AlignBoxes(asciiBoxes []AsciiBox, desiredWith int, options ...func(*BoxOptions)) AsciiBox
	BoxSideBySide(box1 AsciiBox, box2 AsciiBox, options ...func(*BoxOptions)) AsciiBox
	BoxBelowBox(box1 AsciiBox, box2 AsciiBox, options ...func(*BoxOptions)) AsciiBox
}

func NewAsciiBoxWriter(opts ...func(writer *asciiBoxWriter)) AsciiBoxWriter {
	return newAsciiBoxWriter(opts...)
}

func WithAsciiBoxWriterDefaultBoxSet(boxSet BoxSet) func(*asciiBoxWriter) {
	return func(a *asciiBoxWriter) {
		a.defaultBoxSet = boxSet
	}
}

func WithAsciiBoxWriterDefaultColoredBoxes(nameColor, headerColor, footerColor *color.Color) func(*asciiBoxWriter) {
	return func(a *asciiBoxWriter) {
		if nameColor != nil {
			a.namePrinter = nameColor.Sprint
		} else {
			a.namePrinter = fmt.Sprint
		}
		if headerColor != nil {
			a.headerPrinter = headerColor.Sprint
		} else {
			a.headerPrinter = fmt.Sprint
		}
		if footerColor != nil {
			a.footerPrinter = footerColor.Sprint
		} else {
			a.footerPrinter = fmt.Sprint
		}
	}
}

func WithAsciiBoxWriterDisableColoredBoxes() func(*asciiBoxWriter) {
	return func(a *asciiBoxWriter) {
		a.namePrinter = fmt.Sprint
		a.headerPrinter = fmt.Sprint
		a.footerPrinter = fmt.Sprint
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

func (b BoxSet) compressBoxSet() string {
	return b.UpperLeftCorner + b.UpperRightCorner + b.HorizontalLine + b.VerticalLine + b.LowerLeftCorner + b.LowerRightCorner
}

func (b BoxSet) contributeToCompressedBoxSet(box AsciiBox) string {
	actualSet := b.compressBoxSet()
	if box.compressedBoxSet == "" {
		// they have nothing to contribute
		return actualSet
	}
	if actualSet == "" {
		// I have nothing to contribute
		return box.compressedBoxSet
	}
	if strings.Contains(box.compressedBoxSet, actualSet) {
		// we have nothing to add
		return box.compressedBoxSet
	}
	return box.compressedBoxSet + "," + actualSet
}

func combineCompressedBoxSets(box1, box2 AsciiBox) string {
	allSets := make(map[string]any)
	for _, s := range strings.Split(box1.compressedBoxSet, ",") {
		allSets[s] = true
	}
	for _, s := range strings.Split(box2.compressedBoxSet, ",") {
		allSets[s] = true
	}
	var foundSets []string
	for set := range allSets {
		foundSets = append(foundSets, set)
	}
	return strings.Join(foundSets, ",")
}

type asciiBoxWriter struct {
	newLine      rune
	emptyPadding string
	// the name gets prefixed with an extra symbol for indent
	extraNameCharIndent int
	borderWidth         int
	newLineCharWidth    int
	defaultBoxSet       BoxSet
	boxHeaderRegex      *regexp.Regexp
	boxFooterRegex      *regexp.Regexp
	namePrinter         func(a ...any) string
	headerPrinter       func(a ...any) string
	footerPrinter       func(a ...any) string
}

var _ AsciiBoxWriter = (*asciiBoxWriter)(nil)

func newAsciiBoxWriter(opts ...func(writer *asciiBoxWriter)) *asciiBoxWriter {
	a := &asciiBoxWriter{
		newLine:      '\n',
		emptyPadding: " ",
		// the name gets prefixed with an extra symbol for indent
		extraNameCharIndent: 1,
		borderWidth:         1,
		newLineCharWidth:    1,
		defaultBoxSet:       DefaultBoxSet(),
		namePrinter:         color.New(color.FgGreen, color.Bold).Sprint,
		headerPrinter:       color.New(color.FgBlue).Sprint,
		footerPrinter:       color.New(color.FgRed, color.Italic).Sprint,
	}
	for _, opt := range opts {
		opt(a)
	}
	hl := a.defaultBoxSet.HorizontalLine
	a.boxHeaderRegex = regexp.MustCompile(`^` + a.defaultBoxSet.UpperLeftCorner + hl + `(?P<name>[^` + hl + `]+)` + hl + `*` + `(?P<header>[^` + hl + `]+)?` + hl + `*` + a.defaultBoxSet.UpperRightCorner)
	a.boxFooterRegex = regexp.MustCompile(`(?m)^` + a.defaultBoxSet.LowerLeftCorner + hl + `*` + `(?P<footer>[^` + hl + `]+)` + hl + `*` + a.defaultBoxSet.LowerRightCorner)
	return a
}

func (a *asciiBoxWriter) boxString(data string, options ...func(*BoxOptions)) AsciiBox {
	var opts BoxOptions
	opts.BoxSet = a.defaultBoxSet
	for _, opt := range options {
		opt(&opts)
	}
	name := opts.Name
	nameLength := countChars(name)
	if name != "" {
		name = a.namePrinter(name)
	}

	header := opts.Header
	if name != "" && header != "" {
		header = opts.BoxSet.HorizontalLine + opts.BoxSet.HorizontalLine + a.headerPrinter(header) + opts.BoxSet.HorizontalLine // Lazy manipulation to trick calculation below (adds a spacing between name and header)
	}
	headerLength := countChars(header)

	footer := opts.Footer
	if footer != "" {
		footer = a.footerPrinter(footer) + opts.BoxSet.HorizontalLine
	}
	footerLength := countChars(footer)

	charWidth := opts.CharWidth

	data = strings.ReplaceAll(data, "\r\n", "\n") // carriage return just messes with boxes
	data = strings.ReplaceAll(data, "\t", "  ")   // Tabs just don't work well as they distort the boxes so we convert them to a double space

	rawBox := AsciiBox{data, a, opts.BoxSet.compressBoxSet()}
	longestLine := rawBox.Width()
	footerAddOn := 0
	if footer != "" {
		footerAddOn = footerLength + 2
	}
	longestLine = max(longestLine, footerAddOn)
	if charWidth < longestLine {
		if DebugAsciiBox {
			AsciiBoxLogger.Debug().Int("nChars", longestLine-charWidth).Msg("Overflow by nChars chars")
		}
		charWidth = longestLine + a.borderWidth + a.borderWidth
	}
	var boxedString strings.Builder
	boxedString.Grow((a.borderWidth + longestLine + a.borderWidth + a.newLineCharWidth) * rawBox.Height())
	namePadding := int(math.Max(float64(charWidth-nameLength-a.borderWidth-a.extraNameCharIndent-a.borderWidth-headerLength), 0))
	boxedString.WriteString(opts.BoxSet.UpperLeftCorner + opts.BoxSet.HorizontalLine + name + strings.Repeat(opts.BoxSet.HorizontalLine, namePadding) + header + opts.BoxSet.UpperRightCorner)
	boxedString.WriteRune(a.newLine)
	// Name of the header stretches the box so we align to that
	charWidth = a.borderWidth + a.extraNameCharIndent + nameLength + namePadding + headerLength + a.borderWidth
	for _, line := range rawBox.Lines() {
		linePadding := float64(charWidth - boxLineOverheat - countChars(line))
		if linePadding < 0 {
			linePadding = 0
		}
		// TODO: this distorts boxes...
		frontPadding := math.Floor(linePadding / 2.0)
		backPadding := math.Ceil(linePadding / 2.0)
		boxedString.WriteString(opts.BoxSet.VerticalLine + strings.Repeat(a.emptyPadding, int(frontPadding)) + line + strings.Repeat(a.emptyPadding, int(backPadding)) + opts.BoxSet.VerticalLine)
		boxedString.WriteRune(a.newLine)
	}
	bottomPadding := namePadding + nameLength + a.extraNameCharIndent + headerLength - footerLength
	boxedString.WriteString(opts.BoxSet.LowerLeftCorner + strings.Repeat(opts.BoxSet.HorizontalLine, bottomPadding) + footer + opts.BoxSet.LowerRightCorner)
	return AsciiBox{boxedString.String(), a, opts.BoxSet.compressBoxSet()}
}

func (a *asciiBoxWriter) getBoxName(box AsciiBox) string {
	subMatch := a.boxHeaderRegex.FindStringSubmatch(box.String())
	if subMatch == nil {
		return ""
	}
	index := a.boxHeaderRegex.SubexpIndex("name")
	if index < 0 {
		return ""
	}
	return cleanString(subMatch[index])
}

func (a *asciiBoxWriter) getBoxHeader(box AsciiBox) string {
	subMatch := a.boxHeaderRegex.FindStringSubmatch(box.String())
	if subMatch == nil {
		return ""
	}
	index := a.boxHeaderRegex.SubexpIndex("header")
	if index < 0 {
		return ""
	}
	return cleanString(subMatch[index])
}

func (a *asciiBoxWriter) getBoxFooter(box AsciiBox) string {
	subMatch := a.boxFooterRegex.FindStringSubmatch(box.String())
	if subMatch == nil {
		return ""
	}
	index := a.boxFooterRegex.SubexpIndex("footer")
	if index < 0 {
		return ""
	}
	return cleanString(subMatch[index])
}

func (a *asciiBoxWriter) changeBoxName(box AsciiBox, newName string) AsciiBox {
	return a.changeBoxAttributes(box, &newName, nil, nil)
}

func (a *asciiBoxWriter) changeBoxHeader(box AsciiBox, newHeader string) AsciiBox {
	return a.changeBoxAttributes(box, nil, &newHeader, nil)

}

func (a *asciiBoxWriter) changeBoxFooter(box AsciiBox, newFooter string) AsciiBox {
	return a.changeBoxAttributes(box, nil, nil, &newFooter)
}

func (a *asciiBoxWriter) changeBoxAttributes(box AsciiBox, newName, newHeader, newFooter *string) AsciiBox {
	// Current data
	name := box.asciiBoxWriter.getBoxName(box)
	header := box.asciiBoxWriter.getBoxHeader(box)
	footer := box.asciiBoxWriter.getBoxFooter(box)
	// set new metadata
	if newName != nil {
		name = *newName
	}
	if newHeader != nil {
		header = *newHeader
	}
	if newFooter != nil {
		footer = *newFooter
	}
	var newOptions = []func(options *BoxOptions){
		WithAsciiBoxName(name),
		WithAsciiBoxHeader(header),
		WithAsciiBoxFooter(footer),
	}

	if !a.hasBorders(box) { // this means that this is a naked box.
		return a.boxString(box.String(), newOptions...)
	}
	minimumWidth := countChars(a.defaultBoxSet.UpperLeftCorner + a.defaultBoxSet.HorizontalLine + name + a.defaultBoxSet.UpperRightCorner)
	if header != "" { // if we have a header we need to extend that minimum width to make space for the header
		minimumWidth += countChars(a.defaultBoxSet.HorizontalLine + header)
	}
	boxContent := a.unwrap(box)                            // get the content itself ...
	rawWidth := boxContent.Width()                         // ... and look at the width.
	minimumWidth = max(minimumWidth, rawWidth+2)           // check that we have enough space for the content.
	minimumWidth = max(minimumWidth, countChars(footer)+2) // check that we have enough space for the footer.
	newBox := a.BoxString(
		boxContent.String(),
		append(newOptions, WithAsciiBoxCharWidth(minimumWidth))...,
	)
	newBox.compressedBoxSet = a.defaultBoxSet.contributeToCompressedBoxSet(box)
	return newBox
}

func (a *asciiBoxWriter) mergeHorizontal(boxes []AsciiBox) AsciiBox {
	switch len(boxes) {
	case 0:
		return AsciiBox{"", a, a.defaultBoxSet.compressBoxSet()}
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
	return AsciiBox{newBox.String(), a, a.defaultBoxSet.contributeToCompressedBoxSet(box)}
}

func (a *asciiBoxWriter) unwrap(box AsciiBox) AsciiBox {
	if !a.hasBorders(box) {
		return box
	}
	originalLines := box.Lines()
	newLines := make([]string, len(originalLines)-2)
	completeBoxSet := a.defaultBoxSet.contributeToCompressedBoxSet(box)
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
		if !strings.ContainsAny(unwrappedLine, strings.ReplaceAll(completeBoxSet, ",", "")) {
			// only trim boxes witch don't contain other boxes
			unwrappedLine = strings.Trim(unwrappedLine, a.emptyPadding)
		}
		newLines[i-1] = unwrappedLine
	}
	return AsciiBox{strings.Join(newLines, string(a.newLine)), a, completeBoxSet}
}

func (a *asciiBoxWriter) hasBorders(box AsciiBox) bool {
	if len(box.String()) == 0 {
		return false
	}
	// Check if the first char is the upper left corner
	return []rune(box.String())[0] == []rune(a.defaultBoxSet.UpperLeftCorner)[0]
}

func countChars(s string) int {
	return len([]rune(ANSI_PATTERN.ReplaceAllString(s, "")))
}

// cleanString returns the strings minus the control sequences
func cleanString(s string) string {
	regex, _ := regexp.Compile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
	return regex.ReplaceAllString(s, "")
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

func (m AsciiBox) ChangeBoxHeader(newHeader string) AsciiBox {
	return m.asciiBoxWriter.changeBoxHeader(m, newHeader)
}

func (m AsciiBox) ChangeBoxFooter(newFooter string) AsciiBox {
	return m.asciiBoxWriter.changeBoxFooter(m, newFooter)
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
func (a *asciiBoxWriter) BoxBox(box AsciiBox, options ...func(*BoxOptions)) AsciiBox {
	// TODO: if there is a box bigger then others in that this will get distorted
	newBox := a.BoxString(box.data, options...)
	newBox.compressedBoxSet = a.defaultBoxSet.contributeToCompressedBoxSet(box)
	return newBox
}

// BoxString boxes a newline separated string into a beautiful box
func (a *asciiBoxWriter) BoxString(data string, options ...func(*BoxOptions)) AsciiBox {
	return a.boxString(data, options...)
}

// AlignBoxes aligns all boxes to a desiredWidth and orders them from left to right and top to bottom (size will be at min the size of the biggest box)
func (a *asciiBoxWriter) AlignBoxes(boxes []AsciiBox, desiredWidth int, options ...func(*BoxOptions)) AsciiBox {
	if len(boxes) == 0 {
		return AsciiBox{"", a, a.defaultBoxSet.compressBoxSet()}
	}
	actualWidth := desiredWidth
	for _, box := range boxes {
		boxWidth := box.Width()
		if boxWidth > actualWidth {
			if DebugAsciiBox {
				AsciiBoxLogger.Debug().Int("nChars", boxWidth-desiredWidth).Msg("Overflow by nChars chars")
			}
			actualWidth = boxWidth
		}
	}
	if DebugAsciiBox {
		AsciiBoxLogger.Debug().Int("actualWidth", actualWidth).Msg("Working with actualWidth chars")
	}
	bigBox := AsciiBox{"", a, a.defaultBoxSet.compressBoxSet()}
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
func (a *asciiBoxWriter) BoxSideBySide(box1, box2 AsciiBox, options ...func(*BoxOptions)) AsciiBox {
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
	return AsciiBox{aggregateBox.String(), a, combineCompressedBoxSets(box1, box2)}
}

// BoxBelowBox renders two boxes below
func (a *asciiBoxWriter) BoxBelowBox(box1, box2 AsciiBox, options ...func(*BoxOptions)) AsciiBox {
	box1Width := box1.Width()
	box2Width := box2.Width()
	if box1Width < box2Width {
		box1 = a.expandBox(box1, box2Width)
	} else if box2Width < box1Width {
		box2 = a.expandBox(box2, box1Width)
	}
	return AsciiBox{box1.String() + "\n" + box2.String(), a, combineCompressedBoxSets(box1, box2)}
}
