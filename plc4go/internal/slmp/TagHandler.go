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
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
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
		// The selection sits between the address and the type, as it does in plc4j's
		// SlmpTag.ADDRESS_PATTERN: "D100[0..3]:INT".
		addressPattern: regexp.MustCompile(`^(?P<device>[A-Za-z]+?)(?P<hexPrefix>0[xX])?(?P<address>[0-9A-Fa-f]+)` + spiModel.ArrayGroupPattern + `(:(?P<datatype>[A-Za-z_]+))?$`),
	}
}

func (m TagHandler) ParseTag(tagAddress string) (apiModel.PlcTag, error) {
	match := utils.GetSubgroupMatches(m.addressPattern, tagAddress)
	if match == nil {
		// "D100:INT[4]" - the count after the type - no longer parses, so name the form to
		// write rather than reporting only that nothing matched.
		return nil, spiModel.InvalidAddressError(tagAddress,
			"{device}{address}[selection]:{TYPE} - for example D100[0..3]:INT")
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

	// The selection's offset moves the device number and its size is how many devices are read,
	// so "D100[4..7]" is the same read as "D104[0..3]". An SLMP Batch Read covers one contiguous
	// run of devices, so nothing deeper than a single dimension fits.
	quantity := uint64(1)
	explicitRange := false
	if expression := match["array"]; expression != "" {
		dimensions, err := spiModel.ParseArrayExpression(expression, tagAddress, spiModel.SingleDimension)
		if err != nil {
			return nil, err
		}
		dimension := dimensions[0]
		// The offset counts elements; a device number counts 16-bit words. They coincide only for
		// a one-word type, so D100[4]:DINT would otherwise land on D104 instead of D108. This is
		// the same scale the point count below applies, so offset and length cannot disagree.
		deviceNumber += uint64(dimension.GetLowerBound()-dimension.GetBase()) * uint64(dataType.WordsPerElement())
		if deviceNumber > maxDeviceNumber {
			return nil, errors.Errorf("device number %d exceeds the 24-bit SLMP device-address range [0..%d]: %s",
				deviceNumber, maxDeviceNumber, tagAddress)
		}
		quantity = uint64(dimension.GetSize())
		// A range is an array even when it spans one element, which the count cannot say.
		explicitRange = dimension.IsRange()
	}

	// The point count is what the frame asks for, so it - not the quantity - is what has to stay
	// inside the single-frame ceiling. There is no request optimizer to split a bigger transfer.
	numberOfPoints := quantity * uint64(dataType.WordsPerElement())
	if numberOfPoints > maxPoints {
		return nil, errors.Errorf("requested %d words exceeds the single-frame Batch Read/Write ceiling of %d (no optimizer to split): %s",
			numberOfPoints, maxPoints, tagAddress)
	}

	return NewTagWithShape(device.deviceCode, uint32(deviceNumber), dataType, uint16(quantity), explicitRange), nil
}

// ParseQuery is not supported: neither this driver nor plc4j's browses an SLMP device. plc4j's
// SlmpTagHandler.parseQuery throws UnsupportedOperationException for the same reason.
func (m TagHandler) ParseQuery(_ string) (apiModel.PlcQuery, error) {
	return nil, errors.New("SLMP does not support queries")
}
