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

package slmp

import (
	"regexp"
	"strconv"
	"strings"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// deviceAddressing says how a supported device's address is spelled.
type deviceAddressing struct {
	deviceCode readWriteModel.SlmpDeviceCode
	// base is the radix the device number is written in: MELSEC addresses W in hex and D and R in
	// decimal (SH-080008 section 8.1).
	base int
}

// supportedDevices are the devices a tag address may name. plc4j's SlmpTag supports exactly these
// three word devices and rejects everything else, including the bit devices (M, X, Y, B) and the
// timer device (TN) the mspec's device-code enum knows: reading a bit device takes a different
// point unit, and this version only speaks word units.
var supportedDevices = map[string]deviceAddressing{
	"D": {deviceCode: readWriteModel.SlmpDeviceCode_D, base: 10},
	"W": {deviceCode: readWriteModel.SlmpDeviceCode_W, base: 16},
	"R": {deviceCode: readWriteModel.SlmpDeviceCode_R, base: 10},
}

// TagHandler parses slmp tag addresses of the form
//
//	<device><address>[:<datatype>][[<quantity>]]
//
// e.g. D350, R200:REAL[4], W1A:WORD[10] or W0x1A. Ported from plc4j's SlmpTag.ADDRESS_PATTERN.
//
// The device token is matched non-greedily, unlike plc4j's pattern. plc4j matches it as a greedy
// [A-Za-z]+ ahead of a hex address, so on a W address whose first two digits are letters the greedy
// match swallows them into the device token: W1A parses, WA parses, but WAB comes out as device "WA"
// and is rejected as an unsupported device (and WABC as "WAB"). Non-greedy matching takes the
// shortest device token that lets the rest of the address parse, which is the intended one.
type TagHandler struct {
	addressPattern *regexp.Regexp
}

func NewTagHandler() TagHandler {
	return TagHandler{
		addressPattern: regexp.MustCompile(`^(?P<device>[A-Za-z]+?)(?P<hexPrefix>0[xX])?(?P<address>[0-9A-Fa-f]+)(:(?P<datatype>[A-Za-z_]+))?(\[(?P<quantity>\d+)])?$`),
	}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	match := utils.GetSubgroupMatches(m.addressPattern, tagAddress)
	if match == nil {
		return nil, errors.Errorf("Unable to parse SLMP address: %s", tagAddress)
	}

	deviceToken := strings.ToUpper(match["device"])
	device, ok := supportedDevices[deviceToken]
	if !ok {
		return nil, errors.Errorf("device '%s' not supported in this version (word devices D/W/R only): %s", deviceToken, tagAddress)
	}

	// The 0x prefix only says "read the following digits as hex", so on a device that is addressed
	// in decimal it can only be a mistake.
	if match["hexPrefix"] != "" && device.base != 16 {
		return nil, errors.Errorf("0x prefix is only valid for hex devices (W): %s", tagAddress)
	}

	// Parsed as 64 bit so that an address beyond the 24-bit field is reported as out of range
	// rather than as an unparsable one.
	deviceNumber, err := strconv.ParseUint(match["address"], device.base, 64)
	if err != nil {
		radixName := "decimal"
		if device.base == 16 {
			radixName = "hex"
		}
		return nil, errors.Wrapf(err, "Invalid %s device number in: %s", radixName, tagAddress)
	}
	if deviceNumber > maxDeviceNumber {
		return nil, errors.Errorf("device number %d exceeds the 24-bit SLMP device-address range [0..%d]: %s",
			deviceNumber, maxDeviceNumber, tagAddress)
	}

	dataType := DataTypeWORD
	if datatypeToken := match["datatype"]; datatypeToken != "" {
		if dataType, ok = DataTypeByName(strings.ToUpper(datatypeToken)); !ok {
			return nil, errors.Errorf("Unsupported SLMP data type '%s' (supported: %s)",
				datatypeToken, strings.Join(SupportedDataTypeNames, ", "))
		}
	}

	// Parsed as 32 bit, which is well past the point ceiling checked below and keeps the
	// multiplication that computes the point count from overflowing.
	quantity := uint64(1)
	if quantityToken := match["quantity"]; quantityToken != "" {
		if quantity, err = strconv.ParseUint(quantityToken, 10, 32); err != nil {
			return nil, errors.Wrapf(err, "quantity out of range in: %s", tagAddress)
		}
		if quantity < 1 {
			return nil, errors.Errorf("quantity must be >= 1 in: %s", tagAddress)
		}
	}

	// The point count is what the frame asks for, so it - not the quantity - is what has to stay
	// inside the single-frame ceiling. There is no request optimizer to split a bigger transfer.
	numberOfPoints := quantity * uint64(dataType.WordsPerElement())
	if numberOfPoints > maxPoints {
		return nil, errors.Errorf("requested %d words exceeds the single-frame Batch Read/Write ceiling of %d (no optimizer to split): %s",
			numberOfPoints, maxPoints, tagAddress)
	}

	return NewTag(device.deviceCode, uint32(deviceNumber), dataType, uint16(quantity)), nil
}

// ParseQuery is not supported: neither this driver nor plc4j's browses an SLMP device. plc4j's
// SlmpTagHandler.parseQuery throws UnsupportedOperationException for the same reason.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("SLMP does not support queries")
}
