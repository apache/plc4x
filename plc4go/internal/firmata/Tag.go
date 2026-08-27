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

package firmata

import (
	"fmt"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

const (
	// maxDigitalPin is the highest digital pin the wire format can report on. The digital-IO
	// message and the subscribe-digital message both address a block of 8 pins with a 4 bit port
	// number, so pin 127 is the last one a board could ever tell us about.
	maxDigitalPin = 127
	// maxAnalogPin is the highest analog pin the wire format can carry: both the analog-IO message
	// and the subscribe-analog message address the pin with 4 bits. Higher analog pins exist on
	// real boards, but they are only reachable through the extended-analog sysex command, which
	// this driver (like plc4j's) doesn't speak.
	maxAnalogPin = 15
	// maxQuantity is the largest run of pins one tag may cover. A run has to stay inside the
	// digital pin range, so it can never be longer than that range itself.
	maxQuantity = maxDigitalPin + 1
)

// Tag is a parsed firmata tag address: one pin, or a run of consecutive pins, on the board.
//
// It is also a apiModel.PlcSubscriptionTag, because firmata's whole reason for existing in plc4x is
// its subscriptions and the subscription-request builder only accepts tags which are one.
type Tag interface {
	apiModel.PlcSubscriptionTag

	// GetAddress is the number of the first pin this tag covers.
	GetAddress() uint8
	// GetNumberOfElements is how many consecutive pins this tag covers, starting at GetAddress.
	GetNumberOfElements() uint8
}

// digitalTag addresses one or more digital pins, which are the pins that can be written and the
// ones that report a bit. Ported from plc4j's FirmataTagDigital.
type digitalTag struct {
	address  uint8
	quantity uint8
	// pinMode is the mode a subscription configures the pin as. Nil means the plc4j default of
	// INPUT; the address syntax lets the user ask for PULLUP instead.
	pinMode *readWriteModel.PinMode
}

// analogTag addresses one or more analog pins, which report a 14 bit sample. Ported from plc4j's
// FirmataTagAnalog.
type analogTag struct {
	address  uint8
	quantity uint8
}

var (
	_ Tag = digitalTag{}
	_ Tag = analogTag{}
)

// NewDigitalTag builds a digital tag. A nil pinMode means the pin is configured as a plain INPUT
// when it is subscribed to.
func NewDigitalTag(address uint8, quantity uint8, pinMode *readWriteModel.PinMode) Tag {
	return digitalTag{address: address, quantity: quantity, pinMode: pinMode}
}

// NewAnalogTag builds an analog tag.
func NewAnalogTag(address uint8, quantity uint8) Tag {
	return analogTag{address: address, quantity: quantity}
}

func (t digitalTag) GetAddress() uint8 {
	return t.address
}

func (t digitalTag) GetNumberOfElements() uint8 {
	return t.quantity
}

// GetPinMode is the mode a subscription puts the covered pins into. Nil means INPUT.
func (t digitalTag) GetPinMode() *readWriteModel.PinMode {
	return t.pinMode
}

// GetAddressString spells the tag the way the tag handler parses it back. plc4j's
// FirmataTagDigital.getAddressString drops the PULLUP suffix, which makes its address strings
// unparseable back into the same tag; since plc4go re-parses address strings whenever a tag arrives
// wrapped in a DefaultPlcSubscriptionTag, the suffix is kept here.
func (t digitalTag) GetAddressString() string {
	address := fmt.Sprintf("digital:%d%s", t.address, spiModel.RenderArrayExpression(t.GetArrayInfo()))
	if t.pinMode != nil && *t.pinMode == readWriteModel.PinMode_PinModePullup {
		address += ":PULLUP"
	}
	return address
}

func (t digitalTag) GetValueType() apiValues.PlcValueType {
	return apiValues.BOOL
}

func (t digitalTag) GetArrayInfo() []apiModel.ArrayInfo {
	return arrayInfoFor(t.quantity)
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. Firmata boards push a message whenever a subscribed pin changes, so that is a
// change-of-state subscription.
func (t digitalTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable as firmata has no per-tag cyclic subscriptions.
func (t digitalTag) GetDuration() time.Duration {
	return 0
}

func (t digitalTag) String() string {
	return fmt.Sprintf("firmata.digitalTag{address: %d, quantity: %d, pinMode: %s}",
		t.address, t.quantity, pinModeName(t.pinMode))
}

func (t analogTag) GetAddress() uint8 {
	return t.address
}

func (t analogTag) GetNumberOfElements() uint8 {
	return t.quantity
}

func (t analogTag) GetAddressString() string {
	return fmt.Sprintf("analog:%d%s", t.address, spiModel.RenderArrayExpression(t.GetArrayInfo()))
}

// GetValueType mirrors plc4j's FirmataTagAnalog.getPlcValueType. An analog sample is 14 bits wide,
// so it fits into an INT with room to spare.
func (t analogTag) GetValueType() apiValues.PlcValueType {
	return apiValues.INT
}

func (t analogTag) GetArrayInfo() []apiModel.ArrayInfo {
	return arrayInfoFor(t.quantity)
}

// GetPlcSubscriptionType is what a tag which wasn't added through one of the typed builder methods
// defaults to. Analog pins are sampled by the board and reported whenever the sample changes.
func (t analogTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionChangeOfState
}

// GetDuration is not applicable as firmata has no per-tag cyclic subscriptions.
func (t analogTag) GetDuration() time.Duration {
	return 0
}

func (t analogTag) String() string {
	return fmt.Sprintf("firmata.analogTag{address: %d, quantity: %d}", t.address, t.quantity)
}

// arrayInfoFor reports a tag covering a single pin as a scalar and everything else as an array, the
// way plc4j's FirmataTag subclasses do.
// arrayInfoFor reports the shape of the value the caller receives: a run of pins is a list, a
// single pin is a scalar. The indices are relative to the value, not to the address - a firmata
// address is a pin number, so the driver folds the start of the selection into it.
func arrayInfoFor(quantity uint8) []apiModel.ArrayInfo {
	if quantity > 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(quantity) - 1,
				Range:      true,
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

// pinModeName spells an optional pin mode for logging.
func pinModeName(pinMode *readWriteModel.PinMode) string {
	if pinMode == nil {
		return "<default>"
	}
	return pinMode.String()
}
