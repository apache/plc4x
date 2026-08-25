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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// addressPattern is the pin number and the optional run length every firmata address starts with,
// ported verbatim from plc4j's FirmataTag.ADDRESS_PATTERN.
const addressPattern = `(?P<address>\d+)(\[(?P<quantity>\d+)])?`

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
		address, quantity, err := parseAddressAndQuantity(match, maxDigitalPin)
		if err != nil {
			return nil, err
		}
		var pinMode *readWriteModel.PinMode
		if match["mode"] == "PULLUP" {
			pullup := readWriteModel.PinMode_PinModePullup
			pinMode = &pullup
		}
		return NewDigitalTag(address, quantity, pinMode), nil
	}
	if match := utils.GetSubgroupMatches(m.analogPattern, tagAddress); match != nil {
		address, quantity, err := parseAddressAndQuantity(match, maxAnalogPin)
		if err != nil {
			return nil, err
		}
		return NewAnalogTag(address, quantity), nil
	}
	return nil, errors.Errorf("Unable to parse %s", tagAddress)
}

// ParseQuery is not supported: firmata boards can be asked for their capabilities, but neither this
// driver nor plc4j's browses them.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("This driver doesn't support browsing")
}

// parseAddressAndQuantity turns the numbers out of an address into a pin and a run length, refusing
// runs which reach past the last pin the wire format can address, since a pin beyond that one is
// silently truncated to something else when it goes onto the wire.
func parseAddressAndQuantity(match map[string]string, maxPin uint64) (uint8, uint8, error) {
	address, err := strconv.ParseUint(match["address"], 10, 32)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Error parsing address %s", match["address"])
	}
	quantity := uint64(1)
	if quantityString := match["quantity"]; quantityString != "" {
		if quantity, err = strconv.ParseUint(quantityString, 10, 32); err != nil {
			return 0, 0, errors.Wrapf(err, "Error parsing quantity %s", quantityString)
		}
	}
	if quantity == 0 {
		return 0, 0, errors.New("quantity must be greater than zero")
	}
	if quantity > maxQuantity {
		return 0, 0, errors.Errorf("quantity may not be larger than %d. Was %d", maxQuantity, quantity)
	}
	if address > maxPin {
		return 0, 0, errors.Errorf("pin %d is out of range, the highest addressable pin is %d", address, maxPin)
	}
	if address+quantity-1 > maxPin {
		return 0, 0, errors.Errorf("a run of %d pins starting at pin %d reaches past the highest addressable pin %d", quantity, address, maxPin)
	}
	return uint8(address), uint8(quantity), nil
}
