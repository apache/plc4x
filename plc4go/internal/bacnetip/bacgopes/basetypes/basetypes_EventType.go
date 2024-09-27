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

type EventType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewEventType(arg Arg) (*EventType, error) {
	s := &EventType{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"changeOfBitstring": 0,
			"changeOfState":  1,
			"changeOfValue":  2,
			"commandFailure": 3,
			"floatingLimit":  4,
			"outOfRange":     5,
			// -- context tag 7 is deprecated
			"changeOfLifeSafety": 8,
			"extended":           9,
			"bufferReady":        10,
			"unsignedRange":      11,
			// -- enumeration value 12 is reserved for future addenda
			"accessEvent":             13,
			"doubleOutOfRange":        14,
			"signedOutOfRange":        15,
			"unsignedOutOfRange":      16,
			"changeOfCharacterstring": 17,
			"changeOfStatusFlags":     18,
			"changeOfReliability":     19,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
