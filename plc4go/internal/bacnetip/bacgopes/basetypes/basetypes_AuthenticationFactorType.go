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

type AuthenticationFactorType struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAuthenticationFactorType(arg Arg) (*AuthenticationFactorType, error) {
	s := &AuthenticationFactorType{
		enumerations: map[string]uint64{"undefined": 0,
			"error":              1,
			"custom":             2,
			"simpleNumber16":     3,
			"simpleNumber32":     4,
			"simpleNumber56":     5,
			"simpleAlphaNumeric": 6,
			"abaTrack2":          7,
			"wiegand26":          8,
			"wiegand37":          9,
			"wiegand37facility":  10,
			"facility16card32":   11,
			"facility32card32":   12,
			"fascN":              13,
			"fascNbcd":           14,
			"fascNlarge":         15,
			"fascNlargeBcd":      16,
			"gsa75":              17,
			"chuid":              18,
			"chuidFull":          19,
			"guid":               20,
			"cbeffA":             21,
			"cbeffB":             22,
			"cbeffC":             23,
			"userPassword":       24,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
