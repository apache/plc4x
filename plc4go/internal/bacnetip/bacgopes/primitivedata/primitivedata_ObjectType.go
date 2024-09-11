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

package primitivedata

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

// ObjectTypeContract provides a set of functions which can be overwritten by a sub struct
type ObjectTypeContract interface {
	EnumeratedContract
	// SetObjectType is required because we do more stuff in the constructor and can't wait for the substruct to finish that
	SetObjectType(objectType *ObjectType)
}

type ObjectType struct {
	*Enumerated

	enumerations map[string]uint64
}

func NewObjectType(args Args) (*ObjectType, error) {
	o := &ObjectType{
		enumerations: map[string]uint64{
			"accessCredential":      32,
			"accessDoor":            30,
			"accessPoint":           33,
			"accessRights":          34,
			"accessUser":            35,
			"accessZone":            36,
			"accumulator":           23,
			"alertEnrollment":       52,
			"analogInput":           0,
			"analogOutput":          1,
			"analogValue":           2,
			"auditLog":              61,
			"auditReporter":         62,
			"averaging":             18,
			"binaryInput":           3,
			"binaryLightingOutput":  55,
			"binaryOutput":          4,
			"binaryValue":           5,
			"bitstringValue":        39,
			"calendar":              6,
			"channel":               53,
			"characterstringValue":  40,
			"command":               7,
			"credentialDataInput":   37,
			"datePatternValue":      41,
			"dateValue":             42,
			"datetimePatternValue":  43,
			"datetimeValue":         44,
			"device":                8,
			"elevatorGroup":         57,
			"escalator":             58,
			"eventEnrollment":       9,
			"eventLog":              25,
			"file":                  10,
			"globalGroup":           26,
			"group":                 11,
			"integerValue":          45,
			"largeAnalogValue":      46,
			"lifeSafetyPoint":       21,
			"lifeSafetyZone":        22,
			"lift":                  59,
			"lightingOutput":        54,
			"loadControl":           28,
			"loop":                  12,
			"multiStateInput":       13,
			"multiStateOutput":      14,
			"multiStateValue":       19,
			"networkSecurity":       38,
			"networkPort":           56,
			"notificationClass":     15,
			"notificationForwarder": 51,
			"octetstringValue":      47,
			"positiveIntegerValue":  48,
			"program":               16,
			"pulseConverter":        24,
			"schedule":              17,
			"staging":               60,
			"structuredView":        29,
			"timePatternValue":      49,
			"timeValue":             50,
			"timer":                 31,
			"trendLog":              20,
			"trendLogMultiple":      27,
		},
	}
	var enumeratedContract EnumeratedContract = o
	var err error
	var arg0 any = 0
	switch len(args) {
	case 1:
		arg0 = args[0]
		switch arg0 := arg0.(type) {
		case *ObjectType:
			o.Enumerated, _ = NewEnumerated(NA(arg0.Enumerated))
			for k, v := range arg0.enumerations {
				o.enumerations[k] = v
			}
			return o, nil
		}
	case 2:
		switch arg := args[0].(type) {
		case ObjectTypeContract:
			arg.SetObjectType(o)
			enumeratedContract = arg
			argEnumerations := arg.GetEnumerations()
			for k, v := range o.enumerations {
				if _, ok := argEnumerations[k]; !ok {
					argEnumerations[k] = v
				}
			}
			o.enumerations = nil // supper seeded
		default:
			return nil, fmt.Errorf("invalid arg type: %T", arg)
		}
		arg0 = args[1]
	}
	o.Enumerated, err = NewEnumerated(NA(enumeratedContract, arg0))
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return o, nil
}

func (o *ObjectType) GetEnumerations() map[string]uint64 {
	return o.enumerations
}

func (o *ObjectType) SetEnumerated(enumerated *Enumerated) {
	o.Enumerated = enumerated
}

func (o *ObjectType) SetObjectType(_ *ObjectType) {
	panic("must be implemented by substruct")
}

func (o *ObjectType) String() string {
	value := strconv.Itoa(int(o.value))
	if o.valueString != "" {
		value = o.valueString
	}
	return fmt.Sprintf("ObjectType(%v)", value)
}
