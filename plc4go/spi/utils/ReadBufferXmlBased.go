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
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// NewXmlReadBuffer return as ReadBuffer which doesn't validate attributes and lists
func NewXmlReadBuffer(reader io.Reader) ReadBuffer {
	return &xmlReadBuffer{
		Decoder:        xml.NewDecoder(reader),
		pos:            1,
		doValidateList: false,
		doValidateAttr: false,
	}
}

// NewStrictXmlReadBuffer return as ReadBuffer which does validate attributes and lists depending on the setting
func NewStrictXmlReadBuffer(reader io.Reader, validateAttr bool, validateList bool) ReadBuffer {
	return &xmlReadBuffer{
		Decoder:        xml.NewDecoder(reader),
		pos:            1,
		doValidateAttr: validateAttr,
		doValidateList: validateList,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type xmlReadBuffer struct {
	BufferCommons
	*xml.Decoder
	pos            uint
	doValidateAttr bool
	doValidateList bool

	// tokenHistory/historyPos let GetPos()/Reset() genuinely rewind the token
	// stream (xml.Decoder.Token() is forward-only, unlike a byte buffer's
	// cursor). Speculative lookahead - e.g. terminated-array termination
	// checks that parse-then-rewind - replays already-consumed tokens from
	// tokenHistory instead of re-reading (and thereby corrupting) the
	// underlying decoder.
	//
	// posMarks queues, per reported GetPos() value, the historyPos each call
	// observed, in call order; Reset(pos) dequeues the oldest one. A plain
	// "last GetPos() wins" map would misbehave here: every generated type's
	// parse function also calls GetPos() once on entry purely for (unused)
	// diagnostics, and because bit-level fields (e.g. a 3-bit discriminator)
	// don't move x.pos across a full byte, several of those vestigial calls
	// - nested arbitrarily deep - legitimately alias to the exact same
	// reported value as the real caller's own GetPos(). Since the real
	// caller's GetPos() always runs first (its Reset() is always the
	// outermost pending one for that value) and none of the vestigial calls
	// ever call Reset() at all, taking the oldest queued entry for a value
	// is exactly the paired one, regardless of how many discarded lookups
	// share it.
	tokenHistory []xml.Token
	historyPos   int
	posMarks     map[uint32][]int
}

var _ ReadBuffer = (*xmlReadBuffer)(nil)

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (x *xmlReadBuffer) SetByteOrder(binary.ByteOrder) {
}

func (x *xmlReadBuffer) GetByteOrder() binary.ByteOrder {
	return binary.BigEndian
}

func (x *xmlReadBuffer) GetPos() uint32 {
	curPos := uint32(x.pos / 8)
	if x.posMarks == nil {
		x.posMarks = map[uint32][]int{}
	}
	x.posMarks[curPos] = append(x.posMarks[curPos], x.historyPos)
	return curPos
}

func (x *xmlReadBuffer) Reset(pos uint32) {
	x.pos = uint(pos) * 8
	if marks := x.posMarks[pos]; len(marks) > 0 {
		x.historyPos = marks[0]
		x.posMarks[pos] = marks[1:]
	}
}

// nextToken returns the next XML token, transparently replaying tokens
// already consumed from the underlying decoder (see tokenHistory) so a
// Reset() to an earlier GetPos() snapshot can genuinely rewind - the
// embedded xml.Decoder itself only ever reads forward.
func (x *xmlReadBuffer) nextToken() (xml.Token, error) {
	if x.historyPos < len(x.tokenHistory) {
		token := x.tokenHistory[x.historyPos]
		x.historyPos++
		return token, nil
	}
	token, err := x.Decoder.Token()
	if err != nil {
		return nil, err
	}
	x.tokenHistory = append(x.tokenHistory, xml.CopyToken(token))
	x.historyPos++
	return token, nil
}

func (x *xmlReadBuffer) HasMore(bitLength uint8) bool {
	// TODO: work with x.InputOffset() and check if we are at EOF
	return true
}

func (x *xmlReadBuffer) PullContext(logicalName string, readerArgs ...WithReaderArgs) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	if startElement.Name.Local != logicalName {
		return errors.Errorf("Unexpected Start element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	}
	if err := x.validateIfList(readerArgs, startElement); err != nil {
		return err
	}
	return nil
}

func (x *xmlReadBuffer) ReadBit(logicalName string, readerArgs ...WithReaderArgs) (bool, error) {
	var value bool
	err := x.decode(logicalName, rwBitKey, 1, readerArgs, &value)
	if err != nil {
		return false, err
	}
	x.move(1)
	return value, nil
}

func (x *xmlReadBuffer) ReadByte(logicalName string, readerArgs ...WithReaderArgs) (byte, error) {
	var value string
	err := x.decode(logicalName, rwByteKey, 8, readerArgs, &value)
	if err != nil {
		return 0, err
	}
	hexString := value
	if !strings.HasPrefix(hexString, "0x") {
		return 0, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
	}
	hexString = strings.Replace(hexString, "0x", "", 1)
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return 0, err
	}
	x.move(8)
	return decoded[0], nil
}

func (x *xmlReadBuffer) ReadByteArray(logicalName string, numberOfBytes int, readerArgs ...WithReaderArgs) ([]byte, error) {
	var value string
	err := x.decode(logicalName, rwByteKey, uint(numberOfBytes/8), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	hexString := value
	if !strings.HasPrefix(hexString, "0x") {
		return nil, errors.Errorf("Hex string should start with 0x. Actual value %s", hexString)
	}
	hexString = strings.Replace(hexString, "0x", "", 1)
	decoded, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	x.move(uint8(numberOfBytes / 8))
	return decoded, nil
}

func (x *xmlReadBuffer) ReadUint8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint8, error) {
	var value uint8
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint16, error) {
	var value uint16
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint32, error) {
	var value uint32
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadUint64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint64, error) {
	var value uint64
	err := x.decode(logicalName, rwUintKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int8, error) {
	var value int8
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int16, error) {
	var value int16
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int32, error) {
	var value int32
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadInt64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int64, error) {
	var value int64
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigInt(logicalName string, bitLength uint64, readerArgs ...WithReaderArgs) (*big.Int, error) {
	var value big.Int
	err := x.decode(logicalName, rwIntKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(uint8(bitLength))
	return &value, nil
}

func (x *xmlReadBuffer) ReadFloat32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float32, error) {
	var value float32
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadFloat64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float64, error) {
	var value float64
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return 0, err
	}
	x.move(bitLength)
	return value, nil
}

func (x *xmlReadBuffer) ReadBigFloat(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (*big.Float, error) {
	var value big.Float
	err := x.decode(logicalName, rwFloatKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return nil, err
	}
	x.move(bitLength)
	return &value, nil
}

func (x *xmlReadBuffer) ReadString(logicalName string, bitLength uint32, readerArgs ...WithReaderArgs) (string, error) {
	var value string
	// TODO: bitlength too short
	err := x.decode(logicalName, rwStringKey, uint(bitLength), readerArgs, &value)
	if err != nil {
		return "", err
	}
	x.move(uint8(bitLength))
	return value, nil
}

func (x *xmlReadBuffer) CloseContext(logicalName string, _ ...WithReaderArgs) error {
	endElement, err := x.travelToNextEndElement()
	if err != nil {
		return err
	}
	if endElement.Name.Local != logicalName {
		return errors.Errorf("Unexpected End element '%s'. Expected '%s'", endElement.Name.Local, logicalName)
	}
	return nil
}

func (x *xmlReadBuffer) move(bits uint8) {
	x.pos += uint(bits)
}

func (x *xmlReadBuffer) travelToNextStartElement() (xml.StartElement, error) {
	for {
		token, err := x.nextToken()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch token.(type) {
		case xml.StartElement:
			return token.(xml.StartElement), nil
		case xml.EndElement:
			return xml.StartElement{}, errors.Errorf("unexpected end element %s", token.(xml.EndElement).Name)
		}
	}
}

func (x *xmlReadBuffer) travelToNextEndElement() (xml.EndElement, error) {
	var endElement xml.EndElement
findTheEndToken:
	for {
		token, err := x.nextToken()
		if err != nil {
			return xml.EndElement{}, err
		}
		switch token.(type) {
		case xml.EndElement:
			endElement = token.(xml.EndElement)
			break findTheEndToken
		case xml.StartElement:
			return xml.EndElement{}, errors.Errorf("unexpected start element %s", token.(xml.StartElement).Name)
		}
	}
	return endElement, nil
}

func (x *xmlReadBuffer) decode(logicalName string, dataType string, bitLength uint, readerArgs []WithReaderArgs, targetValue any) error {
	startElement, err := x.travelToNextStartElement()
	if err != nil {
		return err
	}
	err = x.validateStartElement(startElement, logicalName, dataType, bitLength, readerArgs...)
	if err != nil {
		return err
	}
	// Deliberately not using x.DecodeElement here: DecodeElement consumes tokens
	// straight from the embedded xml.Decoder, bypassing nextToken()'s replay
	// buffer. That invisible consumption left GetPos()/Reset() (used by
	// terminated-array termination checks, see NoMorePathSegments) unable to
	// truly rewind, since the replayed history didn't include everything that
	// had actually been consumed. Reading the element's text content by hand -
	// through nextToken(), like everything else in this file - keeps the
	// replay buffer complete so a Reset() to an earlier GetPos() is exact.
	text, err := x.readElementText(startElement)
	if err != nil {
		return err
	}
	return assignDecodedText(text, targetValue)
}

// readElementText consumes tokens (via nextToken(), so they remain replayable)
// up to and including the EndElement matching start, returning the
// concatenated character data in between.
func (x *xmlReadBuffer) readElementText(start xml.StartElement) (string, error) {
	var sb strings.Builder
	for {
		token, err := x.nextToken()
		if err != nil {
			return "", err
		}
		switch t := token.(type) {
		case xml.CharData:
			sb.Write(t)
		case xml.EndElement:
			if t.Name != start.Name {
				return "", errors.Errorf("unexpected end element %s. Expected end of %s", t.Name, start.Name)
			}
			return sb.String(), nil
		case xml.StartElement:
			return "", errors.Errorf("unexpected nested start element %s while reading text content of %s", t.Name, start.Name)
		}
	}
}

// assignDecodedText mirrors the primitive-type unmarshalling xml.DecodeElement
// used to perform, for the handful of concrete pointer types this file ever
// passes as a decode target. Numeric/bool/big.* targets are whitespace-trimmed
// first, matching encoding/xml's own behavior (test fixtures sometimes wrap a
// value's text content across a line break for readability); *string targets
// are taken verbatim.
func assignDecodedText(text string, targetValue any) error {
	if _, isString := targetValue.(*string); !isString {
		text = strings.TrimSpace(text)
	}
	switch v := targetValue.(type) {
	case *bool:
		parsed, err := strconv.ParseBool(text)
		if err != nil {
			return errors.Wrap(err, "error parsing bool")
		}
		*v = parsed
	case *string:
		*v = text
	case *uint8:
		parsed, err := strconv.ParseUint(text, 10, 8)
		if err != nil {
			return errors.Wrap(err, "error parsing uint8")
		}
		*v = uint8(parsed)
	case *uint16:
		parsed, err := strconv.ParseUint(text, 10, 16)
		if err != nil {
			return errors.Wrap(err, "error parsing uint16")
		}
		*v = uint16(parsed)
	case *uint32:
		parsed, err := strconv.ParseUint(text, 10, 32)
		if err != nil {
			return errors.Wrap(err, "error parsing uint32")
		}
		*v = uint32(parsed)
	case *uint64:
		parsed, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return errors.Wrap(err, "error parsing uint64")
		}
		*v = parsed
	case *int8:
		parsed, err := strconv.ParseInt(text, 10, 8)
		if err != nil {
			return errors.Wrap(err, "error parsing int8")
		}
		*v = int8(parsed)
	case *int16:
		parsed, err := strconv.ParseInt(text, 10, 16)
		if err != nil {
			return errors.Wrap(err, "error parsing int16")
		}
		*v = int16(parsed)
	case *int32:
		parsed, err := strconv.ParseInt(text, 10, 32)
		if err != nil {
			return errors.Wrap(err, "error parsing int32")
		}
		*v = int32(parsed)
	case *int64:
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return errors.Wrap(err, "error parsing int64")
		}
		*v = parsed
	case *float32:
		parsed, err := strconv.ParseFloat(text, 32)
		if err != nil {
			return errors.Wrap(err, "error parsing float32")
		}
		*v = float32(parsed)
	case *float64:
		parsed, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return errors.Wrap(err, "error parsing float64")
		}
		*v = parsed
	case *big.Int:
		if err := v.UnmarshalText([]byte(text)); err != nil {
			return errors.Wrap(err, "error parsing big.Int")
		}
	case *big.Float:
		if err := v.UnmarshalText([]byte(text)); err != nil {
			return errors.Wrap(err, "error parsing big.Float")
		}
	default:
		return errors.Errorf("unsupported xml decode target type %T", targetValue)
	}
	return nil
}

func (x *xmlReadBuffer) validateIfList(readerArgs []WithReaderArgs, startElement xml.StartElement) error {
	if !x.doValidateList {
		return nil
	}
	if x.IsToBeRenderedAsList(UpcastReaderArgs(readerArgs...)...) {
		for _, attr := range startElement.Attr {
			switch attr.Name.Local {
			case rwIsListKey:
				if attr.Value != "true" {
					return errors.Errorf("Startelement should be marked as %s=true", rwIsListKey)
				}
			}
		}
	}
	return nil
}

func (x *xmlReadBuffer) validateStartElement(startElement xml.StartElement, logicalName string, dataType string, bitLength uint, readerArgs ...WithReaderArgs) error {
	logicalName = x.SanitizeLogicalName(logicalName)
	if startElement.Name.Local != logicalName {
		return errors.Errorf("unexpected start element '%s'. Expected '%s'", startElement.Name.Local, logicalName)
	} else if err := x.validateAttr(startElement.Attr, dataType, bitLength, readerArgs...); err != nil {
		return errors.Wrap(err, "Error validating Attributes")
	}
	return nil
}

func (x *xmlReadBuffer) validateAttr(attr []xml.Attr, dataType string, bitLength uint, _ ...WithReaderArgs) error {
	if !x.doValidateAttr {
		return nil
	}
	dataTypeValidated := false
	bitLengthValidate := false
	for _, attribute := range attr {
		switch attribute.Name.Local {
		case rwDataTypeKey:
			if attribute.Value != dataType {
				return errors.Errorf("Unexpected %s :%s. Want %s", rwDataTypeKey, attribute.Value, dataType)
			}
			dataTypeValidated = true
		case rwBitLengthKey:
			if attribute.Value != fmt.Sprintf("%d", bitLength) {
				return errors.Errorf("Unexpected %s '%s'. Want '%d'", rwBitLengthKey, attribute.Value, bitLength)
			}
			bitLengthValidate = true
		}
	}
	if !dataTypeValidated {
		return errors.Errorf("required attribute %s missing", rwDataTypeKey)
	}
	if !bitLengthValidate {
		return errors.Errorf("required attribute %s missing", rwBitLengthKey)
	}
	return nil
}
