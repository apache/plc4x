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

package object

import (
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type AuditLogObject struct {
	Object
	objectType           string // TODO: migrateme
	properties           []Property
	_object_supports_cov bool
}

func NewAuditLogObject(arg Arg) (*AuditLogObject, error) {
	o := &AuditLogObject{
		objectType: "auditLog",
		properties: []Property{
			NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
			NewReadableProperty("eventState", V2P(NewEventState)),
			NewOptionalProperty("reliability", V2P(NewReliability)),
			NewWritableProperty("enable", V2P(NewBoolean)),
			NewReadableProperty("bufferSize", V2P(NewUnsigned)), // Unsigned32
			NewReadableProperty("logBuffer", ListOfP(NewAuditLogRecord)),
			NewReadableProperty("recordCount", V2P(NewUnsigned)),      // Unsigned64
			NewReadableProperty("totalRecordCount", V2P(NewUnsigned)), // Unsigned64
			NewOptionalProperty("memberOf", V2P(NewDeviceObjectReference)),
			NewOptionalProperty("deleteOnForward", V2P(NewBoolean)),
			NewOptionalProperty("issueConfirmedNotifications", V2P(NewBoolean)),
			NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
			NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
			NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
			NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
			NewOptionalProperty("notifyType", V2P(NewNotifyType)),
			NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
			NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
			NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
			NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		},
	}
	// TODO: @register_object_type
	return o, nil
}
