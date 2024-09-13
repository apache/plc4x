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

type AccumulatorRecordAccumulatorStatus struct {
	*Enumerated
	vendorRange  vendorRange
	enumerations map[string]uint64
}

func NewAccumulatorRecordAccumulatorStatus(arg Arg) (*AccumulatorRecordAccumulatorStatus, error) {
	s := &AccumulatorRecordAccumulatorStatus{
		enumerations: map[string]uint64{"normal": 0,
			"starting":  1,
			"recovered": 2,
			"abnormal":  3,
			"failed":    4,
		},
	}
	panic("enumeratedimplementme")
	return s, nil
}
