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
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type PriorityFilter struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewPriorityFilter(arg Arg) (*PriorityFilter, error) {
	s := &PriorityFilter{
		bitNames: map[string]int{
			"manualLifeSafety":          0,
			"automaticLifeSafety":       1,
			"priority3":                 2,
			"priority4":                 3,
			"criticalEquipmentControls": 4,
			"minimumOnOff":              5,
			"priority7":                 6,
			"manualOperator":            7,
			"priority9":                 8,
			"priority10":                9,
			"priority11":                10,
			"priority12":                11,
			"priority13":                12,
			"priority14":                13,
			"priority15":                14,
			"priority16":                15,
		},
		bitLen: 16,
	}
	panic("implement me")
	return s, nil
}
