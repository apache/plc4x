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

type TimerObject struct {
	Object
}

func NewTimerObject(options ...Option) (*TimerObject, error) {
	t := new(TimerObject)
	objectType := "timer"
	properties := []Property{
		NewReadableProperty("presentValue", V2P(NewUnsigned)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewOptionalProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewOptionalProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("timerState", V2P(NewTimerState)),
		NewReadableProperty("timerRunning", V2P(NewBoolean)),
		NewOptionalProperty("updateTime", V2P(NewDateTime)),
		NewOptionalProperty("lastStateChange", V2P(NewTimerTransition)),
		NewOptionalProperty("expirationTime", V2P(NewDateTime)),
		NewOptionalProperty("initialTimeout", V2P(NewUnsigned)),
		NewOptionalProperty("defaultTimeout", V2P(NewUnsigned)),
		NewOptionalProperty("minPresValue", V2P(NewUnsigned)),
		NewOptionalProperty("maxPresValue", V2P(NewUnsigned)),
		NewOptionalProperty("resolution", V2P(NewUnsigned)),
		NewOptionalProperty("stateChangeValues", ArrayOfP(NewTimerStateChangeValue, 7, 0)),
		NewOptionalProperty("listOfObjectPropertyReferences", ListOfP(NewDeviceObjectPropertyReference)),
		NewOptionalProperty("priorityForWriting", V2P(NewUnsigned)), // 1..16
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("timeDelay", V2P(NewUnsigned)),
		NewOptionalProperty("timeDelayNormal", V2P(NewUnsigned)),
		NewOptionalProperty("alarmValues", ListOfP(NewTimerState)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventAlgorithmInhibitRef", V2P(NewObjectPropertyReference)),
		NewOptionalProperty("eventAlgorithmInhibit", V2P(NewBoolean)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
	}
	var err error
	t.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, t)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return t, nil
}
