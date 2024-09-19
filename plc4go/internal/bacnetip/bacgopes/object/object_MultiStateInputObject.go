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

type MultiStateInputObject struct {
	Object
	objectType           string // TODO: migrateme
	properties           []Property
	_object_supports_cov bool
}

func NewMultiStateInputObject(arg Arg) (*MultiStateInputObject, error) {
	o := &MultiStateInputObject{
		objectType:           "multiStateInput",
		_object_supports_cov: true,
		properties: []Property{
			NewReadableProperty("presentValue", V2P(NewUnsigned)),
			NewOptionalProperty("deviceType", V2P(NewCharacterString)),
			NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
			NewReadableProperty("eventState", V2P(NewEventState)),
			NewOptionalProperty("reliability", V2P(NewReliability)),
			NewReadableProperty("outOfService", V2P(NewBoolean)),
			NewReadableProperty("numberOfStates", V2P(NewUnsigned)),
			NewOptionalProperty("stateText", ArrayOfP(NewCharacterString, 0, 0)),
			NewOptionalProperty("timeDelay", V2P(NewUnsigned)),
			NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
			NewOptionalProperty("alarmValues", ListOfP(NewUnsigned)),
			NewOptionalProperty("faultValues", ListOfP(NewUnsigned)),
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
			NewOptionalProperty("interfaceValue", V2P(NewOptionalUnsigned)),
		},
	}
	// TODO: @register_object_type
	return o, nil
}
