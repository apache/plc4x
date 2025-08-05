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

type AccumulatorObject struct {
	Object
}

func NewAccumulatorObject(options ...Option) (*AccumulatorObject, error) {
	a := new(AccumulatorObject)
	objectType := "accumulator"
	properties := []Property{
		NewReadableProperty("presentValue", V2P(NewUnsigned)),
		NewOptionalProperty("deviceType", V2P(NewCharacterString)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("scale", V2P(NewScale)),
		NewReadableProperty("units", V2P(NewEngineeringUnits)),
		NewOptionalProperty("prescale", V2P(NewPrescale)),
		NewReadableProperty("maxPresValue", V2P(NewUnsigned)),
		NewOptionalProperty("valueChangeTime", V2P(NewDateTime)),
		NewOptionalProperty("valueBeforeChange", V2P(NewUnsigned)),
		NewOptionalProperty("valueSet", V2P(NewUnsigned)),
		NewOptionalProperty("loggingRecord", V2P(NewAccumulatorRecord)),
		NewOptionalProperty("loggingObject", Vs2P(NewObjectIdentifier)),
		NewOptionalProperty("pulseRate", V2P(NewUnsigned)),
		NewOptionalProperty("highLimit", V2P(NewUnsigned)),
		NewOptionalProperty("lowLimit", V2P(NewUnsigned)),
		NewOptionalProperty("limitMonitoringInterval", V2P(NewUnsigned)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("timeDelay", V2P(NewUnsigned)),
		NewOptionalProperty("limitEnable", V2P(NewLimitEnable)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("eventAlgorithmInhibitRef", V2P(NewObjectPropertyReference)),
		NewOptionalProperty("eventAlgorithmInhibit", V2P(NewBoolean)),
		NewOptionalProperty("timeDelayNormal", V2P(NewUnsigned)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		NewOptionalProperty("faultHighLimit", V2P(NewUnsigned)),
		NewOptionalProperty("faultLowLimit", V2P(NewUnsigned)),
	}
	var err error
	a.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, a)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return a, nil
}
