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

type StagingObject struct {
	Object
}

func NewStagingObject(options ...Option) (*StagingObject, error) {
	s := new(StagingObject)
	objectType := "staging"
	properties := []Property{
		NewWritableProperty("presentValue", V2P(NewReal)),
		NewReadableProperty("presentStage", V2P(NewUnsigned)),
		NewReadableProperty("stages", ArrayOfP(NewStageLimitValue, 0, 0)),
		NewOptionalProperty("stageNames", ArrayOfP(NewCharacterString, 0, 0)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("units", V2P(NewEngineeringUnits)),
		NewReadableProperty("targetReferences", ArrayOfP(NewDeviceObjectReference, 0, 0)),
		NewReadableProperty("priorityForWriting", V2P(NewUnsigned)), // 1..16
		NewOptionalProperty("defaultPresentValue", V2P(NewReal)),
		NewReadableProperty("minPresValue", V2P(NewReal)),
		NewReadableProperty("maxPresValue", V2P(NewReal)),
		NewOptionalProperty("covIncrement", V2P(NewReal)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
		NewOptionalProperty("eventEnable", V2P(NewEventTransitionBits)),
		NewOptionalProperty("ackedTransitions", V2P(NewEventTransitionBits)),
		NewOptionalProperty("notifyType", V2P(NewNotifyType)),
		NewOptionalProperty("eventTimeStamps", ArrayOfP(NewTimeStamp, 3, 0)),
		NewOptionalProperty("eventMessageTexts", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventMessageTextsConfig", ArrayOfP(NewCharacterString, 3, 0)),
		NewOptionalProperty("eventDetectionEnable", V2P(NewBoolean)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		NewOptionalProperty("valueSource", V2P(NewValueSource)),
	}
	var err error
	s.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, s)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return s, nil
}
