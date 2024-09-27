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

type AccessEvent struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccessEvent(arg Arg) (*AccessEvent, error) {
	s := &AccessEvent{
		vendorRange: vendorRange{512, 65535},
		enumerations: map[string]uint64{"none": 0,
			"granted":                             1,
			"muster":                              2,
			"passbackDetected":                    3,
			"duress":                              4,
			"trace":                               5,
			"lockoutMaxAttempts":                  6,
			"lockoutOther":                        7,
			"lockoutRelinquished":                 8,
			"lockedByHigherPriority":              9,
			"outOfService":                        10,
			"outOfServiceRelinquished":            11,
			"accompanimentBy":                     12,
			"authenticationFactorRead":            13,
			"authorizationDelayed":                14,
			"verificationRequired":                15,
			"deniedDenyAll":                       128,
			"deniedUnknownCredential":             129,
			"deniedAuthenticationUnavailable":     130,
			"deniedAuthenticationFactorTimeout":   131,
			"deniedIncorrectAuthenticationFactor": 132,
			"deniedZoneNoAccessRights":            133,
			"deniedPointNoAccessRights":           134,
			"deniedNoAccessRights":                135,
			"deniedOutOfTimeRange":                136,
			"deniedThreatLevel":                   137,
			"deniedPassback":                      138,
			"deniedUnexpectedLocationUsage":       139,
			"deniedMaxAttempts":                   140,
			"deniedLowerOccupancyLimit":           141,
			"deniedUpperOccupancyLimit":           142,
			"deniedAuthenticationFactorLost":      143,
			"deniedAuthenticationFactorStolen":    144,
			"deniedAuthenticationFactorDamaged":   145,
			"deniedAuthenticationFactorDestroyed": 146,
			"deniedAuthenticationFactorDisabled":  147,
			"deniedAuthenticationFactorError":     148,
			"deniedCredentialUnassigned":          149,
			"deniedCredentialNotProvisioned":      150,
			"deniedCredentialNotYetActive":        151,
			"deniedCredentialExpired":             152,
			"deniedCredentialManualDisable":       153,
			"deniedCredentialLockout":             154,
			"deniedCredentialMaxDays":             155,
			"deniedCredentialMaxUses":             156,
			"deniedCredentialInactivity":          157,
			"deniedCredentialDisabled":            158,
			"deniedNoAccompaniment":               159,
			"deniedIncorrectAccompaniment":        160,
			"deniedLockout":                       161,
			"deniedVerificationFailed":            162,
			"deniedVerificationTimeout":           163,
			"deniedOther":                         164,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
