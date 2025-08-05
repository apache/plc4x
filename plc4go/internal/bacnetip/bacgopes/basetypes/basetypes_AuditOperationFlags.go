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

type AuditOperationFlags struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewAuditOperationFlags(arg Arg) (*AuditOperationFlags, error) {
	s := &AuditOperationFlags{
		vendorRange: vendorRange{32, 63},
		bitNames: map[string]int{
			"read":              0,
			"write":             1,
			"create":            2,
			"delete":            3,
			"lifeSafety":        4,
			"acknowledgeAlarm":  5,
			"deviceDisableComm": 6,
			"deviceEnableComm":  7,
			"deviceReset":       8,
			"deviceBackup":      9,
			"deviceRestore":     10,
			"subscription":      11,
			"notification":      12,
			"auditingFailure":   13,
			"networkChanges":    14,
			"general":           15,
		},
		bitLen: 16,
	}
	var err error
	s.BitString, err = NewBitStringWithExtension(s, NA(arg))
	if err != nil {
		return nil, errors.Wrap(err, "NewBitStringWithExtension failed")
	}
	return s, nil
}

func (a AuditOperationFlags) GetBitNames() map[string]int {
	return a.bitNames
}

func (a AuditOperationFlags) GetBitLen() int {
	return a.bitLen
}
