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

type ServicesSupported struct {
	*BitString
	bitLen      int
	bitNames    map[string]int
	vendorRange vendorRange
}

func NewServicesSupported(arg Arg) (*ServicesSupported, error) {
	s := &ServicesSupported{
		bitNames: map[string]int{
			"acknowledgeAlarm":           0,
			"confirmedCOVNotification":   1,
			"confirmedEventNotification": 2,
			"getAlarmSummary":            3,
			"getEnrollmentSummary":       4,
			"subscribeCOV":               5,
			"atomicReadFile":             6,
			"atomicWriteFile":            7,
			"addListElement":             8,
			"removeListElement":          9,
			"createObject":               10,
			"deleteObject":               11,
			"readProperty":               12,
			// "readPropertyConditional": 13      // removed in version 1 revision 12,
			"readPropertyMultiple":       14,
			"writeProperty":              15,
			"writePropertyMultiple":      16,
			"deviceCommunicationControl": 17,
			"confirmedPrivateTransfer":   18,
			"confirmedTextMessage":       19,
			"reinitializeDevice":         20,
			"vtOpen":                     21,
			"vtClose":                    22,
			"vtData":                     23,
			// "authenticate": 24                 // removed in version 1 revision 11,
			// "requestKey": 25                   // removed in version 1 revision 11,
			"iAm":                                26,
			"iHave":                              27,
			"unconfirmedCOVNotification":         28,
			"unconfirmedEventNotification":       29,
			"unconfirmedPrivateTransfer":         30,
			"unconfirmedTextMessage":             31,
			"timeSynchronization":                32,
			"whoHas":                             33,
			"whoIs":                              34,
			"readRange":                          35,
			"utcTimeSynchronization":             36,
			"lifeSafetyOperation":                37,
			"subscribeCOVProperty":               38,
			"getEventInformation":                39,
			"writeGroup":                         40,
			"subscribeCOVPropertyMultiple":       41,
			"confirmedCOVNotificationMultiple":   42,
			"unconfirmedCOVNotificationMultiple": 43,
			"confirmedAuditNotification":         44,
			"auditLogQuery":                      45,
			"unconfirmedAuditNotification":       46,
			"whoAmI":                             47,
			"youAre":                             48,
		},
		bitLen: 49,
	}
	var err error
	s.BitString, err = NewBitStringWithExtension(s, NA(arg))
	if err != nil {
		return nil, errors.Wrap(err, "NewBitStringWithExtension failed")
	}
	return s, nil
}

func (s ServicesSupported) GetBitNames() map[string]int {
	return s.bitNames
}

func (s ServicesSupported) GetBitLen() int {
	return s.bitLen
}
