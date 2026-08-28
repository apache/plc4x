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
	"regexp"
	"strconv"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// addressPattern is the pin number and the optional selection every firmata address starts with,
// ported from plc4j's FirmataTag.ADDRESS_PATTERN. Firmata has no type suffix, so the selection
// ends the address - apart from the digital form's PULLUP mode.
var addressPattern = `(?P<address>\d+)` + spiModel.ArrayGroupPattern

// TagHandler parses firmata tag addresses. There are exactly two forms, ported from plc4j's
// FirmataTagDigital.ADDRESS_PATTERN and FirmataTagAnalog.ADDRESS_PATTERN:
//
//	digital:<pin>[<quantity>]:PULLUP
//	analog:<pin>[<quantity>]
//
// where the run length and the PULLUP suffix are both optional.
type TagHandler struct {
	digitalPattern *regexp.Regexp
	analogPattern  *regexp.Regexp

	log zerolog.Logger
}

func NewTagHandler(_options ...options.WithOption) TagHandler {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return TagHandler{
		digitalPattern: regexp.MustCompile(`^digital:` + addressPattern + `(:(?P<mode>PULLUP))?$`),
		analogPattern:  regexp.MustCompile(`^analog:` + addressPattern + `$`),
		log:            customLogger,
	}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	if match := utils.GetSubgroupMatches(m.digitalPattern, tagAddress); match != nil {
		address, quantity, explicitRange, err := parseAddressAndQuantity(match, tagAddress, maxDigitalPin)
		if err != nil {
			return nil, err
		}
		var pinMode *readWriteModel.PinMode
		if match["mode"] == "PULLUP" {
			pullup := readWriteModel.PinMode_PinModePullup
			pinMode = &pullup
		}
		return NewDigitalTagWithShape(address, quantity, pinMode, explicitRange), nil
	}
	if match := utils.GetSubgroupMatches(m.analogPattern, tagAddress); match != nil {
		address, quantity, explicitRange, err := parseAddressAndQuantity(match, tagAddress, maxAnalogPin)
		if err != nil {
			return nil, err
		}
		return NewAnalogTagWithShape(address, quantity, explicitRange), nil
	}
	// "digital:2[4]" still parses, but it now selects the pin at index 4 rather than four pins;
	// a form that no longer parses at all gets the shape it should have been written in.
	return nil, spiModel.InvalidAddressError(tagAddress,
		"digital:{pin}[selection] or analog:{pin}[selection] - for example digital:2[0..3]")
}

// ParseQuery is not supported: firmata boards can be asked for their capabilities, but neither this
// driver nor plc4j's browses them.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("This driver doesn't support browsing")
}

// parseAddressAndQuantity turns the numbers out of an address into a pin and a run length, refusing
// runs which reach past the last pin the wire format can address, since a pin beyond that one is
// silently truncated to something else when it goes onto the wire.
//
// The selection's offset moves the pin, so "digital:2[4..7]" is the same run as "digital:6[0..3]".
// A firmata run is consecutive pins, so nothing deeper than a single dimension fits.
func parseAddressAndQuantity(match map[string]string, address string, maxPin uint64) (uint8, uint8, bool, error) {
	pin, err := strconv.ParseUint(match["address"], 10, 32)
	if err != nil {
		return 0, 0, false, errors.Wrapf(err, "Error parsing address %s", match["address"])
	}
	quantity := uint64(1)
	explicitRange := false
	if expression := match["array"]; expression != "" {
		dimensions, err := spiModel.ParseArrayExpression(expression, address, spiModel.SingleDimension)
		if err != nil {
			return 0, 0, false, err
		}
		pin += uint64(dimensions[0].GetLowerBound() - dimensions[0].GetBase())
		quantity = uint64(dimensions[0].GetSize())
		// A range is an array even when it spans one pin, which the count cannot say.
		explicitRange = dimensions[0].IsRange()
	}
	if quantity > maxQuantity {
		return 0, 0, false, errors.Errorf("quantity may not be larger than %d. Was %d", maxQuantity, quantity)
	}
	if pin > maxPin {
		return 0, 0, false, errors.Errorf("pin %d is out of range, the highest addressable pin is %d", pin, maxPin)
	}
	if pin+quantity-1 > maxPin {
		return 0, 0, false, errors.Errorf("a run of %d pins starting at pin %d reaches past the highest addressable pin %d", quantity, pin, maxPin)
	}
	return uint8(pin), uint8(quantity), explicitRange, nil
}
