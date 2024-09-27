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

type LifeSafetyState struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewLifeSafetyState(arg Arg) (*LifeSafetyState, error) {
	s := &LifeSafetyState{
		enumerations: map[string]uint64{"quiet": 0,
			"preAlarm":        1,
			"alarm":           2,
			"fault":           3,
			"faultPreAlarm":   4,
			"faultAlarm":      5,
			"notReady":        6,
			"active":          7,
			"tamper":          8,
			"testAlarm":       9,
			"testActive":      10,
			"testFault":       11,
			"testFaultAlarm":  12,
			"holdup":          13,
			"duress":          14,
			"tamperAlarm":     15,
			"abnormal":        16,
			"emergencyPower":  17,
			"delayed":         18,
			"blocked":         19,
			"localAlarm":      20,
			"generalAlarm":    21,
			"supervisory":     22,
			"testSupervisory": 23,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
