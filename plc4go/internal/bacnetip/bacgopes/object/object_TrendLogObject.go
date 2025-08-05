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

type TrendLogObject struct {
	Object
}

func NewTrendLogObject(options ...Option) (*TrendLogObject, error) {
	t := new(TrendLogObject)
	objectType := "trendLog"
	properties := []Property{
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewWritableProperty("enable", V2P(NewBoolean)),
		NewOptionalProperty("startTime", V2P(NewDateTime)),
		NewOptionalProperty("stopTime", V2P(NewDateTime)),
		NewOptionalProperty("logDeviceObjectProperty", V2P(NewDeviceObjectPropertyReference)),
		NewOptionalProperty("logInterval", V2P(NewUnsigned)),
		NewOptionalProperty("covResubscriptionInterval", V2P(NewUnsigned)),
		NewOptionalProperty("clientCovIncrement", V2P(NewClientCOV)),
		NewReadableProperty("stopWhenFull", V2P(NewBoolean)),
		NewReadableProperty("bufferSize", V2P(NewUnsigned)),
		NewReadableProperty("logBuffer", ListOfP(NewLogRecord)),
		NewWritableProperty("recordCount", V2P(NewUnsigned)),
		NewReadableProperty("totalRecordCount", V2P(NewUnsigned)),
		NewReadableProperty("loggingType", V2P(NewLoggingType)),
		NewOptionalProperty("alignIntervals", V2P(NewBoolean)),
		NewOptionalProperty("intervalOffset", V2P(NewUnsigned)),
		NewOptionalProperty("trigger", V2P(NewBoolean)),
		NewReadableProperty("statusFlags", V2P(NewStatusFlags)),
		NewOptionalProperty("reliability", V2P(NewReliability)),
		NewOptionalProperty("notificationThreshold", V2P(NewUnsigned)),
		NewOptionalProperty("recordsSinceNotification", V2P(NewUnsigned)),
		NewOptionalProperty("lastNotifyRecord", V2P(NewUnsigned)),
		NewReadableProperty("eventState", V2P(NewEventState)),
		NewOptionalProperty("notificationClass", V2P(NewUnsigned)),
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
	t.Object, err = NewObject(Combine(options, WithObjectType(objectType), WithObjectExtraProperties(properties))...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating object")
	}
	if _, err := RegisterObjectType(NKW(KWCls, t)); err != nil {
		return nil, errors.Wrap(err, "error registering object type")
	}
	return t, nil
}
