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

type Relationship struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewRelationship(arg Arg) (*Relationship, error) {
	s := &Relationship{
		vendorRange: vendorRange{1024, 65535},
		enumerations: map[string]uint64{"unknown": 0,
			"default":           1,
			"contains":          2,
			"containedBy":       3,
			"uses":              4,
			"usedBy":            5,
			"commands":          6,
			"commandedBy":       7,
			"adjusts":           8,
			"adjustedBy":        9,
			"ingress":           10,
			"egress":            11,
			"suppliesAir":       12,
			"receivesAir":       13,
			"suppliesHotAir":    14,
			"receivesHotAir":    15,
			"suppliesCoolAir":   16,
			"receivesCoolAir":   17,
			"suppliesPower":     18,
			"receivesPower":     19,
			"suppliesGas":       20,
			"receivesGas":       21,
			"suppliesWater":     22,
			"receivesWater":     23,
			"suppliesHotWater":  24,
			"receivesHotWater":  25,
			"suppliesCoolWater": 26,
			"receivesCoolWater": 27,
			"suppliesSteam":     28,
			"receivesSteam":     29,
		},
	}
	var err error
	s.Enumerated, err = NewEnumerated(NoArgs)
	if err != nil {
		return nil, errors.Wrap(err, "error creating enumerated")
	}
	return s, nil
}
