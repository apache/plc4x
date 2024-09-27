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

type Reliability struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewReliability(arg Arg) (*Reliability, error) {
	s := &Reliability{
		vendorRange: vendorRange{64, 65535},
		enumerations: map[string]uint64{"noFaultDetected": 0,
			"noSensor":                      1,
			"overRange":                     2,
			"underRange":                    3,
			"openLoop":                      4,
			"shortedLoop":                   5,
			"noOutput":                      6,
			"unreliableOther":               7,
			"processError":                  8,
			"multiStateFault":               9,
			"configurationError":            10,
			"communicationFailure":          12,
			"memberFault":                   13,
			"monitoredObjectFault":          14,
			"tripped":                       15,
			"lampFailure":                   16,
			"activationFailure":             17,
			"renewDHCPFailure":              18,
			"renewFDRegistrationFailure":    19,
			"restartAutoNegotiationFailure": 20,
			"restartFailure":                21,
			"proprietaryCommandFailure":     22,
			"faultsListed":                  23,
			"referencedObjectFault":         24,
			"multiStateOutOfRange":          25,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
