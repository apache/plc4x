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

type BinaryLightingOutputObject struct {
	Object
}

func NewBinaryLightingOutputObject(options ...Option) (*BinaryLightingOutputObject, error) {
	b := new(BinaryLightingOutputObject)
	objectType := "binaryLightingOutput"
	properties := []Property{
		NewWritableProperty("presentValue", V2P(NewBinaryLightingPV)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewOptionalProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("blinkWarnEnable", V2P(NewBoolean)),
		NewReadableProperty("egressTime", V2P(NewUnsigned)),
		NewReadableProperty("egressActive", V2P(NewBoolean)),
		NewOptionalProperty("feedbackValue", V2P(NewBinaryLightingPV)),
		NewReadableProperty("priorityArray", V2P(NewPriorityArray)),
		NewReadableProperty("relinquishDefault", V2P(NewBinaryLightingPV)),
		NewOptionalProperty("power", V2P(NewReal)),
		NewOptionalProperty("polarity", V2P(NewPolarity)),
		NewOptionalProperty("elapsedActiveTime", V2P(NewUnsigned)), // Unsigned32,
		NewOptionalProperty("timeOfActiveTimeReset", V2P(NewDateTime)),
		NewOptionalProperty("strikeCount", V2P(NewUnsigned)),
		NewOptionalProperty("timeOfStrikeCountReset", V2P(NewDateTime)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		NewReadableProperty("currentCommandPriority", V2P(NewOptionalUnsigned)),
		NewOptionalProperty("valueSource", V2P(NewValueSource)),
		NewOptionalProperty("valueSourceArray", ArrayOfP(NewValueSource, 16, 0)),
		NewOptionalProperty("lastCommandTime", V2P(NewTimeStamp)),
		NewOptionalProperty("commandTimeArray", ArrayOfP(NewTimeStamp, 16, 0)),
		NewOptionalProperty("auditablePriorityFilter", V2P(NewOptionalPriorityFilter)),
	}
	var err error
	b.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, b)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return b, nil
}
