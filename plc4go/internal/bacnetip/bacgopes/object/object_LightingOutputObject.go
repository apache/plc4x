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

type LightingOutputObject struct {
	Object
}

func NewLightingOutputObject(options ...Option) (*LightingOutputObject, error) {
	l := new(LightingOutputObject)
	objectType := "lightingOutput"
	_object_supports_cov := true
	properties := []Property{
		NewWritableProperty("presentValue", V2P(NewReal)),
		NewReadableProperty("trackingValue", V2P(NewReal)),
		NewWritableProperty("lightingCommand", V2P(NewLightingCommand)),
		NewReadableProperty("inProgress", V2P(NewLightingInProgress)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewReadableProperty("outOfService", V2P(NewBoolean)),
		NewReadableProperty("blinkWarnEnable", V2P(NewBoolean)),
		NewReadableProperty("egressTime", V2P(NewUnsigned)),
		NewReadableProperty("egressActive", V2P(NewBoolean)),
		NewReadableProperty("defaultFadeTime", V2P(NewUnsigned)),
		NewReadableProperty("defaultRampRate", V2P(NewReal)),
		NewReadableProperty("defaultStepIncrement", V2P(NewReal)),
		NewOptionalProperty("transition", V2P(NewLightingTransition)),
		NewOptionalProperty("feedbackValue", V2P(NewReal)),
		NewReadableProperty("priorityArray", V2P(NewPriorityArray)),
		NewReadableProperty("relinquishDefault", V2P(NewReal)),
		NewOptionalProperty("power", V2P(NewReal)),
		NewOptionalProperty("instantaneousPower", V2P(NewReal)),
		NewOptionalProperty("minActualValue", V2P(NewReal)),
		NewOptionalProperty("maxActualValue", V2P(NewReal)),
		NewReadableProperty("lightingCommandDefaultPriority", V2P(NewUnsigned)),
		NewOptionalProperty("covIncrement", V2P(NewReal)),
		NewOptionalProperty("reliabilityEvaluationInhibit", V2P(NewBoolean)),
		NewOptionalProperty("currentCommandPriority", V2P(NewOptionalUnsigned)),
		NewOptionalProperty("valueSource", V2P(NewValueSource)),
		NewOptionalProperty("valueSourceArray", ArrayOfP(NewValueSource, 16, 0)),
		NewOptionalProperty("lastCommandTime", V2P(NewTimeStamp)),
		NewOptionalProperty("commandTimeArray", ArrayOfP(NewTimeStamp, 16, 0)),
		NewOptionalProperty("auditablePriorityFilter", V2P(NewOptionalPriorityFilter)),
	}
	var err error
	l.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties), WithObjectSupportsCov(_object_supports_cov))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, l)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return l, nil
}
