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

type EventEnrollmentObject struct {
	Object
}

func NewEventEnrollmentObject(options ...Option) (*EventEnrollmentObject, error) {
	e := new(EventEnrollmentObject)
	objectType := "eventEnrollment"
	properties := []Property{
		NewReadableProperty("eventType", V2P(NewEventType)),
		NewReadableProperty("notifyType", V2P(NewNotifyType)),
		NewReadableProperty("eventParameters", V2P(NewEventParameter)),
		NewReadableProperty("objectPropertyReference", V2P(NewDeviceObjectPropertyReference)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewReadableProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewReadableProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewReadableProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewReadableProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("eventAlgorithmInhibitRef", V2P(NewObjectPropertyReference)),
		NewOptionalProperty("eventAlgorithmInhibit", V2P(NewBoolean)),
		NewOptionalProperty("timeDelayNormal", V2P(NewUnsigned)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("reliability", V2P(NewReliability)),
		NewOptionalProperty("faultType", V2P(NewFaultType)),
		NewOptionalProperty("faultParameters", V2P(NewFaultParameter)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
	}
	var err error
	e.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, e)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return e, nil
}
