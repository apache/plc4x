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
	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/basetypes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/primitivedata"
)

type AccessPointObject struct {
	Object
}

func NewAccessPointObject(options ...Option) (*AccessPointObject, error) {
	a := new(AccessPointObject)
	objectType := "accessPoint"
	_object_supports_cov := true
	properties := []Property{
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewReadableProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("authenticationStatus", V2P(NewAuthenticationStatus)),
		NewReadableProperty("activeAuthenticationPolicy", V2P(NewUnsigned)),
		NewReadableProperty("numberOfAuthenticationPolicies", V2P(NewUnsigned)),
		NewOptionalProperty("authenticationPolicyList", ArrayOfP(NewAuthenticationPolicy, 0, 0)),
		NewOptionalProperty("authenticationPolicyNames", ArrayOfP(NewCharacterString, 0, 0)),
		NewReadableProperty("authorizationMode", V2P(NewAuthorizationMode)),
		NewOptionalProperty("verificationTime", V2P(NewUnsigned)),
		NewOptionalProperty("lockout", V2P(NewBoolean)),
		NewOptionalProperty("lockoutRelinquishTime", V2P(NewUnsigned)),
		NewOptionalProperty("failedAttempts", V2P(NewUnsigned)),
		NewOptionalProperty("failedAttemptEvents", ListOfP(NewAccessEvent)),
		NewOptionalProperty("maxFailedAttempts", V2P(NewUnsigned)),
		NewOptionalProperty("failedAttemptsTime", V2P(NewUnsigned)),
		NewOptionalProperty("threatLevel", V2P(NewAccessThreatLevel)),
		NewOptionalProperty("occupancyUpperLimitEnforced", V2P(NewBoolean)),
		NewOptionalProperty("occupancyLowerLimitEnforced", V2P(NewBoolean)),
		NewOptionalProperty("occupancyCountAdjust", V2P(NewBoolean)),
		NewOptionalProperty("accompanimentTime", V2P(NewUnsigned)),
		NewReadableProperty("accessEvent", V2P(NewAccessEvent)),
		NewReadableProperty("accessEventTag", V2P(NewUnsigned)),
		NewReadableProperty("accessEventTime", V2P(NewTimeStamp)),
		NewReadableProperty("accessEventCredential", V2P(NewDeviceObjectReference)),
		NewOptionalProperty("accessEventAuthenticationFactor", V2P(NewAuthenticationFactor)),
		NewReadableProperty("accessDoors", ArrayOfP(NewDeviceObjectReference, 0, 0)),
		NewReadableProperty("priorityForWriting", V2P(NewUnsigned)),
		NewOptionalProperty("musterPoint", V2P(NewBoolean)),
		NewOptionalProperty("zoneTo", V2P(NewDeviceObjectReference)),
		NewOptionalProperty("zoneFrom", V2P(NewDeviceObjectReference)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("transactionNotificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("accessAlarmEvents", ListOfP(NewAccessEvent)),
		NewOptionalProperty("accessTransactionEvents", ListOfP(NewAccessEvent)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("eventAlgorithmInhibitRef", V2P(NewObjectPropertyReference)),
		NewOptionalProperty("eventAlgorithmInhibit", V2P(NewBoolean)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
	}
	var err error
	a.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties), WithObjectSupportsCov(_object_supports_cov))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, a)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return a, nil
}
