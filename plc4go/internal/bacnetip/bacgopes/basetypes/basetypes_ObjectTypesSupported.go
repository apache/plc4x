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

package basetypes

import (
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type ObjectTypesSupported struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewObjectTypesSupported(arg Arg) (*ObjectTypesSupported, error) {
	s := &ObjectTypesSupported{
		bitNames: map[string]int{
			"analogInput":           0,
			"analogOutput":          1,
			"analogValue":           2,
			"binaryInput":           3,
			"binaryOutput":          4,
			"binaryValue":           5,
			"calendar":              6,
			"command":               7,
			"device":                8,
			"eventEnrollment":       9,
			"file":                  10,
			"group":                 11,
			"loop":                  12,
			"multiStateInput":       13,
			"multiStateOutput":      14,
			"notificationClass":     15,
			"program":               16,
			"schedule":              17,
			"averaging":             18,
			"multiStateValue":       19,
			"trendLog":              20,
			"lifeSafetyPoint":       21,
			"lifeSafetyZone":        22,
			"accumulator":           23,
			"pulseConverter":        24,
			"eventLog":              25,
			"globalGroup":           26,
			"trendLogMultiple":      27,
			"loadControl":           28,
			"structuredView":        29,
			"accessDoor":            30,
			"accessCredential":      32,
			"accessPoint":           33,
			"accessRights":          34,
			"accessUser":            35,
			"accessZone":            36,
			"credentialDataInput":   37,
			"networkSecurity":       38, // removed revision 22
			"bitstringValue":        39,
			"characterstringValue":  40,
			"datePatternValue":      41,
			"dateValue":             42,
			"datetimePatternValue":  43,
			"datetimeValue":         44,
			"integerValue":          45,
			"largeAnalogValue":      46,
			"octetstringValue":      47,
			"positiveIntegerValue":  48,
			"timePatternValue":      49,
			"timeValue":             50,
			"notificationForwarder": 51,
			"alertEnrollment":       52,
			"channel":               53,
			"lightingOutput":        54,
			"binaryLightingOutput":  55,
			"networkPort":           56,
			"elevatorGroup":         57,
			"escalator":             58,
			"lift":                  59,
			"staging":               60,
			"auditLog":              61,
			"auditReporter":         62,
		},
		bitLen: 63,
	}

	var err error
	s.BitString, err = NewBitStringWithExtension(s, NA(arg))
	if err != nil {
		return nil, errors.Wrap(err, "NewBitStringWithExtension failed")
	}
	return s, nil
}

func (o ObjectTypesSupported) GetBitNames() map[string]int {
	return o.bitNames
}

func (o ObjectTypesSupported) GetBitLen() int {
	return o.bitLen
}
